package http

import (
	"context"
	"log"
	"net/http"
	"time"

	transactionv1 "github.com/nullexp/finman-api-gateway/internal/adapter/grpc/transaction/v1"
	"github.com/nullexp/finman-api-gateway/internal/port/model"

	httpapi "github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol"
	"github.com/nullexp/finman-api-gateway/pkg/infrastructure/http/protocol/model/openapi"
)

const TransactionBaseURL = "/transactions"

func NewTransaction(client transactionv1.TransactionServiceClient, parser model.SubjectParser) httpapi.Module {
	return TransactionHandler{client: client, parser: parser}
}

type TransactionHandler struct {
	client transactionv1.TransactionServiceClient
	parser model.SubjectParser
}

func (s TransactionHandler) GetRequestHandlers() []*httpapi.RequestDefinition {
	return []*httpapi.RequestDefinition{
		s.CreateTransaction(),
		s.GetTransactionById(),
		s.GetTransactionsByUserId(),
		s.GetOwnTransactionById(),
		s.GetAllTransactions(),
		s.UpdateTransaction(),
		s.DeleteTransaction(),
		s.GetTransactionsWithPagination(),
	}
}

func (s TransactionHandler) GetBaseURL() string {
	return TransactionBaseURL
}

const (
	TransactionManagement  = "Transaction Management"
	TransactionDescription = "Use these APIs to manage transactions"
)

func (s TransactionHandler) GetTag() openapi.Tag {
	return openapi.Tag{
		Name:        TransactionManagement,
		Description: TransactionDescription,
	}
}

func (s TransactionHandler) CreateTransaction() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:       "",
		Method:      http.MethodPost,
		FreeRoute:   false,
		Dto:         &CreateTransactionRequest{},
		Description: "Note that type can be deposit or withdrawal",
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &CreateTransactionResponse{},
			},
		},
		Handler: func(req httpapi.Request) {

			caller := req.MustGetCaller()
			sub := s.parser.MustParseSubject(caller.GetSubject())
			dto := req.MustGetDTO().(*CreateTransactionRequest)
			resp, err := s.client.CreateTransaction(context.Background(), &transactionv1.CreateTransactionRequest{
				UserId:      sub.UserId,
				Type:        dto.Type,
				Amount:      dto.Amount,
				Description: dto.Description,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusCreated, err, CreateTransactionResponse{
				Id: resp.Id,
			})
		},
	}
}

func (s TransactionHandler) GetTransactionById() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:     "/{id}",
		Method:    http.MethodGet,
		FreeRoute: false,
		Dto:       &GetTransactionByIdRequest{},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetTransactionByIdResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*GetTransactionByIdRequest)
			resp, err := s.client.GetTransactionById(context.Background(), &transactionv1.GetTransactionByIdRequest{
				Id: dto.Id,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusOK, err, GetTransactionByIdResponse{
				Transaction: Transaction{
					Id:          resp.Transaction.Id,
					UserId:      resp.Transaction.UserId,
					Type:        resp.Transaction.Type,
					Amount:      resp.Transaction.Amount,
					Date:        MustParseTime(resp.Transaction.Date),
					Description: resp.Transaction.Description,
					CreatedAt:   MustParseTime(resp.Transaction.CreatedAt),
					UpdatedAt:   MustParseTime(resp.Transaction.UpdatedAt),
				},
			})
		},
	}
}

func (s TransactionHandler) GetTransactionsByUserId() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "/user/{userId}",
		Method:         http.MethodGet,
		FreeRoute:      false,
		AnyPermissions: []string{"ManageTransactions"},
		Dto:            &GetTransactionsByUserIdRequest{},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetTransactionsByUserIdResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*GetTransactionsByUserIdRequest)
			resp, err := s.client.GetTransactionsByUserId(context.Background(), &transactionv1.GetTransactionsByUserIdRequest{
				UserId: dto.UserId,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			transactions := make([]Transaction, len(resp.Transactions))
			for i, txn := range resp.Transactions {
				transactions[i] = Transaction{
					Id:          txn.Id,
					UserId:      txn.UserId,
					Type:        txn.Type,
					Amount:      txn.Amount,
					Date:        MustParseTime(txn.Date),
					Description: txn.Description,
					CreatedAt:   MustParseTime(txn.CreatedAt),
					UpdatedAt:   MustParseTime(txn.UpdatedAt),
				}
			}
			req.Negotiate(http.StatusOK, err, GetTransactionsByUserIdResponse{
				Transactions: transactions,
			})
		},
	}
}

