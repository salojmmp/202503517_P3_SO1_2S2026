# Manual Técnico — Proyecto 3
## Sistemas Operativos 1 — Segundo Semestre 2026

**Universidad de San Carlos de Guatemala**  
**Proyecto:** Proyecto 3 — Sistemas distribuidos, almacenamiento y monitoreo  
**Repositorio:** `202503517_P3_SO1_2S2026`  
**Versión del documento:** 1.0  
**Fecha:** Octubre de 2026

---

# Índice

1. Introducción
2. Objetivos
3. Alcance del proyecto
4. Arquitectura general
5. Estructura del repositorio
6. Descripción de los componentes
7. Tecnologías utilizadas
8. Configuración y despliegue
9. Operador Helm y charts personalizados
10. Almacenamiento persistente con Valkey
11. Monitoreo con Prometheus y Grafana
12. API, comunicación gRPC y procesamiento de mensajes
13. Pruebas de funcionamiento
14. Pruebas de carga y escalabilidad
15. Seguridad y consideraciones operativas
16. Solución de problemas
17. Comandos de administración y diagnóstico
18. Conclusiones
19. Referencias técnicas
20. Anexos

---

# 1. Introducción

El presente manual técnico documenta la arquitectura, los componentes, la configuración y las pruebas del Proyecto 3 de Sistemas Operativos 1, correspondiente al segundo semestre de 2026.

El proyecto integra servicios desarrollados en Go y Rust con mecanismos de comunicación entre componentes, procesamiento de mensajes, almacenamiento de información y visualización mediante herramientas de monitoreo.

La estructura del proyecto incluye servicios de API, un servidor gRPC, un consumidor desarrollado en Go, archivos de definición de servicios, configuración de Prometheus, pruebas de carga con Locust y archivos de despliegue mediante Docker y Kubernetes.

También se implementó un operador basado en Helm para administrar instancias de Grafana y Valkey mediante recursos personalizados de Kubernetes.

El objetivo de este documento es facilitar la instalación, administración, verificación y diagnóstico de los componentes, así como proporcionar una referencia técnica para comprender la implementación.

# 2. Objetivos

## 2.1 Objetivo general

Documentar técnicamente la arquitectura y los procedimientos de configuración, despliegue, administración y verificación de los componentes del Proyecto 3.

## 2.2 Objetivos específicos

- Describir la función de cada servicio y su relación con los demás componentes.
- Documentar la estructura del código fuente y los archivos de configuración.
- Explicar el uso de Docker y Kubernetes en la administración de los servicios.
- Documentar la implementación del operador Helm y sus charts personalizados.
- Explicar la configuración de almacenamiento persistente mediante Valkey.
- Describir la integración de Prometheus y Grafana para consultar y visualizar métricas.
- Registrar las pruebas funcionales realizadas.
- Proporcionar comandos de diagnóstico para identificar problemas de ejecución.
- Documentar los aspectos relacionados con las pruebas de carga y escalabilidad.

# 3. Alcance del proyecto

El proyecto contempla los siguientes componentes:

| Componente | Función |
|---|---|
| API en Go | Servicio de API desarrollado en Go. |
| API en Rust | Servicio desarrollado en Rust. |
| Servidor gRPC | Proporciona comunicación mediante gRPC. |
| Consumidor en Go | Componente destinado al procesamiento de mensajes. |
| RabbitMQ | Broker de mensajería contemplado en la arquitectura. |
| Valkey | Almacenamiento de datos en memoria con persistencia. |
| Prometheus | Consulta y recopilación de métricas. |
| Grafana | Visualización de métricas mediante dashboards. |
| Metrics Exporter | Servicio de exportación de métricas. |
| Locust | Herramienta para pruebas de carga. |
| Docker | Construcción y ejecución de contenedores. |
| Kubernetes | Orquestación de los componentes. |
| Operador Helm | Administración declarativa de instancias de Grafana y Valkey. |
| Zot | Componente contemplado en los requisitos de despliegue de imágenes. |

La estructura del repositorio confirma la presencia de los archivos fuente de los servicios, la configuración de Prometheus, los archivos de Locust y los componentes de despliegue. La configuración exacta de cada servicio debe consultarse en sus respectivos archivos.

## 3.1 Entorno de validación

Durante las pruebas documentadas del operador se utilizó un clúster local de Minikube.

Se verificaron las instancias personalizadas de Grafana y Valkey, el estado de sus pods, la reconciliación del operador y la conexión entre Grafana y Prometheus.

**Importante:** la ejecución local no demuestra por sí sola el cumplimiento del requisito de despliegue en Google Kubernetes Engine (GKE). Para declarar dicho requisito como cumplido se necesita verificar el despliegue correspondiente en GKE.

# 4. Arquitectura general

La arquitectura está organizada en componentes independientes que pueden ejecutarse en contenedores y administrarse mediante Kubernetes.

## 4.1 Componentes principales

### Capa de servicios

Incluye la API desarrollada en Go, la API desarrollada en Rust, el servidor gRPC y el consumidor desarrollado en Go.

Cada componente mantiene su propia implementación y se encuentra organizado en archivos separados dentro del repositorio.

### Capa de comunicación y procesamiento

RabbitMQ forma parte de la arquitectura prevista para la distribución de mensajes. El consumidor en Go está destinado a procesar los mensajes recibidos.

El archivo `proto/report.proto` contiene la definición del protocolo utilizado para la comunicación basada en Protocol Buffers y gRPC.

La relación exacta entre los servicios, las colas, los métodos gRPC y los mensajes debe interpretarse a partir del código fuente y de la definición del protocolo.

### Capa de almacenamiento

Valkey proporciona almacenamiento de datos en memoria y admite persistencia en disco.

En el despliegue administrado por el operador, los datos se almacenan en un volumen persistente de Kubernetes montado en el directorio `/data` del contenedor.

### Capa de monitoreo

Prometheus permite consultar métricas y Grafana proporciona una interfaz para representarlas en dashboards.

El proyecto incluye un servicio de exportación de métricas y un archivo de configuración de Prometheus.

### Capa de orquestación

Kubernetes administra los pods, los deployments, los servicios de red y los volúmenes persistentes.

El operador Helm permite declarar instancias de Grafana y Valkey mediante recursos personalizados.

## 4.2 Diagrama lógico

