package handlers

import (
	"testing"

	"github.com/svetsed/todo_cli_app/internal/models"
)

func TestTaskHandler_Add(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks:  []models.Task{},
			NextID: 1,
		},
	}

	taskText := "New Test Task"
	taskPoints := 10

	handler.Add(taskText, taskPoints)

	if len(handler.Todo.Tasks) != 1 {
		t.Fatalf("Expected 1 task in the list, got %d", len(handler.Todo.Tasks))
	}

	addedTask := handler.Todo.Tasks[0]
	if addedTask.Text != taskText {
		t.Errorf("Ecpected test of task '%s', got '%s'", taskText, addedTask.Text)
	}
	if addedTask.TaskPoints != taskPoints {
		t.Errorf("Expected %d TaskPoints, а got %d", taskPoints, addedTask.TaskPoints)
	}
	if addedTask.IsComplete {
		t.Error("New task should not be completed")
	}
	if addedTask.ID != 1 {
		t.Errorf("Expected ID=1, got %d", addedTask.ID)
	}
	if handler.Todo.NextID != 2 {
		t.Errorf("Expected, what NextID becoming 2, but he equal %d", handler.Todo.NextID)
	}
}

func TestTaskHandler_Complete_SuccessfullyCompletedTask(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks: []models.Task{
				{ID: 1, Text: "Test Task", IsComplete: false},
			},
		},
	}

	err := handler.Complete(0)

	if err != nil {
		t.Fatalf("Complete() returned an unexpected error: %v", err)
	}
	if !handler.Todo.Tasks[0].IsComplete {
		t.Error("Expected IsComplete to be true, but it was false")
	}
}

func TestTaskHandler_Complete_AlreadyCompletedTask(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks: []models.Task{
				{ID: 1, Text: "Test Task", IsComplete: true},
			},
		},
	}

	err := handler.Complete(0)

	if err == nil {
		t.Fatal("Expected an error when completing an already completed task, but got nil")
	}
}

func TestTaskHandler_NotCompleted_SuccessfullyUncompletedTask(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks: []models.Task{
				{ID: 1, IsComplete: true},
			},
		},
	}

	err := handler.NotCompleted(0)

	if err != nil {
		t.Fatalf("NotCompleted() returned an unexpected error: %v", err)
	}
	if handler.Todo.Tasks[0].IsComplete {
		t.Error("Expected IsComplete to be false, but it was true")
	}
}

func TestTaskHandler_NotCompleted_AlreadyUncompletedTask(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks: []models.Task{
				{ID: 1, IsComplete: false},
			},
		},
	}

	err := handler.NotCompleted(0)

	if err == nil {
		t.Fatal("Expected an error when un-completing an already un-completed task, but got nil")
	}
}

func TestTaskHandler_Edit_TextSuccessfullyUpdated(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks: []models.Task{
				{ID: 1, Text: "Old Text"},
			},
		},
	}
	newText := "New updated text"

	handler.Edit(0, newText)

	if handler.Todo.Tasks[0].Text != newText {
		t.Errorf("Expected text to be '%s', but got '%s'", newText, handler.Todo.Tasks[0].Text)
	}
}

func TestTaskHandler_EditTaskPoints_PointsSuccessfullyChanged(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks: []models.Task{
				{ID: 1, TaskPoints: 10},
			},
		},
	}
	newPoints := 50

	handler.EditTaskPoints(0, newPoints)

	if handler.Todo.Tasks[0].TaskPoints != newPoints {
		t.Errorf("Expected points to be %d, but got %d", newPoints, handler.Todo.Tasks[0].TaskPoints)
	}
}

func TestTaskHandler_Delete_TaskRemovedAndMoved(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks: []models.Task{
				{ID: 1, Text: "Task 1"},
				{ID: 2, Text: "Task 2 to delete"},
				{ID: 3, Text: "Task 3"},
			},
			DeletedTasks: []models.Task{},
		},
	}

	handler.Delete(1)

	if len(handler.Todo.Tasks) != 2 {
		t.Fatalf("Expected 2 tasks to remain, but got %d", len(handler.Todo.Tasks))
	}
	if len(handler.Todo.DeletedTasks) != 1 {
		t.Fatalf("Expected 1 task in the deleted list, but got %d", len(handler.Todo.DeletedTasks))
	}
	if handler.Todo.Tasks[0].ID != 1 || handler.Todo.Tasks[1].ID != 3 {
		t.Error("The remaining tasks are not the ones expected")
	}
	if handler.Todo.DeletedTasks[0].ID != 2 {
		t.Error("The wrong task was moved to the deleted list")
	}
}

func TestTaskHandler_ClearAllTasks_ListCleared(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks:  []models.Task{{ID: 1}, {ID: 2}},
			NextID: 3,
		},
	}

	handler.ClearAllTasks()

	if len(handler.Todo.Tasks) != 0 {
		t.Errorf("Expected an empty task list, but it has %d elements", len(handler.Todo.Tasks))
	}
	if handler.Todo.NextID != 1 {
		t.Errorf("Expected NextID to be reset to 1, but it is %d", handler.Todo.NextID)
	}
}

func TestTaskHandler_CancelLastDelete_SuccessfullyRestoresTask(t *testing.T) {
	deletedTask := models.Task{ID: 2, Text: "Deleted Task", IsComplete: true}
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks:        []models.Task{{ID: 1, Text: "Existing Task"}},
			DeletedTasks: []models.Task{deletedTask},
			NextID:       3,
		},
	}

	err := handler.CancelLastDelete()

	if err != nil {
		t.Fatalf("CancelLastDelete() returned an unexpected error: %v", err)
	}
	if len(handler.Todo.DeletedTasks) != 0 {
		t.Error("Expected the deleted tasks list to be empty")
	}
	if len(handler.Todo.Tasks) != 2 {
		t.Fatalf("Expected 2 tasks in the main list, but got %d", len(handler.Todo.Tasks))
	}

	restoredTask := handler.Todo.Tasks[1]
	if restoredTask.Text != deletedTask.Text {
		t.Error("The text of the restored task does not match the deleted one")
	}
	if restoredTask.ID != 3 {
		t.Errorf("Expected the restored task to have a new ID of 3, but got %d", restoredTask.ID)
	}
	if restoredTask.IsComplete {
		t.Error("Expected the restored task to be marked as not complete, but it was")
	}
	if handler.Todo.NextID != 4 {
		t.Errorf("Expected NextID to be incremented to 4, but it is %d", handler.Todo.NextID)
	}
}

func TestTaskHandler_CancelLastDelete_WhenNothingToDelete(t *testing.T) {
	handler := &TaskHandler{
		Todo: &models.TodoList{
			Tasks:        []models.Task{{ID: 1}},
			DeletedTasks: []models.Task{},
		},
	}

	err := handler.CancelLastDelete()

	if err == nil {
		t.Fatal("Expected an error when canceling with no deleted tasks, but got nil")
	}
}
