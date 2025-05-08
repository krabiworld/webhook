pub(crate) mod base;
pub(crate) mod check_run;
pub(crate) mod fork;
pub(crate) mod push;
pub(crate) mod release;
pub(crate) mod star;
pub(crate) mod workflow_run;

use crate::errors::Error;
use base::WebhookMessage;

pub trait Event {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error>;
}
