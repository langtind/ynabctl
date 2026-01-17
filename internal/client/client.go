package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.ynab.com/v1"

// Client handles communication with the YNAB API
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

// New creates a new YNAB API client
func New(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token:   token,
		baseURL: baseURL,
	}
}

// Error represents a YNAB API error
type Error struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Detail)
}

// ErrorResponse wraps the error from YNAB API
type ErrorResponse struct {
	Error *Error `json:"error"`
}

// doRequest performs an HTTP request to the YNAB API
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != nil {
			return nil, errResp.Error
		}
		return nil, fmt.Errorf("API error: %s (status %d)", string(respBody), resp.StatusCode)
	}

	return respBody, nil
}

// User types
type User struct {
	ID string `json:"id"`
}

type UserResponse struct {
	Data struct {
		User User `json:"user"`
	} `json:"data"`
}

// GetUser returns the authenticated user
func (c *Client) GetUser() (*User, error) {
	body, err := c.doRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}

	var resp UserResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.User, nil
}

// Budget types
type Budget struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	LastModifiedOn string          `json:"last_modified_on"`
	FirstMonth     string          `json:"first_month"`
	LastMonth      string          `json:"last_month"`
	DateFormat     *DateFormat     `json:"date_format"`
	CurrencyFormat *CurrencyFormat `json:"currency_format"`
}

type DateFormat struct {
	Format string `json:"format"`
}

type CurrencyFormat struct {
	ISOCode          string `json:"iso_code"`
	ExampleFormat    string `json:"example_format"`
	DecimalDigits    int    `json:"decimal_digits"`
	DecimalSeparator string `json:"decimal_separator"`
	SymbolFirst      bool   `json:"symbol_first"`
	GroupSeparator   string `json:"group_separator"`
	CurrencySymbol   string `json:"currency_symbol"`
	DisplaySymbol    bool   `json:"display_symbol"`
}

type BudgetSettings struct {
	DateFormat     DateFormat     `json:"date_format"`
	CurrencyFormat CurrencyFormat `json:"currency_format"`
}

type BudgetsResponse struct {
	Data struct {
		Budgets []Budget `json:"budgets"`
	} `json:"data"`
}

type BudgetDetailResponse struct {
	Data struct {
		Budget Budget `json:"budget"`
	} `json:"data"`
}

type BudgetSettingsResponse struct {
	Data struct {
		Settings BudgetSettings `json:"settings"`
	} `json:"data"`
}

// GetBudgets returns all budgets
func (c *Client) GetBudgets() ([]Budget, error) {
	body, err := c.doRequest("GET", "/budgets", nil)
	if err != nil {
		return nil, err
	}

	var resp BudgetsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Budgets, nil
}

// GetBudget returns a specific budget
func (c *Client) GetBudget(budgetID string) (*Budget, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s", budgetID), nil)
	if err != nil {
		return nil, err
	}

	var resp BudgetDetailResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Budget, nil
}

// GetBudgetSettings returns settings for a specific budget
func (c *Client) GetBudgetSettings(budgetID string) (*BudgetSettings, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/settings", budgetID), nil)
	if err != nil {
		return nil, err
	}

	var resp BudgetSettingsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Settings, nil
}

// Account types
type Account struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	OnBudget               bool   `json:"on_budget"`
	Closed                 bool   `json:"closed"`
	Note                   string `json:"note"`
	Balance                int64  `json:"balance"`
	ClearedBalance         int64  `json:"cleared_balance"`
	UnclearedBalance       int64  `json:"uncleared_balance"`
	TransferPayeeID        string `json:"transfer_payee_id"`
	DirectImportLinked     bool   `json:"direct_import_linked"`
	DirectImportInError    bool   `json:"direct_import_in_error"`
	LastReconciledAt       string `json:"last_reconciled_at"`
	DebtOriginalBalance    int64  `json:"debt_original_balance"`
	DebtInterestRates      map[string]int64 `json:"debt_interest_rates"`
	DebtMinimumPayments    map[string]int64 `json:"debt_minimum_payments"`
	DebtEscrowAmounts      map[string]int64 `json:"debt_escrow_amounts"`
	Deleted                bool   `json:"deleted"`
}

type AccountsResponse struct {
	Data struct {
		Accounts []Account `json:"accounts"`
	} `json:"data"`
}

type AccountResponse struct {
	Data struct {
		Account Account `json:"account"`
	} `json:"data"`
}

