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

func GetPair() ([]Car, error) {

	client, picCollection, err := PicDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	var car []Car
	cur, err := picCollection.Aggregate(context.Background(), bson.A{bson.M{"$sample": bson.M{"size": 2}}})
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.Background(), &car); err != nil {
		return nil, err
	}

	return car, nil
}

func GetCarID(carID string) (*Car, error) {
	client, picCollection, err := PicDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	oid, err := primitive.ObjectIDFromHex(carID)
	if err != nil {
		return nil, err
	}

	var car Car
	err = picCollection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&car)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

func UpdateEloRating(winnerCarID, loserCarID string, newWinnerR, newLoserR float64) error {
	client, picCollection, err := PicDB()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	winnerOID, err := primitive.ObjectIDFromHex(winnerCarID)
	if err != nil {
		return err
	}

	loserOID, err := primitive.ObjectIDFromHex(loserCarID)
	if err != nil {
		return err
	}

	_, err = picCollection.UpdateOne(context.Background(), bson.M{"_id": winnerOID}, bson.M{"$set": bson.M{"rating": newWinnerR}})
	if err != nil {
		return err
	}

	_, err = picCollection.UpdateOne(context.Background(), bson.M{"_id": loserOID}, bson.M{"$set": bson.M{"rating": newLoserR}})
	if err != nil {
		return err
	}

	return nil
}

func GetAllPics() ([]Car, error) {

	client, picCollection, err := PicDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	var cars []Car
	cur, err := picCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.Background(), &cars); err != nil {
		return nil, err
	}

	return cars, nil

}
