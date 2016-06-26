package backend

import (
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type User struct {
	ID       int64 `datastore:"-" goon:"id"`
	Email    string
	Password string `json:"-"`
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