// GetAccounts returns all accounts for a budget
func (c *Client) GetAccounts(budgetID string) ([]Account, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/accounts", budgetID), nil)
	if err != nil {
		return nil, err
	}

	var resp AccountsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Accounts, nil
}

// GetAccount returns a specific account
func (c *Client) GetAccount(budgetID, accountID string) (*Account, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/accounts/%s", budgetID, accountID), nil)
	if err != nil {
		return nil, err
	}

	var resp AccountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Account, nil
}

// CreateAccountRequest represents the request to create an account
type CreateAccountRequest struct {
	Account struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Balance int64  `json:"balance"`
	} `json:"account"`
}

// CreateAccount creates a new account
func (c *Client) CreateAccount(budgetID, name, accountType string, balance int64) (*Account, error) {
	req := CreateAccountRequest{}
	req.Account.Name = name
	req.Account.Type = accountType
	req.Account.Balance = balance

	body, err := c.doRequest("POST", fmt.Sprintf("/budgets/%s/accounts", budgetID), req)
	if err != nil {
		return nil, err
	}

	var resp AccountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Account, nil
}

// Category types
type CategoryGroup struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Hidden     bool       `json:"hidden"`
	Deleted    bool       `json:"deleted"`
	Categories []Category `json:"categories"`
}

type Category struct {
	ID                      string `json:"id"`
	CategoryGroupID         string `json:"category_group_id"`
	CategoryGroupName       string `json:"category_group_name"`
	Name                    string `json:"name"`
	Hidden                  bool   `json:"hidden"`
	OriginalCategoryGroupID string `json:"original_category_group_id"`
	Note                    string `json:"note"`
	Budgeted                int64  `json:"budgeted"`
	Activity                int64  `json:"activity"`
	Balance                 int64  `json:"balance"`
	GoalType                string `json:"goal_type"`
	GoalDay                 int    `json:"goal_day"`
	GoalCadence             int    `json:"goal_cadence"`
	GoalCadenceFrequency    int    `json:"goal_cadence_frequency"`
	GoalCreationMonth       string `json:"goal_creation_month"`
	GoalTarget              int64  `json:"goal_target"`
	GoalTargetMonth         string `json:"goal_target_month"`
	GoalPercentageComplete  int    `json:"goal_percentage_complete"`
	GoalMonthsToBudget      int    `json:"goal_months_to_budget"`
	GoalUnderFunded         int64  `json:"goal_under_funded"`
	GoalOverallFunded       int64  `json:"goal_overall_funded"`
	GoalOverallLeft         int64  `json:"goal_overall_left"`
	Deleted                 bool   `json:"deleted"`
}

type CategoriesResponse struct {
	Data struct {
		CategoryGroups []CategoryGroup `json:"category_groups"`
	} `json:"data"`
}

type CategoryResponse struct {
	Data struct {
		Category Category `json:"category"`
	} `json:"data"`
}

// GetCategories returns all categories for a budget
func (c *Client) GetCategories(budgetID string) ([]CategoryGroup, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/categories", budgetID), nil)
	if err != nil {
		return nil, err
	}

	var resp CategoriesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.CategoryGroups, nil
}

// GetCategory returns a specific category
func (c *Client) GetCategory(budgetID, categoryID string) (*Category, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/categories/%s", budgetID, categoryID), nil)
	if err != nil {
		return nil, err
	}

	var resp CategoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Category, nil
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Category struct {
		Budgeted int64 `json:"budgeted"`
	} `json:"category"`
}

// UpdateCategory updates a category for a specific month
func (c *Client) UpdateCategory(budgetID, categoryID, month string, budgeted int64) (*Category, error) {
	req := UpdateCategoryRequest{}
	req.Category.Budgeted = budgeted

	body, err := c.doRequest("PATCH", fmt.Sprintf("/budgets/%s/months/%s/categories/%s", budgetID, month, categoryID), req)
	if err != nil {
		return nil, err
	}

	var resp CategoryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Category, nil
}

// Payee types
type Payee struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	TransferAccountID string `json:"transfer_account_id"`
	Deleted           bool   `json:"deleted"`
}

type PayeesResponse struct {
	Data struct {
		Payees []Payee `json:"payees"`
	} `json:"data"`
}

type PayeeResponse struct {
	Data struct {
		Payee Payee `json:"payee"`
	} `json:"data"`
}

// GetPayees returns all payees for a budget
func (c *Client) GetPayees(budgetID string) ([]Payee, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/payees", budgetID), nil)
	if err != nil {
		return nil, err
	}

	var resp PayeesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Payees, nil
}

