package handlers

import (
	"MiniProjRamadh/internal/database"
	"MiniProjRamadh/internal/models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-co-op/gocron"
	"github.com/manifoldco/promptui"
	_ "github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

type WikiHandlerImpl struct {
	cfg *models.Config
}

func NewWikiHandlerImpl(cfg *models.Config) *WikiHandlerImpl {
	return &WikiHandlerImpl{cfg: cfg}
}

func (handler *WikiHandlerImpl) AddWiki() error {
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

func (handler *WikiHandlerImpl) UpdateWiki() error {
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

func (handler *WikiHandlerImpl) DeleteWiki() error {
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

func (handler *WikiHandlerImpl) GetWiki() error {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(`
    SELECT id, topic, /*description,*/ created_at, updated_at
    FROM wikis
`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id    int
			topic string
			//description string
			createdAt time.Time
			updatedAt time.Time
		)
		err := rows.Scan(&id, &topic /*&description,*/, &createdAt, &updatedAt)
		if err != nil {
			return err
		}

		fmt.Printf("ID: %d\n", id)
		fmt.Printf("Topic: %s\n", topic)
		//fmt.Printf("Description: %s\n", description)
		fmt.Printf("Created At: %s\n", createdAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated At: %s\n", updatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

func (handler *WikiHandlerImpl) StartWorker() error {
	// create a new scheduler
	scheduler := gocron.NewScheduler(time.Local)

	// add a new job that runs every minute
	scheduler.Every(1).Minute().Do(handler.UpdateDesc)

	// start the scheduler
	scheduler.StartBlocking()

	return nil
}

func (handler *WikiHandlerImpl) UpdateDesc() {
	db, err := database.ConnectDB(handler.cfg)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
        SELECT id, topic
        FROM wikis
        WHERE description IS NULL OR description = ''
    `)
	if err != nil {
		log.Printf("failed to query wikis: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id    int
			topic string
		)
		err := rows.Scan(&id, &topic)
		if err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}

		go func(id int, topic string) {
			resp, err := http.Get(fmt.Sprintf("https://en.wikipedia.org/wiki/%s", topic))
			if err != nil {
				log.Printf("failed to fetch %s: %v", topic, err)
				return
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("failed to parse HTML: %v", err)
				return
			}

			// get the first paragraph of the page
			firstParagraph := doc.Find("div#mw-content-text p").First().Text()

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
				log.Printf("failed to update wiki: %v", err)
				return
			}

			log.Printf("updated wiki %d with description: %s", id, firstParagraph)
		}(id, topic)
	}
}
