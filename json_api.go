package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/Noverload/price-fetcher/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        PriceService
}

func NewJSONAPIServer(listenAddr string, svc PriceService) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) RUN() {
	http.HandleFunc("/", makeHTTPHandleFunc(s.handleFetchPrice))
	http.ListenAndServe(s.listenAddr, nil)
}

func makeHTTPHandleFunc(apiFn APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestID", rand.Intn(10000000))
	return func(w http.ResponseWriter, r *http.Request) {
		if err := apiFn(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")

	price, err := s.svc.FetchPrice(ctx, ticker)

	if err != nil {
		return err
	}
	priceResp := types.PriceResponse{
		Price:  price,
		Ticker: ticker,
	}
	return writeJSON(w, http.StatusOK, &priceResp)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
