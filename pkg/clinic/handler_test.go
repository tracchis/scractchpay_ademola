package clinic

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name             string
		setupFetcherMock func(mock *DataFetcherMock)
		wantCode         int
		wantBody         string
	}{

		{
			name: "error while fetching data",
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).Return(nil, errors.New("random network error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"error\":\"error fetching all clinics\",\"messages\":{}}\n",
		},
		{
			name: "check if user exist success",
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).Return([]Clinic{
					{
						Name:  "Scratchpay Official practice",
						State: "FL",
						Availability: Availability{
							From: "09:00",
							To:   "20:00",
						},
					},
				}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[{\"name\":\"Scratchpay Official practice\",\"state\":\"FL\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}}]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcherMock := &DataFetcherMock{}
			if tt.setupFetcherMock != nil {
				tt.setupFetcherMock(fetcherMock)
			}

			request := httptest.NewRequest(http.MethodGet, "http://www.test.com/v1/clinics/", nil)
			response := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/v1/clinics/", GetAllClinics(fetcherMock))
			r.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)

			assert.True(t, fetcherMock.AssertExpectations(t))
			assert.Equal(t, tt.wantBody, string(body))
			assert.Equal(t, tt.wantCode, response.Code)
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		setupFetcherMock func(mock *DataFetcherMock)
		wantCode         int
		wantBody         string
	}{
		{
			name:     "invalid request body",
			body:     `{`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"error\":\"invalid json params\",\"messages\":{}}\n",
		},
		{
			name:     "invalid params payload",
			body:     `{"name":"Good health", "invalid_key": "sample"`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"error\":\"invalid json params\",\"messages\":{}}\n",
		},
		{
			name: "error while searching",
			body: `{"name": "Good ","state": "FL"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).Return(nil,
					errors.New("random error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"error\":\"error fetching all clinics\",\"messages\":{}}\n",
		},
		{
			name: "search returns no match",
			body: `{"name": "Good ","state": "FL"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).
					Return([]Clinic{
						{
							Name:  "Scratchpay Official practice",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[]\n",
		},
		{
			name: "search matches by name",
			body: `{"name": "Scratchpay Official practice"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).
					Return([]Clinic{
						{
							Name:  "Scratchpay Official practice",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
						{
							Name:  "Good Health",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[{\"name\":\"Scratchpay Official practice\",\"state\":\"FL\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}}]\n",
		},
		{
			name: "search matches by state",
			body: `{"state": "California"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).
					Return([]Clinic{
						{
							Name:  "Scratchpay Official practice",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
						{
							Name:  "Good Health",
							State: "California",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[{\"name\":\"Good Health\",\"state\":\"California\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}}]\n",
		},
		{
			name: "search fails when name and state don't match ",
			body: `{"state": "FL", "name": "Good Health"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).
					Return([]Clinic{
						{
							Name:  "Scratchpay Official practice",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
						{
							Name:  "Good Health",
							State: "California",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[]\n",
		},
		{
			name: "search matches by name & state",
			body: `{"state": "California", "name": "Good Health"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).
					Return([]Clinic{
						{
							Name:  "Scratchpay Official practice",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
						{
							Name:  "Good Health",
							State: "California",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[{\"name\":\"Good Health\",\"state\":\"California\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}}]\n",
		},
		{
			name: "search matches by availability (from & to)",
			body: `{"from": "09:00", "to": "20:00"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).
					Return([]Clinic{
						{
							Name:  "Scratchpay Official practice",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
						{
							Name:  "Good Health",
							State: "California",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[{\"name\":\"Scratchpay Official practice\",\"state\":\"FL\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}},{\"name\":\"Good Health\",\"state\":\"California\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}}]\n",
		},
		{
			name: "search matches by availability within range",
			body: `{"from": "11:00", "to": "16:00"}`,
			setupFetcherMock: func(mock *DataFetcherMock) {
				mock.On("GetClinicData", m.Anything).
					Return([]Clinic{
						{
							Name:  "Scratchpay Official practice",
							State: "FL",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
						{
							Name:  "Good Health",
							State: "California",
							Availability: Availability{
								From: "09:00",
								To:   "20:00",
							},
						},
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: "[{\"name\":\"Scratchpay Official practice\",\"state\":\"FL\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}},{\"name\":\"Good Health\",\"state\":\"California\",\"availability\":{\"from\":\"09:00\",\"to\":\"20:00\"}}]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fetcherMock := &DataFetcherMock{}
			if tt.setupFetcherMock != nil {
				tt.setupFetcherMock(fetcherMock)
			}

			request := httptest.NewRequest(http.MethodPost, "http://www.test.com/v1/clinics/search", strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/v1/clinics/search", Search(fetcherMock))
			r.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)

			assert.True(t, fetcherMock.AssertExpectations(t))
			assert.Equal(t, tt.wantBody, string(body))
			assert.Equal(t, tt.wantCode, response.Code)
		})
	}
}
