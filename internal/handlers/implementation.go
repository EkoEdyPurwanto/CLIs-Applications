package handlers

import (
	"MiniProjRamadh/internal/models"
	"database/sql"
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
	DB *sql.DB
}

func NewWikiHandlerImpl(DB *sql.DB) *WikiHandlerImpl {
	return &WikiHandlerImpl{DB: DB}
}

func (handler *WikiHandlerImpl) AddTopic() error {
	// Meminta USER untk memasukkan topik
	prompt := promptui.Prompt{
		Label: "Topic",
	}
	topic, err := prompt.Run()
	if err != nil {
		return err
	}

	// untk Mendapatkan waktu saat ini
	now := time.Now()

	// Menyiapkan statement SQL untuk memasukkan data baru ke dalam tabel wikis
	stmt, err := handler.DB.Prepare(`
	       INSERT INTO wikis(topic, created_at, updated_at)
	       VALUES($1, $2, $3)
	       RETURNING id
	   `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Mengeksekusi statement SQL dan mengambil ID dari data baru
	var id int
	err = stmt.QueryRow(topic, now, now).Scan(&id)
	if err != nil {
		return err
	}

	// Mencetak pesan yang menunjukkan bahwa topik telah berhasil ditambahkan
	fmt.Printf("Added topic with id %d\n", id)
	return nil
}

func (handler *WikiHandlerImpl) ScrapeIslandByAreaForTopics() error {

	// Scrape data dr Wikipedia
	url := "https://en.wikipedia.org/wiki/List_of_islands_by_area"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parsing HTML dari response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Mengekstrak nama pulau dari tabel di halaman Wikipedia
	var islands []string
	doc.Find("#mw-content-text > div.mw-parser-output > table > tbody > tr > td:nth-child(2) > a").Each(func(i int, s *goquery.Selection) {
		island := s.Text()
		islands = append(islands, island)
	})

	// untk Mendapatkan waktu saat ini
	now := time.Now()

	// Menyiapkan statement SQL untuk memasukkan data baru ke dalam tabel wikis
	stmt, err := handler.DB.Prepare(`
	       INSERT INTO wikis(topic, created_at, updated_at)
	       VALUES($1, $2, $3)
	       RETURNING id
	   `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Melakukan loop untuk memasukkan setiap data pulau ke dalam tabel wikis
	for _, island := range islands {
		_, err = stmt.Exec(island, now, now)
		if err != nil {
			return err
		}
	}

	// Mencetak pesan yang menunjukkan bahwa topik islands telah berhasil ditambahkan
	fmt.Printf("Added %d islands\n", len(islands))
	return nil
}

func (handler *WikiHandlerImpl) AutoGenerateTopics() error {

	// Membuka file CSV yang berisi topik-topik yang akan digenerate
	file, err := os.Open("internal/handlers/topics.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Membaca isi file CSV menggunakan package csv
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Menyiapkan statement SQL untuk memasukkan data baru ke dalam tabel wikis
	stmt, err := handler.DB.Prepare(`
       INSERT INTO wikis(topic, created_at, updated_at)
       VALUES($1, $2, $3)
       RETURNING id
   `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// untk Mendapatkan waktu saat ini
	now := time.Now()

	// Melakukan generate topik berdasarkan isi file CSV
	for _, record := range records {
		topic := record[0]
		var id int
		err = stmt.QueryRow(topic, now, now).Scan(&id)
		if err != nil {
			log.Printf("Failed to insert topic %s: %v", topic, err)
			continue
		}
		// Mencetak pesan yang menunjukkan bahwa topik telah berhasil ditambahkan
		fmt.Printf("Added topic %s with id %d\n", topic, id)
	}

	return nil
}

func (handler *WikiHandlerImpl) UpdateTopic() error {

	// Meminta USER untk memasukkan ID topik yg akan di ubh
	prompt := promptui.Prompt{
		Label: "Topic ID",
	}
	topicID, err := prompt.Run()
	if err != nil {
		return err
	}

	// Meminta USER untk memasukkan topik baru
	prompt = promptui.Prompt{
		Label: "New Topic Name",
	}
	newTopic, err := prompt.Run()
	if err != nil {
		return err
	}

	// untk Mendapatkan waktu saat ini
	now := time.Now()

	// Menyiapkan statement SQL untuk mengupdate data yang sudah di tentukan sebelumnya ke dalam tabel wikis
	stmt, err := handler.DB.Prepare(`
       UPDATE wikis
       SET topic = $1, updated_at = $2
       WHERE id = $3
   `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Eksekusi statement SQL untuk mengupdate topik di tabel wikis
	res, err := stmt.Exec(newTopic, now, topicID)
	if err != nil {
		return err
	}

	// Mengambil jumlah baris yang terupdate
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Mencetak pesan yang menunjukkan bahwa topik telah berhasil di update
	fmt.Printf("Updated %d topic(s)\n", count)
	return nil
}

func (handler *WikiHandlerImpl) DeleteTopic() error {

	// Meminta USER untk memasukkan ID topik yg akan di hapus
	prompt := promptui.Prompt{
		Label: "Topic ID",
	}
	topicID, err := prompt.Run()
	if err != nil {
		return err
	}

	// Menyiapkan statement SQL untuk mendelete data yang sudah di tentukan sebelumnya yang ada di dalam tabel wikis
	stmt, err := handler.DB.Prepare(`
   DELETE FROM wikis
   WHERE id = $1
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Eksekusi statement SQL untuk mendelete topik di tabel wikis
	res, err := stmt.Exec(topicID)
	if err != nil {
		return err
	}

	// Mengambil jumlah baris yang terdelete
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Mencetak pesan yang menunjukkan bahwa topik telah berhasil di delete
	fmt.Printf("Deleted %d topic(s)\n", count)
	return nil
}

func (handler *WikiHandlerImpl) GetWikis() error {

	// Mengeksekusi query SQL untuk mengambil semua baris dari tabel wikis
	rows, err := handler.DB.Query(`
   SELECT id, topic, description, created_at, updated_at
   FROM wikis
`)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Mengambil setiap baris hasil query SQL dan memasukkannya ke dalam objek models.Wiki
	for rows.Next() {
		var wiki models.Wiki
		err := rows.Scan(&wiki.ID, &wiki.Topic, &wiki.Description, &wiki.CreatedAt, &wiki.UpdatedAt)
		if err != nil {
			return err
		}

		// Mencetak informasi dari setiap objek models.Wiki
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
	// Membuat sebuah scheduler baru
	s := gocron.NewScheduler(time.UTC)

	// Menjadwalkan pekerjaan untuk dijalankan setiap menit
	_, err := s.Every(1).Minute().Do(handler.AutoGenerateDescWorker)
	if err != nil {
		return err
	}

	// Memulai scheduler di background
	s.StartAsync()

	// Menunggu scheduler berhenti
	defer s.Stop()

	// Mengeksekusi pekerjaan setiap menit dan mencetak pesan "DONE" jika semua desc sudah di isi
	select {
	case <-time.After(time.Minute):
		log.Println("DONE")
		return nil
	}

}

func (handler *WikiHandlerImpl) AutoGenerateDescWorker() error {

	// Query semua wikis yang memiliki deskripsi null atau kosong
	rows, err := handler.DB.Query(`
       SELECT id, topic
       FROM wikis
       WHERE description IS NULL OR description = ''
   `)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Buat channel untuk sinkronisasi goroutine
	ch := make(chan struct{})

	// Hitung jumlah baris yang dikembalikan oleh query
	count := 0

	// Loop melalui semua baris yang dikembalikan oleh query
	for rows.Next() {
		var wiki models.Wiki
		err := rows.Scan(&wiki.ID, &wiki.Topic)
		if err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}

		count++

		// Goroutine untuk memperbarui deskripsi wiki
		go func(id int, topic string) {
			defer func() {
				// Signal channel ketika goroutine selesai
				ch <- struct{}{}
			}()

			//// Connect ke DB
			//db, err := database.ConnectDB(handler.cfg)
			//if err != nil {
			//	log.Printf("failed to connect to database: %v", err)
			//	return
			//}
			//defer db.Close()

			// Ambil halaman Wikipedia untuk topik
			resp, err := http.Get(fmt.Sprintf("https://id.wikipedia.org/wiki/%s", topic))
			if err != nil {
				log.Printf("failed to fetch %s: %v", topic, err)
				return
			}
			defer resp.Body.Close()

			// Parse HTML dengan goquery
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("failed to parse HTML: %v", err)
				return
			}

			// Dapatkan paragraf pertama halaman
			firstParagraph := doc.Find("div#mw-content-text p").First().Text()

			// Perbarui deskripsi dan timestamp updated_at di database
			stmt, err := handler.DB.Prepare(`
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

	// Tunggu semua goroutine selesai
	for i := 0; i < count; i++ {
		<-ch
	}

	return nil
}
