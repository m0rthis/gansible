package gansible

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
)

type OS string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
)

var ValidOS = []OS{Windows, Linux}

type Result struct {
	Changed       bool           `json:"changed"`
	Failed        bool           `json:"failed"`
	Msg           string         `json:"msg,omitempty"`
	Data          any            `json:"data,omitempty"`
	ModuleName    string         `json:"module_name,omitempty"`
	ModuleVersion string         `json:"module_version,omitempty"`
	ModuleStdout  string         `json:"module_stdout,omitempty"`
	ModuleStderr  string         `json:"module_stderr,omitempty"`
	Warnings      []string       `json:"warnings,omitempty"`
	RC            int            `json:"rc,omitempty"`
	AnsibleFacts  map[string]any `json:"ansible_facts,omitempty"`
}

func ExitJson(res Result) {
	res.Failed = false
	printAndExit(res, 0)
}

func FailJson(res Result) {
	res.Failed = true
	printAndExit(res, 1)
}

func printAndExit(res Result, code int) {
	output, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"failed":true,"msg":"failed to marshal JSON"}`)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, string(output))
	os.Exit(code)
}

func GetOS() (OS, error) {
	switch runtime.GOOS {
	case "windows":
		return Windows, nil
	case "linux":
		return Linux, nil
	default:
		return "", fmt.Errorf("currently unsupported operating system")
	}
}

func LoadArgs[T any](opts LoadOptions) (T, Result) {
	var result Result
	var args T

	if !opts.UseRawInput {
		if len(os.Args) < 2 {
			result.Msg = "Collection failure, missing argument file."
			result.Msg += "\n HINT: Are you running this outside of Ansible?"
			result.Msg += "\n If so, enable USE_RAW_INPUT"
			FailJson(result)
		}

		content, err := os.ReadFile(os.Args[1])
		if err != nil {
			result.Msg = fmt.Sprintf("Error reading file: %s", err)
			FailJson(result)
		}
		cleanInput, err := stripAnsibleFields(content)
		if err != nil {
			result.Msg = "Error cleaning input: " + err.Error()
			FailJson(result)
		}

		decoder := json.NewDecoder(bytes.NewReader(cleanInput))
		decoder.DisallowUnknownFields()

		err = decoder.Decode(&args)
		if err != nil && !errors.Is(err, io.EOF) {
			result.Msg = fmt.Sprintf("issue decoding JSON with error: %s", err)
			FailJson(result)
		}

		return args, result
	}

	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeCharDevice != 0 {
		result.Msg = "no input received on stdin"
		FailJson(result)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		result.Msg = "failed to read stdin"
		FailJson(result)
	}
	cleanInput, err := stripAnsibleFields(input)
	if err != nil {
		result.Msg = "Error cleaning input: " + err.Error()
		FailJson(result)
	}

	decoder := json.NewDecoder(bytes.NewReader(cleanInput))
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&args)
	if err != nil && !errors.Is(err, io.EOF) {
		result.Msg = fmt.Sprintf("issue decoding JSON with error: %s", err)
		FailJson(result)
	}
	return args, result
}
