use log::info;
use once_cell::sync::OnceCell;
use serde::Deserialize;
use std::error::Error;

#[derive(Deserialize, Debug)]
pub struct Config {
    #[serde(default = "default_address")]
    pub address: String,
    #[serde(default = "default_port")]
    pub port: u16,
    pub secret: String,
    pub happy_emoji: String,
    pub success_emoji: String,
    pub failure_emoji: String,
}

static CONFIG: OnceCell<Config> = OnceCell::new();

pub fn init() -> Result<(), Box<dyn Error>> {
    let config = envy::from_env::<Config>()?;
    CONFIG.set(config).map_err(|_| "Config already initialized")?;

    info!("Config initialized");

    Ok(())
}

pub fn get() -> &'static Config {
    CONFIG.get().expect("Config not initialized")
}

fn default_address() -> String {
    "0.0.0.0".to_string()
}

fn default_port() -> u16 {
    8080
}
