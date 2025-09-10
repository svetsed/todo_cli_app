package rewards

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

func AddRewardCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "reward <description> [flags]",
		Short: "Add a new reward with optional points after flag -p",
		Long:  "Add a new reward with optional points, what you will receive after completing the task",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("add reward command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			var price int = cfg.Defaults.RewardPrice
			if cmd.Flags().Changed("price") {
				if priceTmp, err := cmd.Flags().GetInt("price"); err != nil {
					logger.Error("could not parse price flag in add reward command", err, slog.String("flag", "--price"))
					fmt.Println("For this reward will set default price")
				} else if priceTmp < 0 {
					logger.Error("Price must be positive", fmt.Errorf("negative price"))
					fmt.Println("For this reward will set default price")
				} else {
					price = priceTmp
				}
			}

			desrc := strings.Join(args, " ")

			r.AddReward(desrc, price)

			if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
				logger.Error("failed to save reward file after reward addition", err, slog.String("file", cfg.Storage.RewardFile))
			} else {
				logger.Info("added new reward", slog.Int("id", r.RSystem.NextID-1))
				fmt.Printf("Added reward: %d. %s (with price %d points)\n", r.RSystem.NextID-1, desrc, price)
			}
		},
	}
}

func ListRewardCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "reward",
		Short: "Show all reward and your balance of points",
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("list reward command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			needSave := r.RSystem.IsUserPointsUpdate
			r.ListReward(cmd.OutOrStdout())
			if needSave {
				if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
					logger.Error("failed to save reward file after show all reward", err, slog.String("file", cfg.Storage.RewardFile))
				}
			}
		},
	}
}

func BuyRewardCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "buy-reward <ID>",
		Short: "Buy the existing reward by ID of reward",
		Long:  "Buy the existing reward, if your balance of points more or equals the price of reward",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("buy-reward command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			id, err := utils.ValidateID(args[0], r.RSystem.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "buy-reward"))
				return
			}

			var rewardIndexElem int
			if i, err := utils.CheckExistItem(id, r.RSystem.Rewards); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "buy-reward"))
				return
			} else {
				rewardIndexElem = i
			}

			if err := r.BuyRewards(rewardIndexElem); err != nil {
				logger.Error("the reward is not yet available", err)
				fmt.Printf("The reward is not yet available: %v\n", err)
				return
			}

			if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
				logger.Error("failed to save reward file after buying reward", err, slog.String("file", cfg.Storage.RewardFile))
				return
			}

			fmt.Println("Good Job! Here is your reward! Enjoy!")
			fmt.Printf("Receive reward: %s\n", r.RSystem.Rewards[rewardIndexElem].Description)
			fmt.Printf("Now your balance: %d\n", r.RSystem.UserPoints)
			logger.Info("receive reward", slog.Int("reward_id", id), slog.Int("balance", r.RSystem.UserPoints))

		},
	}
}

func EditRewardDescrCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "redescr <ID> <new description>",
		Short: "Edits the description of an existing reward",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("edit redescr command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			id, err := utils.ValidateID(args[0], r.RSystem.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "edit redesrc"))
				return
			}

			var rewardIndexElem int
			if i, err := utils.CheckExistItem(id, r.RSystem.Rewards); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "edit redesrc"))
				return
			} else {
				rewardIndexElem = i
			}

			newDescr := strings.Join(args[1:], " ")

			r.EditDesrcRewards(rewardIndexElem, newDescr)
			if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
				logger.Error("failed to save reward file after editing description of an existing reward", err, slog.String("file", cfg.Storage.RewardFile))
			} else {
				logger.Info("desription of reward has been changed", slog.Int("id", id))
				fmt.Printf("Description of reward has been changed: %d. %s with price %d points\n", id, r.RSystem.Rewards[rewardIndexElem].Description, r.RSystem.Rewards[rewardIndexElem].PriceOfReward)
			}
		},
	}
}

func EditRewardPriceCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "reprice <ID> <new price>",
		Short: "Edits the price of an existing reward",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("edit reprice command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			id, err := utils.ValidateID(args[0], r.RSystem.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "edit reprice"))
				return
			}

			var rewardIndexElem int
			if i, err := utils.CheckExistItem(id, r.RSystem.Rewards); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "edit reprice"))
				return
			} else {
				rewardIndexElem = i
			}

			newPrice, err := utils.ValidatePointsOrPrice(args[1])
			if err != nil {
				logger.Error("incorrect number for price of reward", err, slog.String("command", "edit reprice"))
				return
			}

			r.EditPriceRewards(rewardIndexElem, newPrice)
			if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
				logger.Error("failed to save reward file after editing price of an existing reward", err, slog.String("file", cfg.Storage.RewardFile))
			} else {
				logger.Info("price of reward has been changed", slog.Int("reward_id", id))
				fmt.Printf("Price of reward has been changed: %d. %s with new price %d points\n", id, r.RSystem.Rewards[rewardIndexElem].Description, r.RSystem.Rewards[rewardIndexElem].PriceOfReward)
			}
		},
	}
}

func EditRewardPriceByDefaultCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "repricedef <new price>",
		Short: "Edits the default price, when you add reward",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			newPrice, err := utils.ValidatePointsOrPrice(args[0])
			if err != nil {
				logger.Error("incorrect number for price of reward", err, slog.String("command", "edit repricedef"))
				return
			}

			config.EditPriceOfRewardByDefault(newPrice)

			if err := config.SaveConfig(); err != nil {
				logger.Error("failed to save config file", err)
			} else {
				logger.Info("price of reward by default has been changed", slog.Int("new price by default", newPrice))
				fmt.Printf("Price of reward by default has been changed: %d\n", newPrice)
			}
		},
	}
}

func DeleteRewardCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "reward <ID> [flags]",
		Short: "Delete the reward from list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("delete reward command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			id, err := utils.ValidateID(args[0], r.RSystem.NextID-1)
			if err != nil {
				logger.Error("catch error when checking id", err, slog.String("command", "delete reward"))
				return
			}

			var rewardIndexElem int
			if i, err := utils.CheckExistItem(id, r.RSystem.Rewards); err != nil {
				logger.Error("catch error when searching task by id", err, slog.String("command", "delete reward"))
				return
			} else {
				rewardIndexElem = i
			}

			forceFlag, err := cmd.Flags().GetBool("force")
			if err != nil {
				logger.Error("could not parse force flag", err, slog.String("flag", "--force"), slog.String("command", "delete reward"))
				return
			}

			if !forceFlag {
				fmt.Printf("You want to delete reward: %d. %s\n", id, r.RSystem.Rewards[rewardIndexElem].Description)
				fmt.Print("Are you sure (y/n): ")
				var confirm string
				fmt.Scanln(&confirm)
				if strings.ToLower(confirm) != "y" {
					fmt.Println("The deletion was cancelled!")
					logger.Info("the deletion reward was cancelled by user")
					return
				}
			}
			r.DeleteReward(rewardIndexElem)

			if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
				logger.Error("failed to save reward file after deleting reward", err, slog.String("file", cfg.Storage.RewardFile))
			} else {
				logger.Info("the reward was deleted", slog.Int("reward_id", id))
				fmt.Printf("The reward %d was deleted!\n", id)
			}
		},
	}
}

func ClearRewardCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "reward",
		Short: "Remove all rewards without restore with confirmation",
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("clear reward command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			fmt.Print("The rewards cannot be restored!\nAre you sure you want to delete ALL rewards? (y/n): ")
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Operation was cancelled!")
				logger.Info("Remove all rewards has been cancelled")
				return
			}

			r.ClearAllRewards()

			if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
				logger.Error("failed to save reward file after clearing all rewards", err, slog.String("file", cfg.Storage.RewardFile))
			} else {
				fmt.Println("All rewards have been removed!")
				logger.Info("all rewards have been removed by user")
			}
		},
	}
}

func ResetPointsCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "resetp",
		Short: "Reset to zero all your points without restore with confirmation",
		Run: func(cmd *cobra.Command, args []string) {
			rewardSystem, err := loaders.LoadRewardSystem(cfg.Storage.RewardFile)
			if err != nil {
				logger.Error("resetp command failed", err)
				return
			}

			r := &handlers.RewardHandler{RSystem: rewardSystem}

			fmt.Print("The count of points cannot be restored!\nAre you sure you want to reset to zero your points? (y/n): ")
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Operation was cancelled!")
				logger.Info("Reset to zero all points of user has been cancelled")
				return
			}

			r.ResetPoints()

			if err := storage.Save(cfg.Storage.RewardFile, r.RSystem); err != nil {
				logger.Error("failed to save reward file after reseting points to zero", err, slog.String("file", cfg.Storage.RewardFile))
			} else {
				fmt.Printf("Now your balance of points: %d\n!", r.RSystem.UserPoints)
				logger.Info("User reset to zero all your points")
			}
		},
	}
}
