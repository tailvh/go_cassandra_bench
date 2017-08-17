package models

import "gopkg.in/mgo.v2/bson"

type (
  Channel struct {
    _id  bson.ObjectId `json:"_id" bson:"_id"`
    Id   string        `json:"id" bson:"id"`
    Name string        `json:"name" bson:"name"`
  }
)
