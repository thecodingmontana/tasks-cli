/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
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

		status, statusErr := cmd.Flags().GetString("status")
		if statusErr != nil {
			log.Fatalf("Error getting status flag: %v", statusErr)
		}

		fmt.Println(title, description, status)
	},
}

func init() {
	// Title flag
	addCmd.Flags().StringP("title", "t", "", "Task title.")
	addCmd.MarkFlagRequired("title")

	// Description flag
	addCmd.Flags().StringP("description", "d", "", "Task description.")

	// Status flag
	addCmd.Flags().StringP("status", "s", "", "Task status (pending/completed).")

	rootCmd.AddCommand(addCmd)
}
