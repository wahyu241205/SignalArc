package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/wahyu241205/SignalArc/backend/internal/database"
)

type Trade struct {
	ID               string         `json:"id"`
	UserID           string         `json:"user_id"`
	MarketID         string         `json:"market_id"`
	Outcome          string         `json:"outcome"`
	Side             string         `json:"side"`
	Quantity         string         `json:"quantity"`
	Price            string         `json:"price"`
	CollateralAmount string         `json:"collateral_amount"`
	FeeAmount        string         `json:"fee_amount"`
	Status           string         `json:"status"`
	TxHash           sql.NullString `json:"tx_hash"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type TradesRepository struct {
	db *database.DB
}

func NewTradesRepository(db *database.DB) *TradesRepository {
	return &TradesRepository{db: db}
}

func (r *TradesRepository) GetTradeByID(ctx context.Context, id string) (Trade, error) {
	var trade Trade
	err := r.db.QueryRow(ctx, tradeSelectSQL+`
		WHERE id = $1
	`, id).Scan(
		&trade.ID,
		&trade.UserID,
		&trade.MarketID,
		&trade.Outcome,
		&trade.Side,
		&trade.Quantity,
		&trade.Price,
		&trade.CollateralAmount,
		&trade.FeeAmount,
		&trade.Status,
		&trade.TxHash,
		&trade.CreatedAt,
		&trade.UpdatedAt,
	)

	return trade, err
}

func (r *TradesRepository) ListTradesByUserID(ctx context.Context, userID string, limit int) ([]Trade, error) {
	rows, err := r.db.Query(ctx, tradeSelectSQL+`
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTrades(rows)
}

func (r *TradesRepository) ListTradesByMarketID(ctx context.Context, marketID string, limit int) ([]Trade, error) {
	rows, err := r.db.Query(ctx, tradeSelectSQL+`
		WHERE market_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, marketID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTrades(rows)
}

func scanTrades(rows pgx.Rows) ([]Trade, error) {
	trades := []Trade{}
	for rows.Next() {
		var trade Trade
		if err := rows.Scan(
			&trade.ID,
			&trade.UserID,
			&trade.MarketID,
			&trade.Outcome,
			&trade.Side,
			&trade.Quantity,
			&trade.Price,
			&trade.CollateralAmount,
			&trade.FeeAmount,
			&trade.Status,
			&trade.TxHash,
			&trade.CreatedAt,
			&trade.UpdatedAt,
		); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}

const tradeSelectSQL = `
	SELECT
		id::text,
		user_id::text,
		market_id::text,
		outcome,
		side,
		quantity::text,
		price::text,
		collateral_amount::text,
		fee_amount::text,
		status,
		tx_hash,
		created_at,
		updated_at
	FROM trades
`
