/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/thecodingmontana/tasks-cli/pkg/database"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Long:  `Create a new task by providing a title, optional description, and optional status.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var title, description string
		db := database.GetDB()

		if len(args) > 0 {
			title = args[0]
			fmt.Printf("%s Task title: %s\n", promptui.IconGood, title)
		} else {
			prompt := promptui.Prompt{
				Label: fmt.Sprintf("%s Task title: ", promptui.IconInitial),
				Validate: func(input string) error {
					if len(input) == 0 {
						return fmt.Errorf("project title cannot be empty")
					}
					return nil
				},
			}

			// Run the prompt
			projectTitle, err := prompt.Run()
			if err != nil {
				fmt.Printf("%s Error: %v\n", promptui.IconBad, err)
				os.Exit(1)
			}
			title = projectTitle
		}

		descriptionPrompt := promptui.Prompt{
			Label: fmt.Sprintf("%s Task description (optional): ", promptui.IconInitial),
		}

		descriptionText, err := descriptionPrompt.Run()
		if err != nil {
			fmt.Printf("%s Error: %v\n", promptui.IconBad, err)
			os.Exit(1)
		}

		description = descriptionText

		saveOptionPrompt := promptui.Select{
			Label:     "Where do you wish to save the task",
			Items:     []string{"Database (sqlite)", "CSV File"},
			CursorPos: 0,
		}

		_, saveOption, saveOptionErr := saveOptionPrompt.Run()

		if saveOptionErr != nil {
			fmt.Printf("%s Error: %v\n", promptui.IconBad, saveOptionErr)
			os.Exit(1)
		}

		switch saveOption {
		case "CSV File":
			saveToCSVFile(Task{
				Title:       title,
				Description: description,
			})
		default:
			saveToSqliteDB(db, Task{
				Title:       title,
				Description: description,
			})
		}
		fmt.Printf("%s Task created successfully!\n", promptui.IconGood)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

type Task struct {
	Title       string
	Description string
}

func saveToSqliteDB(db *sql.DB, task Task) {
	createTask := `
	INSERT INTO tasks(title, description)
	VALUES(?, ?)
`
	if _, taskCreateErr := db.Exec(createTask, task.Title, task.Description); taskCreateErr != nil {
		fmt.Printf("Failed to create the task: %v", taskCreateErr)
		os.Exit(1)
	}
}

func saveToCSVFile(task Task) {
	// Open CSV File in append mode create if doesn't exists
	file, err := os.OpenFile("./pkg/database/tasks.csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		fmt.Printf("Failed to create/open CSV file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	// check if CSV file is empty to write headers
	stat, statErr := file.Stat()

	if statErr != nil {
		fmt.Printf("Cannot get file info: %v", statErr)
	}

	nextId := 1
	if stat.Size() <= 0 {
		// File empty write headers.
		writer := csv.NewWriter(file)
		csvHeaders := []string{"ID", "TITLE", "DESCRIPTION", "STATUS", "CREATED AT", "UPDATED AT"}

		if errCSVHeaders := writer.Write(csvHeaders); errCSVHeaders != nil {
			fmt.Printf("Failed to write headers for the CSV File: %v", errCSVHeaders)
			os.Exit(1)
		}
		writer.Flush()
	} else {
		// Read CSV File
		reader := csv.NewReader(file)
		records, errRecords := reader.ReadAll()

		if errRecords != nil {
			fmt.Printf("Error loading CSV File records: %v", errRecords)
			os.Exit(1)
		}

		if len(records) > 0 {
			lastRecord := records[len(records)-1]
			lastIndex, errLastIndex := strconv.Atoi(lastRecord[0])

			if errLastIndex != nil {
				fmt.Printf("Failed to parse the last record index: %v", errLastIndex)
				os.Exit(1)
			}
			nextId = lastIndex + 1
		}
	}

	newRecord := append([]string{strconv.Itoa(nextId)}, task.Title, task.Description, "pending", time.Now().UTC().String(), time.Now().UTC().String())

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(newRecord)
}
