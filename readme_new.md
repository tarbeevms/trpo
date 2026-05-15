# TaskFlow — система управления задачами

## 1. Постановка задачи

**TaskFlow** — это веб-приложение для управления задачами, разработанное в соответствии с принципами объектно-ориентированного программирования. Приложение позволяет пользователям создавать проекты, управлять задачами в рамках проектов, отслеживать статус и приоритет задач, а также экспортировать данные в различных форматах.

### Цель проекта

Разработать полнофункциональное приложение, демонстрирующее все принципы ООП (инкапсуляция, абстракция, полиморфизм, наследование), использование паттернов проектирования, работу с базой данных PostgreSQL, логирование и unit-тестирование.

---

## 2. Функции приложения

### 2.1. Аутентификация и авторизация
- **Регистрация пользователей**: создание учётной записи с уникальным логином
- **Вход в систему**: аутентификация пользователя с использованием JWT-токенов
- **Защищённые маршруты**: доступ к функциям только для авторизованных пользователей

### 2.2. Управление проектами
- **Создание проекта**: название (3-80 символов), описание (до 500 символов)
- **Просмотр списка проектов**: отображение всех проектов текущего пользователя
- **Проекты принадлежат пользователю**: каждый проект имеет владельца

### 2.3. Управление задачами
- **Создание задачи**: название, описание, приоритет, дедлайн, теги; статус при создании автоматически становится `new`, исполнитель — текущий пользователь
- **Просмотр задач**: таблица только с задачами текущего пользователя
- **Изменение статуса**: смена статуса на одно из допустимых значений (`new`, `in_progress`, `done`, `cancelled`)
- **Валидация данных**: проверка корректности вводимых данных

### 2.4. Экспорт отчётов
- **Экспорт в HTML**: табличное отображение задач текущего пользователя
- **Экспорт в JSON**: вывод задач текущего пользователя в формате JSON
- **Экспорт в XML**: вывод задач текущего пользователя в формате XML
- **Отображение на фронтенде**: HTML, JSON и XML отображаются на странице без скачивания

### 2.5. Теги
- **Добавление тегов**: до 10 тегов на задачу
- **Валидация тегов**: название тега 2-30 символов

---

## 3. Формат входных и выходных данных

### 3.1. Входные данные

#### Регистрация пользователя
```json
{
  "login": "user123",
  "password": "securepass"
}
```

#### Создание проекта
```json
{
  "name": "Мой проект",
  "description": "Описание проекта"
}
```

#### Создание задачи
```json
{
  "project_id": 1,
  "title": "Название задачи",
  "description": "Описание задачи",
  "priority": "high",
  "deadline": "2026-06-01",
  "tags": ["study", "urgent"]
}
```

`status` и `assignee_id` пользователь не вводит вручную: обработчик создает задачу со статусом `new`, а исполнителем назначает текущего авторизованного пользователя.

### 3.2. Выходные данные

#### Список задач (JSON)
```json
[
  {
    "ID": 1,
    "ProjectID": 1,
    "Title": "Название задачи",
    "Description": "Описание",
    "Status": "new",
    "Priority": "high",
    "Deadline": "2026-06-01T00:00:00Z",
    "Tags": [{"ID": 1, "Name": "study"}]
  }
]
```

#### Список задач (XML)
```xml
<?xml version="1.0" encoding="UTF-8"?>
<tasks>
  <task>
    <id>1</id>
    <project_id>1</project_id>
    <title>Название задачи</title>
    <description>Описание</description>
    <status>new</status>
    <priority>high</priority>
    <deadline>2026-06-01</deadline>
    <tags>
      <tag>study</tag>
    </tags>
  </task>
</tasks>
```

---

## 4. Ограничения на данные

| Поле | Тип | Ограничение |
|------|-----|-------------|
| Логин пользователя | string | 3-50 символов, без пробелов, уникальный |
| Пароль пользователя | string | 6-72 символа |
| Название проекта | string | 3-80 символов |
| Описание проекта | string | До 500 символов |
| `project_id` задачи | integer | Обязателен, должен ссылаться на проект текущего пользователя |
| Название задачи | string | 3-100 символов |
| Описание задачи | string | До 1000 символов |
| Статус задачи | enum | new, in_progress, done, cancelled |
| Приоритет задачи | enum | low, medium, high, critical |
| Дедлайн | date | Не ранее текущей даты |
| Исполнитель задачи | integer | Автоматически текущий пользователь |
| Теги | array | До 10 тегов на задачу |
| Название тега | string | 2-30 символов |
| Формат экспорта | enum | html, json, xml |

