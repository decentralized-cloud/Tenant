package service_test

import (
	"context"
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

	Describe("Given user going to create a new tenant", func() {
		Context("When create tenant is called", func() {
			It("Then new tenant is created", func() {
				response, err := sut.CreateTenant(ctx, &createRequest)
				Ω(err).Should(BeNil())
				Ω(response.TenantID).ShouldNot(BeNil())
			})
		})
	})

	Describe("Given tenant already exists", func() {
		var (
			tenantID string
		)

		BeforeEach(func() {
			response, _ := sut.CreateTenant(ctx, &createRequest)
			tenantID = response.TenantID
		})

		Context("When user reads the tenant", func() {
			It("Then tenant information is returned", func() {
				response, err := sut.ReadTenant(ctx, &contract.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())
				Ω(response.Tenant).ShouldNot(BeNil())
				Ω(response.Tenant.Name).Should(Equal(createRequest.Tenant.Name))
			})
		})

		Context("When user updates the existing tenant", func() {
			It("Then tenant information is updated", func() {
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

		Context("When user deletes the tenant", func() {
			It("Then tenant is deleted", func() {
				_, err := sut.DeleteTenant(ctx, &contract.DeleteTenantRequest{TenantID: tenantID})
				Ω(err).Should(BeNil())

				response, err := sut.ReadTenant(ctx, &contract.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				notFoundErr, ok := err.(contract.TenantNotFoundError)
				Ω(ok).Should(BeTrue())
				Ω(notFoundErr.TenantID).Should(Equal(tenantID))

			})
		})
	})

	Describe("Given tenant does not exist", func() {

		var (
			tenantID string
		)

		BeforeEach(func() {
			tenantID = cuid.New()
		})

		Context("When user reads the tenant", func() {
			It("Then not found error is returned", func() {
				response, err := sut.ReadTenant(ctx, &contract.ReadTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				notFoundErr, ok := err.(contract.TenantNotFoundError)
				Ω(ok).Should(BeTrue())
				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})

		Context("When user tries to update the tenant", func() {
			It("Then not found error is returned", func() {

				updateRequest := contract.UpdateTenantRequest{
					TenantID: tenantID,
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}
				response, err := sut.UpdateTenant(ctx, &updateRequest)
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				notFoundErr, ok := err.(contract.TenantNotFoundError)
				Ω(ok).Should(BeTrue())
				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})

		Context("When user tries to delete the tenant", func() {
			It("Then not found error is returned", func() {
				response, err := sut.DeleteTenant(ctx, &contract.DeleteTenantRequest{TenantID: tenantID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				notFoundErr, ok := err.(contract.TenantNotFoundError)
				Ω(ok).Should(BeTrue())
				Ω(notFoundErr.TenantID).Should(Equal(tenantID))
			})
		})
	})
})

func TestCreateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TenantRepositoryService Tests")
}
