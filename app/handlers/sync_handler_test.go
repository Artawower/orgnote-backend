package handlers

import "testing"

func TestValidateFilePath_Valid(t *testing.T) {
	tests := []string{
		"notes/my-note.org",
		"folder/subfolder/note.org",
		"readme.md",
		"notes/my.note.org",
		"заметки/файл.org",
	}

	for _, path := range tests {
		err := validateFilePath(path)
		if err != nil {
			t.Errorf("validateFilePath(%q) unexpected error: %v", path, err)
		}
	}
}

func TestValidateFilePath_Invalid(t *testing.T) {
	tests := []struct {
		path string
		desc string
	}{
		{"", "empty path"},
		{"/notes/file.org", "absolute path"},
		{"notes/../secret/file.org", "parent reference"},
		{"notes/..hidden/file.org", "double dots"},
		{"notes/CON.txt", "windows reserved CON"},
		{"PRN/file.org", "windows reserved PRN"},
		{"folder/NUL", "windows reserved NUL"},
		{"COM1.org", "windows reserved COM1"},
		{"notes/LPT1.txt", "windows reserved LPT1"},
		{"notes//file.org", "empty segment"},
	}

	for _, tt := range tests {
		err := validateFilePath(tt.path)
		if err == nil {
			t.Errorf("validateFilePath(%q) expected error for %s", tt.path, tt.desc)
		}
	}
}
