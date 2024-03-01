package main

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
type requestData struct {
	LongUrl   string `json:"longUrl"`
	UserToken string `json:"userToken"`
}
type urlData struct {
	LongUrl    string `bson:"longUrl"`
	UserID     string `bson:"userId"`
	VisitCount int    `bson:"visitCount"`
}
