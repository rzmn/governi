package spendings

import (
	"net/http"
	"verni/internal/common"
	spendingsController "verni/internal/controllers/spendings"
	spendingsRepository "verni/internal/repositories/spendings"
	"verni/internal/schema"
	"verni/internal/services/logging"
	"verni/internal/services/longpoll"
	"verni/internal/services/pushNotifications"
)

type defaultRequestsHandler struct {
	controller     spendingsController.Controller
	pushService    pushNotifications.Service
	pollingService longpoll.Service
	logger         logging.Service
}

func (c *defaultRequestsHandler) AddExpense(
	subject schema.UserId,
	request AddExpenseRequest,
	success func(schema.StatusCode, schema.Response[schema.IdentifiableExpense]),
	failure func(schema.StatusCode, schema.Response[schema.Error]),
) {
	expense, err := c.controller.AddExpense(mapHttpServerExpense(request.Expense), spendingsController.CounterpartyId(subject))
	if err != nil {
		switch err.Code {
		case spendingsController.AddExpenseErrorNoSuchUser:
			failure(
				http.StatusConflict,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeNoSuchUser,
						err.Error(),
					),
				),
			)
		case spendingsController.AddExpenseErrorNotYourExpense:
			failure(
				http.StatusConflict,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeIsNotYourExpense,
						err.Error(),
					),
				),
			)
		default:
			c.logger.LogError("addExpense request %v failed with unknown err: %v", request, err)
			failure(
				http.StatusInternalServerError,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeInternal,
						err.Error(),
					),
				),
			)
		}
		return
	}
	for _, share := range expense.Shares {
		if share.Counterparty == spendingsRepository.CounterpartyId(subject) {
			continue
		}
		c.pushService.NewExpenseReceived(
			pushNotifications.UserId(share.Counterparty),
			pushNotifications.Expense(mapIdentifiableExpense(expense)),
			pushNotifications.UserId(subject),
		)
		c.pollingService.ExpensesUpdated(longpoll.UserId(share.Counterparty), longpoll.UserId(subject))
		c.pollingService.CounterpartiesUpdated(longpoll.UserId(share.Counterparty))
	}
	success(http.StatusOK, schema.Success(mapIdentifiableExpense(expense)))
}

func (c *defaultRequestsHandler) RemoveExpense(
	subject schema.UserId,
	request RemoveExpenseRequest,
	success func(schema.StatusCode, schema.Response[schema.IdentifiableExpense]),
	failure func(schema.StatusCode, schema.Response[schema.Error]),
) {
	expense, err := c.controller.RemoveExpense(spendingsController.ExpenseId(request.ExpenseId), spendingsController.CounterpartyId(subject))
	if err != nil {
		switch err.Code {
		case spendingsController.RemoveExpenseErrorExpenseNotFound:
			failure(
				http.StatusConflict,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeExpenseNotFound,
						err.Error(),
					),
				),
			)
		case spendingsController.RemoveExpenseErrorNotAFriend:
			failure(
				http.StatusConflict,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeNotAFriend,
						err.Error(),
					),
				),
			)
		case spendingsController.RemoveExpenseErrorNotYourExpense:
			failure(
				http.StatusConflict,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeIsNotYourExpense,
						err.Error(),
					),
				),
			)
		default:
			c.logger.LogError("removeExpense request %v failed with unknown err: %v", request, err)
			failure(
				http.StatusInternalServerError,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeInternal,
						err.Error(),
					),
				),
			)
		}
		return
	}
	for _, share := range expense.Shares {
		if share.Counterparty == spendingsRepository.CounterpartyId(subject) {
			continue
		}
		c.pushService.NewExpenseReceived(
			pushNotifications.UserId(share.Counterparty),
			pushNotifications.Expense(mapIdentifiableExpense(expense)),
			pushNotifications.UserId(subject),
		)
		c.pollingService.ExpensesUpdated(longpoll.UserId(share.Counterparty), longpoll.UserId(subject))
		c.pollingService.CounterpartiesUpdated(longpoll.UserId(share.Counterparty))
	}
	success(http.StatusOK, schema.Success(mapIdentifiableExpense(expense)))
}

