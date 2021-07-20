package cmd

import (
	"errors"
	"fmt"
	"github.com/badgerodon/penv"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var envCmd = &cobra.Command{
	Use:     "env",
	Aliases: []string{"e"},
}

func init() {
	envCmd.AddCommand(unsetEnvCmd)
	envCmd.AddCommand(listEnvCmd)
	envCmd.AddCommand(setEnvCmd)
	envCmd.AddCommand(getEnvCmd)
}

var getEnvCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			log.Println(os.Getenv(arg))
		}
	},
}

var setEnvCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			return errors.New("not enough arguments")
		case 1:
			return penv.SetEnv(args[0], "")
		default:
			return penv.SetEnv(args[0], strings.Join(args[1:], ";"))
		}
	},
}

var unsetEnvCmd = &cobra.Command{
	Use:     "unset",
	Aliases: []string{"u"},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			if err := penv.UnsetEnv(arg); err != nil {
				return err
			}
		}
		return nil
	},
}

var listEnvCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			found := false
			for _, arg := range args {
				if strings.Contains(strings.ToLower(pair[0]), strings.ToLower(arg)) {
					found = true
					break
				}
			}
			if found || len(args) == 0 {
				fmt.Println(pair[0], ": ", pair[1])
			}
		}
	},
}
