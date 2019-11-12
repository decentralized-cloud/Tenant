// Package mongodb implements MongoDB repository services
package mongodb

import (
	"context"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/configuration"
	"github.com/decentralized-cloud/tenant/services/repository"
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

const collectionName string = "tenant"

// NewMongodbRepositoryService creates new instance of the RepositoryService, setting up all dependencies and returns the instance
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

// CreateTenant creates a new tenant.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *mongodbRepositoryService) CreateTenant(
	ctx context.Context,
	request *repository.CreateTenantRequest) (*repository.CreateTenantResponse, error) {
	clientOptions := options.Client().ApplyURI(service.connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	defer disconnect(ctx, client)

	collection := client.Database(service.databaseName).Collection(collectionName)
	insertResult, err := collection.InsertOne(ctx, request.Tenant)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Insert tenant failed.", err)
	}

	return &repository.CreateTenantResponse{
		TenantID: insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Tenant:   request.Tenant,
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) ReadTenant(
	ctx context.Context,
	request *repository.ReadTenantRequest) (*repository.ReadTenantResponse, error) {

	clientOptions := options.Client().ApplyURI(service.connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	defer disconnect(ctx, client)

	collection := client.Database(service.databaseName).Collection(collectionName)
	ObjectID, _ := primitive.ObjectIDFromHex(request.TenantID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}
	var tenant models.Tenant

	err = collection.FindOne(ctx, filter).Decode(&tenant)
	if err != nil {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	return &repository.ReadTenantResponse{
		Tenant: tenant,
	}, nil
}

// UpdateTenant update an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) UpdateTenant(
	ctx context.Context,
	request *repository.UpdateTenantRequest) (*repository.UpdateTenantResponse, error) {
	clientOptions := options.Client().ApplyURI(service.connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	defer disconnect(ctx, client)

	collection := client.Database(service.databaseName).Collection(collectionName)
	ObjectID, _ := primitive.ObjectIDFromHex(request.TenantID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}

	newTenant := bson.M{"$set": bson.M{"name": request.Tenant.Name}}
	response, err := collection.UpdateOne(ctx, filter, newTenant)

	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Update tenant failed.", err)
	}

	if response.MatchedCount == 0 {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	return &repository.UpdateTenantResponse{
		Tenant: request.Tenant,
	}, nil
}

// DeleteTenant delete an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) DeleteTenant(
	ctx context.Context,
	request *repository.DeleteTenantRequest) (*repository.DeleteTenantResponse, error) {
	clientOptions := options.Client().ApplyURI(service.connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	defer disconnect(ctx, client)

	collection := client.Database(service.databaseName).Collection(collectionName)
	ObjectID, _ := primitive.ObjectIDFromHex(request.TenantID)
	filter := bson.D{{Key: "_id", Value: ObjectID}}
	response, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Delete tenant failed.", err)
	}

	if response.DeletedCount == 0 {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	return &repository.DeleteTenantResponse{}, nil
}

// Search returns the list of tenants that matched the criteria
// ctx: Mandatory The reference to the context
// request: Mandatory. The request contains the search criteria
// Returns the list of tenants that matched the criteria
func (service *mongodbRepositoryService) Search(
	ctx context.Context,
	request *repository.SearchRequest) (*repository.SearchResponse, error) {

	response := &repository.SearchResponse{
		HasPreviousPage: false,
		HasNextPage:     false,
	}

	clientOptions := options.Client().ApplyURI(service.connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	defer disconnect(ctx, client)

	collection := client.Database(service.databaseName).Collection(collectionName)

	filter := bson.M{}
	ids := make([]primitive.ObjectID, len(request.TenantIDs))
	if len(request.TenantIDs) > 0 {
		for i := range request.TenantIDs {
			ObjectID, _ := primitive.ObjectIDFromHex(request.TenantIDs[i])
			ids[i] = ObjectID
		}
		filter = bson.M{"_id": bson.M{"$in": ids}}
	}

	findOptions := options.Find()
	if request.Pagination.After != "" {
		ObjectID, _ := primitive.ObjectIDFromHex(request.Pagination.After)
		filter = bson.M{
			"_id": bson.M{"$in": ids},
			"$and": []interface{}{
				bson.M{"_id": bson.M{"$gt": ObjectID}},
			},
		}
	}

	if request.Pagination.First > 0 {
		findOptions.SetLimit(int64(request.Pagination.First))
	}

	if request.Pagination.Before != "" {
		ObjectID, _ := primitive.ObjectIDFromHex(request.Pagination.Before)

		filter = bson.M{
			"_id": bson.M{"$in": ids},
			"$and": []interface{}{
				bson.M{"_id": bson.M{"$lt": ObjectID}},
			},
		}
	}

	if request.Pagination.Last > 0 {
		findOptions.SetLimit(int64(request.Pagination.Last))
	}

	if len(request.SortingOptions) > 0 {

		var sortOptionPairs bson.D
		var direction int
		for i := range request.SortingOptions {
			fieldName := request.SortingOptions[i].Name

			switch request.SortingOptions[i].Direction {
			case common.Ascending:
				direction = 1
			case common.Descending:
				direction = -1
			}
			sortOptionPairs = append(sortOptionPairs, bson.E{fieldName, direction})
		}

		findOptions.SetSort(sortOptionPairs)
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not load the data.", err)
	}

	var tenants []models.TenantWithCursor

	for cursor.Next(ctx) {

		var tenant models.Tenant
		//Todo : below line need to be removed, if we pass 'ShowRecordID' in findOption, ObjectID will be available
		var tenantBson bson.M
		err := cursor.Decode(&tenant)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Could not load the data.", err)
		}

		err = cursor.Decode(&tenantBson)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Could not load the data.", err)
		}

		tenantWithCursor := models.TenantWithCursor{
			TenantID: tenantBson["_id"].(primitive.ObjectID).Hex(),
			Tenant:   tenant,
		}

		tenants = append(tenants, tenantWithCursor)
	}

	response.Tenants = tenants

	//TODO : find a way to populate below properties in response
	//HasPreviousPage
	//HasNextPage

	return response, nil
}

func disconnect(ctx context.Context, client *mongo.Client) {
	_ = client.Disconnect(ctx)
}
