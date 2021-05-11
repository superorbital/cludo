package cludo

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
	tmpDir, err := ioutil.TempDir("", "cludo-*")
	require.NoError(t, err, "error creating a temporary test directory")
	testDir, err := os.Getwd()
	require.NoError(t, err, "error getting the current working directory")
	defer os.Chdir(testDir)
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "error changing to the temporary test directory")
	// Copy the SSH key w/ passphrase into our temporary test directory
	keyPP, err := ioutil.ReadFile(filepath.Join(testDir, "../data/id_rsa_TEST_KEY_PP"))
	require.NoError(t, err, "error reading test file: SSH key w/ passphrase")
	err = ioutil.WriteFile(filepath.Join(tmpDir, "id_rsa_TEST_KEY_PP"), keyPP, 0644)
	require.NoError(t, err, "error writing file test: SSH key w/ passphrase")
	defer os.Remove(filepath.Join(tmpDir, "id_rsa_TEST_KEY_PP"))
	// Copy the SSH key w/o passphrase into our temporary test directory
	keyNOPP, err := ioutil.ReadFile(filepath.Join(testDir, "../data/id_rsa_TEST_KEY_NOPP"))
	require.NoError(t, err, "error reading test file: SSH key w/o passphrase")
	err = ioutil.WriteFile(filepath.Join(tmpDir, "id_rsa_TEST_KEY_NOPP"), keyNOPP, 0644)
	require.NoError(t, err, "error writing file test: SSH key w/o passphrase")
	defer os.Remove(filepath.Join(tmpDir, "id_rsa_TEST_KEY_NOPP"))

	// Setup the config file
	t.Run("config file", func(t *testing.T) {
		// Copy the config file into our temporary test directory
		configB, err := ioutil.ReadFile(filepath.Join(testDir, "../data/cludo.yaml"))
		require.NoError(t, err, "error reading test file: cludo.yaml")
		err = os.Mkdir(".cludo", 0755)
		require.NoError(t, err, "error creating temp config directory")
		err = ioutil.WriteFile(filepath.Join(tmpDir, ".cludo/cludo.yaml"), configB, 0644)
		require.NoError(t, err, "error writing file test: cludo.yaml")
		defer os.Remove(filepath.Join(tmpDir, ".cludo/cludo.yaml"))
		defer os.Remove(filepath.Join(tmpDir, ".cludo"))

		// Run ./cludo
		cmd, err := MakeRootCmd()
		require.NoError(t, err, "RootCmd construction error:")
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"exec", "aws", "ec2", "describe-instances"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `
exec ran!
server_url: http://127.0.0.1:8080/
shell_path: /bin/sh
command   : aws ec2 describe-instances
`

		assert.Equal(t, wantOutput, gotOutput, "expected the server url from the config file and the shell_path from the flag default")
	})

	// Set server-url with an environment variable
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
		wantOutput := `
exec ran!
server_url: https://cludo.example.com:8443/
shell_path: /bin/sh
command   : aws ec2 describe-instances
`

		assert.Equal(t, wantOutput, gotOutput, "expected the server_url to use the environment variable value and the shell_path to use the flag default")
	})

	// Set shell-path with a flag
	t.Run("flag", func(t *testing.T) {
		// Run ./cludo exec --shell-path="/usr/local/bin/bash"
		cmd, err := MakeRootCmd()
		require.NoError(t, err, "RootCmd construction error:")
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		// FIXME: I strongly suspect that these tests are being run in parallel and
		// this might be causing issues.
		cmd.SetArgs([]string{"exec", "--shell-path", "/usr/local/bin/bash", "aws", "ec2", "describe-instances"})
		cmd.Execute()

		gotOutput := output.String()
		wantOutput := `
exec ran!
server_url: http://localhost:80/
shell_path: /usr/local/bin/bash
command   : aws ec2 describe-instances
`

		assert.Equal(t, wantOutput, gotOutput, "expected the shell_path to use the flag value and the server_url to use the default flag value")
	})
}
