package config

import (
	"os/exec"
)

type SymbolicLink struct {
	Src  SrcPath `yaml:"src"`
	Dest DstPath `yaml:"dest"`
	Exe  bool    `yaml:"exe"`
}

func (f *SymbolicLink) CreateLink() error {
	cmd := exec.Command("ln", "-s", string(f.Src), string(f.Dest))

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := f.makeExecutable(); err != nil {
		return err
	}

	return nil
}

func (f *SymbolicLink) makeExecutable() error {
	if !f.Exe {
		return nil
	}

	cmd := exec.Command("chmod", "+x", string(f.Src))

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
