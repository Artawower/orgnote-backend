package services

import "testing"

func TestActivationDomainFromClientAddress(t *testing.T) {
	tests := []struct {
		name          string
		clientAddress string
		expected      *string
	}{
		{
			name:          "extracts host from https url",
			clientAddress: "https://dev.org-note.com",
			expected:      stringPointer("dev.org-note.com"),
		},
		{
			name:          "keeps localhost port",
			clientAddress: "http://localhost:9000",
			expected:      stringPointer("localhost:9000"),
		},
		{
			name:          "supports raw host",
			clientAddress: "org-note.com",
			expected:      stringPointer("org-note.com"),
		},
		{
			name:          "trims whitespace",
			clientAddress: "  https://org-note.com/app  ",
			expected:      stringPointer("org-note.com"),
		},
		{
			name:          "ignores empty address",
			clientAddress: "",
			expected:      nil,
		},
		{
			name:          "ignores whitespace-only address",
			clientAddress: "   ",
			expected:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := activationDomainFromClientAddress(test.clientAddress)
			assertStringPointersEqual(t, actual, test.expected)
		})
	}
}

func stringPointer(value string) *string {
	return &value
}

func assertStringPointersEqual(t *testing.T, actual *string, expected *string) {
	t.Helper()

	if actual == nil || expected == nil {
		if actual != expected {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		return
	}

	if *actual != *expected {
		t.Fatalf("expected %q, got %q", *expected, *actual)
	}
}
