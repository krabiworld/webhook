use crate::DISCORD_BASE_URL;
use crate::errors::Error;
use crate::events::base::{Credentials, WebhookMessage};

pub async fn execute_webhook(event: WebhookMessage, creds: Credentials) -> Result<(), Error> {
    let url = format!("{}/webhooks/{}/{}", DISCORD_BASE_URL, creds.id, creds.token);

    let client = reqwest::Client::new();
    let res = client
        .post(&url)
        .header("Content-Type", "application/json")
        .json(&event)
        .send()
        .await?;

    if !res.status().is_success() {
        let body = res.bytes().await?;
        let body_str = std::str::from_utf8(&body)?;
        return Err(Error::DiscordError(body_str.to_string()));
    }

    Ok(())
}
