package objects

type Model interface {
	GetSomething() (*ExampleStructInOtherFileAndPackage, error)
	GetSomethingArray() ([]*ExampleStructInOtherFileAndPackage, error)
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
