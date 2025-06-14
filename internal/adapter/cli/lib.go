package cli

import (
	"encoding/json"
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/spf13/viper"
)

// ValidateConfigFile checks if the specified YAML configuration file is valid
func ValidateConfigFile(path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(yamlExt)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// LoadConfig loads the configuration from the specified YAML file
func LoadConfig(path string) (*config.Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(yamlExt)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg config.Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// PrintConfigJSON prints the config in JSON format
func PrintConfigJSON(cfg *config.Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
