package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"lazygo/utils/model"
	"os"
	"strings"
)

// GetPackageName ...
func GetPackageName(rootFolderName string) string {
	i := strings.LastIndex(rootFolderName, "/")
	return rootFolderName[i+1:]
}

func GetGogenConfig() model.GogenConfig {

	fileInBytes, err := os.ReadFile("./.lazygo/lazygo_config.json")
	if err != nil {
		fmt.Printf(".lazygo/lazygo_config.json is not found. Please run 'lazygo domain' first\n")
		os.Exit(1)
	}

	var cfg model.GogenConfig

	err = json.Unmarshal(fileInBytes, &cfg)
	if err != nil {
		fmt.Printf("fail to unmarshal .lazygo/lazygo_config.json\n")
		os.Exit(1)
	}

	return cfg

}

func GetPackagePath() string {

	var gomodPath string

	file, err := os.Open("go.mod")
	if err != nil {
		fmt.Printf("go.mod is not found. Please create it with command `go mod init your/path/project`\n")
		os.Exit(1)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		if strings.HasPrefix(row, "module") {
			moduleRow := strings.Split(row, " ")
			if len(moduleRow) > 1 {
				gomodPath = moduleRow[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	return strings.Trim(gomodPath, "\"")

}

func GetExecutableName() string {
	pn := GetPackagePath()
	i := strings.LastIndex(pn, "/")
	return pn[i+1:]
}
