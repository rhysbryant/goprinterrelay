package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mgutz/str"
)

func ToolsGetRequest(w http.ResponseWriter, r *http.Request) {
	WriteJsonResponse(w, tools)
}

func executeCommand(w io.Writer, cmd string, args []string) {
	expandedCmd := os.Expand(cmd, func(name string) string {
		argNumber, err := strconv.Atoi(name)
		if err != nil || argNumber < 1 || argNumber > len(args) {
			return ""
		}
		return args[argNumber-1]
	})

	parts := str.ToArgv(expandedCmd)
	head := parts[0]
	//parts = parts[1:len(parts)]

	c := exec.Cmd{}
	c.Path = head
	c.Args = parts
	c.Stdout = w
	c.Stderr = w
	err := c.Run()

	if err != nil {
		fmt.Fprintf(w, "error %s", err)
	}
}

func ToolsFormPostRequest(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		fmt.Fprintln(w, "expected multipart form")
		return
	}
	strToolId := r.PostFormValue("toolId")
	toolId, err := strconv.Atoi(strToolId)
	if err != nil {
		fmt.Fprintln(w, "Tool id not valid")
		return
	}
	if toolId < 0 && toolId >= len(tools) {
		fmt.Fprintln(w, "Tool id not valid")
		return
	}

	tool := tools[toolId]

	cmdArguments := make([]string, len(tool.Arguments))

	for index, arg := range tool.Arguments {
		fieldName := "toolArg" + strconv.Itoa(index)

		if arg.Type == "file" {
			uploadedFile, _, err := r.FormFile(fieldName)
			if err == http.ErrMissingFile && arg.Required {
				fmt.Fprintln(w, arg.Name+" is required")
				return
			} else if err != nil {
				fmt.Fprintf(w, "Error %s processing file upload\n", err)
				return
			}
			//todo use temp directory and generated name
			dstFileName := "tmp-" + fieldName + "upload"
			file, err := os.Create(dstFileName)
			if err != nil {
				fmt.Fprintf(w, "Error %s processing file upload\n", err)
				return
			}
			io.Copy(file, uploadedFile)
			file.Close()
			cmdArguments[index] = dstFileName

		} else if arg.Type == "text" {
			value := r.PostFormValue(fieldName)
			if value == "" && arg.Required {
				fmt.Fprintln(w, arg.Name+" is required")
				return
			}
			cmdArguments[index] = strings.Replace(value, "\"", "\\\"", -1)
		}
	}
	if tool.UsesPrinterSerialPort {
		serialConLock.Lock()
	}

	executeCommand(w, tool.CommandName, cmdArguments)

	if tool.UsesPrinterSerialPort {
		serialConLock.Unlock()
	}
}
