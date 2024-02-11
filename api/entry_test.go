package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/k-kaddal/bank-go/db/mock"
	db "github.com/k-kaddal/bank-go/db/sqlc"
	"github.com/k-kaddal/bank-go/token"
	"github.com/stretchr/testify/require"
)

func TestCreateEntyAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	amount := int64(10)

	testCases := []struct{
		name			string
		body			gin.H
		setupAuth 		func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs		func(store *mockdb.MockStore)
		checkResponse	func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"account_id": account.ID,
				"amount" : amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)

				arg := db.EntryTxParams{
					AccountID: account.ID,
					Amount: amount,
				}

				store.EXPECT().
					EntryTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse:	func(recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnauthorisedUser",
			body: gin.H{
				"account_id": account.ID,
				"amount" : amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse:	func(recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "MissingAmount",
			body: gin.H{
				"account_id": account.ID,
				"amount" : 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
					EntryTx(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			checkResponse:	func(recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "AccountNotFound",
			body: gin.H{
				"account_id": account.ID,
				"amount" : amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse:	func(recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "AccountError",
			body: gin.H{
				"account_id": account.ID,
				"amount" : amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse:	func(recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "EntryTxError",
			body: gin.H{
				"account_id": account.ID,
				"amount" : amount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
					
				store.EXPECT().
					EntryTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.EntryTxResult{}, sql.ErrConnDone)
			},
			checkResponse:	func(recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			
			testCase.buildStubs(store)
	
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
	
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)
	
			url := "/entries"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
	
			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
