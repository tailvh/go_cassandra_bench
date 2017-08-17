package users

import (
  "github.com/valyala/fasthttp"
  "strconv"
)

func FormToUser(ctx *fasthttp.RequestCtx) (User, []string) {
  var user User
  var errStr, ageStr string
  var errs []string
  var err error

  user.FirstName, errStr = processFormField(ctx, "firstname")
  errs = appendError(errs, errStr)
  user.LastName, errStr = processFormField(ctx, "lastname")
  errs = appendError(errs, errStr)
  user.Email, errStr = processFormField(ctx, "email")
  errs = appendError(errs, errStr)
  user.City, errStr = processFormField(ctx, "city")
  errs = appendError(errs, errStr)

  ageStr, errStr = processFormField(ctx, "age")
  if len(errStr) != 0 {
    errs = append(errs, errStr)
  } else {
    user.Age, err = strconv.Atoi(ageStr)
    if err != nil {
      errs = append(errs, "Parameters 'age' not an integer")
    }
  }
  return user, errs
}

func appendError(errs []string, errStr string) []string {
  if len(errStr) > 0 {
    errs = append(errs, errStr)
  }
  return errs
}

func processFormField(ctx *fasthttp.RequestCtx, field string) (string, string) {
  fieldData := ctx.FormValue(field)
  if len(fieldData) == 0 {
    return "", "Missing '" + field + "' parameter, cannot continue"
  }

  return string(fieldData[:]), ""
}
