package service

import (
	"context"
	"testing"

	"taskflow/internal/models"
)

type fakeReportTasks struct {
	tasks []models.Task
}

func (f fakeReportTasks) Create(ctx context.Context, task *models.Task) error {
	return nil
}

func (f fakeReportTasks) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	return f.tasks, nil
}

func (f fakeReportTasks) FindByID(ctx context.Context, id int64) (models.Task, error) {
	return models.Task{}, nil
}

func (f fakeReportTasks) UpdateStatus(ctx context.Context, task models.Task, history models.TaskHistory) error {
	return nil
}

func TestSelectReportStrategy(t *testing.T) {
	tests := []struct {
		reportType models.ReportType
		wantType   models.ReportType
	}{
		{reportType: models.ReportByStatus, wantType: models.ReportByStatus},
		{reportType: models.ReportByPriority, wantType: models.ReportByPriority},
		{reportType: models.ReportByAssignee, wantType: models.ReportByAssignee},
	}

	tasks := []models.Task{
		{Status: models.StatusNew, Priority: models.PriorityHigh, AssigneeID: 1},
		{Status: models.StatusDone, Priority: models.PriorityHigh, AssigneeID: 2},
		{Status: models.StatusNew, Priority: models.PriorityLow, AssigneeID: 1},
	}

	for _, tt := range tests {
		strategy, err := SelectReportStrategy(tt.reportType)
		if err != nil {
			t.Fatalf("expected strategy for %q: %v", tt.reportType, err)
		}
		if strategy.ReportType() != tt.wantType {
			t.Fatalf("expected strategy type %q, got %q", tt.wantType, strategy.ReportType())
		}
		report := strategy.Generate(tasks)
		if report.Type != tt.wantType {
			t.Fatalf("expected report type %q, got %q", tt.wantType, report.Type)
		}
		if len(report.Items) == 0 {
			t.Fatalf("expected report items for %q", tt.reportType)
		}
	}
}

func TestSelectReportStrategyNegative(t *testing.T) {
	if _, err := SelectReportStrategy(models.ReportType("unknown")); err == nil {
		t.Fatal("expected error for unknown report strategy")
	}
}

func TestReportServiceBuildUsesSelectedStrategy(t *testing.T) {
	service := NewReportService(fakeReportTasks{
		tasks: []models.Task{
			{Status: models.StatusNew, Priority: models.PriorityHigh, AssigneeID: 1},
			{Status: models.StatusDone, Priority: models.PriorityHigh, AssigneeID: 2},
			{Status: models.StatusNew, Priority: models.PriorityLow, AssigneeID: 1},
		},
	}, &fakeLogger{})

	report, err := service.Build(context.Background(), models.ReportByPriority, models.TaskFilter{})
	if err != nil {
		t.Fatalf("expected report, got error: %v", err)
	}
	if report.Type != models.ReportByPriority {
		t.Fatalf("expected priority report, got %q", report.Type)
	}
	if len(report.Items) != 2 {
		t.Fatalf("expected two priority groups, got %d", len(report.Items))
	}
}