```text
                    Usuarios y clientes
                           |
                +----------+----------+
                |                     |
             API Go                API Rust
                |                     |
                +----------+----------+
                           |
                  Comunicación interna
                           |
                     Servidor gRPC
                           |
                   Procesamiento Go
                           |
                    RabbitMQ
                           |
                     Consumidores
                           |
                         Valkey
                           |
                 Almacenamiento de datos


              Sistema de monitoreo
              --------------------
                 Metrics Exporter
                        |
                    Prometheus
                        |
                     Grafana
                        |
                     Dashboard


              Administración del clúster
              --------------------------
                     Kubernetes
                        |
                    Operador Helm
                     /        \
                 Grafana     Valkey
```

El diagrama representa las categorías de componentes del proyecto. No debe interpretarse como una descripción exacta de todas las conexiones implementadas: las relaciones definitivas deben confirmarse con el código fuente y la configuración de despliegue.

# 5. Estructura del repositorio

La estructura compartida del proyecto es la siguiente:

```text
202503517_P3_SO1_2S2026/
├── api-go/
│   └── main.go
├── api-rust/
│   ├── Cargo.lock
│   ├── Cargo.toml
│   ├── Dockerfile
│   └── src/
│       └── main.rs
├── consumer-go/
│   └── main.go
├── grpc-server/
│   └── main.go
├── locust/
│   ├── locustfile.py
│   └── requirements.txt
├── metrics-exporter/
│   ├── Dockerfile
│   └── main.go
├── operator/
│   ├── Dockerfile
│   ├── Makefile
│   ├── PROJECT
│   ├── watches.yaml
│   ├── .gitignore
│   ├── bin/
│   │   └── kustomize
│   └── helm-charts/
│       ├── grafanainstance/
│       └── valkeyinstance/
├── prometheus/
│   └── prometheus.yml
├── proto/
│   ├── report.proto
│   ├── report.pb.go
│   └── report_grpc.pb.go
├── docker-compose.yml
├── Dockerfile.go
├── .dockerignore
├── .gitkeep
├── go.mod
└── go.sum
```

La lista muestra los archivos principales compartidos. Los directorios pueden contener otros archivos de configuración y plantillas que no aparecen en esta representación.

## 5.1 Descripción de los archivos principales

| Archivo o directorio | Descripción |
|---|---|
| `api-go/main.go` | Código fuente de la API en Go. |
| `api-rust/src/main.rs` | Punto de entrada de la aplicación Rust. |
| `api-rust/Cargo.toml` | Dependencias y configuración del proyecto Rust. |
| `api-rust/Cargo.lock` | Versiones concretas de las dependencias resueltas. |
| `consumer-go/main.go` | Implementación del consumidor en Go. |
| `grpc-server/main.go` | Implementación del servidor gRPC. |
| `proto/report.proto` | Definición del protocolo y los mensajes gRPC. |
| `proto/report.pb.go` | Código Go generado a partir de Protocol Buffers. |
| `proto/report_grpc.pb.go` | Código Go generado para los servicios gRPC. |
| `metrics-exporter/main.go` | Implementación del exportador de métricas. |
| `metrics-exporter/Dockerfile` | Instrucciones para construir la imagen del exportador. |
| `prometheus/prometheus.yml` | Configuración de Prometheus. |
| `locust/locustfile.py` | Definición de las pruebas de carga. |
| `locust/requirements.txt` | Dependencias de Python para Locust. |
| `Dockerfile.go` | Instrucciones de construcción de una imagen relacionada con los servicios Go. |
| `docker-compose.yml` | Configuración de servicios para ejecución mediante Docker Compose. |
| `operator/Dockerfile` | Construcción de la imagen del operador. |
| `operator/Makefile` | Comandos de construcción y administración del operador. |
| `operator/watches.yaml` | Relación entre recursos personalizados y charts Helm. |
| `operator/helm-charts/` | Charts utilizados por el operador. |

# 6. Descripción de los componentes

## 6.1 API desarrollada en Go

Archivo principal:

```text
api-go/main.go
```

Este componente contiene la implementación de la API en Go.

Su comportamiento específico, las rutas HTTP, los parámetros de entrada y las respuestas deben consultarse en el código fuente.

Para verificar el servicio cuando se encuentre desplegado en Kubernetes:

```bash
kubectl get pods -A
kubectl get services -A
```

Para consultar los logs de su deployment, primero identifica su nombre y namespace:

```bash
kubectl get deployments -A
```

Después ejecuta:

```bash
kubectl logs deployment/api-go -n mumnk8s
```

El último comando corresponde al deployment `api-go` en el namespace `mumnk8s`, utilizado durante las pruebas del entorno local. Debe ajustarse si el nombre o el namespace cambian.

## 6.2 API desarrollada en Rust

Archivo principal:

```text
api-rust/src/main.rs
```

El componente contiene la implementación de la API en Rust.

El archivo `Cargo.toml` declara las dependencias del proyecto y `Cargo.lock` registra las versiones concretas resueltas.

El directorio `target/` contiene archivos generados durante la compilación.

Para compilar y comprobar el proyecto Rust, desde su directorio se puede utilizar:

```bash
cargo check
```

Este comando comprueba la compilación del código sin necesidad de iniciar el servicio.

Para consultar las rutas, los parámetros y los mecanismos de comunicación implementados, se debe revisar `src/main.rs`.

## 6.3 Servidor gRPC

Archivo principal:

```text
grpc-server/main.go
```

El servidor gRPC proporciona comunicación mediante servicios definidos utilizando Protocol Buffers.

La definición del protocolo se encuentra en:

```text
proto/report.proto
```

Los archivos generados para Go son:

```text
proto/report.pb.go
proto/report_grpc.pb.go
```

Estos archivos permiten utilizar los mensajes y las interfaces generadas a partir del protocolo.

Durante las pruebas del entorno Kubernetes se identificó un servicio llamado `grpc-server` en el namespace `mumnk8s`, con un endpoint interno que utilizaba el puerto `50051`.

Para consultar el servicio y sus endpoints:

```bash
kubectl get service grpc-server -n mumnk8s
kubectl get endpoints grpc-server -n mumnk8s
```

La existencia del endpoint demuestra que Kubernetes registró un destino para el servicio, pero no sustituye una prueba completa de todas las llamadas gRPC.

