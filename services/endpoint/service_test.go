package endpoint_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/decentralized-cloud/project/models"
	"github.com/decentralized-cloud/project/services/business"
	businessMock "github.com/decentralized-cloud/project/services/business/mock"
	"github.com/decentralized-cloud/project/services/endpoint"
	gokitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEndpointCreatorService(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

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
		When("project business service is not provided and NewEndpointCreatorService is called", func() {
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
		When("CreateProjectEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.CreateProjectEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.CreateProjectRequest
				response business.CreateProjectResponse
			)

			BeforeEach(func() {
				endpoint = sut.CreateProjectEndpoint()
				request = business.CreateProjectRequest{
					UserEmail: cuid.New() + "@test.com",
					Project: models.Project{
						Name: cuid.New(),
					},
				}

				response = business.CreateProjectResponse{
					ProjectID: cuid.New(),
					Project: models.Project{
						Name: cuid.New(),
					},
					Cursor: cuid.New(),
				}
			})

			Context("CreateProjectEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateProjectResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateProjectResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.CreateProjectRequest{
							Project: models.Project{
								Name: "",
							}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateProjectResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service CreateProject method", func() {
						mockBusinessService.
							EXPECT().
							CreateProject(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.CreateProjectRequest) {
								Ω(mappedRequest.Project).Should(Equal(request.Project))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.CreateProjectResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service CreateProject returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							CreateProject(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service CreateProject returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							CreateProject(gomock.Any(), gomock.Any()).
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
		When("ReadProjectEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.ReadProjectEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.ReadProjectRequest
				response business.ReadProjectResponse
			)

			BeforeEach(func() {
				endpoint = sut.ReadProjectEndpoint()
				request = business.ReadProjectRequest{
					UserEmail: cuid.New() + "@test.com",
					ProjectID: cuid.New(),
				}

				response = business.ReadProjectResponse{
					Project: models.Project{
						Name: cuid.New(),
					},
				}
			})

			Context("ReadProjectEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadProjectResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadProjectResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.ReadProjectRequest{
							ProjectID: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadProjectResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service ReadProject method", func() {
						mockBusinessService.
							EXPECT().
							ReadProject(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.ReadProjectRequest) {
								Ω(mappedRequest.ProjectID).Should(Equal(request.ProjectID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.ReadProjectResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service ReadProject returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							ReadProject(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service ReadProject returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							ReadProject(gomock.Any(), gomock.Any()).
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
		When("UpdateProjectEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.UpdateProjectEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.UpdateProjectRequest
				response business.UpdateProjectResponse
			)

			BeforeEach(func() {
				endpoint = sut.UpdateProjectEndpoint()
				request = business.UpdateProjectRequest{
					UserEmail: cuid.New() + "@test.com",
					ProjectID: cuid.New(),
					Project: models.Project{
						Name: cuid.New(),
					}}

				response = business.UpdateProjectResponse{
					Project: models.Project{
						Name: cuid.New(),
					},
					Cursor: cuid.New(),
				}
			})

			Context("UpdateProjectEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateProjectResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateProjectResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.UpdateProjectRequest{
							ProjectID: "",
							Project: models.Project{
								Name: "",
							}}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateProjectResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service UpdateProject method", func() {
						mockBusinessService.
							EXPECT().
							UpdateProject(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.UpdateProjectRequest) {
								Ω(mappedRequest.ProjectID).Should(Equal(request.ProjectID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.UpdateProjectResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service UpdateProject returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							UpdateProject(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service UpdateProject returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							UpdateProject(gomock.Any(), gomock.Any()).
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
		When("DeleteProjectEndpoint is called", func() {
			It("should return valid function", func() {
				endpoint := sut.DeleteProjectEndpoint()
				Ω(endpoint).ShouldNot(BeNil())
			})

			var (
				endpoint gokitendpoint.Endpoint
				request  business.DeleteProjectRequest
				response business.DeleteProjectResponse
			)

			BeforeEach(func() {
				endpoint = sut.DeleteProjectEndpoint()
				request = business.DeleteProjectRequest{
					UserEmail: cuid.New() + "@test.com",
					ProjectID: cuid.New(),
				}

				response = business.DeleteProjectResponse{}
			})

			Context("DeleteProjectEndpoint function is returned", func() {
				When("endpoint is called with nil context", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(nil, &request)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteProjectResponse)
						assertArgumentNilError("ctx", "", castedResponse.Err)
					})
				})

				When("endpoint is called with nil request", func() {
					It("should return ArgumentNilError", func() {
						returnedResponse, err := endpoint(ctx, nil)

						Ω(err).Should(BeNil())
						Ω(returnedResponse).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteProjectResponse)
						assertArgumentNilError("request", "", castedResponse.Err)
					})
				})

				When("endpoint is called with invalid request", func() {
					It("should return ArgumentNilError", func() {
						invalidRequest := business.DeleteProjectRequest{
							ProjectID: "",
						}
						returnedResponse, err := endpoint(ctx, &invalidRequest)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteProjectResponse)
						validationErr := invalidRequest.Validate()
						assertArgumentError("request", validationErr.Error(), castedResponse.Err, validationErr)
					})
				})

				When("endpoint is called with valid request", func() {
					It("should call business service DeleteProject method", func() {
						mockBusinessService.
							EXPECT().
							DeleteProject(ctx, gomock.Any()).
							Do(func(_ context.Context, mappedRequest *business.DeleteProjectRequest) {
								Ω(mappedRequest.ProjectID).Should(Equal(request.ProjectID))
							}).
							Return(&response, nil)

						returnedResponse, err := endpoint(ctx, &request)

						Ω(err).Should(BeNil())
						Ω(response).ShouldNot(BeNil())
						castedResponse := returnedResponse.(*business.DeleteProjectResponse)
						Ω(castedResponse.Err).Should(BeNil())
					})
				})

				When("business service DeleteProject returns error", func() {
					It("should return the same error", func() {
						expectedErr := errors.New(cuid.New())
						mockBusinessService.
							EXPECT().
							DeleteProject(gomock.Any(), gomock.Any()).
							Return(nil, expectedErr)

						_, err := endpoint(ctx, &request)

						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("business service DeleteProject returns response", func() {
					It("should return the same response", func() {
						mockBusinessService.
							EXPECT().
							DeleteProject(gomock.Any(), gomock.Any()).
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
				endpoint   gokitendpoint.Endpoint
				projectIDs []string
				request    business.SearchRequest
				response   business.SearchResponse
			)

			BeforeEach(func() {
				endpoint = sut.SearchEndpoint()
				projectIDs = []string{}
				for idx := 0; idx < rand.Intn(20)+1; idx++ {
					projectIDs = append(projectIDs, cuid.New())
				}

				request = business.SearchRequest{
					UserEmail: cuid.New() + "@test.com",
					Pagination: common.Pagination{
						After:  convertStringToPointer(cuid.New()),
						First:  convertIntToPointer(rand.Intn(1000)),
						Before: convertStringToPointer(cuid.New()),
						Last:   convertIntToPointer(rand.Intn(1000)),
					},
					SortingOptions: []common.SortingOptionPair{
						{
							Name:      cuid.New(),
							Direction: common.Ascending,
						},
						{
							Name:      cuid.New(),
							Direction: common.Descending,
						},
					},
					ProjectIDs: projectIDs,
				}

				projects := []models.ProjectWithCursor{}

				for idx := 0; idx < rand.Intn(20)+1; idx++ {
					projects = append(projects, models.ProjectWithCursor{
						ProjectID: cuid.New(),
						Project: models.Project{
							Name: cuid.New(),
						},
						Cursor: cuid.New(),
					})
				}

				response = business.SearchResponse{
					HasPreviousPage: (rand.Intn(10) % 2) == 0,
					HasNextPage:     (rand.Intn(10) % 2) == 0,
					TotalCount:      rand.Int63n(1000),
					Projects:        projects,
				}
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
								Ω(mappedRequest.ProjectIDs).Should(Equal(request.ProjectIDs))
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

func convertStringToPointer(str string) *string {
	return &str
}

func convertIntToPointer(i int) *int {
	return &i
}
