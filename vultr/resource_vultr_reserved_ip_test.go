package vultr

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vultr/govultr/v3"
)

func TestAccVultrReservedIPIPv4(t *testing.T) {
	rServerLabel := acctest.RandomWithPrefix("tf-vps-rip4")
	rLabel := acctest.RandomWithPrefix("tf-rip4-rs")
	rLabelUpdated := rLabel + "_updated"
	ipType := "v4"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
			{
				Config: testAccVultrReservedIPConfigAttach(rServerLabel, rLabelUpdated, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabelUpdated),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "instance_id"),
				),
			},
			{
				// test detach by unsetting the attached_id
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
		},
	})
}

func TestAccVultrReservedIPIPv6(t *testing.T) {
	rServerLabel := acctest.RandomWithPrefix("tf-vps-rip6")
	rLabel := acctest.RandomWithPrefix("tf-rip6-rs")
	ipType := "v6"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
			{
				Config: testAccVultrReservedIPConfigAttach(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "instance_id"),
				),
			},
			{
				// test detach by unsetting the attached_id
				Config: testAccVultrReservedIPConfig(rServerLabel, rLabel, ipType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "ip_type", ipType),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "region"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet"),
					resource.TestCheckResourceAttrSet("vultr_reserved_ip.foo", "subnet_size"),
				),
			},
		},
	})
}

func TestAccVultrReservedIPLabelUpdate(t *testing.T) {
	rLabel := acctest.RandomWithPrefix("tf-rip-rs")
	rLabelUpdated := rLabel + "_updated"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVultrReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVultrReservedIPConfigLabel(rLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabel),
				),
			},
			{
				Config: testAccVultrReservedIPConfigLabel(rLabelUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVultrReservedIPExists("vultr_reserved_ip.foo"),
					resource.TestCheckResourceAttr("vultr_reserved_ip.foo", "label", rLabelUpdated),
				),
			},
		},
	})
}

func testAccCheckVultrReservedIPDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vultr_reserved_ip" {
			continue
		}

		ripID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		_, _, err := client.ReservedIP.Get(context.Background(), ripID)
		if err == nil {
			return fmt.Errorf("reserved IP still exists: %s", ripID)
		}
	}
	return nil
}

func testAccCheckVultrReservedIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("reserved IP ID is not set")
		}

		ripID := rs.Primary.ID
		client := testAccProvider.Meta().(*Client).govultrClient()

		if _, _, err := client.ReservedIP.Get(context.Background(), ripID); err != nil {
			return fmt.Errorf("reserved IP does not exist: %s", ripID)
		}

		return nil
	}
}

func testAccVultrReservedIPConfigLabel(label string) string {
	return fmt.Sprintf(`

   resource "vultr_reserved_ip" "foo" {
       label     = "%s"
       region    = "ewr"
       ip_type   = "v4"
   }`, label)
}

func testAccVultrReservedIPConfig(rServerLabel, label, ipType string) string {
	return fmt.Sprintf(`
	resource "vultr_instance" "ip" {
       label       = "%s"
       region      = "ewr"
       plan        = "vc2-1c-2gb"
       os_id       = 167
	   enable_ipv6 = true
   }
   resource "vultr_reserved_ip" "foo" {
       label     = "%s"
       region    = "ewr"
       ip_type   = "%s"
   }`, rServerLabel, label, ipType)
}

func testAccVultrReservedIPConfigAttach(rServerLabel, label, ipType string) string {
	return fmt.Sprintf(`
	resource "vultr_instance" "ip" {
       label = "%s"
       region = "ewr"
       plan = "vc2-1c-2gb"
       os_id = 167
       enable_ipv6 = true
   }
   resource "vultr_reserved_ip" "foo" {
       label       = "%s"
       region      = "ewr"
       ip_type     = "%s"
       instance_id = "${vultr_instance.ip.id}"
   }`, rServerLabel, label, ipType)
}

func TestResourceVultrReservedIPRead(t *testing.T) {
	const foundBody = `{"reserved_ip":{"id":"rip-1","region":"ewr","ip_type":"v4","subnet":"192.0.2.10",` +
		`"subnet_size":32,"label":"example","instance_id":"inst-1"}}`

	tests := []struct {
		name        string
		status      int
		body        string
		contentType string
		wantErr     string
		wantID      string
		wantState   map[string]interface{}
	}{
		{
			name:        "found",
			status:      http.StatusOK,
			body:        foundBody,
			contentType: "application/json",
			wantID:      "rip-1",
			wantState: map[string]interface{}{
				"region":      "ewr",
				"ip_type":     "v4",
				"subnet":      "192.0.2.10",
				"subnet_size": 32,
				"label":       "example",
				"instance_id": "inst-1",
			},
		},
		{
			name:        "empty body",
			status:      http.StatusOK,
			body:        `{}`,
			contentType: "application/json",
			wantErr:     "empty response",
			wantID:      "rip-1",
		},
		{
			name:        "not found",
			status:      http.StatusNotFound,
			body:        `{"error":"reserved ip not found","status":404}`,
			contentType: "application/json",
			wantID:      "",
		},
		{
			name:        "invalid id",
			status:      http.StatusBadRequest,
			body:        `{"error":"Invalid reserved-ip ID","status":400}`,
			contentType: "application/json",
			wantErr:     "Invalid reserved-ip ID",
			wantID:      "rip-1",
		},
		{
			name:        "server error",
			status:      http.StatusInternalServerError,
			body:        `{"error":"internal","status":500}`,
			contentType: "application/json",
			wantErr:     "internal",
			wantID:      "rip-1",
		},
		{
			name:        "not json",
			status:      http.StatusServiceUnavailable,
			body:        `upstream unavailable`,
			contentType: "text/plain",
			wantErr:     "unable to unmarshal api response",
			wantID:      "rip-1",
		},
		{
			name:        "server failed to include a response",
			status:      http.StatusOK,
			body:        `{"error":"(REMOTE) Server failed to include a response","error_type":"server"}`,
			contentType: "application/json",
			wantErr:     "empty response",
			wantID:      "rip-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/v2/reserved-ips/rip-1" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(tt.status)
				if _, err := w.Write([]byte(tt.body)); err != nil {
					t.Errorf("write body: %v", err)
				}
			}))
			defer srv.Close()

			c := govultr.NewClient(nil)
			if err := c.SetBaseURL(srv.URL); err != nil {
				t.Fatalf("set base url: %v", err)
			}
			c.SetRetryLimit(0)
			meta := &Client{client: c}
			d := resourceVultrReservedIP().TestResourceData()
			d.SetId("rip-1")

			diags := resourceVultrReservedIPRead(context.Background(), d, meta)

			if tt.wantErr == "" {
				if diags.HasError() {
					t.Fatalf("unexpected error: %v", diags)
				}
			} else {
				if !diags.HasError() {
					t.Fatalf("expected an error containing %q, got none", tt.wantErr)
				}
				var summaries []string
				for _, diag := range diags {
					summaries = append(summaries, diag.Summary)
				}
				if got := strings.Join(summaries, "\n"); !strings.Contains(got, tt.wantErr) {
					t.Errorf("error summary %q does not contain %q", got, tt.wantErr)
				}
			}

			if d.Id() != tt.wantID {
				t.Errorf("id = %q, want %q", d.Id(), tt.wantID)
			}
			for key, want := range tt.wantState {
				if got := d.Get(key); got != want {
					t.Errorf("%s = %v, want %v", key, got, want)
				}
			}
		})
	}
}
