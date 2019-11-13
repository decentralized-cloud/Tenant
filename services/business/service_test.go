package business_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business"
	repository "github.com/decentralized-cloud/tenant/services/repository"
	repsoitoryMock "github.com/decentralized-cloud/tenant/services/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBusinessService(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

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
						Return(&repository.CreateTenantResponse{
							TenantID: cuid.New(),
							Tenant: models.Tenant{
								Name: cuid.New(),
							}}, nil)

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
					It("should return expected details", func() {
						expectedResponse := repository.CreateTenantResponse{
							TenantID: cuid.New(),
							Tenant: models.Tenant{
								Name: cuid.New(),
							}}

						mockRepositoryService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(&expectedResponse, nil)

						response, err := sut.CreateTenant(ctx, &request)
						Ω(err).Should(BeNil())
						Ω(response.Err).Should(BeNil())
						Ω(response.TenantID).ShouldNot(BeNil())
						Ω(response.TenantID).Should(Equal(expectedResponse.TenantID))
						assertTenant(response.Tenant, expectedResponse.Tenant)
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
					assertTenant(response.Tenant, tenant)
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
				It("should return expected details", func() {
					expectedResponse := repository.UpdateTenantResponse{
						Tenant: models.Tenant{
							Name: cuid.New(),
						}}
					mockRepositoryService.
						EXPECT().
						UpdateTenant(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.UpdateTenant(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					assertTenant(response.Tenant, expectedResponse.Tenant)
				})
			})
		})
	})

	Describe("DeleteTenant is called", func() {
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

	Describe("Search is called", func() {
		var (
			request   business.SearchRequest
			tenantIDs []string
		)

		BeforeEach(func() {
			tenantIDs = []string{}
			for idx := 0; idx < rand.Intn(20)+1; idx++ {
				tenantIDs = append(tenantIDs, cuid.New())
			}

			request = business.SearchRequest{
				Pagination: common.Pagination{
					After:  convertStringToPointer(cuid.New()),
					First:  convertIntToPointer(rand.Intn(1000)),
					Before: convertStringToPointer(cuid.New()),
					Last:   convertIntToPointer(rand.Intn(1000)),
				},
				SortingOptions: []common.SortingOptionPair{
					common.SortingOptionPair{
						Name:      cuid.New(),
						Direction: common.Ascending,
					},
					common.SortingOptionPair{
						Name:      cuid.New(),
						Direction: common.Descending,
					},
				},
				TenantIDs: tenantIDs,
			}
		})

		Context("tenant service is instantiated", func() {
			When("Search is called", func() {
				It("should call tenant repository Search method", func() {
					mockRepositoryService.
						EXPECT().
						Search(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.SearchRequest) {
							Ω(mappedRequest.Pagination).Should(Equal(request.Pagination))
							Ω(mappedRequest.SortingOptions).Should(Equal(request.SortingOptions))
							Ω(mappedRequest.TenantIDs).Should(Equal(request.TenantIDs))
						}).
						Return(&repository.SearchResponse{}, nil)

					response, err := sut.Search(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("tenant repository Search is faced with any other error", func() {
				It("should return UnknownError", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						Search(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.Search(ctx, &request)
					Ω(err).Should(BeNil())
					assertUnknowError(expectedError.Error(), response.Err, expectedError)
				})
			})

			When("tenant repository Search completes successfully", func() {
				It("should return the list of matched tenantIDs", func() {
					tenants := []models.TenantWithCursor{}

					for idx := 0; idx < rand.Intn(20)+1; idx++ {
						tenants = append(tenants, models.TenantWithCursor{
							TenantID: cuid.New(),
							Tenant: models.Tenant{
								Name: cuid.New(),
							},
							Cursor: cuid.New(),
						})
					}

					expectedResponse := repository.SearchResponse{
						HasPreviousPage: (rand.Intn(10) % 2) == 0,
						HasNextPage:     (rand.Intn(10) % 2) == 0,
						TotalCount:      rand.Int63n(1000),
						Tenants:         tenants,
					}

					mockRepositoryService.
						EXPECT().
						Search(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.Search(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					Ω(response.HasPreviousPage).Should(Equal(expectedResponse.HasPreviousPage))
					Ω(response.HasNextPage).Should(Equal(expectedResponse.HasNextPage))
					Ω(response.TotalCount).Should(Equal(expectedResponse.TotalCount))
					Ω(response.Tenants).Should(Equal(expectedResponse.Tenants))
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

func assertTenant(tenant, expectedTenant models.Tenant) {
	Ω(tenant).ShouldNot(BeNil())
	Ω(tenant.Name).Should(Equal(expectedTenant.Name))
}

func convertStringToPointer(str string) *string {
	return &str
}

func convertIntToPointer(i int) *int {
	return &i
}
