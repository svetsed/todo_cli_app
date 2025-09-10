package tasks

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/svetsed/todo_cli_app/internal/config"
	"github.com/svetsed/todo_cli_app/internal/handlers"
	"github.com/svetsed/todo_cli_app/internal/loaders"
	"github.com/svetsed/todo_cli_app/internal/logger"
	"github.com/svetsed/todo_cli_app/internal/storage"
	"github.com/svetsed/todo_cli_app/internal/utils"
)

func AddCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "add <task text> [flags] [count of points with flag -p]",
		Short: "Add a new task with optional points, what you will receive after completing the task",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("add command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}

			var pointsCount int = cfg.Defaults.TaskPoints
			if cmd.Flags().Changed("points") {
				if pointsCountTmp, err := cmd.Flags().GetInt("points"); err != nil {
					logger.Error("could not parse points flag in add command", err)
					fmt.Println("For this task will set default count of points")
				} else if pointsCountTmp < 0 {
					fmt.Println("Count of points must be positive")
					fmt.Println("For this task will set default count of points")
				} else {
					pointsCount = pointsCountTmp
				}
			}

			text := strings.Join(args, " ")

			h.Add(text, pointsCount)
			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list in add command", err, slog.String("file", cfg.Storage.TodoFile))
			} else {
				logger.Info("task added successfully", slog.Int("task_id", h.Todo.NextID-1))
				fmt.Printf("Added task: [ ] %d. %s (%d points)\n", h.Todo.NextID-1, text, h.Todo.Tasks[len(h.Todo.Tasks)-1].TaskPoints)
			}
		},
	}
}

func ListCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list [flags]",
		Short: "Show all tasks and optional points for the task",
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("list command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}

			pointsFlag, err := cmd.Flags().GetBool("points")
			if err != nil {
				logger.Error("could not parse points flag in list command: %v", err)
				return
			}
			h.List(pointsFlag, cmd.OutOrStdout())
		},
	}
}

func CompleteCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "complete <ID> [flags]",
		Short: "Mark the task as completed and/or delete this task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("complete command failed", err)
				return
			}

			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("complete command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}
			r := &handlers.RewardHandler{RSystem: rewardSystem}

			id, err := utils.ValidateID(args[0], h.Todo.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "complete"))
				return
			}

			var taskIndexElem int
			if i, err := utils.CheckExistItem(id, h.Todo.Tasks); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "complete"))
				return
			} else {
				taskIndexElem = i
			}

			if err := h.Complete(taskIndexElem); err != nil {
				logger.Error("could not mark task as completed", err)
				return
			}

			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list after completing the task", err, slog.String("file", cfg.Storage.TodoFile))
				return
			}

			deleteFlag, err := cmd.Flags().GetBool("delete")
			if err != nil {
				logger.Error("could not parse delete flag", err, slog.String("flag", "--delete"), slog.String("command", "complete"))
				return
			}
			forceFlag, err := cmd.Flags().GetBool("force")
			if err != nil {
				logger.Error("could not parse force flag", err, slog.String("flag", "--force"), slog.String("command", "complete"))
				return
			}

			if forceFlag && !deleteFlag {
				logger.Error("incorrect using flags", fmt.Errorf("flag -f can only be used with -d"))
				_ = cmd.Usage()
				return
			}

			if !h.Todo.Tasks[taskIndexElem].IsTaskPointsReceive {
				r.UpdateUserPoints(h.Todo.Tasks[taskIndexElem].TaskPoints)
				if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
					logger.Error("failed to save reward file after updating balance of points", err, slog.String("file", cfg.Storage.RewardFile), slog.String("command", "complete"))
					return
				}
				h.Todo.Tasks[taskIndexElem].IsTaskPointsReceive = true
			}

			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list after completing the task", err, slog.String("file", cfg.Storage.TodoFile))
				return
			}

			printingTask := utils.PrintInfoOfTask(id, taskIndexElem, h.Todo.Tasks)
			if deleteFlag {
				if !forceFlag {
					fmt.Printf("You want to delete completed task: %s", printingTask)
					fmt.Print("Are you sure (y/n): ")
					var confirm string
					fmt.Scanln(&confirm)
					if strings.ToLower(confirm) != "y" {
						logger.Info("the deletion was cancelled by user")
						fmt.Println("The deletion was cancelled!")
						return
					}
				}

				h.Delete(taskIndexElem)

				if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
					logger.Error("failed to save todo list after deleting completed task", err, slog.String("file", cfg.Storage.TodoFile))
					return
				}
				logger.Info("task has been completed and deleted", slog.Int("id", id))
				fmt.Printf("Well done! Task %d has been completed and deleted!\n", id)

			} else {
				logger.Info("task was marked as completed", slog.Int("id", id))
				fmt.Printf("Well done! Task %d was marked as completed!\n", id)
				fmt.Print(printingTask)
			}
		},
	}
}

func NotCompletedCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "not-complete <ID>",
		Short: "Mark the task as not completed",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("not-complete command failed", err)
				return
			}

			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("not-complete command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}
			r := &handlers.RewardHandler{RSystem: rewardSystem}

			id, err := utils.ValidateID(args[0], h.Todo.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "not-complete"))
				return
			}

			var taskIndexElem int
			if i, err := utils.CheckExistItem(id, h.Todo.Tasks); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "not-complete"))
				return
			} else {
				taskIndexElem = i
			}

			if err := h.NotCompleted(taskIndexElem); err != nil {
				logger.Error("could not mark task as NOT completed", err)
				return
			}

			if h.Todo.Tasks[taskIndexElem].IsTaskPointsReceive {
				r.UpdateUserPoints(-(h.Todo.Tasks[taskIndexElem].TaskPoints))
				if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
					logger.Error("failed to save reward file after updating balance of points", err, slog.String("file", cfg.Storage.RewardFile), slog.String("command", "not-complete"))
					return
				}
				h.Todo.Tasks[taskIndexElem].IsTaskPointsReceive = false
			}

			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list by not-complete command", err, slog.String("file", cfg.Storage.TodoFile))
				return
			}

			logger.Info("task has been not completed", slog.Int("id", id))
			fmt.Printf("Now task %d is not completed!\n", id)
			fmt.Print(utils.PrintInfoOfTask(id, taskIndexElem, h.Todo.Tasks))
		},
	}
}

func EditCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <ID> <task text>",
		Short: "Edits the text of an existing task",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("edit command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}

			id, err := utils.ValidateID(args[0], h.Todo.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "edit"))
				return
			}

			var taskIndexElem int
			if i, err := utils.CheckExistItem(id, h.Todo.Tasks); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "edit"))
				return
			} else {
				taskIndexElem = i
			}

			newText := strings.Join(args[1:], " ")

			h.Edit(taskIndexElem, newText)
			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list after editing text of task", err, slog.String("file", cfg.Storage.TodoFile))
			} else {
				logger.Info("text of task has been changed", slog.Int("id", id))
				fmt.Printf("Task is changed: %s", utils.PrintInfoOfTask(id, taskIndexElem, h.Todo.Tasks))
			}
		},
	}
}

func EditTaskPointsByDefaultCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pointsdef <new count of points>",
		Short: "Edits the count of points, what set for tasks by default",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			newTaskPointsByDef, err := utils.ValidatePointsOrPrice(args[0])
			if err != nil {
				logger.Error("incorrect number for count of points", err, slog.String("command", "edit task points by default"))
				return
			}
			config.EditTaskPointsByDefault(newTaskPointsByDef)

			if err := config.SaveConfig(); err != nil {
				logger.Error("failed to save config file", err)
			} else {
				logger.Info("count of points by default has been changed", slog.Int("new count of points", newTaskPointsByDef))
				fmt.Printf("Count of points by default has been changed: %d\n", newTaskPointsByDef)
			}
		},
	}
}

func EditTaskPointsCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "points <ID> <new count of points>",
		Short: "Edits the count of points for existing task",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("edit points command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}

			id, err := utils.ValidateID(args[0], h.Todo.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "edit task points by id"))
				return
			}

			var taskIndexElem int
			if i, err := utils.CheckExistItem(id, h.Todo.Tasks); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "edit task points by id"))
				return
			} else {
				taskIndexElem = i
			}

			newTaskPoints, err := utils.ValidatePointsOrPrice(args[1])
			if err != nil {
				logger.Error("incorrect number for count of points", err, slog.String("command", "edit task points by id"))
				return
			}

			h.EditTaskPoints(taskIndexElem, newTaskPoints)
			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list after editing count of points in task", err, slog.String("file", cfg.Storage.TodoFile))
			} else {
				logger.Info("count of points has been changed for task", slog.Int("id", id))
				fmt.Printf("Count of points has been changed for task %d: %s (%d points)\n", id, h.Todo.Tasks[taskIndexElem].Text, h.Todo.Tasks[taskIndexElem].TaskPoints)
			}
		},
	}
}

func DeleteCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <ID> [flags]",
		Short: "Delete the task from list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("delete command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}

			id, err := utils.ValidateID(args[0], h.Todo.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "delete"))
				return
			}

			var taskIndexElem int
			if i, err := utils.CheckExistItem(id, h.Todo.Tasks); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "delete"))
				return
			} else {
				taskIndexElem = i
			}

			forceFlag, err := cmd.Flags().GetBool("force")
			if err != nil {
				logger.Error("could not parse force flag", err, slog.String("flag", "--force"), slog.String("command", "delete"))
				return
			}
			if !forceFlag {
				fmt.Printf("You want to delete task: %s", utils.PrintInfoOfTask(id, taskIndexElem, h.Todo.Tasks))
				fmt.Print("Are you sure (y/n): ")
				var confirm string
				fmt.Scanln(&confirm)
				if strings.ToLower(confirm) != "y" {
					fmt.Println("The deletion was cancelled!")
					logger.Info("the deletion was cancelled by user")
					return
				}
			}
			h.Delete(taskIndexElem)

			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list  after deleting task", err, slog.String("file", cfg.Storage.TodoFile))
			} else {
				fmt.Printf("Task %d was deleted!\n", id)
				logger.Info("task was deleted", slog.Int("id", id))
			}
		},
	}
}

func ClearCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Remove all tasks without restore with confirmation",
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("clear command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}

			fmt.Print("The tasks cannot be restored!\nAre you sure you want to delete ALL tasks? (y/n): ")
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Operation was cancelled!")
				logger.Info("Remove all tasks has been cancelled")
				return
			}

			h.ClearAllTasks()

			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list after clearing all tasks", err, slog.String("file", cfg.Storage.TodoFile))
			} else {
				fmt.Println("All tasks has been removed!")
				logger.Info("All tasks was removed by user")
			}
		},
	}
}

func CancelLastDeleteCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel-delete",
		Short: "Cancels the last delete and returns the task as not completed in todolist with a new ID",
		Run: func(cmd *cobra.Command, args []string) {
			todoList, err := loaders.LoadTodoList(cfg.Storage.TodoFile)
			if err != nil {
				logger.Error("cancel-delete command failed", err)
				return
			}

			h := &handlers.TaskHandler{Todo: todoList}

			if err := h.CancelLastDelete(); err != nil {
				logger.Error("catch error when try canceling last delete", err)
				return
			}

			if err := storage.Save(cfg.Storage.TodoFile, h.Todo); err != nil {
				logger.Error("failed to save todo list after cancelling last deleted task", err, slog.String("file", cfg.Storage.TodoFile))
			} else {
				logger.Info("the task was restored with new id", slog.Int("new id", h.Todo.NextID-1))
				fmt.Printf("The task was restored with NEW ID: %s", utils.PrintInfoOfTask(h.Todo.NextID-1, len(h.Todo.Tasks)-1, h.Todo.Tasks))
			}
		},
	}
}
