package objects

type ExampleStructInOtherFileAndPackage struct {
	Name string `json:"name"`
}

type ExampleStructInOtherFileAndPackageNested struct {
	X []*ExampleStructInOtherFileAndPackage `json:"x"`
}

type CreateAgentRequest struct {
	Email string `json:"email,omitempty" valid:"runelength(5|255),email,required"`
	Role  string `json:"role,omitempty" valid:"in(agent|admin),required"`
}

type CreateAgentsBatchRequest struct {
	Agents []*CreateAgentRequest `json:"agents" valid:"requiredArrayOfStructs~agents[].email and agents[].role is required or invalid"`
}
