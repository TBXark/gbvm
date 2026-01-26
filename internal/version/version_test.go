package version

import (
	"reflect"
	"testing"

	"github.com/tbxark/gbvm/internal/env"
)

func TestNormalizeVersion(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{input: "(devel)", expected: "v0.0.0"},
		{input: "1.2.3", expected: "v1.2.3"},
		{input: "v1.2.3", expected: "v1.2.3"},
		{input: "v1.2", expected: "v1.2"},
		{input: "1.2.3-rc1", expected: "v1.2.3-rc1"},
		{input: "not-a-version", expected: "v0.0.0"},
	}

	for _, testCase := range cases {
		if got := normalizeVersion(testCase.input); got != testCase.expected {
			t.Fatalf("normalizeVersion(%q) = %q, want %q", testCase.input, got, testCase.expected)
		}
	}
}

func TestProxyCandidates(t *testing.T) {
	cases := []struct {
		proxy     string
		expected  []string
		expectErr bool
	}{
		{proxy: "https://proxy.golang.org", expected: []string{"https://proxy.golang.org"}},
		{proxy: "https://proxy.golang.org/,direct", expected: []string{"https://proxy.golang.org"}},
		{proxy: "https://proxy.golang.org|https://goproxy.io", expected: []string{"https://proxy.golang.org", "https://goproxy.io"}},
		{proxy: "direct", expected: []string{"https://proxy.golang.org"}},
		{proxy: "off", expectErr: true},
		{proxy: "  https://proxy.golang.org/  ", expected: []string{"https://proxy.golang.org"}},
	}

	for _, testCase := range cases {
		got, err := proxyCandidates(env.SplitGoProxy(testCase.proxy))
		if testCase.expectErr {
			if err == nil {
				t.Fatalf("proxyCandidates(%q) expected error", testCase.proxy)
			}
			continue
		}
		if err != nil {
			t.Fatalf("proxyCandidates(%q) unexpected error: %v", testCase.proxy, err)
		}
		if !reflect.DeepEqual(got, testCase.expected) {
			t.Fatalf("proxyCandidates(%q) = %v, want %v", testCase.proxy, got, testCase.expected)
		}
	}
}