// GetPayee returns a specific payee
func (c *Client) GetPayee(budgetID, payeeID string) (*Payee, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/payees/%s", budgetID, payeeID), nil)
	if err != nil {
		return nil, err
	}

	var resp PayeeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Payee, nil
}

// UpdatePayeeRequest represents the request to update a payee
type UpdatePayeeRequest struct {
	Payee struct {
		Name string `json:"name"`
	} `json:"payee"`
}

// UpdatePayee updates a payee
func (c *Client) UpdatePayee(budgetID, payeeID, name string) (*Payee, error) {
	req := UpdatePayeeRequest{}
	req.Payee.Name = name

	body, err := c.doRequest("PATCH", fmt.Sprintf("/budgets/%s/payees/%s", budgetID, payeeID), req)
	if err != nil {
		return nil, err
	}

	var resp PayeeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Payee, nil
}

// Transaction types
type Transaction struct {
	ID                    string        `json:"id"`
	Date                  string        `json:"date"`
	Amount                int64         `json:"amount"`
	Memo                  string        `json:"memo"`
	Cleared               string        `json:"cleared"`
	Approved              bool          `json:"approved"`
	FlagColor             string        `json:"flag_color"`
	FlagName              string        `json:"flag_name"`
	AccountID             string        `json:"account_id"`
	AccountName           string        `json:"account_name"`
	PayeeID               string        `json:"payee_id"`
	PayeeName             string        `json:"payee_name"`
	CategoryID            string        `json:"category_id"`
	CategoryName          string        `json:"category_name"`
	TransferAccountID     string        `json:"transfer_account_id"`
	TransferTransactionID string        `json:"transfer_transaction_id"`
	MatchedTransactionID  string        `json:"matched_transaction_id"`
	ImportID              string        `json:"import_id"`
	ImportPayeeName       string        `json:"import_payee_name"`
	ImportPayeeNameOriginal string      `json:"import_payee_name_original"`
	DebtTransactionType   string        `json:"debt_transaction_type"`
	Deleted               bool          `json:"deleted"`
	Subtransactions       []Subtransaction `json:"subtransactions"`
}

type Subtransaction struct {
	ID                    string `json:"id"`
	TransactionID         string `json:"transaction_id"`
	Amount                int64  `json:"amount"`
	Memo                  string `json:"memo"`
	PayeeID               string `json:"payee_id"`
	PayeeName             string `json:"payee_name"`
	CategoryID            string `json:"category_id"`
	CategoryName          string `json:"category_name"`
	TransferAccountID     string `json:"transfer_account_id"`
	TransferTransactionID string `json:"transfer_transaction_id"`
	Deleted               bool   `json:"deleted"`
}

type TransactionsResponse struct {
	Data struct {
		Transactions []Transaction `json:"transactions"`
	} `json:"data"`
}

type TransactionResponse struct {
	Data struct {
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}

// TransactionFilter contains filters for listing transactions
type TransactionFilter struct {
	SinceDate   string
	Type        string
	AccountID   string
	CategoryID  string
	PayeeID     string
}

// GetTransactions returns transactions for a budget with optional filters
func (c *Client) GetTransactions(budgetID string, filter *TransactionFilter) ([]Transaction, error) {
	path := fmt.Sprintf("/budgets/%s/transactions", budgetID)

	if filter != nil {
		params := url.Values{}
		if filter.SinceDate != "" {
			params.Set("since_date", filter.SinceDate)
		}
		if filter.Type != "" {
			params.Set("type", filter.Type)
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}
	}

	body, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp TransactionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Transactions, nil
}

// GetTransactionsByAccount returns transactions for a specific account
func (c *Client) GetTransactionsByAccount(budgetID, accountID string, sinceDate string) ([]Transaction, error) {
	path := fmt.Sprintf("/budgets/%s/accounts/%s/transactions", budgetID, accountID)
	if sinceDate != "" {
		path += "?since_date=" + sinceDate
	}

	body, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp TransactionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Transactions, nil
}

// GetTransactionsByCategory returns transactions for a specific category
func (c *Client) GetTransactionsByCategory(budgetID, categoryID string, sinceDate string) ([]Transaction, error) {
	path := fmt.Sprintf("/budgets/%s/categories/%s/transactions", budgetID, categoryID)
	if sinceDate != "" {
		path += "?since_date=" + sinceDate
	}

	body, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp TransactionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Transactions, nil
}

// GetTransactionsByPayee returns transactions for a specific payee
func (c *Client) GetTransactionsByPayee(budgetID, payeeID string, sinceDate string) ([]Transaction, error) {
	path := fmt.Sprintf("/budgets/%s/payees/%s/transactions", budgetID, payeeID)
	if sinceDate != "" {
		path += "?since_date=" + sinceDate
	}

	body, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp TransactionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Transactions, nil
}

// GetTransaction returns a specific transaction
func (c *Client) GetTransaction(budgetID, transactionID string) (*Transaction, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/transactions/%s", budgetID, transactionID), nil)
	if err != nil {
		return nil, err
	}

	var resp TransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Transaction, nil
}

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	Transaction SaveTransaction `json:"transaction"`
}

