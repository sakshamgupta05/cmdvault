package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sakshamgupta05/cmdvault/internal/store"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
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

	//// fzf

	selectedIndex, err := fuzzyfinder.Find(
		commands,
		func(i int) string {
			tagsText := ""
			if len(commands[i].Tags) > 0 {
				tagsText = fmt.Sprintf(" [%s]", strings.Join(commands[i].Tags, ", "))
			}
			return fmt.Sprintf("%s: %s%s", commands[i].Collection, commands[i].Name, tagsText)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return formatCommandLong(commands[i])
		}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		return
	}

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
		fmt.Printf("\n%s\n", formatCommandLong(selected))

	case "Execute command":
		fmt.Printf("%s %s\n", yellow("Executing:"), selected.Command)

		// Create shell command
		var cmd *exec.Cmd
		if isWindows() {
			cmd = exec.Command("cmd", "/C", selected.Command)
		} else {
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "sh"
			}
			cmd = exec.Command(shell, "-c", selected.Command)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			} else {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", red(err))
				os.Exit(1)
			}
		}
	}
}

func formatCommandLong(cmd store.Command) string {
	tagsText := ""
	if len(cmd.Tags) > 0 {
		tagsText = fmt.Sprintf("\n   [%s]", strings.Join(cmd.Tags, ", "))
	}
	return fmt.Sprintf("%s\n   %s %s%s\n\n%s\n   %s",
		bold("Description:"),
		cmd.Collection,
		yellow(cmd.Name),
		tagsText,
		bold("Command:"),
		cyan(cmd.Command))
}

// isWindows determines if the current OS is Windows
func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
