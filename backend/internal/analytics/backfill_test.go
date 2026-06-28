package analytics

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/wahyu241205/SignalArc/backend/internal/repository"
)

type fakeLogsClient struct {
	pages []LogsPage
	calls int
}

func (client *fakeLogsClient) FetchAddressLogs(ctx context.Context, address string, pageParams map[string]string) (LogsPage, error) {
	if client.calls >= len(client.pages) {
		return LogsPage{}, nil
	}
	page := client.pages[client.calls]
	client.calls++
	return page, nil
}

type fakeAnalyticsStore struct {
	markets map[string]repository.UpsertAnalyticsMarketInput
	events  map[string]repository.InsertAnalyticsEventInput
	summary repository.AnalyticsSummary
	state   repository.UpdateAnalyticsIndexerStateInput
}

func newFakeAnalyticsStore() *fakeAnalyticsStore {
	return &fakeAnalyticsStore{
		markets: map[string]repository.UpsertAnalyticsMarketInput{},
		events:  map[string]repository.InsertAnalyticsEventInput{},
	}
}

func (store *fakeAnalyticsStore) UpsertAnalyticsMarket(ctx context.Context, input repository.UpsertAnalyticsMarketInput) error {
	store.markets[input.MarketAddress] = input
	return nil
}

func (store *fakeAnalyticsStore) InsertAnalyticsEvent(ctx context.Context, input repository.InsertAnalyticsEventInput) (bool, error) {
	key := input.TransactionHash + ":" + strconv.Itoa(input.LogIndex)
	if _, ok := store.events[key]; ok {
		return false, nil
	}
	store.events[key] = input
	return true, nil
}

func (store *fakeAnalyticsStore) UpdateAnalyticsIndexerState(ctx context.Context, input repository.UpdateAnalyticsIndexerStateInput) error {
	store.state = input
	return nil
}

func (store *fakeAnalyticsStore) RebuildAnalyticsSummaryCache(ctx context.Context, factoryAddress string) (repository.AnalyticsSummary, error) {
	summary := repository.AnalyticsSummary{
		Status:         repository.AnalyticsStatusOK,
		SourceStatus:   repository.AnalyticsSourceIndexed,
		FactoryAddress: factoryAddress,
		GeneratedAt:    time.Date(2026, 6, 28, 15, 0, 0, 0, time.UTC),
		LatestBlock:    sql.NullInt64{Int64: 49152802, Valid: true},
		LatestEventAt:  sql.NullTime{Time: time.Date(2026, 6, 28, 14, 53, 38, 0, time.UTC), Valid: true},
		Metrics: repository.AnalyticsMetrics{
			MarketsCreated:       int64(len(store.events)),
			MarketContractsFound: int64(len(store.markets)),
		},
	}
	store.summary = summary
	return summary, nil
}

func TestBackfillDryRunDoesNotWrite(t *testing.T) {
	store := newFakeAnalyticsStore()
	client := &fakeLogsClient{pages: []LogsPage{{Items: []BlockscoutLog{testMarketDeployedLog("0xtx1", 1)}}}}

	result, err := NewBackfiller(client, store).Run(context.Background(), BackfillOptions{
		FactoryAddress: "0xFactory",
		DryRun:         true,
		PageLimit:      1,
	})
	if err != nil {
		t.Fatalf("dry-run backfill: %v", err)
	}
	if result.EventsParsed != 1 {
		t.Fatalf("expected one parsed event, got %d", result.EventsParsed)
	}
	if len(store.markets) != 0 || len(store.events) != 0 {
		t.Fatalf("dry-run should not write, got markets=%d events=%d", len(store.markets), len(store.events))
	}
}

func TestBackfillIdempotentEventInsertAndSummary(t *testing.T) {
	store := newFakeAnalyticsStore()
	log := testMarketDeployedLog("0xtx1", 2)
	client := &fakeLogsClient{pages: []LogsPage{{Items: []BlockscoutLog{log, log}}}}

	result, err := NewBackfiller(client, store).Run(context.Background(), BackfillOptions{
		FactoryAddress: "0xFactory",
		DryRun:         false,
		PageLimit:      1,
		ChainID:        DefaultChainID,
	})
	if err != nil {
		t.Fatalf("write backfill: %v", err)
	}
	if result.EventsParsed != 2 {
		t.Fatalf("expected two parsed events, got %d", result.EventsParsed)
	}
	if result.EventsInserted != 1 {
		t.Fatalf("expected one inserted event after duplicate, got %d", result.EventsInserted)
	}
	if len(store.markets) != 1 {
		t.Fatalf("expected one market upsert, got %d", len(store.markets))
	}
	if result.Summary.Metrics.MarketsCreated != 1 {
		t.Fatalf("expected rebuilt summary markets_created=1, got %#v", result.Summary.Metrics)
	}
	if store.state.LastIndexedBlock != 49152802 {
		t.Fatalf("expected indexer state latest block, got %d", store.state.LastIndexedBlock)
	}
}

func TestBackfillSkipsLogsBeforeFromBlock(t *testing.T) {
	store := newFakeAnalyticsStore()
	client := &fakeLogsClient{pages: []LogsPage{{Items: []BlockscoutLog{testMarketDeployedLog("0xtx1", 1)}}}}

	result, err := NewBackfiller(client, store).Run(context.Background(), BackfillOptions{
		FactoryAddress: "0xFactory",
		FromBlock:      49152803,
		DryRun:         false,
		PageLimit:      1,
	})
	if err != nil {
		t.Fatalf("backfill: %v", err)
	}
	if result.EventsParsed != 0 || len(store.events) != 0 {
		t.Fatalf("expected old log to be skipped, result=%#v events=%d", result, len(store.events))
	}
}

func testMarketDeployedLog(txHash string, logIndex int) BlockscoutLog {
	return BlockscoutLog{
		BlockNumber:     49152802,
		BlockTimestamp:  "2026-06-28T14:53:38.000000Z",
		Index:           logIndex,
		TransactionHash: txHash,
		Raw:             json.RawMessage(`{"safe":"raw"}`),
		Decoded: &BlockscoutDecoded{
			MethodCall: "MarketDeployed(string indexed marketId, address indexed market, address indexed creator, address resolver, address collateralToken, uint256 closeTimestamp, string question)",
			Parameters: []BlockscoutParameter{
				{Name: "marketId", Value: "0xmarketidhash"},
				{Name: "market", Value: "0x09646deC03f5724C38BD486b0992A8CaF50Fcc59"},
				{Name: "creator", Value: "0xE2BB0d3445f5681994413879f5eF0802B4c2F624"},
				{Name: "resolver", Value: "0xE2BB0d3445f5681994413879f5eF0802B4c2F624"},
				{Name: "collateralToken", Value: "0x3600000000000000000000000000000000000000"},
				{Name: "closeTimestamp", Value: "1782658963"},
				{Name: "question", Value: "SignalArc test market"},
			},
		},
	}
}
