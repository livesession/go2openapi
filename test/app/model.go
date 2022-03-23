package app

type Model interface {
	GetSomething() (*response2, error)
}

type model struct {
}

func newModel() Model {
	return &model{}
}

func (m model) GetSomething() (*response2, error) {
	return &response2{}, nil
}
