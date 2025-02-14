package server

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func UseAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	// hashingKey := os.Getenv("HASHING_KEY")

	// store := sessions.NewCookieStore([]byte(hashingKey))
	// store.MaxAge(maxAge)
	// store.Options.Path = "/"
	// store.Options.HttpOnly = true
	// store.Options.Secure = true
	//
	// gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://127.0.0.1:3000/auth/callback?provider=google", "email", "profile"),
	)
}
