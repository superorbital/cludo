// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/superorbital/cludo/client/role"

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
	// make request and then print result
	if err := printOperationRoleListRolesResult(appCli.Role.ListRoles(params)); err != nil {
		return err
	}
	return nil
}

// printOperationRoleListRolesResult prints output to stdout
func printOperationRoleListRolesResult(resp0 *role.ListRolesOK, respErr error) error {
	if respErr != nil {

		var iResp interface{} = respErr
		defaultResp, ok := iResp.(*role.ListRolesDefault)
		if !ok {
			return respErr
		}
		if defaultResp.Payload != nil {
			msgStr, err := json.Marshal(defaultResp.Payload)
			if err != nil {
				return err
			}
			fmt.Println(string(msgStr))
			return nil
		}

		return respErr
	}

	if resp0.Payload != "" {
		msgStr, err := json.Marshal(resp0.Payload)
		if err != nil {
			return err
		}
		fmt.Println(string(msgStr))
	}

	return nil
}

// registerOperationRoleListRolesParamFlags registers all flags needed to fill params
func registerOperationRoleListRolesParamFlags(cmd *cobra.Command) error {
	return nil
}
