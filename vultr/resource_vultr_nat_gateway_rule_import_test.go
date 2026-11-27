package vultr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

const (
	testImportVPCID  = "11111111-1111-1111-1111-111111111111"
	testImportNATID  = "22222222-2222-2222-2222-222222222222"
	testImportRuleID = "33333333-3333-3333-3333-333333333333"
)

func TestNATGatewayRuleImport(t *testing.T) {
	cases := map[string]struct {
		resource *schema.Resource
		path     string
		body     string
	}{
		"firewall_rule": {
			resource: resourceVultrNATGatewayFirewallRule(),
			path:     "/v2/vpcs/" + testImportVPCID + "/nat-gateway/" + testImportNATID + "/global/firewall-rules/" + testImportRuleID,
			body:     `{"firewall_rule":{"id":"` + testImportRuleID + `","action":"accept","protocol":"tcp","port":"80","subnet":"10.0.0.0","subnet_size":24}}`,
		},
		"port_forwarding_rule": {
			resource: resourceVultrNATGatewayPortForwardingRule(),
			path:     "/v2/vpcs/" + testImportVPCID + "/nat-gateway/" + testImportNATID + "/global/port-forwarding-rules/" + testImportRuleID,
			body:     `{"port_forwarding_rule":{"id":"` + testImportRuleID + `","name":"web","protocol":"tcp","external_port":80,"internal_ip":"10.0.0.5","internal_port":8080,"enabled":true}}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path != tc.path {
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write([]byte(`{"error":"unexpected path ` + r.URL.Path + `","status":404}`))
					return
				}
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			client := govultr.NewClient(server.Client())
			if err := client.SetBaseURL(server.URL); err != nil {
				t.Fatal(err)
			}
			meta := &Client{client: client}

			d := tc.resource.TestResourceData()
			d.SetId(testImportVPCID + "_" + testImportNATID + "_" + testImportRuleID)

			imported, err := tc.resource.Importer.StateContext(context.Background(), d, meta)
			if err != nil {
				t.Fatalf("import failed: %v", err)
			}
			if len(imported) != 1 {
				t.Fatalf("expected 1 imported resource, got %d", len(imported))
			}

			d = imported[0]
			if d.Id() != testImportRuleID {
				t.Fatalf("expected id %q, got %q", testImportRuleID, d.Id())
			}
			if got := d.Get("vpc_id").(string); got != testImportVPCID {
				t.Fatalf("expected vpc_id %q, got %q", testImportVPCID, got)
			}
			if got := d.Get("nat_gateway_id").(string); got != testImportNATID {
				t.Fatalf("expected nat_gateway_id %q, got %q", testImportNATID, got)
			}

			if diags := tc.resource.ReadContext(context.Background(), d, meta); diags.HasError() {
				t.Fatalf("read after import failed: %v", diags)
			}
			if got := d.Get("protocol").(string); got != "tcp" {
				t.Fatalf("expected protocol to be read back as tcp, got %q", got)
			}
		})
	}
}

func TestNATGatewayRuleImportInvalidID(t *testing.T) {
	for _, id := range []string{
		testImportRuleID,
		testImportVPCID + "_" + testImportRuleID,
		testImportVPCID + "__" + testImportRuleID,
	} {
		for name, r := range map[string]*schema.Resource{
			"firewall_rule":        resourceVultrNATGatewayFirewallRule(),
			"port_forwarding_rule": resourceVultrNATGatewayPortForwardingRule(),
		} {
			d := r.TestResourceData()
			d.SetId(id)
			if _, err := r.Importer.StateContext(context.Background(), d, &Client{}); err == nil {
				t.Fatalf("%s: expected an error for import id %q", name, id)
			}
		}
	}
}
