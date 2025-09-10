package models

type Task struct {
	ID                  int    `json:"id"`
	Text                string `json:"text"`
	IsComplete          bool   `json:"isComplete"`
	TaskPoints          int    `json:"taskPoints"`
	IsTaskPointsReceive bool   `json:"isTaskPointsReceive"`
}

type TodoList struct {
	Tasks        []Task `json:"tasks"`
	DeletedTasks []Task `json:"deletedTasks"`
	NextID       int    `json:"nextId"`
}
