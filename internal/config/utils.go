package config

import (
	"errors"
	"fmt"
	"go.yaml.in/yaml/v4"
	"os"
	"path/filepath"
	"strings"
)

const (
	RedColor   = "\033[31m"
	ResetColor = "\033[0m"
	Elipse     = ".................................................."
	ConfigDir  = ".shm"
	ManagerDir = ".shm/managers"
	ConfigFile = "base"
	Extention  = ".yml"
)

type Configuration interface {
	*Config | *PackageManager
}

func GetContextLines(file []byte, targetLine int) (string, error) {
	lines := strings.Split(string(file), "\n")
	totalLines := len(lines)

	if targetLine < 1 || targetLine > totalLines {
		return "", fmt.Errorf("line number %d is out of range (file has %d lines)", targetLine, totalLines)
	}

	start := targetLine - 6
	if start < 0 {
		start = 0
	}

	end := targetLine + 5
	if end > totalLines {
		end = totalLines
	}

	var cl strings.Builder

	cl.WriteString(Elipse + "\n")
	for i := start; i < end; i++ {
		lineNumber := i + 1
		line := lines[i]

		if lineNumber == targetLine {
			// Print the target line in red
			fmt.Fprintf(&cl, "%s%4d | %s%s\n", RedColor, lineNumber, line, ResetColor)
		} else {
			fmt.Fprintf(&cl, "%4d | %s\n", lineNumber, line)
		}
	}
	cl.WriteString(Elipse + "\n")

	return cl.String(), nil
}

func GetFullPath(path *string) error {
	if path == nil {
		return nil
	}

	p := *path

	if len(p) > 0 && p[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		if p == "~" {
			p = home
		} else if len(p) > 1 && p[1] == '/' {
			p = filepath.Join(home, p[2:])
		}
	}

	if !filepath.IsAbs(p) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		p = filepath.Join(cwd, p)
	}
	p = filepath.Clean(p)

	*path = p
	return nil
}

func GetConfigFile(path string) ([]byte, string, error) {
	if err := GetFullPath(&path); err != nil {
		return nil, "", err
	}

	if _, err := os.Stat(path); err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, "", err
	}

	return data, path, nil
}

func GetFromFile[T Configuration](path string, config T) {
	data, path, err := GetConfigFile(path)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		var te *yaml.TypeError
		if errors.As(err, &te) {
			for _, ue := range te.Errors {
				str, clerr := GetContextLines(data, ue.Line)
				if clerr == nil {
					fmt.Printf("Error in: %s\n%s\n%s\n", path, str, ue.Err.Error())
				}
			}
		} else {
			fmt.Printf("%s", err)
		}
		os.Exit(1)
	}
}
