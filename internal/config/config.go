package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Token         string `mapstructure:"token"`
	DefaultBudget string `mapstructure:"default_budget"`
	Format        string `mapstructure:"format"`
}

var configDir string
var configFile string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		configDir = "."
	} else {
		configDir = filepath.Join(home, ".config", "ynabctl")
	}
	configFile = filepath.Join(configDir, "config.toml")
}

// Load reads the configuration from file and environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set config file
	v.SetConfigFile(configFile)
	v.SetConfigType("toml")

	// Environment variable support
	v.SetEnvPrefix("YNAB")
	v.AutomaticEnv()

	// Map environment variables
	v.BindEnv("token", "YNAB_TOKEN")
	v.BindEnv("default_budget", "YNAB_DEFAULT_BUDGET")
	v.BindEnv("format", "YNAB_FORMAT")

	// Set defaults
	v.SetDefault("format", "json")

	// Read config file (ignore error if file doesn't exist)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return error if it's not a "file not found" error
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to file
func Save(cfg *Config) error {
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("toml")

	v.Set("token", cfg.Token)
	v.Set("default_budget", cfg.DefaultBudget)
	v.Set("format", cfg.Format)

	if err := v.WriteConfig(); err != nil {
		// If config file doesn't exist, create it
		if os.IsNotExist(err) {
			return v.SafeWriteConfig()
		}
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SetToken saves the API token to config
func SetToken(token string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.Token = token
	return Save(cfg)
}

// SetDefaultBudget saves the default budget ID to config
func SetDefaultBudget(budgetID string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.DefaultBudget = budgetID
	return Save(cfg)
}

// SetFormat saves the default output format to config
func SetFormat(format string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.Format = format
	return Save(cfg)
}

// GetConfigFile returns the path to the config file
func GetConfigFile() string {
	return configFile
}
