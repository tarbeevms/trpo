package service

import (
	"context"

	"taskflow/internal/models"
)

// ReportService - сервис для экспорта задач пользователя
type ReportService struct {
	tasks     TaskStore
	logger    AppLogger
	exporters *ExporterRegistry
}

func NewReportService(tasks TaskStore, logger AppLogger) *ReportService {
	return &ReportService{
		tasks:     tasks,
		logger:    logger,
		exporters: NewExporterRegistry(),
	}
}

// Build экспортирует все задачи пользователя в указанном формате
func (s *ReportService) Build(ctx context.Context, filter models.TaskFilter, format string) ([]byte, error) {
	tasks, err := s.tasks.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to load tasks for export", "error", err)
		return nil, err
	}

	// Экспортируем задачи в нужном формате
	exporter := s.exporters.Get(format)
	exported, err := exporter.ExportTasks(tasks)
	if err != nil {
		s.logger.Error("failed to export tasks", "error", err)
		return nil, err
	}

	s.logger.Info("tasks exported", "format", format, "count", len(tasks))
	return exported, nil
}

// GetAvailableFormats возвращает список доступных форматов экспорта
func (s *ReportService) GetAvailableFormats() []string {
	return s.exporters.Formats()
}
