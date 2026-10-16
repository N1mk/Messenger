package main

import (
	"fmt"
	"messenger/internal/adapter/postgres"
	"messenger/internal/controller/http"
	"messenger/internal/usecase"
	"messenger/pkg/jwt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("secret_keys.env"); err != nil {
		fmt.Printf("FATAL ERROR: env variables loading error: %s", err.Error())
		os.Exit(1)
	}

	p, err := postgres.NewPostgres("")
	if err != nil {
		fmt.Printf("FATAL ERROR: postgres initialization error: %s", err.Error())
		os.Exit(1)
	}

	uc := usecase.NewUsecase(p)

	hc := http.NewHTTPController(uc)

	hc.UseMiddleware(jwt.IdentifyMiddleware)
	hc.DefindEndpoints()

	hc.ListenAndServe("localhost:8088")
}
