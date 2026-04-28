package mihomo

import (
	"sublink/models"
	"testing"

	"github.com/metacubex/mihomo/constant"
)

func TestBuildQualityURLCandidatesIncludesFallbacksWithoutDuplicates(t *testing.T) {
	candidates := buildQualityURLCandidates("https://my.123169.xyz/v1/info")
	if len(candidates) != 2 {
		t.Fatalf("expected 2 unique candidates, got %d: %#v", len(candidates), candidates)
	}
	if candidates[0] != "https://my.123169.xyz/v1/info" {
		t.Fatalf("expected primary candidate first, got %s", candidates[0])
	}
	if candidates[1] != "https://my.ippure.com/v1/info" {
		t.Fatalf("expected fallback candidate second, got %s", candidates[1])
	}
}

func TestNormalizeQualityResultWithLandingIPDowngradesFailedResult(t *testing.T) {
	result := normalizeQualityResultWithLandingIP(&QualityCheckResult{
		Status: models.QualityStatusFailed,
		Reason: "context deadline exceeded",
	}, "1.1.1.1")

	if result == nil {
		t.Fatal("expected normalized result, got nil")
	}
	if result.Status != models.QualityStatusPartial {
		t.Fatalf("expected partial status, got %s", result.Status)
	}
	if result.Family != models.QualityFamilyIPv4 {
		t.Fatalf("expected ipv4 family, got %s", result.Family)
	}
	if result.IP != "1.1.1.1" {
		t.Fatalf("expected landing ip to be preserved, got %s", result.IP)
	}
	if result.Reason != "quality_api_unreachable" {
		t.Fatalf("expected normalized reason, got %s", result.Reason)
	}
}

func TestNormalizeQualityResultWithLandingIPKeepsSuccessfulResult(t *testing.T) {
	original := &QualityCheckResult{
		Status: models.QualityStatusSuccess,
		IP:     "8.8.8.8",
	}

	result := normalizeQualityResultWithLandingIP(original, "1.1.1.1")
	if result != original {
		t.Fatal("expected successful result to remain unchanged")
	}
}

func TestFetchQualityWithAdapterPrefersLaterSuccessOverEarlierPartial(t *testing.T) {
	originalFetcher := qualityFetcher
	defer func() { qualityFetcher = originalFetcher }()

	results := map[string]*QualityCheckResult{
		"https://my.123169.xyz/v1/info": {Status: models.QualityStatusPartial, Reason: "missing_quality_fields"},
		"https://my.ippure.com/v1/info": {
			Status:         models.QualityStatusSuccess,
			IsBroadcast:    true,
			IsResidential:  false,
			FraudScore:     30,
			HasBroadcast:   true,
			HasResidential: true,
			HasFraudScore:  true,
		},
	}
	qualityFetcher = func(_ constant.Proxy, qualityURL string) *QualityCheckResult {
		return results[qualityURL]
	}

	result := FetchQualityWithAdapter(nil, "https://my.123169.xyz/v1/info")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Status != models.QualityStatusSuccess {
		t.Fatalf("expected success result, got %s", result.Status)
	}
	if result.FraudScore != 30 || !result.IsBroadcast {
		t.Fatalf("expected full success payload, got %+v", result)
	}
}

func TestFetchQualityWithAdapterFallsBackToPartialWhenNoSuccessExists(t *testing.T) {
	originalFetcher := qualityFetcher
	defer func() { qualityFetcher = originalFetcher }()

	results := map[string]*QualityCheckResult{
		"https://my.123169.xyz/v1/info": {Status: models.QualityStatusPartial, Reason: "missing_quality_fields"},
		"https://my.ippure.com/v1/info": {Status: models.QualityStatusFailed, Reason: "status_500"},
	}
	qualityFetcher = func(_ constant.Proxy, qualityURL string) *QualityCheckResult {
		return results[qualityURL]
	}

	result := FetchQualityWithAdapter(nil, "https://my.123169.xyz/v1/info")
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Status != models.QualityStatusPartial {
		t.Fatalf("expected partial result, got %s", result.Status)
	}
	if result.Reason != "missing_quality_fields" {
		t.Fatalf("expected partial reason to be preserved, got %s", result.Reason)
	}
}
