// Package mongodb implements MongoDB repository services
package mongodb

import (
	"context"
	"fmt"

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

// CreateTenant creates a new tenant.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *mongodbRepositoryService) CreateTenant(
	ctx context.Context,
	request *repository.CreateTenantRequest) (*repository.CreateTenantResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	insertResult, err := collection.InsertOne(ctx, request.Tenant)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Insert tenant failed.", err)
	}

	tenantID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	return &repository.CreateTenantResponse{
		TenantID: tenantID,
		Tenant:   request.Tenant,
		Cursor:   tenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) ReadTenant(
	ctx context.Context,
	request *repository.ReadTenantRequest) (*repository.ReadTenantResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

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
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

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
		Cursor: request.TenantID,
	}, nil
}

// DeleteTenant delete an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) DeleteTenant(
	ctx context.Context,
	request *repository.DeleteTenantRequest) (*repository.DeleteTenantResponse, error) {
	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

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

	ids := []primitive.ObjectID{}
	for _, tenantID := range request.TenantIDs {
		objectID, err := primitive.ObjectIDFromHex(tenantID)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError(fmt.Sprintf("Failed to decode the tenantID: %s.", tenantID), err)
		}

		ids = append(ids, objectID)
	}

	filter := bson.M{}
	if len(request.TenantIDs) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}

	client, collection, err := service.createClientAndCollection(ctx)
	if err != nil {
		return nil, err
	}

	defer disconnect(ctx, client)

	response.TotalCount, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Failed to retrieve the number of tenants that match the filter criteria", err)
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

	tenants := []models.TenantWithCursor{}
	for cursor.Next(ctx) {
		var tenant models.Tenant
		//TODO : below line need to be removed, if we pass 'ShowRecordID' in findOption, ObjectID will be available
		var tenantBson bson.M

		err := cursor.Decode(&tenant)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Failed to decode the tenant", err)
		}

		err = cursor.Decode(&tenantBson)
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Could not load the data.", err)
		}

		tenantID := tenantBson["_id"].(primitive.ObjectID).Hex()
		tenantWithCursor := models.TenantWithCursor{
			TenantID: tenantID,
			Tenant:   tenant,
			Cursor:   tenantID,
		}

		tenants = append(tenants, tenantWithCursor)
	}

	response.Tenants = tenants
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
