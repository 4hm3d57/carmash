package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/drive/v3"
	"log"
)

type Car struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	FileID    string             `bson:"file_id"`
	MimeType  string             `bson:"mime_type"`
	Elorating float64            `bson:"rating"`
}

// db connection function
func PicDB() (*mongo.Client, *mongo.Collection, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, nil, err
	}

	picCollection := client.Database("carmash").Collection("images")

	return client, picCollection, nil

}

func InsertImagesLinks(file *drive.File) error {
	client, picCollection, err := PicDB()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// check if the image already exists in the db
	existingFile := Car{}
	filter := bson.M{"file_id": file.Id}
	err = picCollection.FindOne(context.Background(), filter).Decode(&existingFile)
	if err == mongo.ErrNoDocuments {
		filemetadata := Car{
			ID:        primitive.NewObjectID(),
			Name:      file.Name,
			FileID:    file.Id,
			MimeType:  file.MimeType,
			Elorating: 400.0,
		}

		_, err = picCollection.InsertOne(context.Background(), filemetadata)
		if err != nil {
			return err
		}

		log.Println("Inserted file:", filemetadata)
	} else if err != nil {
		return err
	}

	return nil
}
