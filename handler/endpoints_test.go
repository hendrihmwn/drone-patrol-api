package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServer_PostEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	s := NewServer(NewServerOptions{
		Repository: mockRepository,
	})

	e := echo.New()
	id := uuid.New().String()

	testCases := []struct {
		name           string
		requestBody    map[string]int
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "BAD REQUEST",
			requestBody: map[string]int{
				"width":  -1,
				"length": 20,
			},
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Key: 'CreateEstateRequest.Width' Error:Field validation for 'Width' failed on the 'gte' tag"}`,
		},
		{
			name: "INTERNAL_SERVER_ERROR",
			requestBody: map[string]int{
				"width":  10,
				"length": 20,
			},
			setupMocks: func() {
				mockRepository.EXPECT().CreateEstate(gomock.Any(), repository.Estate{
					Width:  10,
					Length: 20,
				}).Return(repository.Estate{
					Id:        id,
					Width:     10,
					Length:    20,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, errors.New(""))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":""}`,
		},
		{
			name: "OK",
			requestBody: map[string]int{
				"width":  10,
				"length": 20,
			},
			setupMocks: func() {
				mockRepository.EXPECT().CreateEstate(gomock.Any(), repository.Estate{
					Width:  10,
					Length: 20,
				}).Return(repository.Estate{
					Id:        id,
					Width:     10,
					Length:    20,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`{"id":"%s"}`, id),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}
			e.POST("/estate", s.PostEstate)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSuffix(rec.Body.String(), "\n"))
		})
	}
}