type SaveTransaction struct {
	AccountID  string `json:"account_id"`
	Date       string `json:"date"`
	Amount     int64  `json:"amount"`
	PayeeID    string `json:"payee_id,omitempty"`
	PayeeName  string `json:"payee_name,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	Memo       string `json:"memo,omitempty"`
	Cleared    string `json:"cleared,omitempty"`
	Approved   bool   `json:"approved,omitempty"`
	FlagColor  string `json:"flag_color,omitempty"`
	ImportID   string `json:"import_id,omitempty"`
}

// CreateTransaction creates a new transaction
func (c *Client) CreateTransaction(budgetID string, txn SaveTransaction) (*Transaction, error) {
	req := CreateTransactionRequest{Transaction: txn}

	body, err := c.doRequest("POST", fmt.Sprintf("/budgets/%s/transactions", budgetID), req)
	if err != nil {
		return nil, err
	}

	var resp TransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Transaction, nil
}

// UpdateTransactionRequest represents the request to update a transaction
type UpdateTransactionRequest struct {
	Transaction SaveTransaction `json:"transaction"`
}

// UpdateTransaction updates an existing transaction
func (c *Client) UpdateTransaction(budgetID, transactionID string, txn SaveTransaction) (*Transaction, error) {
	req := UpdateTransactionRequest{Transaction: txn}

	body, err := c.doRequest("PUT", fmt.Sprintf("/budgets/%s/transactions/%s", budgetID, transactionID), req)
	if err != nil {
		return nil, err
	}

	var resp TransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Transaction, nil
}

// DeleteTransaction deletes a transaction
func (c *Client) DeleteTransaction(budgetID, transactionID string) (*Transaction, error) {
	body, err := c.doRequest("DELETE", fmt.Sprintf("/budgets/%s/transactions/%s", budgetID, transactionID), nil)
	if err != nil {
		return nil, err
	}

	var resp TransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Transaction, nil
}

// ScheduledTransaction types
type ScheduledTransaction struct {
	ID                    string                    `json:"id"`
	DateFirst             string                    `json:"date_first"`
	DateNext              string                    `json:"date_next"`
	Frequency             string                    `json:"frequency"`
	Amount                int64                     `json:"amount"`
	Memo                  string                    `json:"memo"`
	FlagColor             string                    `json:"flag_color"`
	FlagName              string                    `json:"flag_name"`
	AccountID             string                    `json:"account_id"`
	AccountName           string                    `json:"account_name"`
	PayeeID               string                    `json:"payee_id"`
	PayeeName             string                    `json:"payee_name"`
	CategoryID            string                    `json:"category_id"`
	CategoryName          string                    `json:"category_name"`
	TransferAccountID     string                    `json:"transfer_account_id"`
	Deleted               bool                      `json:"deleted"`
	Subtransactions       []ScheduledSubtransaction `json:"subtransactions"`
}

type ScheduledSubtransaction struct {
	ID                string `json:"id"`
	ScheduledTransactionID string `json:"scheduled_transaction_id"`
	Amount            int64  `json:"amount"`
	Memo              string `json:"memo"`
	PayeeID           string `json:"payee_id"`
	CategoryID        string `json:"category_id"`
	TransferAccountID string `json:"transfer_account_id"`
	Deleted           bool   `json:"deleted"`
}

type ScheduledTransactionsResponse struct {
	Data struct {
		ScheduledTransactions []ScheduledTransaction `json:"scheduled_transactions"`
	} `json:"data"`
}

type ScheduledTransactionResponse struct {
	Data struct {
		ScheduledTransaction ScheduledTransaction `json:"scheduled_transaction"`
	} `json:"data"`
}

// GetScheduledTransactions returns all scheduled transactions for a budget
func (c *Client) GetScheduledTransactions(budgetID string) ([]ScheduledTransaction, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/scheduled_transactions", budgetID), nil)
	if err != nil {
		return nil, err
	}

	var resp ScheduledTransactionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.ScheduledTransactions, nil
}

// GetScheduledTransaction returns a specific scheduled transaction
func (c *Client) GetScheduledTransaction(budgetID, scheduledTransactionID string) (*ScheduledTransaction, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/scheduled_transactions/%s", budgetID, scheduledTransactionID), nil)
	if err != nil {
		return nil, err
	}

	var resp ScheduledTransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.ScheduledTransaction, nil
}

// SaveScheduledTransaction represents a scheduled transaction to create or update
type SaveScheduledTransaction struct {
	AccountID  string `json:"account_id"`
	Date       string `json:"date"`
	Frequency  string `json:"frequency"`
	Amount     int64  `json:"amount"`
	PayeeID    string `json:"payee_id,omitempty"`
	PayeeName  string `json:"payee_name,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	Memo       string `json:"memo,omitempty"`
	FlagColor  string `json:"flag_color,omitempty"`
}

