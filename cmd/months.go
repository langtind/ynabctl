package cmd

import (
	"fmt"
	"time"

	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var monthsCmd = &cobra.Command{
	Use:   "months",
	Short: "Manage budget months",
	Long:  `List and view budget month information.`,
}

var monthsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List budget months",
	Long:  `Returns a list of all budget months.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		months, err := apiClient.GetMonths(budgetID)
		if err != nil {
			return fmt.Errorf("failed to get months: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(months)
	},
}

var monthsGetCmd = &cobra.Command{
	Use:   "get [month]",
	Short: "Get budget month details",
	Long: `Returns details for a specific budget month.

The month should be in YYYY-MM-DD format (first day of the month) or "current" for the current month.
If no month is specified, returns the current month.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		month := "current"
		if len(args) > 0 {
			month = args[0]
		}

		if month == "current" {
			// Use first day of current month
			now := time.Now()
			month = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		}

		monthData, err := apiClient.GetMonth(budgetID, month)
		if err != nil {
			return fmt.Errorf("failed to get month: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(monthData)
	},
}

func init() {
	rootCmd.AddCommand(monthsCmd)
	monthsCmd.AddCommand(monthsListCmd)
	monthsCmd.AddCommand(monthsGetCmd)
}
