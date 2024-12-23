/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"

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
				Label: fmt.Sprintf("%s Task title", promptui.IconInitial),
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
			Label: fmt.Sprintf("%s Task description", promptui.IconInitial),
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
			fmt.Println("Coming Soon!")
		default:
			saveToSqliteDB(db, Task{
				Title:       title,
				Description: description,
			})
		}
		fmt.Printf("%s Task created successfully!", promptui.IconGood)
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