## 6.4 Consumidor desarrollado en Go

Archivo principal:

```text
consumer-go/main.go
```

Este componente contiene la implementación del consumidor encargado de procesar mensajes.

Su comportamiento depende de la configuración del sistema de mensajería, de los mensajes recibidos y de la lógica de procesamiento definida en el código.

Para verificar su estado en Kubernetes:

```bash
kubectl get deployments -A
kubectl get pods -A
```

Para consultar sus logs, utiliza el nombre real del deployment:

```bash
kubectl logs deployment/<nombre-del-consumidor> -n <namespace>
```

El resultado debe revisarse para identificar errores de conexión, problemas al procesar mensajes o fallos de almacenamiento.

## 6.5 RabbitMQ

RabbitMQ es el broker de mensajería contemplado en la arquitectura del proyecto.

Su función es facilitar la comunicación asíncrona entre los componentes que publican y consumen mensajes.

La configuración de las conexiones, las colas, los intercambios y las credenciales debe consultarse en los archivos de configuración y en el código de los servicios.

Para verificar los recursos relacionados en Kubernetes:

```bash
kubectl get pods -A
kubectl get services -A
```

Si RabbitMQ se encuentra desplegado en un namespace específico:

```bash
kubectl get pods -n <namespace>
kubectl logs deployment/<deployment-rabbitmq> -n <namespace>
```

La presencia de un pod en estado `Running` no demuestra por sí sola que las publicaciones y los consumos estén funcionando correctamente. Para validar la integración se deben comprobar los mensajes y los logs de los productores y consumidores.

## 6.6 Valkey

Valkey se utiliza como sistema de almacenamiento de datos en memoria.

El chart personalizado de Valkey permite crear una instancia administrada mediante Kubernetes y conservar datos utilizando un volumen persistente.

Durante las pruebas se utilizó una instancia personalizada llamada:

```text
valkeyinstance-sample
```

Esta instancia se desplegó en el namespace:

```text
operator-demo
```

La configuración verificada utilizaba la imagen:

```text
valkey/valkey:8
```

El contenedor ejecutaba Valkey con persistencia AOF y almacenaba sus archivos en `/data`.

Para comprobar el estado de la instancia:

```bash
kubectl get valkeyinstance -A
kubectl get pods -n operator-demo
kubectl get pvc -n operator-demo
```

La configuración de persistencia se explica en la sección 10.

## 6.7 Prometheus

Archivo principal:

```text
prometheus/prometheus.yml
```

Prometheus permite recopilar y consultar métricas expuestas por los servicios configurados.

La configuración de los objetivos de recopilación se encuentra en el archivo anterior.

Para revisar los recursos de Prometheus:

```bash
kubectl get pods -A
kubectl get services -A
```

La conectividad entre Prometheus y sus objetivos debe verificarse desde su interfaz y mediante consultas de métricas.

## 6.8 Metrics Exporter

Archivos principales:

```text
metrics-exporter/main.go
metrics-exporter/Dockerfile
```

El exportador implementa la lógica para exponer métricas que pueden ser consultadas por Prometheus.

Durante las pruebas se consultaron métricas relacionadas con los objetivos disponibles y con información agrupada por país.

Entre los nombres observados se encontraron:

```text
up
mumnk8s_country_reports
mumnk8s_country_warplanes_total
mumnk8s_country_warships_total
mumnk8s_assigned_country_info
```

Los valores y las etiquetas dependen de la información que el exportador proporcione. No deben interpretarse como datos oficiales o necesariamente actualizados sobre países o capacidades militares.

## 6.9 Grafana

Grafana proporciona una interfaz de visualización para construir dashboards a partir de fuentes de datos.

Durante las pruebas se creó una instancia personalizada:

```text
grafanainstance-sample
```

La instancia se desplegó en el namespace:

```text
default
```

El servicio utilizó el puerto `3000`.

Se configuró una fuente de datos Prometheus con la dirección interna:

```text
http://prometheus.mumnk8s.svc.cluster.local:9090
```

También se creó un dashboard denominado:

```text
Monitoreo MUMNK8S
```

El dashboard contiene paneles para visualizar:

- Estado de disponibilidad de los servicios.
- Total de objetivos disponibles.
- Reportes por país.
- Aviones por país.
- Barcos por país.
- País asignado.

La fuente de datos permite consultar las métricas disponibles en Prometheus. La exactitud y la interpretación de cada panel dependen de los datos y las etiquetas expuestos por el exportador.

## 6.10 Locust

Archivos principales:

```text
locust/locustfile.py
locust/requirements.txt
```

Locust es una herramienta para realizar pruebas de carga sobre servicios web.

El archivo `locustfile.py` define el comportamiento de los usuarios simulados y las solicitudes que se ejecutan durante la prueba.

Las dependencias necesarias se encuentran en `requirements.txt`.

La ejecución de Locust debe realizarse conforme a las instrucciones y los parámetros definidos para el proyecto. Los resultados deben documentarse con las métricas reales generadas por la herramienta.

La existencia del archivo de configuración no demuestra por sí sola que se hayan completado las pruebas exigidas por la rúbrica.

# 7. Tecnologías utilizadas

| Tecnología | Uso dentro del proyecto |
|---|---|
| Go | Desarrollo de servicios y componentes de procesamiento. |
| Rust | Desarrollo de la API correspondiente. |
| gRPC | Comunicación entre servicios mediante un protocolo definido. |
| Protocol Buffers | Definición y serialización de mensajes. |
| RabbitMQ | Mensajería entre componentes. |
| Valkey | Almacenamiento de datos en memoria y persistencia. |
| Prometheus | Consulta y recopilación de métricas. |
| Grafana | Visualización y análisis de métricas. |
| Locust | Pruebas de carga. |
| Docker | Construcción y ejecución de contenedores. |
| Docker Compose | Definición de servicios para ejecución coordinada. |
| Kubernetes | Orquestación de aplicaciones y recursos. |
| Helm | Empaquetado y administración de aplicaciones Kubernetes. |
| Operador Helm | Reconciliación de recursos personalizados con charts. |
| Zot | Registro de imágenes contemplado en los requisitos del proyecto. |

# 8. Configuración y despliegue

## 8.1 Requisitos previos

El entorno de trabajo debe contar con las herramientas correspondientes a las tareas que se vayan a realizar:

