package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business/contract"
	"github.com/decentralized-cloud/tenant/services/business/service"
	repositoryContract "github.com/decentralized-cloud/tenant/services/repository/contract"
	repsoitoryMock "github.com/decentralized-cloud/tenant/services/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TenantService Tests", func() {
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
			It("should return ArgumentNilError", func() {
				service, err := service.NewTenantService(nil)
				Ω(service).Should(BeNil())
				assertArgumentNilError("repositoryService", "", err)
			})
		})

		When("all dependencies are resolved and NewTenantService is called", func() {
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
				It("should return ArgumentNilError", func() {
					response, err := sut.CreateTenant(nil, &request)
					Ω(err).Should(BeNil())
					assertArgumentNilError("ctx", "", response.Err)
				})
			})

			When("CreateTenant is called without request", func() {
				It("should return ArgumentNilError", func() {
					response, err := sut.CreateTenant(ctx, nil)
					Ω(err).Should(BeNil())
					assertArgumentNilError("request", "", response.Err)
				})
			})

			When("CreateTenant is called with invalid request", func() {
				It("should return ArgumentNilError", func() {
					invalidRequest := contract.CreateTenantRequest{
						Tenant: models.Tenant{
							Name: "",
						}}

					response, err := sut.CreateTenant(ctx, &invalidRequest)
					Ω(err).Should(BeNil())
					validationErr := invalidRequest.Validate()
					assertArgumentError("request", validationErr.Error(), response.Err, validationErr)
				})
			})

			When("CreateTenant is called with correct input paramters", func() {
				It("should call tenant repository CreateTenant method", func() {
					mockTenantRepositoryService.
						EXPECT().
						CreateTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repositoryContract.CreateTenantRequest) {
							Ω(mappedRequest.Tenant).Should(Equal(request.Tenant))
						}).
						Return(&repositoryContract.CreateTenantResponse{TenantID: cuid.New()}, nil)

					response, err := sut.CreateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})

				When("And tenant repository CreateTenant return TenantAlreadyExistError", func() {
					It("should return TenantAlreadyExistsError", func() {
						expectedError := repositoryContract.NewTenantAlreadyExistsError()
						mockTenantRepositoryService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedError)

						response, err := sut.CreateTenant(ctx, &request)
						Ω(err).Should(BeNil())
						assertTenantAlreadyExistsError(response.Err, expectedError)
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
						Ω(err).Should(BeNil())
						assertUnknowError(expectedError.Error(), response.Err, expectedError)
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
			request contract.ReadTenantRequest
		)

		BeforeEach(func() {
			request = contract.ReadTenantRequest{
				TenantID: cuid.New(),
			}
		})

		Context("tenant service is instantiated", func() {
			When("ReadTenant is called without context", func() {
				It("should return ArgumentNilError", func() {
					response, err := sut.ReadTenant(nil, &request)
					Ω(err).Should(BeNil())
					assertArgumentNilError("ctx", "", response.Err)
				})
			})

			When("ReadTenant is called without request", func() {
				It("should return ArgumentNilError", func() {
					response, err := sut.ReadTenant(ctx, nil)
					Ω(err).Should(BeNil())
					assertArgumentNilError("request", "", response.Err)
				})
			})

			When("ReadTenant is called with invalid request", func() {
				It("should return ArgumentNilError", func() {
					invalidRequest := contract.ReadTenantRequest{
						TenantID: "",
					}

					response, err := sut.ReadTenant(ctx, &invalidRequest)
					Ω(err).Should(BeNil())
					validationErr := invalidRequest.Validate()
					assertArgumentError("request", validationErr.Error(), response.Err, validationErr)
				})
			})

			When("ReadTenant is called with correct input paramters", func() {
				It("should call tenant repository ReadTenant method", func() {
					mockTenantRepositoryService.
						EXPECT().
						ReadTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repositoryContract.ReadTenantRequest) {
							Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
						}).
						Return(&repositoryContract.ReadTenantResponse{Tenant: models.Tenant{Name: cuid.New()}}, nil)

					response, err := sut.ReadTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And tenant repository ReadTenant return TenantNotFoundError", func() {
				It("should return TenantNotFoundError", func() {
					expectedError := repositoryContract.NewTenantNotFoundError(request.TenantID)
					mockTenantRepositoryService.
						EXPECT().
						ReadTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ReadTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertTenantNotFoundError(request.TenantID, response.Err, expectedError)
				})
			})

			When("And tenant repository ReadTenant return any error rather than TenantNotFoundError", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockTenantRepositoryService.
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
					mockTenantRepositoryService.
						EXPECT().
						ReadTenant(gomock.Any(), gomock.Any()).
						Return(&repositoryContract.ReadTenantResponse{Tenant: tenant}, nil)

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
			request contract.UpdateTenantRequest
		)

		BeforeEach(func() {
			request = contract.UpdateTenantRequest{
				TenantID: cuid.New(),
				Tenant:   models.Tenant{Name: cuid.New()},
			}
		})

		Context("tenant service is instantiated", func() {
			When("UpdateTenant is called without context", func() {
				It("should return ArgumentNilError", func() {
					response, err := sut.UpdateTenant(nil, &request)
					Ω(err).Should(BeNil())
					assertArgumentNilError("ctx", "", response.Err)
				})
			})

			When("UpdateTenant is called without request", func() {
				It("should return ArgumentNilError", func() {
					response, err := sut.UpdateTenant(ctx, nil)
					Ω(err).Should(BeNil())
					assertArgumentNilError("request", "", response.Err)
				})
			})

			When("UpdateTenant is called with invalid request", func() {
				It("should return ArgumentNilError", func() {
					invalidRequest := contract.UpdateTenantRequest{
						TenantID: "",
						Tenant:   models.Tenant{Name: ""},
					}

					response, err := sut.UpdateTenant(ctx, &invalidRequest)
					Ω(err).Should(BeNil())
					validationErr := invalidRequest.Validate()
					assertArgumentError("request", validationErr.Error(), response.Err, validationErr)
				})
			})

			When("UpdateTenant is called with correct input paramters", func() {
				It("should call tenant repository UpdateTenant method", func() {
					mockTenantRepositoryService.
						EXPECT().
						UpdateTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repositoryContract.UpdateTenantRequest) {
							Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							Ω(mappedRequest.Tenant.Name).Should(Equal(request.Tenant.Name))
						}).
						Return(&repositoryContract.UpdateTenantResponse{}, nil)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And tenant repository UpdateTenant return TenantNotFoundError", func() {
				It("should return TenantNotFoundError", func() {
					expectedError := repositoryContract.NewTenantNotFoundError(request.TenantID)
					mockTenantRepositoryService.
						EXPECT().
						UpdateTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertTenantNotFoundError(request.TenantID, response.Err, expectedError)
				})
			})

			When("And tenant repository UpdateTenant return any error rather than TenantNotFoundError", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockTenantRepositoryService.
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
					mockTenantRepositoryService.
						EXPECT().
						UpdateTenant(gomock.Any(), gomock.Any()).
						Return(&repositoryContract.UpdateTenantResponse{}, nil)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})
		})
	})

	Describe("DeleteTenant is called.", func() {
		var (
			request contract.DeleteTenantRequest
		)

		BeforeEach(func() {
			request = contract.DeleteTenantRequest{
				TenantID: cuid.New(),
			}
		})

		Context("tenant service is instantiated", func() {
			When("context is null", func() {
				It("should return ArgumentNilError and ArgumentName matches the context argument name", func() {
					response, err := sut.DeleteTenant(nil, &request)
					Ω(err).Should(BeNil())
					assertArgumentNilError("ctx", "", response.Err)
				})
			})

			When("request is null", func() {
				It("should return ArgumentNilError and ArgumentName matches the request argument name", func() {
					response, err := sut.DeleteTenant(ctx, nil)
					Ω(err).Should(BeNil())
					assertArgumentNilError("request", "", response.Err)
				})
			})

			When("request is invalid", func() {
				It("should return ArgumentNilError and both ArgumentName and ErrorMessage are matched", func() {
					invalidRequest := contract.DeleteTenantRequest{
						TenantID: "",
					}

					response, err := sut.DeleteTenant(ctx, &invalidRequest)
					Ω(err).Should(BeNil())
					validationErr := invalidRequest.Validate()
					assertArgumentError("request", validationErr.Error(), response.Err, validationErr)
				})
			})

			When("input paramters are valid", func() {
				It("should call tenant repository DeleteTenant method", func() {
					mockTenantRepositoryService.
						EXPECT().
						DeleteTenant(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repositoryContract.DeleteTenantRequest) {
							Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
						}).
						Return(&repositoryContract.DeleteTenantResponse{}, nil)

					response, err := sut.DeleteTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("tenant repository DeleteTenant can not find matched tenant", func() {
				It("should return TenantNotFoundError", func() {
					expectedError := repositoryContract.NewTenantNotFoundError(request.TenantID)
					mockTenantRepositoryService.
						EXPECT().
						DeleteTenant(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.DeleteTenant(ctx, &request)
					Ω(err).Should(BeNil())
					assertTenantNotFoundError(request.TenantID, response.Err, expectedError)
				})
			})
			When("tenant repository DeleteTenant is faced with any error rather than TenantNotFoundError", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockTenantRepositoryService.
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
					mockTenantRepositoryService.
						EXPECT().
						DeleteTenant(gomock.Any(), gomock.Any()).
						Return(&repositoryContract.DeleteTenantResponse{}, nil)

					response, err := sut.DeleteTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})
		})
	})
})

func TestTenantService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TenantService Tests")
}

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
	Ω(contract.IsUnknownError(err)).Should(BeTrue())

	var unknownErr contract.UnknownError
	_ = errors.As(err, &unknownErr)

	Ω(strings.Contains(unknownErr.Error(), expectedMessage)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertTenantAlreadyExistsError(err error, nestedErr error) {
	Ω(contract.IsTenantAlreadyExistsError(err)).Should(BeTrue())
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}

func assertTenantNotFoundError(expectedTenantID string, err error, nestedErr error) {
	Ω(contract.IsTenantNotFoundError(err)).Should(BeTrue())

	var tenantNotFoundErr contract.TenantNotFoundError
	_ = errors.As(err, &tenantNotFoundErr)

	Ω(tenantNotFoundErr.TenantID).Should(Equal(expectedTenantID))
	Ω(errors.Unwrap(err)).Should(Equal(nestedErr))
}
