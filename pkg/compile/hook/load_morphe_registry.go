package hook

import (
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
)

type LoadMorpheRegistry struct {
	OnRegistryLoadStart   OnRegistryLoadStartHook
	OnRegistryLoadSuccess OnRegistryLoadSuccessHook
	OnRegistryLoadFailure OnRegistryLoadFailureHook
}

type OnRegistryLoadStartHook = func(config cfg.MorpheLoadRegistryConfig) (cfg.MorpheLoadRegistryConfig, error)
type OnRegistryLoadSuccessHook = func(registry registry.Registry) (registry.Registry, error)
type OnRegistryLoadFailureHook = func(config cfg.MorpheLoadRegistryConfig, registry registry.Registry, loadFailure error) error