---

## 5. Структура проекта

```
taskflow/
├── cmd/
│   └── taskflow/
│       └── main.go              # Точка входа в приложение
├── frontend/
│   └── index.html               # Веб-интерфейс (HTML-шаблоны)
├── internal/
│   ├── app/
│   │   ├── project.go           # HTTP-обработчики для проектов
│   │   ├── report.go            # HTTP-обработчики для отчётов
│   │   ├── task.go              # HTTP-обработчики для задач
│   │   └── user.go              # HTTP-обработчики для пользователей
│   ├── models/
│   │   ├── common.go            # Базовые структуры и функции валидации
│   │   ├── project.go           # Модель проекта
│   │   ├── report.go            # Модели для отчётов
│   │   ├── task.go              # Модель задачи
│   │   ├── user.go              # Модель пользователя
│   │   └── *_test.go            # Unit-тесты для моделей
│   ├── repository/
│   │   ├── project.go           # Репозиторий проектов
│   │   ├── report.go            # Заготовка репозитория отчётов
│   │   ├── task.go              # Репозиторий задач
│   │   └── user.go              # Репозиторий пользователей
│   └── service/
│       ├── project.go           # Бизнес-логика проектов
│       ├── report.go            # Бизнес-логика отчётов
│       ├── report_exporter.go   # Экспортёры отчётов (Strategy)
│       ├── task.go              # Бизнес-логика задач (Facade)
│       ├── task_test.go         # Unit-тесты для задач
│       └── user.go              # Бизнес-логика пользователей
├── migrations/
│   └── 001_init.sql             # Миграции базы данных
├── build/
│   └── Dockerfile               # Docker-конфигурация
├── docker-compose.yml           # Запуск PostgreSQL
├── go.mod                       # Зависимости Go
└── PLAN.md                      # Планирование проекта
```

---

## 6. Принципы ООП

### 6.1. Инкапсуляция

**Инкапсуляция** — это принцип ООП, при котором данные и методы, работающие с этими данными, объединяются в единый объект, скрывая внутреннюю реализацию от внешнего мира. В Go инкапсуляция реализуется через экспортируемые (с заглавной буквы) и неэкспортируемые (со строчной буквы) идентификаторы.

#### Приватные поля и методы

В нашем проекте инкапсуляция реализуется следующим образом:

**В моделях (internal/models/):**
- Поля структур с маленькой буквы недоступны извне пакета
- Валидация выполняется через методы структур (например, `Task.Validate()`)

```go
// internal/models/task.go
type Task struct {
    BaseEntity
    AuditInfo
    SoftDelete
    ProjectID   int64     // экспортируемое поле
    Title       string    // экспортируемое поле
    Description string    // экспортируемое поле
    Status      Status    // экспортируемое поле
    Priority    Priority  // экспортируемое поле
    Deadline    time.Time // экспортируемое поле
    AssigneeID  int64     // экспортируемое поле
    Tags        []Tag     // экспортируемое поле
}

// Валидация — метод структуры, инкапсулирующий логику проверки
func (t Task) Validate(now time.Time) error {
    if t.ProjectID <= 0 {
        return fmt.Errorf("project_id is required")
    }
    if err := validateLength("title", t.Title, 3, 100); err != nil {
        return err
    }
    // ... другие проверки
    return nil
}
```

**В сервисном слое (internal/service/):**
- Фасад `TaskFacade` скрывает сложность взаимодействия с несколькими хранилищами
- Приватные поля структуры недоступны извне

```go
// internal/service/task.go
type TaskFacade struct {
    tasks    TaskStore   // приватное поле
    projects ProjectStore // приватное поле
    users    UserStore   // приватное поле
    logger   AppLogger   // приватное поле
    now      func() time.Time // приватное поле
}

// Приватный метод — недоступен извне пакета
func (f *TaskFacade) validateTask(ctx context.Context, task models.Task) error {
    // логика валидации
}
```

**В репозиториях (internal/repository/):**
- Структуры репозиториев имеют приватное поле `db` для работы с БД

```go
// internal/repository/task.go
type TaskRepository struct {
    db *sql.DB // приватное поле — подключение к БД недоступно извне
}
```

#### Зачем нужна инкапсуляция

