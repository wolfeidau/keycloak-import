package data

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/keycloak-import/internal/models"
)

// LoadData load user data
func LoadUserData(path string) ([]*models.User, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	csvr := csv.NewReader(f)

	columns, err := csvr.Read()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse record")
	}

	log.Info().Strs("headers", columns).Msg("reading csv")

	var users []*models.User

	for {
		fields, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse record")
		}

		m := map[string]string{}

		for n, c := range columns {
			m[c] = fields[n]
		}
		user := new(models.User)
		err = mapstructure.Decode(&m, user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal record")
		}

		users = append(users, user)
	}

	return users, nil
}
