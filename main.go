package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	config "github.com/stevebargelt/convertLitterTrips/config"
)

type LitterboxTrip struct {
	CatID                string
	CatName              string
	CatProbability       float64 `firestore:"probability"`
	Direction            string
	DirectionProbability float64
	Photo                string
	TimeStamp            time.Time `firestore:"timestamp"`
}

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	viper.SetDefault("HTTP_RETRY_COUNT", 20)

	var configuration config.Configuration
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	ctx := context.Background()
	sa := option.WithCredentialsFile(configuration.FirebaseCredentials)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	// var litterboxTrips []LitterboxTrip
	var litterboxTrip LitterboxTrip
	fmt.Printf("Source Collection: %s\n", configuration.FirestoreCollectionSource)
	fmt.Printf("Cat Name: %s\n", configuration.CatID)
	fmt.Printf("Cat Name: %s\n", configuration.CatName)
	iter := client.Collection(configuration.FirestoreCollectionSource).Doc(configuration.CatName).Collection("LitterTrips").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc.DataTo(&litterboxTrip)
		// CatID and CatName were not in the original data - adding here to flatten data
		litterboxTrip.CatID = configuration.CatID
		litterboxTrip.CatName = configuration.CatName
		fmt.Printf("Document data: %#v\n", litterboxTrip)
		fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
			litterboxTrip.TimeStamp.Year(), litterboxTrip.TimeStamp.Month(), litterboxTrip.TimeStamp.Day(),
			litterboxTrip.TimeStamp.Hour(), litterboxTrip.TimeStamp.Minute(), litterboxTrip.TimeStamp.Second())

		addLitterBoxTripToFirestore(litterboxTrip, configuration.FirebaseCredentials, configuration.FirestoreCollectionDestination)
	}

}

func addLitterBoxTripToFirestore(litterBoxTrip LitterboxTrip, firebaseCredentials string, firestoreCollection string) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(firebaseCredentials)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	_, _, err = client.Collection(firestoreCollection).Add(ctx, map[string]interface{}{
		"CatName":              litterBoxTrip.CatName,
		"Probability":          litterBoxTrip.CatProbability,
		"Direction":            litterBoxTrip.Direction,
		"DirectionProbability": litterBoxTrip.DirectionProbability,
		"Photo":                litterBoxTrip.Photo, // right now this is the local name. Could be the URL to the photo in Cloud Storage.
		"timestamp":            litterBoxTrip.TimeStamp,
	})
	if err != nil {
		log.Fatalf("Failed adding litterbox trip: %v", err)
	}
}