1. **Защита данных**: предотвращение некорректного изменения состояния объекта
2. **Сокрытие реализации**: внешний код не зависит от внутренней структуры
3. **Упрощение поддержки**: изменение внутренней реализации не ломает внешний код
4. **Контроль доступа**: можно явно определить, какие данные доступны для чтения/записи

---

### 6.2. Абстракция

**Абстракция** — это принцип ООП, позволяющий выделить существенные характеристики объекта, скрывая несущественные детали. В Go абстракция реализуется через интерфейсы.

#### Интерфейсы в проекте

**TaskStore — абстракция хранилища задач:**
```go
// internal/service/task.go
type TaskStore interface {
    Create(ctx context.Context, task *models.Task) error
    List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
    FindByID(ctx context.Context, id int64) (models.Task, error)
    UpdateStatus(ctx context.Context, task models.Task, history models.TaskHistory) error
    Delete(ctx context.Context, id int64) error
}
```

`TaskFilter` в текущей реализации не является пользовательским фильтром по статусу или приоритету. Он используется как техническое ограничение доступа: `AssigneeID` задается из JWT текущего пользователя, чтобы список задач и экспорт не показывали чужие задачи.

**ProjectStore — абстракция хранилища проектов:**
```go
// internal/service/project.go
type ProjectStore interface {
    Create(ctx context.Context, project *models.Project) error
    List(ctx context.Context) ([]models.Project, error)
    ListByOwner(ctx context.Context, ownerID int64) ([]models.Project, error)
    Exists(ctx context.Context, id int64) (bool, error)
    OwnedBy(ctx context.Context, projectID int64, ownerID int64) (bool, error)
    Delete(ctx context.Context, id int64) error
}
```

**UserStore — абстракция хранилища пользователей:**
```go
// internal/service/user.go
type UserStore interface {
    Create(ctx context.Context, user *models.User) error
    List(ctx context.Context) ([]models.User, error)
    Exists(ctx context.Context, id int64) (bool, error)
    LoginExists(ctx context.Context, login string) (bool, error)
    FindByID(ctx context.Context, id int64) (models.User, error)
    FindByLogin(ctx context.Context, login string) (models.User, error)
}
```

**ReportExporter — абстракция экспортёра отчётов:**
```go
// internal/service/report_exporter.go
type ReportExporter interface {
    Format() string
    ExportTasks(tasks []models.Task) ([]byte, error)
}
```

#### Преимущества абстракции

1. **Взаимозаменяемость**: можно подменить реализацию, не меняя код клиента
2. **Независимость от деталей**: клиент работает с абстракцией, а не с конкретной реализацией
3. **Упрощение тестирования**: можно подставить мок-реализацию интерфейса
4. **Гибкость архитектуры**: легко добавить новые реализации

---

### 6.3. Полиморфизм

**Полиморфизм** — это принцип ООП, позволяющий объектам разных типов обрабатываться через единый интерфейс. В Go полиморфизм реализуется через интерфейсы.

#### Реализация полиморфизма в проекте

##### 1. Паттерн Strategy для экспорта отчётов

Интерфейс `ReportExporter` определяет единый контракт для всех экспортёров:

```go
// internal/service/report_exporter.go
type ReportExporter interface {
    Format() string
    ExportTasks(tasks []models.Task) ([]byte, error)
}
```

Три реализации этого интерфейса:

**JSONExporter:**
```go
type JSONExporter struct{}

func (e JSONExporter) Format() string {
    return "json"
}

func (e JSONExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
    return json.MarshalIndent(tasks, "", "  ")
}
```

**XMLExporter:**
```go
type XMLExporter struct{}

func (e XMLExporter) Format() string {
    return "xml"
}

func (e XMLExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
    // сериализация в XML
}
```

**HTMLExporter:**
```go
type HTMLExporter struct{}

func (e HTMLExporter) Format() string {
    return "html"
}

func (e HTMLExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
    // генерация HTML-таблицы
}
```

Клиентский код работает с интерфейсом, не зная о конкретной реализации:

```go
// internal/service/report.go
type ReportService struct {
    tasks     TaskStore
    logger    AppLogger
    exporters *ExporterRegistry
}

func (s *ReportService) Build(ctx context.Context, filter models.TaskFilter, format string) ([]byte, error) {
    tasks, err := s.tasks.List(ctx, filter)
    if err != nil {
        return nil, err
    }
    exporter := s.exporters.Get(format)
    return exporter.ExportTasks(tasks)
}
```

