use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("Json Parse Error: {0}")]
    JsonParseError(#[from] serde_json::Error),

    #[error("Http Request Error: {0}")]
    HttpRequestError(#[from] reqwest::Error),

    #[error("Utf8 Error: {0}")]
    Utf8Error(#[from] std::str::Utf8Error),

    #[error("Discord Error: {0}")]
    DiscordError(String)
}
