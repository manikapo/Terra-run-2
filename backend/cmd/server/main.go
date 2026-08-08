package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/territory-run/api/internal/config"
	"github.com/territory-run/api/internal/db"
	"github.com/territory-run/api/internal/handlers"
	"github.com/territory-run/api/internal/middleware"
	"github.com/territory-run/api/internal/services"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.SupabaseJWTSecret == "" {
		log.Fatal("SUPABASE_JWT_SECRET is required")
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	activitySvc := services.NewActivityService(pool, cfg)
	territorySvc := services.NewTerritoryService(pool, cfg)
	userSvc := services.NewUserService(pool)

	activityHandler := handlers.NewActivityHandler(activitySvc, userSvc)
	territoryHandler := handlers.NewTerritoryHandler(territorySvc)
	userHandler := handlers.NewUserHandler(userSvc)

	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Authorization, Content-Type, X-Internal-Secret, X-Guest-User",
	}))

	// Lightweight — safe for UptimeRobot every 5 min (no DB)
	app.Get("/health", handlers.Health)
	app.Get("/", handlers.Health)

	api := app.Group("/api/v1")
	auth := middleware.AuthOrGuest(cfg.SupabaseJWTSecret, cfg.AllowGuestAuth)

	api.Get("/users/me", auth, userHandler.Me)

	api.Post("/activities", auth, activityHandler.Create)
	api.Post("/activities/:id/points", auth, activityHandler.UploadPoints)
	api.Post("/activities/:id/complete", auth, activityHandler.Complete)
	api.Get("/activities", auth, activityHandler.List)

	api.Get("/territories/tile/:h3_tile", auth, territoryHandler.Tile)

	// Phase 2: decay worker stub
	internal := app.Group("/internal", middleware.InternalJobAuth(cfg.InternalJobSecret))
	internal.Post("/jobs/decay", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "decay job stub — implement in Phase 2"})
	})

	go func() {
		addr := ":" + cfg.Port
		log.Printf("territory-run API listening on %s (env=%s, guest_auth=%v)", addr, cfg.Env, cfg.AllowGuestAuth)
		if err := app.Listen(addr); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")
	_ = app.ShutdownWithTimeout(10 * time.Second)
}
