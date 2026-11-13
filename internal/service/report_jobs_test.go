package service

import (
	"net/url"
	"testing"
)

func TestReportCacheKeyIgnoresParamOrderAndNoCache(t *testing.T) {
	first := url.Values{"status": {"closed"}, "format": {"html"}, "page": {"1"}}
	second := url.Values{"page": {"1"}, "format": {"html"}, "status": {"closed"}}

	key := reportCacheKey("journal", "html", first)
	if key != reportCacheKey("journal", "html", second) {
		t.Fatal("cache key must not depend on parameter order")
	}
	if key == reportCacheKey("journal", "xlsx", first) {
		t.Fatal("cache key must include format")
	}
	if key == reportCacheKey("talons", "html", first) {
		t.Fatal("cache key must include report code")
	}

	withNoCache := url.Values{"status": {"closed"}, "format": {"html"}, "page": {"1"}, "nocache": {"1"}}
	if key != reportCacheKey("journal", "html", withNoCache) {
		t.Fatal("nocache must not affect the key")
	}
}
