use crate::errors::Error;
use crate::structs::WebhookMessage;

mod fork;
mod push;
mod star;
mod workflow_run;

pub trait Event {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error>;
}
