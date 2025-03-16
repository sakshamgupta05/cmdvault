package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/store"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

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
		fmt.Printf("%s %s: %s\n   %s\n\n", bold("•"), cmd.Collection, green(cmd.Name), strings.Join(cmd.Tags, ", "))
	}
}

// InteractiveSearch provides an interactive search interface
func InteractiveSearch(searchTerm string) {
	commands, err := store.SearchCommands(searchTerm)
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
		options[i] = fmt.Sprintf("%s: %s - %s", cmd.Collection, cmd.Name, strings.Join(cmd.Tags, ", "))
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
