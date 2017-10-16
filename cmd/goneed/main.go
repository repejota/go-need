package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	goneed "github.com/repejota/go-need"
)

func isGoFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func getFiles(spath string) []string {
	var files []string
	filepath.Walk(spath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && isGoFile(info) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func main() {
	project := goneed.NewProject("./")
	project.Files = getFiles(project.SourcePath)

	for _, filePath := range project.Files {

		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			log.Fatal(err)
		}

		fs := bufio.NewScanner(f)
		lineno := 0
		for fs.Scan() {
			lineno++
			if strings.Contains(fs.Text(), `// TODO`) {
				todo := &goneed.Line{
					FilePath:    filePath,
					FileLastMod: fi.ModTime(),
					Number:      lineno,
					Src:         strings.TrimSpace(fs.Text()),
				}
				project.ToDos = append(project.ToDos, *todo)
			}
			if strings.Contains(fs.Text(), `// FIXME`) {
				fixme := &goneed.Line{
					FilePath:    filePath,
					FileLastMod: fi.ModTime(),
					Number:      lineno,
					Src:         strings.TrimSpace(fs.Text()),
				}
				project.FixMes = append(project.FixMes, *fixme)
			}
		}
	}

	now := time.Now()
	for _, todo := range project.ToDos {
		if todo.IsOutdated(now) {
			project.ExitCode = 1
			age := todo.GetAge(now)
			fmt.Printf("%s:%d - %v\n%s\n\n", todo.FilePath, todo.Number, age, todo.Src)
		}
	}

	os.Exit(project.ExitCode)
}
