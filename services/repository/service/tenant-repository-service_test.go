package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/repository/contract"
	"github.com/decentralized-cloud/tenant/services/repository/service"
	"github.com/lucsky/cuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TenantRepositoryService Tests", func() {

	var (
		sut           contract.TenantRepositoryServiceContract
		ctx           context.Context
		createRequest contract.CreateTenantRequest
	)

	BeforeEach(func() {
		sut, _ = service.NewTenantRepositoryService()
		ctx = context.Background()
		createRequest = contract.CreateTenantRequest{
			Tenant: models.Tenant{
				Name: cuid.New(),
			}}
	})

	Context("user tries to instantiate TenantRepositoryService", func() {
		When("all dependecies are resolved and NewTenantRepositoryService is called", func() {
			It("should instantiate the new TenantRepositoryService", func() {
				service, err := service.NewTenantRepositoryService()
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
				response, err := sut.ReadTenant(ctx, &contract.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				Ω(response.Tenant).ShouldNot(BeNil())
				Ω(response.Tenant.Name).Should(Equal(createRequest.Tenant.Name))
			})
		})

		When("user updates the existing tenant", func() {
			It("should update the tenant information", func() {
				updateRequest := contract.UpdateTenantRequest{
					TenantID: tenantID,
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}

				_, err := sut.UpdateTenant(ctx, &updateRequest)
				Ω(err).Should(BeNil())

				response, err := sut.ReadTenant(ctx, &contract.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				Ω(response.Tenant).ShouldNot(BeNil())
				Ω(response.Tenant.Name).Should(Equal(updateRequest.Tenant.Name))
			})
		})

		When("user deletes the tenant", func() {
			It("should delete the tenant", func() {
				_, err := sut.DeleteTenant(ctx, &contract.DeleteTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())

				response, err := sut.ReadTenant(ctx, &contract.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(contract.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr contract.TenantNotFoundError
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
			It("should return NotgFoundError", func() {
				response, err := sut.ReadTenant(ctx, &contract.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(contract.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr contract.TenantNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})

		When("user tries to update the tenant", func() {
			It("should return NotgFoundError", func() {
				updateRequest := contract.UpdateTenantRequest{
					TenantID: tenantID,
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}
				response, err := sut.UpdateTenant(ctx, &updateRequest)
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(contract.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr contract.TenantNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})

		When("user tries to delete the tenant", func() {
			It("should return NotgFoundError", func() {
				response, err := sut.DeleteTenant(ctx, &contract.DeleteTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(contract.IsTenantNotFoundError(err)).Should(BeTrue())

				var notFoundErr contract.TenantNotFoundError
				_ = errors.As(err, &notFoundErr)

				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})
	})
})

func TestTenantRepositoryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TenantRepositoryService Tests")
}
