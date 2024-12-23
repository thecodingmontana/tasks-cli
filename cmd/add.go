/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/thecodingmontana/tasks-cli/pkg/database"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task.",
	Long:  `Add a new task by providing a title, description, and optional status.`,
	Run: func(cmd *cobra.Command, args []string) {
		title, titleErr := cmd.Flags().GetString("title")
		if titleErr != nil {
			log.Fatalf("Error getting title flag: %v", titleErr)
		}

		description, descErr := cmd.Flags().GetString("description")
		if descErr != nil {
			log.Fatalf("Error getting description flag: %v", descErr)
		}

		// Insert task to db
		db := database.GetDB()
		query := `
			INSERT INTO tasks(title, description)
			VALUES(?, ?);
		`

		if _, err := db.Exec(query, title, description); err != nil {
			log.Fatalf("Failed to fetch %v", err)
			return
		}
		fmt.Println("Task saved successfully!")
	},
}

func init() {
	// Title flag
	addCmd.Flags().StringP("title", "t", "", "Task title.")
	addCmd.MarkFlagRequired("title")

	// Description flag
	addCmd.Flags().StringP("description", "d", "", "Task description.")

	rootCmd.AddCommand(addCmd)
}
