package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	RedColor   = "\033[31m"
	ResetColor = "\033[0m"
	Elipse     = ".................................................."
	ConfigDir  = ".yhm"
	ManagerDir = ".yhm/managers"
	ConfigFile = "yhm"
	Extention  = ".yml"
)

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

func GetConfigFile(path string, profile string) ([]byte, string, error) {
	var configFile string

	if err := GetFullPath(&path); err != nil {
		return nil, "", err
	}

	if info, err := os.Stat(path); !info.IsDir() || err != nil {
		return nil, "", err
	}

	if profile != "" {
		configFile = path + "/" + ConfigDir + "/" + ConfigFile + "-" + profile + Extention
	} else {
		configFile = path + "/" + ConfigDir + "/" + ConfigFile + Extention
	}

	data, err := os.ReadFile(configFile)

	if err != nil {
		return nil, "", err
	}

	return data, configFile, nil
}

func GetManagerFile(path string, file string) ([]byte, string, error) {
	var managerFile string

	if err := GetFullPath(&path); err != nil {
		return nil, "", err
	}

	if info, err := os.Stat(path); !info.IsDir() || err != nil {
		return nil, "", err
	}

	managerFile = path + "/" + ManagerDir + "/" + file + Extention

	data, err := os.ReadFile(managerFile)

	if err != nil {
		return nil, "", err
	}

	return data, managerFile, nil
}
