package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trainers-backend/internal/config"
	"trainers-backend/internal/database"
	"trainers-backend/internal/handler"
	"trainers-backend/internal/middleware"
	"trainers-backend/internal/repository"
	"trainers-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// 1. Подключение к БД и Redis
	db, err := database.NewPostgres(ctx, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	rdb, err := database.NewRedis(ctx, cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Fatalf("Redis error: %v", err)
	}
	defer rdb.Close()

	// 2. Инициализация слоев
	userRepo := repository.NewUserRepo(db)
	clientRepo := repository.NewClientRepo(db)
	measRepo := repository.NewMeasurementRepo(db)
	noteRepo := repository.NewNoteRepo(db)
	workoutRepo := repository.NewWorkoutRepo(db)
	schedRepo := repository.NewScheduleRepo(db)
	notifRepo := repository.NewNotificationRepo(db)

	authService := service.NewAuthService(userRepo, rdb, cfg.JWTSecret, cfg.BotToken) // Замените на токен бота
	clientService := service.NewClientService(clientRepo, userRepo)
	measService := service.NewMeasurementService(measRepo)
	noteService := service.NewNoteService(noteRepo)
	workoutService := service.NewWorkoutService(workoutRepo)
	schedService := service.NewScheduleService(schedRepo)
	notifService := service.NewNotificationService(notifRepo)

	authHandler := handler.NewAuthHandler(authService)
	clientHandler := handler.NewClientHandler(clientService, measService, noteService)
	resourceHandler := handler.NewResourceHandler(workoutService, schedService, notifService)

	// 3. Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Trainers API v1.0",
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Публичные роуты
	app.Post("/api/auth/telegram", authHandler.TelegramAuth)

	// Защищенные роуты
	api := app.Group("/api", middleware.AuthMiddleware(authService))

	api.Get("/auth/me", authHandler.GetMe)
	api.Get("/dashboard", authHandler.GetDashboard)

	// Клиенты и вложенные ресурсы
	api.Get("/clients", clientHandler.GetAll)
	api.Get("/clients/:id", clientHandler.GetByID)
	api.Post("/clients", clientHandler.Create)
	api.Patch("/clients/:id", clientHandler.Update)
	api.Delete("/clients/:id", clientHandler.Delete)

	api.Get("/clients/:id/measurements", clientHandler.GetMeasurements)
	api.Post("/clients/:id/measurements", clientHandler.CreateMeasurement)
	api.Delete("/clients/:id/measurements/:measurementId", clientHandler.DeleteMeasurement)

	api.Get("/clients/:id/notes", clientHandler.GetNotes)
	api.Post("/clients/:id/notes", clientHandler.CreateNote)
	api.Patch("/clients/:id/notes/:noteId", clientHandler.UpdateNote)
	api.Delete("/clients/:id/notes/:noteId", clientHandler.DeleteNote)

	// Тренировки
	api.Get("/workouts", resourceHandler.GetAllWorkouts)
	api.Get("/workouts/:id", resourceHandler.GetWorkoutByID)
	api.Post("/workouts", resourceHandler.CreateWorkout)
	api.Patch("/workouts/:id", resourceHandler.UpdateWorkout)
	api.Delete("/workouts/:id", resourceHandler.DeleteWorkout)

	// Расписание
	api.Get("/schedule", resourceHandler.GetSchedule)
	api.Post("/schedule", resourceHandler.CreateSchedule)
	api.Patch("/schedule/:id", resourceHandler.UpdateSchedule)
	api.Delete("/schedule/:id", resourceHandler.DeleteSchedule)

	// Уведомления
	api.Get("/notifications", resourceHandler.GetNotifications)
	api.Patch("/notifications/:id", resourceHandler.MarkNotificationRead)

	// 4. Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		_ = app.Shutdown()
	}()

	log.Printf("Server starting on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
