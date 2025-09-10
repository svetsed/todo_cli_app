package utils

import (
	"fmt"
	"strconv"

	"github.com/svetsed/todo_cli_app/internal/models"
)

func ValidatePointsOrPrice(countString string) (int, error) {
	newCount, err := strconv.Atoi(countString)
	if err != nil {
		return -1, fmt.Errorf("incorrect count (%v)", err)
	}
	if newCount < 0 {
		return -1, fmt.Errorf("count must be positive")
	}
	return newCount, nil
}

func ValidateID(idString string, maxIndex int) (int, error) {
	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 || id > maxIndex {
		return -1, fmt.Errorf("incorrect id")
	}
	return id, nil
}

func CalculateIsAvailableReward(userPoints int, price int) (int, bool) {
	if userPoints >= price {
		return userPoints - price, true
	}
	return price - userPoints, false
}

func CheckExistItem(id int, someSliceOfItems any) (int, error) {
	var indexElem int = -1

	switch sliceOfItems := someSliceOfItems.(type) {
	case []models.Task:
		for i, item := range sliceOfItems {
			if item.ID == id {
				indexElem = i
				break
			}
		}
	case []models.Reward:
		for i, item := range sliceOfItems {
			if item.ID == id {
				indexElem = i
				break
			}
		}
	default:
		return -1, fmt.Errorf("give not available type of slice")
	}

	if indexElem == -1 {
		return -1, fmt.Errorf("task %d was not found", id)
	}

	return indexElem, nil
}

func PrintInfoOfTask(id, taskIndexElem int, tasks []models.Task) string {
	status := " "
	if tasks[taskIndexElem].IsComplete {
		status = "✓"
	}

	return fmt.Sprintf("[%s] %d. %s\n", status, id, tasks[taskIndexElem].Text)
}
