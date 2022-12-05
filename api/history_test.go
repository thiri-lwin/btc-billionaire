package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"btc_billionaire/postgres"
	mockdb "btc_billionaire/postgres/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type btcHistoryTestReq struct {
	StartDate string `json:"startDatetime"`
	EndDate   string `json:"endDatetime"`
}

func TestGetBTCHistory(t *testing.T) {
	validStartDate, _ := time.Parse(time.RFC3339, "2020-11-05T13:30:05+03:00")
	validEndDate, _ := time.Parse(time.RFC3339, "2022-07-05T13:30:05+03:00")

	history := []postgres.BTCInfoResult{
		{
			Amount:  7,
			Created: validStartDate,
		},
		{
			Amount:  10,
			Created: validStartDate.Add(time.Hour * 1),
		},
		{
			Amount:  5,
			Created: validStartDate.Add(time.Hour * 2),
		},
	}

	testCases := []struct {
		name          string
		req           btcHistoryTestReq
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			req: btcHistoryTestReq{
				StartDate: "2020-11-05T13:30:05+03:00",
				EndDate:   "2022-07-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBTCHistory(validStartDate, validEndDate).
					Times(1).
					Return(history, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkResponse(t, recorder.Body)
			},
		},
		{
			name: "NoRecordFound",
			req: btcHistoryTestReq{
				StartDate: "2020-10-05T13:30:05+03:00",
				EndDate:   "2020-10-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				validStartDate, _ := time.Parse(time.RFC3339, "2020-10-05T13:30:05+03:00")
				validEndDate, _ := time.Parse(time.RFC3339, "2020-10-05T13:30:05+03:00")
				store.EXPECT().
					GetBTCHistory(validStartDate, validEndDate).
					Times(1).
					Return([]postgres.BTCInfoResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidDate",
			req: btcHistoryTestReq{
				StartDate: "2020-10-05T13:30:05+03:00",
				EndDate:   "2022-00-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBTCHistory(validStartDate, validEndDate).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			req: btcHistoryTestReq{
				StartDate: "2020-11-05T13:30:05+03:00",
				EndDate:   "2022-07-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBTCHistory(validStartDate, validEndDate).
					Times(1).
					Return([]postgres.BTCInfoResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			reqJson, _ := json.Marshal(tc.req)
			request, err := http.NewRequest(http.MethodPost, "/api/btc/history", bytes.NewBuffer(reqJson))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func checkResponse(t *testing.T, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)

	require.NoError(t, err)

	var history []postgres.BTCInfoResult
	err = json.Unmarshal(data, &history)
	require.NoError(t, err)
	require.Greater(t, len(history), 0)
}
