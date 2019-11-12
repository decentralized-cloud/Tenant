package mongodb_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	configurationMock "github.com/decentralized-cloud/tenant/services/configuration/mock"
	"github.com/decentralized-cloud/tenant/services/repository"
	"github.com/decentralized-cloud/tenant/services/repository/mongodb"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
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
			connectionString = "mongodb://localhost:27017"
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
				assertTenant(response.Tenant, createRequest.Tenant)
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
				assertTenant(response.Tenant, createRequest.Tenant)
			})
		})

		When("user updates the existing tenant", func() {
			It("should update the tenant information", func() {
				updateRequest := repository.UpdateTenantRequest{
					TenantID: tenantID,
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}

				updateResponse, err := sut.UpdateTenant(ctx, &updateRequest)
				Ω(err).Should(BeNil())
				assertTenant(updateResponse.Tenant, updateRequest.Tenant)

				readResponse, err := sut.ReadTenant(ctx, &repository.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				assertTenant(readResponse.Tenant, updateRequest.Tenant)
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

	Context("tenant already exists", func() {

		var (
			tenantIDs []string
		)

		BeforeEach(func() {

			tenantIDs = []string{}
			for i := 0; i < 10; i++ {

				tenantName := fmt.Sprintf("%s%d", "Name", i)
				createRequest.Tenant.Name = tenantName
				response, _ := sut.CreateTenant(ctx, &createRequest)
				tenantIDs = append(tenantIDs, response.TenantID)
			}
		})

		When("user searches for tenants with selected tenant Ids ", func() {
			It("should return all tenants which are matched the criteria ", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  "",
						First:  0,
						Before: "",
						Last:   0,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Tenants[i].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		When("user searches for tenants with selected tenant Ids and first 10 tenants ", func() {
			It("should return first 10 tenants", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  "",
						First:  10,
						Before: "",
						Last:   0,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Tenants[i].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		When("user searches for tenants with selected tenant Ids and first 5 tenants ", func() {
			It("should return first 5 tenants", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  "",
						First:  5,
						Before: "",
						Last:   0,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(5))
				for i := 0; i < 5; i++ {
					Ω(response.Tenants[i].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		When("user searches for tenants with selected tenant Ids with After parameter provided.", func() {
			It("should return first 9 tenants after provided tenant id", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  tenantIDs[0],
						First:  9,
						Before: "",
						Last:   0,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(9))
				for i := 1; i < 10; i++ {
					Ω(response.Tenants[i-1].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i-1].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		When("user searches for tenants with selected tenant Ids with After parameter provided.", func() {
			It("should return first 5 tenants after provided tenant id", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  tenantIDs[0],
						First:  5,
						Before: "",
						Last:   0,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(5))
				for i := 1; i < 5; i++ {
					Ω(response.Tenants[i-1].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i-1].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		//TODO : this test does not make sense
		When("user searches for tenants with selected tenant Ids and last 10 tenants ", func() {
			It("should return first 10 tenants", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  "",
						First:  0,
						Before: "",
						Last:   10,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Tenants[i].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		When("user searches for tenants with selected tenant Ids with Before parameter provided.", func() {
			It("should return first 9 tenants before provided tenant id", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  "",
						First:  0,
						Before: tenantIDs[9],
						Last:   9,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(9))
				for i := 0; i < 9; i++ {
					Ω(response.Tenants[i].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		When("user searches for tenants with selected tenant Ids and first 10 tenants with ascending order on name property ", func() {
			It("should return first 10 tenants in adcending order on name field", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  "",
						First:  10,
						Before: "",
						Last:   0,
					},
					SortingOptions: []common.SortingOptionPair{
						{Name: "name", Direction: common.Ascending},
					},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Tenants[i].TenantID).Should(Equal(tenantIDs[i]))
					tenantName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Tenants[i].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

		When("user searches for tenants with selected tenant Ids and first 10 tenants with descending order on name property ", func() {
			It("should return first 10 tenants in descending order on name field", func() {

				searchRequest := repository.SearchRequest{
					TenantIDs: tenantIDs,
					Pagination: common.Pagination{
						After:  "",
						First:  10,
						Before: "",
						Last:   0,
					},
					SortingOptions: []common.SortingOptionPair{
						{Name: "name", Direction: common.Descending},
					},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Tenants).ShouldNot(BeNil())
				Ω(len(response.Tenants)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Tenants[i].TenantID).Should(Equal(tenantIDs[9-i]))
					tenantName := fmt.Sprintf("%s%d", "Name", 9-i)
					Ω(response.Tenants[i].Tenant.Name).Should(Equal(tenantName))
				}
			})
		})

	})

})

func assertTenant(tenant, expectedTenant models.Tenant) {
	Ω(tenant).ShouldNot(BeNil())
	Ω(tenant.Name).Should(Equal(expectedTenant.Name))
}
