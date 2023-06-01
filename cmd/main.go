package main

import (
	"MiniProjRamadh/internal/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
	"log"

	"MiniProjRamadh/internal/handlers"
	"MiniProjRamadh/internal/models"
)

var cfg = &models.Config{
	DBDriver: "postgres",
	DBUser:   "postgres",
	DBPass:   "postgres",
	DBName:   "cli_wikis",
	DBHost:   "localhost",
	DBPort:   5432,
}

var rootCmd = &cobra.Command{
	Use: "myapp",
}

var handler *handlers.WikiHandlerImpl

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(ScrapeIslandCmd)
	rootCmd.AddCommand(AutoGenTopCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(workerCmd)

	db, err := database.GetDBInstance(cfg)
	if err != nil {
		log.Fatal(err)
	}

	handler = handlers.NewWikiHandlerImpl(db)
}

var addCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.AddTopic()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var ScrapeIslandCmd = &cobra.Command{
	Use: "scrapeIsland",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.ScrapeIslandByAreaForTopics()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var AutoGenTopCmd = &cobra.Command{
	Use: "autoGenTop",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.AutoGenerateTopics()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.UpdateTopic()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.DeleteTopic()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var getCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.GetWikis()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var workerCmd = &cobra.Command{
	Use: "worker",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.StartWorker()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	var cfg models.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed when parsing config: %v", err)
	}

	connectDB, err := database.ConnectDB(&cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.Migrate(connectDB)
	if err != nil {
		log.Fatalf("failed to run database migration: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
