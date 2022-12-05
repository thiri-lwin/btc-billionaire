package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"btc_billionaire/postgres"
	mockdb "btc_billionaire/postgres/mock"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type sendBTCReqTest struct {
	Datetime string  `json:"datetime"`
	Amount   float32 `json:"amount"`
}

func TestSendBTC(t *testing.T) {
	validDate, _ := time.Parse(time.RFC3339, "2020-11-05T13:30:05+03:00")
	_, offset := validDate.Zone()
	btcInfo := postgres.BTCInfo{
		Amount:  decimal.NewFromInt32(10),
		Created: validDate,
		Offset:  offset,
	}

	testCases := []struct {
		name          string
		req           sendBTCReqTest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			req: sendBTCReqTest{
				Amount:   10,
				Datetime: "2020-11-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertBTCInfo(btcInfo).
					Times(1).
					Return(btcInfo, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidAmount",
			req: sendBTCReqTest{
				Amount:   0,
				Datetime: "2020-11-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertBTCInfo(btcInfo).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidDate",
			req: sendBTCReqTest{
				Amount:   10,
				Datetime: "2020-00-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertBTCInfo(btcInfo).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "InternalError",
			req: sendBTCReqTest{
				Amount:   10,
				Datetime: "2020-11-05T13:30:05+03:00",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertBTCInfo(btcInfo).
					Times(1).
					Return(postgres.BTCInfo{}, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodPost, "/api/btc", bytes.NewBuffer(reqJson))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}
