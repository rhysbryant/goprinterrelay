package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

type ToolArgument struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type Tool struct {
	DisplayName           string         `json:"name"`
	CommandName           string         `json:"cmd"`
	Arguments             []ToolArgument `json:"formfields"`
	UsesPrinterSerialPort bool           `json:"usesPrinterConnection"`
}

func loadTool(fileName string) (*Tool, error) {
	config := Tool{}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	e := json.Unmarshal(file, &config)
	if e != nil {
		return nil, e
	}

	return &config, nil
}

func getToolFileList(dir string) ([]string, error) {

	fileNames := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileNames = append(fileNames, path.Join(dir, file.Name()))
	}

	return fileNames, nil

}

func loadTools(dir string) ([]Tool, error) {
	tools = []Tool{}

	toolFiles, err := getToolFileList(dir)
	if err != nil {
		fmt.Println("error getting tools list")
		return nil, err
	}
	for i := range toolFiles {
		tool, err := loadTool(toolFiles[i])
		if err != nil {
			fmt.Printf("Failed to load %s error %s\n", toolFiles[i], err.Error())
			continue
		}
		tools = append(tools, *tool)
	}

	return tools, nil
}
