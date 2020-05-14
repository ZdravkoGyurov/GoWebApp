package models

import (
	"context"
	"errors"
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

// InsertNote creates note in databse
func InsertNote(collection *mongo.Collection, n Note) error {
	insertResult, err := collection.InsertOne(context.TODO(), n)
	if err != nil {
		return errors.New("An error occured while inserting note")
	}

	fmt.Println("Inserted a single note: ", insertResult.InsertedID)
	return nil
}

// FindNoteByTitle reads a note from databse by name
func FindNoteByTitle(collection *mongo.Collection, title string) (*Note, error) {
	var note *Note
	filter := bson.M{
		"title": bson.M{
			"$eq": title,
		},
	}

	if err := collection.FindOne(context.TODO(), filter).Decode(&note); err != nil {
		return nil, errors.New("An error occured while finding note by title: " + title)
	}

	return note, nil
}

// FindNotes reads all notes from database
func FindNotes(collection *mongo.Collection) ([]Note, error) {
	findOptions := options.Find()

	var notes []Note

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, errors.New("An error occured while finding notes")
	}

	for cur.Next(context.TODO()) {
		var elem Note

		if err := cur.Decode(&elem); err != nil {
			return nil, errors.New("An error occured while decoding a note")
		}

		notes = append(notes, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal("An error occured while reading notes from MongoDB")
	}

	cur.Close(context.TODO())

	return notes, nil
}

// UpdateNoteByTitle updates a note in the database
func UpdateNoteByTitle(collection *mongo.Collection, title string, n Note) error {
	filter := bson.M{
		"title": bson.M{
			"$eq": title,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"title":       n.Title,
			"body":        n.Body,
			"modified_on": time.Now(),
		},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.New("An error occured while updating note with title: " + title)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}

// DeleteNoteByTitle deletes a note from the database
func DeleteNoteByTitle(collection *mongo.Collection, title string) error {
	filter := bson.M{
		"title": bson.M{
			"$eq": title,
		},
	}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil || deleteResult.DeletedCount == 0 {
		return errors.New("An error occured while deleting note with title: " + title)
	}

	fmt.Println("Deleted: ", deleteResult)
	return nil
}
