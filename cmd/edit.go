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
	"strings"
	"time"

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
		SET title=?, description=?, status=?, updated_at=DATETIME('now')
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
	records, file := openCSVFile()
	defer file.Close()

	data := getDataFromCSVFile(records)

	if len(data) <= 0 {
		fmt.Printf("%v No tasks found in the CSV file\n", promptui.IconGood)
		return
	}

	found := false
	taskIndex := -1

	for index, content := range data {
		if content.ID == id {
			found = true
			taskIndex = index
		}
	}

	if !found && taskIndex == -1 {
		fmt.Printf("%v No task found with ID %d to edit.\n", promptui.IconGood, id)
		return
	}

	taskToBeUpdated := data[taskIndex]
	title, description, status := otherTasksPrompt(taskToBeUpdated)
	if title != data[taskIndex].Title || description != data[taskIndex].Description || status != data[taskIndex].Status {
		updated_at := time.Now().UTC().String()
		created_at, changeErr := TimeDiffToUTC("2 hours ago")
		if changeErr != nil {
			fmt.Printf("%v Failed to change created_at to UTC Time: %v\n", promptui.IconGood, changeErr)
			return
		}
		data[taskIndex] = DBTask{
			ID:          data[taskIndex].ID,
			Title:       title,
			Description: description,
			Status:      data[taskIndex].Status,
			CreatedAt:   created_at,
			UpdatedAt:   updated_at,
		}
	}

	// Convert tasks back to CSV records
	var newRecords [][]string
	// Add headers
	newRecords = append(newRecords, []string{"ID", "TITLE", "DESCRIPTION", "STATUS", "CREATED AT", "UPDATED AT"})
	for _, task := range data {
		created_at, changeErr := TimeDiffToUTC("2 hours ago")
		if changeErr != nil {
			fmt.Printf("%v Failed to change created_at to UTC Time: %v\n", promptui.IconGood, changeErr)
			return
		}

		updated_at, updatedErr := TimeDiffToUTC("2 hours ago")
		if updatedErr != nil {
			fmt.Printf("%v Failed to change created_at to UTC Time: %v\n", promptui.IconGood, updatedErr)
			return
		}

		record := []string{
			strconv.Itoa(task.ID),
			task.Title,
			task.Description,
			task.Status,
			created_at,
			updated_at,
		}
		newRecords = append(newRecords, record)
	}

	// Truncate the file and write from beginning
	if err := file.Truncate(0); err != nil {
		fmt.Printf("%v Error truncating file: %v", promptui.IconBad, err)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("%v Error seeking file: %v", promptui.IconBad, err)
		return
	}

	// Write the updated records
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(newRecords); err != nil {
		fmt.Printf("%v Error writing to CSV: %v", promptui.IconBad, err)
		return
	}

	fmt.Printf("%v Task with ID %d edited successfully.\n", promptui.IconGood, id)
	fmt.Printf("%v CSV Data updated\n", promptui.IconGood)
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

func TimeDiffToUTC(timeDiff string) (string, error) {
	now := time.Now()

	// Parse the relative time description
	if strings.Contains(timeDiff, "just now") {
		return now.UTC().Format("2006-01-02 15:04:05.999999999 +0000 UTC"), nil
	}

	var duration time.Duration
	if strings.Contains(timeDiff, "minute") {
		minutes := 1
		fmt.Sscanf(timeDiff, "%d minutes ago", &minutes)
		duration = time.Duration(-minutes) * time.Minute
	} else if strings.Contains(timeDiff, "hour") {
		hours := 1
		fmt.Sscanf(timeDiff, "%d hours ago", &hours)
		duration = time.Duration(-hours) * time.Hour
	} else if strings.Contains(timeDiff, "day") {
		days := 1
		fmt.Sscanf(timeDiff, "%d days ago", &days)
		duration = time.Duration(-days) * 24 * time.Hour
	}

	targetTime := now.Add(duration)
	return targetTime.UTC().Format("2006-01-02 15:04:05.999999999 +0000 UTC"), nil
}
