package users

import (
  "github.com/gocql/gocql"
)

type User struct {
  ID        gocql.UUID `json:"id"`
  FirstName string     `json:"firstname"`
  LastName  string     `json:"lastname"`
  Email     string     `json:"email"`
  Age       int        `json:"age"`
  City      string     `json:"city"`
}

type GetUsersResponse struct {
  User User `json:"user"`
}

type AllUsersResponse struct {
  Users []User `json:"users"`
  CPU   int64  `json:"micro_cpu"`
}

type NewUserResponse struct {
  ID gocql.UUID `json:"id"`
}

type ErrorResponse struct {
  Errors []string `json:"errors"`
}

type EmptyResponse struct {
  CPU int64 `json:"micro_cpu"`
}
