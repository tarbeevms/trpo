package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"taskflow/internal/models"
)

type ReportStrategy interface {
	ReportType() models.ReportType
	Generate(tasks []models.Task) models.Report
}

type StatusReportStrategy struct{}

func (StatusReportStrategy) ReportType() models.ReportType {
	return models.ReportByStatus
}

func (StatusReportStrategy) Generate(tasks []models.Task) models.Report {
	counts := make(map[string]int)
	for _, task := range tasks {
		counts[string(task.Status)]++
	}
	return buildReport(models.ReportByStatus, counts)
}

type PriorityReportStrategy struct{}

func (PriorityReportStrategy) ReportType() models.ReportType {
	return models.ReportByPriority
}

func (PriorityReportStrategy) Generate(tasks []models.Task) models.Report {
	counts := make(map[string]int)
	for _, task := range tasks {
		counts[string(task.Priority)]++
	}
	return buildReport(models.ReportByPriority, counts)
}

type AssigneeReportStrategy struct{}

func (AssigneeReportStrategy) ReportType() models.ReportType {
	return models.ReportByAssignee
}

func (AssigneeReportStrategy) Generate(tasks []models.Task) models.Report {
	counts := make(map[string]int)
	for _, task := range tasks {
		counts[strconv.FormatInt(task.AssigneeID, 10)]++
	}
	return buildReport(models.ReportByAssignee, counts)
}

type ReportService struct {
	tasks      TaskStore
	logger     AppLogger
	strategies map[string]ReportStrategy
}

func NewReportService(tasks TaskStore, logger AppLogger) *ReportService {
	return &ReportService{
		tasks:      tasks,
		logger:     logger,
		strategies: DefaultReportStrategies(),
	}
}

func (s *ReportService) Build(ctx context.Context, reportType models.ReportType, filter models.TaskFilter) (models.Report, error) {
	strategy, err := selectReportStrategy(reportType, s.strategies)
	if err != nil {
		s.logger.Error("failed to select report strategy", "error", err)
		return models.Report{}, err
	}
	tasks, err := s.tasks.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to load report tasks", "error", err)
		return models.Report{}, err
	}
	report := strategy.Generate(tasks)
	s.logger.Info("report built", "type", report.Type)
	return report, nil
}

func DefaultReportStrategies() map[string]ReportStrategy {
	strategies := []ReportStrategy{
		StatusReportStrategy{},
		PriorityReportStrategy{},
		AssigneeReportStrategy{},
	}
	byType := make(map[string]ReportStrategy, len(strategies))
	for _, strategy := range strategies {
		byType[string(strategy.ReportType())] = strategy
	}
	return byType
}

func SelectReportStrategy(reportType models.ReportType) (ReportStrategy, error) {
	return selectReportStrategy(reportType, DefaultReportStrategies())
}

func selectReportStrategy(reportType models.ReportType, strategies map[string]ReportStrategy) (ReportStrategy, error) {
	reportType = reportType.Normalize()
	strategy, ok := strategies[string(reportType)]
	if !ok {
		return nil, fmt.Errorf("unknown report type")
	}
	return strategy, nil
}

func buildReport(reportType models.ReportType, counts map[string]int) models.Report {
	labels := make([]string, 0, len(counts))
	for label := range counts {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	items := make([]models.ReportItem, 0, len(labels))
	for _, label := range labels {
		items = append(items, models.ReportItem{Label: label, Count: counts[label]})
	}
	return models.Report{Type: reportType, Items: items}
}
