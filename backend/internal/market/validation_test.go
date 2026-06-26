package market

import (
	"strings"
	"testing"
	"time"
)

func TestCreateMarketRequestAcceptsOptionalHTTPSCoverImageURL(t *testing.T) {
	coverImageURL := " https://example.com/market-cover.png "
	request := validCreateMarketRequest()
	request.CoverImageURL = &coverImageURL

	input, err := request.ToRepositoryInput(time.Now())
	if err != nil {
		t.Fatalf("expected valid cover image URL, got %v", err)
	}
	if !input.CoverImageURL.Valid {
		t.Fatal("expected cover image URL to be stored")
	}
	if input.CoverImageURL.String != strings.TrimSpace(coverImageURL) {
		t.Fatalf("expected trimmed cover image URL, got %q", input.CoverImageURL.String)
	}
}

func TestCreateMarketRequestAllowsMissingCoverImageURL(t *testing.T) {
	request := validCreateMarketRequest()

	input, err := request.ToRepositoryInput(time.Now())
	if err != nil {
		t.Fatalf("expected missing cover image URL to be valid, got %v", err)
	}
	if input.CoverImageURL.Valid {
		t.Fatalf("expected missing cover image URL to remain null, got %q", input.CoverImageURL.String)
	}
}

func TestCreateMarketRequestRejectsInvalidCoverImageURL(t *testing.T) {
	testCases := map[string]string{
		"http scheme":  "http://example.com/market-cover.png",
		"base64 data":  "data:image/png;base64,abc123",
		"missing host": "https:///market-cover.png",
		"not a URL":    "market-cover.png",
		"too long":     "https://example.com/" + strings.Repeat("a", 2049),
	}

	for name, coverImageURL := range testCases {
		coverImageURL := coverImageURL
		t.Run(name, func(t *testing.T) {
			request := validCreateMarketRequest()
			request.CoverImageURL = &coverImageURL

			if _, err := request.ToRepositoryInput(time.Now()); err == nil {
				t.Fatalf("expected cover image URL %q to be rejected", coverImageURL)
			}
		})
	}
}

func validCreateMarketRequest() CreateMarketRequest {
	return CreateMarketRequest{
		CreatorUserID: "10000000-0000-4000-8000-000000000001",
		Title:         "Will SignalArc support market images?",
		Chain:         "Arc Testnet",
		ClosesAt:      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
}
