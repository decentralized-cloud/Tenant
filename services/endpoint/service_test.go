package endpoint_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business"
	businessMock "github.com/decentralized-cloud/tenant/services/business/mock"
	"github.com/decentralized-cloud/tenant/services/endpoint"
	gokitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEndpointCreatorService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Endpoint Creator Service Tests")
}

var _ = Describe("Endpoint Creator Service Tests", func() {
	var (
		mockCtrl            *gomock.Controller
		sut                 endpoint.EndpointCreatorContract
		mockBusinessService *businessMock.MockBusinessContract
		ctx                 context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockBusinessService = businessMock.NewMockBusinessContract(mockCtrl)
		sut, _ = endpoint.NewEndpointCreatorService(mockBusinessService)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate EndpointCreatorService", func() {
		When("tenant business service is not provided and NewEndpointCreatorService is called", func() {
			It("should return ArgumentNilError", func() {
				service, err := endpoint.NewEndpointCreatorService(nil)
				Ω(service).Should(BeNil())
				assertArgumentNilError("businessService", "", err)
			})
		})

		When("all dependencies are resolved and NewEndpointCreatorService is called", func() {
			It("should instantiate the new EndpointCreatorService", func() {
				service, err := endpoint.NewEndpointCreatorService(mockBusinessService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("CreateTenantEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.CreateTenantEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.CreateTenantRequest
				response business.CreateTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.CreateTenantEndpoint()
				request = business.CreateTenantRequest{
					Tenant: models.Tenant{
						Name: cuid.New(),
					},
				}

				response = business.CreateTenantResponse{
					TenantID: cuid.New(),
				}
			})

			Context("CreateTenantEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.CreateTenantRequest{
							Tenant: models.Tenant{
								Name: "",
							}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service CreateTenant method", func() {
						mockBusinessService.
							EXPECT().
							CreateTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.CreateTenantRequest) {
								Ω(mappedRequest.Tenant).Should(Equal(request.Tenant))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service CreateTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service CreateTenant returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("ReadTenantEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.ReadTenantEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.ReadTenantRequest
				response business.ReadTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.ReadTenantEndpoint()
				request = business.ReadTenantRequest{
					TenantID: cuid.New(),
				}

				response = business.ReadTenantResponse{
					Tenant: models.Tenant{
						Name: cuid.New(),
					},
				}
			})

			Context("ReadTenantEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.ReadTenantRequest{
							TenantID: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service ReadTenant method", func() {
						mockBusinessService.
							EXPECT().
							ReadTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.ReadTenantRequest) {
								Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service ReadTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							ReadTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service ReadTenant returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							ReadTenant(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("UpdateTenantEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.UpdateTenantEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.UpdateTenantRequest
				response business.UpdateTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.UpdateTenantEndpoint()
				request = business.UpdateTenantRequest{
					TenantID: cuid.New(),
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}

				response = business.UpdateTenantResponse{}
			})

			Context("UpdateTenantEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.UpdateTenantRequest{
							TenantID: "",
							Tenant: models.Tenant{
								Name: "",
							}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service UpdateTenant method", func() {
						mockBusinessService.
							EXPECT().
							UpdateTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.UpdateTenantRequest) {
								Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service UpdateTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							UpdateTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service UpdateTenant returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							UpdateTenant(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("DeleteTenantEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.DeleteTenantEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.DeleteTenantRequest
				response business.DeleteTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.DeleteTenantEndpoint()
				request = business.DeleteTenantRequest{
					TenantID: cuid.New(),
				}

				response = business.DeleteTenantResponse{}
			})

			Context("DeleteTenantEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.DeleteTenantRequest{
							TenantID: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service DeleteTenant method", func() {
						mockBusinessService.
							EXPECT().
							DeleteTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.DeleteTenantRequest) {
								Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service DeleteTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							DeleteTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service DeleteTenant returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							DeleteTenant(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
				})
			})
		})
	})

	Context("EndpointCreatorService is instantiated", func() {
		When("SearchEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.SearchEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint  gokitendpoint.Endpoint
				tenantIDs []string
				request   business.SearchRequest
				response  business.SearchResponse
			)

			BeforeEach(func() {
				endpoint = sut.SearchEndpoint()
				rand.Seed(42)
				tenantIDs = []string{}
				for idx := 0; idx < rand.Intn(20)+1; idx++ {
					tenantIDs = append(tenantIDs, cuid.New())
				}

				request = business.SearchRequest{
					Pagination: common.Pagination{
						After:  cuid.New(),
						First:  rand.Intn(1000),
						Before: cuid.New(),
						Last:   rand.Intn(1000),
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

				response = business.SearchResponse{Tenants: tenants}
			})

			Context("SearchEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.SearchResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.SearchResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service Search method", func() {
						mockBusinessService.
							EXPECT().
							Search(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.SearchRequest) {
								Ω(mappedRequest.Pagination).Should(Equal(request.Pagination))
								Ω(mappedRequest.SortingOptions).Should(Equal(request.SortingOptions))
								Ω(mappedRequest.TenantIDs).Should(Equal(request.TenantIDs))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.SearchResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service Search returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							Search(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service Search returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							Search(gomock.Any(), gomock.Any()).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).Should(Equal(&response))
					})
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
