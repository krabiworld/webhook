use crate::errors::Error;
use crate::events::base::{Repository, WebhookMessage, Workflow, WorkflowRun};
use crate::events::Event;
use serde::Deserialize;

const IGNORED_WORKFLOWS: [&str; 3] = [
    "CodeQL",
    "Automatic Dependency Submission",
    "Dependabot Updates",
];

#[derive(Deserialize)]
pub struct WorkflowRunEvent {
    pub action: String,
    pub workflow: Workflow,
    pub workflow_run: WorkflowRun,
    pub repository: Repository,
}

impl Event for WorkflowRunEvent {
    fn handle(&self) -> Result<Option<WebhookMessage>, Error> {
        if self.action != "completed" {
            return Ok(None);
        }

        if let Some(conclusion) = &self.workflow_run.conclusion {
            if IGNORED_WORKFLOWS.contains(&self.workflow.name.as_str()) && conclusion == "success" {
                return Ok(None);
            }

            let emoji = if conclusion == "failure" {
                "<:catscream:1325122976575655936>"
            } else {
                "<:pepethinking:1330806911141941249>"
            };

            let branch_name = self
                .workflow_run
                .head_branch
                .as_deref()
                .unwrap_or("unknown");
            let formatted = format!(
                "{emoji} Workflow [{conclusion}](<{}>) on [{}](<{}>)/[{branch_name}](<{}/tree/{branch_name}>)",
                self.workflow_run.html_url,
                self.repository.name,
                self.repository.html_url,
                self.repository.html_url,
            );

            return Ok(Some(WebhookMessage {
                content: formatted,
                username: self.workflow.name.clone(),
                avatar_url: self.repository.owner.avatar_url.clone(),
            }));
        }

        Ok(None)
    }
}
