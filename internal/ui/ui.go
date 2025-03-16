package ui

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

	// Use fuzzyfinder to select a command

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
		Options: []string{"Execute command", "Copy to clipboard", "Show details"},
	}
	survey.AskOne(actionPrompt, &action)

	switch action {
	case "Execute command":
		if len(selected.Parameters) > 0 {
			selected.Command = interactiveParameters(selected)
		}
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

	case "Copy to clipboard":
		if len(selected.Parameters) > 0 {
			selected.Command = interactiveParameters(selected)
		}

		if err := clipboard.WriteAll(selected.Command); err != nil {
			fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\n", err)
			return
		}
		fmt.Println(green("Command copied to clipboard!"))

	case "Show details":
		fmt.Printf("\n%s\n", formatCommandLong(selected))
	}
}

func interactiveParameters(cmd store.Command) string {
	cmdStr := cmd.Command

	for _, param := range cmd.Parameters {
		defaultValueStr := ""
		if param.DefaultValue != "" {
			defaultValueStr = fmt.Sprintf(" (%s)", param.DefaultValue)
		}

		mandatoryStr := ""
		if !param.Optional && param.DefaultValue == "" {
			mandatoryStr = bold("*")
		}
		promptStr := fmt.Sprintf("%s%s%s:", param.Name, mandatoryStr, defaultValueStr)
		prompt := &survey.Input{Message: promptStr}

		var value string
		var err error

		if !param.Optional && param.DefaultValue == "" {
			// Required parameter with no default value
			err = survey.AskOne(prompt, &value, survey.WithValidator(survey.Required))
		} else {
			err = survey.AskOne(prompt, &value)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		cmdStr = replaceParameter(cmdStr, param, value)
	}

	return cmdStr
}

func replaceParameter(cmdStr string, p store.Parameter, value string) string {
	if p.Optional {
		// Handle optional parameter with regex
		placeholder := fmt.Sprintf("{{%s}}", p.Name)
		// Match {? ... {{param}} ... ?}
		re := regexp.MustCompile(`\{\?(?:[^?]|\?(?:[^}]|$))*\{\{` + p.Name + `\}\}(?:[^?]|\?(?:[^}]|$))*\?\}`)

		if value == "" {
			// If value is empty, remove the entire optional block
			return re.ReplaceAllString(cmdStr, "")
		} else {
			// Replace parameter and remove the {? ?} markers
			matches := re.FindAllStringSubmatch(cmdStr, -1)
			result := cmdStr

			for _, match := range matches {
				if len(match) > 0 {
					originalText := match[0]
					newText := strings.Replace(originalText, "{?", "", 1)
					newText = strings.Replace(newText, "?}", "", 1)
					newText = strings.Replace(newText, placeholder, value, -1)
					result = strings.Replace(result, originalText, newText, -1)
				}
			}

			return result
		}
	} else {
		if value == "" {
			value = p.DefaultValue
		}
		placeholder := fmt.Sprintf("{{%s}}", p.Name)
		return strings.Replace(cmdStr, placeholder, value, -1)
	}
}

func formatCommandLong(cmd store.Command) string {
	tagsText := ""
	if len(cmd.Tags) > 0 {
		tagsText = fmt.Sprintf("\n   [%s]", strings.Join(cmd.Tags, ", "))
	}

	descriptionText := ""
	if cmd.Description != "" {
		descriptionText = fmt.Sprintf("\n   %s", cmd.Description)
	}

	parametersText := ""
	if len(cmd.Parameters) > 0 {
		parametersText = fmt.Sprintf("\n\n%s", bold("Parameters:"))
		for _, param := range cmd.Parameters {
			parametersText += fmt.Sprintf("\n   %s: %s", param.Name, param.Description)
		}
	}

	return fmt.Sprintf("%s\n   %s %s%s%s\n\n%s\n   %s%s",
		bold("Description:"),
		cmd.Collection,
		yellow(cmd.Name),
		descriptionText,
		tagsText,
		bold("Command:"),
		cyan(cmd.Command),
		parametersText)
}

// isWindows determines if the current OS is Windows
func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
