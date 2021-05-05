// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/superorbital/cludo/client/environment"
	"github.com/superorbital/cludo/models"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationEnvironmentGenerateEnvironmentCmd returns a cmd to handle operation generateEnvironment
func makeOperationEnvironmentGenerateEnvironmentCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "generate-environment",
		Short: `Generate a temporary environment (set of environment variables)`,
		RunE:  runOperationEnvironmentGenerateEnvironment,
	}

	if err := registerOperationEnvironmentGenerateEnvironmentParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationEnvironmentGenerateEnvironment uses cmd flags to call endpoint api
func runOperationEnvironmentGenerateEnvironment(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := environment.NewGenerateEnvironmentParams()
	if err, _ := retrieveOperationEnvironmentGenerateEnvironmentBodyFlag(params, "", cmd); err != nil {
		return err
	}
	if dryRun {

		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationEnvironmentGenerateEnvironmentResult(appCli.Environment.GenerateEnvironment(params))
	if err != nil {
		return err
	}
	if !debug {

		fmt.Println(msgStr)
	}
	return nil
}

// registerOperationEnvironmentGenerateEnvironmentParamFlags registers all flags needed to fill params
func registerOperationEnvironmentGenerateEnvironmentParamFlags(cmd *cobra.Command) error {
	if err := registerOperationEnvironmentGenerateEnvironmentBodyParamFlags("", cmd); err != nil {
		return err
	}
	return nil
}

func registerOperationEnvironmentGenerateEnvironmentBodyParamFlags(cmdPrefix string, cmd *cobra.Command) error {

	var bodyFlagName string
	if cmdPrefix == "" {
		bodyFlagName = "body"
	} else {
		bodyFlagName = fmt.Sprintf("%v.body", cmdPrefix)
	}

	_ = cmd.PersistentFlags().String(bodyFlagName, "", "Optional json string for [body]. Temporary Environment Request definition")

	// add flags for body
	if err := registerModelModelsEnvironmentRequestFlags(0, "modelsEnvironmentRequest", cmd); err != nil {
		return err
	}

	return nil
}

func retrieveOperationEnvironmentGenerateEnvironmentBodyFlag(m *environment.GenerateEnvironmentParams, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false
	if cmd.Flags().Changed("body") {
		// Read body string from cmd and unmarshal
		bodyValueStr, err := cmd.Flags().GetString("body")
		if err != nil {
			return err, false
		}

		bodyValue := models.ModelsEnvironmentRequest{}
		if err := json.Unmarshal([]byte(bodyValueStr), &bodyValue); err != nil {
			return fmt.Errorf("cannot unmarshal body string in models.ModelsEnvironmentRequest: %v", err), false
		}
		m.Body = &bodyValue
	}
	bodyValueModel := m.Body
	if swag.IsZero(bodyValueModel) {
		bodyValueModel = &models.ModelsEnvironmentRequest{}
	}
	err, added := retrieveModelModelsEnvironmentRequestFlags(0, bodyValueModel, "modelsEnvironmentRequest", cmd)
	if err != nil {
		return err, false
	}
	if added {
		m.Body = bodyValueModel
	}
	if dryRun && debug {

		bodyValueDebugBytes, err := json.Marshal(m.Body)
		if err != nil {
			return err, false
		}
		logDebugf("Body dry-run payload: %v", string(bodyValueDebugBytes))
	}
	retAdded = retAdded || added

	return nil, retAdded
}

// parseOperationEnvironmentGenerateEnvironmentResult parses request result and return the string content
func parseOperationEnvironmentGenerateEnvironmentResult(resp0 *environment.GenerateEnvironmentOK, respErr error) (string, error) {
	if respErr != nil {

		var iRespD interface{} = respErr
		respD, ok := iRespD.(*environment.GenerateEnvironmentDefault)
		if ok {
			if !swag.IsZero(respD.Payload) {
				msgStr, err := json.Marshal(respD.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp0 interface{} = respErr
		resp0, ok := iResp0.(*environment.GenerateEnvironmentOK)
		if ok {
			if !swag.IsZero(resp0.Payload) {
				msgStr, err := json.Marshal(resp0.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		var iResp1 interface{} = respErr
		resp1, ok := iResp1.(*environment.GenerateEnvironmentBadRequest)
		if ok {
			if !swag.IsZero(resp1.Payload) {
				msgStr, err := json.Marshal(resp1.Payload)
				if err != nil {
					return "", err
				}
				return string(msgStr), nil
			}
		}

		return "", respErr
	}

	if !swag.IsZero(resp0.Payload) {
		msgStr, err := json.Marshal(resp0.Payload)
		if err != nil {
			return "", err
		}
		return string(msgStr), nil
	}

	return "", nil
}
