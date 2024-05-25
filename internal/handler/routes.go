package handler

import (
	"credit-plus/internal/model/request"
	"credit-plus/internal/model/response"
	"credit-plus/middleware"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

func (s *Server) RegisterRoutes(basicAuth *middleware.BasicAuthMiddleware) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)

	mux.HandleFunc("/health", s.healthHandler)

	mux.Handle("POST /customers", basicAuth.ValidateBasicAuth(http.HandlerFunc(s.CustomerHandler)))
	mux.Handle("POST /contracts", basicAuth.ValidateBasicAuth(http.HandlerFunc(s.ContractHandler)))
	return mux
}

func (s *Server) CustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.Write(resp)
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Customer Handler Failed to Read Body", "Error", err)
		return
	}

	req := request.CreateCustomerRequest{}
	json.Unmarshal(body, &req)

	err = s.service.CreateCustomer(r.Context(), req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.Write(resp)
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Customer Handler Failed to Read Body", "Error", err)

		return
	}
	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Creating Customer"})
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
	return
}

func (s *Server) ContractHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.Write(resp)
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Customer Handler Failed to Read Body", "Error", err)
		return
	}

	req := request.CreateContactRequest{}
	json.Unmarshal(body, &req)

	code, err := s.service.CreateContract(r.Context(), req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(code)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Creating Contract"})
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
	return
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.service.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
