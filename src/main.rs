mod structs;
mod parser;
mod client;

use actix_web::{post, web, App, HttpRequest, HttpResponse, HttpServer, Responder};
use actix_web::http::StatusCode;
use crate::parser::parse_event;
use crate::structs::{Credentials};

const DISCORD_BASE_URL: &str = "https://discord.com/api";
const GITHUB_EVENT: &str = "X-GitHub-Event";

#[post("/{id}/{token}")]
async fn webhook(req: HttpRequest, creds: web::Path<Credentials>, body: web::Bytes) -> impl Responder {
    let event = req.headers().get(GITHUB_EVENT).and_then(|v| v.to_str().ok()).unwrap_or("").to_string();
    if creds.id.trim().is_empty() || creds.token.trim().is_empty() || event.trim().is_empty() {
        return HttpResponse::Ok().body("Header or credentials are empty");
    }

    tokio::spawn(async move {
        parse_event(event, body, creds.into_inner()).await;
    });

    HttpResponse::Ok().status(StatusCode::NO_CONTENT).finish()
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    unsafe {
        std::env::set_var("RUST_LOG", "info");
    }
    env_logger::init();

    HttpServer::new(|| {
        App::new().service(webhook)
    })
        .bind(("0.0.0.0", 8080))?
        .run()
        .await
}
