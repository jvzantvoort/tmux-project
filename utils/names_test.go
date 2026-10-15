package utils

import "testing"

func TestCheckPathName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "plain", input: "JIRA-1234", wantErr: false},
		{name: "with dot", input: "release.1", wantErr: false},
		{name: "empty", input: "", wantErr: true},
		{name: "dot", input: ".", wantErr: true},
		{name: "dot-dot", input: "..", wantErr: true},
		{name: "traversal", input: "../etc", wantErr: true},
		{name: "slash", input: "a/b", wantErr: true},
		{name: "backslash", input: `a\b`, wantErr: true},
		{name: "newline", input: "a\nb", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPathName("project", tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPathName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestCheckSessionName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "plain", input: "JIRA-1234", wantErr: false},
		{name: "underscore", input: "my_project", wantErr: false},
		{name: "unicode", input: "projekt-ü", wantErr: false},
		{name: "dot", input: "release.1", wantErr: true},
		{name: "colon", input: "a:b", wantErr: true},
		{name: "space", input: "a b", wantErr: true},
		{name: "tab", input: "a\tb", wantErr: true},
		{name: "leading dash", input: "-a", wantErr: true},
		{name: "slash", input: "a/b", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSessionName("project", tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSessionName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
