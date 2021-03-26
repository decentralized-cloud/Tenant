package mongodb_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/decentralized-cloud/project/models"
	configurationMock "github.com/decentralized-cloud/project/services/configuration/mock"
	"github.com/decentralized-cloud/project/services/repository"
	"github.com/decentralized-cloud/project/services/repository/mongodb"
	"github.com/golang/mock/gomock"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMongodbRepositoryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Repository Service Tests")
}

var _ = Describe("Mongodb Repository Service Tests", func() {
	var (
		mockCtrl      *gomock.Controller
		sut           repository.RepositoryContract
		ctx           context.Context
		createRequest repository.CreateProjectRequest
	)

	BeforeEach(func() {
		connectionString := os.Getenv("DATABASE_CONNECTION_STRING")
		if strings.Trim(connectionString, " ") == "" {
			connectionString = "mongodb://mongodb:27017"
		}

		mockCtrl = gomock.NewController(GinkgoT())
		mockConfigurationService := configurationMock.NewMockConfigurationContract(mockCtrl)
		mockConfigurationService.
			EXPECT().
			GetDatabaseConnectionString().
			Return(connectionString, nil)

		mockConfigurationService.
			EXPECT().
			GetDatabaseName().
			Return("project", nil)

		mockConfigurationService.
			EXPECT().
			GetDatabaseCollectionName().
			Return("project", nil)

		sut, _ = mongodb.NewMongodbRepositoryService(mockConfigurationService)
		ctx = context.Background()
		createRequest = repository.CreateProjectRequest{
			UserEmail: cuid.New() + "@test.com",
			Project: models.Project{
				Name: cuid.New(),
			}}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("user tries to instantiate RepositoryService", func() {
		When("all dependencies are resolved and NewRepositoryService is called", func() {
			It("should instantiate the new RepositoryService", func() {
				mockConfigurationService := configurationMock.NewMockConfigurationContract(mockCtrl)
				mockConfigurationService.
					EXPECT().
					GetDatabaseConnectionString().
					Return(cuid.New(), nil)

				mockConfigurationService.
					EXPECT().
					GetDatabaseName().
					Return(cuid.New(), nil)

				mockConfigurationService.
					EXPECT().
					GetDatabaseCollectionName().
					Return(cuid.New(), nil)

				service, err := mongodb.NewMongodbRepositoryService(mockConfigurationService)
				Ω(err).Should(BeNil())
				Ω(service).ShouldNot(BeNil())
			})
		})
	})

	Context("user going to create a new project", func() {
		When("create project is called", func() {
			It("should create the new project", func() {
				response, err := sut.CreateProject(ctx, &createRequest)
				Ω(err).Should(BeNil())
				Ω(response.ProjectID).ShouldNot(BeNil())
				Ω(response.Cursor).Should(Equal(response.ProjectID))
				assertProject(response.Project, createRequest.Project)
			})
		})
	})

	Context("project already exists", func() {
		var (
			projectID string
		)

		BeforeEach(func() {
			response, _ := sut.CreateProject(ctx, &createRequest)
			projectID = response.ProjectID
		})

		When("user reads a project by Id", func() {
			It("should return a project", func() {
				response, err := sut.ReadProject(ctx, &repository.ReadProjectRequest{UserEmail: createRequest.UserEmail, ProjectID: projectID})
				Ω(err).Should(BeNil())
				assertProject(response.Project, createRequest.Project)
			})
		})

		When("user updates the existing project", func() {
			It("should update the project information", func() {
				updateRequest := repository.UpdateProjectRequest{
					UserEmail: createRequest.UserEmail,
					ProjectID: projectID,
					Project: models.Project{
						Name: cuid.New(),
					}}

				updateResponse, err := sut.UpdateProject(ctx, &updateRequest)
				Ω(err).Should(BeNil())
				Ω(updateResponse.Cursor).Should(Equal(projectID))
				assertProject(updateResponse.Project, updateRequest.Project)

				readResponse, err := sut.ReadProject(ctx, &repository.ReadProjectRequest{UserEmail: createRequest.UserEmail, ProjectID: projectID})
				Ω(err).Should(BeNil())
				assertProject(readResponse.Project, updateRequest.Project)
			})
		})

		When("user deletes the project", func() {
			It("should delete the project", func() {
				_, err := sut.DeleteProject(ctx, &repository.DeleteProjectRequest{UserEmail: createRequest.UserEmail, ProjectID: projectID})
				Ω(err).Should(BeNil())

				response, err := sut.ReadProject(ctx, &repository.ReadProjectRequest{UserEmail: createRequest.UserEmail, ProjectID: projectID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(commonErrors.IsNotFoundError(err)).Should(BeTrue())

				var notFoundErr commonErrors.NotFoundError
				_ = errors.As(err, &notFoundErr)
			})
		})
	})

	Context("project does not exist", func() {
		var (
			projectID string
		)

		BeforeEach(func() {
			projectID = cuid.New()
		})

		When("user reads the project", func() {
			It("should return NotFoundError", func() {
				response, err := sut.ReadProject(ctx, &repository.ReadProjectRequest{UserEmail: createRequest.UserEmail, ProjectID: projectID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(commonErrors.IsNotFoundError(err)).Should(BeTrue())

				var notFoundErr commonErrors.NotFoundError
				_ = errors.As(err, &notFoundErr)
			})
		})

		When("user tries to update the project", func() {
			It("should return NotFoundError", func() {
				updateRequest := repository.UpdateProjectRequest{
					UserEmail: createRequest.UserEmail,
					ProjectID: projectID,
					Project: models.Project{
						Name: cuid.New(),
					}}

				response, err := sut.UpdateProject(ctx, &updateRequest)
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(commonErrors.IsNotFoundError(err)).Should(BeTrue())

				var notFoundErr commonErrors.NotFoundError
				_ = errors.As(err, &notFoundErr)
			})
		})

		When("user tries to delete the project", func() {
			It("should return NotFoundError", func() {
				response, err := sut.DeleteProject(ctx, &repository.DeleteProjectRequest{UserEmail: createRequest.UserEmail, ProjectID: projectID})
				Ω(err).Should(HaveOccurred())
				Ω(response).Should(BeNil())

				Ω(commonErrors.IsNotFoundError(err)).Should(BeTrue())

				var notFoundErr commonErrors.NotFoundError
				_ = errors.As(err, &notFoundErr)
			})
		})
	})

	Context("project already exists", func() {
		var (
			projectIDs []string
		)

		BeforeEach(func() {
			projectIDs = []string{}

			for i := 0; i < 10; i++ {
				projectName := fmt.Sprintf("%s%d", "Name", i)
				createRequest.Project.Name = projectName
				response, _ := sut.CreateProject(ctx, &createRequest)
				projectIDs = append(projectIDs, response.ProjectID)
			}
		})

		When("user searches for projects with selected project Ids and first 10 projects", func() {
			It("should return first 10 projects", func() {
				first := 10
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Projects[i].ProjectID).Should(Equal(projectIDs[i]))
					projectName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Projects[i].Project.Name).Should(Equal(projectName))
				}
			})
		})

		When("user searches for projects with selected project Ids and first 5 projects", func() {
			It("should return first 5 projects", func() {
				first := 5
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(5))
				for i := 0; i < 5; i++ {
					Ω(response.Projects[i].ProjectID).Should(Equal(projectIDs[i]))
					projectName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Projects[i].Project.Name).Should(Equal(projectName))
				}
			})
		})

		When("user searches for projects with selected project Ids with After parameter provided.", func() {
			It("should return first 9 projects after provided project id", func() {
				first := 9
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						After: &projectIDs[0],
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(9))
				for i := 1; i < 10; i++ {
					Ω(response.Projects[i-1].ProjectID).Should(Equal(projectIDs[i]))
					projectName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Projects[i-1].Project.Name).Should(Equal(projectName))
				}
			})
		})

		When("user searches for projects with selected project Ids with After parameter provided.", func() {
			It("should return first 5 projects after provided project id", func() {
				first := 5
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						After: &projectIDs[0],
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(5))
				for i := 1; i < 5; i++ {
					Ω(response.Projects[i-1].ProjectID).Should(Equal(projectIDs[i]))
					projectName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Projects[i-1].Project.Name).Should(Equal(projectName))
				}
			})
		})

		When("user searches for projects with selected project Ids and last 10 projects", func() {
			It("should return first 10 projects", func() {
				last := 10
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						Before: nil,
						Last:   &last,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Projects[i].ProjectID).Should(Equal(projectIDs[i]))
					projectName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Projects[i].Project.Name).Should(Equal(projectName))
				}
			})
		})

		When("user searches for projects with selected project Ids with Before parameter provided.", func() {
			It("should return first 9 projects before provided project id", func() {
				last := 9
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						Before: &projectIDs[9],
						Last:   &last,
					},
					SortingOptions: []common.SortingOptionPair{},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(9))
				for i := 0; i < 9; i++ {
					Ω(response.Projects[i].ProjectID).Should(Equal(projectIDs[i]))
					projectName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Projects[i].Project.Name).Should(Equal(projectName))
				}
			})
		})

		When("user searches for projects with selected project Ids and first 10 projects with ascending order on name property", func() {
			It("should return first 10 projects in adcending order on name field", func() {
				first := 10
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{
						{Name: "name", Direction: common.Ascending},
					},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Projects[i].ProjectID).Should(Equal(projectIDs[i]))
					projectName := fmt.Sprintf("%s%d", "Name", i)
					Ω(response.Projects[i].Project.Name).Should(Equal(projectName))
				}
			})
		})

		When("user searches for projects with selected project Ids and first 10 projects with descending order on name property", func() {
			It("should return first 10 projects in descending order on name field", func() {
				first := 10
				searchRequest := repository.SearchRequest{
					UserEmail:  createRequest.UserEmail,
					ProjectIDs: projectIDs,
					Pagination: common.Pagination{
						After: nil,
						First: &first,
					},
					SortingOptions: []common.SortingOptionPair{
						{Name: "name", Direction: common.Descending},
					},
				}

				response, err := sut.Search(ctx, &searchRequest)
				Ω(err).Should(BeNil())
				Ω(response.Projects).ShouldNot(BeNil())
				Ω(len(response.Projects)).Should(Equal(10))
				for i := 0; i < 10; i++ {
					Ω(response.Projects[i].ProjectID).Should(Equal(projectIDs[9-i]))
					projectName := fmt.Sprintf("%s%d", "Name", 9-i)
					Ω(response.Projects[i].Project.Name).Should(Equal(projectName))
				}
			})
		})

	})

})

func assertProject(project, expectedProject models.Project) {
	Ω(project).ShouldNot(BeNil())
	Ω(project.Name).Should(Equal(expectedProject.Name))
}
