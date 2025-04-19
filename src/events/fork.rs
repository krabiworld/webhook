use crate::errors::Error;
use crate::events::Event;
use crate::events::base::{Forkee, User, WebhookMessage};
use serde::Deserialize;

#[derive(Deserialize)]
pub struct ForkEvent {
    pub sender: User,
    pub forkee: Forkee,
}

impl Event for ForkEvent {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error> {
        Ok(Some(WebhookMessage {
            content: format!(
                "[{}](<{}>) forked [{}](<{}>)",
                self.sender.login, self.sender.html_url, self.forkee.name, self.forkee.html_url
            ),
            username: self.sender.login.clone(),
            avatar_url: self.sender.avatar_url.clone(),
        }))
    }
}
