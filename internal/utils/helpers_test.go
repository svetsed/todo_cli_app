package utils

import (
	"testing"

	"github.com/svetsed/todo_cli_app/internal/models"
)

func TestValidateID(t *testing.T) {
	testCases := []struct {
		name      string
		inputID   string
		maxIndex  int
		wantID    int
		shouldErr bool
	}{
		{
			name:      "Correct ID in the middle of the range",
			inputID:   "5",
			maxIndex:  10,
			wantID:    5,
			shouldErr: false,
		},
		{
			name:      "Limit value: min ID",
			inputID:   "1",
			maxIndex:  10,
			wantID:    1,
			shouldErr: false,
		},
		{
			name:      "Limit value: max ID",
			inputID:   "10",
			maxIndex:  10,
			wantID:    10,
			shouldErr: false,
		},
		{
			name:      "Incorrect ID: zero",
			inputID:   "0",
			maxIndex:  10,
			wantID:    -1,
			shouldErr: true,
		},
		{
			name:      "Incorrect ID: more then max index",
			inputID:   "11",
			maxIndex:  10,
			wantID:    -1,
			shouldErr: true,
		},
		{
			name:      "Incorrect value: not a number",
			inputID:   "abc",
			maxIndex:  10,
			wantID:    -1,
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotID, err := ValidateID(tc.inputID, tc.maxIndex)

			if (err != nil) != tc.shouldErr {
				t.Fatalf("ValidateID() return = %v, expected shouldErr=%v", err, tc.shouldErr)
			}

			if !tc.shouldErr && gotID != tc.wantID {
				t.Errorf("ValidateID() = %d, expected %d", gotID, tc.wantID)
			}
		})
	}
}

func TestValidatePointsOrPrice(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		output    int
		shouldErr bool
	}{
		{
			name:      "Correct value for points or price",
			input:     "20",
			output:    20,
			shouldErr: false,
		},
		{
			name:      "Incorrect value: negative number",
			input:     "-5",
			output:    -1,
			shouldErr: true,
		},
		{
			name:      "Incorrect value: not a number",
			input:     "abc",
			output:    -1,
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count, err := ValidatePointsOrPrice(tc.input)

			if err != nil && !tc.shouldErr {
				t.Fatalf("ValidatePointsOrPrice() return error %v, expected shouldErr=%v", err, tc.shouldErr)
			}

			if err == nil && tc.shouldErr {
				t.Fatalf("ValidatePointsOrPrice() NOT return error %v, expected shouldErr=%v", err, tc.shouldErr)
			}

			if !tc.shouldErr && tc.output != count {
				t.Errorf("ValidatePointsOrPrice() = %d, expected %d", count, tc.output)
			}
		})
	}
}

func TestCheckExistItem(t *testing.T) {
	sliceOfTasks := []models.Task{
		{ID: 1, Text: "Купить хлеб"},
		{ID: 2, Text: "Забрать заказ"},
		{ID: 3, Text: "Посмотреть видеоурок"},
	}

	sliceOfRewards := []models.Reward{
		{ID: 1, Description: "Кофебрейк"},
		{ID: 2, Description: "Посмотреть сериал"},
		{ID: 3, Description: "Порисовать"},
	}

	testCases := []struct {
		name          string
		searchID      int
		inputSlice1   []models.Task
		inputSlice2   []models.Reward
		expectedIndex int
		shouldErr     bool
	}{
		{
			name:          "Searching exist element",
			searchID:      1,
			inputSlice1:   sliceOfTasks,
			inputSlice2:   sliceOfRewards,
			expectedIndex: 0,
			shouldErr:     false,
		},
		{
			name:          "Searching exist element in the middle",
			searchID:      2,
			inputSlice1:   sliceOfTasks,
			inputSlice2:   sliceOfRewards,
			expectedIndex: 1,
			shouldErr:     false,
		},
		{
			name:          "Searching not exist element",
			searchID:      99,
			inputSlice1:   sliceOfTasks,
			inputSlice2:   sliceOfRewards,
			expectedIndex: -1,
			shouldErr:     true,
		},
		{
			name:          "Searching in empty slice",
			searchID:      1,
			inputSlice1:   []models.Task{},
			inputSlice2:   []models.Reward{},
			expectedIndex: -1,
			shouldErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// sliceOfTasks
			gotIndexTasks, err := CheckExistItem(tc.searchID, tc.inputSlice1)

			if err != nil && !tc.shouldErr {
				t.Fatalf("CheckExistItem() return error %v, expected shouldErr=%v", err, tc.shouldErr)
			}

			if err == nil && tc.shouldErr {
				t.Fatalf("CheckExistItem() NOT return error %v, expected shouldErr=%v", err, tc.shouldErr)
			}

			if !tc.shouldErr && gotIndexTasks != tc.expectedIndex {
				t.Errorf("CheckExistItem() return index %d, expected %d", gotIndexTasks, tc.expectedIndex)
			}

			// sliceOfRewards
			gotIndexRewards, err := CheckExistItem(tc.searchID, tc.inputSlice2)

			if err != nil && !tc.shouldErr {
				t.Fatalf("CheckExistItem() return error %v, expected shouldErr=%v", err, tc.shouldErr)
			}

			if err == nil && tc.shouldErr {
				t.Fatalf("CheckExistItem() NOT return error %v, expected shouldErr=%v", err, tc.shouldErr)
			}

			if !tc.shouldErr && gotIndexRewards != tc.expectedIndex {
				t.Errorf("CheckExistItem() return index %d, expected %d", gotIndexRewards, tc.expectedIndex)
			}
		})
	}
}
