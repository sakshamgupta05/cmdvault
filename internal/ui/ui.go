package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/sakshamgupta05/cmdvault/internal/store"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

// AddCommandPrompt prompts the user to add a new command
func AddCommandPrompt(collection string) {
	// var cmd store.Command

	// questions := []*survey.Question{
	// 	{
	// 		Name: "name",
	// 		Prompt: &survey.Input{
	// 			Message: "Command name:",
	// 		},
	// 		Validate: survey.Required,
	// 	},
	// 	{
	// 		Name: "description",
	// 		Prompt: &survey.Input{
	// 			Message: "Description (optional):",
	// 		},
	// 	},
	// 	{
	// 		Name: "command",
	// 		Prompt: &survey.Input{
	// 			Message: "Command:",
	// 		},
	// 		Validate: survey.Required,
	// 	},
	// 	{
	// 		Name: "tags",
	// 		Prompt: &survey.Input{
	// 			Message: "Tags (comma separated, optional):",
	// 		},
	// 	},
	// }
	cmd := store.Command{
		Name:        "name",
		Description: "description",
		Command:     "command",
		Tags:        []string{"tags"},
	}

	// err := survey.Ask(questions, &cmd)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 	return
	// }

	// Process tags
	// if cmd.Tags == nil {
	cmd.Tags = []string{}
	// } else if len(cmd.Tags) == 1 {
	// 	tagString := cmd.Tags[0]
	// 	cmd.Tags = []string{}
	// 	for _, tag := range strings.Split(tagString, ",") {
	// 		tag = strings.TrimSpace(tag)
	// 		if tag != "" {
	// 			cmd.Tags = append(cmd.Tags, tag)
	// 		}
	// 	}
	// }

	if err := store.SaveCommand(cmd, collection); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving command: %v\n", err)
		return
	}

	fmt.Printf("%s Command \"%s\" added successfully!\n", green("✓"), cmd.Name)
}

// ListCommands lists all commands in a collection
func ListCommands(collection string) {
	commands, err := store.GetCommands(collection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if len(commands) == 0 {
		fmt.Println(yellow("No commands found."))
		return
	}

	fmt.Printf("%s Commands in collection \"%s\":\n\n", bold("•"), collection)

	for _, cmd := range commands {
		desc := cmd.Description
		if desc == "" {
			desc = "(No description)"
		}
		fmt.Printf("%s %s\n   %s\n\n", bold("•"), green(cmd.Name), desc)
	}
}

// InteractiveSearch provides an interactive search interface
func InteractiveSearch(searchTerm string) {
	collection := config.GetDefaultCollection()
	commands, err := store.SearchCommands(searchTerm, collection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if len(commands) == 0 {
		fmt.Println(yellow("No commands found."))
		return
	}

	// Prepare options
	options := make([]string, len(commands))
	for i, cmd := range commands {
		desc := cmd.Description
		if desc == "" {
			desc = "(No description)"
		}
		options[i] = fmt.Sprintf("%s - %s", cmd.Name, desc)
	}

	// Prompt for command selection
	var selectedIndex int
	prompt := &survey.Select{
		Message: "Select a command:",
		Options: options,
	}
	survey.AskOne(prompt, &selectedIndex)

	selected := commands[selectedIndex]

	// Prompt for action
	var action string
	actionPrompt := &survey.Select{
		Message: "What would you like to do?",
		Options: []string{"Copy to clipboard", "Show details", "Execute command"},
	}
	survey.AskOne(actionPrompt, &action)

	switch action {
	case "Copy to clipboard":
		if err := clipboard.WriteAll(selected.Command); err != nil {
			fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\n", err)
			return
		}
		fmt.Println(green("Command copied to clipboard!"))

	case "Show details":
		fmt.Println()
		fmt.Printf("%s %s\n\n", bold("Command:"), cyan(selected.Command))
		fmt.Printf("%s %s\n\n", bold("Description:"), selected.Description)
		fmt.Printf("%s %s\n", bold("Tags:"), strings.Join(selected.Tags, ", "))

	case "Execute command":
		fmt.Printf("%s %s\n", yellow("Executing:"), selected.Command)

		// Create shell command
		var cmd *exec.Cmd
		if isWindows() {
			cmd = exec.Command("cmd", "/C", selected.Command)
		} else {
			cmd = exec.Command("sh", "-c", selected.Command)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		}
	}
}

// isWindows determines if the current OS is Windows
func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
