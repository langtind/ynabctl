# ynabctl

A command-line interface for [You Need A Budget (YNAB)](https://www.ynab.com/).

## Installation

### Homebrew (macOS/Linux)

```bash
brew install langtind/tap/ynabctl
```

### From Source

```bash
go install github.com/langtind/ynabctl@latest
```

### Manual Download

Download the latest release from the [releases page](https://github.com/langtind/ynabctl/releases).

## Getting Started

### 1. Get your YNAB API Token

1. Go to the [YNAB web app](https://app.ynab.com/)
2. Click on your account name (top left)
3. Go to "Account Settings"
4. Click on "Developer Settings"
5. Create a new Personal Access Token

### 2. Configure ynabctl

```bash
ynabctl config set-token <your-token>
```

### 3. Set a default budget (optional)

```bash
# List available budgets
ynabctl budgets list

# Set a default budget
ynabctl config set-default-budget <budget-id>
```

## Usage

### Configuration

```bash
# Show current configuration
ynabctl config show

# Set API token
ynabctl config set-token <token>

# Set default budget
ynabctl config set-default-budget <budget-id>

# Set default output format
ynabctl config set-format <json|table>
```

### Budgets

```bash
# List all budgets
ynabctl budgets list

# Get budget details
ynabctl budgets get [budget-id]

# Get budget settings
ynabctl budgets settings [budget-id]
```

### Accounts

```bash
# List all accounts
ynabctl accounts list

# Get account details
ynabctl accounts get <account-id>

# Create a new account
ynabctl accounts create --name "Checking" --type checking --balance 1000.00
```

### Categories

```bash
# List all categories
ynabctl categories list

# Get category details
ynabctl categories get <category-id>

# Update category budget
ynabctl categories update <category-id> --budgeted 500.00 --month 2024-01-01
```

### Transactions

```bash
# List transactions
ynabctl transactions list
ynabctl transactions list --since 2024-01-01
ynabctl transactions list --account <account-id>
ynabctl transactions list --category <category-id>

# Get transaction details
ynabctl transactions get <transaction-id>

# Create a transaction
ynabctl transactions create --account <account-id> --amount -50.00 --payee-name "Coffee Shop" --memo "Morning coffee"

# Update a transaction
ynabctl transactions update <transaction-id> --amount -55.00

# Delete a transaction
ynabctl transactions delete <transaction-id>
```

### Payees

```bash
# List all payees
ynabctl payees list

# Get payee details
ynabctl payees get <payee-id>

# Rename a payee
ynabctl payees update <payee-id> --name "New Name"
```

### Scheduled Transactions

```bash
# List scheduled transactions
ynabctl scheduled list

# Get scheduled transaction details
ynabctl scheduled get <scheduled-transaction-id>

# Create a scheduled transaction
ynabctl scheduled create --account <account-id> --amount -100.00 --frequency monthly --date-first 2024-01-01

# Update a scheduled transaction
ynabctl scheduled update <scheduled-transaction-id> --amount -150.00

# Delete a scheduled transaction
ynabctl scheduled delete <scheduled-transaction-id>
```

### Months

```bash
# List budget months
ynabctl months list

# Get month details
ynabctl months get 2024-01-01
ynabctl months get current
```

### User

```bash
# Get authenticated user info
ynabctl user
```

## Global Flags

```
--budget, -b    Budget ID to use (overrides default)
--format, -f    Output format (json, table)
```

## Configuration

Configuration is stored in `~/.config/ynabctl/config.toml`.

You can also use environment variables:
- `YNAB_TOKEN` - API token
- `YNAB_DEFAULT_BUDGET` - Default budget ID
- `YNAB_FORMAT` - Default output format

## Currency

YNAB uses milliunits internally (1000 = $1.00). This CLI automatically converts between regular currency amounts and milliunits for display and input.

## License

MIT
