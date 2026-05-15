# TaskFlow Documentation

## 1. Постановка задачи

`TaskFlow` - простое веб-приложение на Go для управления личными задачами и проектами.

Приложение позволяет регистрироваться, входить через JWT, создавать проекты и задачи на персональной странице пользователя, менять статус задачи, просматривать список своих задач, фильтровать задачи и строить простые отчеты. Данные сохраняются в PostgreSQL, который поднимается через Docker.

Проект реализован строго под требования КМ-4 и КМ-5: есть ООП-модель, не менее 8 структур, связи между объектами, загрузка данных из БД, 2 паттерна проектирования, логирование, unit-тесты и описание реализации.

## 2. Функции приложения

- регистрация пользователя;
- вход пользователя через JWT;
- выход пользователя;
- персональная страница пользователя `/me`;
- создание проекта;
- создание задачи;
- автоматическое назначение задачи текущему пользователю;
- изменение статуса задачи;
- просмотр списка задач;
- фильтрация задач по статусу и приоритету;
- построение отчета по статусам, приоритетам или исполнителям;
- логирование основных операций;
- сохранение и загрузка данных из PostgreSQL.

## 3. Формат входных и выходных данных

Пользователь работает с приложением через простой HTML-интерфейс [`frontend/index.html`](frontend/index.html).

Входные данные передаются через HTML-формы:

- регистрация: `login`, `password`;
- вход: `login`, `password`;
- проект: `name`, `description`;
- задача: `project_id`, `title`, `description`, `priority`, `deadline`, `tags`;
- смена статуса: `task_id`, `status`;
- отчет: `report_type`.

`report_type` определяет, какой алгоритм построения отчета будет выбран:

- `status` - отчет считает задачи по статусам: `new`, `in_progress`, `done`, `cancelled`;
- `priority` - отчет считает задачи по приоритетам: `low`, `medium`, `high`, `critical`;
- `assignee` - отчет считает задачи по ID исполнителя.

