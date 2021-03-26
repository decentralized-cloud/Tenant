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

type project struct {
	UserEmail string `bson:"userEmail" json:"userEmail"`
	Name      string `bson:"name" json:"name"`
}

type mongodbRepositoryService struct {
	connectionString       string
	databaseName           string
	databaseCollectionName string
}

// NewMongodbRepositoryService creates new instance of the mongodbRepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewMongodbRepositoryService(
	configurationService configuration.ConfigurationContract) (repository.RepositoryContract, error) {
	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	connectionString, err := configurationService.GetDatabaseConnectionString()
	if err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Failed to get connection string to mongodb", err)
	}

	databaseName, err := configurationService.GetDatabaseName()
	if err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Failed to get the database name", err)
	}

	databaseCollectionName, err := configurationService.GetDatabaseCollectionName()
	if err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Failed to get the database collection name", err)
	}

	return &mongodbRepositoryService{
		connectionString:       connectionString,
		databaseName:           databaseName,
		databaseCollectionName: databaseCollectionName,
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

	insertResult, err := collection.InsertOne(ctx, mapToInternalProject(request.UserEmail, request.Project))
	if err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Insert project failed.", err)
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
	filter := bson.D{{Key: "_id", Value: ObjectID}, {Key: "userEmail", Value: request.UserEmail}}
	var project project

	err = collection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		return nil, commonErrors.NewNotFoundError()
	}

	return &repository.ReadProjectResponse{
		Project: mapFromInternalProject(project),
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
	filter := bson.D{{Key: "_id", Value: ObjectID}, {Key: "userEmail", Value: request.UserEmail}}

	newProject := bson.M{"$set": bson.M{"name": request.Project.Name}}
	response, err := collection.UpdateOne(ctx, filter, newProject)

	if err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Update project failed.", err)
	}

	if response.MatchedCount == 0 {
		return nil, commonErrors.NewNotFoundError()
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
	filter := bson.D{{Key: "_id", Value: ObjectID}, {Key: "userEmail", Value: request.UserEmail}}
	response, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Delete project failed.", err)
	}

	if response.DeletedCount == 0 {
		return nil, commonErrors.NewNotFoundError()
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
			return nil, commonErrors.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the projectID: %s.", projectID), err)
		}

		ids = append(ids, objectID)
	}

	filter := bson.M{}
	if len(request.ProjectIDs) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}

	filter["$and"] = []interface{}{
		bson.M{"userEmail": bson.M{"$eq": request.UserEmail}},
	}

	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	response.TotalCount, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Failed to retrieve the number of projects that match the filter criteria", err)
	}

	if response.TotalCount == 0 {
		// No tennat matched the filter criteria
		return response, nil
	}

	if request.Pagination.After != nil {
		after := *request.Pagination.After
		objectID, err := primitive.ObjectIDFromHex(after)
		if err != nil {
			return nil, commonErrors.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the After: %s.", after), err)
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
			return nil, commonErrors.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the Before: %s.", before), err)
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
		return nil, commonErrors.NewUnknownErrorWithError("Failed to call the Find function on the collection.", err)
	}

	projects := []models.ProjectWithCursor{}
	for cursor.Next(ctx) {
		var project project
		var projectBson bson.M

		err := cursor.Decode(&project)
		if err != nil {
			return nil, commonErrors.NewUnknownErrorWithError("Failed to decode the project", err)
		}

		err = cursor.Decode(&projectBson)
		if err != nil {
			return nil, commonErrors.NewUnknownErrorWithError("Could not load the data.", err)
		}

		projectID := projectBson["_id"].(primitive.ObjectID).Hex()
		projectWithCursor := models.ProjectWithCursor{
			ProjectID: projectID,
			Project:   mapFromInternalProject(project),
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
		return nil, nil, commonErrors.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	return client, client.Database(service.databaseName).Collection(service.databaseCollectionName), nil
}

func disconnect(ctx context.Context, client *mongo.Client) {
	_ = client.Disconnect(ctx)
}

func mapToInternalProject(email string, from models.Project) project {
	return project{
		UserEmail: email,
		Name:      from.Name,
	}
}

func mapFromInternalProject(from project) models.Project {
	return models.Project{
		Name: from.Name,
	}
}
