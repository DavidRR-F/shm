package config

import (
	"os"
	"os/exec"
)

type PackageManager struct {
	Name        string   `yaml:"name"`
	Command     []string `yaml:"command"`
	Args        []string `yaml:"args"`
	PreInstall  []string `yaml:"pre-install"`
	PostInstall []string `yaml:"post-install"`
	Packages    []string `yaml:"packages"`
}

func (pm *PackageManager) InstallPackages() error {
	if _, err := exec.LookPath(pm.Name); err != nil {
		if len(pm.PreInstall) != 0 {
			if err := pm.runPreInstallCommand(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := pm.runPostInstallCommand(); err != nil {
		return err
	}

	if err := pm.installPackages(); err != nil {
		return err
	}

	return nil
}

func (pm *PackageManager) installPackages() error {
	var fullArgs []string
	if len(pm.Command) > 1 {
		fullArgs = append(fullArgs, pm.Command[1:]...)
	}

	fullArgs = append(fullArgs, pm.Args...)

	fullArgs = append(fullArgs, pm.Packages...)

	cmd := exec.Command(pm.Command[0], fullArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (pm *PackageManager) runPreInstallCommand() error {

	cmd := exec.Command(pm.PreInstall[0], pm.PreInstall[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (pm *PackageManager) runPostInstallCommand() error {

	if len(pm.PostInstall) == 0 {
		return nil
	}

	cmd := exec.Command(pm.PostInstall[0], pm.PostInstall[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
