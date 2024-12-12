package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/nottgy/discord-oauth-proxy/pkg/discord"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no env file found")
	}
}

func env() (string, string, string) {
	var res []string
	for _, e := range []string{
		"API_KEY",
		"DISCORD_ID",
		"DISCORD_SECRET",
		"REDIRECT_URI",
	} {
		a, exists := os.LookupEnv(e)
		if !exists {
			log.Fatalf("Missing %s", e)
		}
		res = append(res, a)
	}
	return res[0], res[1], res[2]
}

func retrieveUsername(code string) (string, error) {
	_, DISCORD_ID, DISCORD_SECRET, REDIRECT_URI := env()

	res, err := discord.Auth(
		code,
		DISCORD_ID,
		DISCORD_SECRET,
		REDIRECT_URI,
	)
	if err != nil {
		return "", err
	}

	access_token := res.AccessToken
	token_type := res.TokenType

	res2, err = discord.Me(
		access_token, token_type,
	)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf(
		"%s#%s",
		res2.Username,
		res2.Discriminator,
	)

	return name, nil
}

func main() {
	API_KEY, _, _, _ := env()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {

		if len(r.Header["Api-Key"]) < 1 || r.Header["Api-Key"][0] != API_KEY {
			http.Error(w, "Authenticate", http.StatusForbidden)
			return
		}
		hasCode := r.URL.Query().Has("code")
		if !hasCode {
			http.Error(
				w,
				"Usage /auth?code=<code>",
				http.StatusBadRequest,
			)
			return
		}
		code := r.URL.Query().Get("code")
		nickname, err := retrieveUsername(code)
		if err != nil {
			http.Error(
				w,
				"Discord API error",
				http.StatusInternalServerError,
			)
			return
		}
		io.WriteString(w, nickname)
	})

	port := 4000
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Println(err)
	}
}
