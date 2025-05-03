use std::sync::Arc;

use crate::events::base::Credentials;
use crate::parser::parse_event;
use crate::{GITHUB_EVENT, GITHUB_SIG};
use actix_web::http::StatusCode;
use actix_web::{post, web, HttpRequest, HttpResponse, Responder};
use hmac::{Hmac, Mac};
use log::error;
use reqwest::Client;
use sha2::Sha256;
use subtle::ConstantTimeEq;

type HmacSha256 = Hmac<Sha256>;

fn no_content() -> HttpResponse {
    HttpResponse::Ok()
        .status(StatusCode::NO_CONTENT)
        .finish()
}

#[post("/{id}/{token}")]
pub async fn webhook(
    req: HttpRequest,
    creds: web::Path<Credentials>,
    body: web::Bytes,
    secret: web::Data<Arc<String>>,
    client: web::Data<Client>,
) -> impl Responder {
    // get and check github signature
    let sig = req.headers().get(GITHUB_SIG)
        .and_then(|h| h.to_str().ok())
        .filter(|s| s.starts_with("sha256="))
        .map(|s| &s[7..]);

    let sig_hex = match sig {
        Some(s) => s,
        None => return no_content(),
    };

    let mut mac = match HmacSha256::new_from_slice(secret.as_bytes()) {
        Ok(o) => o,
        Err(_) => return no_content(),
    };
    mac.update(&body);
    let expected_hex = hex::encode(mac.finalize().into_bytes());

    // constant-time equality check
    if expected_hex.as_bytes().ct_eq(sig_hex.as_bytes()).unwrap_u8() != 1 {
        return no_content();
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
        let c = client.get_ref().clone();
        async move {
            if let Err(e) = parse_event(event, body, creds.into_inner(), c).await {
                error!("{}", e);
            }
        }
    });

    no_content()
}
