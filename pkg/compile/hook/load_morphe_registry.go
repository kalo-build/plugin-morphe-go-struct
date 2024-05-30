package hook

type LoadMorpheRegistry struct {
	OnRegistryLoadStart   OnRegistryLoadStartHook
	OnRegistryLoadSuccess OnRegistryLoadSuccessHook
	OnRegistryLoadFailure OnRegistryLoadFailureHook
}
