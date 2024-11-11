package infrastructure

import (
	"log"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	e    *echo.Echo
	db   *gorm.DB
	auth *FirebaseAuth
}

func NewServer(db *gorm.DB, auth *FirebaseAuth) *Server {
	e := echo.New()
	return &Server{
		e:    e,
		db:   db,
		auth: auth,
	}
}

func (s *Server) Serve(addr string) {
	s.setupMiddleware()
	s.setupRoutes()

	log.Printf("Server running on %s", addr)
	if err := s.e.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) setupMiddleware() {
	s.e.Use(echomiddleware.Logger())
	s.e.Use(echomiddleware.Recover())

	s.e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		Skipper:      echomiddleware.DefaultCORSConfig.Skipper,
		AllowOrigins: echomiddleware.DefaultCORSConfig.AllowOrigins,
		AllowMethods: echomiddleware.DefaultCORSConfig.AllowMethods,
		AllowHeaders: []string{"Content-Type", "Accept", "Origin", "X-Token", "Authorization"},
	}))
}

func (s *Server) setupRoutes() {
	// 認証ミドルウェアの作成
	authMiddleware := NewAuthMiddleware(s.auth)

	api := s.e.Group("/api")

	// // 認証不要のエンドポイント
	// public := api.Group("")
	// {
	// 	public.GET("/health", func(c echo.Context) error {
	// 		return c.JSON(200, map[string]string{"status": "ok"})
	// 	})
	// }

	// // 認証必要のエンドポイント
	protected := api.Group("")
	protected.Use(authMiddleware.VerifyToken)
	// {
	// 	// Tripハンドラーの初期化
	// 	tripRepo := NewTripRepository(s.db)
	// 	tripService := NewTripService(tripRepo)
	// 	tripHandler := NewTripHandler(tripService)

	// 	trips := protected.Group("/trips")
	// 	{
	// 		trips.POST("", tripHandler.CreateTrip)
	// 		trips.GET("", tripHandler.GetTrips)
	// 		trips.GET("/:id", tripHandler.GetTrip)
	// 	}
	// }
}