- Go para compilar y ejecutar los servicios desarrollados en Go.
- Rust y Cargo para compilar el servicio Rust.
- Docker para construir imágenes.
- Docker Compose para ejecutar los servicios definidos en el archivo correspondiente.
- `kubectl` para interactuar con Kubernetes.
- Helm para validar y administrar charts.
- Un clúster Kubernetes accesible.
- Locust y sus dependencias para ejecutar las pruebas de carga.

La disponibilidad de cada herramienta se puede comprobar con:

```bash
go version
cargo --version
docker --version
docker compose version
kubectl version --client
helm version
```

Para consultar el estado del clúster activo:

```bash
kubectl cluster-info
kubectl get nodes
```

Estos comandos permiten comprobar el entorno configurado. No crean ni modifican recursos.

## 8.2 Preparación del código fuente

El repositorio contiene archivos de Go y Rust, definiciones de gRPC, configuración de Prometheus y componentes de despliegue.

Antes de iniciar el sistema, se recomienda revisar:

- Las dependencias de cada servicio.
- Las variables de entorno requeridas.
- Las direcciones de los servicios internos.
- Los puertos utilizados.
- Las credenciales necesarias.
- La configuración de RabbitMQ y Valkey.
- La configuración de las imágenes de contenedor.
- Los archivos de Kubernetes y Helm.

Las variables de entorno y las credenciales deben obtenerse de la configuración real del proyecto. No se deben sustituir por valores inventados.

## 8.3 Construcción de imágenes

Los Dockerfiles contienen las instrucciones de construcción de las imágenes.

En particular, el repositorio incluye:

```text
Dockerfile.go
api-rust/Dockerfile
metrics-exporter/Dockerfile
operator/Dockerfile
```

Para construir una imagen a partir de un Dockerfile, se utiliza el siguiente patrón:

```bash
docker build -t <nombre-imagen>:<etiqueta> -f <ruta-dockerfile> .
```

El directorio de contexto debe contener los archivos que el Dockerfile necesita copiar.

Antes de construir una imagen, debe comprobarse que el contexto y las rutas coincidan con la estructura real del proyecto.

## 8.4 Docker Compose

El archivo:

```text
docker-compose.yml
```

define servicios para su ejecución coordinada mediante Docker Compose.

Para inspeccionar la configuración:

```bash
docker compose config
```

Para iniciar los servicios definidos:

```bash
docker compose up -d
```

Para consultar su estado:

```bash
docker compose ps
```

Para consultar los logs:

```bash
docker compose logs
```

**Advertencia:** estos comandos pueden iniciar, crear o modificar contenedores. Deben ejecutarse únicamente cuando se quiera utilizar el entorno definido en Docker Compose y después de revisar su configuración.

No se debe asumir que Docker Compose despliega automáticamente los recursos del clúster Kubernetes.

## 8.5 Despliegue en Kubernetes

Kubernetes administra los recursos utilizados por los servicios.

Entre los recursos que pueden intervenir se encuentran:

- Pods.
- Deployments.
- Services.
- ConfigMaps.
- Secrets.
- PersistentVolumeClaims.
- Recursos personalizados.
- Controladores del operador.

Para consultar los recursos existentes:

```bash
kubectl get pods -A
kubectl get deployments -A
kubectl get services -A
kubectl get pvc -A
```

Para inspeccionar un recurso concreto:

```bash
kubectl describe pod <nombre-pod> -n <namespace>
```

La aplicación de manifiestos debe realizarse utilizando las rutas y los archivos reales del proyecto. No se debe aplicar una configuración completa sin revisar previamente qué recursos modificará.

## 8.6 Despliegue en GKE

La rúbrica exige que el sistema esté desplegado y funcione en Google Kubernetes Engine.

La configuración local de Minikube no sustituye ese requisito.

Para verificar el contexto de Kubernetes seleccionado:

```bash
kubectl config current-context
kubectl cluster-info
kubectl get nodes
```

Estos comandos permiten identificar el contexto activo, consultar el clúster y comprobar los nodos.

Para documentar un despliegue en GKE se deben registrar el clúster utilizado, la región o zona, los namespaces, las imágenes desplegadas y el estado de los componentes.

No se incluyen aquí comandos de creación de infraestructura en GCP porque deben adaptarse al clúster real y a su configuración.

# 9. Operador Helm y charts personalizados

## 9.1 Descripción

El operador Helm permite administrar aplicaciones a partir de recursos personalizados de Kubernetes.

En este proyecto se utiliza para gestionar instancias de Grafana y Valkey.

El directorio del operador contiene:

```text
operator/
├── Dockerfile
├── Makefile
├── PROJECT
├── watches.yaml
└── helm-charts/
    ├── grafanainstance/
    └── valkeyinstance/
```

El archivo `watches.yaml` relaciona cada tipo de recurso personalizado con el chart que debe administrar.

La configuración utilizada fue:

```yaml
- group: infra.mumnk8s.local
  version: v1alpha1
  kind: ValkeyInstance
  chart: helm-charts/valkeyinstance

- group: infra.mumnk8s.local
  version: v1alpha1
  kind: GrafanaInstance
  chart: helm-charts/grafanainstance
```

Los recursos personalizados pertenecen al grupo:

```text
infra.mumnk8s.local
```

y utilizan la versión:

```text
v1alpha1
```

## 9.2 Construcción de la imagen del operador

El Dockerfile del operador utiliza la imagen base:

```dockerfile
FROM quay.io/operator-framework/helm-operator:v1.42.3

ENV HOME=/opt/helm

COPY watches.yaml ${HOME}/watches.yaml
COPY helm-charts ${HOME}/helm-charts

WORKDIR ${HOME}
```

La imagen incorpora el archivo de observación y los charts utilizados por el operador.

Para construir la imagen desde el directorio del operador:

```bash
docker build -t mumnk8s-operator:dev .
```

En el entorno local de Minikube, la imagen construida se cargó mediante:

```bash
minikube image load mumnk8s-operator:dev
```

Estos comandos describen el procedimiento utilizado para la imagen local. Para GKE se debe utilizar una imagen accesible desde el clúster y una etiqueta acorde con el despliegue correspondiente.

## 9.3 Validación de los charts

Los charts se encuentran en:

```text
helm-charts/grafanainstance
helm-charts/valkeyinstance
```

