package project

import (
	"github.com/jvzantvoort/tmux-project/archive"
	"github.com/jvzantvoort/tmux-project/utils"
)

func (proj Project) ListFiles() []string {
	utils.LogStart()
	defer utils.LogEnd()

	retv := []string{}

	utils.Debugf("  add %#v", proj.Directory)
	retv = append(retv, proj.Directory)

	utils.Debugf("  add %#v", proj.ProjectConfigFile())
	retv = append(retv, proj.ProjectConfigFile())

	for _, proj_target := range proj.Targets {
		dest := proj.CalcDestination(proj_target.Destination)
		utils.Debugf("  add %s (%s)", dest, proj.Directory)
		retv = append(retv, dest)
	}
	return retv
}

func (proj Project) Archive(archivedescr string) error {
	utils.LogStart()
	defer utils.LogEnd()

	if archivedescr == ".tar.gz" {
		utils.Abort("no name for archive")
	}

	archivename, _ := utils.Expand(proj.Parse(archivedescr))

	utils.Debugf("archive name: %s", archivename)

	tar := archive.NewTarArchive(archivename)

	tar.AddFiles(proj.ListFiles())

	return tar.CreateArchive()

}
