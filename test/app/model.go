package app

type Model interface {
	GetSomething() (*request, error)
}

type model struct {
}

func newModel() Model {
	return &model{}
}

func (m model) GetSomething() (*request, error) {
	return &request{}, nil
}
