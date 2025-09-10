package tasks

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/svetsed/todo_cli_app/internal/config"
	"github.com/svetsed/todo_cli_app/internal/logger"
	"github.com/svetsed/todo_cli_app/internal/models"
)

// executeCommand is a helper function that simulates running a cobra command
// and captures its output for testing.
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return strings.TrimSpace(out.String()), err
}

func TestIntegration_AddCmd_SuccessfullyAddsTask(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	todoFile := filepath.Join(tempDir, "test_todo.json")

	t.Setenv("STORAGE_TODO_FILE", todoFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config for test: %v", err)
	}

	todoList := &models.TodoList{
		Tasks:  []models.Task{},
		NextID: 1,
	}

	initialData, err := json.Marshal(todoList)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := os.WriteFile(cfg.Storage.TodoFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write in file %s: %v", cfg.Storage.TodoFile, err)
	}

	addCmd := AddCmd(cfg)
	taskText := "Integration test task"

	_, err = executeCommand(addCmd, taskText)
	if err != nil {
		t.Fatalf("Command AddCmd return error: %v", err)
	}

	data, err := os.ReadFile(todoFile)
	if err != nil {
		t.Fatalf("Could not read temp. file with tasks: %v", err)
	}

	var resultList models.TodoList
	if err := json.Unmarshal(data, &resultList); err != nil {
		t.Fatalf("Could not parse JSON from file: %v", err)
	}

	if len(resultList.Tasks) != 1 {
		t.Fatalf("Expected 1 task in file, got %d", len(resultList.Tasks))
	}
	if resultList.Tasks[0].Text != taskText {
		t.Errorf("Expected text of task like this = '%s', got '%s'", taskText, resultList.Tasks[0].Text)
	}
	if resultList.Tasks[0].IsComplete {
		t.Error("New task in file should not to be comleted")
	}
}

func TestIntegration_ListCmd_SuccessfullyShowsTasks(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	todoFile := filepath.Join(tempDir, "test_todo.json")

	t.Setenv("STORAGE_TODO_FILE", todoFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config for test: %v", err)
	}

	tasks := []models.Task{
		{ID: 1, Text: "First test task", TaskPoints: 10, IsComplete: false},
		{ID: 2, Text: "Second test task", TaskPoints: 20, IsComplete: true},
	}
	todoList := &models.TodoList{
		Tasks:        tasks,
		DeletedTasks: []models.Task{},
		NextID:       3,
	}

	initialData, err := json.Marshal(todoList)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := os.WriteFile(cfg.Storage.TodoFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write in file %s: %v", cfg.Storage.TodoFile, err)
	}

	listCmd := ListCmd(cfg)
	listCmd.Flags().BoolP("points", "p", false, "Show info about points, what you can receive for the task")

	output, err := executeCommand(listCmd)
	if err != nil {
		t.Fatalf("Command ListCmd return error: %v", err)
	}

	if !strings.Contains(output, "First test task") {
		t.Errorf("Output should contain 'First test task', but it doesn't. Got: \n%s", output)
	}
	if !strings.Contains(output, "Second test task") {
		t.Errorf("Output should contain 'Second test task', but it doesn't. Got: \n%s", output)
	}
	if !strings.Contains(output, "✓") {
		t.Errorf("Output should contain the completed symbol '✓' for the second task. Got: \n%s", output)
	}
}

func TestIntegration_ListCmd_ShowsPointsWithFlag(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	todoFile := filepath.Join(tempDir, "test_todo.json")

	t.Setenv("STORAGE_TODO_FILE", todoFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config for test: %v", err)
	}

	tasks := []models.Task{
		{ID: 1, Text: "A task with points", TaskPoints: 50, IsComplete: false},
	}
	todoList := &models.TodoList{Tasks: tasks}

	initialData, err := json.Marshal(todoList)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := os.WriteFile(cfg.Storage.TodoFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write in file %s: %v", cfg.Storage.TodoFile, err)
	}

	listCmd := ListCmd(cfg)
	listCmd.Flags().BoolP("points", "p", false, "Show info about points")

	output, err := executeCommand(listCmd, "-p")
	if err != nil {
		t.Fatalf("ListCmd with --points flag finished with an unexpected error: %v", err)
	}

	if !strings.Contains(output, "50") {
		t.Errorf("Output should contain the points '50', but it doesn't. Got: \n%s", output)
	}
}

