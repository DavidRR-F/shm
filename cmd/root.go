package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DavidRR-F/shm/internal/config"
	"github.com/spf13/cobra"
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
	p := config.Config{}
	m := []config.PackageManager{}

	config.GetFromFile(args[0]+"/"+config.ConfigDir+"/"+config.ConfigFile+config.Extention, &c)

	if profile != "" {
		config.GetFromFile(args[0]+"/"+config.ConfigDir+"/"+profile+config.Extention, &p)
	}

	if install {
		for _, file := range append(c.Managers, p.Managers...) {
			var tmp config.PackageManager
			config.GetFromFile(args[0]+"/"+config.ManagerDir+"/"+file+config.Extention, &tmp)
			m = append(m, tmp)
		}
	}

	for _, link := range append(c.Links, p.Links...) {
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
				if len(manager.Install.Pre) != 0 {
					fmt.Printf("\tPre Install: %s\n", strings.Join(manager.Install.Pre, " "))
				}
				if len(manager.Install.Post) != 0 {
					fmt.Printf("\tPost Install: %s\n", strings.Join(manager.Install.Post, " "))
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
