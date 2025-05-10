use crate::errors::Error;
use crate::events::Event;
use crate::events::base::{PushCommit, Repository, User, WebhookMessage};
use regex::Regex;
use serde::Deserialize;

#[derive(Deserialize)]
pub struct PushEvent {
    pub commits: Vec<PushCommit>,
    #[serde(rename = "ref")]
    pub ref_: String,
    pub pusher: User,
    pub sender: User,
    pub repository: Repository,
}

impl Event for PushEvent {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error> {
        if self.commits.len() == 0 {
            return Ok(None);
        }

        let link_re = Regex::new(r"\[([^]]+)]\((https?://[^)]+)\)")?;
        let md_re = Regex::new(r"(?m)^\s*#{1,3}\s+")?;
        let mut commits = String::new();

        for c in &self.commits {
            let mut lines = c.message.lines();
            let first = lines.next().unwrap_or("");
            let commit_msg = if lines.next().is_some() {
                format!("{first}...")
            } else {
                first.to_string()
            };

            let clean_link_msg = link_re
                .replace_all(&commit_msg, |caps: &regex::Captures| {
                    format!("[{}](<{}>)", &caps[1], &caps[2])
                })
                .to_string();
            let clean_md_msg = md_re.replace_all(&clean_link_msg, "").to_string();

            commits.push_str(&format!(
                "[`{}`](<{}>) {}\n",
                &c.id[..7],
                c.url,
                clean_md_msg
            ));
        }

        let branch = self.ref_.strip_prefix("refs/heads/").unwrap_or(&self.ref_);
        let footer = format!(
            "\n- [{}](<{}>) on [{}](<{}>)/[{}](<{}/tree/{}>)",
            self.pusher.name,
            self.sender.html_url,
            self.repository.name,
            self.repository.html_url,
            branch,
            self.repository.html_url,
            branch,
        );

        let limit = 2000 - (footer.chars().count() + "...".len() + 1);
        if commits.chars().count() > limit {
            let mut truncated = commits.chars().take(limit).collect::<String>() + "...";
            if !truncated.ends_with(">)") {
                let lines: Vec<&str> = truncated.lines().collect();
                truncated = lines[..lines.len().saturating_sub(1)].join("\n");
            }
            commits = truncated + "\n";
        }

        Ok(Some(WebhookMessage {
            content: format!("{}{}", commits, footer),
            username: self.pusher.name.clone(),
            avatar_url: self.sender.avatar_url.clone(),
        }))
    }
}
