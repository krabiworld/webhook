use log::error;
use crate::DISCORD_BASE_URL;
use crate::structs::{Credentials, Discord};

pub async fn execute_webhook(content: String, username: String, avatar: String, creds: Credentials) {
    let body = Discord {
        content: content.to_string(),
        username: username.to_string(),
        avatar_url: avatar.to_string(),
    };

    let url = format!(
        "{}/webhooks/{}/{}",
        DISCORD_BASE_URL, creds.id, creds.token
    );

    let client = reqwest::Client::new();
    match client.post(&url)
        .header("Content-Type", "application/json")
        .json(&body)
        .send().await {
        Ok(res) => {
            if !res.status().is_success() {
                let body = res.bytes().await.unwrap();
                let body_str = std::str::from_utf8(&body).unwrap();
                error!("discord request failed: {}", body_str);
            }
        }
        Err(e) => {
            error!("http request failed: {}", e);
            return;
        }
    };
}
