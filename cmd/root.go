package cmd

import (
	"fmt"
	"os"

	"github.com/langtind/ynabctl/internal/client"
	"github.com/langtind/ynabctl/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	outputFormat string
	budgetID     string

	// Shared client instance
	apiClient *client.Client

	// Config instance
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "ynabctl",
	Short: "A CLI tool for interacting with the YNAB API",
	Long: `ynabctl is a command-line interface for You Need A Budget (YNAB).

It allows you to manage budgets, accounts, transactions, categories,
and more directly from your terminal.

To get started, set your YNAB API token:
  ynabctl config set-token <your-token>

You can obtain a token from YNAB: Account Settings > Developer Settings`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for commands that don't need it
		if cmd.Name() == "version" || cmd.Name() == "help" || cmd.Name() == "ai" {
			return nil
		}
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			// Allow config commands to run without full initialization
			return nil
		}

		// Load configuration
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Set output format from config if not specified via flag
		if outputFormat == "" {
			outputFormat = cfg.Format
		}
		if outputFormat == "" {
			outputFormat = "json"
		}

		// Set budget ID from config if not specified via flag
		if budgetID == "" {
			budgetID = cfg.DefaultBudget
		}

		// Initialize API client for commands that need it
		if requiresAuth(cmd) {
			if cfg.Token == "" {
				return fmt.Errorf("YNAB API token not configured. Run 'ynabctl config set-token <token>' to set it")
			}
			apiClient = client.New(cfg.Token)
		}

		return nil
	},
}

// requiresAuth returns true if the command needs API authentication
func requiresAuth(cmd *cobra.Command) bool {
	// Config commands don't need auth
	if cmd.Name() == "show" && cmd.Parent() != nil && cmd.Parent().Name() == "config" {
		return false
	}
	if cmd.Name() == "set-token" || cmd.Name() == "set-default-budget" {
		return false
	}
	return true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "", "Output format (json, table)")
	rootCmd.PersistentFlags().StringVarP(&budgetID, "budget", "b", "", "Budget ID to use")
}

// getBudgetID returns the budget ID to use, checking flag first, then config default
func getBudgetID() (string, error) {
	if budgetID != "" {
		return budgetID, nil
	}
	if cfg != nil && cfg.DefaultBudget != "" {
		return cfg.DefaultBudget, nil
	}
	return "", fmt.Errorf("no budget specified. Use --budget flag or set a default with 'ynabctl config set-default-budget <id>'")
}

// getOutputFormat returns the output format to use
func getOutputFormat() string {
	if outputFormat != "" {
		return outputFormat
	}
	return "json"
}
