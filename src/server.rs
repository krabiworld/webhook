use crate::events::base::Credentials;
use crate::parser::parse_event;
use crate::{config, GITHUB_EVENT, GITHUB_SIG};
use actix_web::{HttpRequest, HttpResponse, Responder, get, post, web};
use hmac::{Hmac, Mac};
use log::error;
use reqwest::Client;
use sha2::Sha256;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use subtle::ConstantTimeEq;

type HmacSha256 = Hmac<Sha256>;

fn no_content() -> HttpResponse {
    HttpResponse::NoContent().finish()
}

#[get("/health")]
pub async fn health() -> impl Responder {
    HttpResponse::Ok().body("OK")
}

#[post("/{id}/{token}")]
pub async fn webhook(
    req: HttpRequest,
    creds: web::Path<Credentials>,
    body: web::Bytes,
    client: web::Data<Client>,
    star_jail: web::Data<Arc<Mutex<HashMap<String, bool>>>>,
) -> impl Responder {
    if config::get().signature_check {
        // get and check GitHub signature
        let sig = req
            .headers()
            .get(GITHUB_SIG)
            .and_then(|h| h.to_str().ok())
            .filter(|s| s.starts_with("sha256="))
            .map(|s| &s[7..]);

        let sig_hex = match sig {
            Some(s) => s,
            None => return no_content(),
        };

        let mut mac = match HmacSha256::new_from_slice(config::get().secret.as_bytes()) {
            Ok(o) => o,
            Err(e) => {
                error!("{}", e);
                return no_content()
            },
        };
        mac.update(&body);
        let expected_hex = hex::encode(mac.finalize().into_bytes());

        // constant-time equality check
        if expected_hex
            .as_bytes()
            .ct_eq(sig_hex.as_bytes())
            .unwrap_u8()
            != 1
        {
            return no_content();
        }
    }

    // get and check event header
    let event = match req.headers().get(GITHUB_EVENT) {
        Some(h) => h.to_str().unwrap_or("").to_string(),
        None => return no_content(),
    };

    // check credentials
    if !creds.is_valid() {
        return no_content();
    }

    tokio::spawn({
        let c = client.clone();
        let j = star_jail.clone();
        async move {
            if let Err(e) = parse_event(event, body, creds.into_inner(), &*c, &*j).await {
                error!("{}", e);
            }
        }
    });

    no_content()
}