HTTP-обработчик `/report` передает в `Build` фильтр `models.TaskFilter{AssigneeID: user.ID}`. Поэтому экспорт строится только по задачам текущего пользователя.

##### 2. Полиморфизм через интерфейсы хранилищ

Один и тот же код сервиса может работать с разными реализациями хранилищ:

```go
// Создание фасада с конкретными реализациями
facade := service.NewTaskFacade(
    repository.NewTaskRepository(db),    // реализация TaskStore
    repository.NewProjectRepository(db), // реализация ProjectStore
    repository.NewUserRepository(db),    // реализация UserStore
    logger,
)

// В тестах можно подставить мок-реализации
mockTaskStore := &MockTaskStore{}
facade := service.NewTaskFacade(mockTaskStore, mockProjectStore, mockUserStore, logger)
```

#### Преимущества полиморфизма

1. **Расширяемость**: добавление новых форматов экспорта не требует изменения существующего кода
2. **Унификация**: единый интерфейс для работы с разными реализациями
3. **Тестируемость**: легко подменить реализацию на мок

---

### 6.4. Наследование

**Наследование** — это принцип ООП, позволяющий создавать новые классы на основе существующих. В Go наследование реализуется через **встраивание структур** (struct embedding).

#### Одиночное наследование

В Go одиночное наследование реализуется через встраивание одной структуры в другую:

```go
// internal/models/common.go
type BaseEntity struct {
    ID        int64
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (e BaseEntity) GetID() int64 {
    return e.ID
}

func (e BaseEntity) GetCreatedAt() time.Time {
    return e.CreatedAt
}
```

Структура `Task` наследует от `BaseEntity`:

```go
// internal/models/task.go
type Task struct {
    BaseEntity           // встраивание — наследование
    ProjectID   int64
    Title       string
    // ...
}
```

Теперь `Task` имеет доступ к полям и методам `BaseEntity`:

```go
task := Task{BaseEntity: BaseEntity{ID: 1}, Title: "Test"}
id := task.BaseEntity.GetID() // явный вызов метода встроенной структуры
createdAt := task.CreatedAt  // доступ к унаследованному полю
```

#### Множественное наследование

Go поддерживает множественное наследование через встраивание нескольких структур:

```go
// internal/models/common.go

// BaseEntity — базовая сущность с ID и временными метками
type BaseEntity struct {
    ID        int64
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (e BaseEntity) GetID() int64 {
    return e.ID
}

func (e BaseEntity) GetCreatedAt() time.Time {
    return e.CreatedAt
}

// AuditInfo — информация о создателе записи
type AuditInfo struct {
    CreatedBy int64
}

// GetID возвращает ID создателя записи
func (a AuditInfo) GetID() int64 {
    return a.CreatedBy
}

// GetCreatedBy возвращает ID создателя (уникальный метод)
func (a AuditInfo) GetCreatedBy() int64 {
    return a.CreatedBy
}

// SoftDelete — мягкое удаление
type SoftDelete struct {
    DeletedAt *time.Time
}

// GetID возвращает статус удаления (0 — не удалён, -1 — удалён)
func (s SoftDelete) GetID() int64 {
    if s.DeletedAt != nil {
        return -1
    }
    return 0
}

// IsDeleted проверяет, удалена ли запись (уникальный метод)
func (s SoftDelete) IsDeleted() bool {
    return s.DeletedAt != nil
}
```

Структура `Task` наследует от трёх структур:

```go
// internal/models/task.go
type Task struct {
    BaseEntity   // наследование #1
    AuditInfo    // наследование #2
    SoftDelete   // наследование #3
    ProjectID    int64
    Title        string
    Description  string
    Status       Status
    Priority     Priority
    Deadline     time.Time
    AssigneeID   int64
    Tags         []Tag
}
```

#### Разрешение конфликтов методов

При множественном наследовании может возникнуть конфликт, когда встроенные структуры имеют методы с одинаковыми названиями. В нашем случае все три структуры имеют метод `GetID()`:

- `BaseEntity.GetID()` — возвращает ID сущности
- `AuditInfo.GetID()` — возвращает ID создателя
- `SoftDelete.GetID()` — возвращает статус удаления

**Как Go разрешает конфликты:**

1. Если несколько встроенных структур имеют метод с одинаковым именем на одном уровне, вызов `task.GetID()` будет неоднозначным и не скомпилируется.
2. Для доступа к конфликтующим методам нужно явно указать встроенную структуру: `task.BaseEntity.GetID()`, `task.AuditInfo.GetID()`, `task.SoftDelete.GetID()`.

