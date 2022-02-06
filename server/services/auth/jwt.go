package auth

import (
	"context"
	"net/http"
)

type User struct {
	Id uint
	Email string
	Name string
	Password string
}

type ContextUser string

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// TODO retrieve the user data from the JWT token
		user := User{
			Id: 1,
			Email: "wilcokuyper@hotmail.com",
			Name: "Wilco",
		}
		r = r.WithContext(context.WithValue(r.Context(), ContextUser("user"), user))
		next(w, r)
	})
}