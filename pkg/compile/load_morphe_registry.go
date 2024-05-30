package compile

import (
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
)

func loadMorpheRegistry(hooks hook.LoadMorpheRegistry, config cfg.MorpheLoadRegistryConfig) (*registry.Registry, error) {
	config, loadStartErr := triggerLoadRegistryStart(hooks, config)
	if loadStartErr != nil {
		return nil, triggerLoadRegistryFailure(hooks, config, nil, loadStartErr)
	}

	r := registry.NewRegistry()

	loadErr := loadConfiguredRegistry(config, r)
	if loadErr != nil {
		return nil, triggerLoadRegistryFailure(hooks, config, r, loadErr)
	}

	loadSuccessErr := triggerLoadRegistrySuccess(hooks, r)
	if loadSuccessErr != nil {
		return nil, triggerLoadRegistryFailure(hooks, config, r, loadSuccessErr)
	}

	return r, nil
}

func loadConfiguredRegistry(config cfg.MorpheLoadRegistryConfig, r *registry.Registry) error {
	modelsErr := r.LoadModelsFromDirectory(config.RegistryModelsDirPath)
	if modelsErr != nil {
		return modelsErr
	}

	entitiesErr := r.LoadEntitiesFromDirectory(config.RegistryEntitiesDirPath)
	if entitiesErr != nil {
		return entitiesErr
	}

	return nil
}

func triggerLoadRegistryStart(hooks hook.LoadMorpheRegistry, config cfg.MorpheLoadRegistryConfig) (cfg.MorpheLoadRegistryConfig, error) {
	if hooks.OnRegistryLoadStart == nil {
		return config, nil
	}

	updatedConfig, startErr := hooks.OnRegistryLoadStart(config)
	if startErr != nil {
		return cfg.MorpheLoadRegistryConfig{}, startErr
	}

	return updatedConfig, nil
}

func triggerLoadRegistrySuccess(hooks hook.LoadMorpheRegistry, r *registry.Registry) error {
	if hooks.OnRegistryLoadSuccess == nil {
		return nil
	}
	if r == nil {
		return ErrRegistryNotInitialized
	}
	registry, successErr := hooks.OnRegistryLoadSuccess(*r)
	if successErr != nil {
		return successErr
	}
	*r = registry
	return nil
}

func triggerLoadRegistryFailure(hooks hook.LoadMorpheRegistry, config cfg.MorpheLoadRegistryConfig, r *registry.Registry, failureErr error) error {
	if hooks.OnRegistryLoadFailure == nil {
		return failureErr
	}

	if r == nil {
		return hooks.OnRegistryLoadFailure(config, registry.Registry{}, failureErr)
	}

	return hooks.OnRegistryLoadFailure(config, *r, failureErr)
}
