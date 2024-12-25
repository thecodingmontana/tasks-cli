/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/thecodingmontana/tasks-cli/pkg/database"
)

// Task represents a task entity
type Task struct {
	Title       string
	Description string
}

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

		task := Task{
			Title:       title,
			Description: description,
		}

		var saveErr error
		switch saveOption {
		case "CSV File":
			saveErr = saveToCSVFile(task)
		default:
			saveErr = saveToSqliteDB(db, task)
		}

		if saveErr != nil {
			fmt.Printf("%s Error: %v\n", promptui.IconBad, saveErr)
			os.Exit(1)
		}

		fmt.Printf("%s Task created successfully!\n", promptui.IconGood)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func saveToSqliteDB(db *sql.DB, task Task) error {
	createTask := `
	INSERT INTO tasks(title, description)
	VALUES(?, ?)
`
	_, err := db.Exec(createTask, task.Title, task.Description)
	if err != nil {
		return fmt.Errorf("failed to create the task: %w", err)
	}
	return nil
}

func saveToCSVFile(task Task) error {
	csvFilePath := "./pkg/database/tasks.csv"

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(csvFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open CSV File in append mode
	file, err := os.OpenFile(csvFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	nextID, err := getNextID(file)
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	newRecord := []string{
		strconv.Itoa(nextID),
		task.Title,
		task.Description,
		"pending",
		timestamp,
		timestamp,
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(newRecord); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}

func getNextID(file *os.File) (int, error) {
	// Get file info
	stat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	// If file is empty, write headers and return initial ID
	if stat.Size() == 0 {
		writer := csv.NewWriter(file)
		headers := []string{"ID", "TITLE", "DESCRIPTION", "STATUS", "CREATED AT", "UPDATED AT"}
		if err := writer.Write(headers); err != nil {
			return 0, fmt.Errorf("failed to write headers: %w", err)
		}
		writer.Flush()
		return 1, nil
	}

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return 0, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Read all records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV records: %w", err)
	}

	// If we only have headers, start with ID 1
	if len(records) <= 1 {
		return 1, nil
	}

	// Get last record (excluding header)
	lastRecord := records[len(records)-1]
	if len(lastRecord) == 0 {
		return 0, fmt.Errorf("invalid record format: empty record")
	}

	lastID, err := strconv.Atoi(lastRecord[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse last record ID: %w", err)
	}

	// Reset file pointer to end for appending
	if _, err := file.Seek(0, 2); err != nil {
		return 0, fmt.Errorf("failed to move file pointer to end: %w", err)
	}

	return lastID + 1, nil
}
