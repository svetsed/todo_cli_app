package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Storage struct {
		TodoFile   string `mapstructure:"todo_file"`
		RewardFile string `mapstructure:"reward_file"`
	} `mapstructure:"storage"`
	Defaults struct {
		TaskPoints  int `mapstructure:"task_points"`
		RewardPrice int `mapstructure:"reward_price"`
	} `mapstructure:"defaults"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("storage.todo_file", "todo.json")
	viper.SetDefault("storage.reward_file", "rewards.json")
	viper.SetDefault("defaults.task_points", 20)
	viper.SetDefault("defaults.reward_price", 20)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found. Creating a new one 'config.yaml' with default values.")

			if writeErr := viper.WriteConfigAs("./config.yaml"); writeErr != nil {
				return nil, fmt.Errorf("could not create config file: %w", writeErr)
			}
		} else {
			return nil, fmt.Errorf("could not reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

func SaveConfig() error {
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("could not save settings to config file: %v", err)
	}
	return nil
}

func EditPriceOfRewardByDefault(newPriceOfRewardByDef int) {
	viper.Set("defaults.reward_price", newPriceOfRewardByDef)
}

func EditTaskPointsByDefault(newTaskPointsByDef int) {
	viper.Set("defaults.task_points", newTaskPointsByDef)
}
