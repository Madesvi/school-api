package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"rest-api-app/internal/api/handlers/execs"
	"rest-api-app/internal/api/handlers/students"
	"rest-api-app/internal/api/handlers/teachers"
	mw "rest-api-app/internal/api/middlewares"
	"rest-api-app/internal/api/router"
	"rest-api-app/internal/repositories/postgre"

	"github.com/joho/godotenv"
)

// This is interprise software where
// admin will create new user

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env files", "err", err)
	}

	// === Load pprof ===
	// pprofAddr := os.Getenv("PPROF_ADDR")
	// go func() {
	// 	slog.Info("pprof server started on port", "port", pprofAddr)
	// 	if err := http.ListenAndServe(pprofAddr, nil); err != nil {
	// 		slog.Error("pprof server error", "err", err)
	// 	}
	// }()
	// === Load pprof ===

	// Load logger
	logLevel := os.Getenv("LOG_LEVEL")

	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}

	if os.Getenv("APP_ENV") == "development" {
		handler = slog.NewTextHandler(os.Stderr, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	// Load logger

	slog.Debug("log level set", "value", logLevel)

	db, err := postgre.ConnectDB()
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}

	// providerT := postgre.NewTeacherProvider(db)
	providerS := postgre.NewStudentProvider(db)
	providerT := postgre.NewTeacherProvider(db)
	providerE := postgre.NewExecProvider(db)

	// root := handlers.RootHandler()

	h := router.Handlers{
		Teachers: teachers.NewAPI(providerT),
		Students: students.NewAPI(providerS),
		Execs:    execs.NewAPI(providerE),
	}

	port := os.Getenv("SERVER_PORT")
	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	cfg := struct {
		JWTSecret string
	}{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	// rl := mw.NewRateLimiter(5, time.Minute)
	//
	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	// secureMux := mw.Cors(rl.Middleware((mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	// secureMux := utils.ApplyMiddleWares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, rl.Middleware, mw.Cors)
	router := router.Router(h)
	jwtMiddleware := mw.ExcludeRoutesMiddleware(
		mw.JWTMiddleware([]byte(cfg.JWTSecret)),
		"/execs/login",
		"/login",
		"/public/",
		"/assets",
		"/register",
		"/public/favicon_io/favicon.ico",
	)
	secureMux := jwtMiddleware(mw.SecurityHeaders(router))
	// secureMux := mw.SecurityHeaders(router)

	// Create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	slog.Info("Server is running on port", "port", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		slog.Error("Error starting the server", "err", err)
	}
}
