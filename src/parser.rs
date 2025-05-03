use crate::client::execute_webhook;
use crate::errors::Error;
use crate::events::Event;
use crate::events::base::Credentials;
use crate::events::fork::ForkEvent;
use crate::events::push::PushEvent;
use crate::events::release::ReleaseEvent;
use crate::events::star::StarEvent;
use crate::events::workflow_run::WorkflowRunEvent;
use actix_web::web;
use reqwest::Client;

pub async fn parse_event(
    event: String,
    body: web::Bytes,
    creds: Credentials,
    client: &Client,
) -> Result<(), Error> {
    let event_result = match event.as_str() {
        "push" => serde_json::from_slice::<PushEvent>(&body)?.handle()?,
        "workflow_run" => serde_json::from_slice::<WorkflowRunEvent>(&body)?.handle()?,
        "star" => serde_json::from_slice::<StarEvent>(&body)?.handle()?,
        "fork" => serde_json::from_slice::<ForkEvent>(&body)?.handle()?,
        "release" => serde_json::from_slice::<ReleaseEvent>(&body)?.handle()?,
        _ => None,
    };

    if let Some(event_result) = event_result {
        execute_webhook(event_result, creds, client).await?;
    }

    Ok(())
}
