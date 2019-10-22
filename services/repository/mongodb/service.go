// Package mongodb implements MongoDB repository services
package mongodb

import (
	"context"
	"fmt"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/configuration"
	"github.com/decentralized-cloud/tenant/services/repository"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type mongodbRepositoryService struct {
	dbConnectionString string
}

// NewMongodbRepositoryService creates new instance of the RepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewMongodbRepositoryService(configurationService configuration.ConfigurationContract) (repository.RepositoryContract, error) {

	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	dbConnectionString, err := configurationService.GetDbConnectionString()
	if err != nil {
		return nil, err
	}

	return &mongodbRepositoryService{
		dbConnectionString: dbConnectionString,
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) ReadTenant(
	ctx context.Context,
	request *repository.ReadTenantRequest) (*repository.ReadTenantResponse, error) {	

		// Todo : move line 47-57 to a method
		clientOptions := options.Client().ApplyURI(service.dbConnectionString)
		client, err := mongo.Connect(context.TODO(), clientOptions)		
		if err != nil {
			return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
		}

		//Todo : make sure if ping is necessary
		err = client.Ping(context.TODO(), nil)
		if err != nil {
    		return nil, repository.NewUnknownErrorWithError("database is not accessible.", err)
		}

		//Todo : Do we need to have GetDbName in configuration service ? for now it is just hard coded
		collection := client.Database("tenantdb").Collection("tenants")

		filter := bson.D{{"TenantID", request.TenantID}}
		var tenant models.Tenant

		err = collection.FindOne(context.TODO(), filter).Decode(&tenant)
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

	clientOptions := options.Client().ApplyURI(service.dbConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)		
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}

	//Todo : make sure if ping is necessary
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("database is not accessible.", err)
	}

 	//Todo : Do we need to have GetDbName in configuration service ? for now it is just hard coded
	collection := client.Database("tenantdb").Collection("tenants")
	insertResult, err := collection.InsertOne(context.TODO(), request.Tenant)

	if err != nil {
    	return nil, repository.NewUnknownErrorWithError("Insert tenant failed.", err)
	}

	return &repository.CreateTenantResponse{
		TenantID: insertResult.InsertedID.(string),
	}, nil
}


// UpdateTenant update an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *mongodbRepositoryService) UpdateTenant(
	ctx context.Context,
	request *repository.UpdateTenantRequest) (*repository.UpdateTenantResponse, error) {

	clientOptions := options.Client().ApplyURI(service.dbConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)		
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}
	
	//Todo : make sure if ping is necessary
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("database is not accessible.", err)
	}
	
	 //Todo : Do we need to have GetDbName in configuration service ? for now it is just hard coded
	collection := client.Database("tenantdb").Collection("tenants")
	filter := bson.D{{"TenantID", request.TenantID}}

	newTenant := bson.D{
		{"$set", request.Tenant},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, newTenant)
	if err != nil {
    	return nil, repository.NewUnknownErrorWithError("Update tenant failed.", err)
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

	clientOptions := options.Client().ApplyURI(service.dbConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)		
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Could not connect to mongodb database.", err)
	}
		
	//Todo : make sure if ping is necessary
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("database is not accessible.", err)
	}
		
	 //Todo : Do we need to have GetDbName in configuration service ? for now it is just hard coded
	collection := client.Database("tenantdb").Collection("tenants")
	filter := bson.D{{"TenantID", request.TenantID}}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, repository.NewUnknownErrorWithError("Delete tenant failed.", err)
	}
	
	return &repository.DeleteTenantResponse{}, nil
}