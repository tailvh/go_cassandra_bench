package models

import "gopkg.in/mgo.v2/bson"

type (
  Card struct {
    _id         bson.ObjectId `json:"_id" bson:"_id"`
    Id          string        `json:"id" bson:"id"`
    Telco_code  string        `json:"telco_code" bson:"telco_code"`
    Name        string        `json:"name" bson:"name"`
    ChannelMeta string        `json:"channelMeta" bson:"channelMeta"`
    Code        string        `json:"code" bson:"code"`
    Key         string        `json:"key" bson:"key"`
  }
)
