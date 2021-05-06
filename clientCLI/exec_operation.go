package clientCLI

import (
	"fmt"
	"strings"

	"github.com/superorbital/cludo/pkg/config"

	"github.com/spf13/cobra"
	//"github.com/superorbital/cludo/models"
)

// runOperationExec uses cmd flags to call endpoint api
func runOperationExec(cmd *cobra.Command, args []string, conf config.ClientConfig) string {
	//appCli, err := makeClient(cmd, args)
	//if err != nil {
	//	return err
	//}
	// retrieve flag values from cmd and fill params
	//params := environment.NewGenerateEnvironmentParams()
	//if err, _ := retrieveOperationExecGenerateExecBodyFlag(params, //"", cmd); err != nil {
	//	return err
	//}

	if debug {
		fmt.Printf("\n%#v\n", conf)
	}

	if dryRun {

		logDebugf("dry-run flag specified. Skip sending request.")
		return ""
	}
	// make request and then print result
	msgStr, err := "\nexec ran!\nserver_url: "+conf.ServerURL+"\nkey_path  : "+conf.KeyPath+"\ncommand   : "+string(strings.Join(args, " ")), *new(error)

	if err != nil {
		return ""
	}

	return msgStr
}
