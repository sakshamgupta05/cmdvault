mod config;
mod store;
mod ui;

use clap::{Parser, Subcommand};
use colored::*;
use std::process;

#[derive(Parser)]
#[command(name = "cmdvault")]
#[command(about = "Store and retrieve shell commands")]
#[command(
    long_about = "CmdVault helps you store and retrieve commonly used shell commands.\nStore commands with descriptions and tags, then search and use them when needed."
)]
struct Cli {
    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    /// List all commands in a collection
    List {
        /// Collection to list
        #[arg(short, long)]
        collection: String,
    },
    /// Search for commands
    Search {
        /// Search term
        term: Option<String>,
    },
    /// List all collections
    Collections,
}

fn main() {
    let cli = Cli::parse();

    let app_config = match config::init_config() {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Error initializing config: {}", e);
            process::exit(1);
        }
    };

    match cli.command {
        None => {
            // Default: interactive search
            ui::interactive_search(&app_config.collection_dirs, "");
        }
        Some(Commands::List { collection }) => {
            ui::list_commands(&app_config.collection_dirs, &collection);
        }
        Some(Commands::Search { term }) => {
            let search_term = term.unwrap_or_default();
            ui::interactive_search(&app_config.collection_dirs, &search_term);
        }
        Some(Commands::Collections) => match store::list_collections(&app_config.collection_dirs) {
            Ok(collections) => {
                println!("{} Available collections:\n", "\u{2022}".bold());
                for collection in &collections {
                    println!("  {}", collection);
                }
            }
            Err(e) => {
                eprintln!("Error: {}", e);
            }
        },
    }
}
