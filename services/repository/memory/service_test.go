// Package repository implements different repository services required by the tenant service
package memory_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/repository"
	"github.com/decentralized-cloud/tenant/services/repository/memory"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
	"github.com/thoas/go-funk"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInMemoryRepositoryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In-Memory Repository Service Tests")
}

var _ = Describe("In-Memory Repository Service Tests", func() {
	var (
		sut           repository.RepositoryContract
		ctx           context.Context
		createRequest repository.CreateTenantRequest
	)

	BeforeEach(func() {
		sut, _ = memory.NewRepositoryService()
		ctx = context.Background()
		createRequest = repository.CreateTenantRequest{
			Tenant: models.Tenant{
				Name: cuid.New(),
			}}
	})

	Context("user tries to instantiate RepositoryService", func() {
		When("all dependecies are resolved and NewRepositoryService is called", func() {
			It("should instantiate the new RepositoryService", func() {
				service, err := memory.NewRepositoryService()
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

		When("user reads the tenant", func() {
			It("should return the tenant information", func() {
				response, err := sut.ReadTenant(ctx, &repository.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				Ω(response.Tenant).ShouldNot(BeNil())
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

	Context("tenants exist", func() {
		var (
			tenantIDs []string
			names     []string
		)

		BeforeEach(func() {
			rand.Seed(42)
			names = []string{}
			for idx := 0; idx < rand.Intn(20)+10; idx++ {
				names = append(names, cuid.New())
			}

			tenantIDs = funk.Map(names, func(name string) string {
				response, _ := sut.CreateTenant(ctx, &repository.CreateTenantRequest{
					Tenant: models.Tenant{
						Name: name,
					},
				})

				return response.TenantID
			}).([]string)
		})

		When("user search for tenants without any tenant ID provided", func() {
			It("should return all tenants", func() {
				response, err := sut.Search(ctx, &repository.SearchRequest{})
				Ω(err).Should(BeNil())
				Ω(response.Tenants).Should(HaveLen(len(tenantIDs)))

				filteredTenants := funk.Filter(response.Tenants, func(tenantWithCursor models.TenantWithCursor) bool {
					return !funk.Contains(tenantIDs, tenantWithCursor.TenantID)
				}).([]models.TenantWithCursor)

				Ω(filteredTenants).Should(HaveLen(0))
			})

			It("should sort the result ascending when no sorting direction is provided", func() {
				response, _ := sut.Search(ctx, &repository.SearchRequest{})
				names := funk.Map(response.Tenants, func(tenantWithCursor models.TenantWithCursor) string {
					return tenantWithCursor.Tenant.Name
				}).([]string)

				for idx := range names[:len(names)-1] {
					Ω(names[idx] < names[idx+1]).Should(BeTrue())
				}
			})

			It("should sort the result ascending when sorting direction is ascending", func() {
				response, _ := sut.Search(ctx, &repository.SearchRequest{
					SortingOptions: []common.SortingOptionPair{
						common.SortingOptionPair{
							Name:      "name",
							Direction: common.Ascending,
						}}})
				names := funk.Map(response.Tenants, func(tenantWithCursor models.TenantWithCursor) string {
					return tenantWithCursor.Tenant.Name
				}).([]string)

				for idx := range names[:len(names)-1] {
					Ω(names[idx] < names[idx+1]).Should(BeTrue())
				}
			})

			It("should sort the result descending when sorting direction is descending", func() {
				response, _ := sut.Search(ctx, &repository.SearchRequest{
					SortingOptions: []common.SortingOptionPair{
						common.SortingOptionPair{
							Name:      "name",
							Direction: common.Descending,
						}}})
				names := funk.Map(response.Tenants, func(tenantWithCursor models.TenantWithCursor) string {
					return tenantWithCursor.Tenant.Name
				}).([]string)

				for idx := range names[:len(names)-1] {
					Ω(names[idx] > names[idx+1]).Should(BeTrue())
				}
			})
		})

		When("user search for tenants with tenant IDs provided", func() {
			var (
				numberOfTenantIDs  int
				shuffeledTenantIDs []string
			)

			BeforeEach(func() {
				shuffeledTenantIDs = funk.ShuffleString(tenantIDs)
				numberOfTenantIDs = rand.Intn(10)
			})

			It("should return filtered tenant list", func() {
				response, err := sut.Search(ctx, &repository.SearchRequest{
					TenantIDs: shuffeledTenantIDs[:numberOfTenantIDs],
				})
				Ω(err).Should(BeNil())
				Ω(response.Tenants).Should(HaveLen(numberOfTenantIDs))

				filteredTenants := funk.Filter(response.Tenants, func(tenantWithCursor models.TenantWithCursor) bool {
					return !funk.Contains(tenantIDs, tenantWithCursor.TenantID)
				}).([]models.TenantWithCursor)

				Ω(filteredTenants).Should(HaveLen(0))
			})

			It("should sort the result ascending when no sorting direction is provided", func() {
				response, _ := sut.Search(ctx, &repository.SearchRequest{
					TenantIDs: shuffeledTenantIDs[:numberOfTenantIDs],
				})
				names := funk.Map(response.Tenants, func(tenantWithCursor models.TenantWithCursor) string {
					return tenantWithCursor.Tenant.Name
				}).([]string)

				for idx := range names[:len(names)-1] {
					Ω(names[idx] < names[idx+1]).Should(BeTrue())
				}
			})

			It("should sort the result ascending when sorting direction is ascending", func() {
				response, _ := sut.Search(ctx, &repository.SearchRequest{
					SortingOptions: []common.SortingOptionPair{
						common.SortingOptionPair{
							Name:      "name",
							Direction: common.Ascending,
						},
					},
					TenantIDs: shuffeledTenantIDs[:numberOfTenantIDs],
				})
				names := funk.Map(response.Tenants, func(tenantWithCursor models.TenantWithCursor) string {
					return tenantWithCursor.Tenant.Name
				}).([]string)

				for idx := range names[:len(names)-1] {
					Ω(names[idx] < names[idx+1]).Should(BeTrue())
				}
			})

			It("should sort the result descending when sorting direction is descending", func() {
				response, _ := sut.Search(ctx, &repository.SearchRequest{
					SortingOptions: []common.SortingOptionPair{
						common.SortingOptionPair{
							Name:      "name",
							Direction: common.Descending,
						},
					},
					TenantIDs: shuffeledTenantIDs[:numberOfTenantIDs],
				})
				names := funk.Map(response.Tenants, func(tenantWithCursor models.TenantWithCursor) string {
					return tenantWithCursor.Tenant.Name
				}).([]string)

				for idx := range names[:len(names)-1] {
					Ω(names[idx] > names[idx+1]).Should(BeTrue())
				}
			})
		})
	})
})
