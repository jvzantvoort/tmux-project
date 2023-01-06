package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/jvzantvoort/tmux-project/sessions"
	log "github.com/sirupsen/logrus"
)

type Project struct {
	Name    string
	AbsPath string
	Branch  string
	Expected bool
	Status  map[string]int
}

func (p Project) PrintLine() {
	br_col := color.New(BranchChangedColor)
	if stringInSlice(p.Branch, []string{"master", "main", "develop"}) {
		br_col = color.New(BranchDefaultColor)
	}
	var stat_str string
	for status, amount := range p.Status {
		stat_str = fmt.Sprintf("%s %s:%d", stat_str, status, amount)
		stat_str = strings.TrimSpace(stat_str)
	}
	if len(stat_str) != 0 {
		stat_str = fmt.Sprintf(" [%s]", stat_str)
	}

	br_str := br_col.Sprint(p.Branch)

	fmt.Printf("   %-32s %s%s\n", p.Name, br_str, stat_str)
	// fmt.Printf("%q\n", p)
}


func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

}

func printTitle(title string) {

	purple := color.New(TitleColor)
	fmt.Printf("\n  %s:\n\n", purple.Sprint(title))

}

func printInfo(itype, ival string) {
	infNameCol := color.New(InfoNameColor)
	infValCol := color.New(InfoValueColor)
	fmt.Printf("%-24s %s\n", infNameCol.Sprint(itype) + ":", infValCol.Sprint(ival))

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {

	flag.Parse()

	sessionname := os.Getenv("SESSIONNAME")
	session := sessions.NewTmuxSession(sessionname)
	printInfo("Sessionname", sessionname)
	printInfo("Projectdir", session.Workdir)
	printInfo("Description", session.Description)

	targets, err := ioutil.ReadDir(session.Workdir)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%#v\n"

	projects := []Project{}
	reserved := []string{}
	rest := []Project{}

	for _, target := range targets {
		target_name := target.Name()
		target_dir := path.Join(session.Workdir, target_name)

		target_stat, _ := os.Stat(target_dir)
		if !target_stat.IsDir() {
			continue
		}

		git := NewGitCmd(target_dir)
		if git.IsGit() {
			log.Debugf("%s is a gitdir", target_dir)
		} else {
			continue
		}

		proj := Project{}
		proj.Name = target_name
		proj.Branch, err = git.Branch()
		proj.Status = git.GetStatus()
		proj.AbsPath = target_dir
		if err != nil {
			log.Errorf("%v", err)
			continue
		}

		// we only want the session names
		if strings.HasSuffix(target_name, "_domain") {
			proj.Expected = true
			projects = append(projects, proj)
		} else if stringInSlice(target_name, reserved) {
			proj.Expected = true
			projects = append(projects, proj)
		} else {
			rest = append(rest, proj)
		}
	}

	sort.Slice(projects, func(i, j int) bool { return projects[i].Name < projects[j].Name })

	printTitle("Projects")
	for _, proj := range projects {
		if strings.HasSuffix(proj.Name, "_domain") {
			proj.PrintLine()
		}
	}

	printTitle("Extra")
	for _, proj := range projects {
		if !strings.HasSuffix(proj.Name, "_domain") {
			proj.PrintLine()
		}
	}

	printTitle("Rest")
	for _, proj := range rest {
		proj.PrintLine()
	}

}
