package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	"github.com/golang-jwt/jwt/v5"
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

	// Корневой endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Trainers API is running",
			"version": "1.0.0",
		})
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	app.Get("/test/db", func(c *fiber.Ctx) error {
		var count int
		err := db.QueryRow(c.Context(), "SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":  err.Error(),
				"status": "failed",
			})
		}

		// Получим всех пользователей
		rows, err := db.Query(c.Context(), "SELECT id, telegram_id, username FROM users")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var users []fiber.Map
		for rows.Next() {
			var id, username string
			var telegramID int64
			rows.Scan(&id, &telegramID, &username)
			users = append(users, fiber.Map{
				"id":          id,
				"telegram_id": telegramID,
				"username":    username,
			})
		}

		return c.JSON(fiber.Map{
			"users_count": count,
			"users":       users,
			"status":      "connected",
		})
	})

	// Диагностический endpoint 2: Поиск через репозиторий
	app.Get("/test/find-user/:telegram_id", func(c *fiber.Ctx) error {
		var telegramID int64
		fmt.Sscanf(c.Params("telegram_id"), "%d", &telegramID)

		log.Printf("🔍 Searching for telegram_id: %d", telegramID)

		user, err := userRepo.GetByTelegramID(c.Context(), telegramID)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":       err.Error(),
				"telegram_id": telegramID,
			})
		}

		if user == nil {
			log.Printf("⚠️ User not found")
			return c.Status(404).JSON(fiber.Map{
				"error":       "user not found",
				"telegram_id": telegramID,
			})
		}

		log.Printf("✅ User found: %s", user.ID)
		return c.JSON(fiber.Map{
			"user":   user,
			"status": "found",
		})
	})

	// ТЕСТОВЫЙ endpoint для получения токена без Telegram (УДАЛИТЬ В ПРОДАКШЕНЕ!)
	app.Post("/test/login", func(c *fiber.Ctx) error {
		var req struct {
			TelegramID int64 `json:"telegram_id"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		log.Printf("🔍 /test/login: looking for telegram_id=%d", req.TelegramID)

		user, err := userRepo.GetByTelegramID(c.Context(), req.TelegramID)
		if err != nil {
			log.Printf("❌ DB error: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if user == nil {
			log.Printf("⚠️ User not found in DB")
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}

		log.Printf("✅ User found: %s", user.ID)

		claims := &service.Claims{
			UserID:     user.ID,
			Role:       user.Role,
			TelegramID: user.TelegramID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"token": tokenStr,
			"user": fiber.Map{
				"id":         user.ID,
				"username":   user.Username,
				"first_name": user.FirstName,
				"role":       user.Role,
			},
		})
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*", // Для тестирования разрешаем все
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:  "Origin,Content-Type,Accept,Authorization",
		ExposeHeaders: "Content-Length",
		MaxAge:        86400,
	}))
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
