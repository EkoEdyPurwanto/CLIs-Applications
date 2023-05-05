package handlers

import (
	"MiniProjRamadh/internal/database"
	"MiniProjRamadh/internal/models"
	"fmt"
	"github.com/manifoldco/promptui"
	_ "github.com/spf13/cobra"
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
