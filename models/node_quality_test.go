package models

import "testing"

func TestNormalizeNodeForImportBackfillsQualityAvailabilityForSuccess(t *testing.T) {
	node := &Node{
		FraudScore:    18,
		QualityStatus: QualityStatusSuccess,
	}

	NormalizeNodeForImport(node)

	if !node.QualityHasBroadcast || !node.QualityHasResidential || !node.QualityHasFraudScore {
		t.Fatalf("expected success node to backfill all quality availability flags: %+v", node)
	}
}

func TestGetNodeIPTypeValueUsesPartialAvailability(t *testing.T) {
	node := Node{
		QualityStatus:       QualityStatusPartial,
		QualityHasBroadcast: true,
		IsBroadcast:         false,
	}

	if got := getNodeIPTypeValue(node); got != "native" {
		t.Fatalf("expected native for partial node with broadcast info, got %s", got)
	}
}

func TestGetNodeResidentialTypeValueUsesPartialAvailability(t *testing.T) {
	node := Node{
		QualityStatus:         QualityStatusPartial,
		QualityHasResidential: true,
		IsResidential:         true,
	}

	if got := getNodeResidentialTypeValue(node); got != "residential" {
		t.Fatalf("expected residential for partial node with residential info, got %s", got)
	}
}
