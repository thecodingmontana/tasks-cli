/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/thecodingmontana/tasks-cli/pkg/database"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing task",
	Long:  `Edit the details of an existing task`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")

		if err != nil {
			fmt.Printf("Failed to get id flag: %v", err)
			os.Exit(1)
		}

		prompt := promptui.Select{
			Label:     "Where would you like to edit from?",
			Items:     []string{"Database (sqlite)", "CSV File"},
			CursorPos: 0,
		}

		_, choice, promptErr := prompt.Run()

		if promptErr != nil {
			fmt.Printf("%s Error: %v\n", promptui.IconBad, promptErr)
			os.Exit(1)
		}

		parsedInt, parsedErr := strconv.Atoi(id)
		if parsedErr != nil {
			fmt.Printf("%s ID needs to be an interger %v\n", promptui.IconBad, parsedErr)
			os.Exit(1)
		}

		switch choice {
		case "CSV File":
			editFromCSVFile(parsedInt)
		default:
			editFromDB(parsedInt)
		}
	},
}

func init() {
	editCmd.Flags().StringP("id", "d", "", "Edit a tast by its id")
	editCmd.MarkFlagRequired("id")

	editCmd.RegisterFlagCompletionFunc("id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.AddCommand(editCmd)
}

func editFromDB(id int) {
	db := database.GetDB()

	query := `SELECT * FROM tasks WHERE id=?`

	row := db.QueryRow(query, id)

	var task DBTask
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("%v No task found with ID %d\n", promptui.IconGood, id)
		} else {
			fmt.Printf("%v Failed to get task data for ID %d: %v\n", promptui.IconBad, id, err)
		}
		return
	}

	title, description, status := otherTasksPrompt(task)

	updateQuery := `
		UPDATE tasks
		SET title=?, description=?, status=?
		WHERE id=?
	`
	if _, updateErr := db.Exec(updateQuery, title, description, status, id); updateErr != nil {
		fmt.Printf("%v Failed to create the task: %v\n", promptui.IconBad, updateErr)
		return
	}

	fmt.Printf("%v Successfully updated task with ID %d\n", promptui.IconGood, id)
	fmt.Printf("%v Database (sqlite) Data updated\n", promptui.IconGood)
}

func editFromCSVFile(id int) {
	fmt.Println(id)
}

func otherTasksPrompt(task DBTask) (title, description, status string) {
	titlePrompt := promptui.Prompt{
		Label:   "Task title: ",
		Default: task.Title,
	}

	title, titleErr := titlePrompt.Run()
	if titleErr != nil {
		fmt.Printf("%s Error: %v\n", promptui.IconBad, titleErr)
		os.Exit(1)
	}

	descriptionPrompt := promptui.Prompt{
		Label:   "Task description: ",
		Default: task.Description,
	}

	description, descriptionErr := descriptionPrompt.Run()
	if titleErr != nil {
		fmt.Printf("%s Error: %v\n", promptui.IconBad, descriptionErr)
		os.Exit(1)
	}

	statusPrompt := promptui.Prompt{
		Label:   "Task status (pending/completed): ",
		Default: task.Status,
	}

	status, statusErr := statusPrompt.Run()
	if titleErr != nil {
		fmt.Printf("%s Error: %v\n", promptui.IconBad, statusErr)
		os.Exit(1)
	}
	return title, description, status
}
