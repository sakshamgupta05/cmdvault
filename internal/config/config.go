package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	CollectionDirs []string `yaml:"collectionDirs"`
}

var (
	configDir      string
	collectionDirs []string
	configFile     string
	currentConfig  Config
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

	configFile = filepath.Join(configDir, "config.yaml")

	// Create directories if they don't exist
	os.MkdirAll(configDir, 0755)

	// Load or create default config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		currentConfig = Config{}
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

	collectionDirs = currentConfig.CollectionDirs
	if len(collectionDirs) == 0 {
		collectionDirs = []string{filepath.Join(configDir, "collections")}
	}

	// Create collection directories if they don't exist
	for _, dir := range collectionDirs {
		os.MkdirAll(dir, 0755)
	}
}

// GetConfig returns the current configuration
func GetConfig() Config {
	return currentConfig
}

// GetCollectionPath returns the path to a collection file
func GetCollectionPath(collection string) (string, error) {
	collection = collection + ".yaml"
	// Traverse collection paths in reverse order
	for i := len(collectionDirs) - 1; i >= 0; i-- {
		dir := collectionDirs[i]
		filePath := filepath.Join(dir, collection)

		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}
	}

	return "", fmt.Errorf("collection %s not found in any configured paths", collection)
}

// GetConfigDir returns the config directory
func GetConfigDir() string {
	return configDir
}

// GetCollectionDirs returns the list of collection directories
func GetCollectionDirs() []string {
	return collectionDirs
}
