package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"

	"github.com/egor-denisov/wallet-infotecs/internal/entity"
	mock_usecase "github.com/egor-denisov/wallet-infotecs/internal/usecase/mocks"
	"github.com/egor-denisov/wallet-infotecs/pkg/logger"
)

func Test_createNewWallet(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_usecase.MockWallet, id string)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mock_usecase.MockWallet, name string) {
				r.EXPECT().CreateNewWalletWithDefaultBalance(context.Background()).Return(&entity.Wallet{
					ID: "5b53700ed469fa6a09ea72bb78f36fd9",
					Balance: 100.0,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"id":"5b53700ed469fa6a09ea72bb78f36fd9","balance":100}`,
		},
		{
			name: "Something went wrong",
			mockBehavior: func(r *mock_usecase.MockWallet, name string) {
				r.EXPECT().CreateNewWalletWithDefaultBalance(context.Background()).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.name)
			handler := walletRoutes{
				w: repo,
				l: logger.New(""),
			}
			// Init Endpoint
			r := gin.New()
			r.POST("/", handler.createNewWallet)
			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", nil)
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func Test_sendFunds(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest)

	tests := []struct {
		name                 string
		id                   string
		transactionRequest   entity.TransactionRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			transactionRequest: entity.TransactionRequest{
				To: "eb376add88bf8e70f80787266a0801d5",
				Amount: 100.0,
			},
			mockBehavior: func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest) {
				r.EXPECT().SendFunds(context.Background(), id, transactionRequest.To, transactionRequest.Amount).Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: "",
		},
		{
			name: "Not found",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			transactionRequest: entity.TransactionRequest{
				To: "eb376add88bf8e70f80787266a0801d5",
				Amount: 100.0,
			},
			mockBehavior: func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest) {
				r.EXPECT().SendFunds(context.Background(), id, transactionRequest.To, transactionRequest.Amount).Return(entity.ErrWalletNotFound)
			},
			expectedStatusCode: 404,
			expectedResponseBody: "",
		},
		{
			name: "Wrong input - without reciever id",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			transactionRequest: entity.TransactionRequest{
				Amount: 100.0,
			},
			mockBehavior: func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest) {
				r.EXPECT().SendFunds(context.Background(), id, transactionRequest.To, transactionRequest.Amount).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: "",
		},
		{
			name: "Wrong input - without amount",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			transactionRequest: entity.TransactionRequest{
				To: "eb376add88bf8e70f80787266a0801d5",
			},
			mockBehavior: func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest) {
				r.EXPECT().SendFunds(context.Background(), id, transactionRequest.To, transactionRequest.Amount).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: "",
		},
		{
			name: "Wrong input - empty request body",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			transactionRequest: entity.TransactionRequest{},
			mockBehavior: func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest) {
				r.EXPECT().SendFunds(context.Background(), id, transactionRequest.To, transactionRequest.Amount).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: "",
		},
		{
			name: "Sender is reciever",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			transactionRequest: entity.TransactionRequest{
				To: "5b53700ed469fa6a09ea72bb78f36fd9",
				Amount: 100.0,
			},
			mockBehavior: func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest) {
				r.EXPECT().SendFunds(context.Background(), id, transactionRequest.To, transactionRequest.Amount).Return(entity.ErrSenderIsReceiver)
			},
			expectedStatusCode: 400,
			expectedResponseBody: "",
		},
		{
			name: "Something went wrong",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			transactionRequest: entity.TransactionRequest{
				To: "eb376add88bf8e70f80787266a0801d5",
				Amount: 100.0,
			},
			mockBehavior: func(r *mock_usecase.MockWallet, id string, transactionRequest entity.TransactionRequest) {
				r.EXPECT().SendFunds(context.Background(), id, transactionRequest.To, transactionRequest.Amount).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.id, test.transactionRequest)
			handler := walletRoutes{
				w: repo,
				l: logger.New(""),
			}
			// Init Endpoint
			r := gin.New()
			r.POST("/:walletId/send", handler.sendFunds)
			// Create Request
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(test.transactionRequest)
			req := httptest.NewRequest("POST", fmt.Sprintf("/%s/send", test.id), bytes.NewBuffer(reqBody))
			
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func Test_getWalletHistoryById(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_usecase.MockWallet, id string)

	tests := []struct {
		name                 string
		id                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok - history exists",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			mockBehavior: func(r *mock_usecase.MockWallet, id string) {
				t, _ := time.Parse(time.RFC3339, "2024-02-04T17:25:35.448Z")

				r.EXPECT().GetWalletHistoryById(context.Background(), id).Return([]entity.Transaction{
					{
						Time: t,
						From: "5b53700ed469fa6a09ea72bb78f36fd9",
						To: "eb376add88bf8e70f80787266a0801d5",
						Amount: 30.0,
					},
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[{"time":"2024-02-04T17:25:35.448Z","from":"5b53700ed469fa6a09ea72bb78f36fd9","to":"eb376add88bf8e70f80787266a0801d5","amount":30}]`,
		},
		{
			name: "Ok - history is empty",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			mockBehavior: func(r *mock_usecase.MockWallet, id string) {
				r.EXPECT().GetWalletHistoryById(context.Background(), id).Return([]entity.Transaction{}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[]`,
		},
		{
			name: "Not Found",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			mockBehavior: func(r *mock_usecase.MockWallet, id string) {
				r.EXPECT().GetWalletHistoryById(context.Background(), id).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode: 404,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.id)
			handler := walletRoutes{
				w: repo,
				l: logger.New(""),
			}
			// Init Endpoint
			r := gin.New()
			r.GET("/:walletId/history", handler.getWalletHistoryById)
			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/%s/history", test.id), nil)
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func Test_getWalletById(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_usecase.MockWallet, id string)

	tests := []struct {
		name                 string
		id                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			id: "5b53700ed469fa6a09ea72bb78f36fd9",
			mockBehavior: func(r *mock_usecase.MockWallet, id string) {
				r.EXPECT().GetWalletById(context.Background(), id).Return(&entity.Wallet{
					ID: id,
					Balance: 100.0,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{"id":"5b53700ed469fa6a09ea72bb78f36fd9","balance":100}`,
		},
		{
			name: "Not Found",
			id: "abc",
			mockBehavior: func(r *mock_usecase.MockWallet, id string) {
				r.EXPECT().GetWalletById(context.Background(), id).Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode: 404,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockWallet(c)
			test.mockBehavior(repo, test.id)
			handler := walletRoutes{
				w: repo,
				l: logger.New(""),
			}
			// Init Endpoint
			r := gin.New()
			r.GET("/:walletId", handler.getWalletById)
			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/%s", test.id), nil)
			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}