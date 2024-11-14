package server

import (
	"log"

	"firebase.google.com/go/v4/db"
	"github.com/Mire0726/safe_travel/backend/api/handler"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore/datastoresql"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	e    *echo.Echo
	db   *gorm.DB
	auth *firebase.FirebaseAuth
}

func NewServer(db *gorm.DB, auth *firebase.FirebaseAuth) *Server {
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
	authMiddleware := firebase.NewAuthMiddleware(s.auth)

	// 認証不要のエンドポイント
	public := s.e.Group("")
	{
		public.GET("/health", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"status": "ok"})
		})

		authClient, err := firebase.NewFirebaseAuth()
		if err != nil {
			log.Fatalf("Failed to create firebase auth client: %v", err)
		}

		

		data := datastoresql.NewStore(dbclient, log.Default())

		// ハンドラーの初期化
		handlerCmd := handler.NewHandler(*authClient, data)
		public.POST("/signUp", handlerCmd.SignUp)
	}

	// // 認証必要のエンドポイント
	protected := s.e.Group("/user")
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
