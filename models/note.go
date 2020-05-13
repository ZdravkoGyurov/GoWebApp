package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Note with title and body
type Note struct {
	Title      string    `bson:"title"`
	Body       string    `bson:"body"`
	ModifiedOn time.Time `bson:"modified_on"`
}

// NewNote creates note
func NewNote(title, body string, modifiedOn time.Time) Note {
	return Note{title, body, modifiedOn}
}

// ToBson converts note to bson object
func (n Note) ToBson() bson.D {
	return bson.D{
		{Key: "title", Value: n.Title},
		{Key: "body", Value: n.Body},
		{Key: "modified_on", Value: n.ModifiedOn},
	}
}

// InsertNote creates note in databse
func InsertNote(collection *mongo.Collection, n Note) {
	insertResult, err := collection.InsertOne(context.TODO(), n.ToBson())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single note: ", insertResult.InsertedID)
}

// FindNoteByTitle reads a note from databse by name
func FindNoteByTitle(collection *mongo.Collection, title string) Note {
	var note Note
	filter := bson.M{
		"title": bson.M{
			"$eq": title,
		},
	}

	if err := collection.FindOne(context.TODO(), filter).Decode(&note); err != nil {
		log.Fatal(err)
	}

	return note
}

// FindNotes reads all notes from database
func FindNotes(collection *mongo.Collection) []Note {
	findOptions := options.Find()

	var notes []Note

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Note

		if err := cur.Decode(&elem); err != nil {
			log.Fatal(err)
		}

		notes = append(notes, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return notes
}

// UpdateNoteByTitle updates a note in the database
func UpdateNoteByTitle(collection *mongo.Collection, title string, n Note) {
	filter := bson.M{
		"title": bson.M{
			"$eq": title,
		},
	}

	update := bson.M{"$set": bson.M{"title": n.Title, "body": n.Body, "modified_on": time.Now()}}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

// DeleteNoteByTitle deletes a note from the database
func DeleteNoteByTitle(collection *mongo.Collection, title string) {
	filter := bson.M{
		"title": bson.M{
			"$eq": title,
		},
	}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted: ", deleteResult)
}
