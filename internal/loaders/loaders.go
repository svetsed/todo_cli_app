package loaders

import (
	"fmt"

	"github.com/svetsed/todo_cli_app/internal/models"
	"github.com/svetsed/todo_cli_app/internal/storage"
)

func LoadTodoList(filePath string) (*models.TodoList, error) {
	var todoList models.TodoList
	if err := storage.Load(filePath, &todoList); err != nil {
		return nil, fmt.Errorf("failed to load todo list from %s: %w", filePath, err)
	}

	return &todoList, nil
}

func LoadRewardSystem(filePath string) (*models.RewardSystem, error) {
	var rewardSystem models.RewardSystem
	if err := storage.Load(filePath, &rewardSystem); err != nil {
		return nil, fmt.Errorf("failed to load reward system from %s: %w", filePath, err)
	}

	return &rewardSystem, nil
}