func (s TransactionHandler) GetOwnTransactionById() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:     "/own/{id}",
		Method:    http.MethodGet,
		FreeRoute: false,
		Dto:       &GetOwnTransactionByIdRequest{},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetOwnTransactionByIdResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			caller := req.MustGetCaller()
			sub := s.parser.MustParseSubject(caller.GetSubject())
			dto := req.MustGetDTO().(*GetOwnTransactionByIdRequest)
			resp, err := s.client.GetOwnTransactionById(context.Background(), &transactionv1.GetOwnTransactionByIdRequest{
				Id:     dto.Id,
				UserId: sub.UserId,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.Negotiate(http.StatusOK, err, GetOwnTransactionByIdResponse{
				Transaction: Transaction{
					Id:          resp.Transaction.Id,
					UserId:      resp.Transaction.UserId,
					Type:        resp.Transaction.Type,
					Amount:      resp.Transaction.Amount,
					Date:        MustParseTime(resp.Transaction.Date),
					Description: resp.Transaction.Description,
					CreatedAt:   MustParseTime(resp.Transaction.CreatedAt),
					UpdatedAt:   MustParseTime(resp.Transaction.UpdatedAt),
				},
			})
		},
	}
}

func (s TransactionHandler) GetAllTransactions() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "",
		Method:         http.MethodGet,
		FreeRoute:      false,
		AnyPermissions: []string{"ManageTransactions"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetAllTransactionsResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			resp, err := s.client.GetAllTransactions(context.Background(), &transactionv1.GetAllTransactionsRequest{})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			transactions := make([]Transaction, len(resp.Transactions))
			for i, txn := range resp.Transactions {
				transactions[i] = Transaction{
					Id:          txn.Id,
					UserId:      txn.UserId,
					Type:        txn.Type,
					Amount:      txn.Amount,
					Date:        MustParseTime(txn.Date),
					Description: txn.Description,
					CreatedAt:   MustParseTime(txn.CreatedAt),
					UpdatedAt:   MustParseTime(txn.UpdatedAt),
				}
			}
			req.Negotiate(http.StatusOK, err, GetAllTransactionsResponse{
				Transactions: transactions,
			})
		},
	}
}

func (s TransactionHandler) UpdateTransaction() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "/{id}",
		Method:         http.MethodPut,
		FreeRoute:      false,
		Dto:            &UpdateTransactionRequest{},
		AnyPermissions: []string{"ManageTransactions"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusNoContent,
				Description: "If everything is fine",
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*UpdateTransactionRequest)
			_, err := s.client.UpdateTransaction(context.Background(), &transactionv1.UpdateTransactionRequest{
				Id:          dto.Id,
				UserId:      dto.UserId,
				Type:        dto.Type,
				Amount:      dto.Amount,
				Description: dto.Description,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.ReturnStatus(http.StatusNoContent, err)
		},
	}
}

func (s TransactionHandler) DeleteTransaction() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "/{id}",
		Method:         http.MethodDelete,
		FreeRoute:      false,
		Dto:            &DeleteTransactionRequest{},
		AnyPermissions: []string{"ManageTransactions"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusNoContent,
				Description: "If everything is fine",
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*DeleteTransactionRequest)
			_, err := s.client.DeleteTransaction(context.Background(), &transactionv1.DeleteTransactionRequest{
				Id: dto.Id,
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			req.ReturnStatus(http.StatusNoContent, err)
		},
	}
}

