package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
)

type YamlValidationError struct {
	Line  int
	Value interface{}
	Err   error
}

func (e *YamlValidationError) Error() string {
	return fmt.Sprintf("yaml validation error on line %d: %v (value: %#v)", e.Line, e.Err, e.Value)
}

func (e *YamlValidationError) Unwrap() error {
	return e.Err
}

type DstPath string

func (p *DstPath) UnmarshalYAML(value *yaml.Node) error {
	var raw string
	if err := value.Decode(&raw); err != nil {
		return err
	}

	if err := GetFullPath(&raw); err != nil {
		return err
	}

	if raw == "" {
		return errors.New("path is empty")
	}

	if strings.ContainsRune(raw, '\x00') {
		return errors.New("path contains null byte")
	}

	clean := filepath.Clean(raw)

	if strings.Contains(clean, "\x00") {
		return &YamlValidationError{
			Line:  value.Line,
			Value: value.Value,
			Err:   errors.New("path contains null byte"),
		}
	}

	*p = DstPath(raw)
	return nil
}

type SrcPath string

func (p *SrcPath) UnmarshalYAML(value *yaml.Node) error {
	var raw string
	if err := value.Decode(&raw); err != nil {
		return err
	}

	fmt.Printf("%s\n", raw)
	if err := GetFullPath(&raw); err != nil {
		return err
	}

	if raw == "" {
		return errors.New("path is empty")
	}

	if strings.ContainsRune(raw, '\x00') {
		return errors.New("path contains null byte")
	}

	clean := filepath.Clean(raw)

	if strings.Contains(clean, "\x00") {
		return &YamlValidationError{
			Line:  value.Line,
			Value: value.Value,
			Err:   errors.New("path contains null byte"),
		}
	}

	if _, err := os.Stat(clean); err != nil {
		return &YamlValidationError{
			Line:  value.Line,
			Value: value.Value,
			Err:   err,
		}
	}

	*p = SrcPath(raw)
	return nil
}
