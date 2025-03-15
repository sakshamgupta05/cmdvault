package cmd

import (
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/sakshamgupta05/cmdvault/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdvault",
	Short: "Store and retrieve shell commands",
	Long: `CmdVault helps you store and retrieve commonly used shell commands.
Store commands with descriptions and tags, then search and use them when needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, run interactive search
		if len(args) == 0 {
			ui.InteractiveSearch("")
		}
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(collectionsCmd)
}
