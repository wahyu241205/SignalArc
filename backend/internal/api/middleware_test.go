package api

import "testing"

func TestIsCORSOriginAllowedIncludesLocalDefaults(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	if !isCORSOriginAllowed("http://localhost:3000") {
		t.Fatal("expected localhost frontend origin to be allowed")
	}
	if !isCORSOriginAllowed("http://127.0.0.1:3000") {
		t.Fatal("expected loopback frontend origin to be allowed")
	}
}

func TestIsCORSOriginAllowedUsesConfiguredExactOrigins(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://example.signalarc.fun, https://agents.signalarc.fun")

	if !isCORSOriginAllowed("https://example.signalarc.fun") {
		t.Fatal("expected configured origin to be allowed")
	}
	if !isCORSOriginAllowed("https://agents.signalarc.fun") {
		t.Fatal("expected configured origin with surrounding whitespace to be allowed")
	}
}

func TestIsCORSOriginAllowedDeniesUnknownAndWildcardOrigins(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "*,https://example.signalarc.fun")

	if isCORSOriginAllowed("https://evil.example") {
		t.Fatal("expected unknown origin to be denied")
	}
	if isCORSOriginAllowed("*") {
		t.Fatal("expected wildcard origin to be denied")
	}
}
