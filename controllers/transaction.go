package controllers

import (
  "fmt"
  "time"  
  "math/rand"
  "strconv"
  "crypto/md5"
  "encoding/hex"

  "github.com/valyala/fasthttp"
  
  zmq "github.com/pebbe/zmq4"

  "zmqlog/zinc"  
)

type (
  TransactionController struct {
  }
)

func NewTransactionController() *TransactionController {
  return &TransactionController{}
}

func makeTimestamp() int64 {
  return time.Now().UnixNano()
}

func (uc TransactionController) GetFile(ctx *fasthttp.RequestCtx) {
  fmt.Fprintf(ctx, "%s", "loaderio-1d07256e549d80c5701bcb1cf8f9fb44")
}

func HelpGenTransId (product_id interface{}) string {  
  uuid    := RandStringBytesMaskImprSrc(5)
  tuuid   := product_id.(string) + strconv.FormatInt(time.Now().UnixNano(), 10) + uuid

  return tuuid;
}

func HelpGetMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}

func HelpWriteLog(data map[string]interface{}) bool {
  return true
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandStringBytesMaskImprSrc(n int) string {
  b := make([]rune, n)
  for i := range b {
      b[i] = letters[rand.Intn(len(letters))]
  }
  return string(b)
}

func (uc TransactionController) GetRequest(ctx *fasthttp.RequestCtx) {

  // Variables
  startTime := makeTimestamp()

  // Connect to socket 
  socket, _ := zmq.NewSocket(zmq.REQ)
  defer socket.Close()  

  // Send message to socket
  socket.Connect("tcp://127.0.0.1:5555")
  socket.Send("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 0)
  msg, _ := socket.Recv(0)

  if msg != "1" {
    ctx.Error("Unsupported path", fasthttp.StatusNotFound)
    return
  }

  // Response
  deltaTime := (makeTimestamp() - startTime) / 1000
  zinc.RxResponse(ctx, map[string]interface{} {"status": 1, "msg": "Transaction success", "amount": 10000, "debug": deltaTime})
  return
}