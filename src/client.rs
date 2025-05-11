use crate::DISCORD_BASE_URL;
use crate::errors::Error;
use crate::events::base::{Credentials, WebhookMessage};
use reqwest::Client;

pub async fn execute_webhook(
    event: WebhookMessage,
    creds: Credentials,
    client: &Client,
) -> Result<(), Error> {
    let url = format!("{DISCORD_BASE_URL}/webhooks/{}/{}", creds.id, creds.token);

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
