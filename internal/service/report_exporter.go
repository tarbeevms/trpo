package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"

	"taskflow/internal/models"
)

// ReportExporter - интерфейс для экспорта задач в разных форматах
// Паттерн "Strategy" - каждая реализация представляет собой стратегию экспорта
type ReportExporter interface {
	// Format возвращает идентификатор формата (json, xml, html)
	Format() string
	// ExportTasks сериализует задачи в нужный формат
	ExportTasks(tasks []models.Task) ([]byte, error)
}

// JSONExporter - экспорт задач в формате JSON
type JSONExporter struct{}

func (e JSONExporter) Format() string {
	return "json"
}

func (e JSONExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
	return json.MarshalIndent(tasks, "", "  ")
}

// XMLExporter - экспорт задач в формате XML
type XMLExporter struct{}

func (e XMLExporter) Format() string {
	return "xml"
}

func (e XMLExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
	type XMLTask struct {
		ID          int64    `xml:"id"`
		ProjectID   int64    `xml:"project_id"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Status      string   `xml:"status"`
		Priority    string   `xml:"priority"`
		Deadline    string   `xml:"deadline"`
		AssigneeID  int64    `xml:"assignee_id"`
		CreatedAt   string   `xml:"created_at"`
		UpdatedAt   string   `xml:"updated_at"`
		Tags        []string `xml:"tags>tag"`
	}

	xmlTasks := make([]XMLTask, len(tasks))
	for i, task := range tasks {
		tags := make([]string, len(task.Tags))
		for j, tag := range task.Tags {
			tags[j] = tag.Name
		}
		xmlTasks[i] = XMLTask{
			ID:          task.ID,
			ProjectID:   task.ProjectID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Priority:    string(task.Priority),
			Deadline:    task.Deadline.Format("2006-01-02"),
			AssigneeID:  task.AssigneeID,
			CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Tags:        tags,
		}
	}

	type XMLTasks struct {
		XMLName string    `xml:"tasks"`
		Tasks   []XMLTask `xml:"task"`
	}

	return xml.MarshalIndent(XMLTasks{Tasks: xmlTasks}, "", "  ")
}

// HTMLExporter - экспорт задач в формате HTML (для отображения на фронте)
type HTMLExporter struct{}

func (e HTMLExporter) Format() string {
	return "html"
}

func (e HTMLExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
	if len(tasks) == 0 {
		return []byte("<p>Нет задач для отображения</p>"), nil
	}

	table := "<table><thead><tr><th>ID</th><th>Название</th><th>Описание</th><th>Статус</th><th>Приоритет</th><th>Дедлайн</th><th>Теги</th></tr></thead><tbody>"
	for _, task := range tasks {
		tags := ""
		for _, tag := range task.Tags {
			tags += html.EscapeString(tag.Name) + " "
		}
		table += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			task.ID,
			html.EscapeString(task.Title),
			html.EscapeString(task.Description),
			html.EscapeString(string(task.Status)),
			html.EscapeString(string(task.Priority)),
			task.Deadline.Format("2006-01-02"),
			tags,
		)
	}
	table += "</tbody></table>"

	return []byte(table), nil
}

// ExporterRegistry - реестр доступных экспортёров
// Позволяет выбрать нужный экспортёр по формату
type ExporterRegistry struct {
	exporters map[string]ReportExporter
}

// NewExporterRegistry создаёт реестр с доступными экспортёрами
func NewExporterRegistry() *ExporterRegistry {
	return &ExporterRegistry{
		exporters: map[string]ReportExporter{
			"json": JSONExporter{},
			"xml":  XMLExporter{},
			"html": HTMLExporter{},
		},
	}
}

// Get возвращает экспортёр по формату. Если формат не найден - возвращает HTMLExporter по умолчанию
func (r *ExporterRegistry) Get(format string) ReportExporter {
	if exp, ok := r.exporters[format]; ok {
		return exp
	}
	// По умолчанию возвращаем HTML
	return HTMLExporter{}
}

// Formats возвращает список доступных форматов
func (r *ExporterRegistry) Formats() []string {
	formats := make([]string, 0, len(r.exporters))
	for f := range r.exporters {
		formats = append(formats, f)
	}
	return formats
}
