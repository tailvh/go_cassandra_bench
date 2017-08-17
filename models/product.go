package models

import "gopkg.in/mgo.v2/bson"

type (
  Product struct {
    Id   bson.ObjectId `json:"id" bson:"_id"`
    Code string        `json:"code" bson:"code"`
    Key  string        `json:"key" bson:"key"`
  }
)
