package main

import (
	"fmt"
	postgres "messenger/internal/adapter/repository"
	"messenger/internal/controller/http"
	"messenger/internal/usecase"
	"messenger/pkg/jwt"
	"messenger/pkg/logger"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("FATAL ERROR: executable path getting error: %s\n", err.Error())
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)

	l, err := logger.NewLogger(filepath.Join(exeDir, "/logs/app.log"), filepath.Join(exeDir, "/logs/old"))
	if err != nil {
		fmt.Printf("FATAL ERROR: logger initialization error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := godotenv.Load("secret.env"); err != nil {
		l.LogFatal(fmt.Sprintf("FATAL ERROR: env variables loading error: %s\n", err.Error()))
		os.Exit(1)
	}

	r, err := postgres.NewPostgres("")
	if err != nil {
		l.LogFatal(fmt.Sprintf("FATAL ERROR: postgres initialization error: %s\n", err.Error()))
		os.Exit(1)
	}

	j := jwt.NewJWT()

	uc := usecase.NewUsecase(r, l, j)

	hc := http.NewHTTPController(uc)

	hc.UseMiddleware(j.IdentifyMiddleware)
	hc.DefindEndpoints()

	hc.ListenAndServe("localhost:8088")
}
