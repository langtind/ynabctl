package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "Output context for AI assistants",
	Long:  `Prints documentation and examples to help AI assistants use ynabctl effectively.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(aiContext)
	},
}

func init() {
	rootCmd.AddCommand(aiCmd)
}

const aiContext = `# ynabctl - AI Assistant Context

## Overview

CLI for interacting with You Need A Budget (YNAB) API. This document helps AI assistants use ynabctl effectively.

## IMPORTANT: YNAB Concepts

### Budgets
A budget is a container for all financial data. Users typically have one budget, but may have multiple (personal, business, etc.). Most commands require a budget ID.

### Accounts
Bank accounts, credit cards, cash, etc. Types: checking, savings, cash, creditCard, lineOfCredit, mortgage, autoLoan, studentLoan, etc.

### Categories
Budget categories organized in groups (e.g., "Bills" group containing "Rent", "Utilities"). Each category has:
- Budgeted: Amount assigned this month
- Activity: Spending/income this month
- Balance: Available amount

### Transactions
Individual financial transactions with date, amount, payee, category, and memo.

### Milliunits
**CRITICAL**: YNAB API uses milliunits internally (1000 = $1.00).
- $50.00 = 50000 milliunits
- -$25.50 = -25500 milliunits
- ynabctl handles conversion automatically for display and input

---

## Setup

` + "```bash" + `
# Set API token (get from YNAB: Account Settings > Developer Settings)
ynabctl config set-token <token>

# Set default budget (optional, avoids --budget flag)
ynabctl budgets list                           # Find budget ID
ynabctl config set-default-budget <budget-id>

# Set default output format
ynabctl config set-format table                # or "json" (default)

# View current config
ynabctl config show
` + "```" + `

---

## Quick Reference

### Configuration

` + "```bash" + `
ynabctl config show                            # Show current config
ynabctl config set-token <token>               # Set API token
ynabctl config set-default-budget <id>         # Set default budget
ynabctl config set-format <json|table>         # Set output format
` + "```" + `

### Budgets

` + "```bash" + `
ynabctl budgets list                           # List all budgets
ynabctl budgets get                            # Get default budget details
ynabctl budgets get <budget-id>                # Get specific budget
ynabctl budgets settings                       # Get budget settings (currency, date format)
` + "```" + `

### Accounts

` + "```bash" + `
ynabctl accounts list                          # List all accounts
ynabctl accounts get <account-id>              # Get account details
ynabctl accounts create --name "Checking" --type checking --balance 1000.00
` + "```" + `

Account types: checking, savings, cash, creditCard, lineOfCredit, otherAsset, otherLiability, mortgage, autoLoan, studentLoan, personalLoan, medicalDebt, otherDebt

### Categories

` + "```bash" + `
ynabctl categories list                        # List all category groups and categories
ynabctl categories get <category-id>           # Get category details
ynabctl categories update <id> --budgeted 500  # Update budgeted amount
ynabctl categories update <id> --budgeted 500 --month 2024-01-01
` + "```" + `

### Transactions

` + "```bash" + `
# List transactions
ynabctl transactions list                      # All transactions
ynabctl transactions list --since 2024-01-01   # Since date
ynabctl transactions list --account <id>       # By account
ynabctl transactions list --category <id>      # By category
ynabctl transactions list --payee <id>         # By payee
ynabctl transactions list --type unapproved    # Unapproved only
ynabctl transactions list --type uncategorized # Uncategorized only

# Get single transaction
ynabctl transactions get <transaction-id>

# Create transaction
ynabctl transactions create \
  --account <account-id> \
  --amount -50.00 \
  --payee-name "Coffee Shop" \
  --category <category-id> \
  --memo "Morning coffee" \
  --date 2024-01-15

# Update transaction
ynabctl transactions update <id> --amount -55.00
ynabctl transactions update <id> --memo "Updated memo"
ynabctl transactions update <id> --category <new-category-id>

# Delete transaction
ynabctl transactions delete <transaction-id>
` + "```" + `

