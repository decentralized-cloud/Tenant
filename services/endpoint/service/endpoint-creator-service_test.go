package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-kit/kit/endpoint"
	"github.com/lucsky/cuid"

	"github.com/decentralized-cloud/tenant/models"
	businessContract "github.com/decentralized-cloud/tenant/services/business/contract"
	businessMock "github.com/decentralized-cloud/tenant/services/business/mock"
	"github.com/decentralized-cloud/tenant/services/endpoint/contract"
	"github.com/decentralized-cloud/tenant/services/endpoint/service"
	"github.com/golang/mock/gomock"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEndpointCreatorService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EndpointCreatorService Tests")
}

var _ = Describe("EndpointCreatorService Tests", func() {
	var (
		mockCtrl                  *gomock.Controller
		sut                       contract.EndpointCreatorContract
		mockTenantBusinessService *businessMock.MockTenantServiceContract
		ctx                       context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())

		mockTenantBusinessService = businessMock.NewMockTenantServiceContract(mockCtrl)
		sut, _ = service.NewEndpointCreatorService(mockTenantBusinessService)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate EndpointCreatorService", func() {
		When("tenant business service is not provided and NewEndpointCreatorService is called", func() {
			It("should return ArgumentNilError", func() {
				service, err := service.NewEndpointCreatorService(nil)
				Ω(service).Should(BeNil())
				assertArgumentNilError("businessService", "", err)
			})
		})

		When("all dependencies are resolved and NewEndpointCreatorService is called", func() {
			It("should instantiate the new EndpointCreatorService", func() {
				service, err := service.NewEndpointCreatorService(mockTenantBusinessService)
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
				endpoint endpoint.Endpoint
				request  businessContract.CreateTenantRequest
				response businessContract.CreateTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.CreateTenantEndpoint()
				request = businessContract.CreateTenantRequest{
					Tenant: models.Tenant{
						Name: cuid.New(),
					},
				}

				response = businessContract.CreateTenantResponse{
					TenantID: cuid.New(),
				}
			})

			Context("CreateTenantEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.CreateTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.CreateTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := businessContract.CreateTenantRequest{
							Tenant: models.Tenant{
								Name: "",
							}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.CreateTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service CreateTenant method", func() {
						mockTenantBusinessService.
							EXPECT().
							CreateTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *businessContract.CreateTenantRequest) {
								Ω(mappedRequest.Tenant).Should(Equal(request.Tenant))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.CreateTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service CreateTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockTenantBusinessService.
							EXPECT().
							CreateTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service CreateTenant returns response", func() {
					It("should return the same response", func() {
						mockTenantBusinessService.
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
				endpoint endpoint.Endpoint
				request  businessContract.ReadTenantRequest
				response businessContract.ReadTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.ReadTenantEndpoint()
				request = businessContract.ReadTenantRequest{
					TenantID: cuid.New(),
				}

				response = businessContract.ReadTenantResponse{
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
						castedResponse := returnedResponse.(*businessContract.ReadTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.ReadTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := businessContract.ReadTenantRequest{
							TenantID: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.ReadTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service ReadTenant method", func() {
						mockTenantBusinessService.
							EXPECT().
							ReadTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *businessContract.ReadTenantRequest) {
								Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.ReadTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service ReadTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockTenantBusinessService.
							EXPECT().
							ReadTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service ReadTenant returns response", func() {
					It("should return the same response", func() {
						mockTenantBusinessService.
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
				endpoint endpoint.Endpoint
				request  businessContract.UpdateTenantRequest
				response businessContract.UpdateTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.UpdateTenantEndpoint()
				request = businessContract.UpdateTenantRequest{
					TenantID: cuid.New(),
					Tenant: models.Tenant{
						Name: cuid.New(),
					}}

				response = businessContract.UpdateTenantResponse{}
			})

			Context("UpdateTenantEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.UpdateTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.UpdateTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := businessContract.UpdateTenantRequest{
							TenantID: "",
							Tenant: models.Tenant{
								Name: "",
							}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.UpdateTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service UpdateTenant method", func() {
						mockTenantBusinessService.
							EXPECT().
							UpdateTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *businessContract.UpdateTenantRequest) {
								Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.UpdateTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service UpdateTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockTenantBusinessService.
							EXPECT().
							UpdateTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service UpdateTenant returns response", func() {
					It("should return the same response", func() {
						mockTenantBusinessService.
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
				endpoint endpoint.Endpoint
				request  businessContract.DeleteTenantRequest
				response businessContract.DeleteTenantResponse
			)

			BeforeEach(func() {
				endpoint = sut.DeleteTenantEndpoint()
				request = businessContract.DeleteTenantRequest{
					TenantID: cuid.New(),
				}

				response = businessContract.DeleteTenantResponse{}
			})

			Context("DeleteTenantEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.DeleteTenantResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.DeleteTenantResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := businessContract.DeleteTenantRequest{
							TenantID: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.DeleteTenantResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service DeleteTenant method", func() {
						mockTenantBusinessService.
							EXPECT().
							DeleteTenant(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *businessContract.DeleteTenantRequest) {
								Ω(mappedRequest.TenantID).Should(Equal(request.TenantID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*businessContract.DeleteTenantResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service DeleteTenant returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockTenantBusinessService.
							EXPECT().
							DeleteTenant(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service DeleteTenant returns response", func() {
					It("should return the same response", func() {
						mockTenantBusinessService.
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
