package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	reportpb "example.com/mumnk8s/proto"
)

type Report struct {
	ReportID        string `json:"report_id"`
	CountryCode     string `json:"country_code"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	TimestampUnix   int64  `json:"timestamp_unix"`
}

type reportServer struct {
	reportpb.UnimplementedReportServiceServer
	channel   *amqp.Channel
	queueName string
}

func generateID() (string, error) {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(b[:]), nil
}

func (s *reportServer) SubmitReport(
	ctx context.Context,
	req *reportpb.ReportRequest,
) (*reportpb.ReportResponse, error) {
	country := strings.ToUpper(strings.TrimSpace(req.GetCountryCode()))

	allowedCountries := map[string]bool{
		"USA": true,
		"RUS": true,
		"CHN": true,
		"ESP": true,
		"GTM": true,
	}

	if !allowedCountries[country] {
		return nil, status.Error(codes.InvalidArgument, "country_code debe ser USA, RUS, CHN, ESP o GTM")
	}

	if req.GetWarplanesInAir() < 0 || req.GetWarplanesInAir() > 50 {
		return nil, status.Error(codes.InvalidArgument, "warplanes_in_air debe estar entre 0 y 50")
	}

	if req.GetWarshipsInWater() < 0 || req.GetWarshipsInWater() > 30 {
		return nil, status.Error(codes.InvalidArgument, "warships_in_water debe estar entre 0 y 30")
	}

	timestamp := req.GetTimestampUnix()
	if timestamp <= 0 {
		timestamp = time.Now().Unix()
	}

	id, err := generateID()
	if err != nil {
		return nil, status.Error(codes.Internal, "no se pudo generar report_id")
	}

	report := Report{
		ReportID:        id,
		CountryCode:     country,
		WarplanesInAir:  req.GetWarplanesInAir(),
		WarshipsInWater: req.GetWarshipsInWater(),
		TimestampUnix:   timestamp,
	}

	body, err := json.Marshal(report)
	if err != nil {
		return nil, status.Error(codes.Internal, "no se pudo serializar el reporte")
	}

	if err := s.channel.PublishWithContext(
		ctx,
		"",
		s.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		log.Printf("Error publicando reporte: %v", err)
		return nil, status.Error(codes.Unavailable, "no se pudo publicar el reporte")
	}

	log.Printf("Reporte publicado: id=%s país=%s", id, country)

	return &reportpb.ReportResponse{
		Accepted: true,
		Message:  "Reporte aceptado y enviado a RabbitMQ",
		ReportId: id,
	}, nil
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://p3user:p3localpass@localhost:5672/"
	}

	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "reports"
	}

	addr := os.Getenv("GRPC_ADDR")
	if addr == "" {
		addr = ":50051"
	}

	connection, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("No se pudo conectar a RabbitMQ: %v", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("No se pudo abrir el canal RabbitMQ: %v", err)
	}
	defer channel.Close()

	if _, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalf("No se pudo declarar la cola: %v", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("No se pudo escuchar en %s: %v", addr, err)
	}
	defer listener.Close()

	server := grpc.NewServer()
	reportpb.RegisterReportServiceServer(server, &reportServer{
		channel:   channel,
		queueName: queueName,
	})

	fmt.Printf("Servidor gRPC escuchando en %s\n", addr)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error del servidor gRPC: %v", err)
	}
}
