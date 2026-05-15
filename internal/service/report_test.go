package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"taskflow/internal/models"
)

type fakeReportTasks struct {
	tasks      []models.Task
	lastFilter models.TaskFilter
}

func (f *fakeReportTasks) Create(ctx context.Context, task *models.Task) error {
	return nil
}

func (f *fakeReportTasks) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	f.lastFilter = filter
	return f.tasks, nil
}

func (f *fakeReportTasks) FindByID(ctx context.Context, id int64) (models.Task, error) {
	return models.Task{}, nil
}

func (f *fakeReportTasks) UpdateStatus(ctx context.Context, task models.Task, history models.TaskHistory) error {
	return nil
}

func (f *fakeReportTasks) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestExporterRegistrySelectsExporter(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{format: "html", want: "html"},
		{format: "json", want: "json"},
		{format: "xml", want: "xml"},
		{format: "unknown", want: "html"},
	}

	registry := NewExporterRegistry()
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			exporter := registry.Get(tt.format)
			if exporter.Format() != tt.want {
				t.Fatalf("expected exporter %q, got %q", tt.want, exporter.Format())
			}
		})
	}
}

func TestReportExportersPositive(t *testing.T) {
	tasks := []models.Task{{
		BaseEntity:  models.BaseEntity{ID: 10, CreatedAt: time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)},
		ProjectID:   1,
		Title:       "Write report",
		Description: "Prepare project description",
		Status:      models.StatusNew,
		Priority:    models.PriorityHigh,
		Deadline:    time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
		AssigneeID:  7,
		Tags:        []models.Tag{{Name: "study"}},
	}}

	tests := []struct {
		name     string
		exporter ReportExporter
		wantPart string
	}{
		{name: "html", exporter: HTMLExporter{}, wantPart: "<table>"},
		{name: "json", exporter: JSONExporter{}, wantPart: `"Title": "Write report"`},
		{name: "xml", exporter: XMLExporter{}, wantPart: "<title>Write report</title>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.exporter.ExportTasks(tasks)
			if err != nil {
				t.Fatalf("expected export, got error: %v", err)
			}
			if !strings.Contains(string(data), tt.wantPart) {
				t.Fatalf("expected exported data to contain %q, got %s", tt.wantPart, string(data))
			}
		})
	}
}

func TestReportServiceBuildUsesUserFilterAndExporter(t *testing.T) {
	tasks := &fakeReportTasks{
		tasks: []models.Task{{
			BaseEntity: models.BaseEntity{ID: 10},
			Title:      "Write report",
			Status:     models.StatusNew,
			Priority:   models.PriorityHigh,
			Deadline:   time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
			AssigneeID: 7,
		}},
	}
	service := NewReportService(tasks, &fakeLogger{})

	data, err := service.Build(context.Background(), models.TaskFilter{AssigneeID: 7}, "json")
	if err != nil {
		t.Fatalf("expected report export, got error: %v", err)
	}
	if tasks.lastFilter.AssigneeID != 7 {
		t.Fatalf("expected assignee filter 7, got %d", tasks.lastFilter.AssigneeID)
	}
	if !strings.Contains(string(data), `"AssigneeID": 7`) {
		t.Fatalf("expected exported data for current user, got %s", string(data))
	}
}
