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

var importFlags struct {
	Data         string `type:"path"`
	User         string
	Password     string
	UserPassword string
}

func main() {
	kong.Parse(&importFlags)

	users, err := data.LoadUserData(importFlags.Data)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load users from file")
	}

	start := time.Now()

	client := gocloak.NewClient("https://localhost:18443")
	restyClient := client.RestyClient()

	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	token, err := client.LoginAdmin(importFlags.User, importFlags.Password, "master")
	if err != nil {
		log.Error().Err(err).Msg("failed to authenticate")
		return
	}

	var wg sync.WaitGroup

	for _, u := range users {

		wg.Add(1)

		go func(u *models.User) {

			defer wg.Done()

			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Second)

			log.Info().Str("Name", u.Name).Msg("reading user")

			user := gocloak.User{
				FirstName: gocloak.StringP(u.GivenName),
				LastName:  gocloak.StringP(u.Surname),
				Email:     gocloak.StringP(u.Mail),
				Enabled:   gocloak.BoolP(true),
				Username:  gocloak.StringP(u.SamAccountName),
			}
			id, err := client.CreateUser(token.AccessToken, "master", user)
			if err != nil {
				log.Error().Err(err).Msg("failed to create user")
			}

			users, err := client.GetUsers(token.AccessToken, "master", gocloak.GetUsersParams{Email: gocloak.StringP(u.Mail)})
			if err != nil {
				log.Error().Err(err).Msg("failed to create user")
				return
			}

			if len(users) != 1 {
				log.Warn().Str("email", u.Mail).Msg("failed to locate user")

				return
			}

			err = client.SetPassword(token.AccessToken, gocloak.PString(users[0].ID), "master", importFlags.UserPassword, false)
			if err != nil {
				log.Error().Err(err).Msg("failed to set user password")
			}

			log.Info().Str("id", id).Msg("user created")
		}(u)
	}

	wg.Wait()

	log.Info().Dur("total", time.Since(start)).Msg("creation complete")

}
