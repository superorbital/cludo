package cludo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/superorbital/cludo/pkg/config"

	// "github.com/go-openapi/strfmt"
	// "github.com/superorbital/cludo/client"
	// "github.com/superorbital/cludo/client/environment"
	"github.com/spf13/cobra"
)

// runOperationExec uses cmd flags to call endpoint api
func runOperationExec(cmd *cobra.Command, args []string, conf config.ClientConfig) string {

	u, err := url.Parse(conf.ServerURL)
	if err != nil {
		fmt.Println("[ERROR] Unable to parse server url.")
	}
	hostport := u.Host
	apipath := u.Path
	scheme := u.Scheme
	transport := httptransport.New(hostport, apipath, []string{scheme})
	transport.Consumers["application/json"] = runtime.JSONConsumer()
	transport.Producers["application/json"] = runtime.JSONProducer()
	// Insantiate the client that we will use to talk to the Todo server
	//client := client.New(transport, strfmt.Default)

	//params := environment.NewGenerateEnvironmentParams()

	//auth, err := conf.NewDefaultSigner()
	// Let's make sure we can talk to the server now
	//_, err = client.Environment.GenerateEnvironment(params, authInfo, opts)
	//if err != nil {
	//	fmt.Println("[ERROR] Unable to generate environment.")
	//}

	//return client, nil

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
