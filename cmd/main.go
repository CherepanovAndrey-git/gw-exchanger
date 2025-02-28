package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/CherepanovAndrey-git/gw-exchange-grpc/proto/exchange"
	_ "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	grpcReflection "google.golang.org/grpc/reflection"
	"gw-exchanger/internal/database"
)

type Server struct {
	exchange.UnimplementedExchangeServiceServer
	queries *database.Queries
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		queries: database.New(db),
	}
}

func (s *Server) GetExchangeRates(ctx context.Context, req *exchange.Empty) (*exchange.ExchangeRatesResponse, error) {
	log.Println("Handling GetExchangeRates request")

	rates, err := s.queries.GetAllRates(ctx)
	if err != nil {
		log.Printf("Database error: %v", err)
		return nil, fmt.Errorf("failed to get rates: %w", err)
	}

	log.Printf("Fetched %d rates from DB", len(rates))

	response := &exchange.ExchangeRatesResponse{
		Rates: make(map[string]float32),
	}
	for _, rate := range rates {
		key := fmt.Sprintf("%s_%s", rate.FromCurrency, rate.ToCurrency)
		log.Printf("Processing rate: %s = %s", key, rate.Rate)

		rateValue, err := strconv.ParseFloat(rate.Rate, 32)
		if err != nil {
			log.Printf("Invalid rate format: %s = %s (error: %v)", key, rate.Rate, err)
			continue
		}

		response.Rates[key] = float32(rateValue)
	}
	log.Printf("Returning %d rates", len(response.Rates))
	return response, nil
}

func (s *Server) GetExchangeRateForCurrency(ctx context.Context, req *exchange.CurrencyRequest) (*exchange.ExchangeRateResponse, error) {
	rate, err := s.queries.GetRate(ctx, database.GetRateParams{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get rate: %w", err)
	}

	rateValue, err := strconv.ParseFloat(rate, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rate: %w", err)
	}

	return &exchange.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         float32(rateValue),
	}, nil
}

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	port := os.Getenv("EXCHANGER_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	exchange.RegisterExchangeServiceServer(s, NewServer(db))
	grpcReflection.Register(s)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
