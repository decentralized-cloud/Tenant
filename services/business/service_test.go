package business_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business"
	repository "github.com/decentralized-cloud/tenant/services/repository"
	repsoitoryMock "github.com/decentralized-cloud/tenant/services/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBusinessService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Business Service Tests")
}

var _ = Describe("Business Service Tests", func() {
	var (
		mockCtrl              *gomock.Controller
		sut                   business.BusinessContract
		mockRepositoryService *repsoitoryMock.MockRepositoryContract
		ctx                   context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockRepositoryService = repsoitoryMock.NewMockRepositoryContract(mockCtrl)
		sut, _ = business.NewBusinessService(mockRepositoryService)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate BusinessService", func() {
		When("tenant repository service is not provided and NewBusinessService is called", func() {
			It("should return ArgumentNilError", func() {
				service, err := business.NewBusinessService(nil)
				Ω(service).Should(BeNil())
				assertArgumentNilError("repositoryService", "", err)
			})
		})

		When("all dependencies are resolved and NewBusinessService is called", func() {
			It("should instantiate the new BusinessService", func() {
				service, err := business.NewBusinessService(mockRepositoryService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Describe("CreateTenant", func() {
		var (
			request business.CreateTenantRequest
		)

		BeforeEach(func() {
			request = business.CreateTenantRequest{
				Tenant: models.Tenant{
					Name: cuid.New(),
				}}
		})

		Context("tenant service is instantiated", func() {
			When("CreateTenant is called", func() {
				It("should call tenant repository CreateTenant method", func() {
					mockRepositoryService.
						EXPECT().
						CreateTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.CreateTenantRequest) {
							Ω(mappedRequest.Tenant).Should(Equal(request.Tenant))
						}).
						Return(&repository.CreateTenantResponse{TenantID: cuid.New()}, nil)

					response, err := sut.CreateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})

				When("And tenant repository CreateTenant return TenantAlreadyExistError", func() {
					It("should return TenantAlreadyExistsError", func() {
						expectedError := repository.NewTenantAlreadyExistsError()
						mockRepositoryService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedError)

						response, err := sut.CreateTenant(ctx, &request)
						Ω(err).Should(BeNil())
						assertTenantAlreadyExistsError(response.Err, expectedError)
					})
				})

				When("And tenant repository CreateTenant return any other error", func() {
					It("should return UnknownError", func() {
						expectedError := errors.New(cuid.New())
						mockRepositoryService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedError)

						response, err := sut.CreateTenant(ctx, &request)
						Ω(err).Should(BeNil())
						assertUnknowError(expectedError.Error(), response.Err, expectedError)
					})
				})

				When("And tenant repository CreateTenant return no error", func() {
					It("should return the new tenantID", func() {
						tenantID := cuid.New()
						mockRepositoryService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(&repository.CreateTenantResponse{TenantID: tenantID}, nil)

						response, err := sut.CreateTenant(ctx, &request)
						Ω(err).Should(BeNil())
						Ω(response.Err).Should(BeNil())
						Ω(response.TenantID).ShouldNot(BeNil())
						Ω(response.TenantID).Should(Equal(tenantID))
					})
				})
			})
		})
	})

	Describe("ReadTenant", func() {
		var (
			request business.ReadTenantRequest
		)

		BeforeEach(func() {
			request = business.ReadTenantRequest{
				TenantID: cuid.New(),
			}
		})

		Context("tenant service is instantiated", func() {
			When("ReadTenant is called", func() {
				It("should call tenant repository ReadTenant method", func() {
					mockRepositoryService.
						EXPECT().
						ReadTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.ReadTenantRequest) {
							Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
						}).
						Return(&repository.ReadTenantResponse{Tenant: models.Tenant{Name: cuid.New()}}, nil)

					response, err := sut.ReadTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And tenant repository ReadTenant cannot find provided tenant", func() {
				It("should return TenantNotFoundError", func() {
					expectedError := repository.NewTenantNotFoundError(request.TenantID)
					mockRepositoryService.
						EXPECT().
						ReadTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ReadTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertTenantNotFoundError(request.TenantID, response.Err, expectedError)
				})
			})

			When("And tenant repository ReadTenant return any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						ReadTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ReadTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("And tenant repository ReadTenant return no error", func() {
				It("should return the tenantID", func() {
					tenant := models.Tenant{Name: cuid.New()}
					mockRepositoryService.
						EXPECT().
						ReadTenant(gomock.Any(), gomock.Any()).
						Return(&repository.ReadTenantResponse{Tenant: tenant}, nil)

					response, err := sut.ReadTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					Ω(response.Tenant).ShouldNot(BeNil())
					Ω(response.Tenant.Name).Should(Equal(tenant.Name))
				})
			})
		})
	})

	Describe("UpdateTenant", func() {
		var (
			request business.UpdateTenantRequest
		)

		BeforeEach(func() {
			request = business.UpdateTenantRequest{
				TenantID: cuid.New(),
				Tenant:   models.Tenant{Name: cuid.New()},
			}
		})

		Context("tenant service is instantiated", func() {
			When("UpdateTenant is called", func() {
				It("should call tenant repository UpdateTenant method", func() {
					mockRepositoryService.
						EXPECT().
						UpdateTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.UpdateTenantRequest) {
							Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							Ω(mappedRequest.Tenant.Name).Should(Equal(request.Tenant.Name))
						}).
						Return(&repository.UpdateTenantResponse{}, nil)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And tenant repository UpdateTenant cannot find provided tenant", func() {
				It("should return TenantNotFoundError", func() {
					expectedError := repository.NewTenantNotFoundError(request.TenantID)
					mockRepositoryService.
						EXPECT().
						UpdateTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertTenantNotFoundError(request.TenantID, response.Err, expectedError)
				})
			})

			When("And tenant repository UpdateTenant return any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						UpdateTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("And tenant repository UpdateTenant return no error", func() {
				It("should return no error", func() {
					mockRepositoryService.
						EXPECT().
						UpdateTenant(gomock.Any(), gomock.Any()).
						Return(&repository.UpdateTenantResponse{}, nil)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})
		})
	})

	Describe("DeleteTenant is called.", func() {
		var (
			request business.DeleteTenantRequest
		)

		BeforeEach(func() {
			request = business.DeleteTenantRequest{
				TenantID: cuid.New(),
			}
		})

		Context("tenant service is instantiated", func() {
			When("DeleteTenant is called", func() {
				It("should call tenant repository DeleteTenant method", func() {
					mockRepositoryService.
						EXPECT().
						DeleteTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.DeleteTenantRequest) {
							Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
						}).
						Return(&repository.DeleteTenantResponse{}, nil)

					response, err := sut.DeleteTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("tenant repository DeleteTenant cannot find provided tenant", func() {
				It("should return TenantNotFoundError", func() {
					expectedError := repository.NewTenantNotFoundError(request.TenantID)
					mockRepositoryService.
						EXPECT().
						DeleteTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.DeleteTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertTenantNotFoundError(request.TenantID, response.Err, expectedError)
				})
			})
			When("tenant repository DeleteTenant is faced with any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						DeleteTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.DeleteTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("tenant repository DeleteTenant completes successfully", func() {
				It("should return no error", func() {
					mockRepositoryService.
						EXPECT().
						DeleteTenant(gomock.Any(), gomock.Any()).
						Return(&repository.DeleteTenantResponse{}, nil)

					response, err := sut.DeleteTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})
		})
	})
})

func assertArgumentNilError(expectedArgumentName, expectedMessage string, err error) {
	Ω(commonErrors.IsArgumentNilError(err)).Should(BeTrue())

	var argumentNilErr commonErrors.ArgumentNilError
	_ = errors.As(err, &argumentNilErr)

	if expectedArgumentName != "" {
		Ω(argumentNilErr.ArgumentName).Should(Equal(expectedArgumentName))
	}

	if expectedMessage != "" {
		Ω(strings.Contains(argumentNilErr.Error(), expectedMessage)).Should(BeTrue())
	}
}

func assertArgumentError(expectedArgumentName, expectedMessage string, err error, nestedErr error) {
	Ω(commonErrors.IsArgumentError(err)).Should(BeTrue())

	var argumentErr commonErrors.ArgumentError
	_ = errors.As(err, &argumentErr)

	Ω(argumentErr.ArgumentName).Should(Equal(expectedArgumentName))
	Ω(strings.Contains(argumentErr.Error(), expectedMessage)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertUnknowError(expectedMessage string, err error, nestedErr error) {
	Ω(business.IsUnknownError(err)).Should(BeTrue())

	var unknownErr business.UnknownError
	_ = errors.As(err, &unknownErr)

	Ω(strings.Contains(unknownErr.Error(), expectedMessage)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertTenantAlreadyExistsError(err error, nestedErr error) {
	Ω(business.IsTenantAlreadyExistsError(err)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertTenantNotFoundError(expectedTenantID string, err error, nestedErr error) {
	Ω(business.IsTenantNotFoundError(err)).Should(BeTrue())

	var tenantNotFoundErr business.TenantNotFoundError
	_ = errors.As(err, &tenantNotFoundErr)

	Ω(tenantNotFoundErr.TenantID).Should(Equal(expectedTenantID))
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}
