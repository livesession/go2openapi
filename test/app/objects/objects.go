package objects

import (
	"time"

	"github.com/livesession/restflix/test/app/objects/outside"
	apikit_objects "github.com/livesession/restflix/test/app/objects/outside2"
)

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

type CustomType string
type CustomType2 ExampleEmbeddedParent
type CustomType3 outside.Outside
type CustomType4 *outside.Outside
type CustomType5 outside.String

type CustomStruct struct {
	Custom  CustomType  `json:"custom"`
	Custom2 CustomType2 `json:"custom2"`
	Custom3 CustomType3 `json:"custom3"`
	Custom4 CustomType4 `json:"custom4"`
	Custom5 CustomType5 `json:"custom5"`
}

type CreateExportRequest struct {
	WebsiteID string                      `json:"website_id" valid:"required"`
	Name      string                      `json:"name" valid:"required"`
	DateRange *apikit_objects.DateRange   `json:"date_range"`
	Format    apikit_objects.ExportFormat `json:"format" valid:"in(csv|json),required"`
	Type      apikit_objects.ExportType   `json:"type" valid:"in(visitors|events),required"`
	//Filters   *apikit_objects.Filters     `json:"filters"`
}
