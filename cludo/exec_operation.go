package cludo

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/strfmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/superorbital/cludo/client"
	"github.com/superorbital/cludo/client/environment"
	"github.com/superorbital/cludo/pkg/auth"
	"github.com/superorbital/cludo/pkg/config"
)

// runOperationExec uses cmd flags to call endpoint api
func runOperationExec(cmd *cobra.Command, args []string, conf config.Config) string {

	cc := conf.Client["default"]
	u, err := url.Parse(cc.ServerURL)
	if err != nil {
		fmt.Printf("[ERROR] Unable to parse server url. (%s)", err)
	}
	hostport := u.Host
	apipath := u.Path
	scheme := u.Scheme
	transport := httptransport.New(hostport, apipath, []string{scheme})
	transport.Consumers["application/json"] = runtime.JSONConsumer()
	transport.Producers["application/json"] = runtime.JSONProducer()
	// Insantiate the client that we will use to talk to the Todo server
	client := client.New(transport, strfmt.Default)

	params := environment.NewGenerateEnvironmentParams()

	keyPath, err := homedir.Expand(cc.KeyPath)
	if err != nil {
		fmt.Printf("[ERROR] Could not expand SSH Private Key path. (%s)", err)
		return ""
	}

	rawKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Printf("[ERROR] Reading SSH Private Key. (%s)", err)
		return ""
	}

	key, err := auth.DecodePrivateKey(rawKey, []byte(cc.Passphrase))
	if err != nil {
		fmt.Printf("[ERROR] Decoding SSH Private Key. (%s)", err)
		return ""
	}

	signer := auth.NewDefaultSigner(key)

	// Let's make sure we can talk to the server now
	_, err = client.Environment.GenerateEnvironment(params, signer.CludoAuth())
	if err != nil {
		fmt.Printf("[ERROR] Unable to generate environment. (%s)", err)
		return ""
	}

	//return client, nil

	if debug {
		fmt.Printf("\n%#v\n", cc)
	}

	if dryRun {

		logDebugf("dry-run flag specified. Skip sending request.")
		return ""
	}

	// make request and then print result
	msgStr, err := "\nexec ran!\nserver_url: "+cc.ServerURL+"\nkey_path  : "+cc.KeyPath+"\ncommand   : "+string(strings.Join(args, " ")), *new(error)
	//msgStr, err := "\nexec ran!\nserver_url: "+conf["server_url"]+"\nkey_path  : "+conf["key_path"]+"\ncommand   : "+string(strings.Join(args, " ")), *new(error)

	if err != nil {
		return ""
	}

	return msgStr
}
