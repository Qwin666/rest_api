package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rest_api/internal/pkg/logging"
	"rest_api/internal/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) ( string,  error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil{
		return  "", fmt.Errorf ("file to create user due to erro %v",err)
	}
	d.logger.Debug("convert insertedID to ObjactID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("filed to convert objectid to hex. oid: %s", oid)
}
func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
  oid, err := primitive.ObjectIDFromHex(id)
   if err != nil {
	   return u, fmt.Errorf("faled to convert objactid, Hex:%s", id)
   }
   filter:= bson.M{"_id": oid}
   result:= d.collection.FindOne(ctx, filter)
   if result.Err()!= nil{
	   if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		   return u, fmt.Errorf("ErrorEntityNotFound")
	   }
	   return u, fmt.Errorf("failed to finde one user bu id:%s bue to error %v",id, err)
   }
   if err := result.Decode(&u); err !=nil{
	   return u, fmt.Errorf("failed to decod user id:%s bue to error %v",id, err)
   }
return u, nil
}
func (d *db) Update(ctx context.Context, user user.User) error {
 objectID, err := primitive.ObjectIDFromHex(user.ID)
 if err != nil {
	 return fmt.Errorf("faled to convert user Id to objectId, ID=%s", user.ID)
 }
	 filter := bson.M{"_id":objectID}

	 userBytes, err := bson.Marshal(user)
	 if err!=nil {
		 return fmt.Errorf("faled to marshal user. err6 %v", err)
	 }
		 var updateUserObj bson.M
		 err = bson.Unmarshal(userBytes, updateUserObj)
		 if err != nil {
			 return fmt.Errorf("faled to unmarshal user bytes, err: %v",err)
		 }
		 delete(updateUserObj, "_id")
		 update := bson.M{
			 "$set": updateUserObj,
		 }
		 result, err := d.collection.UpdateOne(ctx, filter, update)
		if  err != nil {
			return fmt.Errorf("faled to execute update user query. err: %v",err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("matched %d documents and modified %d",result.MatchedCount, result.ModifiedCount)
	return nil

}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("faled to convert user Id to objectId, ID=%s", id)
	}
	filter := bson.M{"_id":objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. err: %v",err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("deleted %d documents ",result.DeletedCount)
	return nil

}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger: logger,
	}
	
}
