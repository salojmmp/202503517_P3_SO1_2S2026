use axum::{
    Json, Router,
    extract::State,
    http::StatusCode,
    routing::{get, post},
};
use serde::{Deserialize, Serialize};
use std::{env, net::SocketAddr};
use tokio::net::TcpListener;

#[derive(Clone)]
struct AppState {
    go_api_url: String,
    client: reqwest::Client,
}

#[derive(Deserialize, Serialize)]
struct ReportRequest {
    country_code: String,
    warplanes_in_air: i32,
    warships_in_water: i32,
    timestamp_unix: Option<i64>,
}

#[derive(Serialize, Deserialize)]
struct ApiResponse {
    accepted: bool,
    message: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    report_id: Option<String>,
}

async fn health() -> Json<serde_json::Value> {
    Json(serde_json::json!({
        "status": "ok",
        "service": "mumnk8s-api-rust"
    }))
}

async fn submit_report(
    State(state): State<AppState>,
    Json(mut report): Json<ReportRequest>,
) -> (StatusCode, Json<ApiResponse>) {
    report.country_code = report.country_code.trim().to_uppercase();

    let allowed_countries = ["USA", "RUS", "CHN", "ESP", "GTM"];

    if !allowed_countries.contains(&report.country_code.as_str()) {
        return (
            StatusCode::BAD_REQUEST,
            Json(ApiResponse {
                accepted: false,
                message: "country_code debe ser USA, RUS, CHN, ESP o GTM".into(),
                report_id: None,
            }),
        );
    }

    if !(0..=50).contains(&report.warplanes_in_air) {
        return (
            StatusCode::BAD_REQUEST,
            Json(ApiResponse {
                accepted: false,
                message: "warplanes_in_air debe estar entre 0 y 50".into(),
                report_id: None,
            }),
        );
    }

    if !(0..=30).contains(&report.warships_in_water) {
        return (
            StatusCode::BAD_REQUEST,
            Json(ApiResponse {
                accepted: false,
                message: "warships_in_water debe estar entre 0 y 30".into(),
                report_id: None,
            }),
        );
    }

    let go_url = format!("{}/reports", state.go_api_url.trim_end_matches('/'));

    let response = match state.client.post(&go_url).json(&report).send().await {
        Ok(response) => response,
        Err(error) => {
            eprintln!("Error conectando con API Go: {error}");

            return (
                StatusCode::BAD_GATEWAY,
                Json(ApiResponse {
                    accepted: false,
                    message: "No se pudo conectar con la API Go".into(),
                    report_id: None,
                }),
            );
        }
    };

    let status = response.status();

    if status.is_success() {
        match response.json::<ApiResponse>().await {
            Ok(body) => (StatusCode::ACCEPTED, Json(body)),
            Err(error) => {
                eprintln!("Respuesta inválida de API Go: {error}");

                (
                    StatusCode::BAD_GATEWAY,
                    Json(ApiResponse {
                        accepted: false,
                        message: "La API Go devolvió una respuesta inválida".into(),
                        report_id: None,
                    }),
                )
            }
        }
    } else {
        let message = response
            .text()
            .await
            .unwrap_or_else(|_| "Error desconocido en API Go".into());

        (
            if status.is_client_error() {
                StatusCode::BAD_REQUEST
            } else {
                StatusCode::BAD_GATEWAY
            },
            Json(ApiResponse {
                accepted: false,
                message,
                report_id: None,
            }),
        )
    }
}

#[tokio::main]
async fn main() {
    let go_api_url = env::var("GO_API_URL").unwrap_or_else(|_| "http://127.0.0.1:8080".into());

    let address: SocketAddr = env::var("RUST_API_ADDR")
        .unwrap_or_else(|_| "0.0.0.0:8000".into())
        .parse()
        .expect("RUST_API_ADDR no es una dirección válida");

    let state = AppState {
        go_api_url: go_api_url.clone(),
        client: reqwest::Client::new(),
    };

    let app = Router::new()
        .route("/health", get(health))
        .route("/reports", post(submit_report))
        .with_state(state);

    let listener = TcpListener::bind(address)
        .await
        .expect("No se pudo abrir el puerto de la API Rust");

    println!("API REST Rust escuchando en {address}");
    println!("API Go configurada en {go_api_url}");

    axum::serve(listener, app)
        .await
        .expect("Error ejecutando la API REST Rust");
}
