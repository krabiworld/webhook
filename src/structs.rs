use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct Discord {
    #[serde(rename = "content")]
    pub content: String,

    #[serde(rename = "username")]
    pub username: String,

    #[serde(rename = "avatar_url")]
    pub avatar_url: String,
}

#[derive(Serialize, Deserialize)]
pub struct Credentials {
    pub id: String,
    pub token: String,
}

#[derive(Deserialize)]
pub struct PushCommit {
    pub id: String,
    pub url: String,
    pub message: String,
}

#[derive(Deserialize)]
pub struct PushEvent {
    pub commits: Vec<PushCommit>,
    #[serde(rename = "ref")]
    pub ref_: String,
    pub pusher: User,
    pub sender: User,
    pub repository: Repository,
}

#[derive(Deserialize)]
pub struct WorkflowRunEvent {
    pub action: String,
    pub workflow: Workflow,
    pub workflow_run: WorkflowRun,
    pub repository: Repository,
    pub sender: User,
}

#[derive(Deserialize)]
pub struct StarEvent {
    pub action: String,
    pub sender: User,
    pub repository: Repository,
}

#[derive(Deserialize)]
pub struct ForkEvent {
    pub sender: User,
    pub forkee: Forkee,
}

#[derive(Deserialize)]
pub struct Workflow {
    #[serde(default)]
    pub name: String,
}

#[derive(Deserialize)]
pub struct WorkflowRun {
    pub conclusion: Option<String>,
    #[serde(default)]
    pub html_url: String,
    pub head_branch: Option<String>,
}

#[derive(Deserialize)]
pub struct Repository {
    #[serde(default)]
    pub name: String,
    #[serde(default)]
    pub html_url: String,
}

#[derive(Deserialize)]
pub struct User {
    #[serde(default)]
    pub name: String,
    #[serde(default)]
    pub login: String,
    #[serde(default)]
    pub avatar_url: String,
    #[serde(default)]
    pub html_url: String,
}

#[derive(Deserialize)]
pub struct Forkee {
    #[serde(default)]
    pub name: String,
    #[serde(default)]
    pub html_url: String,
}
