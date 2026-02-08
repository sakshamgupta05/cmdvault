use anyhow::{Context, Result};
use serde::Deserialize;
use std::fs;
use std::path::PathBuf;

use crate::config;

#[derive(Deserialize, Clone, Debug)]
pub struct Parameter {
    #[serde(default)]
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub optional: bool,
    #[serde(default, rename = "defaultValue")]
    pub default_value: String,
}

#[derive(Deserialize, Clone, Debug)]
pub struct Command {
    #[serde(skip)]
    pub collection: String,
    #[serde(default)]
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub command: String,
    #[serde(default)]
    pub tags: Vec<String>,
    #[serde(default)]
    pub parameters: Vec<Parameter>,
}

#[derive(Deserialize)]
pub struct Collection {
    #[serde(default)]
    pub commands: Vec<Command>,
}

/// List all collection names by scanning collection directories for .yaml files.
pub fn list_collections(collection_dirs: &[PathBuf]) -> Result<Vec<String>> {
    let mut collections = Vec::new();

    for dir in collection_dirs {
        let entries = fs::read_dir(dir)
            .with_context(|| format!("Error reading commands directory: {}", dir.display()))?;

        for entry in entries {
            let entry = entry?;
            let path = entry.path();
            if !path.is_dir() {
                if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
                    if name.ends_with(".yaml") {
                        collections.push(name.trim_end_matches(".yaml").to_string());
                    }
                }
            }
        }
    }

    Ok(collections)
}

/// Get all commands from a specific collection.
pub fn get_commands(collection_dirs: &[PathBuf], collection: &str) -> Result<Vec<Command>> {
    let collection_path = config::get_collection_path(collection_dirs, collection)?;

    let data = fs::read_to_string(&collection_path)
        .with_context(|| format!("Error: could not read file {}", collection_path.display()))?;

    let mut col: Collection = serde_yaml::from_str(&data)
        .with_context(|| format!("Error: could not parse file {}", collection_path.display()))?;

    for cmd in &mut col.commands {
        cmd.collection = collection.to_string();
    }

    Ok(col.commands)
}

/// Get all commands from all collections.
pub fn get_all_commands(collection_dirs: &[PathBuf]) -> Result<Vec<Command>> {
    let collections = list_collections(collection_dirs)?;
    let mut all_commands = Vec::new();

    for collection in &collections {
        let cmds = get_commands(collection_dirs, collection)?;
        all_commands.extend(cmds);
    }

    Ok(all_commands)
}

/// Search commands by term (case-insensitive substring match on collection+name and tags).
pub fn search_commands(collection_dirs: &[PathBuf], search_term: &str) -> Result<Vec<Command>> {
    let commands = get_all_commands(collection_dirs)?;

    if search_term.is_empty() {
        return Ok(commands);
    }

    let term = search_term.to_lowercase();
    let mut results = Vec::new();

    for cmd in commands {
        let name_match = format!("{} {}", cmd.collection, cmd.name)
            .to_lowercase()
            .contains(&term);

        if name_match {
            results.push(cmd);
            continue;
        }

        let tag_match = cmd
            .tags
            .iter()
            .any(|tag| tag.to_lowercase().contains(&term));

        if tag_match {
            results.push(cmd);
        }
    }

    Ok(results)
}
