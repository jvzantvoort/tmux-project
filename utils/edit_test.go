package utils

import (
	"reflect"
	"testing"
)

func TestEditArgs(t *testing.T) {
	tests := []struct {
		name       string
		visual     string
		editor     string
		wantEditor string
		wantArgs   []string
	}{
		{name: "default vim splits", wantEditor: "vim", wantArgs: []string{"-O", "a", "b"}},
		{name: "EDITOR nano", editor: "nano", wantEditor: "nano", wantArgs: []string{"a", "b"}},
		{name: "VISUAL wins", visual: "nvim", editor: "nano", wantEditor: "nvim", wantArgs: []string{"-O", "a", "b"}},
		{name: "editor with arguments", editor: "code --wait", wantEditor: "code", wantArgs: []string{"--wait", "a", "b"}},
		{name: "full path to vim", editor: "/usr/bin/vim", wantEditor: "/usr/bin/vim", wantArgs: []string{"-O", "a", "b"}},
		{name: "plain vi has no -O", editor: "vi", wantEditor: "vi", wantArgs: []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("VISUAL", tt.visual)
			t.Setenv("EDITOR", tt.editor)
			editor, args := EditArgs("a", "b")
			if editor != tt.wantEditor || !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("EditArgs() = %q %q, want %q %q", editor, args, tt.wantEditor, tt.wantArgs)
			}
		})
	}
}
