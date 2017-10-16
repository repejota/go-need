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
func (l *Line) IsOutdated() bool {
	now := time.Now()
	outdated := false
	d := time.Duration(time.Hour * 24 * 30 * 6) // 6 months
	monthago := now.Add(-d)
	diff := monthago.Sub(l.FileLastMod)
	if diff.Hours() > 1 {
		outdated = true
	}
	return outdated
}

// GetFileAge ...
func (l *Line) GetFileAge() time.Time {
	out, err := exec.Command("git", "log", "-n", "1", "--pretty=format:%at", "--", l.FilePath).Output()
	if err != nil {
		log.Fatal(err)
	}
	slastmod := strings.Split(string(out), "\n")[0]
	i, err := strconv.ParseInt(slastmod, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	lastmod := time.Unix(i, 0)
	return lastmod
}

// GetLineAge ...
func (l *Line) GetLineAge() time.Time {
	linefile := fmt.Sprintf("%d,%d:%s", l.Number, l.Number, l.FilePath)
	out, err := exec.Command("git", "log", "-n", "1", "--pretty=format:%at", "-L", linefile).Output()
	if err != nil {
		log.Fatal(err)
	}
	slastmod := strings.Split(string(out), "\n")[0]
	i, err := strconv.ParseInt(slastmod, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	lastmod := time.Unix(i, 0)
	return lastmod
}

// Project ...
type Project struct {
	SourcePath string
	ExitCode   int
	Files      []string
	ToDos      []Line
}

// NewProject ...
func NewProject(spath string) *Project {
	project := Project{
		SourcePath: spath,
		ExitCode:   0,
	}
	return &project
}