func TestServer_GetEstateIdStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	s := NewServer(NewServerOptions{
		Repository: mockRepository,
	})
	e := echo.New()
	id := uuid.New().String()

	testCases := []struct {
		name           string
		requestId      string
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "BAD_REQUEST",
			requestId: "11",
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Key: 'IdPath.ID' Error:Field validation for 'ID' failed on the 'uuid4' tag"}`,
		},
		{
			name:      "ESTATE_NOT_FOUND",
			requestId: id,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     10,
					Length:    20,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, errors.New("sql: no rows in result set"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"estate is not found"}`,
		},
		{
			name:      "INTERNAL_SERVER_ERROR",
			requestId: id,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     10,
					Length:    20,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, errors.New(""))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":""}`,
		},
		{
			name:      "INTERNAL_SERVER_ERROR",
			requestId: id,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     10,
					Length:    20,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().ListTreesByEstateId(gomock.Any(), repository.ListTreesByEstateIdInput{
					EstateId: id,
				}).Return([]repository.Tree{{
					Id:        id,
					X:         5,
					Y:         1,
					Height:    15,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}}, errors.New(""))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":""}`,
		},
		{
			name:      "OK_WITH_ZERO_STATS",
			requestId: id,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     10,
					Length:    20,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().ListTreesByEstateId(gomock.Any(), repository.ListTreesByEstateIdInput{
					EstateId: id,
				}).Return([]repository.Tree{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"count":0,"max":0,"median":0,"min":0}`,
		},
		{
			name:      "OK",
			requestId: id,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     10,
					Length:    20,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().ListTreesByEstateId(gomock.Any(), repository.ListTreesByEstateIdInput{
					EstateId: id,
				}).Return([]repository.Tree{
					{
						Id:        uuid.New().String(),
						EstateId:  id,
						X:         2,
						Y:         1,
						Height:    5,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						Id:        uuid.New().String(),
						EstateId:  id,
						X:         3,
						Y:         1,
						Height:    2,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						Id:        uuid.New().String(),
						EstateId:  id,
						X:         4,
						Y:         1,
						Height:    10,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"count":3,"max":10,"median":3,"min":2}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			e.GET("/estate/:id/stats", func(c echo.Context) error {
				return s.GetEstateIdStats(c, tc.requestId)
			})

			req := httptest.NewRequest(http.MethodGet, "/estate/"+tc.requestId+"/stats", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSuffix(rec.Body.String(), "\n"))
		})
	}
}

func TestServer_PostEstateIdTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	s := NewServer(NewServerOptions{
		Repository: mockRepository,
	})

	e := echo.New()
	id := uuid.New().String()
	treeId := uuid.New().String()

	testCases := []struct {
		name           string
		pathId         string
		requestBody    map[string]int
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "BAD_REQUEST_VALIDATION_PATH",
			pathId:      "123",
			requestBody: map[string]int{},
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Key: 'IdPath.ID' Error:Field validation for 'ID' failed on the 'uuid4' tag"}`,
		},
		{
			name:   "BAD_REQUEST_VALIDATION_REQUEST",
			pathId: id,
			requestBody: map[string]int{
				"x":      5,
				"y":      1,
				"height": -1,
			},
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Key: 'CreateTreeRequest.Height' Error:Field validation for 'Height' failed on the 'gte' tag"}`,
		},
		{
			name:   "ESTATE_NOT_FOUND",
			pathId: id,
			requestBody: map[string]int{
				"x":      5,
				"y":      1,
				"height": 10,
			},
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{}, errors.New("sql: no rows in result set"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"estate is not found"}`,
		},
		{
			name:   "INTERNAL_SERVER_ERROR",
			pathId: id,
			requestBody: map[string]int{
				"x":      5,
				"y":      1,
				"height": 10,
			},
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{}, errors.New(""))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":""}`,
		},
		{
			name:   "BAD_REQUEST_INDEX_OUT_OF_BOUND",
			pathId: id,
			requestBody: map[string]int{
				"x":      5,
				"y":      1,
				"height": 10,
			},
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     1,
					Length:    4,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"index out of bound"}`,
		},
		{
			name:   "BAD_REQUEST_PLOT_ALREADY_EXIST",
			pathId: id,
			requestBody: map[string]int{
				"x":      5,
				"y":      1,
				"height": 10,
			},
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     1,
					Length:    5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().GetTreeByPlot(gomock.Any(), repository.GetTreeByPlot{
					EstateId: id,
					X:        5,
					Y:        1,
				}).Return(repository.Tree{
					Id:       uuid.New().String(),
					EstateId: id,
					X:        5,
					Y:        1,
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"plot already exist"}`,
		},
		{
			name:   "BAD_REQUEST",
			pathId: id,
			requestBody: map[string]int{
				"x":      5,
				"y":      1,
				"height": 10,
			},
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     1,
					Length:    5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().GetTreeByPlot(gomock.Any(), repository.GetTreeByPlot{
					EstateId: id,
					X:        5,
					Y:        1,
				}).Return(repository.Tree{}, errors.New("not exist"))
				mockRepository.EXPECT().CreateTree(gomock.Any(), repository.Tree{
					EstateId: id,
					X:        5,
					Y:        1,
					Height:   10,
				}).Return(repository.Tree{}, errors.New(""))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":""}`,
		},
		{
			name:   "OK",
			pathId: id,
			requestBody: map[string]int{
				"x":      5,
				"y":      1,
				"height": 10,
			},
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     1,
					Length:    5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().GetTreeByPlot(gomock.Any(), repository.GetTreeByPlot{
					EstateId: id,
					X:        5,
					Y:        1,
				}).Return(repository.Tree{}, errors.New("not exist"))
				mockRepository.EXPECT().CreateTree(gomock.Any(), repository.Tree{
					EstateId: id,
					X:        5,
					Y:        1,
					Height:   10,
				}).Return(repository.Tree{
					Id:       treeId,
					EstateId: id,
					X:        5,
					Y:        1,
					Height:   10,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`{"id":"%s"}`, treeId),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			e.POST("/estate/:id/tree", func(c echo.Context) error {
				return s.PostEstateIdTree(c, tc.pathId)
			})

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/estate/"+tc.pathId+"/tree", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSuffix(rec.Body.String(), "\n"))
		})
	}
}

