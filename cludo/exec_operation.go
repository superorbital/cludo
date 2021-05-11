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

	if debug {
		fmt.Printf("\n%#v\n", cc)
	}

	u, err := url.Parse(cc.ServerURL)
	if err != nil {
		return fmt.Sprintf("[ERROR] Unable to parse server url. (%s)", err)
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
		return fmt.Sprintf("[ERROR] Could not expand SSH Private Key path. (%s)", err)
	}

	rawKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Sprintf("[ERROR] Reading SSH Private Key. (%s)", err)
	}

	key, err := auth.DecodePrivateKey(rawKey, []byte(cc.Passphrase))
	if err != nil {
		return fmt.Sprintf("[ERROR] Decoding SSH Private Key. (%s)", err)
	}

	signer := auth.NewDefaultSigner(key)

	if dryRun {
		logDebugf("dry-run flag specified. Skip sending request.")
		return ""
	}

	// Let's make sure we can talk to the server now
	_, err = client.Environment.GenerateEnvironment(params, signer.CludoAuth())
	if err != nil {
		//return fmt.Sprintf("[ERROR] Unable to generate environment. (%s)", err)
	}

	// At the moment seeing this error from the above code:
	// [ERROR] Unable to generate environment. (Post "http://127.0.0.1:8080/environment": net/http: invalid header field value "\x92\u007f\xbd\xb3H6)B\nH:mtVduNdGxp6gh96ULfHVPKnA9DAIeJ6QPkk3eU4/WDR7APrEqrt0w8Ey4QcPMQzSXeXR0Hixpzd9U6ruC9OIHbyawQDS3lm+ho0X1pFeHmGNQSGJkFitsyYjbV42y0wrkbXaww9LhYgt7OCKg0wLRAm2ezANuKwVLJDYGeQnON6BukMEfSfUvQElAnCmbXfyj6Qet+pKKxUV6aifjX/JYiAV9IK3+A3HmiiVQCAIyNfZxiCr7NTSR3OqEIxtH/upsK1Lrvmtz9h/ug9OVYeDXgOU1Vd+PqaTAzoPg3C1B/zBqnqMtG/Ucw3KPhXTe9aIp++cfPHiJESgTgA7y+xq9A==" for key X-Cludo-Key)

	// make request and then print result
	msgStr := "\nexec ran!\nserver_url: " + cc.ServerURL + "\nshell_path: " + cc.ShellPath + "\ncommand   : " + string(strings.Join(args, " "))

	return msgStr
}