Para validar el chart de Grafana:

```bash
helm lint helm-charts/grafanainstance
```

Para validar el chart de Valkey:

```bash
helm lint helm-charts/valkeyinstance
```

Durante las pruebas realizadas, ambos comandos finalizaron con el resultado:

```text
1 chart(s) linted, 0 chart(s) failed
```

Se mostró una recomendación sobre el ícono de `Chart.yaml`, pero no se reportaron errores de validación.

La validación de Helm comprueba aspectos del chart, pero no garantiza por sí sola que la aplicación pueda desplegarse correctamente en cualquier clúster.

## 9.4 Verificación del operador

Para comprobar los recursos personalizados:

```bash
kubectl get grafanainstance -A
kubectl get valkeyinstance -A
```

Durante las pruebas se identificaron:

```text
default          grafanainstance-sample
operator-demo    valkeyinstance-sample
```

Para consultar los logs del operador:

```bash
kubectl logs deployment/mumnk8s-operator-controller-manager \
  -n mumnk8s-operator-system --since=5m
```

En los logs se observaron mensajes `Reconciled release` para ambas instancias.

Estos mensajes confirmaron que el operador había reconciliado las versiones administradas por Helm.

## 9.5 Permisos del operador

El operador necesita permisos de Kubernetes para administrar los recursos creados por los charts.

Durante la implementación se identificó que el operador no tenía permisos suficientes sobre ConfigMaps.

La regla correspondiente se agregó directamente al ClusterRole instalado:

```yaml
- apiGroups:
    - ""
  resources:
    - configmaps
  verbs:
    - get
    - list
    - watch
    - create
    - update
    - patch
    - delete
```

Después del ajuste, los permisos para consultar y crear ConfigMaps fueron verificados y el operador logró reconciliar Grafana correctamente.

**Nota de mantenimiento:** el ajuste se aplicó directamente al recurso del clúster. Debe comprobarse que la misma regla esté incorporada en el manifiesto fuente de RBAC para conservarla al reconstruir o reinstalar el operador. Esta sincronización no se debe dar por realizada automáticamente.

# 10. Almacenamiento persistente con Valkey

## 10.1 Configuración

La instancia de Valkey administrada por el operador utiliza un Deployment y un volumen persistente.

Durante las pruebas se observó la imagen:

```text
valkey/valkey:8
```

El servidor se inició con argumentos equivalentes a:

```text
valkey-server
--appendonly
yes
--dir
/data
```

La opción `--appendonly yes` habilita la persistencia basada en el archivo AOF.

El directorio `/data` se utiliza para almacenar los archivos de persistencia.

## 10.2 PersistentVolumeClaim

La instancia de demostración creó el siguiente recurso:

```text
valkeyinstance-sample
```

en el namespace:

```text
operator-demo
```

El PVC observado durante las pruebas tenía las siguientes características:

| Propiedad | Valor observado |
|---|---|
| Nombre | `valkeyinstance-sample` |
| Estado | `Bound` |
| Capacidad | `1Gi` |
| Modo de acceso | `RWO` |
| StorageClass | `standard` |
| Punto de montaje | `/data` |

El estado `Bound` indica que la solicitud de almacenamiento quedó asociada a un volumen.

## 10.3 Verificación del almacenamiento

Para consultar el PVC:

```bash
kubectl get pvc -n operator-demo
```

Para consultar el pod:

```bash
kubectl get pods -n operator-demo
```

Para comprobar el estado del Deployment:

```bash
kubectl get deployment valkeyinstance-sample -n operator-demo
```

## 10.4 Prueba funcional

Durante una prueba anterior se comprobó que Valkey respondía correctamente a `PING`:

```text
PONG
```

También se almacenó y recuperó un valor de prueba:

```text
Clave: persistencia
Valor: Dato persistente
```

La consulta del valor devolvió:

```text
Dato persistente
```

Estas pruebas confirmaron el funcionamiento de las operaciones básicas de lectura y escritura en el momento de la prueba.

La prueba de persistencia debe distinguirse de una prueba de recuperación después de reiniciar el pod. Para documentar esta última se necesita registrar que el dato continuó disponible después del reinicio, sin eliminar el PVC.

## 10.5 Protección de los datos

El PVC contiene los datos persistentes de la instancia.

Por ello:

- No debe eliminarse el PVC como procedimiento rutinario de diagnóstico.
- Antes de modificar el chart se deben revisar los cambios que puedan afectar al volumen.
- La eliminación del pod no debe confundirse con la eliminación del almacenamiento.
- Cualquier operación destructiva sobre el volumen debe contar con un procedimiento de respaldo y recuperación.

# 11. Monitoreo con Prometheus y Grafana

## 11.1 Configuración de Prometheus

La configuración de Prometheus se encuentra en:

```text
prometheus/prometheus.yml
```

Este archivo define los parámetros de recopilación de métricas utilizados por la instalación correspondiente.

Los objetivos configurados deben mantenerse accesibles desde Prometheus y deben exponer las métricas en el formato esperado.

Para verificar los recursos relacionados:

```bash
kubectl get pods -A
kubectl get services -A
```

## 11.2 Integración de Grafana con Prometheus

El chart de Grafana incluye un ConfigMap para provisionar la fuente de datos.

El archivo de plantilla es:

```text
helm-charts/grafanainstance/templates/datasources-configmap.yaml
```

El ConfigMap contiene una configuración equivalente a:

```yaml
apiVersion: 1
datasources:
  - name: Prometheus
    uid: prometheus
    type: prometheus
    access: proxy
    url: http://prometheus.mumnk8s.svc.cluster.local:9090
    isDefault: true
    editable: true
```

El Deployment de Grafana monta el archivo de configuración en:

```text
/etc/grafana/provisioning/datasources/prometheus.yaml
```

El montaje permite que Grafana cargue la fuente de datos mediante su mecanismo de aprovisionamiento.

La URL utiliza el DNS interno de Kubernetes. Por tanto, su funcionamiento depende de que Prometheus esté disponible en el namespace y servicio indicados.

## 11.3 Dashboard de monitoreo

Durante las pruebas se creó el dashboard:

```text
Monitoreo MUMNK8S
```

Los paneles configurados fueron:

