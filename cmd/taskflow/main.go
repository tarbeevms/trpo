package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"taskflow/internal/app"
	"taskflow/internal/repository"
	"taskflow/internal/service"
	"taskflow/pkg/auth"
	"taskflow/pkg/logger"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://taskflow:taskflow@localhost:5432/taskflow?sslmode=disable"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "taskflow-dev-secret"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	logg := logger.New()
	users := repository.NewUserRepository(db)
	projects := repository.NewProjectRepository(db)
	tasks := repository.NewTaskRepository(db)

	userService := service.NewUserService(users, logg)
	projectService := service.NewProjectService(projects, users, logg)
	taskFacade := service.NewTaskFacade(tasks, projects, users, logg)
	reportService := service.NewReportService(tasks, logg)
	authManager := auth.NewManager(jwtSecret, 24*time.Hour)

	server := app.NewServer(userService, projectService, taskFacade, reportService, authManager)
	logg.Info("taskflow started", "addr", ":8080")
	if err := http.ListenAndServe(":8080", server.Routes()); err != nil {
		log.Fatal(err)
	}
}
