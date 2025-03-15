package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	DefaultCollection string   `yaml:"defaultCollection"`
	Collections       []string `yaml:"collections"`
}

var (
	configDir     string
	commandsDir   string
	configFile    string
	currentConfig Config
)

// InitConfig initializes the configuration
func InitConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	// Check if CMDVAULT_CONFIG environment variable is set
	envConfigDir := os.Getenv("CMDVAULT_CONFIG")
	if envConfigDir != "" {
		configDir = envConfigDir
	} else {
		configDir = filepath.Join(homeDir, ".config", "cmdvault")
	}

	commandsDir = filepath.Join(configDir, "commands")
	configFile = filepath.Join(configDir, "config.yaml")

	// Create directories if they don't exist
	os.MkdirAll(configDir, 0755)
	os.MkdirAll(commandsDir, 0755)

	// Load or create default config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		currentConfig = Config{
			DefaultCollection: "default",
			Collections:       []string{"default"},
		}
		SaveConfig()
	} else {
		data, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}

		if err := yaml.Unmarshal(data, &currentConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Ensure default collection directory exists
	os.MkdirAll(filepath.Join(commandsDir, currentConfig.DefaultCollection), 0755)
}

// GetConfig returns the current configuration
func GetConfig() Config {
	return currentConfig
}

// SaveConfig saves the current configuration
func SaveConfig() error {
	data, err := yaml.Marshal(currentConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

// GetDefaultCollection returns the default collection
func GetDefaultCollection() string {
	return currentConfig.DefaultCollection
}

// SetDefaultCollection sets the default collection
func SetDefaultCollection(collection string) error {
	currentConfig.DefaultCollection = collection
	return SaveConfig()
}

// AddCollection adds a new collection
func AddCollection(collection string) error {
	// Check if collection already exists
	for _, c := range currentConfig.Collections {
		if c == collection {
			return nil
		}
	}

	currentConfig.Collections = append(currentConfig.Collections, collection)
	os.MkdirAll(filepath.Join(commandsDir, collection), 0755)
	return SaveConfig()
}

// GetCollectionPath returns the path to a collection directory
func GetCollectionPath(collection string) string {
	return filepath.Join(commandsDir, collection)
}

// GetConfigDir returns the config directory
func GetConfigDir() string {
	return configDir
}

// GetCommandsDir returns the commands directory
func GetCommandsDir() string {
	return commandsDir
}
