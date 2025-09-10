package handlers

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/svetsed/todo_cli_app/internal/models"
	"github.com/svetsed/todo_cli_app/internal/utils"
)

type RewardHandler struct {
	RSystem *models.RewardSystem `json:"rewardSystem"`
}

func (r *RewardHandler) AddReward(desrc string, price int) {
	if len(r.RSystem.Rewards) == 0 {
		r.RSystem.NextID = 1
	}
	_, isAvailable := utils.CalculateIsAvailableReward(r.RSystem.UserPoints, price)
	reward := models.Reward{
		ID:            r.RSystem.NextID,
		Description:   desrc,
		PriceOfReward: price,
		IsAvailable:   isAvailable,
	}
	r.RSystem.Rewards = append(r.RSystem.Rewards, reward)
	r.RSystem.NextID++
}

func (r *RewardHandler) ListReward(writer io.Writer) {
	fmt.Printf("Your balance of points: %d\n", r.RSystem.UserPoints)

	if len(r.RSystem.Rewards) == 0 {
		fmt.Println("The rewards was not added")
		return
	}

	w := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Available\tID\tDescription\tPrice\n")

	if r.RSystem.IsUserPointsUpdate {
		r.UpdateIsAvailableRewards()
	}

	for _, reward := range r.RSystem.Rewards {
		status := "no"
		if reward.IsAvailable {
			status = "yes"
		}
		fmt.Fprintf(w, "  %s\t%d.\t%s\t%d\n", status, reward.ID, reward.Description, reward.PriceOfReward)
	}
	w.Flush()
}

func (r *RewardHandler) BuyRewards(indexBuyElem int) error {
	points, available := utils.CalculateIsAvailableReward(r.RSystem.UserPoints, r.RSystem.Rewards[indexBuyElem].PriceOfReward)
	if available {
		r.RSystem.UserPoints = points
		r.RSystem.IsUserPointsUpdate = true
	} else {
		return fmt.Errorf("not enough %d points", points)
	}
	return nil
}

func (r *RewardHandler) UpdateUserPoints(countPoints int) {
	if countPoints > 0 || countPoints < 0 {
		r.RSystem.UserPoints += countPoints
		r.RSystem.IsUserPointsUpdate = true
	}

}

func (r *RewardHandler) EditDesrcRewards(indexEditElem int, newDesrc string) {
	r.RSystem.Rewards[indexEditElem].Description = newDesrc
}

func (r *RewardHandler) EditPriceRewards(indexEditElem int, newPrice int) {
	r.RSystem.Rewards[indexEditElem].PriceOfReward = newPrice
	_, r.RSystem.Rewards[indexEditElem].IsAvailable = utils.CalculateIsAvailableReward(r.RSystem.UserPoints, newPrice)
}

func (r *RewardHandler) DeleteReward(indexDelElem int) {
	r.RSystem.Rewards = append(r.RSystem.Rewards[:indexDelElem], r.RSystem.Rewards[indexDelElem+1:]...)
}

func (r *RewardHandler) UpdateIsAvailableRewards() {
	userPoints := r.RSystem.UserPoints
	for i := range r.RSystem.Rewards {
		_, r.RSystem.Rewards[i].IsAvailable = utils.CalculateIsAvailableReward(userPoints, r.RSystem.Rewards[i].PriceOfReward)
	}

	r.RSystem.IsUserPointsUpdate = false
}

func (r *RewardHandler) ClearAllRewards() {
	r.RSystem.Rewards = []models.Reward{}
	r.RSystem.NextID = 1
}

func (r *RewardHandler) ResetPoints() {
	r.RSystem.UserPoints = 0
	r.RSystem.IsUserPointsUpdate = true
}
