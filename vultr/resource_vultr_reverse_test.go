package vultr

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v3"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func testGovultrClient(t *testing.T, handler func(*http.Request) (int, string)) *govultr.Client {
	t.Helper()

	client := govultr.NewClient(&http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			status, body := handler(r)
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
				Request:    r,
			}, nil
		}),
	})

	if err := client.SetBaseURL("https://example.com"); err != nil {
		t.Fatalf("error setting base URL: %v", err)
	}

	return client
}

func TestResourceVultrReverseIPV4ReadMissingParentInstance(t *testing.T) {
	client := testGovultrClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/v2/instances/missing-instance/ipv4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return http.StatusNotFound, `{"status":404,"error":"instance not found"}`
	})

	resource := resourceVultrReverseIPV4()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"instance_id": "missing-instance",
		"ip":          "192.0.2.10",
		"reverse":     "node.example.com",
	})
	data.SetId("192.0.2.10")

	diags := resourceVultrReverseIPV4Read(context.Background(), data, &Client{client: client})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if data.Id() != "" {
		t.Fatalf("expected resource ID to be cleared, got %q", data.Id())
	}
}

func TestResourceVultrReverseIPV4ReadMissingReverseRecord(t *testing.T) {
	client := testGovultrClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/v2/instances/existing-instance/ipv4" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return http.StatusOK, `{"ipv4s":[],"meta":{"links":{"next":""}}}`
	})

	resource := resourceVultrReverseIPV4()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"instance_id": "existing-instance",
		"ip":          "192.0.2.10",
		"reverse":     "node.example.com",
	})
	data.SetId("192.0.2.10")

	diags := resourceVultrReverseIPV4Read(context.Background(), data, &Client{client: client})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if data.Id() != "" {
		t.Fatalf("expected resource ID to be cleared, got %q", data.Id())
	}
}

func TestResourceVultrReverseIPV6ReadMissingParentInstance(t *testing.T) {
	client := testGovultrClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/v2/instances/missing-instance/ipv6/reverse" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return http.StatusNotFound, `{"status":404,"error":"instance not found"}`
	})

	resource := resourceVultrReverseIPV6()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"instance_id": "missing-instance",
		"ip":          "2001:db8::1",
		"reverse":     "node.example.com",
	})
	data.SetId("2001:0db8:0000:0000:0000:0000:0000:0001")

	diags := resourceVultrReverseIPV6Read(context.Background(), data, &Client{client: client})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if data.Id() != "" {
		t.Fatalf("expected resource ID to be cleared, got %q", data.Id())
	}
}
