package server

import (
	"database/sql"
	"log"

	"github.com/Mire0726/safe_travel/backend/api/handler"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/datastore/datastoresql"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e    *echo.Echo
	db   *sql.DB
	auth *firebase.FirebaseAuth
}

func NewServer(db *sql.DB, auth *firebase.FirebaseAuth) *Server {
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

		dbCfg, err := infrastructure.LoadDBConfig()
		if err != nil {
			log.Fatalf("Failed to load db config: %v", err)
		}

		dbClient, err := infrastructure.NewDB(dbCfg)
		if err != nil {
			log.Fatalf("Failed to connect to db: %v", err)
		}

		data := datastoresql.NewStore(dbClient, log.Default())

		// ハンドラーの初期化
		handlerCmd := handler.NewHandler(*authClient, data)
		public.POST("/signUp", handlerCmd.SignUp)
		public.POST("/signIn", handlerCmd.SignIn)
	}

	// 認証必要のエンドポイント
	protected := s.e.Group("/user")
	protected.Use(authMiddleware.VerifyToken)
	{
		authClient, err := firebase.NewFirebaseAuth()
		if err != nil {
			log.Fatalf("Failed to create firebase auth client: %v", err)
		}

		dbCfg, err := infrastructure.LoadDBConfig()
		if err != nil {
			log.Fatalf("Failed to load db config: %v", err)
		}

		dbClient, err := infrastructure.NewDB(dbCfg)
		if err != nil {
			log.Fatalf("Failed to connect to db: %v", err)
		}

		data := datastoresql.NewStore(dbClient, log.Default())

		// ハンドラーの初期化
		handlerCmd := handler.NewHandler(*authClient, data)
		protected.DELETE("/delete", handlerCmd.Delete)
	}
}
