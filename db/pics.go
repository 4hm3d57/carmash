package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/drive/v3"
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

	picCollection := client.Database("photos").Collection("images")

	return client, picCollection, nil

}

// function to insert the image links to the database
func InsertImagesLinks(file *drive.File) error {

	// create database connection
	client, picCollection, err := PicDB()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	filemetadata := Car{
		ID:        primitive.NewObjectID(),
		Name:      file.Name,
		FileID:    file.Id,
		MimeType:  file.MimeType,
		Elorating: 400.0,
	}

	// insert the links to the database
	_, err = picCollection.InsertOne(context.Background(), filemetadata)
	if err != nil {
		return err
	}

	return nil

}
