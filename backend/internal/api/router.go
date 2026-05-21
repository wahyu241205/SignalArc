package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/agent"
	"github.com/wahyu241205/SignalArc/backend/internal/config"
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
	agentWalletsRepository := repository.NewAgentWalletsRepository(db)
	agentIntentStore := agent.NewStore()
	cfg := config.Load()
	circleExecutor := agent.NewCircleCLIExecutor(agent.CircleCLIExecutorConfig{
		Enabled:      cfg.CircleAgentWalletExecutionEnabled,
		CLIPath:      cfg.CircleCLIPath,
		Chain:        cfg.CircleAgentWalletChain,
		Timeout:      time.Duration(cfg.CircleAgentWalletTimeoutSeconds) * time.Second,
		AgentFactory: agent.AgentFactoryAddress,
	})

	registerStatusRoutes(router, db)
	registerArcRoutes(router)
	registerMarketRoutes(router, marketsRepository)
	registerTradeRoutes(router, tradesRepository, marketsRepository)
	registerPositionRoutes(router, positionsRepository)
	registerResolutionRoutes(router, resolutionsRepository)
	registerSettlementRoutes(router, settlementsRepository)
	registerAgentIntentRoutes(router, agentIntentStore, agentWalletsRepository, circleExecutor)

	return router
}
