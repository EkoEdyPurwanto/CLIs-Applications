package handlers

import (
	"MiniProjRamadh/internal/database"
	"MiniProjRamadh/internal/models"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-co-op/gocron"
	"github.com/manifoldco/promptui"
	_ "github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

type WikiHandlerImpl struct {
	cfg *models.Config
}

func NewWikiHandlerImpl(cfg *models.Config) *WikiHandlerImpl {
	return &WikiHandlerImpl{cfg: cfg}
}

func (handler *WikiHandlerImpl) AddTopic() error {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	prompt := promptui.Prompt{
		Label: "Topic",
	}
	topic, err := prompt.Run()
	if err != nil {
		return err
	}

	now := time.Now()

	stmt, err := db.Prepare(`
	       INSERT INTO wikis(topic, created_at, updated_at)
	       VALUES($1, $2, $3)
	       RETURNING id
	   `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(topic, now, now).Scan(&id)
	if err != nil {
		return err
	}

	fmt.Printf("Added topic with id %d\n", id)
	return nil
}

func (handler *WikiHandlerImpl) ScrapeIslandByAreaForTopics() error {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	// Scrape data from Wikipedia
	url := "https://en.wikipedia.org/wiki/List_of_islands_by_area"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Extract island names from the table
	var islands []string
	doc.Find("#mw-content-text > div.mw-parser-output > table > tbody > tr > td:nth-child(2) > a").Each(func(i int, s *goquery.Selection) {
		island := s.Text()
		islands = append(islands, island)
	})

	// Insert islands into database
	now := time.Now()

	stmt, err := db.Prepare(`
	       INSERT INTO wikis(topic, created_at, updated_at)
	       VALUES($1, $2, $3)
	       RETURNING id
	   `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, island := range islands {
		_, err = stmt.Exec(island, now, now)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Added %d islands\n", len(islands))
	return nil
}

func (handler *WikiHandlerImpl) AutoGenerateTopics() error {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	// Open the CSV file
	file, err := os.Open("internal/handlers/topics.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
        INSERT INTO wikis(topic, created_at, updated_at)
        VALUES($1, $2, $3)
        RETURNING id
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Loop through the records and insert each topic into the database
	now := time.Now()
	for _, record := range records {
		topic := record[0]
		var id int
		err = stmt.QueryRow(topic, now, now).Scan(&id)
		if err != nil {
			log.Printf("Failed to insert topic %s: %v", topic, err)
			continue
		}
		fmt.Printf("Added topic %s with id %d\n", topic, id)
	}

	return nil
}

func (handler *WikiHandlerImpl) UpdateTopic() error {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	prompt := promptui.Prompt{
		Label: "Topic ID",
	}
	topicID, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "New Topic Name",
	}
	newTopic, err := prompt.Run()
	if err != nil {
		return err
	}

	now := time.Now()

	stmt, err := db.Prepare(`
        UPDATE wikis
        SET topic = $1, updated_at = $2
        WHERE id = $3
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(newTopic, now, topicID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Updated %d topic(s)\n", count)
	return nil
}

func (handler *WikiHandlerImpl) DeleteTopic() error {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	prompt := promptui.Prompt{
		Label: "Topic ID",
	}
	topicID, err := prompt.Run()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`
    DELETE FROM wikis
    WHERE id = $1
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(topicID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %d topic(s)\n", count)
	return nil
}

func (handler *WikiHandlerImpl) GetWikis() error {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(`
    SELECT id, topic, description, created_at, updated_at
    FROM wikis
`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var wiki models.Wiki
		err := rows.Scan(&wiki.ID, &wiki.Topic, &wiki.Description, &wiki.CreatedAt, &wiki.UpdatedAt)
		if err != nil {
			return err
		}

		fmt.Printf("ID: %d\n", wiki.ID)
		fmt.Printf("Topic: %s\n", wiki.Topic)
		fmt.Printf("Description: %s\n", wiki.Description)
		fmt.Printf("Created At: %s\n", wiki.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated At: %s\n", wiki.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

func (handler *WikiHandlerImpl) StartWorker() error {
	// Create a new scheduler
	s := gocron.NewScheduler(time.UTC)

	// Schedule the job to run every minute
	_, err := s.Every(1).Minute().Do(handler.AutoGenerateDescWorker)
	if err != nil {
		return err
	}

	// Start the scheduler in the background
	s.StartAsync()

	// Wait for the scheduler to stop
	defer s.Stop()

	// Wait indefinitely
	select {
	case <-time.After(time.Minute):
		log.Println("DONE")
		return nil
	}

}

func (handler *WikiHandlerImpl) AutoGenerateDescWorker() error {
	// Connect to the database
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	// Query all wikis that have a null or empty description
	rows, err := db.Query(`
        SELECT id, topic
        FROM wikis
        WHERE description IS NULL OR description = ''
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create a channel to synchronize the goroutines
	ch := make(chan struct{})

	// Keep track of the number of rows returned by the query
	count := 0

	// Concurrently update each wiki's description
	for rows.Next() {
		var wiki models.Wiki
		err := rows.Scan(&wiki.ID, &wiki.Topic)
		if err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}

		count++

		go func(id int, topic string) {
			defer func() {
				// Signal the channel when the goroutine completes
				ch <- struct{}{}
			}()

			// Connect to the database
			db, err := database.ConnectDB(handler.cfg)
			if err != nil {
				log.Printf("failed to connect to database: %v", err)
				return
			}
			defer db.Close()

			// Fetch the Wikipedia page for the topic
			resp, err := http.Get(fmt.Sprintf("https://id.wikipedia.org/wiki/%s", topic))
			if err != nil {
				log.Printf("failed to fetch %s: %v", topic, err)
				return
			}
			defer resp.Body.Close()

			// Parse the HTML with goquery
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("failed to parse HTML: %v", err)
				return
			}

			// Get the first paragraph of the page
			firstParagraph := doc.Find("div#mw-content-text p").First().Text()

			// Update the wiki's description and updated_at timestamp in the database
			stmt, err := db.Prepare(`
                UPDATE wikis
                SET description = $1, updated_at = $2
                WHERE id = $3
            `)
			if err != nil {
				log.Printf("failed to prepare statement: %v", err)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(firstParagraph, time.Now(), id)
			if err != nil {
				log.Printf("failed to execute statement: %v", err)
				return
			}

			// Notification for a successful update
			log.Printf("Wiki with ID %d and topic %s has been updated.", id, topic)
		}(wiki.ID, wiki.Topic)
	}

	// Wait for all the goroutines to complete
	for i := 0; i < count; i++ {
		<-ch
	}

	return nil
}
