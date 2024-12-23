/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
	"github.com/thecodingmontana/tasks-cli/pkg/database"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `Display all tasks with their titles, descriptions, and status`,
	Run: func(cmd *cobra.Command, args []string) {
		// DB Connect
		db := database.GetDB()

		listAllQuery := `
			SELECT * FROM tasks;
		`
		rows, listAllErr := db.Query(listAllQuery)

		if listAllErr != nil {
			log.Fatalf("Failed to fetch all tasks: %v", listAllErr)
		}
		defer rows.Close()

		w := tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', 0)

		// Headers
		fmt.Fprintln(w, "ID\tTITLE\tDESCRIPTION\tSTATUS\tCREATED AT\tUPDATED AT")

		// Separator line using dashes, adjusted to match column widths
		fmt.Fprintln(w, strings.Repeat("-", 3)+"\t"+
			strings.Repeat("-", 20)+"\t"+
			strings.Repeat("-", 30)+"\t"+
			strings.Repeat("-", 19)+"\t"+
			strings.Repeat("-", 12)+"\t"+
			strings.Repeat("-", 12))

		for rows.Next() {
			var id int
			var title, description, status, created_at, updated_at string

			err := rows.Scan(&id, &title, &description, &status, &created_at, &updated_at)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}

			parsedCreatedDate, err := time.Parse("2006-01-02T15:04:05Z", created_at)
			if err != nil {
				log.Printf("Error parsing created_at: %v", err)
				continue
			}

			parsedUpdatedDate, err := time.Parse("2006-01-02T15:04:05Z", updated_at)
			if err != nil {
				log.Printf("Error parsing updated_at: %v", err)
				continue
			}

			if len(title) > 20 {
				title = title[:17] + "..."
			}
			if len(description) > 30 {
				description = description[:27] + "..."
			}

			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
				id,
				title,
				description,
				status,
				timediff.TimeDiff(parsedCreatedDate),
				timediff.TimeDiff(parsedUpdatedDate))
		}

		w.Flush()
	},
}

func init() {
	listCmd.Flags().StringP("search", "s", "", "Search in title and description")
	listCmd.Flags().IntP("limit", "n", 0, "Limit number of results")
	listCmd.Flags().StringP("sort", "o", "", "Sort by: title, created_asc (default: created_desc)")
	listCmd.Flags().StringP("format", "f", "table", "Output format: table, json, csv")

	rootCmd.AddCommand(listCmd)
}
