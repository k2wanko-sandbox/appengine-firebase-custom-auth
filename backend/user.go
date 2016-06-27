package backend

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type User struct {
	ID    int64 `datastore:"-" goon:"id"`
	Email string
}

func (u *User) RefreshToken() (string, error) {
	tok := jwt.New(jwt.SigningMethodRS256)
	cl := tok.Claims.(jwt.MapClaims)
	cl["type"] = "refresh_token"
	cl["uid"] = fmt.Sprintf("%d", u.ID)
	return tok.SignedString(privateKey)
}

func verifyRefreshToken(token string) bool {
	tok, err := parseTokenString(token)
	if err != nil {
		return false
	}
	if !tok.Valid {
		return false
	}

	return tok.Claims["type"] == "refresh_token"
}

func userFromEmail(ctx context.Context, email string) (*User, error) {
	g := goon.FromContext(ctx)
	q := datastore.NewQuery("User").
		Filter("Email =", email)
	for t := g.Run(q); ; {
		u := new(User)
		_, err := t.Next(u)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
}
