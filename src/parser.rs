use actix_web::web;
use log::error;
use regex::Regex;
use crate::client::execute_webhook;
use crate::structs::{Credentials, ForkEvent, PushEvent, StarEvent, WorkflowRunEvent};

pub async fn parse_event(event: String, body: web::Bytes, creds: Credentials) {
    let mut content = String::new();
    let mut username = String::new();
    let mut avatar = String::new();

    match event.as_str() {
        "push" => {
            let parsed = serde_json::from_slice::<PushEvent>(&body);
            match parsed {
                Ok(e) => {
                    let re = Regex::new(r"(?m)^\s*\n").unwrap();
                    let mut commits = String::new();

                    for c in &e.commits {
                        commits.push_str(&format!(
                            "[`{}`](<{}>) {}\n",
                            &c.id[..7],
                            c.url,
                            re.replace_all(&c.message, "").to_string()
                        ));
                    }

                    let branch = e.ref_.strip_prefix("refs/heads/").unwrap_or(&e.ref_);
                    let footer = format!(
                        "\n- [{}](<{}>) on [{}](<{}>)/[{}](<{}>)",
                        e.pusher.name,
                        e.sender.html_url,
                        e.repository.name,
                        e.repository.html_url,
                        branch,
                        format!("{}/tree/{}", e.repository.html_url, branch),
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

                    content = format!("{}{}", commits, footer);
                    username = e.pusher.name.to_string();
                    avatar = e.sender.avatar_url.to_string();
                }
                Err(err) => {
                    error!("{}", err);
                    return;
                }
            }
        }
        "workflow_run" => {
            let parsed = serde_json::from_slice::<WorkflowRunEvent>(&body);
            match parsed {
                Ok(e) => {
                    if e.action != "completed" {
                        return;
                    }

                    if let Some(conclusion) = &e.workflow_run.conclusion {
                        if (e.workflow.name.starts_with("CodeQL") || e.workflow.name == "Dependabot Updates")
                            && conclusion == "success"
                        {
                            return;
                        }

                        let emoji = if conclusion == "failure" {
                            "<:catscream:1325122976575655936>"
                        } else {
                            "<:pepethinking:1330806911141941249>"
                        };

                        let formatted = format!(
                            "{} Workflow [{}](<{}>) on [{}](<{}>)/[{}](<{}>)",
                            emoji,
                            conclusion,
                            e.workflow_run.html_url,
                            e.repository.name,
                            e.repository.html_url,
                            e.workflow_run.head_branch.as_deref().unwrap_or("unknown"),
                            format!("{}/tree/{}", e.repository.html_url, e.workflow_run.head_branch.as_deref().unwrap_or("unknown"))
                        );

                        content = formatted;
                        username = e.workflow.name.clone();
                        avatar = e.sender.avatar_url.clone();
                    }
                }
                Err(err) => {
                    error!("{}", err);
                    return;
                }
            }
        }
        "star" => {
            let parsed = serde_json::from_slice::<StarEvent>(&body);
            match parsed {
                Ok(e) => {
                    if e.action != "created" {
                        return;
                    }
                    content = format!(
                        "[{}](<{}>) starred [{}](<{}>) <:foxtada:1311327105300172882>",
                        e.sender.login, e.sender.html_url,
                        e.repository.name, e.repository.html_url
                    );
                    username = e.sender.login.clone();
                    avatar = e.sender.avatar_url.clone();
                }
                Err(err) => {
                    error!("{}", err);
                    return;
                }
            }
        }
        "fork" => {
            let parsed = serde_json::from_slice::<ForkEvent>(&body);
            match parsed {
                Ok(e) => {
                    content = format!(
                        "[{}](<{}>) forked [{}](<{}>)",
                        e.sender.login, e.sender.html_url,
                        e.forkee.name, e.forkee.html_url
                    );
                    username = e.sender.login.clone();
                    avatar = e.sender.avatar_url.clone();
                }
                Err(err) => {
                    error!("{}", err);
                    return;
                }
            }
        }
        _ => {}
    }

    if !content.is_empty() {
        execute_webhook(content, username, avatar, creds).await;
    }
}
