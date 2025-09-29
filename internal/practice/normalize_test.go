package practice

import "testing"

func TestNormalizeFolderName(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"test", "test"},
		{"test/test", "test_test"},
		{"test\\test", "test_test"},
		{"默认", DefaultFolderDir},
		{"", DefaultFolderDir},
		{"../unsafe", "unsafe"},
	}

	for _, tc := range cases {
		got, _ := NormalizeFolderName(tc.input)
		if got != tc.expected {
			t.Fatalf("normalize %q => %q, expected %q", tc.input, got, tc.expected)
		}
	}
}
