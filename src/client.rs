use crate::DISCORD_BASE_URL;
use crate::errors::Error;
use crate::structs::{Credentials, Discord};

pub async fn execute_webhook(content: String, username: String, avatar_url: String, creds: Credentials) -> Result<(), Error> {
    let body = Discord {
        content,
        username,
        avatar_url,
    };

    let url = format!(
        "{}/webhooks/{}/{}",
        DISCORD_BASE_URL, creds.id, creds.token
    );

    let client = reqwest::Client::new();
    let res = client.post(&url)
        .header("Content-Type", "application/json")
        .json(&body)
        .send().await?;

    if !res.status().is_success() {
        let body = res.bytes().await?;
        let body_str = std::str::from_utf8(&body)?;
        return Err(Error::DiscordError(body_str.to_string()));
    }

    Ok(())
}
