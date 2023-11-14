package database

import (
	"context"
	"log"
	"time"

	"github.com/maksimulitin/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "127.0.0.1:27017"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}

}

func (db *DB) CreateJobListing(jobInf model.CreateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, concel := context.WithTimeout(context.Background(), 30*time.Second)
	defer concel()

	inserg, err := jobCollec.InsertOne(ctx, bson.M{"title": jobInf.Title, "description": jobInf.Description, "company": jobInf.Company, "url": jobInf.URL})

	if err != nil {
		log.Fatal(err)
	}

	insertedId := inserg.InsertedID.(primitive.ObjectID).Hex()

	returnJobListing := model.JobListing{ID: insertedId, Title: jobInf.Title, Description: jobInf.Description, Company: jobInf.Company, URL: jobInf.URL}
	return &returnJobListing
}

func (db *DB) UpdateJobListing(jobId string, jobInf model.UpdateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("grapql-job-board").Collection("jobs")
	ctx, concel := context.WithTimeout(context.Background(), 30*time.Second)
	defer concel()
	updateJobInfo := bson.M{}

	if jobInf.Title != nil {
		updateJobInfo["title"] = jobInf.Title

	}
	if jobInf.Description != nil {
		updateJobInfo["description"] = jobInf.Description

	}
	if jobInf.URL != nil {
		updateJobInfo["url"] = jobInf.URL
	}

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	result := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var jobListing model.JobListing

	if err := result.Decode(&jobListing); err != nil {
		log.Fatal(err)
	}

	return &jobListing
}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, concel := context.WithTimeout(context.Background(), 30*time.Second)
	defer concel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	_, err := jobCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	return &model.DeleteJobResponse{DeletedJobID: jobId}
}

func (db *DB) GetJob(id string) *model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, concel := context.WithTimeout(context.Background(), 30*time.Second)
	defer concel()
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var jobListing model.JobListing

	if err := jobCollec.FindOne(ctx, filter).Decode(jobListing); err != nil {
		log.Fatal(err)
	}
	return &jobListing
}

func (db *DB) GetJobs() []*model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, concel := context.WithTimeout(context.Background(), 30*time.Second)
	defer concel()
	var jobListings []*model.JobListing
	cursor, err := jobCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &jobListings); err != nil {
		log.Fatal(err)
	}
	return jobListings
}
