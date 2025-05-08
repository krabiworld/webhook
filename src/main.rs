mod client;
mod errors;
mod events;
mod parser;
mod server;

use crate::server::{health, webhook};
use actix_web::{App, HttpServer, web};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

const DISCORD_BASE_URL: &str = "https://discord.com/api";
const GITHUB_EVENT: &str = "X-GitHub-Event";
const GITHUB_SIG: &str = "X-Hub-Signature-256";

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let addr = std::env::var("ADDR").unwrap_or_else(|_| "0.0.0.0".into());
    let port = std::env::var("PORT")
        .ok()
        .and_then(|p| p.parse().ok())
        .unwrap_or(8080);
    let secret = std::env::var("SECRET").expect("env variable `SECRET` should be set");
    let secret = Arc::new(secret);

    env_logger::init();

    let client = reqwest::Client::new();
    let star_jail = Arc::new(Mutex::new(HashMap::<String, bool>::new()));

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(secret.clone()))
            .app_data(web::Data::new(client.clone()))
            .app_data(web::Data::new(star_jail.clone()))
            .service(health)
            .service(webhook)
    })
    .bind((addr, port))?
    .run()
    .await
}
