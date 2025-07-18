package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jvzantvoort/tmux-project/git"
	"github.com/jvzantvoort/tmux-project/utils"
)

// ProjectDef represents a project definition with various attributes such as path, name, branch, and status.
// It includes methods to initialize the project, retrieve git information, and format the project fields for display.
type ProjectDef struct {
	Path       string         // Path to the project directory
	ProjectDir string         // Project directory
	Name       string         // Project name
	AbsPath    string         // Absolute path to the project
	Branch     string         // Current git branch
	Expected   bool           // Whether the project is expected to be found
	Status     map[string]int // Git status information
	SubPath    string         // Sub-path within the project directory
	Chapter    string         // Chapter or section of the project
	Info       os.FileInfo    // File information for the project directory
}

// GetFields returns a slice of strings representing the fields of the project.
// It includes the project name, branch, and status information formatted as a string.
// The branch is colored based on its status, with default colors for common branches like "master", "main", and "develop".
// If the project is not found, an empty string is returned.
// The status information is formatted as a string showing the counts of different git statuses (e.g., modified, added).
// This method is useful for displaying project information in a structured format, such as in a table
func (p ProjectDef) GetFields() []string {
	retv := []string{}
	retv = append(retv, p.Name)

	br_col := color.New(BranchChangedColor)
	if utils.StringInSlice(p.Branch, []string{"master", "main", "develop"}) {
		br_col = color.New(BranchDefaultColor)
	}
	var stat_str string
	for status, amount := range p.Status {
		stat_str = fmt.Sprintf("%s %s:%d", stat_str, status, amount)
		stat_str = strings.TrimSpace(stat_str)
	}
	if len(stat_str) != 0 {
		retv = append(retv, fmt.Sprintf("[%s]", stat_str))
	} else {
		retv = append(retv, " ")
	}

	retv = append(retv, br_col.Sprint(p.Branch))

	return retv
}

// GetGitInfo retrieves the git information for the project.
// It uses the git package to check if the project is a git repository and retrieves the current branch and status.
// The branch and status are stored in the ProjectDef struct.
// This method is essential for projects that use git for version control, allowing the command to display the current branch and status.
func (p *ProjectDef) GetGitInfo() {
	gitcmnd := git.NewGitCmd(p.Path)
	if gitcmnd.IsGit() {
		p.Branch, _ = gitcmnd.Branch()
		p.Status = gitcmnd.GetStatus()
	}
}

func NewProjectDef(projdir, dirname string) *ProjectDef {
	retv := &ProjectDef{}

	retv.Path = dirname
	retv.ProjectDir = projdir

	return retv
}

func (pd *ProjectDef) Init() {
	utils.LogStart()
	defer utils.LogEnd()

	pd.GetGitInfo()

	pd.Info, _ = os.Lstat(pd.Path)
	if pd.ProjectDir == pd.Path {
		pd.Name = "."
		pd.Chapter = "root"
	} else {
		pd.Name = pd.Path[len(pd.ProjectDir)+1:]
		pd.Chapter = "rest"
	}

}
