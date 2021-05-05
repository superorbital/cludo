package clientCLI

import (
	"fmt"
	"strings"

	"github.com/superorbital/cludo/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//"github.com/superorbital/cludo/models"
)

// runOperationExec uses cmd flags to call endpoint api
func runOperationExec(cmd *cobra.Command, args []string) error {
	//appCli, err := makeClient(cmd, args)
	//if err != nil {
	//	return err
	//}
	// retrieve flag values from cmd and fill params
	//params := environment.NewGenerateEnvironmentParams()
	//if err, _ := retrieveOperationExecGenerateExecBodyFlag(params, //"", cmd); err != nil {
	//	return err
	//}

	conf := &config.Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	if debug {
		fmt.Printf("%#v", conf.Client["default"])
	}

	if dryRun {

		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := "exec ran! (server: "+conf.Client["default"].ServerURL+"command: "+string(strings.Join(args, " ")), *new(error)

	if err != nil {
		return err
	}
	if !debug {

		fmt.Println(msgStr)
	} else {
		fmt.Println(msgStr)
	}
	return nil
}
