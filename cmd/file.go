package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
	"strings"
)

var fileCmd = &cobra.Command{
	Use:     "file",
	Aliases: []string{"f"},
}

func init() {
	fileCmd.AddCommand(listFilesCmd)
	fileCmd.AddCommand(readFilesCmd)
	fileCmd.AddCommand(writeFileCmd)
}

var listFilesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			files, err := ioutil.ReadDir(path.Join(".", arg))
			if err != nil {
				return err
			}
			for _, f := range files {
				if !f.IsDir() {
					fmt.Println(f.Name())
				}
			}
		}
		return nil
	},
}

var readFilesCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing arguments")
		}
		for _, arg := range args {
			dat, err := ioutil.ReadFile(path.Join(".", arg))
			if err != nil {
				return err
			}
			fmt.Print(string(dat))
		}
		return nil
	},
}

var writeFileCmd = &cobra.Command{
	Use:     "write",
	Aliases: []string{"w"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing arguments")
		}
		return ioutil.WriteFile(path.Join(".", args[0]), []byte(strings.Replace(strings.Join(args[1:], " "), "\\n", "\n", -1)), 0644)
	},
}
