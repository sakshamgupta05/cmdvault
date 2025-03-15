package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sakshamgupta05/cmdvault/internal/config"
	"gopkg.in/yaml.v3"
)

// Command represents a stored shell command
type Command struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Command     string   `yaml:"command"`
	Tags        []string `yaml:"tags"`
}

// SaveCommand saves a command to the specified collection
func SaveCommand(cmd Command, collection string) error {
	// Create safe filename from command name
	filename := strings.ToLower(cmd.Name)
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ReplaceAll(filename, "/", "-")
	filename += ".yaml"

	collectionPath := config.GetCollectionPath(collection)
	filePath := filepath.Join(collectionPath, filename)

	data, err := yaml.Marshal(cmd)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// GetCommands returns all commands from the specified collection
func GetCommands(collection string) ([]Command, error) {
	collectionPath := config.GetCollectionPath(collection)
	files, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, err
	}

	var commands []Command
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		filePath := filepath.Join(collectionPath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read file %s: %v\n", filePath, err)
			continue
		}

		var cmd Command
		if err := yaml.Unmarshal(data, &cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not parse file %s: %v\n", filePath, err)
			continue
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}

// SearchCommands searches for commands matching the search term
func SearchCommands(searchTerm, collection string) ([]Command, error) {
	commands, err := GetCommands(collection)
	if err != nil {
		return nil, err
	}

	if searchTerm == "" {
		return commands, nil
	}

	searchTerm = strings.ToLower(searchTerm)
	var results []Command
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd.Name), searchTerm) ||
			strings.Contains(strings.ToLower(cmd.Description), searchTerm) {
			results = append(results, cmd)
			continue
		}

		// Search tags
		for _, tag := range cmd.Tags {
			if strings.Contains(strings.ToLower(tag), searchTerm) {
				results = append(results, cmd)
				break
			}
		}
	}

	return results, nil
}

// ExportCommands exports all commands to a directory
func ExportCommands(exportDir string) error {
	cfg := config.GetConfig()

	// Create export directory
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return err
	}

	// Export config
	configData, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(exportDir, "config.yaml"), configData, 0644); err != nil {
		return err
	}

	// Export each collection
	for _, collection := range cfg.Collections {
		collectionDir := filepath.Join(exportDir, collection)
		if err := os.MkdirAll(collectionDir, 0755); err != nil {
			return err
		}

		sourcePath := config.GetCollectionPath(collection)
		files, err := os.ReadDir(sourcePath)
		if err != nil {
			// Skip if collection directory doesn't exist yet
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			sourceFile := filepath.Join(sourcePath, file.Name())
			destFile := filepath.Join(collectionDir, file.Name())

			data, err := os.ReadFile(sourceFile)
			if err != nil {
				return err
			}

			if err := os.WriteFile(destFile, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

// ImportCommands imports commands from a directory
func ImportCommands(importDir string) error {
	// Import config
	configPath := filepath.Join(importDir, "config.yaml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	var importedConfig config.Config
	if err := yaml.Unmarshal(configData, &importedConfig); err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// Add each collection from imported config
	for _, collection := range importedConfig.Collections {
		config.AddCollection(collection)

		// Import commands from this collection
		sourceDir := filepath.Join(importDir, collection)
		targetDir := config.GetCollectionPath(collection)

		// Skip if source directory doesn't exist
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			continue
		}

		files, err := os.ReadDir(sourceDir)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
				continue
			}

			sourceFile := filepath.Join(sourceDir, file.Name())
			targetFile := filepath.Join(targetDir, file.Name())

			data, err := os.ReadFile(sourceFile)
			if err != nil {
				return err
			}

			if err := os.WriteFile(targetFile, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
