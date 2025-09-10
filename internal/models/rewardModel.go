package models

type Reward struct {
	ID            int    `json:"id"`
	Description   string `json:"description"`
	PriceOfReward int    `json:"priceOfReward"`
	IsAvailable   bool   `json:"isAvailable"`
}

type RewardSystem struct {
	Rewards            []Reward `json:"rewards"`
	UserPoints         int      `json:"userPoints"`
	IsUserPointsUpdate bool     `json:"isUserPointsUpdate"`
	NextID             int      `json:"nextId"`
}
