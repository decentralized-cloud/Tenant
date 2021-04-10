package business_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/decentralized-cloud/project/models"
	"github.com/decentralized-cloud/project/services/business"
	repository "github.com/decentralized-cloud/project/services/repository"
	repsoitoryMock "github.com/decentralized-cloud/project/services/repository/mock"
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
		When("project repository service is not provided and NewBusinessService is called", func() {
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

	Describe("CreateProject", func() {
		var (
			request business.CreateProjectRequest
		)

		BeforeEach(func() {
			request = business.CreateProjectRequest{
				Project: models.Project{
					Name: cuid.New(),
				}}
		})

		Context("project service is instantiated", func() {
			When("CreateProject is called", func() {
				It("should call project repository CreateProject method", func() {
					mockRepositoryService.
						EXPECT().
						CreateProject(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.CreateProjectRequest) {
							Ω(mappedRequest.Project).Should(Equal(request.Project))
						}).
						Return(&repository.CreateProjectResponse{}, nil)

					response, err := sut.CreateProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})

				When("And project repository CreateProject returns error", func() {
					It("should return the same error", func() {
						expectedError := errors.New(cuid.New())
						mockRepositoryService.
							EXPECT().
							CreateProject(gomock.Any(), gomock.Any()).
							Return(nil, expectedError)

						response, err := sut.CreateProject(ctx, &request)
						Ω(err).Should(BeNil())
						Ω(response.Err).Should(Equal(expectedError))
					})
				})

				When("And project repository CreateProject return no error", func() {
					It("should return expected details", func() {
						expectedResponse := repository.CreateProjectResponse{
							ProjectID: cuid.New(),
							Project: models.Project{
								Name: cuid.New(),
							},
							Cursor: cuid.New(),
						}

						mockRepositoryService.
							EXPECT().
							CreateProject(gomock.Any(), gomock.Any()).
							Return(&expectedResponse, nil)

						response, err := sut.CreateProject(ctx, &request)
						Ω(err).Should(BeNil())
						Ω(response.Err).Should(BeNil())
						Ω(response.ProjectID).ShouldNot(BeNil())
						Ω(response.ProjectID).Should(Equal(expectedResponse.ProjectID))
						Ω(response.Project).Should(Equal(expectedResponse.Project))
					})
				})
			})
		})
	})

	Describe("ReadProject", func() {
		var (
			request business.ReadProjectRequest
		)

		BeforeEach(func() {
			request = business.ReadProjectRequest{
				ProjectID: cuid.New(),
			}
		})

		Context("project service is instantiated", func() {
			When("ReadProject is called", func() {
				It("should call project repository ReadProject method", func() {
					mockRepositoryService.
						EXPECT().
						ReadProject(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.ReadProjectRequest) {
							Ω(mappedRequest.ProjectID).Should(Equal(request.ProjectID))
						}).
						Return(&repository.ReadProjectResponse{}, nil)

					response, err := sut.ReadProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And project repository ReadProject returns error", func() {
				It("should return the same error", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						ReadProject(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ReadProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(Equal(expectedError))
				})
			})

			When("And project repository ReadProject return no error", func() {
				It("should return the project details", func() {
					expectedResponse := repository.ReadProjectResponse{
						Project: models.Project{Name: cuid.New()},
					}

					mockRepositoryService.
						EXPECT().
						ReadProject(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.ReadProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					Ω(response.Project).Should(Equal(expectedResponse.Project))
				})
			})
		})
	})

	Describe("UpdateProject", func() {
		var (
			request business.UpdateProjectRequest
		)

		BeforeEach(func() {
			request = business.UpdateProjectRequest{
				ProjectID: cuid.New(),
				Project:   models.Project{Name: cuid.New()},
			}
		})

		Context("project service is instantiated", func() {
			When("UpdateProject is called", func() {
				It("should call project repository UpdateProject method", func() {
					mockRepositoryService.
						EXPECT().
						UpdateProject(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.UpdateProjectRequest) {
							Ω(mappedRequest.ProjectID).Should(Equal(request.ProjectID))
							Ω(mappedRequest.Project.Name).Should(Equal(request.Project.Name))
						}).
						Return(&repository.UpdateProjectResponse{}, nil)

					response, err := sut.UpdateProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("And project repository UpdateProject returns error", func() {
				It("should return the same error", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						UpdateProject(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.UpdateProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(Equal(expectedError))
				})
			})

			When("And project repository UpdateProject return no error", func() {
				It("should return expected details", func() {
					expectedResponse := repository.UpdateProjectResponse{
						Project: models.Project{
							Name: cuid.New(),
						},
						Cursor: cuid.New(),
					}
					mockRepositoryService.
						EXPECT().
						UpdateProject(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.UpdateProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					Ω(response.Project).Should(Equal(expectedResponse.Project))
				})
			})
		})
	})

	Describe("DeleteProject is called", func() {
		var (
			request business.DeleteProjectRequest
		)

		BeforeEach(func() {
			request = business.DeleteProjectRequest{
				ProjectID: cuid.New(),
			}
		})

		Context("project service is instantiated", func() {
			When("DeleteProject is called", func() {
				It("should call project repository DeleteProject method", func() {
					mockRepositoryService.
						EXPECT().
						DeleteProject(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.DeleteProjectRequest) {
							Ω(mappedRequest.ProjectID).Should(Equal(request.ProjectID))
						}).
						Return(&repository.DeleteProjectResponse{}, nil)

					response, err := sut.DeleteProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("project repository DeleteProject returns error", func() {
				It("should return the same error", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						DeleteProject(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.DeleteProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(Equal(expectedError))
				})
			})

			When("project repository DeleteProject completes successfully", func() {
				It("should return no error", func() {
					mockRepositoryService.
						EXPECT().
						DeleteProject(gomock.Any(), gomock.Any()).
						Return(&repository.DeleteProjectResponse{}, nil)

					response, err := sut.DeleteProject(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})
		})
	})

	Describe("ListProjects is called", func() {
		var (
			request    business.ListProjectsRequest
			projectIDs []string
		)

		BeforeEach(func() {
			projectIDs = []string{}
			for idx := 0; idx < rand.Intn(20)+1; idx++ {
				projectIDs = append(projectIDs, cuid.New())
			}

			request = business.ListProjectsRequest{
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
				ProjectIDs: projectIDs,
			}
		})

		Context("project service is instantiated", func() {
			When("ListProjects is called", func() {
				It("should call project repository ListProjects method", func() {
					mockRepositoryService.
						EXPECT().
						ListProjects(ctx, gomock.Any()).
						Do(func(_ context.Context, mappedRequest *repository.ListProjectsRequest) {
							Ω(mappedRequest.Pagination).Should(Equal(request.Pagination))
							Ω(mappedRequest.SortingOptions).Should(Equal(request.SortingOptions))
							Ω(mappedRequest.ProjectIDs).Should(Equal(request.ProjectIDs))
						}).
						Return(&repository.ListProjectsResponse{}, nil)

					response, err := sut.ListProjects(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
				})
			})

			When("project repository ListProjects returns error", func() {
				It("should return the same error", func() {
					expectedError := errors.New(cuid.New())
					mockRepositoryService.
						EXPECT().
						ListProjects(gomock.Any(), gomock.Any()).
						Return(nil, expectedError)

					response, err := sut.ListProjects(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(Equal(expectedError))
				})
			})

			When("project repository ListProjects completes successfully", func() {
				It("should return the list of matched projectIDs", func() {
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

					expectedResponse := repository.ListProjectsResponse{
						HasPreviousPage: (rand.Intn(10) % 2) == 0,
						HasNextPage:     (rand.Intn(10) % 2) == 0,
						TotalCount:      rand.Int63n(1000),
						Projects:        projects,
					}

					mockRepositoryService.
						EXPECT().
						ListProjects(gomock.Any(), gomock.Any()).
						Return(&expectedResponse, nil)

					response, err := sut.ListProjects(ctx, &request)
					Ω(err).Should(BeNil())
					Ω(response.Err).Should(BeNil())
					Ω(response.HasPreviousPage).Should(Equal(expectedResponse.HasPreviousPage))
					Ω(response.HasNextPage).Should(Equal(expectedResponse.HasNextPage))
					Ω(response.TotalCount).Should(Equal(expectedResponse.TotalCount))
					Ω(response.Projects).Should(Equal(expectedResponse.Projects))
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

func convertStringToPointer(str string) *string {
	return &str
}

func convertIntToPointer(i int) *int {
	return &i
}
