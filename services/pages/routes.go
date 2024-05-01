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
)

type LandingPageResponse struct {
	TotalRaised         float64                `json:"total_raised"`
	HighestTransactions *[]cmodels.Transaction `json:"top_donors"`
	RecentTransactions  *[]cmodels.Transaction `json:"recent_transactions"`
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

	landingPageResponse := LandingPageResponse{
		TotalRaised:         SumTransactionAmounts(*transactions),
		HighestTransactions: GetMostRecentTransactions(*transactions, ResponseTransLimit),
		RecentTransactions:  GetHighestTransactions(*transactions, ResponseTransLimit),
	}

	return ctx.JSON(http.StatusOK, landingPageResponse)
}

func SumTransactionAmounts(transactions []cmodels.Transaction) float64 {
	total := 0.0
	for _, transaction := range transactions {
		total += transaction.Amount
	}
	return total
}

func GetMostRecentTransactions(transactions []cmodels.Transaction, limit int) *[]cmodels.Transaction {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Created.Time().After(transactions[j].Created.Time())
	})
	slicedTransactions := transactions
	if len(transactions) > limit {
		slicedTransactions = transactions[:limit]
	}
	return &slicedTransactions
}

func GetHighestTransactions(transactions []cmodels.Transaction, limit int) *[]cmodels.Transaction {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Amount > transactions[j].Amount
	})
	slicedTransactions := transactions
	if len(transactions) > limit {
		slicedTransactions = transactions[:limit]
	}
	return &slicedTransactions
}