func (s TransactionHandler) GetTransactionsWithPagination() *httpapi.RequestDefinition {
	return &httpapi.RequestDefinition{
		Route:          "/paginate",
		Method:         http.MethodGet,
		FreeRoute:      false,
		Dto:            &GetTransactionsWithPaginationRequest{},
		AnyPermissions: []string{"ManageTransactions"},
		ResponseDefinitions: []httpapi.ResponseDefinition{
			{
				Status:      http.StatusOK,
				Description: "If everything is fine",
				Dto:         &GetTransactionsWithPaginationResponse{},
			},
		},
		Handler: func(req httpapi.Request) {
			dto := req.MustGetDTO().(*GetTransactionsWithPaginationRequest)
			resp, err := s.client.GetTransactionsWithPagination(context.Background(), &transactionv1.GetTransactionsWithPaginationRequest{
				Offset: int32(dto.Offset),
				Limit:  int32(dto.Limit),
			})
			if err != nil {
				req.SetBadRequest(PleaseReadTheErrorCode, err.Error())
				return
			}
			transactions := make([]Transaction, len(resp.Transactions))
			for i, txn := range resp.Transactions {
				transactions[i] = Transaction{
					Id:          txn.Id,
					UserId:      txn.UserId,
					Type:        txn.Type,
					Amount:      txn.Amount,
					Date:        MustParseTime(txn.Date),
					Description: txn.Description,
					CreatedAt:   MustParseTime(txn.CreatedAt),
					UpdatedAt:   MustParseTime(txn.UpdatedAt),
				}
			}
			req.Negotiate(http.StatusOK, err, GetTransactionsWithPaginationResponse{
				Transactions: transactions,
			})
		},
	}
}

func MustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		log.Println("Error parsing date: %v", err)
	}
	return t
}

// dto

type Transaction struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	Type        string    `json:"type"`
	Amount      int64     `json:"amount"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateTransactionRequest struct {
	Type        string `json:"type" validate:"required,oneof=deposit withdrawal" example:"withdrawal"`
	Amount      int64  `json:"amount" validate:"required,gt=0" `
	Description string `json:"description" example:"samole"`
}

func (dto CreateTransactionRequest) Validate(ctx context.Context) error {
	return nil
}

type CreateTransactionResponse struct {
	Id string `json:"id"`
}

type GetTransactionByIdRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto GetTransactionByIdRequest) Validate(ctx context.Context) error {
	return nil
}

type GetTransactionByIdResponse struct {
	Transaction Transaction `json:"transaction"`
}

type GetOwnTransactionByIdRequest struct {
	Id     string `json:"id" validate:"required,uuid"`
	UserId string `json:"userId" validate:"required,uuid"`
}

func (dto GetOwnTransactionByIdRequest) Validate(ctx context.Context) error {
	return nil
}

type GetOwnTransactionByIdResponse struct {
	Transaction Transaction `json:"transaction"`
}

type GetAllTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type UpdateTransactionRequest struct {
	Id          string `json:"id" validate:"required,uuid"`
	UserId      string `json:"userId" validate:"required,uuid"`
	Type        string `json:"type" validate:"required,oneof=deposit withdrawal"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Description string `json:"description"`
}

func (dto UpdateTransactionRequest) Validate(ctx context.Context) error {
	return nil
}

type DeleteTransactionRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto DeleteTransactionRequest) Validate(ctx context.Context) error {
	return nil
}

type GetTransactionsByUserIdRequest struct {
	UserId string `json:"userId" validate:"required,uuid"`
}

func (dto GetTransactionsByUserIdRequest) Validate(ctx context.Context) error {
	return nil
}

type GetTransactionsByUserIdResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type GetTransactionsWithPaginationRequest struct {
	Offset int `json:"offset" validate:"gte=0"`
	Limit  int `json:"limit" validate:"gt=0"`
}

func (dto GetTransactionsWithPaginationRequest) Validate(ctx context.Context) error {
	return nil
}

type GetTransactionsWithPaginationResponse struct {
	Transactions []Transaction `json:"transactions"`
}
