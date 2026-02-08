use anyhow::{Context, Result};
use serde::Deserialize;
use std::env;
use std::fs;
use std::path::PathBuf;

#[derive(Deserialize, Default)]
pub struct Config {
    #[serde(default, rename = "collectionDirs")]
    pub collection_dirs: Vec<String>,
}

pub struct AppConfig {
    #[allow(dead_code)]
    pub config_dir: PathBuf,
    pub collection_dirs: Vec<PathBuf>,
}

pub fn init_config() -> Result<AppConfig> {
    let home_dir = dirs::home_dir().context("Could not find home directory")?;

    let config_dir = match env::var("CMDVAULT_CONFIG") {
        Ok(dir) if !dir.is_empty() => PathBuf::from(dir),
        _ => home_dir.join(".config").join("cmdvault"),
    };

    let config_file = config_dir.join("config.yaml");

    // Create config directory if it doesn't exist
    fs::create_dir_all(&config_dir).with_context(|| {
        format!(
            "Could not create config directory: {}",
            config_dir.display()
        )
    })?;

    // Load or create default config
    let config: Config = if config_file.exists() {
        let data = fs::read_to_string(&config_file)
            .with_context(|| format!("Error reading config file: {}", config_file.display()))?;
        serde_yaml::from_str(&data)
            .with_context(|| format!("Error parsing config file: {}", config_file.display()))?
    } else {
        Config::default()
    };

    let collection_dirs: Vec<PathBuf> = if config.collection_dirs.is_empty() {
        vec![config_dir.join("collections")]
    } else {
        config.collection_dirs.iter().map(PathBuf::from).collect()
    };

    // Create collection directories if they don't exist
    for dir in &collection_dirs {
        fs::create_dir_all(dir)
            .with_context(|| format!("Could not create collection directory: {}", dir.display()))?;
    }

    Ok(AppConfig {
        config_dir,
        collection_dirs,
    })
}

/// Find the path to a collection YAML file, searching dirs in reverse order.
pub fn get_collection_path(collection_dirs: &[PathBuf], collection: &str) -> Result<PathBuf> {
    let filename = format!("{}.yaml", collection);

    for dir in collection_dirs.iter().rev() {
        let file_path = dir.join(&filename);
        if file_path.exists() {
            return Ok(file_path);
        }
    }

    anyhow::bail!(
        "collection {} not found in any configured paths",
        collection
    )
}