| Panel | Consulta | Visualización |
|---|---|---|
| Estado de disponibilidad de los servicios | `up` | Stat |
| Total de objetivos disponibles | `sum(up)` | Stat |
| Reportes por país | `mumnk8s_country_reports` | Bar chart |
| Aviones por país | `mumnk8s_country_warplanes_total` | Bar chart |
| Barcos por país | `mumnk8s_country_warships_total` | Bar chart |
| País asignado | `mumnk8s_assigned_country_info` | Stat |

Los paneles agrupados por país utilizan la etiqueta `country_name` para mostrar la identificación de cada serie.

## 11.4 Prueba de conexión

Durante la validación se ejecutó la consulta:

```promql
up
```

Se obtuvo una serie con el siguiente resultado:

```text
up{instance="metrics-exporter:9100", job="mumnk8s-metrics"} 1
```

El valor `1` indica que Prometheus consideraba disponible el objetivo correspondiente en el momento de la consulta.

Esta prueba confirmó que Grafana podía consultar datos de Prometheus y que Prometheus tenía un objetivo disponible para el exportador.

## 11.5 Acceso a Grafana en el entorno local

Durante las pruebas, Grafana estuvo disponible mediante el servicio:

```text
grafanainstance-sample
```

en el namespace:

```text
default
```

El servicio utilizó el puerto `3000`.

Para abrir una conexión temporal desde el equipo local:

```bash
kubectl port-forward -n default \
  svc/grafanainstance-sample 13000:3000
```

Después se puede acceder desde el navegador a:

```text
http://localhost:13000
```

Las credenciales deben obtenerse de la configuración vigente de la instalación. No deben asumirse credenciales predeterminadas ni reutilizarse contraseñas sin comprobarlas.

# 12. API, comunicación gRPC y procesamiento de mensajes

## 12.1 Definición del protocolo

El archivo:

```text
proto/report.proto
```

contiene la definición del protocolo de comunicación basado en Protocol Buffers.

Los archivos generados para Go son:

```text
proto/report.pb.go
proto/report_grpc.pb.go
```

Los mensajes y métodos disponibles deben consultarse en `report.proto`.

## 12.2 Comunicación entre servicios

La arquitectura contempla servicios separados para la API en Go, la API en Rust, el servidor gRPC y el consumidor en Go.

La configuración de red debe permitir que cada componente se comunique con los servicios de los que depende.

Para diagnosticar problemas de comunicación deben comprobarse:

1. El estado de los pods.
2. La existencia de los Services.
3. Los endpoints asociados.
4. Los puertos configurados.
5. Las variables de entorno de conexión.
6. Los logs de ambos extremos de la comunicación.

Ejemplo para revisar el servicio gRPC:

```bash
kubectl get service grpc-server -n mumnk8s
kubectl get endpoints grpc-server -n mumnk8s
```

Si un Service no tiene endpoints, debe revisarse si existen pods disponibles que coincidan con sus selectores.

## 12.3 Procesamiento de mensajes

RabbitMQ forma parte del sistema de mensajería y el consumidor en Go contiene la lógica de procesamiento.

Para demostrar la integración completa es necesario comprobar que:

- El broker se encuentra disponible.
- El productor puede publicar mensajes.
- Las colas o los intercambios están configurados correctamente.
- El consumidor puede recibir los mensajes.
- El procesamiento finaliza según lo esperado.
- Los datos resultantes se almacenan en el destino correspondiente.

Estas comprobaciones deben respaldarse con logs, resultados de consultas o evidencias del funcionamiento. La existencia de los archivos fuente no sustituye esas pruebas.

# 13. Pruebas de funcionamiento

Esta sección registra las verificaciones documentadas durante el trabajo en el entorno local.

## 13.1 Validación de los charts Helm

Se ejecutaron:

```bash
helm lint helm-charts/grafanainstance
helm lint helm-charts/valkeyinstance
```

Resultados observados:

```text
1 chart(s) linted, 0 chart(s) failed
```

en ambos casos.

**Conclusión:** los dos charts superaron la validación realizada con `helm lint`.

## 13.2 Estado de las instancias personalizadas

Se ejecutaron:

```bash
kubectl get grafanainstance -A
kubectl get valkeyinstance -A
```

Se encontraron las instancias:

```text
default          grafanainstance-sample
operator-demo    valkeyinstance-sample
```

**Conclusión:** ambos recursos personalizados existían en el clúster durante la verificación.

## 13.3 Estado de los pods

Se ejecutaron:

```bash
kubectl get pods -n default
kubectl get pods -n operator-demo
```

En la verificación documentada, el pod de Grafana y el pod de Valkey aparecieron con el estado:

```text
1/1 Running
```

También apareció un pod de prueba de Prometheus con estado `Completed`. Ese estado corresponde a una ejecución finalizada y no indica por sí solo un fallo.

**Conclusión:** los pods principales estaban disponibles durante la comprobación.

## 13.4 Reconciliación del operador

Se ejecutó:

```bash
kubectl logs deployment/mumnk8s-operator-controller-manager \
  -n mumnk8s-operator-system --since=5m
```

Los logs mostraron mensajes como:

```text
Reconciled release
namespace="default"
name="grafanainstance-sample"
kind="GrafanaInstance"
```

y:

```text
Reconciled release
namespace="operator-demo"
name="valkeyinstance-sample"
kind="ValkeyInstance"
```

**Conclusión:** el operador reconcilió ambas instancias durante el periodo consultado y no se observaron errores en esa salida.

## 13.5 Prueba de disponibilidad en Prometheus

Se ejecutó la consulta:

```promql
up
```

El resultado incluyó:

```text
up{instance="metrics-exporter:9100", job="mumnk8s-metrics"} 1
```

**Conclusión:** el objetivo indicado estaba disponible según la métrica de Prometheus en el momento de la consulta.

## 13.6 Prueba de almacenamiento en Valkey

Se comprobó la respuesta:

```text
PING
PONG
```

También se verificó la recuperación de un valor de prueba:

```text
Dato persistente
```

**Conclusión:** las operaciones básicas de comunicación y recuperación de datos funcionaron durante la prueba documentada.

## 13.7 Resumen de resultados

