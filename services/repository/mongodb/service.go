// Package mongodb implements MongoDB repository services
package mongodb

import (
	"context"
	"fmt"

	"github.com/decentralized-cloud/project/models"
	"github.com/decentralized-cloud/project/services/configuration"
	"github.com/decentralized-cloud/project/services/repository"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbRepositoryService struct {
	connectionString string
	databaseName     string
}

const collectionName string = "project"

// NewMongodbRepositoryService creates new instance of the mongodbRepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewMongodbRepositoryService(
	configurationService configuration.ConfigurationContract) (repository.RepositoryContract, error) {
	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	connectionString, err := configurationService.GetDatabaseConnectionString()
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to get connection string to mongodb", err)
	}

	databaseName, err := configurationService.GetDatabaseName()
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to get the database name", err)
	}

	return &mongodbRepositoryService{
		connectionString: connectionString,
		databaseName:     databaseName,
	}, nil
}

// CreateProject creates a new project.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new project
// Returns either the result of creating new project or error if something goes wrong.
func (service *mongodbRepositoryService) CreateProject(
	ctx context.Context,
	request *repository.CreateProjectRequest) (*repository.CreateProjectResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	insertResult, err := collection.InsertOne(ctx, request.Project)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Insert project failed.", err)
	}

	projectID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	return &repository.CreateProjectResponse{
		ProjectID: projectID,
		Project:   request.Project,
		Cursor:    projectID,
	}, nil
}

// ReadProject read an existing project
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing project
// Returns either the result of reading an existing project or error if something goes wrong.
func (service *mongodbRepositoryService) ReadProject(
	ctx context.Context,
	request *repository.ReadProjectRequest) (*repository.ReadProjectResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	ObjectID, _ := primitive.ObjectIDFromHex(request.ProjectID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}
	var project models.Project

	err = collection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		return nil, repository.NewProjectNotFoundError(request.ProjectID)
	}

	return &repository.ReadProjectResponse{
		Project: project,
	}, nil
}

// UpdateProject update an existing project
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing project
// Returns either the result of updateing an existing project or error if something goes wrong.
func (service *mongodbRepositoryService) UpdateProject(
	ctx context.Context,
	request *repository.UpdateProjectRequest) (*repository.UpdateProjectResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	ObjectID, _ := primitive.ObjectIDFromHex(request.ProjectID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}

	newProject := bson.M{"$set": bson.M{"name": request.Project.Name}}
	response, err := collection.UpdateOne(ctx, filter, newProject)

	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Update project failed.", err)
	}

	if response.MatchedCount == 0 {
		return nil, repository.NewProjectNotFoundError(request.ProjectID)
	}

	return &repository.UpdateProjectResponse{
		Project: request.Project,
		Cursor:  request.ProjectID,
	}, nil
}

// DeleteProject delete an existing project
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing project
// Returns either the result of deleting an existing project or error if something goes wrong.
func (service *mongodbRepositoryService) DeleteProject(
	ctx context.Context,
	request *repository.DeleteProjectRequest) (*repository.DeleteProjectResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	ObjectID, _ := primitive.ObjectIDFromHex(request.ProjectID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}
	response, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Delete project failed.", err)
	}

	if response.DeletedCount == 0 {
		return nil, repository.NewProjectNotFoundError(request.ProjectID)
	}

	return &repository.DeleteProjectResponse{}, nil
}

// Search returns the list of projects that matched the criteria
// ctx: Mandatory The reference to the context
// request: Mandatory. The request contains the search criteria
// Returns the list of projects that matched the criteria
func (service *mongodbRepositoryService) Search(
	ctx context.Context,
	request *repository.SearchRequest) (*repository.SearchResponse, error) {
	response := &repository.SearchResponse{
		HasPreviousPage: false,
		HasNextPage:     false,
	}

	ids := []primitive.ObjectID{}
	for _, projectID := range request.ProjectIDs {
		objectID, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the projectID: %s.", projectID), err)
		}

		ids = append(ids, objectID)
	}

	filter := bson.M{}
	if len(request.ProjectIDs) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}

	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	response.TotalCount, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to retrieve the number of projects that match the filter criteria", err)
	}

	if response.TotalCount == 0 {
		// No tennat matched the filter criteria
		return response, nil
	}

	if request.Pagination.After != nil {
		after := *request.Pagination.After
		objectID, err := primitive.ObjectIDFromHex(after)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the After: %s.", after), err)
		}

		if len(filter) > 0 {
			filter["$and"] = []interface{}{
				bson.M{"_id": bson.M{"$gt": objectID}},
			}
		} else {
			filter["_id"] = bson.M{"$gt": objectID}
		}
	}

	if request.Pagination.Before != nil {
		before := *request.Pagination.Before
		objectID, err := primitive.ObjectIDFromHex(before)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the Before: %s.", before), err)
		}

		if len(filter) > 0 {
			filter["$and"] = []interface{}{
				bson.M{"_id": bson.M{"$lt": objectID}},
			}
		} else {
			filter["_id"] = bson.M{"$lt": objectID}
		}
	}

	findOptions := options.Find()

	if request.Pagination.First != nil {
		findOptions.SetLimit(int64(*request.Pagination.First))
	}

	if request.Pagination.Last != nil {
		findOptions.SetLimit(int64(*request.Pagination.Last))
	}

	if len(request.SortingOptions) > 0 {
		var sortOptionPairs bson.D

		for _, sortingOption := range request.SortingOptions {
			direction := 1
			if sortingOption.Direction == common.Descending {
				direction = -1
			}

			sortOptionPairs = append(
				sortOptionPairs,
				bson.E{
					Key:   sortingOption.Name,
					Value: direction,
				})
		}

		findOptions.SetSort(sortOptionPairs)
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to call the Find function on the collection.", err)
	}

	projects := []models.ProjectWithCursor{}
	for cursor.Next(ctx) {
		var project models.Project
		//TODO : below line need to be removed, if we pass 'ShowRecordID' in findOption, ObjectID will be available
		var projectBson bson.M

		err := cursor.Decode(&project)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Failed to decode the project", err)
		}

		err = cursor.Decode(&projectBson)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Could not load the data.", err)
		}

		projectID := projectBson["_id"].(primitive.ObjectID).Hex()
		projectWithCursor := models.ProjectWithCursor{
			ProjectID: projectID,
			Project:   project,
			Cursor:    projectID,
		}

		projects = append(projects, projectWithCursor)
	}

	response.Projects = projects
	if (request.Pagination.After != nil && request.Pagination.First != nil && int64(*request.Pagination.First) < response.TotalCount) ||
		(request.Pagination.Before != nil && request.Pagination.Last != nil && int64(*request.Pagination.Last) < response.TotalCount) {
		response.HasNextPage = true
		response.HasPreviousPage = true
	} else if request.Pagination.After == nil && request.Pagination.First != nil && int64(*request.Pagination.First) < response.TotalCount {
		response.HasNextPage = true
		response.HasPreviousPage = false
	} else if request.Pagination.Before == nil && request.Pagination.Last != nil && int64(*request.Pagination.Last) < response.TotalCount {
		response.HasNextPage = false
		response.HasPreviousPage = true
	}

	return response, nil
}

func (service *mongodbRepositoryService) createClientAndCollection(ctx context.Context) (*mongo.Client, *mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(service.connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	return client, client.Database(service.databaseName).Collection(collectionName), nil
}

func disconnect(ctx context.Context, client *mongo.Client) {
	_ = client.Disconnect(ctx)
}
