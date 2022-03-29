package objects

type Model interface {
	GetSomething() (*ExampleStructInOtherFileAndPackage, error)
	GetSomethingArray() ([]*ExampleStructInOtherFileAndPackage, error)
	GetSomethingNested() (*ExampleEmbeddedParent, error)
	Errors() ([]error, error)
}

type model struct {
}

func NewModel() Model {
	return &model{}
}

func (m model) GetSomething() (*ExampleStructInOtherFileAndPackage, error) {
	return &ExampleStructInOtherFileAndPackage{}, nil
}

func (m model) GetSomethingArray() ([]*ExampleStructInOtherFileAndPackage, error) {
	return []*ExampleStructInOtherFileAndPackage{}, nil
}

func (m model) GetSomethingNested() (*ExampleEmbeddedParent, error) {
	return &ExampleEmbeddedParent{}, nil
}

func (m model) Errors() ([]error, error) {
	return nil, nil
}
