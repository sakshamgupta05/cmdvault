package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
		desc := cmd.Description
		if desc == "" {
			desc = "(No description)"
		}
		fmt.Printf("%s %s\n   %s\n\n", bold("•"), green(cmd.Name), desc)
	}
}

// Command represents a command for UI display
type Command struct {
	store.Command
}

// FilterValue implements list.Item for Command
func (c Command) FilterValue() string {
	return c.Name
}

// Title returns the title for list display
func (c Command) Title() string {
	return c.Name
}

// Description returns the description for list display
func (c Command) Description() string {
	if c.Command.Description == "" {
		return "(No description)"
	}
	return c.Command.Description
}

type model struct {
	list         list.Model
	commands     []Command
	selectedCmd  *Command
	showingState string // "selection", "action", "details", "executed", "copied"
	quitting     bool
}

func initialModel(commands []store.Command) model {
	// Convert to UI commands
	items := make([]list.Item, len(commands))
	uiCommands := make([]Command, len(commands))

	for i, cmd := range commands {
		uiCommands[i] = Command{cmd}
		items[i] = uiCommands[i]
	}

	// Create list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Commands"

	return model{
		list:         l,
		commands:     uiCommands,
		showingState: "selection",
	}
}

// Actions for the selected command
type actionModel struct {
	list     list.Model
	command  *Command
	quitting bool
	executed bool
}

func newActionModel(cmd *Command) actionModel {
	actions := []list.Item{
		actionItem{"Copy to clipboard"},
		actionItem{"Show details"},
		actionItem{"Execute command"},
	}

	l := list.New(actions, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Action"

	return actionModel{
		list:    l,
		command: cmd,
	}
}

type actionItem struct {
	title string
}

func (a actionItem) Title() string       { return a.title }
func (a actionItem) Description() string { return "" }
func (a actionItem) FilterValue() string { return a.title }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.showingState == "selection" && m.list.SelectedItem() != nil {
				selectedIndex := m.list.Index()
				m.selectedCmd = &m.commands[selectedIndex]
				m.showingState = "action"

				actionModel := newActionModel(m.selectedCmd)
				return actionModel, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	return m.list.View()
}

func (m actionModel) Init() tea.Cmd {
	return nil
}

func (m actionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			selected := m.list.SelectedItem().(actionItem).Title()

			switch selected {
			case "Copy to clipboard":
				if err := clipboard.WriteAll(m.command.Command); err != nil {
					fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\n", err)
				} else {
					fmt.Println(green("Command copied to clipboard!"))
				}
				return m, tea.Quit

			case "Show details":
				fmt.Println()
				fmt.Printf("%s %s\n\n", bold("Command:"), cyan(m.command.Command))
				fmt.Printf("%s %s\n\n", bold("Description:"), m.command.Description)
				fmt.Printf("%s %s\n", bold("Tags:"), strings.Join(m.command.Tags, ", "))
				return m, tea.Quit

			case "Execute command":
				fmt.Printf("%s %s\n", yellow("Executing:"), m.command.Command)

				// Create shell command
				var cmd *exec.Cmd
				if isWindows() {
					cmd = exec.Command("cmd", "/C", m.command.Command)
				} else {
					cmd = exec.Command("sh", "-c", m.command.Command)
				}

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin

				// We'll execute after quitting the TUI
				m.executed = true
				m.quitting = true

				go func() {
					if err := cmd.Run(); err != nil {
						fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
					}
				}()

				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m actionModel) View() string {
	if m.quitting {
		return ""
	}

	return m.list.View()
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

	m := initialModel(commands)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running UI: %v", err)
	}
}

// isWindows determines if the current OS is Windows
func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
