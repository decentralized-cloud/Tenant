package service_test

import (
	"context"
	"testing"

	commonErrors "github.com/decentralized-cloud/tenant/common/errors"
	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business/contract"
	"github.com/decentralized-cloud/tenant/services/business/service"
	repsoitoryMocks "github.com/decentralized-cloud/tenant/services/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateTenant Tests", func() {
	var (
		mockCtrl                    *gomock.Controller
		tenantService               contract.TenantServiceContract
		mockTenantRepositoryService *repsoitoryMocks.MockTenantRepositoryServiceContract
		ctx                         context.Context
		request                     contract.CreateTenantRequest
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockTenantRepositoryService = repsoitoryMocks.NewMockTenantRepositoryServiceContract(mockCtrl)
		tenantService, _ = service.NewTenantService(mockTenantRepositoryService)

		ctx = context.Background()
		request = contract.CreateTenantRequest{
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
			request = contract.CreateTenantRequest{
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