Пример использования в методе `GetTaskInfo()`:

```go
// internal/models/task.go
func (t Task) GetTaskInfo() string {
    // Уникальные методы — работают напрямую без конфликтов
    creatorID := t.GetCreatedBy()    // метод из AuditInfo
    isDeleted := t.IsDeleted()       // метод из SoftDelete
    createdAt := t.GetCreatedAt()    // метод из BaseEntity

    // Методы с одинаковым названием — требуют явного указания структуры
    entityID := t.BaseEntity.GetID()          // ID задачи из BaseEntity
    creatorIDFromGetID := t.AuditInfo.GetID() // ID создателя из AuditInfo
    deleteID := t.SoftDelete.GetID()          // статус удаления из SoftDelete

    return fmt.Sprintf("Task[id=%d, creator=%d, deleted=%v, created=%v, entityID=%d, creatorID=%d, deleteStatus=%d]",
        t.ID, creatorID, isDeleted, createdAt, entityID, creatorIDFromGetID, deleteID)
}
```

#### Преимущества наследования в Go

1. **Повторное использование кода**: общая логика выносится в базовые структуры
2. **Композиция вместо наследования**: Go поощряет композицию, но поддерживает наследование через встраивание
3. **Гибкость**: можно наследовать от нескольких структур
4. **Явное разрешение конфликтов**: разработчик контролирует, какой метод вызывать

---

## 7. Паттерны проектирования

### 7.1. Facade (Фасад)

**Назначение:** предоставить унифицированный интерфейс к набору интерфейсов в подсистеме, упрощая работу с ней.

**Суть паттерна:** скрыть сложность взаимодействия нескольких компонентов за простым интерфейсом.

#### Реализация в проекте

`TaskFacade` — фасад для работы с задачами, который скрывает взаимодействие с тремя хранилищами:

```go
// internal/service/task.go
type TaskFacade struct {
    tasks    TaskStore
    projects ProjectStore
    users    UserStore
    logger   AppLogger
    now      func() time.Time
}

func NewTaskFacade(tasks TaskStore, projects ProjectStore, users UserStore, logger AppLogger) *TaskFacade {
    return &TaskFacade{
        tasks:    tasks,
        projects: projects,
        users:    users,
        logger:   logger,
        now:      time.Now,
    }
}
```

**Что скрывает фасад:**

1. **Валидация задачи:** проверка проекта, владельца, исполнителя
2. **Проверка существования:** верификация, что проект и пользователь существуют
3. **Проверка прав:** задача должна принадлежать проекту пользователя
4. **Логирование:** все операции логируются

```go
func (f *TaskFacade) Create(ctx context.Context, task models.Task) (models.Task, error) {
    // 1. Валидация задачи
    if err := task.Validate(f.now()); err != nil {
        f.logger.Error("task validation failed", "error", err)
        return models.Task{}, err
    }

    // 2. Проверка существования проекта
    projectExists, err := f.projects.Exists(ctx, task.ProjectID)
    if err != nil {
        f.logger.Error("failed to check task project", "error", err)
        return models.Task{}, err
    }
    if !projectExists {
        return models.Task{}, fmt.Errorf("project_id must reference an existing project")
    }

    // 3. Проверка прав на проект
    if task.CreatedBy > 0 {
        projectOwned, err := f.projects.OwnedBy(ctx, task.ProjectID, task.CreatedBy)
        if err != nil {
            return models.Task{}, err
        }
        if !projectOwned {
            return models.Task{}, fmt.Errorf("project_id must reference current user's project")
        }
    }

    // 4. Проверка существования исполнителя
    assigneeExists, err := f.users.Exists(ctx, task.AssigneeID)
    if err != nil {
        return models.Task{}, err
    }
    if !assigneeExists {
        return models.Task{}, fmt.Errorf("assignee_id must reference an existing user")
    }

    // 5. Создание задачи
    if err := f.tasks.Create(ctx, &task); err != nil {
        f.logger.Error("failed to create task", "error", err)
        return models.Task{}, err
    }

    f.logger.Info("task created", "task_id", task.ID)
    return task, nil
}
```

**Преимущества использования Facade:**

1. **Упрощение клиентского кода:** клиент работает с одним методом вместо нескольких
2. **Слабая связанность:** клиент не зависит от внутренней структуры
3. **Изоляция сложности:** вся логика проверок скрыта в фасаде
4. **Единая точка входа:** все операции проходят через фасад

