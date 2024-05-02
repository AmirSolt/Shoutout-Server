package pages

import (
	"basedpocket/base"
	"basedpocket/cmodels"
	"fmt"
	"net/http"
	"sort"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type LandingPageResponse struct {
	TotalRaised         float64           `json:"total_raised"`
	HighestTransactions *[]TransactionAPI `json:"highest_transactions"`
	RecentTransactions  *[]TransactionAPI `json:"recent_transactions"`
}

type TransactionAPI struct {
	Created  types.DateTime `db:"created" json:"created"`
	Amount   float64        `db:"amount" json:"amount"`
	UserName string         `db:"user_name" json:"user_name"`
}

const ResponseTransLimit = 5

func handleLandingPageData(app core.App, ctx echo.Context, env *base.Env) error {

	characterName := ctx.QueryParam("character_name")

	if characterName == "" {
		return fmt.Errorf("characterName query param is empty")
	}

	var model *cmodels.Transaction
	var transactions *[]cmodels.Transaction
	if err := app.Dao().ModelQuery(model).
		AndWhere(dbx.HashExp{"character_name": characterName}).
		All(&transactions); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	transactionsAPI := convertTransactionsForAPI(*transactions)

	landingPageResponse := LandingPageResponse{
		TotalRaised:         SumTransactionAmounts(transactionsAPI),
		HighestTransactions: GetMostRecentTransactions(transactionsAPI, ResponseTransLimit),
		RecentTransactions:  GetHighestTransactions(transactionsAPI, ResponseTransLimit),
	}

	return ctx.JSON(http.StatusOK, landingPageResponse)
}

func SumTransactionAmounts(transactions []TransactionAPI) float64 {
	total := 0.0
	for _, transaction := range transactions {
		total += transaction.Amount
	}
	return total
}

func GetMostRecentTransactions(transactions []TransactionAPI, limit int) *[]TransactionAPI {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Created.Time().After(transactions[j].Created.Time())
	})
	slicedTransactions := transactions
	if len(transactions) > limit {
		slicedTransactions = transactions[:limit]
	}
	return &slicedTransactions
}

func GetHighestTransactions(transactions []TransactionAPI, limit int) *[]TransactionAPI {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Amount > transactions[j].Amount
	})
	slicedTransactions := transactions
	if len(transactions) > limit {
		slicedTransactions = transactions[:limit]
	}
	return &slicedTransactions
}

func convertTransactionsForAPI(transactions []cmodels.Transaction) []TransactionAPI {
	transAPI := []TransactionAPI{}

	for _, trans := range transactions {
		transAPI = append(transAPI, TransactionAPI{
			Created:  trans.Created,
			Amount:   trans.Amount,
			UserName: trans.UserName,
		})
	}

	return transAPI
}
