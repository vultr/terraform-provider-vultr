package vultr

import (
	"context"
	"testing"
)

func TestResourceVultrKubernetesStateUpgradeV0ToV1_NodePoolsNil(t *testing.T) {
	rawState := map[string]interface{}{
		"id":     "test-cluster-id",
		"label":  "test-cluster",
		"region": "ewr",
	}

	upgraded, err := resourceVultrKubernetesStateUpgradeV0ToV1(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if upgraded["id"] != "test-cluster-id" {
		t.Errorf("expected id to be preserved, got: %v", upgraded["id"])
	}
}

func TestResourceVultrKubernetesStateUpgradeV0ToV1_NodePoolsEmpty(t *testing.T) {
	rawState := map[string]interface{}{
		"id":         "test-cluster-id",
		"label":      "test-cluster",
		"node_pools": []interface{}{},
	}

	upgraded, err := resourceVultrKubernetesStateUpgradeV0ToV1(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if upgraded["id"] != "test-cluster-id" {
		t.Errorf("expected id to be preserved, got: %v", upgraded["id"])
	}
}

func TestResourceVultrKubernetesStateUpgradeV0ToV1_EmptyRawState(t *testing.T) {
	rawState := map[string]interface{}{}

	upgraded, err := resourceVultrKubernetesStateUpgradeV0ToV1(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(upgraded) != 0 {
		t.Errorf("expected empty map, got: %v", upgraded)
	}
}

func TestResourceVultrKubernetesNodePoolsStateUpgradeV0ToV1_LabelsNil(t *testing.T) {
	rawState := map[string]interface{}{
		"id":         "test-np-id",
		"cluster_id": "test-cluster-id",
		"labels":     nil,
		"taints":     nil,
	}

	upgraded, err := resourceVultrKubernetesNodePoolsStateUpgradeV0ToV1(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if upgraded["id"] != "test-np-id" {
		t.Errorf("expected id to be preserved, got: %v", upgraded["id"])
	}
}

func TestResourceVultrKubernetesNodePoolsStateUpgradeV0ToV1_LabelsEmpty(t *testing.T) {
	rawState := map[string]interface{}{
		"id":         "test-np-id",
		"cluster_id": "test-cluster-id",
		"labels":     map[string]interface{}{},
		"taints":     []interface{}{},
	}

	upgraded, err := resourceVultrKubernetesNodePoolsStateUpgradeV0ToV1(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if upgraded["id"] != "test-np-id" {
		t.Errorf("expected id to be preserved, got: %v", upgraded["id"])
	}
}

func TestResourceVultrKubernetesNodePoolsStateUpgradeV0ToV1_EmptyRawState(t *testing.T) {
	rawState := map[string]interface{}{}

	upgraded, err := resourceVultrKubernetesNodePoolsStateUpgradeV0ToV1(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(upgraded) != 0 {
		t.Errorf("expected empty map, got: %v", upgraded)
	}
}

func TestResourceVultrKubernetesNodePoolsStateUpgradeV0ToV1_TaintsNil(t *testing.T) {
	rawState := map[string]interface{}{
		"id":         "test-np-id",
		"cluster_id": "test-cluster-id",
		"labels":     nil,
		"taints":     nil,
	}

	upgraded, err := resourceVultrKubernetesNodePoolsStateUpgradeV0ToV1(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if upgraded["id"] != "test-np-id" {
		t.Errorf("expected id to be preserved, got: %v", upgraded["id"])
	}
	// labels should be set to empty slice when nil
	if upgraded["labels"] == nil {
		t.Error("expected labels to be set to empty slice, got nil")
	}
}