func TestIntegration_CompleteCmd_SuccessfullyCompletedTask(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	todoFile := filepath.Join(tempDir, "test_todo.json")
	rewardFile := filepath.Join(tempDir, "test_rewards.json")

	t.Setenv("STORAGE_TODO_FILE", todoFile)
	t.Setenv("STORAGE_REWARD_FILE", rewardFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config for test: %v", err)
	}

	todoList := &models.TodoList{
		Tasks: []models.Task{
			{ID: 1, Text: "Task to complete", IsComplete: false, TaskPoints: 10, IsTaskPointsReceive: false},
		},
		NextID: 2,
	}

	rewardSystem := &models.RewardSystem{Rewards: []models.Reward{}}

	initialDataTask, err := json.Marshal(todoList)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := os.WriteFile(cfg.Storage.TodoFile, initialDataTask, 0666); err != nil {
		t.Fatalf("Failed to write in file %s: %v", cfg.Storage.TodoFile, err)
	}

	initialDataReward, err := json.Marshal(rewardSystem)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := os.WriteFile(cfg.Storage.RewardFile, initialDataReward, 0666); err != nil {
		t.Fatalf("Failed to write in file %s: %v", cfg.Storage.RewardFile, err)
	}

	completeCmd := CompleteCmd(cfg)
	completeCmd.Flags().BoolP("delete", "d", false, "Delete task after completion")
	completeCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation (only with -d)")

	_, err = executeCommand(completeCmd, "1")
	if err != nil {
		t.Fatalf("CompleteCmd command finished with an unexpected error: %v", err)
	}

	todoData, err := os.ReadFile(cfg.Storage.TodoFile)
	if err != nil {
		t.Fatalf("Could not read file %s with tasks: %v", cfg.Storage.TodoFile, err)
	}

	var resultTodoList models.TodoList
	if err := json.Unmarshal(todoData, &resultTodoList); err != nil {
		t.Fatalf("Could not parse json file %s with tasks: %v", cfg.Storage.TodoFile, err)
	}

	if len(resultTodoList.Tasks) != 1 {
		t.Fatal("The number of tasks in the file should not have changed")
	}
	if !resultTodoList.Tasks[0].IsComplete {
		t.Error("Expected the task in the file to be marked as completed, but it was not")
	}

	var resultRewardSystem models.RewardSystem
	rewardData, err := os.ReadFile(cfg.Storage.RewardFile)
	if err != nil {
		t.Fatalf("Could not read file %s with tasks: %v", cfg.Storage.RewardFile, err)
	}
	if err := json.Unmarshal(rewardData, &resultRewardSystem); err != nil {
		t.Fatalf("Could not parse json file %s with tasks: %v", cfg.Storage.RewardFile, err)
	}

	if resultRewardSystem.UserPoints != resultTodoList.Tasks[0].TaskPoints {
		t.Errorf("Expected user points balance to be %d, but it is %d", resultTodoList.Tasks[0].TaskPoints, resultRewardSystem.UserPoints)
	}
}

func TestIntegration_DeleteCmd_SuccessfullyDeletesTask(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	todoFile := filepath.Join(tempDir, "test_todo_for_delete.json")
	t.Setenv("STORAGE_TODO_FILE", todoFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config for test: %v", err)
	}

	tasks := []models.Task{
		{ID: 1, Text: "Task to keep"},
		{ID: 2, Text: "Task to delete"},
	}
	todoList := &models.TodoList{
		Tasks:        tasks,
		NextID:       3,
		DeletedTasks: []models.Task{},
	}

	initialData, err := json.Marshal(todoList)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := os.WriteFile(cfg.Storage.TodoFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write in file %s: %v", cfg.Storage.TodoFile, err)
	}

	deleteCmd := DeleteCmd(cfg)
	deleteCmd.Flags().BoolP("force", "f", false, "Force delete")

	_, err = executeCommand(deleteCmd, "2", "--force")
	if err != nil {
		t.Fatalf("DeleteCmd command finished with an unexpected error: %v", err)
	}

	var resultList models.TodoList
	todoData, err := os.ReadFile(cfg.Storage.TodoFile)
	if err != nil {
		t.Fatalf("Could not read file %s with tasks: %v", cfg.Storage.TodoFile, err)
	}
	if err := json.Unmarshal(todoData, &resultList); err != nil {
		t.Fatalf("Could not parse json file %s with tasks: %v", cfg.Storage.TodoFile, err)
	}

	if len(resultList.Tasks) != 1 {
		t.Fatalf("Expected 1 task to remain in the list, but got %d", len(resultList.Tasks))
	}
	if resultList.Tasks[0].ID != 1 {
		t.Errorf("Expected the remaining task to have ID 1, but got %d", resultList.Tasks[0].ID)
	}
	if len(resultList.DeletedTasks) != 1 {
		t.Fatalf("Expected 1 task in the deleted list, but got %d", len(resultList.DeletedTasks))
	}
	if resultList.DeletedTasks[0].ID != 2 {
		t.Errorf("Expected the deleted task to have ID 2, but got %d", resultList.DeletedTasks[0].ID)
	}
}

func TestIntegration_EditCmd_SuccessfullyEditsTask(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	todoFile := filepath.Join(tempDir, "test_todo_for_edit.json")
	t.Setenv("STORAGE_TODO_FILE", todoFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config for test: %v", err)
	}

	initialTask := models.Task{ID: 1, Text: "Original text"}
	todoList := &models.TodoList{Tasks: []models.Task{initialTask}, NextID: 2}

	initialData, err := json.Marshal(todoList)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if err := os.WriteFile(cfg.Storage.TodoFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write in file %s: %v", cfg.Storage.TodoFile, err)
	}

	editCmd := EditCmd(cfg)
	newText := "This text has been updated"

	_, err = executeCommand(editCmd, "1", newText)
	if err != nil {
		t.Fatalf("EditCmd command finished with an unexpected error: %v", err)
	}

	var resultList models.TodoList
	todoData, err := os.ReadFile(cfg.Storage.TodoFile)
	if err != nil {
		t.Fatalf("Could not read file %s with tasks: %v", cfg.Storage.TodoFile, err)
	}
	if err := json.Unmarshal(todoData, &resultList); err != nil {
		t.Fatalf("Could not parse json file %s with tasks: %v", cfg.Storage.TodoFile, err)
	}

	if len(resultList.Tasks) != 1 {
		t.Fatal("The number of tasks should not have changed")
	}
	if resultList.Tasks[0].Text != newText {
		t.Errorf("Expected task text to be '%s', but got '%s'", newText, resultList.Tasks[0].Text)
	}
}
