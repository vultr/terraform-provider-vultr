package vultr

import (
	"testing"
)

func TestSuppressIPDiff(t *testing.T) {
	cases := []struct {
		name     string
		old      string
		new      string
		suppress bool
	}{
		// leading zero in v6 group
		{
			name:     "v6 leading zero in group",
			old:      "2001:db8:1000:3b79:5400:5ff:fedf:fade",
			new:      "2001:db8:1000:3b79:5400:05ff:fedf:fade",
			suppress: true,
		},
		// same addr both sides -- no-op
		{
			name:     "v6 identical",
			old:      "2001:db8::1",
			new:      "2001:db8::1",
			suppress: true,
		},
		// full expansion vs compressed -- same addr
		{
			name:     "v6 expanded vs compressed",
			old:      "2001:db8:0:0:0:0:0:1",
			new:      "2001:db8::1",
			suppress: true,
		},
		// actually different addrs
		{
			name:     "v6 different addrs",
			old:      "2001:db8::1",
			new:      "2001:db8::2",
			suppress: false,
		},
		// v4 sanity -- should be a noop but make sure it doesn't break
		{
			name:     "v4 identical",
			old:      "10.0.0.1",
			new:      "10.0.0.1",
			suppress: true,
		},
		{
			name:     "v4 different",
			old:      "10.0.0.1",
			new:      "10.0.0.2",
			suppress: false,
		},
		// garbage in -- fall through to string compare
		{
			name:     "unparseable falls through",
			old:      "not-an-ip",
			new:      "also-not-an-ip",
			suppress: false,
		},
		// mixed case hex -- same addr
		{
			name:     "v6 mixed case",
			old:      "2001:db8:1000:3b79:5400:5FF:FEDF:FADE",
			new:      "2001:db8:1000:3b79:5400:5ff:fedf:fade",
			suppress: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := suppressIPDiff("subnet", tc.old, tc.new, nil)
			if got != tc.suppress {
				t.Errorf("suppressIPDiff(%q, %q) = %v, want %v",
					tc.old, tc.new, got, tc.suppress)
			}
		})
	}
}

func TestCanonicalizeIP(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		// the actual bug case
		{
			name: "strips leading zero from v6 group",
			in:   "2001:db8:1000:3b79:5400:05ff:fedf:fade",
			want: "2001:db8:1000:3b79:5400:5ff:fedf:fade",
		},
		// already canonical -- passthrough
		{
			name: "already canonical v6",
			in:   "2001:db8::1",
			want: "2001:db8::1",
		},
		// v4 passthrough
		{
			name: "v4 passthrough",
			in:   "10.0.0.1",
			want: "10.0.0.1",
		},
		// full form gets compressed
		{
			name: "v6 full form compresses",
			in:   "2001:0db8:0000:0000:0000:0000:0000:0001",
			want: "2001:db8::1",
		},
		// bad input -- pass through unchanged
		{
			name: "bad input passthrough",
			in:   "not-an-ip",
			want: "not-an-ip",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := canonicalizeIP(tc.in)
			if got != tc.want {
				t.Errorf("canonicalizeIP(%q) = %q, want %q",
					tc.in, got, tc.want)
			}
		})
	}
}

func TestCanonicalizeUserDataLineEndings(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "lf unchanged",
			in:   "#cloud-config\nmanage_etc_hosts: false\n",
			want: "#cloud-config\nmanage_etc_hosts: false\n",
		},
		{
			name: "crlf normalized",
			in:   "#cloud-config\r\nmanage_etc_hosts: false\r\n",
			want: "#cloud-config\nmanage_etc_hosts: false\n",
		},
		{
			name: "cr normalized",
			in:   "#cloud-config\rmanage_etc_hosts: false\r",
			want: "#cloud-config\nmanage_etc_hosts: false\n",
		},
		{
			name: "content preserved",
			in:   "#cloud-config\nusers:\n  - name: terraform\n",
			want: "#cloud-config\nusers:\n  - name: terraform\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := canonicalizeUserDataLineEndings(tc.in)
			if got != tc.want {
				t.Errorf("canonicalizeUserDataLineEndings(%q) = %q, want %q",
					tc.in, got, tc.want)
			}
		})
	}
}
