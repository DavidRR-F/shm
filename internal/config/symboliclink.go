package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type SymbolicLink struct {
	Src  SrcPath `yaml:"src"`
	Dest DstPath `yaml:"dst"`
	Exe  bool    `yaml:"exe"`
}

func (f *SymbolicLink) CreateLink() error {
	dest := string(f.Dest)
	src := string(f.Src)

	if info, err := os.Lstat(dest); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(dest); err == nil {
				resolvedTarget := filepath.Clean(target)
				resolvedSource := filepath.Clean(src)
				if resolvedTarget == resolvedSource {
					return nil
				}
			}
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("failed to remove existing symlink %s: %w", dest, err)
			}
		} else {
			return fmt.Errorf("destination %s exists and is not a symlink", dest)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat %s: %w", dest, err)
	}

	if err := os.Symlink(src, dest); err != nil {
		return fmt.Errorf("failed to create symlink from %s to %s: %w", dest, src, err)
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