На уровне кода эти значения описаны отдельным доменным типом [`ReportType`](internal/models/report.go#L3-L9), а не случайными строками.

Выходные данные отображаются на HTML-странице:

- персональная страница текущего пользователя;
- список проектов текущего пользователя;
- список задач текущего пользователя;
- сообщение об успешной операции;
- сообщение об ошибке;
- отчет по задачам.

## 4. Ограничения на данные

`User`:

- `login` - обязательный, уникальный, от 3 до 50 символов, без пробелов, проверяется в [`User.Validate`](internal/models/user.go#L16-L24) и [`UserService.Create`](internal/service/user.go#L38-L59);
- `password` - от 6 до 72 символов, проверяется в [`validatePassword`](internal/service/user.go#L99-L104);
- `password_hash` - хранится в БД вместо открытого пароля, формируется в [`hashPassword`](internal/service/user.go#L106-L113).

`Project`:

- `name` - от 3 до 80 символов, проверяется в [`Project.Validate`](internal/models/project.go#L15-L26);
- `description` - до 500 символов, проверяется в [`Project.Validate`](internal/models/project.go#L15-L26);
- `owner_id` - обязательный, должен ссылаться на существующего пользователя, проверяется в [`ProjectService.Create`](internal/service/project.go#L26-L47).

`Task`:

- `title` - от 3 до 100 символов, проверяется в [`Task.Validate`](internal/models/task.go#L59-L90);
- `description` - до 1000 символов, проверяется в [`Task.Validate`](internal/models/task.go#L59-L90);
- `status` - только `new`, `in_progress`, `done`, `cancelled`, проверяется в [`Status.IsValid`](internal/models/task.go#L107-L109);
- `priority` - только `low`, `medium`, `high`, `critical`, проверяется в [`Priority.IsValid`](internal/models/task.go#L111-L113);
- `deadline` - не раньше текущей даты, проверяется в [`Task.Validate`](internal/models/task.go#L75-L77);
- `project_id` - обязательный, должен ссылаться на существующий проект текущего пользователя, проверяется в [`TaskFacade.Create`](internal/service/task.go#L40-L86);
- `assignee_id` - обязательный, должен ссылаться на существующего пользователя, проверяется в [`TaskFacade.Create`](internal/service/task.go#L40-L86);
- `tags` - не больше 10 тегов, проверяется в [`Task.Validate`](internal/models/task.go#L81-L88).

`Tag`:

- `name` - от 2 до 30 символов, проверяется в [`Tag.Validate`](internal/models/task.go#L115-L117).

`Report`:

- `report_type` - только `status`, `priority`, `assignee`, значения описаны в [`ReportType`](internal/models/report.go#L3-L9);
- пустой `report_type` автоматически превращается в `status` через [`ReportType.Normalize`](internal/models/report.go#L21-L26).

## 5. Структура проекта

- [`cmd/taskflow/main.go`](cmd/taskflow/main.go) - точка входа приложения;
- [`build/Dockerfile`](build/Dockerfile) - Docker-файл PostgreSQL;
- [`docker-compose.yml`](docker-compose.yml) - запуск PostgreSQL;
- [`frontend/index.html`](frontend/index.html) - простой веб-интерфейс;
- [`internal/app`](internal/app) - HTTP-обработчики;
- [`internal/models`](internal/models) - модели предметной области;
- [`internal/service`](internal/service) - бизнес-логика, фасад и стратегии;
- [`internal/repository`](internal/repository) - методы работы с PostgreSQL;
- [`pkg/auth`](pkg/auth) - генерация и проверка JWT;
- [`pkg/logger`](pkg/logger) - логирование;
- [`migrations/001_init.sql`](migrations/001_init.sql) - создание таблиц БД;
- [`readme.md`](readme.md) - документация.

## 6. Структуры и методы

В Go нет классов в классическом смысле. В этом проекте роль классов выполняют структуры, методы структур и интерфейсы.

### 6.1. Модели предметной области

| Структура или тип | Где реализовано | За что отвечает | Методы |
|---|---|---|---|
| `BaseEntity` | [`common.go`](internal/models/common.go#L9-L13) | Общие поля сущностей: `ID`, `CreatedAt`, `UpdatedAt`. Используется для простого наследования через embedding. | Методов нет |
| `AuditInfo` | [`common.go`](internal/models/common.go#L15-L17) | Информация об авторе создания записи: `CreatedBy`. Используется как часть множественного наследования через embedding. | Методов нет |
| `SoftDelete` | [`common.go`](internal/models/common.go#L19-L21) | Поле мягкого удаления `DeletedAt`. Используется как часть множественного наследования через embedding. | Методов нет |
| `User` | [`user.go`](internal/models/user.go#L8-L14) | Пользователь системы. Может быть владельцем проекта и исполнителем задачи. Хранит логин и хеш пароля. | [`Validate`](internal/models/user.go#L16-L24) |
| `Project` | [`project.go`](internal/models/project.go#L5-L13) | Проект, внутри которого создаются задачи. Хранит владельца и может агрегировать список задач. | [`Validate`](internal/models/project.go#L15-L26) |
| `Status` | [`task.go`](internal/models/task.go#L8-L15) | Тип статуса задачи: `new`, `in_progress`, `done`, `cancelled`. | [`IsValid`](internal/models/task.go#L107-L109) |
| `Priority` | [`task.go`](internal/models/task.go#L17-L24) | Тип приоритета задачи: `low`, `medium`, `high`, `critical`. | [`IsValid`](internal/models/task.go#L111-L113) |
| `Task` | [`task.go`](internal/models/task.go#L26-L39) | Главная сущность приложения. Хранит проект, название, описание, статус, приоритет, дедлайн, исполнителя, теги и историю. | [`Validate`](internal/models/task.go#L59-L90), [`ChangeStatus`](internal/models/task.go#L92-L105) |
| `TaskFilter` | [`task.go`](internal/models/task.go#L41-L45) | Параметры фильтрации задач по статусу, приоритету и исполнителю. | Методов нет |
| `Tag` | [`task.go`](internal/models/task.go#L46-L49) | Тег задачи. Используется для дополнительной классификации задач. | [`Validate`](internal/models/task.go#L115-L117) |
| `TaskHistory` | [`task.go`](internal/models/task.go#L51-L57) | Запись истории изменения статуса задачи. | Методов нет |
| `ReportType` | [`report.go`](internal/models/report.go#L3-L9) | Доменный тип отчета. Может быть `status`, `priority`, `assignee`. | [`Normalize`](internal/models/report.go#L21-L26) |
| `Report` | [`report.go`](internal/models/report.go#L11-L14) | Итоговый отчет: тип отчета и набор агрегированных строк. | Методов нет |
| `ReportItem` | [`report.go`](internal/models/report.go#L16-L19) | Одна строка отчета: подпись и количество задач. | Методов нет |

Дополнительные функции моделей:

- [`validateLength`](internal/models/common.go#L23-L29) проверяет длину строки в диапазоне;
- [`validateMaxLength`](internal/models/common.go#L31-L37) проверяет максимальную длину строки;
- [`dateOnly`](internal/models/common.go#L39-L41) приводит дату к началу дня, чтобы дедлайн сравнивался без учета времени.

### 6.2. Сервисный слой

| Структура или интерфейс | Где реализовано | За что отвечает | Методы |
|---|---|---|---|
| `AppLogger` | [`user.go`](internal/service/user.go#L10-L13) | Абстракция логгера, чтобы сервисы не зависели от конкретной реализации. | `Info`, `Error` |
| `UserStore` | [`user.go`](internal/service/user.go#L20-L27) | Интерфейс репозитория пользователей. | `Create`, `List`, `Exists`, `LoginExists`, `FindByID`, `FindByLogin` |
| `UserService` | [`user.go`](internal/service/user.go#L29-L36) | Бизнес-логика пользователей: валидация, регистрация, хеширование пароля, вход, список. | [`Create`](internal/service/user.go#L38-L59), [`Register`](internal/service/user.go#L61-L76), [`Login`](internal/service/user.go#L78-L90), [`FindByID`](internal/service/user.go#L92-L94), [`List`](internal/service/user.go#L96-L98) |
| `ProjectStore` | [`project.go`](internal/service/project.go#L10-L15) | Интерфейс репозитория проектов. | `Create`, `List`, `ListByOwner`, `Exists` |
| `ProjectService` | [`project.go`](internal/service/project.go#L16-L24) | Бизнес-логика проектов: валидация, проверка владельца, создание, список. | [`Create`](internal/service/project.go#L26-L47), [`List`](internal/service/project.go#L49-L51) |
| `TaskStore` | [`task.go`](internal/service/task.go#L11-L16) | Интерфейс репозитория задач. | `Create`, `List`, `FindByID`, `UpdateStatus` |
| `TaskFacade` | [`task.go`](internal/service/task.go#L18-L24) | Фасад для сценариев работы с задачами. Координирует валидацию, проверки, БД, историю и логирование. | [`Create`](internal/service/task.go#L40-L86), [`List`](internal/service/task.go#L88-L90), [`ChangeStatus`](internal/service/task.go#L92-L109), [`ChangeStatusForAssignee`](internal/service/task.go#L111-L133), [`SetNow`](internal/service/task.go#L36-L38) |
| `ReportStrategy` | [`report.go`](internal/service/report.go#L12-L15) | Интерфейс стратегии построения отчета. | `ReportType`, `Generate` |
| `StatusReportStrategy` | [`report.go`](internal/service/report.go#L17-L29) | Стратегия отчета по статусам задач. | [`ReportType`](internal/service/report.go#L19-L21), [`Generate`](internal/service/report.go#L23-L29) |
| `PriorityReportStrategy` | [`report.go`](internal/service/report.go#L31-L43) | Стратегия отчета по приоритетам задач. | [`ReportType`](internal/service/report.go#L33-L35), [`Generate`](internal/service/report.go#L37-L43) |
| `AssigneeReportStrategy` | [`report.go`](internal/service/report.go#L45-L57) | Стратегия отчета по исполнителям задач. | [`ReportType`](internal/service/report.go#L47-L49), [`Generate`](internal/service/report.go#L51-L57) |
| `ReportService` | [`report.go`](internal/service/report.go#L59-L71) | Загружает задачи, хранит набор стратегий отчета и строит отчет через выбранную стратегию. | [`Build`](internal/service/report.go#L73-L87) |

Дополнительные функции отчетов:

- [`DefaultReportStrategies`](internal/service/report.go#L89-L100) создает набор доступных стратегий;
- [`SelectReportStrategy`](internal/service/report.go#L102-L104) выбирает стратегию по типу отчета;
- [`buildReport`](internal/service/report.go#L115-L127) собирает итоговый `Report` из подсчитанных значений.

### 6.3. Репозитории

| Структура | Где реализовано | За что отвечает | Методы |
|---|---|---|---|
| `UserRepository` | [`repository/user.go`](internal/repository/user.go#L10-L16) | SQL-операции с таблицей `users`. | [`Create`](internal/repository/user.go#L18-L24), [`List`](internal/repository/user.go#L26-L47), [`Exists`](internal/repository/user.go#L49-L53), `LoginExists`, `FindByID`, `FindByLogin` |
| `ProjectRepository` | [`repository/project.go`](internal/repository/project.go#L10-L16) | SQL-операции с таблицей `projects`. | [`Create`](internal/repository/project.go#L18-L24), [`List`](internal/repository/project.go#L26-L47), [`Exists`](internal/repository/project.go#L49-L53) |
| `TaskRepository` | [`repository/task.go`](internal/repository/task.go#L12-L18) | SQL-операции с задачами, тегами и историей статусов. | [`Create`](internal/repository/task.go#L20-L65), [`List`](internal/repository/task.go#L67-L116), [`FindByID`](internal/repository/task.go#L118-L142), [`UpdateStatus`](internal/repository/task.go#L144-L167) |

### 6.4. HTTP-слой и логгер

| Структура | Где реализовано | За что отвечает | Методы |
|---|---|---|---|
| `Server` | [`app/task.go`](internal/app/task.go#L18-L24) | HTTP-сервер приложения. Хранит сервисы, JWT-менеджер и связывает маршруты с обработчиками. | [`Routes`](internal/app/task.go#L44-L55), [`index`](internal/app/task.go#L58-L71), [`me`](internal/app/task.go#L73-L92), [`createTask`](internal/app/task.go#L94-L131), [`changeTaskStatus`](internal/app/task.go#L133-L152), [`render`](internal/app/task.go#L154-L201) |
| `pageData` | [`app/task.go`](internal/app/task.go#L26-L38) | Данные, которые передаются в HTML-шаблон. Содержит флаг авторизации и текущего пользователя. | Методов нет |
| `Manager` | [`jwt.go`](pkg/auth/jwt.go#L14-L17) | Генератор и валидатор JWT. | [`Generate`](pkg/auth/jwt.go#L28-L48), [`Validate`](pkg/auth/jwt.go#L50-L77) |
| `Claims` | [`jwt.go`](pkg/auth/jwt.go#L19-L23) | Данные, которые кладутся в JWT: ID пользователя, логин и срок действия. | Методов нет |
| `Logger` | [`logger.go`](pkg/logger/logger.go#L8-L16) | Обертка над стандартным `log/slog`. | [`Info`](pkg/logger/logger.go#L18-L20), [`Error`](pkg/logger/logger.go#L22-L24) |

## 7. Как реализованы принципы ООП

Инкапсуляция означает, что объект или пакет скрывает детали реализации и дает внешнему коду ограниченный набор методов для работы. В Go это делается не модификаторами `private/public`, а регистром первой буквы: `Name` экспортируется из пакета, а `name` доступен только внутри пакета.

В проекте инкапсуляция показана на нескольких уровнях:

- доменные типы сами проверяют свои ограничения: `User.Validate` [`user.go`](internal/models/user.go#L16-L24), `Project.Validate` [`project.go`](internal/models/project.go#L15-L26), `Task.Validate` [`task.go`](internal/models/task.go#L59-L90);
- изменение статуса задачи вынесено в метод [`Task.ChangeStatus`](internal/models/task.go#L92-L105), поэтому вместе со сменой статуса всегда создается `TaskHistory`;
- зависимости сервисов скрыты в неэкспортируемых полях: например `ReportService.tasks`, `ReportService.logger`, `ReportService.strategies` [`report.go`](internal/service/report.go#L59-L63). Их нельзя поменять напрямую из другого пакета;
- создавать `ReportService` нужно через конструктор [`NewReportService`](internal/service/report.go#L65-L71), где сразу задается корректный набор стратегий;
- выбор стратегии спрятан в неэкспортируемой функции [`selectReportStrategy`](internal/service/report.go#L106-L113), а внешний код пользуется более простым методом [`Build`](internal/service/report.go#L73-L87);
- тип отчета оформлен как `ReportType`, а значение по умолчанию задается методом [`Normalize`](internal/models/report.go#L21-L26), поэтому правило "пустой тип отчета = отчет по статусам" находится в одном месте.

Практическая польза инкапсуляции здесь такая: HTTP-слой не знает, как устроено хеширование пароля, как выбирается стратегия отчета, как пишется история статусов и как устроены SQL-запросы. Он вызывает сервисные методы, а детали остаются внутри своих пакетов.

Абстракция означает работу через общие интерфейсы без знания конкретной реализации.

В проекте это реализовано через:

- [`UserStore`](internal/service/user.go#L20-L27);
- [`ProjectStore`](internal/service/project.go#L10-L14);
- [`TaskStore`](internal/service/task.go#L11-L16);
- [`AppLogger`](internal/service/user.go#L10-L13);
- [`ReportStrategy`](internal/service/report.go#L12-L15).

Полиморфизм означает, что разные типы могут использоваться через один общий интерфейс.

В проекте это реализовано через интерфейс `ReportStrategy`: `StatusReportStrategy`, `PriorityReportStrategy` и `AssigneeReportStrategy` реализуют одинаковые методы `ReportType` и `Generate` [`report.go`](internal/service/report.go#L12-L57). `ReportService` хранит их в поле `strategies` как значения интерфейса [`ReportStrategy`](internal/service/report.go#L59-L63), выбирает нужную стратегию по `report_type` и вызывает [`Generate`](internal/service/report.go#L73-L87). В этот момент конкретная структура может быть разной, но код сервиса остается одинаковым.

## 8. Наследование и множественное наследование в Go

В Go нет классического наследования как в Java, C++ или C#. Вместо этого используется композиция и встраивание структур, то есть embedding.

Embedding работает так: если одна структура содержит другую структуру без имени поля, то поля и методы вложенной структуры становятся доступны как будто они принадлежат внешней структуре.

Простое наследование в проекте:

- `User` встраивает `BaseEntity`, поэтому получает поля `ID`, `CreatedAt`, `UpdatedAt`: [`User`](internal/models/user.go#L8-L14);
- `Project` встраивает `BaseEntity`: [`Project`](internal/models/project.go#L5-L13);
- `Task` встраивает `BaseEntity`: [`Task`](internal/models/task.go#L26-L39);
- `Tag` и `TaskHistory` тоже встраивают `BaseEntity`: [`Tag`](internal/models/task.go#L46-L49), [`TaskHistory`](internal/models/task.go#L51-L57).

Множественное наследование в проекте:

- `User` одновременно встраивает `BaseEntity`, `AuditInfo`, `SoftDelete`: [`User`](internal/models/user.go#L8-L14);
- `Project` одновременно встраивает `BaseEntity`, `AuditInfo`, `SoftDelete`: [`Project`](internal/models/project.go#L5-L13);
- `Task` одновременно встраивает `BaseEntity`, `AuditInfo`, `SoftDelete`: [`Task`](internal/models/task.go#L26-L39).

Что это дает:

- от `BaseEntity` сущности получают технические поля идентификатора и дат;
- от `AuditInfo` получают поле `CreatedBy`;
- от `SoftDelete` получают поле `DeletedAt`;
- код не дублирует одинаковые поля в каждой сущности.

Пример работы embedding в коде: репозиторий может записывать `task.ID`, `task.CreatedAt`, `task.CreatedBy`, хотя эти поля объявлены во встроенных структурах [`BaseEntity`](internal/models/common.go#L9-L13) и [`AuditInfo`](internal/models/common.go#L15-L17). Это используется при сохранении задачи в [`TaskRepository.Create`](internal/repository/task.go#L20-L65).

## 9. Связи между объектами

Ассоциация - это связь, при которой один объект знает о другом, но не владеет его жизненным циклом.

В проекте ассоциация реализована так:

- `Task.AssigneeID` связывает задачу с пользователем-исполнителем: [`Task`](internal/models/task.go#L26-L39);
- в БД это внешний ключ `assignee_id` на таблицу `users`: [`tasks`](migrations/001_init.sql#L23-L36);
- `TaskFacade.Create` проверяет, что проект существует, принадлежит текущему пользователю и что исполнитель существует: [`TaskFacade.Create`](internal/service/task.go#L48-L79).

Агрегация - это связь “целое содержит части”, но части могут существовать отдельно.

В проекте агрегация реализована так:

- `Project` содержит поле `Tasks []Task`: [`Project`](internal/models/project.go#L5-L13);
- задача при этом хранится в отдельной таблице `tasks` и может быть загружена отдельно: [`TaskRepository.List`](internal/repository/task.go#L67-L116).

Композиция - это более сильная связь, когда часть логически принадлежит целому.

В проекте композиция реализована так:

- `Task` содержит `History []TaskHistory`: [`Task`](internal/models/task.go#L26-L39);
- `Task.ChangeStatus` создает запись истории при изменении статуса: [`ChangeStatus`](internal/models/task.go#L92-L105);
- в БД история связана с задачей через `ON DELETE CASCADE`, то есть при удалении задачи история удалится вместе с ней: [`task_history`](migrations/001_init.sql#L51-L57).

Связь многие-ко-многим означает, что одна задача может иметь несколько тегов, а один тег может принадлежать нескольким задачам.

В проекте это реализовано через:

- `Task.Tags []Tag`: [`Task`](internal/models/task.go#L26-L39);
- таблицу `tags`: [`tags`](migrations/001_init.sql#L38-L43);
- связующую таблицу `task_tags`: [`task_tags`](migrations/001_init.sql#L45-L49);
- сохранение тегов в [`TaskRepository.Create`](internal/repository/task.go#L37-L54).

Зависимость означает, что один слой использует другой для выполнения своей работы.

В проекте зависимости такие:

- `app` зависит от сервисов через структуру [`Server`](internal/app/task.go#L16-L21);
- `service` зависит от интерфейсов репозиториев, например [`TaskStore`](internal/service/task.go#L11-L16);
- `repository` зависит от `database/sql` и PostgreSQL;
- `cmd/taskflow/main.go` собирает зависимости вместе в точке входа приложения: [`main`](cmd/taskflow/main.go).

## 10. Паттерны проектирования

### 10.1. Facade

`Facade` - паттерн, который дает простой интерфейс к сложной подсистеме. Вместо того чтобы обработчик HTTP-запроса сам выполнял валидацию, проверки в БД, сохранение, историю и логирование, он вызывает один фасад.

В проекте фасад реализован структурой [`TaskFacade`](internal/service/task.go#L18-L24).

Что делает `TaskFacade.Create`:

- назначает статус `new`, если статус не передан: [`Create`](internal/service/task.go#L40-L43);
- валидирует задачу через `Task.Validate`: [`Create`](internal/service/task.go#L44-L47);
- проверяет существование проекта: [`Create`](internal/service/task.go#L48-L57);
- проверяет, что проект принадлежит текущему пользователю: [`Create`](internal/service/task.go#L58-L69);
- проверяет существование исполнителя: [`Create`](internal/service/task.go#L70-L79);
- сохраняет задачу через репозиторий: [`Create`](internal/service/task.go#L80-L85);
- пишет лог.

Что делает `TaskFacade.ChangeStatusForAssignee`:

- загружает задачу по ID: [`ChangeStatusForAssignee`](internal/service/task.go#L111-L116);
- проверяет, что задача принадлежит текущему пользователю: [`ChangeStatusForAssignee`](internal/service/task.go#L117-L121);
- меняет статус через `Task.ChangeStatus`: [`ChangeStatusForAssignee`](internal/service/task.go#L122-L126);
- сохраняет новый статус и историю: [`ChangeStatusForAssignee`](internal/service/task.go#L127-L133).

Почему паттерн подходит: операции над задачей затрагивают несколько подсистем, и фасад скрывает эту сложность от HTTP-слоя.

### 10.2. Strategy

`Strategy` - паттерн, который выносит разные алгоритмы в отдельные классы/структуры с общим интерфейсом. Это позволяет заменять алгоритм без переписывания сервиса.

`report_type` - это значение из формы отчета [`frontend/index.html`](frontend/index.html#L231-L235). Оно показывает, какой именно отчет нужно построить:

| `report_type` | Какая стратегия выбирается | Что считает |
|---|---|---|
| `status` | [`StatusReportStrategy`](internal/service/report.go#L17-L29) | Количество задач в каждом статусе |
| `priority` | [`PriorityReportStrategy`](internal/service/report.go#L31-L43) | Количество задач с каждым приоритетом |
| `assignee` | [`AssigneeReportStrategy`](internal/service/report.go#L45-L57) | Количество задач по исполнителям |

Пустое значение `report_type` считается значением `status`; это правило находится в [`ReportType.Normalize`](internal/models/report.go#L21-L26).

В проекте общий интерфейс стратегии:

- [`ReportStrategy`](internal/service/report.go#L12-L15).

Конкретные стратегии:

- [`StatusReportStrategy`](internal/service/report.go#L17-L29) считает задачи по статусам;
- [`PriorityReportStrategy`](internal/service/report.go#L31-L43) считает задачи по приоритетам;
- [`AssigneeReportStrategy`](internal/service/report.go#L45-L57) считает задачи по исполнителям.

Набор стратегий создается в [`DefaultReportStrategies`](internal/service/report.go#L89-L100). Выбор стратегии выполняется в [`selectReportStrategy`](internal/service/report.go#L106-L113). Использование выбранной стратегии находится в [`ReportService.Build`](internal/service/report.go#L73-L87): сервис вызывает `strategy.Generate(tasks)` через интерфейс, поэтому это явный пример полиморфизма.

Как `report_type` связан с полиморфизмом:

1. HTTP-слой получает строку из формы и превращает ее в `models.ReportType`: [`task.go`](internal/app/task.go#L88), [`report.go`](internal/app/report.go#L21).
2. `ReportService.Build` передает этот тип в `selectReportStrategy`: [`report.go`](internal/service/report.go#L73-L75).
3. `selectReportStrategy` возвращает значение интерфейса `ReportStrategy`: [`report.go`](internal/service/report.go#L106-L113).
4. Дальше сервис вызывает `strategy.Generate(tasks)`: [`report.go`](internal/service/report.go#L84).

Именно четвертый шаг является полиморфизмом: переменная `strategy` имеет интерфейсный тип `ReportStrategy`, но внутри нее может лежать `StatusReportStrategy`, `PriorityReportStrategy` или `AssigneeReportStrategy`. Go сам вызывает метод конкретной структуры.

Почему паттерн подходит: отчеты имеют одинаковую форму результата `Report`, но разные правила подсчета.

Почему `Facade` и `Strategy` сочетаются: `Facade` управляет сложными пользовательскими сценариями, а `Strategy` отвечает за заменяемый алгоритм отчета. Они решают разные задачи и не конфликтуют.

## 11. JWT-авторизация

Авторизация реализована через JWT без ролей. Роли были удалены, потому что они не влияли на поведение приложения.

Как работает регистрация:

- пользователь вводит `login` и `password` в форме регистрации: [`frontend/index.html`](frontend/index.html#L84-L95);
- обработчик [`register`](internal/app/user.go#L10-L25) вызывает [`UserService.Register`](internal/service/user.go#L61-L76);
- пароль проверяется функцией [`validatePassword`](internal/service/user.go#L100-L105);
- пароль превращается в хеш функцией [`hashPassword`](internal/service/user.go#L107-L114);
- после успешной регистрации сервер создает JWT и кладет его в httpOnly cookie: [`setAuthCookie`](internal/app/user.go#L60-L74).

Как работает вход:

- пользователь вводит `login` и `password`: [`frontend/index.html`](frontend/index.html#L97-L107);
- обработчик [`login`](internal/app/user.go#L27-L42) вызывает [`UserService.Login`](internal/service/user.go#L78-L90);
- пароль проверяется через [`verifyPassword`](internal/service/user.go#L116-L131);
- при успехе сервер выдает JWT-cookie.

Как работает персональная страница:

- JWT хранится в cookie `taskflow_token`;
- метод [`currentUser`](internal/app/task.go#L251-L265) читает cookie, проверяет JWT и загружает пользователя из БД;
- маршрут `/me` показывает страницу текущего пользователя: [`me`](internal/app/task.go#L73-L92);
- проекты загружаются только по `owner_id` текущего пользователя: [`render`](internal/app/task.go#L167-L183);
- задачи и отчеты фильтруются по `AssigneeID` текущего пользователя: [`render`](internal/app/task.go#L154-L183).

JWT реализован в [`pkg/auth`](pkg/auth/jwt.go): [`Generate`](pkg/auth/jwt.go#L28-L48) создает токен, [`Validate`](pkg/auth/jwt.go#L50-L77) проверяет подпись, срок действия и данные пользователя.

## 12. Логирование

Логирование реализовано в [`pkg/logger`](pkg/logger/logger.go) на базе стандартного пакета `log/slog`.

Структура [`Logger`](pkg/logger/logger.go#L8-L16) содержит `slog.Logger`. Методы:

- [`Info`](pkg/logger/logger.go#L18-L20) пишет информационные сообщения;
- [`Error`](pkg/logger/logger.go#L22-L24) пишет ошибки.

Логирование используется в сервисах:

- регистрация и вход пользователя: [`UserService`](internal/service/user.go#L38-L90);
- создание проекта: [`ProjectService.Create`](internal/service/project.go#L26-L47);
- создание задачи и смена статуса: [`TaskFacade`](internal/service/task.go#L40-L133);
- построение отчета: [`ReportService.Build`](internal/service/report.go#L73-L87).

## 13. База данных

Данные хранятся в PostgreSQL. PostgreSQL поднимается через [`docker-compose.yml`](docker-compose.yml) и [`build/Dockerfile`](build/Dockerfile).

Схема БД находится в [`migrations/001_init.sql`](migrations/001_init.sql).

Структура таблиц:

| Таблица | Где описана | Назначение | Основные поля |
|---|---|---|---|
| `users` | [`001_init.sql`](migrations/001_init.sql#L1-L9) | Хранит пользователей приложения. Пользователь может быть владельцем проекта или исполнителем задачи. | `id`, `login`, `password_hash`, `created_at`, `updated_at`, `created_by`, `deleted_at` |
| `projects` | [`001_init.sql`](migrations/001_init.sql#L11-L20) | Хранит проекты. Проект объединяет задачи и имеет владельца. | `id`, `name`, `description`, `owner_id`, `created_at`, `updated_at`, `created_by`, `deleted_at` |
| `tasks` | [`001_init.sql`](migrations/001_init.sql#L22-L35) | Хранит задачи. Задача принадлежит проекту и назначается пользователю. | `id`, `project_id`, `title`, `description`, `status`, `priority`, `deadline`, `assignee_id`, `created_at`, `updated_at`, `created_by`, `deleted_at` |
| `tags` | [`001_init.sql`](migrations/001_init.sql#L37-L42) | Хранит уникальные теги задач. | `id`, `name`, `created_at`, `updated_at` |
| `task_tags` | [`001_init.sql`](migrations/001_init.sql#L44-L48) | Связующая таблица для связи задач и тегов многие-ко-многим. | `task_id`, `tag_id` |
| `task_history` | [`001_init.sql`](migrations/001_init.sql#L50-L56) | Хранит историю изменения статусов задач. | `id`, `task_id`, `old_status`, `new_status`, `changed_at` |

Связи в БД:

- `projects.owner_id` ссылается на `users.id`, то есть у проекта есть владелец;
- `tasks.project_id` ссылается на `projects.id`, то есть задача принадлежит проекту;
- `tasks.assignee_id` ссылается на `users.id`, то есть задача назначается пользователю;
- `task_tags.task_id` ссылается на `tasks.id`;
- `task_tags.tag_id` ссылается на `tags.id`;
- `task_history.task_id` ссылается на `tasks.id` и удаляется каскадно вместе с задачей.

Загрузка и сохранение данных реализованы в репозиториях:

- [`UserRepository`](internal/repository/user.go#L10-L59);
- [`ProjectRepository`](internal/repository/project.go#L10-L88);
- [`TaskRepository`](internal/repository/task.go#L12-L195).

## 14. Веб-интерфейс

Фронтенд находится в [`frontend/index.html`](frontend/index.html).

Интерфейс содержит:

- форму регистрации;
- форму входа;
- персональную страницу пользователя;
- кнопку выхода;
- форму создания проекта;
- форму создания задачи;
- форму фильтрации задач;
- таблицу задач;
- форму смены статуса;
- блок отчета.

HTTP-маршруты задаются в [`Server.Routes`](internal/app/task.go#L44-L55):

- `/` - главная страница;
- `/register` - регистрация пользователя;
- `/login` - вход пользователя;
- `/logout` - выход пользователя;
- `/me` - персональная страница пользователя;
- `/projects` - создание проекта;
- `/tasks` - создание задачи;
- `/tasks/status` - смена статуса;
- `/report` - построение отчета.

## 15. Unit-тестирование

По требованию КМ-5 unit-тесты должны проверять отдельные модули исходного кода. В проекте тестируется нетривиальная бизнес-логика, а не работа браузера или реальной БД.

Тесты моделей:

| Тест | Где реализован | Что проверяет | Вид теста |
|---|---|---|---|
| `TestUserValidatePositive` | [`user_test.go`](internal/models/user_test.go#L8-L13) | Пользователь с корректным `login` проходит валидацию. | Позитивный |
| `TestUserValidateNegativeLoginWithSpace` | [`user_test.go`](internal/models/user_test.go#L15-L20) | Логин с пробелом возвращает ошибку. | Негативный |
| `TestUserValidateBoundaryLogin` | [`user_test.go`](internal/models/user_test.go#L22-L33) | `User.login` ровно 3 и 50 символов допустим. | Граничный |
| `TestProjectValidatePositive` | [`project_test.go`](internal/models/project_test.go#L8-L13) | Проект с корректными данными проходит валидацию. | Позитивный |
| `TestProjectValidateNegativeOwner` | [`project_test.go`](internal/models/project_test.go#L15-L20) | `owner_id = 0` возвращает ошибку. | Негативный |
| `TestProjectValidateBoundaryName` | [`project_test.go`](internal/models/project_test.go#L22-L33) | `Project.name` ровно 3 и 80 символов допустим. | Граничный |
| `TestTaskValidatePositive` | [`task_test.go`](internal/models/task_test.go#L22-L28) | Задача с корректными данными проходит валидацию. | Позитивный |
| `TestTaskValidateNegativeStatus` | [`task_test.go`](internal/models/task_test.go#L30-L37) | Недопустимый статус возвращает ошибку. | Негативный |
| `TestTaskValidateNegativeDeadline` | [`task_test.go`](internal/models/task_test.go#L39-L46) | Дедлайн в прошлом возвращает ошибку. | Негативный |
| `TestTaskValidateBoundaryFields` | [`task_test.go`](internal/models/task_test.go#L48-L65) | Границы `title`, `description`, 10 тегов и дедлайн сегодня. | Граничный |
| `TestTaskValidateNegativeTooManyTags` | [`task_test.go`](internal/models/task_test.go#L67-L77) | 11 тегов возвращают ошибку. | Негативный |
| `TestTaskChangeStatus` | [`task_test.go`](internal/models/task_test.go#L79-L94) | Смена статуса меняет статус и создает историю. | Позитивный |
| `TestTagValidateBoundary` | [`task_test.go`](internal/models/task_test.go#L96-L102) | `Tag.name` ровно 2 и 30 символов допустим. | Граничный |
| `TestReportTypeNormalize` | [`report_test.go`](internal/models/report_test.go#L5-L22) | Пустой тип отчета превращается в `status`, явный тип сохраняется. | Позитивный, граничный |

Тесты сервисов и паттернов:

| Тест | Где реализован | Что проверяет | Вид теста |
|---|---|---|---|
| `TestSelectReportStrategy` | [`report_test.go`](internal/service/report_test.go#L30-L59) | Выбор стратегий `status`, `priority`, `assignee` и построение отчета. | Позитивный, Strategy |
| `TestSelectReportStrategyNegative` | [`report_test.go`](internal/service/report_test.go#L61-L65) | Неизвестный тип отчета возвращает ошибку. | Негативный, Strategy |
| `TestReportServiceBuildUsesSelectedStrategy` | [`report_test.go`](internal/service/report_test.go#L70-L89) | `ReportService.Build` выбирает стратегию и строит отчет через общий интерфейс `ReportStrategy`. | Позитивный, Strategy, полиморфизм |
| `TestTaskFacadeCreatePositive` | [`task_test.go`](internal/service/task_test.go#L87-L115) | `TaskFacade` успешно создает задачу при корректных данных. | Позитивный, Facade |
| `TestTaskFacadeCreateNegativeValidation` | [`task_test.go`](internal/service/task_test.go#L117-L145) | При ошибке валидации задача не сохраняется и ошибка логируется. | Негативный, Facade |
| `TestTaskFacadeChangeStatusWritesHistory` | [`task_test.go`](internal/service/task_test.go#L147-L180) | Смена статуса через фасад сохраняет статус и историю. | Позитивный, Facade |
| `TestTaskFacadeChangeStatusForAssigneeRejectsAnotherUser` | [`task_test.go`](internal/service/task_test.go#L198-L225) | Пользователь не может изменить статус чужой задачи. | Негативный, Facade |
| `TestUserServiceRegisterAndLoginPositive` | [`user_test.go`](internal/service/user_test.go#L71-L93) | Регистрация хеширует пароль, а вход с правильным паролем успешен. | Позитивный, Auth |
| `TestUserServiceRegisterNegativePassword` | [`user_test.go`](internal/service/user_test.go#L95-L100) | Слишком короткий пароль возвращает ошибку. | Негативный, Auth |
| `TestUserServiceLoginNegativePassword` | [`user_test.go`](internal/service/user_test.go#L102-L111) | Неверный пароль при входе возвращает ошибку. | Негативный, Auth |
| `TestManagerGenerateAndValidate` | [`jwt_test.go`](pkg/auth/jwt_test.go#L8-L21) | JWT создается и успешно проверяется. | Позитивный, JWT |
| `TestManagerValidateNegativeSignature` | [`jwt_test.go`](pkg/auth/jwt_test.go#L23-L33) | JWT с неверной подписью отклоняется. | Негативный, JWT |

Для тестов сервисного слоя используются fake-репозитории и fake-логгер:

- [`fakeLogger`](internal/service/task_test.go#L11-L19);
- [`fakeUsers`](internal/service/task_test.go#L21-L40);
- [`fakeProjects`](internal/service/task_test.go#L42-L57);
- [`fakeTasks`](internal/service/task_test.go#L59-L85).

Так тесты остаются unit-тестами и не требуют PostgreSQL.

Запуск тестов:

```bash
go test ./...
```

Если Go в текущей среде не может писать в стандартный cache, можно использовать:

```bash
GOCACHE=/private/tmp/go-build GOMODCACHE=/private/tmp/go-mod go test ./...
```

## 16. Итог реализации

В проекте реализовано:

- веб-приложение на Go;
- регистрация и вход через JWT;
- персональная страница пользователя;
- стандартная точка входа [`cmd/taskflow/main.go`](cmd/taskflow/main.go);
- PostgreSQL через Docker;
- модели предметной области;
- репозитории для работы с БД;
- сервисный слой;
- паттерн `Facade`;
- паттерн `Strategy`;
- логирование;
- минимальный HTML-интерфейс;
- unit-тесты;
- документация по требованиям КМ-4 и КМ-5.
