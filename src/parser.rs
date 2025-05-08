use crate::client::execute_webhook;
use crate::errors::Error;
use crate::events::Event;
use crate::events::base::Credentials;
use crate::events::check_run::CheckRunEvent;
use crate::events::fork::ForkEvent;
use crate::events::push::PushEvent;
use crate::events::release::ReleaseEvent;
use crate::events::star::StarEvent;
use crate::events::workflow_run::WorkflowRunEvent;
use actix_web::web;
use reqwest::Client;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::Duration;
use tokio::time::sleep;

pub async fn parse_event(
    event: String,
    body: web::Bytes,
    creds: Credentials,
    client: &Client,
    star_jail: &Arc<Mutex<HashMap<String, bool>>>,
) -> Result<(), Error> {
    let event_result = match event.as_str() {
        "push" => serde_json::from_slice::<PushEvent>(&body)?.handle()?,
        "workflow_run" => serde_json::from_slice::<WorkflowRunEvent>(&body)?.handle()?,
        "star" => {
            let star = serde_json::from_slice::<StarEvent>(&body)?;

            // moved from event
            if star.action != "created" {
                return Ok(());
            }

            let mut map = star_jail
                .lock()
                .map_err(|e| Error::MutexPoisonError(e.to_string()))?;
            let key = format!("{}-{}", star.sender.login, star.repository.name);
            match map.get(key.as_str()) {
                Some(_) => return Ok(()),
                None => {
                    let k = key.clone();

                    map.insert(key, true);

                    let m = star_jail.clone();
                    tokio::spawn(async move {
                        sleep(Duration::from_secs(86400)).await;
                        m.lock().unwrap().remove(&k);
                    });
                }
            }

            star.handle()?
        }
        "fork" => serde_json::from_slice::<ForkEvent>(&body)?.handle()?,
        "release" => serde_json::from_slice::<ReleaseEvent>(&body)?.handle()?,
        "check_run" => serde_json::from_slice::<CheckRunEvent>(&body)?.handle()?,
        _ => None,
    };

    if let Some(event_result) = event_result {
        execute_webhook(event_result, creds, client).await?;
    }

    Ok(())
}
