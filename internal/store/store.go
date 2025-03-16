package store

import (
	"fmt"
	"os"
	"strings"

	"github.com/sakshamgupta05/cmdvault/internal/config"
	"gopkg.in/yaml.v3"
)

// Command represents a stored shell command
type Command struct {
	Collection  string
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Command     string   `yaml:"command"`
	Tags        []string `yaml:"tags"`
}

type Collection struct {
	Commands []Command `yaml:"commands"`
}

func ListCollections() ([]string, error) {
	// List all collections by finding all YAML files in the commands directory
	collections := make([]string, 0)
	dirs := config.GetCollectionDirs()
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading commands directory: %v\n", err)
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
				collections = append(collections, strings.TrimSuffix(entry.Name(), ".yaml"))
			}
		}
	}

	return collections, nil
}

func GetAllCommands() ([]Command, error) {
	collections, err := ListCollections()
	if err != nil {
		return nil, err
	}

	var commands []Command
	for _, collection := range collections {
		cmds, err := GetCommands(collection)
		if err != nil {
			return nil, err
		}

		commands = append(commands, cmds...)
	}

	return commands, nil
}

// GetCommands returns all commands from the specified collection
func GetCommands(collection string) ([]Command, error) {
	collectionPath, err := config.GetCollectionPath(collection)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(collectionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %s: %v\n", collectionPath, err)
		return nil, err
	}

	var collectionStruct Collection
	if err := yaml.Unmarshal(data, &collectionStruct); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not parse file %s: %v\n", collectionPath, err)
		return nil, err
	}
	for i := range collectionStruct.Commands {
		collectionStruct.Commands[i].Collection = collection
	}

	return collectionStruct.Commands, nil
}

// SearchCommands searches for commands matching the search term
func SearchCommands(searchTerm string) ([]Command, error) {
	commands, err := GetAllCommands()
	if err != nil {
		return nil, err
	}

	if searchTerm == "" {
		return commands, nil
	}

	searchTerm = strings.ToLower(searchTerm)
	var results []Command
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd.Collection)+" "+strings.ToLower(cmd.Name), searchTerm) {
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
