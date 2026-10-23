package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Report struct {
	ReportID        string `json:"report_id"`
	CountryCode     string `json:"country_code"`
	WarplanesInAir  int    `json:"warplanes_in_air"`
	WarshipsInWater int    `json:"warships_in_water"`
	TimestampUnix   int64  `json:"timestamp_unix"`
}

var allowedCountries = map[string]bool{
	"USA": true,
	"RUS": true,
	"CHN": true,
	"ESP": true,
	"GTM": true,
}

var saveReportScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[1]) == 1 then
    return 0
end

local report_id = ARGV[1]
local report_json = ARGV[2]
local timestamp = tonumber(ARGV[3])
local country = ARGV[4]
local planes = tonumber(ARGV[5])
local ships = tonumber(ARGV[6])

-- Reporte completo e idempotencia
redis.call('SET', KEYS[2], report_json)
redis.call('SET', KEYS[1], '1')

-- Historiales global y por país
redis.call('ZADD', KEYS[3], timestamp, report_id)
redis.call('ZADD', KEYS[4], timestamp, report_id)

-- Totales globales
redis.call('HINCRBY', KEYS[5], 'warplanes_in_air', planes)
redis.call('HINCRBY', KEYS[5], 'warships_in_water', ships)
redis.call('HINCRBY', KEYS[5], 'total_reports', 1)

-- Totales por país
redis.call('HINCRBY', KEYS[6], 'warplanes_in_air', planes)
redis.call('HINCRBY', KEYS[6], 'warships_in_water', ships)
redis.call('HINCRBY', KEYS[6], 'total_reports', 1)

-- Distribuciones de valores para máximos, mínimos y modas
redis.call('ZINCRBY', KEYS[7], 1, tostring(planes))
redis.call('ZINCRBY', KEYS[8], 1, tostring(ships))
redis.call('ZINCRBY', KEYS[9], 1, tostring(planes))
redis.call('ZINCRBY', KEYS[10], 1, tostring(ships))

-- Cantidad de reportes por país
redis.call('HINCRBY', KEYS[11], country, 1)

return 1
`)

func validateReport(report Report) error {
	report.CountryCode = strings.ToUpper(strings.TrimSpace(report.CountryCode))

	if report.ReportID == "" {
		return fmt.Errorf("report_id es obligatorio")
	}

	if !allowedCountries[report.CountryCode] {
		return fmt.Errorf("country_code no permitido: %s", report.CountryCode)
	}

	if report.WarplanesInAir < 0 || report.WarplanesInAir > 50 {
		return fmt.Errorf("warplanes_in_air debe estar entre 0 y 50")
	}

	if report.WarshipsInWater < 0 || report.WarshipsInWater > 30 {
		return fmt.Errorf("warships_in_water debe estar entre 0 y 30")
	}

	if report.TimestampUnix <= 0 {
		return fmt.Errorf("timestamp_unix debe ser positivo")
	}

	return nil
}

func storeReport(
	ctx context.Context,
	client *redis.Client,
	report Report,
) (bool, error) {
	report.CountryCode = strings.ToUpper(strings.TrimSpace(report.CountryCode))

	data, err := json.Marshal(report)
	if err != nil {
		return false, fmt.Errorf("serializar reporte: %w", err)
	}

	country := report.CountryCode

	keys := []string{
		"processed:" + report.ReportID,
		"report:" + report.ReportID,
		"reports:timeline",
		"reports:country:" + country + ":timeline",
		"reports:totals",
		"reports:country:" + country + ":totals",
		"reports:distribution:warplanes",
		"reports:distribution:warships",
		"reports:country:" + country + ":distribution:warplanes",
		"reports:country:" + country + ":distribution:warships",
		"reports:country:counts",
	}

	args := []interface{}{
		report.ReportID,
		string(data),
		report.TimestampUnix,
		country,
		report.WarplanesInAir,
		report.WarshipsInWater,
	}

	result, err := saveReportScript.Run(ctx, client, keys, args...).Int64()
	if err != nil {
		return false, fmt.Errorf("guardar reporte en Valkey: %w", err)
	}

	return result == 1, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://p3user:p3localpass@localhost:5672/"
	}

	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr == "" {
		valkeyAddr = "localhost:6379"
	}

	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "reports"
	}

	valkey := redis.NewClient(&redis.Options{
		Addr: valkeyAddr,
	})
	defer valkey.Close()

	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	err := valkey.Ping(pingCtx).Err()
	pingCancel()

	if err != nil {
		log.Fatalf("No se pudo conectar a Valkey: %v", err)
	}

	log.Println("Conexión con Valkey establecida")

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

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("No se pudo declarar la cola: %v", err)
	}

	if err := channel.Qos(10, 0, false); err != nil {
		log.Fatalf("No se pudo configurar QoS: %v", err)
	}

	messages, err := channel.Consume(
		queue.Name,
		"mumnk8s-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("No se pudo consumir la cola: %v", err)
	}

	log.Printf("Consumidor activo. Cola: %s", queue.Name)

	for {
		select {
		case <-ctx.Done():
			log.Println("Cerrando consumidor...")
			return

		case message, open := <-messages:
			if !open {
				log.Println("El canal de mensajes se cerró")
				return
			}

			var report Report

			if err := json.Unmarshal(message.Body, &report); err != nil {
				log.Printf("Mensaje JSON inválido: %v", err)
				_ = message.Nack(false, false)
				continue
			}

			if err := validateReport(report); err != nil {
				log.Printf("Reporte rechazado: %v", err)
				_ = message.Nack(false, false)
				continue
			}

			storeCtx, storeCancel := context.WithTimeout(ctx, 5*time.Second)
			stored, err := storeReport(storeCtx, valkey, report)
			storeCancel()

			if err != nil {
				log.Printf("Error al guardar reporte: %v", err)
				_ = message.Nack(false, true)
				continue
			}

			if err := message.Ack(false); err != nil {
				log.Printf("Error al confirmar mensaje: %v", err)
				continue
			}

			if !stored {
				log.Printf("Reporte duplicado ignorado: id=%s", report.ReportID)
				continue
			}

			log.Printf(
				"Reporte procesado: id=%s país=%s aviones=%d barcos=%d",
				report.ReportID,
				report.CountryCode,
				report.WarplanesInAir,
				report.WarshipsInWater,
			)
		}
	}
}
