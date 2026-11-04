package main

import (
	"fmt"
	postgres "messenger/internal/adapter/repository"
	"messenger/internal/controller/http"
	"messenger/internal/usecase"
	"messenger/pkg/jwt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("secret.env"); err != nil {
		fmt.Printf("FATAL ERROR: env variables loading error: %s\n", err.Error())
		os.Exit(1)
	}

	r, err := postgres.NewPostgres("")
	if err != nil {
		fmt.Printf("FATAL ERROR: postgres initialization error: %s\n", err.Error())
		os.Exit(1)
	}

	j := jwt.NewJWT()

	uc := usecase.NewUsecase(r, j)

	hc := http.NewHTTPController(uc)

	hc.UseMiddleware(j.IdentifyMiddleware)
	hc.DefindEndpoints()

	hc.ListenAndServe("localhost:8088")
}
