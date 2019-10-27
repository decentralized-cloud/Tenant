// Package repository implements different repository services required by the tenant service
package mongodb_test

import (
	"context"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/repository"
	"github.com/decentralized-cloud/tenant/services/repository/mongodb"
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
		sut           repository.RepositoryContract
		ctx           context.Context
		createRequest repository.CreateTenantRequest
	)

	BeforeEach(func() {
		sut, _ = mongodb.NewMongodbRepositoryService()
		ctx = context.TODO()
		createRequest = repository.CreateTenantRequest{
			Tenant: models.Tenant{
				Name: cuid.New(),
			}}
	})

	Context("user tries to instantiate RepositoryService", func() {
		When("all dependecies are resolved and NewRepositoryService is called", func() {
			It("should instantiate the new RepositoryService", func() {
				service, err := mongodb.NewMongodbRepositoryService()
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
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
				Ω(tenantID).ShouldNot(BeNil())
				response, err := sut.ReadTenant(ctx, &repository.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				Ω(response.Tenant.Name).Should(Equal(createRequest.Tenant.Name))
			})
		})
	})
})