---

### 7.2. Strategy (Стратегия)

**Назначение:** определить семейство алгоритмов, инкапсулировать каждый из них и сделать их взаимозаменяемыми.

**Суть паттерна:** вынести поведение в отдельные классы, которые можно менять в runtime.

#### Реализация в проекте

Интерфейс `ReportExporter` определяет стратегию экспорта:

```go
// internal/service/report_exporter.go
type ReportExporter interface {
    Format() string
    ExportTasks(tasks []models.Task) ([]byte, error)
}
```

Три реализации стратегии:

```go
// JSON — экспорт в JSON
type JSONExporter struct{}
func (e JSONExporter) Format() string { return "json" }
func (e JSONExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
    return json.MarshalIndent(tasks, "", "  ")
}

// XML — экспорт в XML
type XMLExporter struct{}
func (e XMLExporter) Format() string { return "xml" }
func (e XMLExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
    // сериализация в XML
}

// HTML — экспорт для отображения на фронтенде
type HTMLExporter struct{}
func (e HTMLExporter) Format() string { return "html" }
func (e HTMLExporter) ExportTasks(tasks []models.Task) ([]byte, error) {
    // генерация HTML
}
```

**Реестр стратегий** — позволяет выбрать нужную стратегию:

```go
type ExporterRegistry struct {
    exporters map[string]ReportExporter
}

func NewExporterRegistry() *ExporterRegistry {
    return &ExporterRegistry{
        exporters: map[string]ReportExporter{
            "json": JSONExporter{},
            "xml":  XMLExporter{},
            "html": HTMLExporter{},
        },
    }
}

func (r *ExporterRegistry) Get(format string) ReportExporter {
    if exp, ok := r.exporters[format]; ok {
        return exp
    }
    return HTMLExporter{} // стратегия по умолчанию
}
```

**Использование в сервисе:**

```go
func (s *ReportService) Build(ctx context.Context, filter models.TaskFilter, format string) ([]byte, error) {
    tasks, err := s.tasks.List(ctx, filter)
    if err != nil {
        return nil, err
    }
    exporter := s.exporters.Get(format) // выбор стратегии
    return exporter.ExportTasks(tasks)
}
```

В обработчике отчета фильтр формируется из текущего пользователя:

```go
filter := models.TaskFilter{AssigneeID: user.ID}
exported, err := s.reports.Build(r.Context(), filter, format)
```

Это важно для безопасности: пользователь получает экспорт только своих задач.

**Преимущества использования Strategy:**

1. **Взаимозаменяемость:** можно менять формат экспорта без изменения клиентского кода
2. **Расширяемость:** добавление нового формата не требует изменения существующего кода
3. **Изоляция алгоритмов:** каждый экспортёр инкапсулирует свою логику
4. **Тестируемость:** можно тестировать каждый экспортёр отдельно

---

### 7.3. Сочетание паттернов

**Facade + Strategy** — эти паттерны отлично сочетаются:

1. **Facade** (`TaskFacade`) предоставляет простой интерфейс к сложной подсистеме управления задачами
2. **Strategy** (`ReportExporter`) позволяет гибко выбирать алгоритм экспорта

Клиентский код работает через Facade, который внутри может использовать различные Strategy:

```go
// HTTP-обработчик использует Facade
handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    format := r.URL.Query().Get("format")
    filter := models.TaskFilter{AssigneeID: userID}
    data, err := reportService.Build(r.Context(), filter, format)
    // ...
})
```

---

## 8. Связи между объектами

### 8.1. Ассоциация (Association)

**Определение:** отношение, при котором объекты одного класса связаны с объектами другого класса. Это самая общая форма связи.

**В коде:** реализуется через поля структур, ссылающиеся на другие структуры.

#### Примеры в проекте

**Task → User (Assignee):**
```go
// Задача ссылается на исполнителя через ID
type Task struct {
    AssigneeID int64  // ссылка на пользователя
    // ...
}
```

**Project → User (Owner):**
```go
type Project struct {
    OwnerID int64  // ссылка на владельца
    // ...
}
```

**В БД:**
```sql
-- Внешний ключ связывает задачу с пользователем
ALTER TABLE tasks ADD CONSTRAINT fk_task_assignee
    FOREIGN KEY (assignee_id) REFERENCES users(id);
```

---

### 8.2. Агрегация (Aggregation)

