package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Country struct {
	Code string
	Name string
}

var countries = []Country{
	{Code: "USA", Name: "Estados Unidos"},
	{Code: "RUS", Name: "Rusia"},
	{Code: "CHN", Name: "China"},
	{Code: "ESP", Name: "España"},
	{Code: "GTM", Name: "Guatemala"},
}

var assignedCountry = "ESP"

var valkey *redis.Client

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getHashInt(
	ctx context.Context,
	key string,
	field string,
) (int64, error) {
	value, err := valkey.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func getDistributionStats(
	ctx context.Context,
	key string,
) (min, max, mode int64, err error) {
	values, err := valkey.ZRangeWithScores(
		ctx,
		key,
		0,
		-1,
	).Result()
	if err != nil {
		return 0, 0, 0, err
	}

	if len(values) == 0 {
		return 0, 0, 0, nil
	}

	var highestFrequency float64 = -1

	for i, item := range values {
		// Member contiene el valor reportado.
		value, parseErr := strconv.ParseInt(
			fmt.Sprint(item.Member),
			10,
			64,
		)
		if parseErr != nil {
			return 0, 0, 0, fmt.Errorf(
				"valor inválido en distribución %s: %v",
				key,
				item.Member,
			)
		}

		// Score contiene la frecuencia del valor.
		frequency := item.Score

		if i == 0 || value < min {
			min = value
		}
		if i == 0 || value > max {
			max = value
		}

		// En caso de empate, elegir el valor menor.
		if frequency > highestFrequency ||
			(frequency == highestFrequency && value < mode) {
			mode = value
			highestFrequency = frequency
		}
	}

	return min, max, mode, nil
}

func writeMetric(
	builder *strings.Builder,
	name string,
	help string,
	labels string,
	value int64,
) {
	fmt.Fprintf(builder, "# HELP %s %s\n", name, help)
	fmt.Fprintf(builder, "# TYPE %s gauge\n", name)

	if labels == "" {
		fmt.Fprintf(builder, "%s %d\n", name, value)
	} else {
		fmt.Fprintf(builder, "%s{%s} %d\n", name, labels, value)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var output strings.Builder

	// Totales globales.
	totalReports, err := getHashInt(
		ctx,
		"reports:totals",
		"total_reports",
	)
	if err != nil {
		http.Error(w, "Error leyendo totales de Valkey", http.StatusInternalServerError)
		log.Printf("Error leyendo total_reports: %v", err)
		return
	}

	totalPlanes, err := getHashInt(
		ctx,
		"reports:totals",
		"warplanes_in_air",
	)
	if err != nil {
		http.Error(w, "Error leyendo totales de Valkey", http.StatusInternalServerError)
		return
	}

	totalShips, err := getHashInt(
		ctx,
		"reports:totals",
		"warships_in_water",
	)
	if err != nil {
		http.Error(w, "Error leyendo totales de Valkey", http.StatusInternalServerError)
		return
	}

	writeMetric(
		&output,
		"mumnk8s_total_reports",
		"Total de reportes procesados",
		"",
		totalReports,
	)
	writeMetric(
		&output,
		"mumnk8s_total_warplanes",
		"Suma de aviones reportados",
		"",
		totalPlanes,
	)
	writeMetric(
		&output,
		"mumnk8s_total_warships",
		"Suma de barcos reportados",
		"",
		totalShips,
	)

	// Estadísticas globales.
	for _, stat := range []struct {
		name string
		key  string
		help string
	}{
		{
			name: "warplanes",
			key:  "reports:distribution:warplanes",
			help: "aviones",
		},
		{
			name: "warships",
			key:  "reports:distribution:warships",
			help: "barcos",
		},
	} {
		min, max, mode, statsErr := getDistributionStats(ctx, stat.key)
		if statsErr != nil {
			http.Error(w, "Error leyendo distribuciones de Valkey", http.StatusInternalServerError)
			log.Printf("Error leyendo %s: %v", stat.key, statsErr)
			return
		}

		writeMetric(
			&output,
			"mumnk8s_"+stat.name+"_min",
			"Mínimo global de "+stat.help,
			"",
			min,
		)
		writeMetric(
			&output,
			"mumnk8s_"+stat.name+"_max",
			"Máximo global de "+stat.help,
			"",
			max,
		)
		writeMetric(
			&output,
			"mumnk8s_"+stat.name+"_mode",
			"Moda global de "+stat.help,
			"",
			mode,
		)
	}

	// Totales y estadísticas por país.
	for _, country := range countries {
		labels := fmt.Sprintf(
			`country_code="%s",country_name="%s"`,
			country.Code,
			country.Name,
		)

		baseKey := "reports:country:" + country.Code + ":totals"

		reportCount, err := getHashInt(ctx, baseKey, "total_reports")
		if err != nil {
			http.Error(w, "Error leyendo estadísticas por país", http.StatusInternalServerError)
			return
		}

		planes, err := getHashInt(ctx, baseKey, "warplanes_in_air")
		if err != nil {
			http.Error(w, "Error leyendo estadísticas por país", http.StatusInternalServerError)
			return
		}

		ships, err := getHashInt(ctx, baseKey, "warships_in_water")
		if err != nil {
			http.Error(w, "Error leyendo estadísticas por país", http.StatusInternalServerError)
			return
		}

		writeMetric(
			&output,
			"mumnk8s_country_reports",
			"Cantidad de reportes por país",
			labels,
			reportCount,
		)
		writeMetric(
			&output,
			"mumnk8s_country_warplanes_total",
			TotalAvionesPorPaisHelp(),
			labels,
			planes,
		)
		writeMetric(
			&output,
			"mumnk8s_country_warships_total",
			"Total de barcos reportados por país",
			labels,
			ships,
		)

		for _, stat := range []struct {
			name string
		}{
			{name: "warplanes"},
			{name: "warships"},
		} {
			key := "reports:country:" + country.Code +
				":distribution:" + stat.name

			min, max, mode, statsErr := getDistributionStats(ctx, key)
			if statsErr != nil {
				http.Error(w, "Error leyendo estadísticas por país", http.StatusInternalServerError)
				return
			}

			writeMetric(
				&output,
				"mumnk8s_country_"+stat.name+"_min",
				"Mínimo por país",
				labels,
				min,
			)
			writeMetric(
				&output,
				"mumnk8s_country_"+stat.name+"_max",
				"Máximo por país",
				labels,
				max,
			)
			writeMetric(
				&output,
				"mumnk8s_country_"+stat.name+"_mode",
				"Moda por país",
				labels,
				mode,
			)
		}
	}

	// Identificación del país asignado.
	for _, country := range countries {
		if country.Code == assignedCountry {
			writeMetric(
				&output,
				"mumnk8s_assigned_country_info",
				"País asignado al proyecto",
				fmt.Sprintf(
					`country_code="%s",country_name="%s"`,
					country.Code,
					country.Name,
				),
				1,
			)
			break
		}
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	_, _ = w.Write([]byte(output.String()))
}

func TotalAvionesPorPaisHelp() string {
	return "Total de aviones reportados por país"
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := valkey.Ping(ctx).Err(); err != nil {
		http.Error(w, "Valkey no disponible", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func main() {
	valkeyAddr := envOrDefault("VALKEY_ADDR", "localhost:6379")
	listenAddr := envOrDefault("METRICS_ADDR", ":9100")
	assignedCountry = strings.ToUpper(
		envOrDefault("ASSIGNED_COUNTRY", "ESP"),
	)

	valkey = redis.NewClient(&redis.Options{
		Addr: valkeyAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := valkey.Ping(ctx).Err()
	cancel()

	if err != nil {
		log.Fatalf("No se pudo conectar a Valkey en %s: %v", valkeyAddr, err)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Printf(
		"Exportador de métricas escuchando en %s; Valkey: %s; país asignado: %s",
		listenAddr,
		valkeyAddr,
		assignedCountry,
	)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
