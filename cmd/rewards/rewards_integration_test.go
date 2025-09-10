package rewards

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

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return strings.TrimSpace(out.String()), err
}

func TestIntegration_AddRewardCmd_SuccessfullyAddsReward(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	rewardFile := filepath.Join(tempDir, "test_rewards.json")
	t.Setenv("STORAGE_REWARD_FILE", rewardFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config for test: %v", err)
	}

	initialReward := &models.RewardSystem{
		Rewards:            []models.Reward{},
		UserPoints:         0,
		IsUserPointsUpdate: false,
		NextID:             1,
	}
	initialData, err := json.Marshal(initialReward)
	if err != nil {
		t.Fatalf("Error marshaling initial reward data: %v", err)
	}

	if err := os.WriteFile(cfg.Storage.RewardFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write initial reward file: %v", err)
	}

	addRewardCmd := AddRewardCmd(cfg)
	addRewardCmd.Flags().IntP("price", "p", 0, "Price in points for the reward")

	desc := "Integration test reward"
	price := "15"

	_, err = executeCommand(addRewardCmd, desc, "-p", price)
	if err != nil {
		t.Fatalf("AddRewardCmd returned error: %v", err)
	}

	data, err := os.ReadFile(cfg.Storage.RewardFile)
	if err != nil {
		t.Fatalf("Could not read reward file: %v", err)
	}

	var rewards models.RewardSystem
	if err := json.Unmarshal(data, &rewards); err != nil {
		t.Fatalf("Could not parse reward JSON: %v", err)
	}

	if len(rewards.Rewards) != 1 {
		t.Fatalf("Expected 1 reward, got %d", len(rewards.Rewards))
	}

	if rewards.Rewards[0].Description != desc {
		t.Errorf("Expected reward description %s, got %s", desc, rewards.Rewards[0].Description)
	}
	if rewards.Rewards[0].PriceOfReward != 15 {
		t.Errorf("Expected reward price 15, got %d", rewards.Rewards[0].PriceOfReward)
	}
}

func TestIntegration_ListRewardCmd_ShowsRewards(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	rewardFile := filepath.Join(tempDir, "test_rewards.json")
	t.Setenv("STORAGE_REWARD_FILE", rewardFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config for test: %v", err)
	}

	rewards := &models.RewardSystem{
		Rewards: []models.Reward{
			{ID: 1, Description: "Sample Reward", PriceOfReward: 10, IsAvailable: true},
		},
		UserPoints:         10,
		IsUserPointsUpdate: false,
		NextID:             2,
	}
	initialData, err := json.Marshal(rewards)
	if err != nil {
		t.Fatalf("Error marshaling initial reward data: %v", err)
	}

	if err := os.WriteFile(cfg.Storage.RewardFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write initial reward file: %v", err)
	}

	listCmd := ListRewardCmd(cfg)

	output, err := executeCommand(listCmd)
	if err != nil {
		t.Fatalf("ListRewardCmd returned error: %v", err)
	}

	if !strings.Contains(output, "Sample Reward") {
		t.Errorf("Expected output to contain 'Sample Reward', got:\n%s", output)
	}
	if !strings.Contains(output, "10") {
		t.Errorf("Expected balance info in output, got:\n%s", output)
	}
}

func TestIntegration_BuyRewardCmd_SuccessfulPurchase(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	rewardFile := filepath.Join(tempDir, "test_rewards.json")
	t.Setenv("STORAGE_REWARD_FILE", rewardFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config for test: %v", err)
	}

	rewards := &models.RewardSystem{
		Rewards: []models.Reward{
			{ID: 1, Description: "Reward to Buy", PriceOfReward: 10, IsAvailable: true},
		},
		UserPoints:         20,
		IsUserPointsUpdate: false,
		NextID:             2,
	}
	initialData, err := json.Marshal(rewards)
	if err != nil {
		t.Fatalf("Error marshaling initial reward data: %v", err)
	}

	if err := os.WriteFile(cfg.Storage.RewardFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write initial reward file: %v", err)
	}

	buyCmd := BuyRewardCmd(cfg)
	_, err = executeCommand(buyCmd, "1")
	if err != nil {
		t.Fatalf("BuyRewardCmd returned error: %v", err)
	}

	data, err := os.ReadFile(cfg.Storage.RewardFile)
	if err != nil {
		t.Fatalf("Could not read reward file: %v", err)
	}

	var updatedRewards models.RewardSystem
	if err := json.Unmarshal(data, &updatedRewards); err != nil {
		t.Fatalf("Could not parse reward JSON: %v", err)
	}

	if updatedRewards.UserPoints != 10 {
		t.Errorf("Expected user points to be 10 after purchase, got %d", updatedRewards.UserPoints)
	}
}

func TestIntegration_DeleteRewardCmd_SuccessfulDelete(t *testing.T) {
	logger.Init(slog.LevelDebug, io.Discard)

	tempDir := t.TempDir()
	rewardFile := filepath.Join(tempDir, "test_rewards.json")
	t.Setenv("STORAGE_REWARD_FILE", rewardFile)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load config for test: %v", err)
	}

	rewards := &models.RewardSystem{
		Rewards: []models.Reward{
			{ID: 1, Description: "Reward to Delete", PriceOfReward: 10, IsAvailable: true},
		},
		UserPoints:         10,
		IsUserPointsUpdate: false,
		NextID:             2,
	}
	initialData, err := json.Marshal(rewards)
	if err != nil {
		t.Fatalf("Error marshaling initial reward data: %v", err)
	}

	if err := os.WriteFile(cfg.Storage.RewardFile, initialData, 0666); err != nil {
		t.Fatalf("Failed to write initial reward file: %v", err)
	}

	deleteCmd := DeleteRewardCmd(cfg)
	deleteCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")

	_, err = executeCommand(deleteCmd, "1", "--force")
	if err != nil {
		t.Fatalf("DeleteRewardCmd returned error: %v", err)
	}

	data, err := os.ReadFile(cfg.Storage.RewardFile)
	if err != nil {
		t.Fatalf("Could not read reward file: %v", err)
	}

	var updatedRewards models.RewardSystem
	if err := json.Unmarshal(data, &updatedRewards); err != nil {
		t.Fatalf("Could not parse reward JSON: %v", err)
	}

	if len(updatedRewards.Rewards) != 0 {
		t.Errorf("Expected 0 rewards after deletion, got %d", len(updatedRewards.Rewards))
	}
}
