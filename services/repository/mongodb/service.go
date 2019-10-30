// Package mongodb implements MongoDB repository services
package mongodb

import (
	"context"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/configuration"
	"github.com/decentralized-cloud/tenant/services/repository"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbRepositoryService struct {
	configurationService configuration.ConfigurationContract
	collectionName       string
}

// NewMongodbRepositoryService creates new instance of the RepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewMongodbRepositoryService(
	configurationService configuration.ConfigurationContract) (repository.RepositoryContract, error) {

	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	return &mongodbRepositoryService{
		configurationService: configurationService,
		collectionName:       "tenants",
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) ReadTenant(
	ctx context.Context,
	request *repository.ReadTenantRequest) (*repository.ReadTenantResponse, error) {

	dbConnectionString, err := service.configurationService.GetDbConnectionString()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetDbConnectionString", "Database connection String is required.")
	}

	clientOptions := options.Client().ApplyURI(dbConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	tenantDbName, err := service.configurationService.GetTenantDbName()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetTenantDbName", "Database name is required.")
	}

	collection := client.Database(tenantDbName).Collection(service.collectionName)
	ObjectID, _ := primitive.ObjectIDFromHex(request.TenantID)
	filter := bson.D{{"_id", ObjectID}}
	var tenant models.Tenant

	err = collection.FindOne(ctx, filter).Decode(&tenant)
	if err != nil {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	return &repository.ReadTenantResponse{Tenant: tenant}, nil
}

// CreateTenant creates a new tenant.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *mongodbRepositoryService) CreateTenant(
	ctx context.Context,
	request *repository.CreateTenantRequest) (*repository.CreateTenantResponse, error) {

	dbConnectionString, err := service.configurationService.GetDbConnectionString()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetDbConnectionString", "Database connection String is required.")
	}

	clientOptions := options.Client().ApplyURI(dbConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	tenantDbName, err := service.configurationService.GetTenantDbName()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetTenantDbName", "Database name is required.")
	}

	collection := client.Database(tenantDbName).Collection(service.collectionName)
	insertResult, err := collection.InsertOne(ctx, request.Tenant)

	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Insert tenant failed.", err)
	}

	return &repository.CreateTenantResponse{
		TenantID: insertResult.InsertedID.(primitive.ObjectID).Hex(),
	}, nil
}

// UpdateTenant update an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) UpdateTenant(
	ctx context.Context,
	request *repository.UpdateTenantRequest) (*repository.UpdateTenantResponse, error) {

	dbConnectionString, err := service.configurationService.GetDbConnectionString()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetDbConnectionString", "Database connection String is required.")
	}

	clientOptions := options.Client().ApplyURI(dbConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	tenantDbName, err := service.configurationService.GetTenantDbName()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetTenantDbName", "Database name is required.")
	}

	collection := client.Database(tenantDbName).Collection(service.collectionName)
	ObjectID, _ := primitive.ObjectIDFromHex(request.TenantID)
	filter := bson.D{{"_id", ObjectID}}

	newTenant := bson.M{"$set": bson.M{"name": request.Tenant.Name}}
	response, err := collection.UpdateOne(context.TODO(), filter, newTenant)

	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Update tenant failed.", err)
	}

	if response.MatchedCount == 0 {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	return &repository.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) DeleteTenant(
	ctx context.Context,
	request *repository.DeleteTenantRequest) (*repository.DeleteTenantResponse, error) {

	dbConnectionString, err := service.configurationService.GetDbConnectionString()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetDbConnectionString", "Database connection String is required.")
	}

	clientOptions := options.Client().ApplyURI(dbConnectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	tenantDbName, err := service.configurationService.GetTenantDbName()
	if err != nil {
		return nil, commonErrors.NewArgumentNilError("GetTenantDbName", "Database name is required.")
	}

	collection := client.Database(tenantDbName).Collection(service.collectionName)

	ObjectID, _ := primitive.ObjectIDFromHex(request.TenantID)
	filter := bson.D{{"_id", ObjectID}}
	response, err := collection.DeleteOne(context.TODO(), filter)

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

	response := &repository.SearchResponse{}

	return response, nil
}
