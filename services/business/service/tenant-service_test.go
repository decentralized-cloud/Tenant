package service_test

import (
	"context"
	"errors"
	"testing"

	commonErrors "github.com/decentralized-cloud/tenant/common/errors"
	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business/contract"
	"github.com/decentralized-cloud/tenant/services/business/service"
	repositoryContract "github.com/decentralized-cloud/tenant/services/repository/contract"
	repsoitoryMock "github.com/decentralized-cloud/tenant/services/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TenantTest Tests", func() {
	var (
		mockCtrl                    *gomock.Controller
		sut                         contract.TenantServiceContract
		mockTenantRepositoryService *repsoitoryMock.MockTenantRepositoryServiceContract
		ctx                         context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockTenantRepositoryService = repsoitoryMock.NewMockTenantRepositoryServiceContract(mockCtrl)
		sut, _ = service.NewTenantService(mockTenantRepositoryService)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate TenantService", func() {
		When("tenant repository service is not provided and NewTenantService is called", func() {
			It("should return ArgumentError", func() {
				service, err := service.NewTenantService(nil)
				Ω(err).Should(HaveOccurred())
				Ω(service).Should(BeNil())

				argumentErr, ok := err.(commonErrors.ArgumentError)
				Ω(ok).Should(BeTrue())
				Ω(argumentErr.ArgumentName).Should(Equal("repositoryService"))
			})
		})

		When("all dependecies are resolved and NewTenantService is called", func() {
			It("should instantiate the new TenantService", func() {
				service, err := service.NewTenantService(mockTenantRepositoryService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Describe("CreateTenant", func() {
		var (
			request contract.CreateTenantRequest
		)

		BeforeEach(func() {
			request = contract.CreateTenantRequest{
				Tenant: models.Tenant{
					Name: cuid.New(),
				}}
		})

		Context("tenant service is instantiated", func() {
			When("CreateTenant is called without context", func() {
				It("should return ArgumentError", func() {
					response, err := sut.CreateTenant(nil, &request)
					Ω(err).Should(HaveOccurred())
					Ω(response).Should(BeNil())

					argumentErr, ok := err.(commonErrors.ArgumentError)
					Ω(ok).Should(BeTrue())
					Ω(argumentErr.ArgumentName).Should(Equal("ctx"))
				})

				When("CreateTenant is called without request", func() {
					It("should return ArgumentError", func() {
						response, err := sut.CreateTenant(ctx, nil)
						Ω(err).Should(HaveOccurred())
						Ω(response).Should(BeNil())

						argumentErr, ok := err.(commonErrors.ArgumentError)
						Ω(ok).Should(BeTrue())
						Ω(argumentErr.ArgumentName).Should(Equal("request"))
					})
				})

				When("CreateTenant is called with invalid request", func() {
					It("should return ArgumentError", func() {
						invalidRequest := contract.CreateTenantRequest{
							Tenant: models.Tenant{
								Name: "",
							}}

						response, err := sut.CreateTenant(ctx, &invalidRequest)
						Ω(err).Should(HaveOccurred())
						Ω(response).Should(BeNil())

						argumentErr, ok := err.(commonErrors.ArgumentError)
						Ω(ok).Should(BeTrue())
						Ω(argumentErr.ArgumentName).Should(Equal("request"))
						Ω(argumentErr.ErrorMessage).Should(Equal(invalidRequest.Validate().Error()))
					})
				})

				When("CreateTenant is called with correct input paramters", func() {
					It("should call tenant respository CreateTenant method", func() {
						mockTenantRepositoryService.
							EXPECT().
							CreateTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *repositoryContract.CreateTenantRequest) {
								Ω(mappedRequest.Tenant).Should(Equal(request.Tenant))
							}).
							Return(&repositoryContract.CreateTenantResponse{TenantID: cuid.New()}, nil)

						_, err := sut.CreateTenant(ctx, &request)
						Ω(err).Should(BeNil())
					})

					When("And tenant repository CreateTenant return TenantAlreadyExistError", func() {
						It("should return TenantAlreadyExistsError", func() {
							mockTenantRepositoryService.
								EXPECT().
								CreateTenant(gomock.Any(), gomock.Any()).
								Return(nil, repositoryContract.NewTenantAlreadyExistsError())

							response, err := sut.CreateTenant(ctx, &request)
							Ω(err).Should(HaveOccurred())
							Ω(response).Should(BeNil())

							_, ok := err.(contract.TenantAlreadyExistsError)
							Ω(ok).Should(BeTrue())
						})
					})

					When("And tenant repository CreateTenant return any error rather than TenantAlreadyExistsError", func() {
						It("should return UnknownError", func() {
							expectedError := errors.New(cuid.New())
							mockTenantRepositoryService.
								EXPECT().
								CreateTenant(gomock.Any(), gomock.Any()).
								Return(nil, expectedError)

							response, err := sut.CreateTenant(ctx, &request)
							Ω(err).Should(HaveOccurred())
							Ω(response).Should(BeNil())

							unknownErr, ok := err.(contract.UnknownError)
							Ω(ok).Should(BeTrue())
							Ω(unknownErr.ErrorMessage).Should(Equal(expectedError.Error()))
						})
					})

					When("And tenant repository CreateTenant return no error", func() {
						It("should return the new tenantID", func() {
							tenantID := cuid.New()
							mockTenantRepositoryService.
								EXPECT().
								CreateTenant(gomock.Any(), gomock.Any()).
								Return(&repositoryContract.CreateTenantResponse{TenantID: tenantID}, nil)

							response, err := sut.CreateTenant(ctx, &request)
							Ω(err).Should(BeNil())
							Ω(response).ShouldNot(BeNil())
							Ω(response.TenantID).Should(Equal(tenantID))
						})
					})

				})
			})
		})
	})
})

func TestCreateTenant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TenantService Tests")
}
