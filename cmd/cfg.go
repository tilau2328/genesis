package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

var cfgCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c"},
}

func init() {
	addPathFlag(unsetCfgCmd, listCfgCmd, setCfgCmd)
	cfgCmd.AddCommand(unsetCfgCmd)
	cfgCmd.AddCommand(listCfgCmd)
	cfgCmd.AddCommand(setCfgCmd)
}

func addPathFlag(commands ...*cobra.Command) {
	for _, cmd := range commands {
		cmd.Flags().StringP("path", "p", ".", "Path to look")
	}
}

var listCfgCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return
		}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			fmt.Println(f.Name())
		}
	},
}

var setCfgCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := cmd.Flags().GetString("path")
		if err != nil {
			return errors.New("invalid path")
		}
		if len(args) < 3 || len(args)%2 == 0 {
			return errors.New("insufficient args")
		}
		p = path.Join(p, args[0])
		viper.SetConfigFile(p)

		exists, err := afero.Exists(afero.NewOsFs(), p)
		if err != nil {
			return err
		}
		if exists {
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		}

		args = args[1:]
		for i := 0; i < len(args); i += 2 {
			values := strings.Split(args[i+1], ",")
			parts := strings.Split(args[i], ".")
			value := values[0]
			key := parts[0]

			if len(parts) == 1 {
				if len(values) == 1 {
					viper.Set(key, value)
					return nil
				}
				viper.Set(key, values)
				return nil
			}

			// Get previous value
			before := viper.Get(key)

			// Check if it's a map
			cursor, ok := before.(map[string]interface{})
			if !ok {
				// Make it one if it ain't
				cursor = make(map[string]interface{})
			}

			root := cursor
			parts = parts[1:]
			for i, part := range parts {
				// Last Step
				if i == len(parts)-1 {
					// Store single sub value
					if len(values) == 1 {
						cursor[part] = value
						break
					}
					// Store multiple sub values
					cursor[part] = values
					break
				}
				// Check if subpart already exists and is map
				cast, ok := cursor[part].(map[string]interface{})
				if !ok {
					// Make it a map if it ain't
					cast = make(map[string]interface{})
					cursor[part] = cast
				}
				// Move cursor«
				cursor = cast
			}
			log.Println(key)
			viper.Set(key, root)
		}
		return viper.WriteConfig()
	},
}

var unsetCfgCmd = &cobra.Command{
	Use:     "unset",
	Aliases: []string{"u"},
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := cmd.Flags().GetString("path")
		if err != nil {
			return errors.New("invalid path")
		}

		if len(args) < 2 {
			return errors.New("insufficient args")
		}
		p = path.Join(p, args[0])
		viper.SetConfigFile(p)

		exists, err := afero.Exists(afero.NewOsFs(), p)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("missing config")
		}

		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		for _, key := range args[1:] {
			parts := strings.Split(key, ".")
			if len(parts) == 1 {
				viper.Set(key, nil)
				return nil
			}

			key = parts[0]

			// Get previous value
			before := viper.Get(key)

			// Check if it's a map
			cursor, ok := before.(map[string]interface{})
			if !ok {
				return nil
			}

			root := cursor
			parts = parts[1:]
			for i, part := range parts {
				// Last Step
				if i == len(parts)-1 {
					delete(cursor, part)
					break
				}
				// Check if subpart already exists and is map
				cast, ok := cursor[part].(map[string]interface{})

				if !ok {
					return nil
				}
				// Move cursor
				cursor = cast
			}
			viper.Set(key, root)
		}
		return viper.WriteConfig()
	},
}