| Prueba | Resultado documentado |
|---|---|
| Validación Helm de Grafana | Correcta |
| Validación Helm de Valkey | Correcta |
| Existencia de `GrafanaInstance` | Confirmada |
| Existencia de `ValkeyInstance` | Confirmada |
| Estado de los pods principales | `Running` durante la verificación |
| Reconciliación del operador | Mensajes exitosos |
| Consulta de métricas en Prometheus | Objetivo disponible |
| Consulta de datos en Valkey | Respuesta correcta |

Estos resultados corresponden a las pruebas realizadas. No deben extrapolarse automáticamente a un despliegue diferente.

# 14. Pruebas de carga y escalabilidad

## 14.1 Herramienta de pruebas

El proyecto incluye Locust mediante los archivos:

```text
locust/locustfile.py
locust/requirements.txt
```

El archivo `locustfile.py` define el comportamiento de los usuarios simulados y las solicitudes que se ejecutan durante la prueba.

## 14.2 Objetivo de las pruebas

Las pruebas de carga permiten estudiar el comportamiento del sistema ante diferentes cantidades de solicitudes concurrentes.

La rúbrica solicita pruebas comparativas con una y dos réplicas para los componentes indicados, además de análisis del rendimiento.

La documentación de cada ejecución debería registrar:

- Fecha y entorno de ejecución.
- Endpoint probado.
- Número de usuarios simulados.
- Tasa de generación de usuarios.
- Duración de la prueba.
- Total de solicitudes.
- Solicitudes fallidas.
- Tiempo de respuesta promedio.
- Percentiles de tiempo de respuesta.
- Rendimiento en solicitudes por segundo.
- Número de réplicas de cada componente.
- Consumo de CPU y memoria, cuando esté disponible.

## 14.3 Comparación de réplicas

La rúbrica requiere comparar el comportamiento con una y dos réplicas para los componentes especificados.

La comparación debe realizarse bajo condiciones equivalentes, manteniendo constantes la duración, el número de usuarios, la tasa de generación y el tipo de solicitudes.

La tabla siguiente sirve para registrar resultados reales:

| Métrica | 1 réplica | 2 réplicas |
|---|---:|---:|
| Solicitudes totales | Pendiente de registrar | Pendiente de registrar |
| Solicitudes fallidas | Pendiente de registrar | Pendiente de registrar |
| Tiempo de respuesta promedio | Pendiente de registrar | Pendiente de registrar |
| Percentil 95 | Pendiente de registrar | Pendiente de registrar |
| Solicitudes por segundo | Pendiente de registrar | Pendiente de registrar |
| CPU utilizada | Pendiente de registrar | Pendiente de registrar |
| Memoria utilizada | Pendiente de registrar | Pendiente de registrar |

No se deben completar estas celdas con estimaciones inventadas.

## 14.4 HPA

Horizontal Pod Autoscaler (HPA) permite ajustar el número de réplicas de un recurso compatible en función de métricas configuradas.

Para consultar los HPA existentes:

```bash
kubectl get hpa -A
```

Para inspeccionar uno:

```bash
kubectl describe hpa <nombre-hpa> -n <namespace>
```

La presencia de un HPA debe verificarse en el clúster. No se debe afirmar que la escalabilidad automática funciona sin comprobar su configuración y observar su comportamiento.

## 14.5 Análisis

El análisis debe comparar los resultados obtenidos con una y dos réplicas.

Un mayor número de réplicas no garantiza por sí solo una mejora. El resultado depende de los recursos disponibles, los límites de CPU y memoria, el tipo de carga y la capacidad de los componentes para procesar solicitudes en paralelo.

La conclusión final debe redactarse después de completar las pruebas y revisar sus resultados reales.

# 15. Seguridad y consideraciones operativas

## 15.1 Credenciales

Las contraseñas, los tokens y las cadenas de conexión no deben incluirse directamente en el manual técnico.

Cuando una aplicación necesita credenciales, estas deben proporcionarse mediante el mecanismo de configuración correspondiente, como variables de entorno o Secrets de Kubernetes.

El manual puede describir la existencia de dichas configuraciones sin revelar sus valores.

## 15.2 Permisos de Kubernetes

El operador requiere permisos para consultar y administrar los recursos que gestiona.

Los permisos deben limitarse a los necesarios para realizar su función.

Las modificaciones directas de RBAC deben mantenerse sincronizadas con los manifiestos fuente para que puedan reproducirse en futuras instalaciones.

## 15.3 Persistencia

El PVC de Valkey debe protegerse frente a eliminaciones accidentales.

Antes de realizar operaciones sobre el almacenamiento se debe comprobar qué recursos serán modificados y qué impacto pueden tener sobre los datos.

## 15.4 Recursos de otros componentes

El clúster puede contener recursos que no pertenezcan directamente al operador personalizado.

Antes de aplicar cambios globales se deben revisar el namespace, los selectores y los nombres de los recursos afectados.

## 15.5 Entorno local y entorno de producción

Minikube resulta útil para desarrollo y pruebas locales, mientras que GKE proporciona un entorno Kubernetes administrado en Google Cloud.

Las configuraciones de imágenes, almacenamiento, red, permisos y credenciales deben adaptarse al entorno de destino.

El funcionamiento en Minikube no constituye evidencia suficiente del funcionamiento en GKE.

# 16. Solución de problemas

## 16.1 El pod no inicia

Consultar:

```bash
kubectl get pods -A
kubectl describe pod <nombre-pod> -n <namespace>
```

Revisar los eventos para identificar problemas de imágenes, configuración, almacenamiento, permisos o recursos insuficientes.

## 16.2 El contenedor se reinicia

Consultar:

```bash
kubectl logs <nombre-pod> -n <namespace>
kubectl logs <nombre-pod> -n <namespace> --previous
```

La opción `--previous` permite consultar los logs del contenedor anterior cuando corresponde.

## 16.3 El operador no reconcilia una instancia

Consultar:

```bash
kubectl get grafanainstance -A
kubectl get valkeyinstance -A
```

Después revisar los logs:

```bash
kubectl logs deployment/mumnk8s-operator-controller-manager \
  -n mumnk8s-operator-system --since=10m
```

Verificar la existencia del recurso personalizado, los permisos RBAC y la correspondencia entre `watches.yaml` y el chart asociado.

## 16.4 Grafana no muestra datos

Comprobar:

