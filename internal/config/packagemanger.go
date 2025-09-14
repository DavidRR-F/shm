package config

import (
	"os"
	"os/exec"
)

type Install struct {
	Pre  []string `yaml:"pre"`
	Post []string `yaml:"post"`
}

type PackageManager struct {
	Name     string   `yaml:"name"`
	Command  []string `yaml:"command"`
	Args     []string `yaml:"args"`
	Install  Install  `yaml:"install"`
	Packages []string `yaml:"packages"`
}

func (pm *PackageManager) InstallPackages() error {
	if _, err := exec.LookPath(pm.Name); err != nil {
		if len(pm.Install.Pre) != 0 {
			if err := pm.runPreInstallCommand(); err != nil {
				return err
			}
		}
		if len(pm.Install.Post) != 0 {
			if err := pm.runPostInstallCommand(); err != nil {
				return err
			}
		}
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

	cmd := exec.Command(pm.Install.Pre[0], pm.Install.Pre[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (pm *PackageManager) runPostInstallCommand() error {

	cmd := exec.Command(pm.Install.Post[0], pm.Install.Post[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
