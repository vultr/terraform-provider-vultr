package vultr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

func TestDatabaseSubResourcesReadParentGone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"invalid database ID","status":404}`))
	}))
	defer server.Close()

	client := govultr.NewClient(server.Client())
	if err := client.SetBaseURL(server.URL); err != nil {
		t.Fatal(err)
	}
	meta := &Client{client: client}

	cases := map[string]struct {
		resource *schema.Resource
		read     schema.ReadContextFunc
		id       string
	}{
		"database_user":            {resourceVultrDatabaseUser(), resourceVultrDatabaseUserRead, "user"},
		"database_db":              {resourceVultrDatabaseDB(), resourceVultrDatabaseDBRead, "db"},
		"database_topic":           {resourceVultrDatabaseTopic(), resourceVultrDatabaseTopicRead, "topic"},
		"database_connection_pool": {resourceVultrDatabaseConnectionPool(), resourceVultrDatabaseConnectionPoolRead, "pool"},
		"database_quota":           {resourceVultrDatabaseQuota(), resourceVultrDatabaseQuotaRead, "client|user"},
		"database_connector":       {resourceVultrDatabaseConnector(), resourceVultrDatabaseConnectorRead, "connector"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, tc.resource.Schema, map[string]interface{}{
				"database_id": "11111111-2222-3333-4444-555555555555",
			})
			d.SetId(tc.id)

			if diags := tc.read(context.Background(), d, meta); diags.HasError() {
				t.Fatalf("expected no error when the parent database is gone, got: %v", diags)
			}
			if d.Id() != "" {
				t.Fatalf("expected resource to be removed from state, id is still %q", d.Id())
			}
		})
	}
}
