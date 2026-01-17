package cmd

import (
	"fmt"
	"time"

	"github.com/langtind/ynabctl/internal/client"
	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Manage transactions",
	Long:  `List, view, create, update, and delete transactions.`,
}

var (
	txnSinceDate  string
	txnType       string
	txnAccountID  string
	txnCategoryID string
	txnPayeeID    string
)

var transactionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List transactions",
	Long: `Returns a list of transactions for the budget.

Use filters to narrow down results:
  --since: Only return transactions on or after this date (YYYY-MM-DD)
  --type: Filter by transaction type (uncategorized, unapproved)
  --account: Filter by account ID
  --category: Filter by category ID
  --payee: Filter by payee ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		var transactions []client.Transaction

		// Use specific endpoint if filtering by account, category, or payee
		if txnAccountID != "" {
			transactions, err = apiClient.GetTransactionsByAccount(budgetID, txnAccountID, txnSinceDate)
		} else if txnCategoryID != "" {
			transactions, err = apiClient.GetTransactionsByCategory(budgetID, txnCategoryID, txnSinceDate)
		} else if txnPayeeID != "" {
			transactions, err = apiClient.GetTransactionsByPayee(budgetID, txnPayeeID, txnSinceDate)
		} else {
			filter := &client.TransactionFilter{
				SinceDate: txnSinceDate,
				Type:      txnType,
			}
			transactions, err = apiClient.GetTransactions(budgetID, filter)
		}

		if err != nil {
			return fmt.Errorf("failed to get transactions: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transactions)
	},
}

var transactionsGetCmd = &cobra.Command{
	Use:   "get <transaction-id>",
	Short: "Get transaction details",
	Long:  `Returns details for a specific transaction.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		transaction, err := apiClient.GetTransaction(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

var (
	newTxnAccountID  string
	newTxnDate       string
	newTxnAmount     float64
	newTxnPayeeID    string
	newTxnPayeeName  string
	newTxnCategoryID string
	newTxnMemo       string
	newTxnCleared    string
	newTxnApproved   bool
	newTxnFlagColor  string
)

var transactionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new transaction",
	Long: `Create a new transaction in the budget.

Required flags:
  --account: Account ID
  --amount: Transaction amount (positive for inflow, negative for outflow)

Optional flags:
  --date: Transaction date (YYYY-MM-DD, default: today)
  --payee-id: Payee ID
  --payee-name: Payee name (creates new payee if needed)
  --category: Category ID
  --memo: Transaction memo
  --cleared: Cleared status (cleared, uncleared, reconciled)
  --approved: Whether the transaction is approved
  --flag: Flag color (red, orange, yellow, green, blue, purple)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		if newTxnAccountID == "" {
			return fmt.Errorf("account ID is required (--account)")
		}

		date := newTxnDate
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		txn := client.SaveTransaction{
			AccountID:  newTxnAccountID,
			Date:       date,
			Amount:     client.AmountToMilliunits(newTxnAmount),
			PayeeID:    newTxnPayeeID,
			PayeeName:  newTxnPayeeName,
			CategoryID: newTxnCategoryID,
			Memo:       newTxnMemo,
			Cleared:    newTxnCleared,
			Approved:   newTxnApproved,
			FlagColor:  newTxnFlagColor,
		}

		transaction, err := apiClient.CreateTransaction(budgetID, txn)
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

var transactionsUpdateCmd = &cobra.Command{
	Use:   "update <transaction-id>",
	Short: "Update a transaction",
	Long: `Update an existing transaction.

All update flags are optional. Only specified fields will be updated.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		// First get the existing transaction
		existing, err := apiClient.GetTransaction(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get existing transaction: %w", err)
		}

		// Build update with existing values, override with any provided flags
		txn := client.SaveTransaction{
			AccountID:  existing.AccountID,
			Date:       existing.Date,
			Amount:     existing.Amount,
			PayeeID:    existing.PayeeID,
			CategoryID: existing.CategoryID,
			Memo:       existing.Memo,
			Cleared:    existing.Cleared,
			Approved:   existing.Approved,
			FlagColor:  existing.FlagColor,
		}

		if cmd.Flags().Changed("account") {
			txn.AccountID = newTxnAccountID
		}
		if cmd.Flags().Changed("date") {
			txn.Date = newTxnDate
		}
		if cmd.Flags().Changed("amount") {
			txn.Amount = client.AmountToMilliunits(newTxnAmount)
		}
		if cmd.Flags().Changed("payee-id") {
			txn.PayeeID = newTxnPayeeID
		}
		if cmd.Flags().Changed("payee-name") {
			txn.PayeeName = newTxnPayeeName
		}
		if cmd.Flags().Changed("category") {
			txn.CategoryID = newTxnCategoryID
		}
		if cmd.Flags().Changed("memo") {
			txn.Memo = newTxnMemo
		}
		if cmd.Flags().Changed("cleared") {
			txn.Cleared = newTxnCleared
		}
		if cmd.Flags().Changed("approved") {
			txn.Approved = newTxnApproved
		}
		if cmd.Flags().Changed("flag") {
			txn.FlagColor = newTxnFlagColor
		}

		transaction, err := apiClient.UpdateTransaction(budgetID, args[0], txn)
		if err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

var transactionsDeleteCmd = &cobra.Command{
	Use:   "delete <transaction-id>",
	Short: "Delete a transaction",
	Long:  `Delete a transaction from the budget.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		transaction, err := apiClient.DeleteTransaction(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to delete transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

func init() {
	rootCmd.AddCommand(transactionsCmd)
	transactionsCmd.AddCommand(transactionsListCmd)
	transactionsCmd.AddCommand(transactionsGetCmd)
	transactionsCmd.AddCommand(transactionsCreateCmd)
	transactionsCmd.AddCommand(transactionsUpdateCmd)
	transactionsCmd.AddCommand(transactionsDeleteCmd)

	// List filters
	transactionsListCmd.Flags().StringVar(&txnSinceDate, "since", "", "Filter transactions since date (YYYY-MM-DD)")
	transactionsListCmd.Flags().StringVar(&txnType, "type", "", "Filter by type (uncategorized, unapproved)")
	transactionsListCmd.Flags().StringVar(&txnAccountID, "account", "", "Filter by account ID")
	transactionsListCmd.Flags().StringVar(&txnCategoryID, "category", "", "Filter by category ID")
	transactionsListCmd.Flags().StringVar(&txnPayeeID, "payee", "", "Filter by payee ID")

	// Create/Update flags
	transactionsCreateCmd.Flags().StringVar(&newTxnAccountID, "account", "", "Account ID (required)")
	transactionsCreateCmd.Flags().StringVar(&newTxnDate, "date", "", "Transaction date (YYYY-MM-DD)")
	transactionsCreateCmd.Flags().Float64Var(&newTxnAmount, "amount", 0, "Amount (positive=inflow, negative=outflow)")
	transactionsCreateCmd.Flags().StringVar(&newTxnPayeeID, "payee-id", "", "Payee ID")
	transactionsCreateCmd.Flags().StringVar(&newTxnPayeeName, "payee-name", "", "Payee name")
	transactionsCreateCmd.Flags().StringVar(&newTxnCategoryID, "category", "", "Category ID")
	transactionsCreateCmd.Flags().StringVar(&newTxnMemo, "memo", "", "Memo")
	transactionsCreateCmd.Flags().StringVar(&newTxnCleared, "cleared", "", "Cleared status")
	transactionsCreateCmd.Flags().BoolVar(&newTxnApproved, "approved", false, "Approved")
	transactionsCreateCmd.Flags().StringVar(&newTxnFlagColor, "flag", "", "Flag color")

	transactionsUpdateCmd.Flags().StringVar(&newTxnAccountID, "account", "", "Account ID")
	transactionsUpdateCmd.Flags().StringVar(&newTxnDate, "date", "", "Transaction date (YYYY-MM-DD)")
	transactionsUpdateCmd.Flags().Float64Var(&newTxnAmount, "amount", 0, "Amount")
	transactionsUpdateCmd.Flags().StringVar(&newTxnPayeeID, "payee-id", "", "Payee ID")
	transactionsUpdateCmd.Flags().StringVar(&newTxnPayeeName, "payee-name", "", "Payee name")
	transactionsUpdateCmd.Flags().StringVar(&newTxnCategoryID, "category", "", "Category ID")
	transactionsUpdateCmd.Flags().StringVar(&newTxnMemo, "memo", "", "Memo")
	transactionsUpdateCmd.Flags().StringVar(&newTxnCleared, "cleared", "", "Cleared status")
	transactionsUpdateCmd.Flags().BoolVar(&newTxnApproved, "approved", false, "Approved")
	transactionsUpdateCmd.Flags().StringVar(&newTxnFlagColor, "flag", "", "Flag color")
}
