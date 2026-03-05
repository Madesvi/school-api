package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	_ "log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"rest-api-app/internal/api/handlers"
	mw "rest-api-app/internal/api/middlewares"
	"rest-api-app/internal/api/router"
	"rest-api-app/internal/repositories/postgre"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("No .env files")
	}

	// === Load pprof ===
	go func() {
		pprofAddr := os.Getenv("PPROF_ADDR")
		log.Info().Msgf("pprof server started on: %s", pprofAddr)
		if err := http.ListenAndServe(pprofAddr, nil); err != nil {
			log.Printf("pprof server error: %v", err)
		}
	}()
	// === Load pprof ===

	db, err := postgre.ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Critical: database unavailable")
	}

	provider := postgre.NewTeacherProvider(db)

	env := handlers.NewEnv(provider)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	port := os.Getenv("SERVER_PORT")
	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
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
	router := router.Router(env)
	secureMux := mw.SecurityHeaders(router)

	// Create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux, // HEADERS ON POSTMAN
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Error().Err(errors.New("Error")).Msg("Error starting the server")
	}
}
