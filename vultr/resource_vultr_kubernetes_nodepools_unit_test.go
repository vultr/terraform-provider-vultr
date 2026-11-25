package vultr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func TestExpandNodePoolCreateReqVPCOnly(t *testing.T) {
	for _, vpcOnly := range []bool{true, false} {
		d := schema.TestResourceDataRaw(t, resourceVultrKubernetesNodePools().Schema, map[string]interface{}{
			"cluster_id":    "cluster",
			"label":         "pool",
			"plan":          "vc2-2c-4gb",
			"node_quantity": 1,
			"vpc_only":      vpcOnly,
		})

		req := expandNodePoolCreateReq(d)
		if req.VPCOnly == nil {
			t.Fatalf("vpc_only = %v: VPCOnly was not set on the create request", vpcOnly)
		}
		if *req.VPCOnly != vpcOnly {
			t.Fatalf("VPCOnly = %v, want %v", *req.VPCOnly, vpcOnly)
		}
	}
}

func TestResourceVultrKubernetesNodePoolsReadVPCOnly(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/kubernetes/clusters/cluster/node-pools/pool", jsonResponse(`{"node_pool":{"id":"pool","label":"pool","plan":"vc2-2c-4gb","status":"active","node_quantity":1,"vpc_only":true,"nodes":[]}}`))
	mux.HandleFunc("/v2/kubernetes/clusters/cluster/node-pools/pool/labels", jsonResponse(`{"labels":[]}`))
	mux.HandleFunc("/v2/kubernetes/clusters/cluster/node-pools/pool/taints", jsonResponse(`{"taints":[]}`))
	server := httptest.NewServer(mux)
	defer server.Close()

	client := govultr.NewClient(server.Client())
	if err := client.SetBaseURL(server.URL); err != nil {
		t.Fatal(err)
	}

	// Simulate an import: only the IDs are known, vpc_only is not in state yet.
	d := schema.TestResourceDataRaw(t, resourceVultrKubernetesNodePools().Schema, map[string]interface{}{
		"cluster_id": "cluster",
	})
	d.SetId("pool")

	if diags := resourceVultrKubernetesNodePoolsRead(context.Background(), d, &Client{client: client}); diags.HasError() {
		t.Fatalf("unexpected read error: %v", diags)
	}

	if got := d.Get("vpc_only").(bool); !got {
		t.Fatalf("vpc_only = %v, want true", got)
	}
}

func jsonResponse(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}
}
