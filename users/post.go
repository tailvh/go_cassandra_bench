package users

import (
  "encoding/json"
  "fmt"
  "github.com/gocql/gocql"
  "github.com/valyala/fasthttp"
  "gocassandra/cassandra"
  "gocassandra/zinc"
)

func Post(ctx *fasthttp.RequestCtx) {
  var errs []string
  var gocqlUuid gocql.UUID

  user, errs := FormToUser(ctx)
  var created bool = false

  if len(errs) == 0 {
    fmt.Println("Creating a new user")
    gocqlUuid = gocql.TimeUUID()

    // write data to cassandra
    if err := cassandra.Session.Query(
      `INSERT INTO users (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
      gocqlUuid, user.FirstName, user.LastName, user.Email, user.City, user.Age).Exec(); err != nil {
      errs = append(errs, err.Error())
    } else {
      created = true
    }
  }

  if created {
    fmt.Println("user_id", gocqlUuid)
    json.NewEncoder(ctx).Encode(NewUserResponse{ID: gocqlUuid})
  } else {
    fmt.Println("errors", errs)
    json.NewEncoder(ctx).Encode(ErrorResponse{Errors: errs})
  }
}

func Get(ctx *fasthttp.RequestCtx) {
  var startTime = zinc.RxTimestamp()

  var userList []User
  m := map[string]interface{}{}

  query := "SELECT id, age, firstname, lastname, city, email FROM users"
  iterable := cassandra.Session.Query(query).Iter()
  for iterable.MapScan(m) {
    userList = append(userList, User{
      ID:        m["id"].(gocql.UUID),
      Age:       m["age"].(int),
      FirstName: m["firstname"].(string),
      LastName:  m["lastname"].(string),
      Email:     m["email"].(string),
      City:      m["city"].(string),
    })
    m = map[string]interface{}{}
  }

  json.NewEncoder(ctx).Encode(AllUsersResponse{Users: userList, CPU: (zinc.RxTimestamp() - startTime) / 1000})
}

func Empty(ctx *fasthttp.RequestCtx) {
  var startTime = zinc.RxTimestamp()
  json.NewEncoder(ctx).Encode(EmptyResponse{CPU: (zinc.RxTimestamp() - startTime) / 1000})
}
