package services_test

import (
	"context"
	"testing"

	"github.com/decentralized-cloud/tenant/business/contracts"
	"github.com/decentralized-cloud/tenant/business/services"
	commonErrors "github.com/decentralized-cloud/tenant/common/errors"
	"github.com/decentralized-cloud/tenant/models"
	repsoitoryMocks "github.com/decentralized-cloud/tenant/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateTenant Tests", func() {
	var (
		mockCtrl                    *gomock.Controller
		tenantService               contracts.TenantServiceContract
		mockTenantRepositoryService *repsoitoryMocks.MockTenantRepositoryServiceContract
		ctx                         context.Context
		request                     contracts.CreateTenantRequest
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockTenantRepositoryService = repsoitoryMocks.NewMockTenantRepositoryServiceContract(mockCtrl)
		tenantService, _ = services.NewTenantService(mockTenantRepositoryService)

		ctx = context.Background()
		request = contracts.CreateTenantRequest{
			Tenant: models.Tenant{
				Name: cuid.New(),
			}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Input Parameters", func() {
		It("should return error when ctx is nil", func() {
			_, err := tenantService.CreateTenant(nil, &request)

			Ω(err).Should(HaveOccurred())

			_, ok := err.(commonErrors.ArgumentError)

			Ω(ok).Should(BeTrue())
		})

		It("should return error when request is nil", func() {
			_, err := tenantService.CreateTenant(ctx, nil)

			Ω(err).Should(HaveOccurred())

			_, ok := err.(commonErrors.ArgumentError)

			Ω(ok).Should(BeTrue())
		})

		It("should return error when request is invalid", func() {
			request = contracts.CreateTenantRequest{
				Tenant: models.Tenant{
					Name: "",
				}}
			_, err := tenantService.CreateTenant(ctx, &request)

			Ω(err).Should(HaveOccurred())

			_, ok := err.(commonErrors.ArgumentError)

			Ω(ok).Should(BeTrue())
		})
	})
})

func TestCreateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateTenant Tests")
}
