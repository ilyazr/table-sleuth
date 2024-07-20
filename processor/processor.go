package processor

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type TableToSvcUsage struct {
	TableName string   `json:"tableName"`
	Services  []string `json:"services"`
}

type SvcToTableUsage struct {
	ServiceName string   `json:"serviceName"`
	Tables      []string `json:"tables"`
}

type ProcessedProject struct {
	ServiceName string
	Tables      []string
}

func ProcessDirOfProjects(dir string) []*TableToSvcUsage {
	childrenDirs, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	preparedPaths := make([]string, 0)
	for _, childrenDir := range childrenDirs {
		preparedPaths = append(preparedPaths, filepath.Join(dir, childrenDir.Name()))
	}
	processed := ProcessProjects(preparedPaths)
	m := make(map[string][]string)
	for _, v := range processed {
		for _, t := range v.Tables {
			if _, ok := m[t]; !ok {
				m[t] = make([]string, 0)
			}
			m[t] = append(m[t], v.ServiceName)
		}
	}
	res := make([]*TableToSvcUsage, 0)
	for k, v := range m {
		res = append(res, &TableToSvcUsage{
			TableName: k,
			Services:  v,
		})
	}
	return res
}

func ProcessProjects(paths []string) []*SvcToTableUsage {
	var result []*SvcToTableUsage
	var wg sync.WaitGroup
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			parsed := ProcessSingleProject(p)
			q := &SvcToTableUsage{
				ServiceName: parsed.ServiceName,
				Tables:      parsed.Tables,
			}
			result = append(result, q)
		}(path)
	}
	wg.Wait()
	return result
}

func ProcessSingleProject(path string) *ProcessedProject {
	files := make([]string, 0)
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".java") {
			return nil
		}
		files = append(files, path)
		return nil
	})

	pattern := regexp.MustCompile(`name\s*=\s*"([a-zA-Z0-9_-]+)"`)

	tableNames := make([]string, 0)

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			openedFile, err := os.Open(f)
			if err != nil {
				log.Println("can't open file: " + err.Error())
				return
			}
			defer openedFile.Close()

			scanner := bufio.NewScanner(openedFile)
			for scanner.Scan() {
				text := scanner.Text()
				// not an annotation
				if !strings.HasPrefix(text, "@Table") {
					continue
				}
				// we reach a class definition
				if strings.HasPrefix(text, "public") {
					break
				}

				//@Table(name = "table name")
				submatch := pattern.FindStringSubmatch(text)

				if submatch == nil || len(submatch) < 2 {
					continue
				}

				tableName := submatch[1]
				tableNames = append(tableNames, tableName)

			}

			if err := scanner.Err(); err != nil {
				log.Println("error scanning file: " + err.Error())
			}

		}(file)
	}
	wg.Wait()
	return &ProcessedProject{
		ServiceName: filepath.Base(path),
		Tables:      tableNames,
	}
}