**Определение:** отношение "часть-целое", при котором часть может существовать независимо от целого.

**В коде:** реализуется через встраивание структур или поля со ссылками.

#### Примеры в проекте

**Task → Tags:**
```go
type Task struct {
    Tags []Tag  // теги принадлежат задаче, но могут существовать независимо
}
```

**Project → Tasks:**
```go
type Project struct {
    Tasks []Task  // задачи принадлежат проекту
}
```

**В БД:**
```sql
-- Связь многие-ко-многим через промежуточную таблицу
CREATE TABLE task_tags (
    task_id BIGINT REFERENCES tasks(id),
    tag_id BIGINT REFERENCES tags(id),
    PRIMARY KEY (task_id, tag_id)
);
```

---

### 8.3. Композиция (Composition)

**Определение:** отношение "часть-целое", при котором часть не может существовать без целого.

**В коде:** реализуется через встраивание структур (embedded structs).

#### Примеры в проекте

**Task встраивает BaseEntity, AuditInfo, SoftDelete:**
```go
type Task struct {
    BaseEntity   // Task "содержит" BaseEntity
    AuditInfo    // Task "содержит" AuditInfo
    SoftDelete   // Task "содержит" SoftDelete
    // ...
}
```

Это композиция, потому что:
- `BaseEntity` не имеет смысла без сущности (Task, Project, User)
- `AuditInfo` и `SoftDelete` — это миксины, добавляющие функциональность

**В БД:**
```sql
-- Поля из встроенных структур хранятся в той же таблице
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    -- остальные поля задачи
);
```

---

### 8.4. Зависимость (Dependency)

**Определение:** отношение, при котором изменение одного класса влияет на другой.

**В коде:** реализуется через параметры методов, интерфейсы.

#### Примеры в проекте

**TaskFacade зависит от интерфейсов:**
```go
type TaskFacade struct {
    tasks    TaskStore     // зависимость от интерфейса
    projects ProjectStore  // зависимость от интерфейса
    users    UserStore     // зависимость от интерфейса
}
```

**Методы принимают интерфейсы:**
```go
func (f *TaskFacade) Create(ctx context.Context, task models.Task) (models.Task, error) {
    // метод зависит от контекста
}
```

---

### 8.5. Реализация (Implementation)

**Определение:** отношение между интерфейсом и его реализацией.

**В коде:** структура реализует интерфейс.

#### Примеры в проекте

**TaskRepository реализует TaskStore:**
```go
type TaskRepository struct {
    db *sql.DB
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
    // реализация
}

func (r *TaskRepository) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
    // реализация
}
```

`TaskRepository` реализует `TaskStore` неявно: в Go структура считается реализацией интерфейса, если у нее есть все методы интерфейса.

---

## 9. Тестирование

### 9.1. Виды тестов

В проекте реализованы следующие виды unit-тестов:

#### Позитивные тесты (Positive Tests)

Проверка корректной работы при валидных данных:

```go
// internal/models/task_test.go
func TestTaskValidatePositive(t *testing.T) {
    now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
    task := validTask(now)
    if err := task.Validate(now); err != nil {
        t.Fatalf("expected valid task, got error: %v", err)
    }
}
```

#### Негативные тесты (Negative Tests)

Проверка обработки ошибок при некорректных данных:

```go
// internal/models/task_test.go
func TestTaskValidateNegativeStatus(t *testing.T) {
    now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
    task := validTask(now)
    task.Status = "bad" // некорректный статус
    if err := task.Validate(now); err == nil {
        t.Fatal("expected status validation error")
    }
}

func TestTaskValidateNegativeDeadline(t *testing.T) {
    now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
    task := validTask(now)
    task.Deadline = now.AddDate(0, 0, -1) // дедлайн в прошлом
    if err := task.Validate(now); err == nil {
        t.Fatal("expected deadline validation error")
    }
}
```

#### Граничные тесты (Boundary Tests)

Проверка значений на границах допустимых диапазонов:

```go
// internal/models/task_test.go
func TestTaskValidateBoundaryFields(t *testing.T) {
    now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
    task := validTask(now)
    
    // Минимальная длина названия (3 символа)
    task.Title = strings.Repeat("a", 3)
    // Максимальная длина описания (1000 символов)
    task.Description = strings.Repeat("b", 1000)
    // Максимальное количество тегов (10)
    task.Tags = make([]Tag, 10)
    
    if err := task.Validate(now); err != nil {
        t.Fatalf("expected boundary values to be valid: %v", err)
    }
    
    // Максимальная длина названия (100 символов)
    task.Title = strings.Repeat("a", 100)
    if err := task.Validate(now); err != nil {
        t.Fatalf("expected maximum title to be valid: %v", err)
    }
}
```

