package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/store"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [path]",
	Short: "Import collections from a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		if err := store.ImportCommands(path); err != nil {
			fmt.Printf("Error importing commands: %v\n", err)
			return
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Commands imported successfully!\n", green("✓"))
	},
}