// CreateScheduledTransactionRequest represents the request to create a scheduled transaction
type CreateScheduledTransactionRequest struct {
	ScheduledTransaction SaveScheduledTransaction `json:"scheduled_transaction"`
}

// CreateScheduledTransaction creates a new scheduled transaction
func (c *Client) CreateScheduledTransaction(budgetID string, st SaveScheduledTransaction) (*ScheduledTransaction, error) {
	req := CreateScheduledTransactionRequest{ScheduledTransaction: st}

	body, err := c.doRequest("POST", fmt.Sprintf("/budgets/%s/scheduled_transactions", budgetID), req)
	if err != nil {
		return nil, err
	}

	var resp ScheduledTransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.ScheduledTransaction, nil
}

// UpdateScheduledTransaction updates an existing scheduled transaction
func (c *Client) UpdateScheduledTransaction(budgetID, scheduledTransactionID string, st SaveScheduledTransaction) (*ScheduledTransaction, error) {
	req := CreateScheduledTransactionRequest{ScheduledTransaction: st}

	body, err := c.doRequest("PUT", fmt.Sprintf("/budgets/%s/scheduled_transactions/%s", budgetID, scheduledTransactionID), req)
	if err != nil {
		return nil, err
	}

	var resp ScheduledTransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.ScheduledTransaction, nil
}

// DeleteScheduledTransaction deletes a scheduled transaction
func (c *Client) DeleteScheduledTransaction(budgetID, scheduledTransactionID string) (*ScheduledTransaction, error) {
	body, err := c.doRequest("DELETE", fmt.Sprintf("/budgets/%s/scheduled_transactions/%s", budgetID, scheduledTransactionID), nil)
	if err != nil {
		return nil, err
	}

	var resp ScheduledTransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.ScheduledTransaction, nil
}

// Month types
type Month struct {
	Month      string     `json:"month"`
	Note       string     `json:"note"`
	Income     int64      `json:"income"`
	Budgeted   int64      `json:"budgeted"`
	Activity   int64      `json:"activity"`
	ToBeBudgeted int64    `json:"to_be_budgeted"`
	AgeOfMoney int        `json:"age_of_money"`
	Deleted    bool       `json:"deleted"`
	Categories []Category `json:"categories"`
}

type MonthsResponse struct {
	Data struct {
		Months []Month `json:"months"`
	} `json:"data"`
}

type MonthResponse struct {
	Data struct {
		Month Month `json:"month"`
	} `json:"data"`
}

// GetMonths returns all budget months
func (c *Client) GetMonths(budgetID string) ([]Month, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/months", budgetID), nil)
	if err != nil {
		return nil, err
	}

	var resp MonthsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Data.Months, nil
}

// GetMonth returns a specific budget month
func (c *Client) GetMonth(budgetID, month string) (*Month, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/budgets/%s/months/%s", budgetID, month), nil)
	if err != nil {
		return nil, err
	}

	var resp MonthResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Data.Month, nil
}

// Currency conversion helpers

// MilliunitsToAmount converts YNAB milliunits to a float amount
func MilliunitsToAmount(milliunits int64) float64 {
	return float64(milliunits) / 1000.0
}

// AmountToMilliunits converts a float amount to YNAB milliunits
func AmountToMilliunits(amount float64) int64 {
	return int64(amount * 1000)
}