**Amount convention**: Negative = outflow (spending), Positive = inflow (income)

### Payees

` + "```bash" + `
ynabctl payees list                            # List all payees
ynabctl payees get <payee-id>                  # Get payee details
ynabctl payees update <id> --name "New Name"   # Rename payee
` + "```" + `

### Scheduled Transactions

` + "```bash" + `
ynabctl scheduled list                         # List scheduled transactions
ynabctl scheduled get <id>                     # Get details

# Create scheduled transaction
ynabctl scheduled create \
  --account <account-id> \
  --amount -100.00 \
  --frequency monthly \
  --date-first 2024-01-01 \
  --payee-name "Landlord" \
  --memo "Rent"

# Update
ynabctl scheduled update <id> --amount -150.00

# Delete
ynabctl scheduled delete <id>
` + "```" + `

Frequency options: never, daily, weekly, everyOtherWeek, twiceAMonth, every4Weeks, monthly, everyOtherMonth, every3Months, every4Months, twiceAYear, yearly, everyOtherYear

### Budget Months

` + "```bash" + `
ynabctl months list                            # List all budget months
ynabctl months get current                     # Current month details
ynabctl months get 2024-01-01                  # Specific month
` + "```" + `

Month response includes: income, budgeted, activity, to_be_budgeted, age_of_money

### User

` + "```bash" + `
ynabctl user                                   # Get authenticated user info
` + "```" + `

---

## Global Flags

` + "```bash" + `
--budget, -b <id>     # Use specific budget (overrides default)
--format, -f <fmt>    # Output format: json (default) or table
` + "```" + `

---

## Output Formats

### JSON (default)
Best for parsing and scripting:
` + "```bash" + `
ynabctl accounts list | jq '.[].name'
ynabctl transactions list | jq '.[] | select(.amount < 0) | .payee_name'
` + "```" + `

### Table
Human-readable:
` + "```bash" + `
ynabctl accounts list -f table
ynabctl transactions list -f table --since 2024-01-01
` + "```" + `

---

## Common Workflows

### Check Budget Status
` + "```bash" + `
ynabctl months get current -f table            # See to_be_budgeted
ynabctl categories list -f table               # See category balances
` + "```" + `

### Record a Purchase
` + "```bash" + `
# 1. Find account and category IDs
ynabctl accounts list | jq '.[] | {id, name}'
ynabctl categories list | jq '.[].categories[] | {id, name}'

# 2. Create transaction
ynabctl transactions create \
  --account <account-id> \
  --amount -42.50 \
  --payee-name "Amazon" \
  --category <category-id>
` + "```" + `

### Review Spending
` + "```bash" + `
# This month's transactions
ynabctl transactions list --since $(date +%Y-%m-01) -f table

# By category
ynabctl transactions list --category <id> -f table

# Unapproved transactions needing review
ynabctl transactions list --type unapproved -f table
` + "```" + `

### Update Category Budget
` + "```bash" + `
# Increase grocery budget for current month
ynabctl categories update <grocery-category-id> --budgeted 600
` + "```" + `

---

## Environment Variables

` + "```bash" + `
YNAB_TOKEN           # API token (alternative to config file)
YNAB_DEFAULT_BUDGET  # Default budget ID
YNAB_FORMAT          # Default output format
` + "```" + `

---

## Error Handling

Common errors:
- **401 Unauthorized**: Invalid or expired token
- **404 Not Found**: Invalid budget/account/transaction ID
- **400 Bad Request**: Invalid parameters (check date format, amount, etc.)

---

## Tips

1. **Always set a default budget** to avoid specifying --budget on every command
2. **Use jq for JSON parsing** when scripting
3. **Amounts are in regular currency** - ynabctl handles milliunit conversion
4. **Date format is YYYY-MM-DD** for all date parameters
5. **IDs are UUIDs** - copy them exactly from list commands

---

## Help

` + "```bash" + `
ynabctl --help                                 # All commands
ynabctl <command> --help                       # Command help
ynabctl version                                # Version info
` + "```" + `
`
