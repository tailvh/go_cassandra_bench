package zinc

import (
  "bytes"
  "encoding/json"
  "fmt"
  "strings"
  "time"

  "github.com/valyala/fasthttp"
)

var rxtimezone *time.Location

type OrderedMap map[string]interface{}

func RxGet(data map[string]interface{}, params ...interface{}) interface{} {

  var returnData interface{} = ""
  var indexName interface{} = "default"
  var defaultValue interface{} = ""

  if len(params) >= 1 {
    indexName = params[0]
  }

  if len(params) >= 2 {
    defaultValue = params[1]
  }

  if val, ok := data[indexName.(string)]; ok {
    returnData = val
  } else {
    returnData = defaultValue
  }

  return returnData
}

func RxGetMapMix(_data map[string][]string, _param string, _default []string) []string {

  var returnData = []string{"", ""}
  var indexName string = _param
  var defaultValue = []string{"", ""}

  if len(_default) > 0 {
    defaultValue = _default
  }

  if val, ok := _data[indexName]; ok {
    returnData = val
  } else {
    returnData = defaultValue
  }

  return returnData
}

func RxJsonDecode(jsonstr string) []interface{} {
  returnData := []interface{}{}
  uj := []interface{}{}

  if err := json.Unmarshal([]byte(jsonstr), &uj); err == nil {
    returnData = uj
  }

  return returnData
}

func RxJsonDecodeObj(jsonstr string) map[string]interface{} {
  returnData := map[string]interface{}{}
  uj := map[string]interface{}{}

  if err := json.Unmarshal([]byte(jsonstr), &uj); err == nil {
    returnData = uj
  }

  return returnData
}

func RxJson(data map[string]interface{}) string {
  returnData := "{}"
  if uj, err := json.Marshal(data); err == nil {
    returnData = string(uj)
  }

  return returnData
}

func (om OrderedMap) ToJson(order ...string) string {
  buf := &bytes.Buffer{}
  buf.Write([]byte{'{'})
  l := len(order)
  for i, k := range order {
    if om[k] != nil {
      fmt.Fprintf(buf, "\"%s\": \"%v\"", k, om[k])
      if i < l-1 {
        buf.WriteByte(',')
      }
    }
  }
  buf.Write([]byte{'}'})

  // replace ,}
  return strings.Replace(buf.String(), ",}", "}", -1)
}

func RxIsIn(a interface{}, list []interface{}) bool {
  for _, b := range list {
    if b == a {
      return true
    }
  }

  return false
}

func RxResponse(ctx *fasthttp.RequestCtx, data OrderedMap) bool {

  // Default value for data
  data["status"] = RxGet(data, "status", 0)
  data["msg"] = RxGet(data, "msg", "Not success")
  data["amount"] = RxGet(data, "amount", 0)
  data["trans_id"] = RxGet(data, "trans_id", "")

  // Data return
  uj := data.ToJson("status", "msg", "amount", "trans_id", "sign")
  ctx.Response.Header.Set("Content-Type", "application/json")
  ctx.SetStatusCode(200)
  fmt.Fprintf(ctx, "%s", uj)

  return true
}

func RxTimestamp() int64 {
  rxtimezone, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
  return time.Now().In(rxtimezone).UnixNano()
}
