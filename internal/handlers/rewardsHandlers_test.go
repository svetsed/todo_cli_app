package handlers

import (
	"testing"

	"github.com/svetsed/todo_cli_app/internal/models"
)

func TestRewardHandler_AddReward(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			Rewards: []models.Reward{},
			NextID:  1,
		},
	}

	r.AddReward("Test reward", 10)

	if len(r.RSystem.Rewards) != 1 {
		t.Fatalf("Expected 1 reward, got %d", len(r.RSystem.Rewards))
	}

	reward := r.RSystem.Rewards[0]
	if reward.Description != "Test reward" {
		t.Errorf("Expected description 'Test reward', got '%s'", reward.Description)
	}
	if reward.PriceOfReward != 10 {
		t.Errorf("Expected price 10, got %d", reward.PriceOfReward)
	}
	if reward.ID != 1 {
		t.Errorf("Expected ID 1, got %d", reward.ID)
	}
}

func TestRewardHandler_BuyRewards(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			UserPoints: 20,
			Rewards: []models.Reward{
				{ID: 1, Description: "Reward 1", PriceOfReward: 15, IsAvailable: true},
			},
		},
	}

	err := r.BuyRewards(0)
	if err != nil {
		t.Fatalf("Unexpected error buying reward: %v", err)
	}

	if r.RSystem.UserPoints != 5 {
		t.Errorf("Expected user points 5 after purchase, got %d", r.RSystem.UserPoints)
	}

	r.RSystem.UserPoints = 10
	err = r.BuyRewards(0)
	if err == nil {
		t.Fatal("Expected error when buying reward with insufficient points, got nil")
	}
}

func TestRewardHandler_EditDesrcRewards(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			Rewards: []models.Reward{
				{ID: 1, Description: "Old Description", PriceOfReward: 10},
			},
		},
	}

	r.EditDesrcRewards(0, "New Description")

	if r.RSystem.Rewards[0].Description != "New Description" {
		t.Errorf("Expected description 'New Description', got '%s'", r.RSystem.Rewards[0].Description)
	}
}

func TestRewardHandler_EditPriceRewards(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			UserPoints: 10,
			Rewards: []models.Reward{
				{ID: 1, Description: "Reward", PriceOfReward: 10, IsAvailable: true},
			},
		},
	}

	r.EditPriceRewards(0, 5)
	if r.RSystem.Rewards[0].PriceOfReward != 5 {
		t.Errorf("Expected price 5, got %d", r.RSystem.Rewards[0].PriceOfReward)
	}

	if r.RSystem.Rewards[0].IsAvailable != true {
		t.Errorf("Expected reward to be available after price change")
	}

	r.EditPriceRewards(0, 15)
	if r.RSystem.Rewards[0].IsAvailable != false {
		t.Errorf("Expected reward to be not available after price increase")
	}
}

func TestRewardHandler_DeleteReward(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			Rewards: []models.Reward{
				{ID: 1, Description: "Reward to delete", PriceOfReward: 10},
			},
		},
	}

	r.DeleteReward(0)
	if len(r.RSystem.Rewards) != 0 {
		t.Errorf("Expected 0 rewards after deletion, got %d", len(r.RSystem.Rewards))
	}
}

func TestRewardHandler_UpdateUserPoints(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			UserPoints:         10,
			IsUserPointsUpdate: false,
		},
	}

	r.UpdateUserPoints(5)
	if r.RSystem.UserPoints != 15 {
		t.Errorf("Expected UserPoints 15, got %d", r.RSystem.UserPoints)
	}
	if !r.RSystem.IsUserPointsUpdate {
		t.Error("Expected IsUserPointsUpdate to be true after update")
	}

	r.UpdateUserPoints(-3)
	if r.RSystem.UserPoints != 12 {
		t.Errorf("Expected UserPoints 12, got %d", r.RSystem.UserPoints)
	}
}

func TestRewardHandler_ClearAllRewards(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			Rewards: []models.Reward{
				{ID: 1, Description: "Reward", PriceOfReward: 10},
			},
			NextID: 2,
		},
	}

	r.ClearAllRewards()
	if len(r.RSystem.Rewards) != 0 {
		t.Errorf("Expected 0 rewards after clear, got %d", len(r.RSystem.Rewards))
	}
	if r.RSystem.NextID != 1 {
		t.Errorf("Expected NextID reset to 1 after clear, got %d", r.RSystem.NextID)
	}
}

func TestRewardHandler_ResetPoints(t *testing.T) {
	r := &RewardHandler{
		RSystem: &models.RewardSystem{
			UserPoints:         10,
			IsUserPointsUpdate: false,
		},
	}

	r.ResetPoints()
	if r.RSystem.UserPoints != 0 {
		t.Errorf("Expected UserPoints 0 after reset, got %d", r.RSystem.UserPoints)
	}
	if !r.RSystem.IsUserPointsUpdate {
		t.Error("Expected IsUserPointsUpdate to be true after reset")
	}
}
