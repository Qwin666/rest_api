package mongobd

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authBD string) (bd *mongo.Database, err error){
var mongoBDURL string
var isAuth bool
if username == "" && password == "" {
	mongoBDURL = fmt.Sprintf("mongobd://%s:%s", host, port)
} else{
	isAuth = true
	mongoBDURL = fmt.Sprintf("mongobd://%s:%s@%s:%s", username, password, host, port)
}
 clientOption := options.Client().ApplyURI(mongoBDURL)
 if isAuth {
	 if authBD==""{
		 authBD = "database"
	 }
	 clientOption.SetAuth(options.Credential{
		AuthSource: authBD,
		Username: username,
		Password: password,
	 })
}

client, err :=mongo.Connect(ctx, clientOption)
if err != nil {
	return nil, fmt.Errorf("error to connect to mongoDB due to error %v", err)
}

if err = client.Ping(ctx, nil); err!= nil {
	return nil, fmt.Errorf("error to ping to mongoDB due to error %v", err)
}
 return client.Database(database), nil
}

