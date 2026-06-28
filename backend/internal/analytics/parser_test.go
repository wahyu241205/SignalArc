package analytics

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseMarketDeployed(t *testing.T) {
	raw := json.RawMessage(`{"transaction_hash":"0xabc"}`)
	log := BlockscoutLog{
		BlockNumber:     49152802,
		BlockTimestamp:  "2026-06-28T14:53:38.000000Z",
		Index:           8,
		TransactionHash: "0xa99473b770f5a708d814cf6609661097ef49cc5e8442875e4d40986d5fd4c4fc",
		Raw:             raw,
		Decoded: &BlockscoutDecoded{
			MethodCall: "MarketDeployed(string indexed marketId, address indexed market, address indexed creator, address resolver, address collateralToken, uint256 closeTimestamp, string question)",
			Parameters: []BlockscoutParameter{
				{Name: "marketId", Value: "0xbbc76d704eef842ab160610b9f45bf52e55d37a68386f9f342cc8bf932436b02"},
				{Name: "market", Value: "0x09646deC03f5724C38BD486b0992A8CaF50Fcc59"},
				{Name: "creator", Value: "0xE2BB0d3445f5681994413879f5eF0802B4c2F624"},
				{Name: "resolver", Value: "0xE2BB0d3445f5681994413879f5eF0802B4c2F624"},
				{Name: "collateralToken", Value: "0x3600000000000000000000000000000000000000"},
				{Name: "closeTimestamp", Value: "1782658963"},
				{Name: "question", Value: "SignalArc test market"},
			},
		},
	}

	event, matched, err := ParseMarketDeployed("0x02555FC5EE3c53938f2F0356e963865503442A56", log)
	if err != nil {
		t.Fatalf("parse MarketDeployed: %v", err)
	}
	if !matched {
		t.Fatal("expected MarketDeployed log to match")
	}
	if event.MarketAddress != "0x09646deC03f5724C38BD486b0992A8CaF50Fcc59" {
		t.Fatalf("unexpected market address %q", event.MarketAddress)
	}
	if event.CreatorAddress != "0xE2BB0d3445f5681994413879f5eF0802B4c2F624" {
		t.Fatalf("unexpected creator address %q", event.CreatorAddress)
	}
	if event.CloseTimestamp != "1782658963" {
		t.Fatalf("unexpected close timestamp %q", event.CloseTimestamp)
	}
	if event.BlockTimestamp.IsZero() {
		t.Fatal("expected parsed block timestamp")
	}
	if event.BlockTimestamp.Format(time.RFC3339) != "2026-06-28T14:53:38Z" {
		t.Fatalf("unexpected block timestamp %s", event.BlockTimestamp.Format(time.RFC3339))
	}
}

func TestParseMarketDeployedIgnoresOtherEvents(t *testing.T) {
	log := BlockscoutLog{
		Decoded: &BlockscoutDecoded{MethodCall: "PositionOpened(address indexed user, uint8 indexed side, uint256 amount)"},
	}

	_, matched, err := ParseMarketDeployed("0xfactory", log)
	if err != nil {
		t.Fatalf("expected unrelated event to be ignored without error: %v", err)
	}
	if matched {
		t.Fatal("expected unrelated event not to match")
	}
}
