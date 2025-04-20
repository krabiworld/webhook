use crate::errors::Error;
use crate::events::Event;
use crate::events::base::{Release, User, WebhookMessage};
use serde::Deserialize;

#[derive(Deserialize)]
pub struct ReleaseEvent {
    pub action: String,
    pub release: Release,
    pub sender: User,
}

impl Event for ReleaseEvent {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error> {
        if self.action != "published" {
            return Ok(None);
        }

        Ok(Some(WebhookMessage {
            content: format!(
                "[{}](<{}>) published release [{}](<{}>)",
                self.sender.login,
                self.sender.html_url,
                self.release.tag_name,
                self.release.html_url
            ),
            username: self.sender.login.clone(),
            avatar_url: self.sender.avatar_url.clone(),
        }))
    }
}
