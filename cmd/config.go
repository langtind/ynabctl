package cmd

import (
	"fmt"

	"github.com/langtind/ynabctl/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage ynabctl configuration",
	Long:  `View and modify ynabctl configuration settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Printf("Config file: %s\n\n", config.GetConfigFile())

		// Mask the token for security
		token := cfg.Token
		if token != "" {
			if len(token) > 8 {
				token = token[:4] + "..." + token[len(token)-4:]
			} else {
				token = "****"
			}
		} else {
			token = "(not set)"
		}

		fmt.Printf("Token:          %s\n", token)
		fmt.Printf("Default Budget: %s\n", valueOrNotSet(cfg.DefaultBudget))
		fmt.Printf("Format:         %s\n", valueOrNotSet(cfg.Format))

		return nil
	},
}

var configSetTokenCmd = &cobra.Command{
	Use:   "set-token <token>",
	Short: "Set the YNAB API token",
	Long: `Set the YNAB API token for authentication.

You can obtain a token from YNAB:
  1. Go to YNAB web app
  2. Click on your account name (top left)
  3. Go to "Account Settings"
  4. Click on "Developer Settings"
  5. Create a new Personal Access Token`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]
		if err := config.SetToken(token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
		fmt.Println("Token saved successfully.")
		return nil
	},
}

var configSetDefaultBudgetCmd = &cobra.Command{
	Use:   "set-default-budget <budget-id>",
	Short: "Set the default budget ID",
	Long: `Set the default budget ID to use for commands.

This budget will be used when the --budget flag is not specified.
You can find budget IDs by running: ynabctl budgets list`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID := args[0]
		if err := config.SetDefaultBudget(budgetID); err != nil {
			return fmt.Errorf("failed to save default budget: %w", err)
		}
		fmt.Printf("Default budget set to: %s\n", budgetID)
		return nil
	},
}

var configSetFormatCmd = &cobra.Command{
	Use:   "set-format <format>",
	Short: "Set the default output format",
	Long: `Set the default output format (json or table).

This format will be used when the --format flag is not specified.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format := args[0]
		if format != "json" && format != "table" {
			return fmt.Errorf("invalid format: %s (must be 'json' or 'table')", format)
		}
		if err := config.SetFormat(format); err != nil {
			return fmt.Errorf("failed to save format: %w", err)
		}
		fmt.Printf("Default format set to: %s\n", format)
		return nil
	},
}

func valueOrNotSet(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetTokenCmd)
	configCmd.AddCommand(configSetDefaultBudgetCmd)
	configCmd.AddCommand(configSetFormatCmd)
}
