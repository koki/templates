package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/glog"
	"github.com/koki/templates/json"
	"github.com/koki/templates/yaml"
	"github.com/spf13/cobra"
)

// RootCmd root cobra command.
var RootCmd = &cobra.Command{
	Use: "templates <subcommand>",
	Short: "templates fills in ${VARIABLE_NAME} from a params file " +
		"with the field VARIABLE_NAME",
}

func doFill(templatePath, paramsPath string, fill func([]byte, []byte) ([]byte, error)) {
	var template, params, filled []byte
	var err error
	template, err = ioutil.ReadFile(templatePath)
	if err != nil {
		glog.Fatalf("couldn't load template file (%s): %v",
			templatePath, err)
	}

	params, err = ioutil.ReadFile(paramsPath)
	if err != nil {
		glog.Fatalf("couldn't load params file (%s): %v",
			paramsPath, err)
	}

	filled, err = yaml.Fill(template, params)
	if err != nil {
		glog.Fatalf("couldn't fill template: %v", err)
	}

	fmt.Println(string(filled))
}

var jsonCmd = &cobra.Command{
	Use:   "json <template.json> <parameters.json>",
	Short: "fill a json template from json params",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		doFill(args[0], args[1], json.Fill)
	},
}

var yamlCmd = &cobra.Command{
	Use:   "yaml <template.yaml> <parameters.yaml>",
	Short: "fill a yaml template from yaml params",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		doFill(args[0], args[1], yaml.Fill)
	},
}

func init() {
	RootCmd.AddCommand(jsonCmd, yamlCmd)
}
