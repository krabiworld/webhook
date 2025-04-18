use crate::errors::Error;
use crate::events::Event;
use crate::structs::{StarEvent, WebhookMessage};

impl Event for StarEvent {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error> {
        if self.action != "created" {
            return Ok(None);
        }

        Ok(Some(WebhookMessage {
            content: format!(
                "[{}](<{}>) starred [{}](<{}>) <:foxtada:1311327105300172882>",
                self.sender.login,
                self.sender.html_url,
                self.repository.name,
                self.repository.html_url
            ),
            username: self.sender.login.clone(),
            avatar_url: self.sender.avatar_url.clone(),
        }))
    }
}
