package models

type ReportType string

const (
	ReportByStatus   ReportType = "status"
	ReportByPriority ReportType = "priority"
)

type Report struct {
	Type  ReportType
	Items []ReportItem
}

type ReportItem struct {
	Label string
	Count int
}

func (t ReportType) Normalize() ReportType {
	if t == "" {
		return ReportByStatus
	}
	return t
}
