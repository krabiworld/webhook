mod client;
mod errors;
mod events;
mod parser;
mod server;
mod config;

use crate::server::{health, webhook};
use actix_web::{App, HttpServer, web};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

const DISCORD_BASE_URL: &str = "https://discord.com/api";
const GITHUB_EVENT: &str = "X-GitHub-Event";
const GITHUB_SIG: &str = "X-Hub-Signature-256";

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenvy::dotenv().ok();
    
    env_logger::init();
    
    config::init().expect("Failed to initialize config");

    let client = reqwest::Client::new();
    let star_jail = Arc::new(Mutex::new(HashMap::<String, bool>::new()));

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(client.clone()))
            .app_data(web::Data::new(star_jail.clone()))
            .service(health)
            .service(webhook)
    })
    .bind((config::get().address.clone(), config::get().port))?
    .run()
    .await
}