1. Que el pod de Grafana esté disponible.
2. Que Prometheus esté disponible.
3. Que la URL configurada para la fuente de datos sea correcta.
4. Que el objetivo de Prometheus esté disponible.
5. Que la consulta utilizada por el panel devuelva datos.

Una consulta de referencia es:

```promql
up
```

Si la consulta no devuelve datos, deben revisarse primero Prometheus y sus objetivos antes de modificar el dashboard.

## 16.5 Valkey no responde

Consultar:

```bash
kubectl get pods -n operator-demo
kubectl get pvc -n operator-demo
kubectl logs deployment/valkeyinstance-sample -n operator-demo
```

Comprobar que el PVC esté asociado a un volumen y que el contenedor tenga acceso al directorio de datos.

No eliminar el PVC como método de diagnóstico.

## 16.6 El Service no tiene endpoints

Consultar:

```bash
kubectl get services -n <namespace>
kubectl get endpoints -n <namespace>
kubectl get pods -n <namespace> --show-labels
```

Comparar los selectores del Service con las etiquetas de los pods.

## 16.7 Un chart falla durante la validación

Ejecutar:

```bash
helm lint <ruta-del-chart>
```

Revisar el resultado y corregir los errores señalados antes de desplegar cambios.

La validación local ayuda a detectar errores de estructura y plantillas, pero no reemplaza la prueba en Kubernetes.

# 17. Comandos de administración y diagnóstico

Los siguientes comandos permiten consultar el estado del entorno.

## 17.1 Estado general

```bash
kubectl cluster-info
kubectl get nodes
kubectl get pods -A
kubectl get deployments -A
kubectl get services -A
kubectl get pvc -A
```

## 17.2 Recursos personalizados

```bash
kubectl get grafanainstance -A
kubectl get valkeyinstance -A
```

## 17.3 Logs del operador

```bash
kubectl logs deployment/mumnk8s-operator-controller-manager \
  -n mumnk8s-operator-system --since=5m
```

## 17.4 Validación de charts

```bash
helm lint helm-charts/grafanainstance
helm lint helm-charts/valkeyinstance
```

## 17.5 Acceso local a Grafana

```bash
kubectl port-forward -n default \
  svc/grafanainstance-sample 13000:3000
```

Acceso desde el navegador:

```text
http://localhost:13000
```

## 17.6 Revisión de Locust

Antes de ejecutar pruebas de carga, revisar el archivo:

```text
locust/locustfile.py
```

y consultar las dependencias:

```text
locust/requirements.txt
```

La dirección del servicio objetivo y los parámetros de carga deben corresponder al entorno que se desea evaluar.

**Nota:** los comandos de esta sección son principalmente de consulta o validación. Los comandos que despliegan servicios, reconstruyen imágenes o modifican recursos deben utilizarse solo después de revisar sus efectos.

# 18. Conclusiones

El proyecto integra componentes desarrollados en diferentes tecnologías con herramientas de mensajería, almacenamiento, monitoreo y orquestación.

La estructura del repositorio separa los servicios en Go y Rust, el servidor gRPC, el consumidor, el exportador de métricas, la configuración de Prometheus y las herramientas de pruebas de carga.

Durante las pruebas del entorno local se validaron los charts Helm de Grafana y Valkey, se verificó la existencia de sus recursos personalizados y se observaron reconciliaciones exitosas del operador.

También se comprobó que Grafana podía consultar métricas de Prometheus y que Valkey respondía a operaciones básicas de lectura y escritura.

Estas verificaciones constituyen evidencia del funcionamiento de los componentes probados. Para completar la evaluación integral deben considerarse por separado los requisitos de despliegue en GKE, funcionamiento de todos los servicios, mensajería, pruebas de carga, escalabilidad y documentación de resultados.

El manual proporciona una referencia para comprender la organización del proyecto, verificar los componentes y diagnosticar problemas sin depender exclusivamente de la memoria de los desarrolladores.

# 19. Referencias técnicas

- Documentación oficial de Kubernetes: https://kubernetes.io/docs/
- Documentación oficial de Helm: https://helm.sh/docs/
- Documentación de Prometheus: https://prometheus.io/docs/
- Documentación de Grafana: https://grafana.com/docs/
- Documentación de Valkey: https://valkey.io/
- Documentación de RabbitMQ: https://www.rabbitmq.com/docs
- Documentación de Go: https://go.dev/doc/
- Documentación de Rust: https://doc.rust-lang.org/
- Documentación de gRPC: https://grpc.io/docs/
- Documentación de Protocol Buffers: https://protobuf.dev/
- Documentación de Locust: https://docs.locust.io/
- Documentación de Docker: https://docs.docker.com/
- Documentación de Google Kubernetes Engine: https://cloud.google.com/kubernetes-engine/docs

# 20. Anexos

## Anexo A. Evidencias recomendadas

Para respaldar el manual durante la evaluación, se recomienda conservar evidencias de:

1. Estructura del repositorio.
2. Validación de ambos charts Helm.
3. Recursos personalizados de Grafana y Valkey.
4. Pods y servicios en estado disponible.
5. Logs del operador con reconciliaciones exitosas.
6. Estado del PVC de Valkey.
7. Respuesta `PONG` de Valkey.
8. Consulta de datos almacenados en Valkey.
9. Configuración de la fuente de datos de Prometheus en Grafana.
10. Dashboard con sus paneles.
11. Resultado de la consulta `up` en Prometheus.
12. Despliegue en GKE, si corresponde al entorno final de entrega.
13. Ejecuciones de Locust y sus resultados.
14. Configuración y comportamiento del HPA.
15. Comparación de rendimiento con una y dos réplicas.

Las evidencias deben corresponder a pruebas reales y mostrar claramente el entorno y el componente evaluado.

## Anexo B. Información que debe completarse si existe evidencia adicional

Si se dispone de los resultados correspondientes, conviene añadir:

- Dirección o identificación del clúster GKE utilizado.
- Procedimiento exacto para desplegar todos los servicios.
- Configuración final de Zot.
- Rutas HTTP y métodos gRPC implementados.
- Configuración real de las colas y los consumidores de RabbitMQ.
- Resultados de las pruebas de carga.
- Configuración del HPA.
- Comparación del rendimiento con diferentes números de réplicas.
- Procedimiento de respaldo y recuperación de Valkey.

Estos datos deben extraerse de la configuración y de las pruebas reales del proyecto, no de valores supuestos.
