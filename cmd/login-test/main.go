package main

import (
	"crypto/tls"
	"math/rand"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v5"
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/keycloak-import/internal/data"
	"github.com/wolfeidau/keycloak-import/internal/models"
)

var loginFlags struct {
	Data         string `type:"path"`
	User         string
	Password     string
	ClientID     string
	Secret       string
	UserPassword string
}

func main() {
	kong.Parse(&loginFlags)

	users, err := data.LoadUserData(loginFlags.Data)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load users from file")
	}

	start := time.Now()

	adminClient := gocloak.NewClient("https://localhost:18443")
	restyClient := adminClient.RestyClient()
	//restyClient.SetDebug(true)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	var wg sync.WaitGroup

	for _, u := range users {

		wg.Add(1)

		go func(u *models.User) {

			defer wg.Done()

			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Second)

			log.Info().Str("Name", u.Name).Msg("reading user")

			_, err = adminClient.Login(loginFlags.ClientID, loginFlags.Secret, "master", u.Mail, loginFlags.UserPassword)
			if err != nil {
				log.Error().Err(err).Msg("failed to authenticate")
				return
			}

		}(u)
	}

	wg.Wait()

	log.Info().Dur("total", time.Since(start)).Msg("creation complete")

}
