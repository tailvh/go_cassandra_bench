package main

import (
  "bufio"
  "fmt"
  "io"
  "net"
  "strings"
  "time"
  //"strconv"
  //"sync/atomic"
  "runtime"

  //"zmqlog/component"
  "zmqlog/logger"
)

//var CSlice1 component.ConcurrentSlice
func main() {
  runtime.GOMAXPROCS(16)
  logger.Init("./log", 1000000000, 0, 100, false)

  /////////////////
  // R O U T E S //
  /////////////////
  // Run write to file
  go CheckQueue()

  // Start server listen
  ln, err := net.Listen("tcp", ":5555")
  if err != nil {
    // handle error
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      // handle error
      fmt.Println(err)
    }
    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  for {

    // If closed
    if _, err := conn.Read([]byte{}); err == io.EOF {
      conn.Close()
      conn = nil
      return

    } else {
      message, err := bufio.NewReader(conn).ReadString('\n')
      if err != nil {
        conn.Close()
        conn = nil
        return

      }

      //CSlice1.Append("1")
      fmt.Fprintf(conn, "success\n")
      go logger.Info(strings.TrimSpace(message))
    }

  }
}

func CheckQueue() {
  quitChan := make(chan bool)
  t := time.NewTicker(time.Second)
  func() {
    for {
      select {
      case <-t.C:
      case <-quitChan:
        t.Stop()
        return
      }
      //fmt.Println("CountAll: " + strconv.Itoa(CSlice1.Len()))
    }
  }()
}
