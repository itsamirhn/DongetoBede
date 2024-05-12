package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DongCollectionName = "dongs"
)

type Dong struct {
	ID           *primitive.ObjectID `bson:"_id,omitempty"`
	IssuerUserID int64               `bson:"issuer_user_id"`
	Amount       int                 `bson:"amount"`
	CardNumber   string              `bson:"card_number"`
	TotalPeople  int                 `bson:"total_people"`
	PaidUserIDs  []int64             `bson:"paid_user_ids"`
	MessageID    string              `bson:"message_id"`
}
