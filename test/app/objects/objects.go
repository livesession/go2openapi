package objects

import "time"

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

type ExampleEmbeddedChild struct {
	Lastname string `json:"lastname"`
}

type RecordingElement struct {
	Element          string    `json:"element"`
	MatchType        string    `json:"match_type"`
	Enabled          bool      `json:"enabled"`
	CreatedByAgentID string    `json:"created_by_agent_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type RecordingElementProjection struct {
	*RecordingElement
	CreatedBy string `json:"created_by"`
}

type RecordingElementsProjection struct {
	List   []*RecordingElementProjection `json:"list"`
	Deeper *RecordingElementProjection   `json:"deeper"`
}

type ExampleEmbeddedParent struct {
	*ExampleEmbeddedChild
	Firstname         string                       `json:"firstname"`
	RecordingElements *RecordingElementsProjection `json:"example"`
}
