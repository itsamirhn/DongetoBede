package entities

const (
	UserCollectionName = "users"
)

type User struct {
	ID         int64  `bson:"_id"`
	FirstName  string `bson:"first_name"`
	LastName   string `bson:"last_name"`
	Username   string `bson:"username"`
	CardNumber string `bson:"card_number"`
	State      string `bson:"state"`
}