### 9.2. Тестовые файлы

| Файл | Описание |
|------|----------|
| `internal/models/task_test.go` | Тесты валидации задач, изменения статуса |
| `internal/models/project_test.go` | Тесты валидации проектов |
| `internal/models/user_test.go` | Тесты валидации пользователей |
| `internal/models/report_test.go` | Тесты нормализации типа отчёта |
| `internal/service/task_test.go` | Тесты бизнес-логики задач и фасада |
| `internal/service/report_test.go` | Тесты экспортёров, выбора стратегии экспорта и передачи `AssigneeID` в `TaskFilter` |

### 9.3. Запуск тестов

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск тестов с детальным выводом
go test -v ./internal/models/...
```

---

## 10. База данных

### 10.1. Схема БД

```sql
-- Пользователи
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL
);

-- Проекты
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(80) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    owner_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL
);

-- Задачи
CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id),
    title VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    deadline DATE NOT NULL,
    assignee_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL
);

-- Теги
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Связь задач и тегов (многие-ко-многим)
CREATE TABLE task_tags (
    task_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);

-- История изменения статусов
CREATE TABLE task_history (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE,
    old_status VARCHAR(20) NOT NULL,
    new_status VARCHAR(20) NOT NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 10.2. Запуск PostgreSQL

```bash
# Запуск через docker-compose
docker-compose up -d postgres

# Миграция копируется Dockerfile в /docker-entrypoint-initdb.d
# и применяется автоматически при первом создании volume.
```

---

## 11. Веб-интерфейс

### 11.1. Структура HTML-страницы

Фронтенд реализован как единый HTML-файл с серверным рендерингом через Go-шаблоны.

**Секции страницы:**

1. **Регистрация** — форма создания нового пользователя
2. **Вход** — форма аутентификации
3. **Проекты** — создание и просмотр проектов
4. **Задачи** — создание и просмотр задач текущего пользователя
5. **Статусы** — изменение статуса задачи
6. **Отчёты** — экспорт задач в JSON/XML/HTML

### 11.2. Функциональность фронтенда

- Асинхронные запросы через fetch API
- Отображение HTML, JSON и XML без скачивания
- Валидация форм на стороне клиента
- Адаптивный дизайн

---

## 12. Логирование

### 12.1. Реализация

Логирование реализовано через интерфейс `AppLogger`:

```go
// internal/service/user.go
type AppLogger interface {
    Info(msg string, args ...any)
    Error(msg string, args ...any)
}
```

### 12.2. Использование

```go
func (f *TaskFacade) Create(ctx context.Context, task models.Task) (models.Task, error) {
    f.logger.Info("creating task", "title", task.Title)
    
    if err := task.Validate(f.now()); err != nil {
        f.logger.Error("task validation failed", "error", err)
        return models.Task{}, err
    }
    // ...
}
```

---

## 13. Итог реализации

### 13.1. Выполненные требования

| Требование | Статус |
|------------|--------|
| Не менее 8 классов | ✅ 8+ моделей |
| Инкапсуляция | ✅ Приватные поля, методы |
| Абстракция | ✅ Интерфейсы TaskStore, ProjectStore, UserStore, ReportExporter |
| Полиморфизм | ✅ Реализация через интерфейсы |
| Наследование одиночное | ✅ BaseEntity → Task, Project, User |
| Наследование множественное | ✅ User, Project и Task встраивают несколько структур |
| Конфликт методов | ✅ GetID() в нескольких структурах с разрешением |
| Связи объектов | ✅ Ассоциация, агрегация, композиция, зависимость |
| Паттерн Facade | ✅ TaskFacade |
| Паттерн Strategy | ✅ ReportExporter (JSON, XML, HTML) |
| Логирование | ✅ AppLogger |
| Загрузка из БД | ✅ PostgreSQL |
| Unit-тестирование | ✅ Позитивные, негативные, граничные тесты |

### 13.2. Технологический стек

- **Язык:** Go 1.21+
- **База данных:** PostgreSQL
- **Фронтенд:** HTML/CSS/JavaScript
- **Аутентификация:** JWT
- **Тестирование:** Go testing package
