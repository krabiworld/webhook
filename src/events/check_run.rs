use crate::errors::Error;
use crate::events::Event;
use crate::events::base::{CheckRun, Repository, WebhookMessage};
use serde::Deserialize;

#[derive(Deserialize)]
pub struct CheckRunEvent {
    pub action: String,
    pub check_run: CheckRun,
    pub repository: Repository,
}

impl Event for CheckRunEvent {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error> {
        if self.action != "completed" {
            return Ok(None);
        }

        if self.repository.private {
            return Ok(None);
        }

        if let Some(conclusion) = &self.check_run.conclusion {
            if self.check_run.app.slug != "cloudflare-workers-and-pages" {
                return Ok(None);
            }

            let emoji = if conclusion == "failure" {
                "<:catscream:1325122976575655936>"
            } else {
                "<:pepethinking:1330806911141941249>"
            };

            let branch_name = self.check_run
                .check_suite
                .head_branch
                .as_deref()
                .unwrap_or("unknown");
            let formatted = format!(
                "{emoji} Check [{conclusion}](<{}>) on [{}](<{}>)/[{branch_name}](<{}/tree/{branch_name}>)",
                self.check_run.html_url,
                self.repository.name,
                self.repository.html_url,
                self.repository.html_url,
            );

            return Ok(Some(WebhookMessage {
                content: formatted,
                username: "Cloudflare Pages".to_string(),
                avatar_url: self.check_run.app.owner.avatar_url.clone(),
            }));
        }

        Ok(None)
    }
}
