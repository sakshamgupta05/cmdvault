package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/store"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [path]",
	Short: "Export all collections to a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		if err := store.ExportCommands(path); err != nil {
			fmt.Printf("Error exporting commands: %v\n", err)
			return
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Commands exported to %s\n", green("✓"), path)
	},
}
