package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svetsed/todo_cli_app/cmd/rewards"
	"github.com/svetsed/todo_cli_app/cmd/tasks"
	"github.com/svetsed/todo_cli_app/internal/config"
)

func RootCmd(cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{Use: "todo", Short: "A todo list for the terminal"}

	completeCmd := tasks.CompleteCmd(cfg)
	completeCmd.Flags().BoolP("delete", "d", false, "Delete task after completion")
	completeCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation (only with -d)")

	deleteRewardCmd := rewards.DeleteRewardCmd(cfg)
	deleteRewardCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")
	deleteCmd := tasks.DeleteCmd(cfg)
	deleteCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")
	deleteCmd.AddCommand(deleteRewardCmd)

	addRewardCmd := rewards.AddRewardCmd(cfg)
	addRewardCmd.Flags().IntP("price", "p", 0, "Price in points for the reward")
	addCmd := tasks.AddCmd(cfg)
	addCmd.Flags().IntP("points", "p", 0, "Counts of points, what you will receive after completing the task")
	addCmd.AddCommand(addRewardCmd)

	listRewardCmd := rewards.ListRewardCmd(cfg)
	listCmd := tasks.ListCmd(cfg)
	listCmd.Flags().BoolP("points", "p", false, "Show info about points, what you can receive for the task")
	listCmd.AddCommand(listRewardCmd)

	editRewardPriceByDefault := rewards.EditRewardPriceByDefaultCmd()
	editRewardDescrCmd := rewards.EditRewardDescrCmd(cfg)
	editRewardPriceCmd := rewards.EditRewardPriceCmd(cfg)
	editTaskPointsByDefault := tasks.EditTaskPointsByDefaultCmd()
	editTaskPoints := tasks.EditTaskPointsCmd(cfg)
	editCmd := tasks.EditCmd(cfg)
	editCmd.AddCommand(
		editRewardDescrCmd,
		editRewardPriceCmd,
		editRewardPriceByDefault,
		editTaskPointsByDefault,
		editTaskPoints,
	)

	clearRewardCmd := rewards.ClearRewardCmd(cfg)
	clearCmd := tasks.ClearCmd(cfg)
	clearCmd.AddCommand(clearRewardCmd)

	rootCmd.AddCommand(
		addCmd,
		completeCmd,
		deleteCmd,
		listCmd,
		editCmd,
		clearCmd,
		tasks.NotCompletedCmd(cfg),
		tasks.CancelLastDeleteCmd(cfg),
		rewards.BuyRewardCmd(cfg),
		rewards.ResetPointsCmd(cfg),
	)

	return rootCmd
}
