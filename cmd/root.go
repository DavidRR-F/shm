package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DavidRR-F/shm/internal/config"
	"github.com/DavidRR-F/shm/internal/utils"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

var (
	install bool
	test    bool
	profile string
)

var rootCmd = &cobra.Command{
	Use:   "shm [directory]",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	Run:   runYhm,
}

func runYhm(cmd *cobra.Command, args []string) {
	c := config.Config{}
	m := []config.PackageManager{}

	data, path, err := utils.GetConfigFile(args[0], profile)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		var te *yaml.TypeError
		if errors.As(err, &te) {
			for _, ue := range te.Errors {
				str, clerr := utils.GetContextLines(data, ue.Line)
				if clerr == nil {
					fmt.Printf("Error in: %s\n%s\n%s\n", path, str, ue.Err.Error())
				}
			}
		} else {
			fmt.Printf("%s", err)
		}
		os.Exit(1)
	}

	if install {
		for _, file := range c.Managers {
			data, path, err := utils.GetManagerFile(args[0], file)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			var tmp config.PackageManager
			if err := yaml.Unmarshal(data, &tmp); err != nil {
				var te *yaml.TypeError
				if errors.As(err, &te) {
					for _, ue := range te.Errors {
						str, clerr := utils.GetContextLines(data, ue.Line)
						if clerr == nil {
							fmt.Printf("Error in: %s\n%s\n%s\n", path, str, ue.Err.Error())
						}
					}
				} else {
					fmt.Printf("%s", err)
				}
				os.Exit(1)
			}
			m = append(m, tmp)
		}
	}

	for _, link := range c.Links {
		if !test {
			if err := link.CreateLink(); err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Link: ln -s %s %s\n", link.Src, link.Dest)
		}
	}

	if install {
		for _, manager := range m {
			if !test {
				if err := manager.InstallPackages(); err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				}
			} else {
				if _, err := exec.LookPath(manager.Name); err != nil {
					fmt.Printf("\n%s: Manager not installed pre-install will be invoked (if exists)\n", manager.Name)
				} else {
					fmt.Printf("\n%s: Manager installed pre-install will not be invoked (if exists)\n", manager.Name)
				}
				if len(manager.PreInstall) != 0 {
					fmt.Printf("\tPre Install: %s\n", strings.Join(manager.PreInstall, " "))
				}
				if len(manager.PostInstall) != 0 {
					fmt.Printf("\tPost Install: %s\n", strings.Join(manager.PostInstall, " "))
				}
				fmt.Printf("\tInstall: %s %s %s\n", manager.Name, strings.Join(manager.Args, " "), strings.Join(manager.Packages, " "))
			}
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&install, "install-packages", "i", false, "Enable package installations")
	rootCmd.Flags().BoolVar(&test, "dry-run", false, "Run command without invoking commands")
	rootCmd.Flags().StringVarP(&profile, "profile", "p", "", "Choose configuration profile")
}
