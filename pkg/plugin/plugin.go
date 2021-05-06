package plugin

import "github.com/superorbital/cludo/models"

type Plugin interface {
	GenerateEnvironment() (*models.ModelsEnvironmentResponse, error)
}

type PluginRegistry map[string]*Plugin
