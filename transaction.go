package main

import (
  "fmt"
  "github.com/valyala/fasthttp"
  //"log"
  "gocassandra/cassandra"
  "gocassandra/users"
  "runtime"
)

type heartbeatResponse struct {
  Status string `json:"status"`
  Code   int    `json:"code"`
}

func main() {
  runtime.GOMAXPROCS(16)
  CassandraSession := cassandra.Session
  defer CassandraSession.Close()

  /////////////////
  // R O U T E S //
  /////////////////
  m := func(ctx *fasthttp.RequestCtx) {
    switch string(ctx.Path()) {
    case "loaderio-2af1a3c64e36faa4782f1b876a209252":
      fmt.Fprintf(ctx, "%s", "loaderio-2af1a3c64e36faa4782f1b876a209252")
    case "/index":
      fmt.Fprintf(ctx, "%s", "Homepage here")
    case "/empty":
      users.Empty(ctx)
    case "/users/new":
      users.Post(ctx)
    case "/users":
      users.Get(ctx)
    default:
      ctx.Error("Not found", fasthttp.StatusNotFound)
    }
  }

  fasthttp.ListenAndServe(":5000", m)
}
