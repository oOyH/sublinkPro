package scheduler

import (
	"sublink/models"
	"sublink/services/mihomo"
	"testing"
)

func TestApplyNodeQualityInfoPreservesPartialFields(t *testing.T) {
	node := &models.Node{}
	quality := &mihomo.QualityCheckResult{
		Status:         models.QualityStatusPartial,
		Family:         models.QualityFamilyIPv4,
		IsBroadcast:    true,
		FraudScore:     12,
		HasBroadcast:   true,
		HasResidential: false,
		HasFraudScore:  true,
		Reason:         "missing_quality_fields",
	}

	applyNodeQualityInfo(node, quality)

	if node.QualityStatus != models.QualityStatusPartial {
		t.Fatalf("expected partial status, got %s", node.QualityStatus)
	}
	if node.QualityFamily != models.QualityFamilyIPv4 {
		t.Fatalf("expected ipv4 family, got %s", node.QualityFamily)
	}
	if !node.IsBroadcast {
		t.Fatal("expected broadcast flag to be preserved")
	}
	if node.QualityHasResidential {
		t.Fatal("expected residential availability to remain false")
	}
	if node.IsResidential {
		t.Fatal("expected residential flag to stay false when field is missing")
	}
	if node.FraudScore != 12 {
		t.Fatalf("expected fraud score 12, got %d", node.FraudScore)
	}
	if !node.QualityHasBroadcast || !node.QualityHasFraudScore {
		t.Fatal("expected available partial fields to be marked as present")
	}
	if node.QualityReason != "missing_quality_fields" {
		t.Fatalf("expected quality reason to be preserved, got %s", node.QualityReason)
	}
}
