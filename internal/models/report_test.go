package models

import "testing"

func TestReportTypeNormalize(t *testing.T) {
	tests := []struct {
		name       string
		reportType ReportType
		want       ReportType
	}{
		{name: "empty defaults to status", reportType: "", want: ReportByStatus},
		{name: "explicit priority stays priority", reportType: ReportByPriority, want: ReportByPriority},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reportType.Normalize(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
