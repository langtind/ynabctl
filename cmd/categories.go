package cmd

import (
	"fmt"
	"time"

	"github.com/langtind/ynabctl/internal/client"
	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage categories",
	Long:  `List, view, and update budget categories.`,
}

var categoriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all categories",
	Long:  `Returns a list of all category groups and categories for the budget.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := getBudgetID()
		if err != nil {
			return err
		}

		categories, err := apiClient.GetCategories(id)
		if err != nil {
			return fmt.Errorf("failed to get categories: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(categories)
	},
}

var categoriesGetCmd = &cobra.Command{
	Use:   "get <category-id>",
	Short: "Get category details",
	Long:  `Returns details for a specific category.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		category, err := apiClient.GetCategory(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get category: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(category)
	},
}

var (
	categoryMonth    string
	categoryBudgeted float64
)

var categoriesUpdateCmd = &cobra.Command{
	Use:   "update <category-id>",
	Short: "Update category budgeted amount",
	Long: `Update the budgeted amount for a category in a specific month.

The month should be in YYYY-MM-DD format (first day of the month) or "current" for the current month.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		month := categoryMonth
		if month == "" || month == "current" {
			month = time.Now().Format("2006-01-01")
		}

		budgeted := client.AmountToMilliunits(categoryBudgeted)

		category, err := apiClient.UpdateCategory(budgetID, args[0], month, budgeted)
		if err != nil {
			return fmt.Errorf("failed to update category: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(category)
	},
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
	categoriesCmd.AddCommand(categoriesListCmd)
	categoriesCmd.AddCommand(categoriesGetCmd)
	categoriesCmd.AddCommand(categoriesUpdateCmd)

	categoriesUpdateCmd.Flags().StringVar(&categoryMonth, "month", "current", "Budget month (YYYY-MM-DD or 'current')")
	categoriesUpdateCmd.Flags().Float64Var(&categoryBudgeted, "budgeted", 0, "Budgeted amount")
}
