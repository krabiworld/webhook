mod client;
mod errors;
mod events;
mod parser;
mod server;

use crate::server::webhook;
use actix_web::{App, HttpServer};

const DISCORD_BASE_URL: &str = "https://discord.com/api";
const GITHUB_EVENT: &str = "X-GitHub-Event";

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let addr = std::env::var("ADDR").unwrap_or_else(|_| "0.0.0.0".into());
    let port = std::env::var("PORT")
        .ok()
        .and_then(|p| p.parse().ok())
        .unwrap_or(8080);

    env_logger::init();

    HttpServer::new(|| App::new().service(webhook))
        .bind((addr, port))?
        .run()
        .await
}
