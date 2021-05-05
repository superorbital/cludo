// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/superorbital/cludo/client/role"

	"github.com/go-openapi/swag"
	"github.com/spf13/cobra"
)

// makeOperationRoleListRolesCmd returns a cmd to handle operation listRoles
func makeOperationRoleListRolesCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "list-roles",
		Short: `List all roles available to current user`,
		RunE:  runOperationRoleListRoles,
	}

	if err := registerOperationRoleListRolesParamFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// runOperationRoleListRoles uses cmd flags to call endpoint api
func runOperationRoleListRoles(cmd *cobra.Command, args []string) error {
	appCli, err := makeClient(cmd, args)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := role.NewListRolesParams()
	if dryRun {

		logDebugf("dry-run flag specified. Skip sending request.")
		return nil
	}
	// make request and then print result
	msgStr, err := parseOperationRoleListRolesResult(appCli.Role.ListRoles(params))
	if err != nil {
		return err
	}
	if !debug {

		fmt.Println(msgStr)
	}
	return nil
}

// registerOperationRoleListRolesParamFlags registers all flags needed to fill params
func registerOperationRoleListRolesParamFlags(cmd *cobra.Command) error {
	return nil
}

// parseOperationRoleListRolesResult parses request result and return the string content
func parseOperationRoleListRolesResult(resp0 *role.ListRolesOK, respErr error) (string, error) {
	if respErr != nil {

		var iRespD interface{} = respErr
		respD, ok := iRespD.(*role.ListRolesDefault)
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
		resp0, ok := iResp0.(*role.ListRolesOK)
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
		resp1, ok := iResp1.(*role.ListRolesBadRequest)
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
		msgStr := fmt.Sprintf("%v", resp0.Payload)
		return string(msgStr), nil
	}

	return "", nil
}
