/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/manifoldco/promptui"
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

		prompt := promptui.Select{
			Label:     "Which database should we list the data from?",
			Items:     []string{"Database (sqlite)", "CSV File"},
			CursorPos: 0,
		}

		_, listChoice, promptErr := prompt.Run()

		if promptErr != nil {
			fmt.Printf("%s Error: %v\n", promptui.IconBad, promptErr)
			os.Exit(1)
		}

		format, formatErr := cmd.Flags().GetString("format")

		if formatErr != nil {
			fmt.Printf("Failed to get format flag: %v", formatErr)
		}

		switch listChoice {
		case "CSV File":
			listFromCSVFile(cmd, args)
		default:
			listFromDatabase(format)
		}
	},
}

func init() {
	listCmd.Flags().StringP("format", "f", "table", "Output format: table, json")
	rootCmd.AddCommand(listCmd)
}

func listFromDatabase(format string) {
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

	switch format {
	case "json":

	default:
		data := getRowData(rows)
		formatInTable(data)
	}
}

func listFromCSVFile(cmd *cobra.Command, args []string) {
	// Read CSV File Data

}

func formatInTable(data []DBTask) {
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

	for _, task := range data {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			task.ID,
			task.Title,
			task.Description,
			task.Status,
			task.CreatedAt,
			task.UpdatedAt,
		)
	}
	w.Flush()
}

type DBTask struct {
	ID          int
	Title       string
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

func getRowData(rows *sql.Rows) []DBTask {
	tasks := make([]DBTask, 0)
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

		task := DBTask{
			ID:          id,
			Title:       title,
			Description: description,
			Status:      status,
			CreatedAt:   timediff.TimeDiff(parsedCreatedDate),
			UpdatedAt:   timediff.TimeDiff(parsedUpdatedDate),
		}

		tasks = append(tasks, task)
	}
	return tasks
}
