package vultr

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

type mockLBService struct {
	govultr.LoadBalancerService
	lb  *govultr.LoadBalancer
	err error
}

func (m *mockLBService) Get(ctx context.Context, lbID string) (*govultr.LoadBalancer, *http.Response, error) {
	if m.err != nil {
		return nil, nil, m.err
	}
	return m.lb, nil, nil
}

func TestNewLBStateRefresh_IPv4Pending(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceVultrLoadBalancer().Schema, map[string]interface{}{})
	d.SetId("test-lb-id")

	// Case 1: Status is active but IPV4 is empty -> should report "installing" to continue waiting
	mockSvc := &mockLBService{
		lb: &govultr.LoadBalancer{
			ID:     "test-lb-id",
			Status: "active",
			IPV4:   "",
		},
	}
	meta := &Client{client: &govultr.Client{LoadBalancer: mockSvc}}

	refresh := newLBStateRefresh(context.Background(), d, meta, "status")
	res, state, err := refresh()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != "installing" {
		t.Errorf("expected state 'installing' while IPV4 is empty, got '%s'", state)
	}
	if res == nil {
		t.Errorf("expected non-nil resource")
	}

	// Case 2: Status is active and IPV4 is assigned -> should report "active"
	mockSvc.lb.IPV4 = "192.0.2.10"
	res, state, err = refresh()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != "active" {
		t.Errorf("expected state 'active' when IPV4 is assigned, got '%s'", state)
	}
	if res == nil {
		t.Errorf("expected non-nil resource")
	}

	// Case 3: Status is pending -> should report "pending"
	mockSvc.lb.Status = "pending"
	mockSvc.lb.IPV4 = ""
	res, state, err = refresh()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != "pending" {
		t.Errorf("expected state 'pending', got '%s'", state)
	}
}
