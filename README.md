# Proyecto 3 — Sistemas Operativos 1

**Universidad de San Carlos de Guatemala**  
**Segundo Semestre 2026**

## Descripción

Proyecto de sistemas distribuidos que integra servicios desarrollados en Go y Rust, comunicación mediante gRPC, procesamiento de mensajes con RabbitMQ, almacenamiento con Valkey y monitoreo mediante Prometheus y Grafana.

También incorpora contenedores Docker, orquestación con Kubernetes, un operador basado en Helm y pruebas de carga con Locust.

## Tecnologías

- **Go y Rust:** desarrollo de servicios.
- **gRPC y Protocol Buffers:** comunicación entre servicios.
- **RabbitMQ:** mensajería.
- **Valkey:** almacenamiento de datos.
- **Prometheus y Grafana:** monitoreo y visualización.
- **Docker y Kubernetes:** contenedores y orquestación.
- **Helm:** administración de instancias mediante un operador.
- **Locust:** pruebas de carga.

## Estructura del proyecto

```text
├── api-go/             # API en Go
├── api-rust/           # API en Rust
├── consumer-go/        # Consumidor en Go
├── grpc-server/        # Servidor gRPC
├── metrics-exporter/   # Exportador de métricas
├── locust/             # Pruebas de carga
├── operator/           # Operador y charts Helm
├── prometheus/         # Configuración de Prometheus
├── proto/              # Definiciones gRPC
├── docker-compose.yml
└── Dockerfile.go
```

## Documentación

El archivo `MANUAL_TECNICO.md` contiene la descripción de la arquitectura, los componentes, la configuración, las pruebas y los procedimientos de diagnóstico.

## Validación del operador

Desde el directorio `operator`, se pueden validar los charts con:

```bash
helm lint helm-charts/grafanainstance
helm lint helm-charts/valkeyinstance
```

## Autores

Proyecto académico desarrollado para el curso de **Sistemas Operativos 1**.
