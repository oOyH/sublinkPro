package mihomo

import (
	"sublink/models"
	"testing"
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

