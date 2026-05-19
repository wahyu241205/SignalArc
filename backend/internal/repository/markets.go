package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/wahyu241205/SignalArc/backend/internal/database"
)

type Market struct {
	ID               string         `json:"id"`
	CreatorUserID    string         `json:"creator_user_id"`
	Title            string         `json:"title"`
	Description      sql.NullString `json:"description"`
	Category         sql.NullString `json:"category"`
	Status           string         `json:"status"`
	OutcomeYesLabel  string         `json:"outcome_yes_label"`
	OutcomeNoLabel   string         `json:"outcome_no_label"`
	CollateralAsset  string         `json:"collateral_asset"`
	Chain            string         `json:"chain"`
	ResolutionSource sql.NullString `json:"resolution_source"`
	OpensAt          sql.NullTime   `json:"opens_at"`
	ClosesAt         time.Time      `json:"closes_at"`
	ResolvedAt       sql.NullTime   `json:"resolved_at"`
	SettledAt        sql.NullTime   `json:"settled_at"`
	WinningOutcome   sql.NullString `json:"winning_outcome"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type MarketsRepository struct {
	db *database.DB
}

type CreateMarketInput struct {
	CreatorUserID    string
	Title            string
	Description      sql.NullString
	Category         sql.NullString
	OutcomeYesLabel  string
	OutcomeNoLabel   string
	CollateralAsset  string
	Chain            string
	ResolutionSource sql.NullString
	OpensAt          sql.NullTime
	ClosesAt         time.Time
}

func NewMarketsRepository(db *database.DB) *MarketsRepository {
	return &MarketsRepository{db: db}
}

func (r *MarketsRepository) GetMarketByID(ctx context.Context, id string) (Market, error) {
	var market Market
	err := r.db.QueryRow(ctx, marketSelectSQL+`
		WHERE id = $1
	`, id).Scan(
		&market.ID,
		&market.CreatorUserID,
		&market.Title,
		&market.Description,
		&market.Category,
		&market.Status,
		&market.OutcomeYesLabel,
		&market.OutcomeNoLabel,
		&market.CollateralAsset,
		&market.Chain,
		&market.ResolutionSource,
		&market.OpensAt,
		&market.ClosesAt,
		&market.ResolvedAt,
		&market.SettledAt,
		&market.WinningOutcome,
		&market.CreatedAt,
		&market.UpdatedAt,
	)

	return market, err
}

func (r *MarketsRepository) CreateMarket(ctx context.Context, input CreateMarketInput) (Market, error) {
	var market Market
	err := r.db.QueryRow(ctx, `
		INSERT INTO markets (
			creator_user_id,
			title,
			description,
			category,
			outcome_yes_label,
			outcome_no_label,
			collateral_asset,
			chain,
			resolution_source,
			opens_at,
			closes_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING
			id::text,
			creator_user_id::text,
			title,
			description,
			category,
			status,
			outcome_yes_label,
			outcome_no_label,
			collateral_asset,
			chain,
			resolution_source,
			opens_at,
			closes_at,
			resolved_at,
			settled_at,
			winning_outcome,
			created_at,
			updated_at
	`,
		input.CreatorUserID,
		input.Title,
		input.Description,
		input.Category,
		input.OutcomeYesLabel,
		input.OutcomeNoLabel,
		input.CollateralAsset,
		input.Chain,
		input.ResolutionSource,
		input.OpensAt,
		input.ClosesAt,
	).Scan(
		&market.ID,
		&market.CreatorUserID,
		&market.Title,
		&market.Description,
		&market.Category,
		&market.Status,
		&market.OutcomeYesLabel,
		&market.OutcomeNoLabel,
		&market.CollateralAsset,
		&market.Chain,
		&market.ResolutionSource,
		&market.OpensAt,
		&market.ClosesAt,
		&market.ResolvedAt,
		&market.SettledAt,
		&market.WinningOutcome,
		&market.CreatedAt,
		&market.UpdatedAt,
	)

	return market, err
}

func (r *MarketsRepository) ListMarkets(ctx context.Context, limit int) ([]Market, error) {
	rows, err := r.db.Query(ctx, marketSelectSQL+`
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	markets := []Market{}
	for rows.Next() {
		var market Market
		if err := rows.Scan(
			&market.ID,
			&market.CreatorUserID,
			&market.Title,
			&market.Description,
			&market.Category,
			&market.Status,
			&market.OutcomeYesLabel,
			&market.OutcomeNoLabel,
			&market.CollateralAsset,
			&market.Chain,
			&market.ResolutionSource,
			&market.OpensAt,
			&market.ClosesAt,
			&market.ResolvedAt,
			&market.SettledAt,
			&market.WinningOutcome,
			&market.CreatedAt,
			&market.UpdatedAt,
		); err != nil {
			return nil, err
		}
		markets = append(markets, market)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return markets, nil
}

const marketSelectSQL = `
	SELECT
		id::text,
		creator_user_id::text,
		title,
		description,
		category,
		status,
		outcome_yes_label,
		outcome_no_label,
		collateral_asset,
		chain,
		resolution_source,
		opens_at,
		closes_at,
		resolved_at,
		settled_at,
		winning_outcome,
		created_at,
		updated_at
	FROM markets
`
