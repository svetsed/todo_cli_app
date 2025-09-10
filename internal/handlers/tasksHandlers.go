package handlers

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/svetsed/todo_cli_app/internal/models"
)

type TaskHandler struct {
	Todo *models.TodoList `json:"Todolist"`
}

func (h *TaskHandler) Add(text string, points int) {
	if len(h.Todo.Tasks) == 0 {
		h.Todo.NextID = 1
	}
	task := models.Task{
		ID:                  h.Todo.NextID,
		Text:                text,
		IsComplete:          false,
		TaskPoints:          points,
		IsTaskPointsReceive: false,
	}
	h.Todo.Tasks = append(h.Todo.Tasks, task)
	h.Todo.NextID++
}

func (h *TaskHandler) List(pointsFlag bool, writer io.Writer) {
	if len(h.Todo.Tasks) == 0 {
		fmt.Fprintln(writer, "No tasks, well done!")
		return
	}

	w := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)

	if !pointsFlag {
		fmt.Fprintf(w, "Done\tID\tTask\n")
		for _, task := range h.Todo.Tasks {
			status := " "
			if task.IsComplete {
				status = "✓"
			}
			fmt.Fprintf(w, "[%s]\t%d.\t%s\n", status, task.ID, task.Text)
		}
	} else {
		fmt.Fprintf(w, "Done\tID\tTask\tPoints for task\n")
		for _, task := range h.Todo.Tasks {
			status := " "
			if task.IsComplete {
				status = "✓"
			}
			fmt.Fprintf(w, "[%s]\t%d.\t%s\t%d\n", status, task.ID, task.Text, task.TaskPoints)
		}
	}

	w.Flush()
}

func (h *TaskHandler) Complete(indexCompElem int) error {
	if h.Todo.Tasks[indexCompElem].IsComplete {
		return fmt.Errorf("the task has already been completed")
	}
	h.Todo.Tasks[indexCompElem].IsComplete = true
	return nil
}

func (h *TaskHandler) NotCompleted(indexNotCompElem int) error {
	if !h.Todo.Tasks[indexNotCompElem].IsComplete {
		return fmt.Errorf("the task has not been completed yet")
	}
	h.Todo.Tasks[indexNotCompElem].IsComplete = false
	return nil
}

func (h *TaskHandler) Edit(indexEditElem int, newText string) {
	h.Todo.Tasks[indexEditElem].Text = newText
}

func (h *TaskHandler) EditTaskPoints(indexPointsElem int, newTaskPoints int) {
	h.Todo.Tasks[indexPointsElem].TaskPoints = newTaskPoints
}

func (h *TaskHandler) ClearAllTasks() {
	h.Todo.Tasks = []models.Task{}
	h.Todo.NextID = 1
}

func (h *TaskHandler) Delete(indexDelElem int) {
	h.Todo.DeletedTasks = append(h.Todo.DeletedTasks, h.Todo.Tasks[indexDelElem])
	h.Todo.Tasks = append(h.Todo.Tasks[:indexDelElem], h.Todo.Tasks[indexDelElem+1:]...)
}

func (h *TaskHandler) CancelLastDelete() error {
	if len(h.Todo.DeletedTasks) == 0 {
		return fmt.Errorf("no deleted tasks")
	}

	last := h.Todo.DeletedTasks[len(h.Todo.DeletedTasks)-1]
	last.ID = h.Todo.NextID
	h.Todo.NextID++
	last.IsComplete = false

	h.Todo.Tasks = append(h.Todo.Tasks, last)
	h.Todo.DeletedTasks = []models.Task{}

	return nil
}
