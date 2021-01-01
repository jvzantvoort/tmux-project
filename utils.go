package tmuxproject

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	HomeDir           = ""
	TmuxDir           = ""
	ProjTypeConfigDir = ""
)

type ListTable struct {
	Name        string
	Description string
	Workdir     string
}

// targetExists return true if target exists
func targetExists(targetpath string) bool {
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Which return path
func Which(command string) string {
	Path := strings.Split(os.Getenv("PATH"), ":")
	var retv string
	for _, dirname := range Path {
		fullpath := path.Join(dirname, command)
		if targetExists(fullpath) {
			retv = fullpath
			break
		}
	}
	return retv
}

// ExitOnError check error and exit if not nil
func ExitOnError(err error) {
	if err != nil {
		log.Errorf("error %v\n", err)
		os.Exit(1)
	}
}

// GetHomeDir simple wrapper function to keep from calling the same functions
// over and over again.
func GetHomeDir() string {
	if len(HomeDir) > 0 {
		return HomeDir
	}
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	HomeDir = usr.HomeDir
	return HomeDir
}

func Ask(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", question)
	text, _ := reader.ReadString('\n')
	return text
}

func GetTmuxDir() string {
	if len(TmuxDir) == 0 {
		homedir := GetHomeDir()
		TmuxDir = path.Join(homedir, ".bash", "tmux.d")
	}
	return TmuxDir
}

func GetProjTypeConfigDir() string {
	if len(ProjTypeConfigDir) == 0 {
		homedir := GetHomeDir()
		ProjTypeConfigDir = path.Join(homedir, ".tmux-project")
	}
	return ProjTypeConfigDir
}

func ListTmuxConfigs() []ListTable {
	var retv []ListTable
	tmuxdir := GetTmuxDir()

	targets, err := ioutil.ReadDir(tmuxdir)
	if err != nil {
		log.Fatal(err)
	}

	for _, target := range targets {
		target_name := target.Name()

		// we only want the session names
		if !strings.HasSuffix(target_name, ".rc") {
			continue
		}

		// "common" is shared by all others
		if target_name == "common.rc" {
			continue
		}
		description := ""
		workdir := ""

		target_name = strings.TrimSuffix(target_name, ".rc")
		description, err = GetDescription(target_name)
		if err != nil {
			log.Fatal(err)
		}
		workdir, err = GetWorkdir(target_name)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("description: %s", description)
		log.Debugf("workdir: %s", workdir)
		t := ListTable{}
		t.Name = target_name
		t.Description = description

		workdir, err = ExpandHome(workdir)
		if err != nil {
			log.Fatal(err)
		}

		t.Workdir = workdir
		retv = append(retv, t)
	}
	return retv
}

// ExpandHome expand the tilde in a given path.
func ExpandHome(pathstr string) (string, error) {
	if len(pathstr) == 0 {
		return pathstr, nil
	}

	if pathstr[0] != '~' {
		return pathstr, nil
	}

	homedir := GetHomeDir()
	return filepath.Join(homedir, pathstr[1:]), nil

}

func LoadStringLines(target string) ([]string, error) {
	var retv []string
	content, err := ioutil.ReadFile(target)
	if err != nil {
		return retv, err
	}

	for _, line := range strings.Split(string(content), "\n") {
		retv = append(retv, line)
	}
	return retv, nil
}

// GetMatches parses each line and sticks the findings in a map
func GetMatches(regEx string, lines []string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)

	paramsMap = make(map[string]string)

	for _, line := range lines {
		match := compRegEx.FindStringSubmatch(line)
		for i, name := range compRegEx.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
	}
	return
}
