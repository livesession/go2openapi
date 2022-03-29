package app

func (a *api) methodErrorsV2() ([]error, error) {
	return a.modelOutside.Errors()
}
