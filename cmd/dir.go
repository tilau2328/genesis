package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var dirCmd = &cobra.Command{
	Use:     "dir",
	Aliases: []string{"d"},
}

func init() {
	fileCmd.AddCommand(mkDirCmd)
	fileCmd.AddCommand(listDirsCmd)
}

var mkDirCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println(args)
		return os.Mkdir(args[0], 0644)
	},
}

var listDirsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			dir := path.Join(".", arg)
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, f := range files {
				if f.IsDir() {
					fmt.Println(f.Name())
				}
			}
		}
		return nil
	},
}
