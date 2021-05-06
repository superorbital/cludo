/*
Copyright © 2021 SuperOrbital, LLC <info@superorbital.io>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// FIXME: Completely broken code, but here for the moment as it would be good to re-create as we get the basic code working in the new layout.

package clientCLI

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecedence(t *testing.T) {
	// Run the tests in a temporary directory
	tmpDir, err := ioutil.TempDir("", "cludo")
	require.NoError(t, err, "error creating a temporary test directory")
	testDir, err := os.Getwd()
	require.NoError(t, err, "error getting the current working directory")
	defer os.Chdir(testDir)
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "error changing to the temporary test directory")

	// Set favorite-color with the config file
	t.Run("config file", func(t *testing.T) {
		// Copy the config file into our temporary test directory
		configB, err := ioutil.ReadFile(filepath.Join(testDir, "../data/cludo.yaml"))
		require.NoError(t, err, "error reading test config file")
		err = os.Mkdir(".cludo", 0755)
		require.NoError(t, err, "error creating temp config directory")
		err = ioutil.WriteFile(filepath.Join(tmpDir, ".cludo/cludo.yaml"), configB, 0644)
		require.NoError(t, err, "error writing test config file")
		defer os.Remove(filepath.Join(tmpDir, ".cludo/cludo.yaml"))

		// Run ./cludo
		cmd, err := MakeRootCmd()
		require.NoError(t, err, "RootCmd construction error:")
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"exec", "aws", "ec2", "describe-instances"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `exec ran!
server_url: https://www.example.com/
key_path  : ~/.ssh/id_rsa
command   : aws ec2 describe-instances
`

		assert.Equal(t, wantOutput, gotOutput, "expected the server url from the config file and the key_path from the flag default")
	})

	// Set favorite-color with an environment variable
	t.Run("env var", func(t *testing.T) {
		// Run CLUDO_SERVER_URL=https://cludo.example.com:8443/ ./cludo
		os.Setenv("CLUDO_SERVER_URL", "https://cludo.example.com:8443/")
		defer os.Unsetenv("CLUDO_SERVER_URL")

		cmd, err := MakeRootCmd()
		require.NoError(t, err, "RootCmd construction error:")
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"exec", "aws", "ec2", "describe-instances"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `exec ran!
server_url: https://cludo.example.com:8443/
key_path  : ~/.ssh/id_rsa
command   : aws ec2 describe-instances
`

		assert.Equal(t, wantOutput, gotOutput, "expected the server_url to use the environment variable value and the key_path to use the flag default")
	})

	// Set number with a flag
	t.Run("flag", func(t *testing.T) {
		// Run ./cludo exec --key-path="~/.ssh/id_ed25519"
		cmd, err := MakeRootCmd()
		require.NoError(t, err, "RootCmd construction error:")
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"exec", "--key-path", "~/.ssh/id_ed25519", "aws", "ec2", "describe-instances"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `exec ran!
server_url: https://cludo.example.com:8443/
key_path  : ~/.ssh/id_ed25519
command   : aws ec2 describe-instances
`

		assert.Equal(t, wantOutput, gotOutput, "expected the key_path to use the flag value and the server_url to use the flag default")
	})
}