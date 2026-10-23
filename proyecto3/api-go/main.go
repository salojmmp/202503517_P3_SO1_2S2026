package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	reportpb "example.com/mumnk8s/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ReportRequest struct {
	CountryCode     string `json:"country_code"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	TimestampUnix   int64  `json:"timestamp_unix"`
}

type ReportResponse struct {
	Accepted bool   `json:"accepted"`
	Message  string `json:"message"`
	ReportID string `json:"report_id,omitempty"`
}

type API struct {
	client reportpb.ReportServiceClient
}

func writeJSON(w http.ResponseWriter, statusCode int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "mumnk8s-api-go",
	})
}

func (a *API) submitReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input ReportRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, ReportResponse{
			Accepted: false,
			Message:  "JSON inválido: " + err.Error(),
		})
		return
	}

	input.CountryCode = strings.ToUpper(strings.TrimSpace(input.CountryCode))

	if input.CountryCode == "" {
		writeJSON(w, http.StatusBadRequest, ReportResponse{
			Accepted: false,
			Message:  "country_code es obligatorio",
		})
		return
	}

	if input.WarplanesInAir < 0 || input.WarplanesInAir > 50 {
		writeJSON(w, http.StatusBadRequest, ReportResponse{
			Accepted: false,
			Message:  "warplanes_in_air debe estar entre 0 y 50",
		})
		return
	}

	if input.WarshipsInWater < 0 || input.WarshipsInWater > 30 {
		writeJSON(w, http.StatusBadRequest, ReportResponse{
			Accepted: false,
			Message:  "warships_in_water debe estar entre 0 y 30",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	response, err := a.client.SubmitReport(ctx, &reportpb.ReportRequest{
		CountryCode:     input.CountryCode,
		WarplanesInAir:  input.WarplanesInAir,
		WarshipsInWater: input.WarshipsInWater,
		TimestampUnix:   input.TimestampUnix,
	})
	if err != nil {
		code := http.StatusBadGateway
		message := "No se pudo enviar el reporte al servicio gRPC"

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			code = http.StatusGatewayTimeout
			message = "Tiempo de espera agotado al enviar el reporte"
		} else if grpcStatus, ok := status.FromError(err); ok {
			if grpcStatus.Code().String() == "InvalidArgument" {
				code = http.StatusBadRequest
				message = grpcStatus.Message()
			}
		}

		log.Printf("Error gRPC: %v", err)

		writeJSON(w, code, ReportResponse{
			Accepted: false,
			Message:  message,
		})
		return
	}

	writeJSON(w, http.StatusAccepted, ReportResponse{
		Accepted: response.GetAccepted(),
		Message:  response.GetMessage(),
		ReportID: response.GetReportId(),
	})
}

func main() {
	grpcAddress := os.Getenv("GRPC_SERVER_ADDR")
	if grpcAddress == "" {
		grpcAddress = "localhost:50051"
	}

	httpAddress := os.Getenv("HTTP_ADDR")
	if httpAddress == "" {
		httpAddress = ":8080"
	}

	connection, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("No se pudo crear el cliente gRPC: %v", err)
	}
	defer connection.Close()

	api := &API{
		client: reportpb.NewReportServiceClient(connection),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", api.health)
	mux.HandleFunc("/reports", api.submitReport)

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("API Go escuchando en %s; servidor gRPC: %s", httpAddress, grpcAddress)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error de la API Go: %v", err)
	}
}
