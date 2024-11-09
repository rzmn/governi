package spendings

import (
	"verni/internal/common"
	spendingsRepository "verni/internal/repositories/spendings"
	"verni/internal/services/logging"
)

type CounterpartyId spendingsRepository.CounterpartyId
type ExpenseId spendingsRepository.ExpenseId
type Expense spendingsRepository.Expense
type IdentifiableExpense spendingsRepository.IdentifiableExpense
type Balance spendingsRepository.Balance
type Repository spendingsRepository.Repository

type Controller interface {
	AddExpense(expense Expense, actor CounterpartyId) (IdentifiableExpense, *common.CodeBasedError[AddExpenseErrorCode])
	RemoveExpense(expenseId ExpenseId, actor CounterpartyId) (IdentifiableExpense, *common.CodeBasedError[RemoveExpenseErrorCode])
	GetExpense(expenseId ExpenseId, actor CounterpartyId) (IdentifiableExpense, *common.CodeBasedError[GetExpenseErrorCode])
	GetExpensesWith(counterparty CounterpartyId, actor CounterpartyId) ([]IdentifiableExpense, *common.CodeBasedError[GetExpensesErrorCode])
	GetBalance(actor CounterpartyId) ([]Balance, *common.CodeBasedError[GetBalanceErrorCode])
}

func DefaultController(repository Repository, logger logging.Service) Controller {
	return &defaultController{
		repository: repository,
		logger:     logger,
	}
}
