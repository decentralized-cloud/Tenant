// Package repository implements different repository services required by the tenant service
package mongodb_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	configurationMock "github.com/decentralized-cloud/tenant/services/configuration/mock"
	"github.com/decentralized-cloud/tenant/services/repository"
	"github.com/decentralized-cloud/tenant/services/repository/mongodb"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMongodbRepositoryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Repository Service Tests")
}

var _ = Describe("Mongodb Repository Service Tests", func() {
	var (
		mockCtrl      *gomock.Controller
		sut           repository.RepositoryContract
		ctx           context.Context
		createRequest repository.CreateTenantRequest
	)

	BeforeEach(func() {
		connectionString := os.Getenv("DATABASE_CONNECTION_STRING")
		if strings.Trim(connectionString, " ") == "" {
			connectionString = "mongodb://mongodb:27017"
		}

		mockCtrl = gomock.NewController(GinkgoT())
		mockConfigurationService := configurationMock.NewMockConfigurationContract(mockCtrl)
		mockConfigurationService.
			EXPECT().
			GetDatabaseConnectionString().
			Return(connectionString, nil)

		mockConfigurationService.
			EXPECT().
			GetDatabaseName().
			Return("tenant", nil)

		sut, _ = mongodb.NewMongodbRepositoryService(mockConfigurationService)
		ctx = context.Background()
		createRequest = repository.CreateTenantRequest{
			Tenant: models.Tenant{
				Name: cuid.New(),
			}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate RepositoryService", func() {
		When("all dependencies are resolved and NewRepositoryService is called", func() {
			It("should instantiate the new RepositoryService", func() {
				mockConfigurationService := configurationMock.NewMockConfigurationContract(mockCtrl)
				mockConfigurationService.
					EXPECT().
					GetDatabaseConnectionString().
					Return(cuid.New(), nil)

				mockConfigurationService.
					EXPECT().
					GetDatabaseName().
					Return(cuid.New(), nil)
				service, err := mongodb.NewMongodbRepositoryService(mockConfigurationService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Context("user going to create a new tenant", func() {
		When("create tenant is called", func() {
			It("should create the new tenant", func() {
				response, err := sut.CreateTenant(ctx, &createRequest)
				Ω(err).Should(BeNil())
				Ω(response.TenantID).ShouldNot(BeNil())
			})
		})
	})

	Context("tenant already exists", func() {
		var (
			tenantID string
		)

		BeforeEach(func() {
			response, _ := sut.CreateTenant(ctx, &createRequest)
			tenantID = response.TenantID
		})

		When("user reads a tenant by Id", func() {
			It("should return a tenant", func() {
				response, err := sut.ReadTenant(ctx, &repository.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				Ω(response.Tenant.Name).Should(Equal(createRequest.Tenant.Name))
			})
		})

		When("user updates the existing tenant", func() {
			It("should update the tenant information", func() {
				updateRequest := repository.UpdateTenantRequest{
					TenantID: tenantID,
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}

				_, err := sut.UpdateTenant(ctx, &updateRequest)
				Ω(err).Should(BeNil())

				response, err := sut.ReadTenant(ctx, &repository.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				Ω(response.Tenant).ShouldNot(BeNil())
				Ω(response.Tenant.Name).Should(Equal(updateRequest.Tenant.Name))
			})
		})

		When("user deletes the tenant", func() {
			It("should delete the tenant", func() {
				_, err := sut.DeleteTenant(ctx, &repository.DeleteTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())

				response, err := sut.ReadTenant(ctx, &repository.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.TenantNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})
	})

	Context("tenant does not exist", func() {
		var (
			tenantID string
		)

		BeforeEach(func() {
			tenantID = cuid.New()
		})

		When("user reads the tenant", func() {
			It("should return NotFoundError", func() {
				response, err := sut.ReadTenant(ctx, &repository.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.TenantNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})

		When("user tries to update the tenant", func() {
			It("should return NotFoundError", func() {
				updateRequest := repository.UpdateTenantRequest{
					TenantID: tenantID,
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}

				response, err := sut.UpdateTenant(ctx, &updateRequest)
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.TenantNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})

		When("user tries to delete the tenant", func() {
			It("should return NotFoundError", func() {
				response, err := sut.DeleteTenant(ctx, &repository.DeleteTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(repository.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr repository.TenantNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})
	})

})
