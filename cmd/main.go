package main

import (
	"github.com/spf13/cobra"
	"log"

	"MiniProjRamadh/internal/handlers"
	"MiniProjRamadh/internal/models"
)

var cfg = &models.Config{
	DBDriver: "postgres",
	DBUser:   "postgres",
	DBPass:   "1903",
	DBName:   "CLI_Wikis",
	DBHost:   "localhost",
	DBPort:   5432,
}

var rootCmd = &cobra.Command{
	Use: "myapp",
}

var addCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.AddTopic()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var ScrapeIslandCmd = &cobra.Command{
	Use: "scrapeIsland",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.ScrapeIslandByAreaForTopics()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var AutoGenTopCmd = &cobra.Command{
	Use: "autoGenTop",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.AutoGenerateTopics()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.UpdateTopic()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.DeleteTopic()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var getCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.GetWikis()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var workerCmd = &cobra.Command{
	Use: "worker",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.StartWorker()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(ScrapeIslandCmd)
	rootCmd.AddCommand(AutoGenTopCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(workerCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
