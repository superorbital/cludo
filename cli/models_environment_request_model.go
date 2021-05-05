// Code generated by go-swagger; DO NOT EDIT.

package cli

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superorbital/cludo/models"
)

// Schema cli for ModelsEnvironmentRequest

// register flags to command
func registerModelModelsEnvironmentRequestFlags(depth int, cmdPrefix string, cmd *cobra.Command) error {

	if err := registerModelsEnvironmentRequestRoleID(depth, cmdPrefix, cmd); err != nil {
		return err
	}

	return nil
}

func registerModelsEnvironmentRequestRoleID(depth int, cmdPrefix string, cmd *cobra.Command) error {
	if depth > maxDepth {
		return nil
	}

	// warning: roleID []string array type is not supported by go-swagger cli yet

	return nil
}

// retrieve flags from commands, and set value in model. Return true if any flag is passed by user to fill model field.
func retrieveModelModelsEnvironmentRequestFlags(depth int, m *models.ModelsEnvironmentRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	retAdded := false

	err, roleIdAdded := retrieveModelsEnvironmentRequestRoleIDFlags(depth, m, cmdPrefix, cmd)
	if err != nil {
		return err, false
	}
	retAdded = retAdded || roleIdAdded

	return nil, retAdded
}

func retrieveModelsEnvironmentRequestRoleIDFlags(depth int, m *models.ModelsEnvironmentRequest, cmdPrefix string, cmd *cobra.Command) (error, bool) {
	if depth > maxDepth {
		return nil, false
	}
	retAdded := false

	roleIdFlagName := fmt.Sprintf("%v.roleID", cmdPrefix)
	if cmd.Flags().Changed(roleIdFlagName) {
		// warning: roleID array type []string is not supported by go-swagger cli yet
	}

	return nil, retAdded
}
