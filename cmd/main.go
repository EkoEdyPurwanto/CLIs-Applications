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
		err := handler.AddWiki()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.UpdateWiki()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.DeleteWiki()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var getCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		handler := handlers.NewWikiHandlerImpl(cfg)
		err := handler.GetWiki()
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
