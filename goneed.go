package goneed

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Line ...
type Line struct {
	FilePath    string
	FileLastMod time.Time
	Number      int
	Src         string
}

// IsOutdated ...
func (l *Line) IsOutdated(now time.Time) bool {
	outdated := false
	diff := now.Sub(l.FileLastMod)
	if diff.Hours() > 1 {
		outdated = true
	}
	return outdated
}

// GetAge ...
func (l *Line) GetAge(now time.Time) time.Duration {
	linefile := fmt.Sprintf("%d,%d:%s", l.Number, l.Number, l.FilePath)
	out, err := exec.Command("git", "log", "--pretty=format:%at", "-L", linefile).Output()
	if err != nil {
		log.Fatal(err)
	}
	slastmod := strings.Split(string(out), "\n")[0]
	i, err := strconv.ParseInt(slastmod, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	lastmod := time.Unix(i, 0)
	diff := now.Sub(lastmod)
	return diff
}

// Project ...
type Project struct {
	SourcePath string
	ExitCode   int
	Files      []string
	ToDos      []Line
	FixMes     []Line
}

// NewProject ...
func NewProject(spath string) *Project {
	project := Project{
		SourcePath: spath,
		ExitCode:   0,
	}
	return &project
}
