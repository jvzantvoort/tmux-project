package project

import (
	"os"
	"testing"

	"github.com/jvzantvoort/tmux-project/internal/testenv"
)

func TestMain(m *testing.M) { os.Exit(testenv.IsolateHome(m)) }