func TestServer_GetEstateIdDronePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	s := NewServer(NewServerOptions{
		Repository: mockRepository,
	})

	e := echo.New()
	id := uuid.New().String()
	//treeId := uuid.New().String()
	negInt := -1
	posInt1 := 40
	posInt2 := 90

	testCases := []struct {
		name           string
		pathId         string
		params         *int
		setupMocks     func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "BAD_REQUEST_VALIDATION_PATH",
			pathId: "123",
			params: nil,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Key: 'IdPath.ID' Error:Field validation for 'ID' failed on the 'uuid4' tag"}`,
		},
		{
			name:   "BAD_REQUEST_VALIDATION_PARAMS",
			pathId: id,
			params: &negInt,
			setupMocks: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid max distance"}`,
		},
		{
			name:   "ESTATE_NOT_FOUND",
			pathId: id,
			params: nil,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{}, errors.New("sql: no rows in result set"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"estate is not found"}`,
		},
		{
			name:   "INTERNAL_SERVER_ERROR",
			pathId: id,
			params: nil,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{}, errors.New(""))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":""}`,
		},
		{
			name:   "INTERNAL_SERVER_ERROR",
			pathId: id,
			params: nil,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     1,
					Length:    4,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().ListTreesByEstateId(gomock.Any(), repository.ListTreesByEstateIdInput{
					EstateId: id,
				}).Return([]repository.Tree{}, errors.New(""))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":""}`,
		},
		{
			name:   "OK",
			pathId: id,
			params: nil,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     2,
					Length:    5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().ListTreesByEstateId(gomock.Any(), repository.ListTreesByEstateIdInput{
					EstateId: id,
				}).Return([]repository.Tree{{
					Id:       uuid.New().String(),
					EstateId: id,
					X:        3,
					Y:        1,
					Height:   5,
				}, {
					Id:       uuid.New().String(),
					EstateId: id,
					X:        3,
					Y:        2,
					Height:   5,
				}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"distance":112,"rest":{"x":1,"y":2}}`,
		},
		{
			name:   "OK_WITH_MAX_DISTANCE_40",
			pathId: id,
			params: &posInt1,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     2,
					Length:    5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().ListTreesByEstateId(gomock.Any(), repository.ListTreesByEstateIdInput{
					EstateId: id,
				}).Return([]repository.Tree{{
					Id:       uuid.New().String(),
					EstateId: id,
					X:        3,
					Y:        1,
					Height:   5,
				}, {
					Id:       uuid.New().String(),
					EstateId: id,
					X:        3,
					Y:        2,
					Height:   5,
				}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"distance":40,"rest":{"x":3,"y":1}}`,
		},
		{
			name:   "OK_WITH_MAX_DISTANCE_90",
			pathId: id,
			params: &posInt2,
			setupMocks: func() {
				mockRepository.EXPECT().GetEstateById(gomock.Any(), repository.GetEstateByIdInput{
					Id: id,
				}).Return(repository.Estate{
					Id:        id,
					Width:     2,
					Length:    5,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRepository.EXPECT().ListTreesByEstateId(gomock.Any(), repository.ListTreesByEstateIdInput{
					EstateId: id,
				}).Return([]repository.Tree{{
					Id:       uuid.New().String(),
					EstateId: id,
					X:        3,
					Y:        1,
					Height:   5,
				}, {
					Id:       uuid.New().String(),
					EstateId: id,
					X:        3,
					Y:        2,
					Height:   5,
				}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"distance":90,"rest":{"x":3,"y":2}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			e.GET("/estate/:id/drone-plan", func(c echo.Context) error {
				return s.GetEstateIdDronePlan(c, tc.pathId, generated.GetEstateIdDronePlanParams{
					MaxDistance: tc.params,
				})
			})

			req := httptest.NewRequest(http.MethodGet, "/estate/"+tc.pathId+"/drone-plan", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimSuffix(rec.Body.String(), "\n"))
		})
	}
}
