package vultr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCredentialAttributesAreSensitive(t *testing.T) {
	cases := []struct {
		resource string
		schema   map[string]*schema.Schema
		attrs    []string
	}{
		{"vultr_database", resourceVultrDatabase().Schema, []string{"password", "ferretdb_credentials", "access_key"}},
		{"vultr_database.read_replicas", resourceVultrDatabase().Schema["read_replicas"].Elem.(*schema.Resource).Schema, []string{"password", "ferretdb_credentials"}},
		{"vultr_database_replica", resourceVultrDatabaseReplica().Schema, []string{"password", "ferretdb_credentials"}},
		{"vultr_database_user", resourceVultrDatabaseUser().Schema, []string{"password", "access_key"}},
		{"vultr_user", resourceVultrUsers().Schema, []string{"api_key"}},
		{"vultr_oidc_token", resourceVultrOIDCToken().Schema, []string{"refresh_token", "access_token", "id_token"}},
		{"vultr_organization_role_session", resourceVultrOrganizationRoleSession().Schema, []string{"token"}},
	}

	for _, tc := range cases {
		for _, attr := range tc.attrs {
			s, ok := tc.schema[attr]
			if !ok {
				t.Errorf("%s: attribute %q not found", tc.resource, attr)
				continue
			}
			if !s.Sensitive {
				t.Errorf("%s: attribute %q must be marked Sensitive", tc.resource, attr)
			}
		}
	}
}