func (c *defaultRequestsHandler) GetBalance(
	subject schema.UserId,
	success func(schema.StatusCode, schema.Response[[]schema.Balance]),
	failure func(schema.StatusCode, schema.Response[schema.Error]),
) {
	balance, err := c.controller.GetBalance(spendingsController.CounterpartyId(subject))
	if err != nil {
		switch err.Code {
		default:
			c.logger.LogError("getBalance request failed with unknown err: %v", err)
			failure(
				http.StatusInternalServerError,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeInternal,
						err.Error(),
					),
				),
			)
		}
		return
	}
	success(http.StatusOK, schema.Success(common.Map(balance, mapBalance)))
}

func (c *defaultRequestsHandler) GetExpenses(
	subject schema.UserId,
	request GetExpensesRequest,
	success func(schema.StatusCode, schema.Response[[]schema.IdentifiableExpense]),
	failure func(schema.StatusCode, schema.Response[schema.Error]),
) {
	expenses, err := c.controller.GetExpensesWith(spendingsController.CounterpartyId(request.Counterparty), spendingsController.CounterpartyId(subject))
	if err != nil {
		switch err.Code {
		default:
			c.logger.LogError("getExpenses request %v failed with unknown err: %v", request, err)
			failure(
				http.StatusInternalServerError,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeInternal,
						err.Error(),
					),
				),
			)
		}
		return
	}
	success(http.StatusOK, schema.Success(common.Map(expenses, mapIdentifiableExpense)))
}

func (c *defaultRequestsHandler) GetExpense(
	subject schema.UserId,
	request GetExpenseRequest,
	success func(schema.StatusCode, schema.Response[schema.IdentifiableExpense]),
	failure func(schema.StatusCode, schema.Response[schema.Error]),
) {
	expense, err := c.controller.GetExpense(spendingsController.ExpenseId(request.Id), spendingsController.CounterpartyId(subject))
	if err != nil {
		switch err.Code {
		case spendingsController.GetExpenseErrorExpenseNotFound:
			failure(
				http.StatusConflict,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeExpenseNotFound,
						err.Error(),
					),
				),
			)
		case spendingsController.GetExpenseErrorNotYourExpense:
			failure(
				http.StatusConflict,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeIsNotYourExpense,
						err.Error(),
					),
				),
			)
		default:
			c.logger.LogError("getExpense request %v failed with unknown err: %v", request, err)
			failure(
				http.StatusInternalServerError,
				schema.Failure(
					common.NewErrorWithDescriptionValue(
						schema.CodeInternal,
						err.Error(),
					),
				),
			)
		}
		return
	}
	success(http.StatusOK, schema.Success(mapIdentifiableExpense(expense)))
}

func mapHttpServerExpense(expense schema.Expense) spendingsController.Expense {
	return spendingsController.Expense{
		Timestamp: expense.Timestamp,
		Details:   expense.Details,
		Total:     spendingsRepository.Cost(expense.Total),
		Currency:  spendingsRepository.Currency(expense.Currency),
		Shares: common.Map(expense.Shares, func(share schema.ShareOfExpense) spendingsRepository.ShareOfExpense {
			return spendingsRepository.ShareOfExpense{
				Counterparty: spendingsRepository.CounterpartyId(share.UserId),
				Cost:         spendingsRepository.Cost(share.Cost),
			}
		}),
	}
}

func mapIdentifiableExpense(expense spendingsController.IdentifiableExpense) schema.IdentifiableExpense {
	return schema.IdentifiableExpense{
		Id:      schema.ExpenseId(expense.Id),
		Expense: mapExpense(spendingsController.Expense(expense.Expense)),
	}
}

func mapExpense(expense spendingsController.Expense) schema.Expense {
	return schema.Expense{
		Timestamp:   expense.Timestamp,
		Details:     expense.Details,
		Total:       schema.Cost(expense.Total),
		Attachments: []schema.ExpenseAttachment{},
		Currency:    schema.Currency(expense.Currency),
		Shares:      common.Map(expense.Shares, mapShareOfExpense),
	}
}

func mapShareOfExpense(share spendingsRepository.ShareOfExpense) schema.ShareOfExpense {
	return schema.ShareOfExpense{
		UserId: schema.UserId(share.Counterparty),
		Cost:   schema.Cost(share.Cost),
	}
}

func mapBalance(balance spendingsController.Balance) schema.Balance {
	currencies := map[schema.Currency]schema.Cost{}
	for currency, cost := range balance.Currencies {
		currencies[schema.Currency(currency)] = schema.Cost(cost)
	}
	return schema.Balance{
		Counterparty: string(balance.Counterparty),
		Currencies:   currencies,
	}
}