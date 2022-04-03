package outside2

type ExportType string
type ExportStatus string
type ExportFormat string

type DateRange struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}
