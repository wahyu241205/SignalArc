package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/database"
	"github.com/wahyu241205/SignalArc/backend/internal/repository"
)

const (
	defaultListLimit    = 50
	defaultMarketsLimit = defaultListLimit
)

func NewRouter(db *database.DB) http.Handler {
	router := chi.NewRouter()
	router.Use(requestIDMiddleware, localCORSMiddleware, requestLoggingMiddleware, recovererMiddleware)

	marketsRepository := repository.NewMarketsRepository(db)
	positionsRepository := repository.NewPositionsRepository(db)
	resolutionsRepository := repository.NewResolutionsRepository(db)
	settlementsRepository := repository.NewSettlementsRepository(db)
	tradesRepository := repository.NewTradesRepository(db)

	registerStatusRoutes(router, db)
	registerArcRoutes(router)
	registerMarketRoutes(router, marketsRepository)
	registerTradeRoutes(router, tradesRepository, marketsRepository)
	registerPositionRoutes(router, positionsRepository)
	registerResolutionRoutes(router, resolutionsRepository)
	registerSettlementRoutes(router, settlementsRepository)

	return router
}
