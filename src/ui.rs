use colored::*;
use inquire::{Select, Text};
use regex::Regex;
use skim::prelude::*;
use std::path::PathBuf;
use std::process;
use std::sync::Arc;

use crate::store::{self, Command, Parameter};

/// A wrapper around Command that implements SkimItem for fuzzy finder display and preview.
struct CommandItem {
    display: String,
    preview: String,
}

impl SkimItem for CommandItem {
    fn text(&self) -> Cow<'_, str> {
        Cow::Borrowed(&self.display)
    }

    fn preview(&self, _context: PreviewContext) -> ItemPreview {
        ItemPreview::Text(self.preview.clone())
    }
}

/// Display all commands in a collection.
pub fn list_commands(collection_dirs: &[PathBuf], collection: &str) {
    let commands = match store::get_commands(collection_dirs, collection) {
        Ok(cmds) => cmds,
        Err(e) => {
            eprintln!("Error: {}", e);
            return;
        }
    };

    if commands.is_empty() {
        println!("{}", "No commands found.".yellow());
        return;
    }

    println!(
        "{} Commands in collection \"{}\":\n",
        "\u{2022}".bold(),
        collection
    );

    for cmd in &commands {
        println!(
            "{} {}: {}\n   {}\n",
            "\u{2022}".bold(),
            cmd.collection,
            cmd.name.green(),
            cmd.tags.join(", ")
        );
    }
}

/// Interactive search using skim fuzzy finder.
pub fn interactive_search(collection_dirs: &[PathBuf], search_term: &str) {
    let commands = match store::search_commands(collection_dirs, search_term) {
        Ok(cmds) => cmds,
        Err(e) => {
            eprintln!("Error: {}", e);
            return;
        }
    };

    if commands.is_empty() {
        println!("{}", "No commands found.".yellow());
        return;
    }

    // Build display strings and skim items
    let display_strings: Vec<String> = commands
        .iter()
        .map(|cmd| {
            let tags_text = if cmd.tags.is_empty() {
                String::new()
            } else {
                format!(" [{}]", cmd.tags.join(", "))
            };
            format!("{}: {}{}", cmd.collection, cmd.name, tags_text)
        })
        .collect();

    let items: Vec<Arc<dyn SkimItem>> = commands
        .iter()
        .enumerate()
        .map(|(i, cmd)| {
            Arc::new(CommandItem {
                display: display_strings[i].clone(),
                preview: format_command_plain(cmd),
            }) as Arc<dyn SkimItem>
        })
        .collect();

    let options = SkimOptionsBuilder::default()
        .height(Some("100%"))
        .multi(false)
        .reverse(true)
        .preview(Some(""))
        .build()
        .unwrap();

    let (tx, rx): (SkimItemSender, SkimItemReceiver) = unbounded();
    for item in items {
        let _ = tx.send(item);
    }
    drop(tx);

    let output = Skim::run_with(&options, Some(rx));

    let selected_index = match output {
        Some(out) if !out.is_abort && !out.selected_items.is_empty() => {
            let selected_text = out.selected_items[0].output().to_string();
            // Match back to our display strings to find the index
            display_strings
                .iter()
                .position(|d| *d == selected_text)
                .unwrap_or_else(|| {
                    eprintln!("Error: could not match selected item");
                    process::exit(1);
                })
        }
        _ => return,
    };

    let selected = &commands[selected_index];

    // Prompt for action
    let action_options = vec!["Execute command", "Copy to clipboard", "Show details"];
    let action = match Select::new("What would you like to do?", action_options).prompt() {
        Ok(a) => a,
        Err(_) => return,
    };

    match action {
        "Execute command" => {
            let mut cmd_str = selected.command.clone();
            if !selected.parameters.is_empty() {
                cmd_str = interactive_parameters(selected);
            }
            println!("{} {}", "Executing:".yellow(), cmd_str);
            execute_command(&cmd_str);
        }
        "Copy to clipboard" => {
            let mut cmd_str = selected.command.clone();
            if !selected.parameters.is_empty() {
                cmd_str = interactive_parameters(selected);
            }
            match cli_clipboard::set_contents(cmd_str) {
                Ok(_) => println!("{}", "Command copied to clipboard!".green()),
                Err(e) => eprintln!("Error copying to clipboard: {}", e),
            }
        }
        "Show details" => {
            println!("\n{}", format_command_long(selected));
        }
        _ => {}
    }
}

/// Prompt user for parameter values and return the command with substitutions applied.
fn interactive_parameters(cmd: &Command) -> String {
    let mut cmd_str = cmd.command.clone();

    for param in &cmd.parameters {
        let default_str = if !param.default_value.is_empty() {
            format!(" ({})", param.default_value)
        } else {
            String::new()
        };

        let mandatory_str = if !param.optional && param.default_value.is_empty() {
            "*".bold().to_string()
        } else {
            String::new()
        };

        let prompt_str = format!("{}{}{}:", param.name, mandatory_str, default_str);

        let is_required = !param.optional && param.default_value.is_empty();

        let value = if is_required {
            loop {
                match Text::new(&prompt_str).prompt() {
                    Ok(v) if !v.trim().is_empty() => break v,
                    Ok(_) => {
                        eprintln!("This field is required.");
                        continue;
                    }
                    Err(e) => {
                        eprintln!("Error reading input: {}", e);
                        process::exit(1);
                    }
                }
            }
        } else {
            match Text::new(&prompt_str).prompt() {
                Ok(v) => v,
                Err(e) => {
                    eprintln!("Error reading input: {}", e);
                    process::exit(1);
                }
            }
        };

        cmd_str = replace_parameter(&cmd_str, param, &value);
    }

    cmd_str
}

/// Replace a parameter placeholder in the command string.
fn replace_parameter(cmd_str: &str, p: &Parameter, value: &str) -> String {
    if p.optional {
        let placeholder = format!("{{{{{}}}}}", p.name);
        // Match {? ... {{param}} ... ?}
        let pattern = format!(
            r"\{{\?(?:[^?]|\?(?:[^}}]|$))*\{{\{{{}\}}\}}(?:[^?]|\?(?:[^}}]|$))*\?\}}",
            regex::escape(&p.name)
        );
        let re = Regex::new(&pattern).unwrap();

        if value.is_empty() {
            re.replace_all(cmd_str, "").to_string()
        } else {
            let mut result = cmd_str.to_string();
            for mat in re.find_iter(cmd_str) {
                let original = mat.as_str();
                let mut new_text = original.replacen("{?", "", 1);
                if let Some(pos) = new_text.rfind("?}") {
                    new_text = format!("{}{}", &new_text[..pos], &new_text[pos + 2..]);
                }
                new_text = new_text.replace(&placeholder, value);
                result = result.replace(original, &new_text);
            }
            result
        }
    } else {
        let actual_value = if value.is_empty() {
            &p.default_value
        } else {
            value
        };
        let placeholder = format!("{{{{{}}}}}", p.name);
        cmd_str.replace(&placeholder, actual_value)
    }
}

/// Format command details for the preview window (plain text, no ANSI colors).
fn format_command_plain(cmd: &Command) -> String {
    let tags_text = if cmd.tags.is_empty() {
        String::new()
    } else {
        format!("\n   [{}]", cmd.tags.join(", "))
    };

    let description_text = if cmd.description.is_empty() {
        String::new()
    } else {
        format!("\n   {}", cmd.description)
    };

    let parameters_text = if cmd.parameters.is_empty() {
        String::new()
    } else {
        let mut text = "\n\nParameters:".to_string();
        for param in &cmd.parameters {
            text.push_str(&format!("\n   {}: {}", param.name, param.description));
        }
        text
    };

    format!(
        "Description:\n   {} {}{}{}\n\nCommand:\n   {}{}",
        cmd.collection, cmd.name, description_text, tags_text, cmd.command, parameters_text
    )
}

/// Format command details for display with colors (used in "Show details").
fn format_command_long(cmd: &Command) -> String {
    let tags_text = if cmd.tags.is_empty() {
        String::new()
    } else {
        format!("\n   [{}]", cmd.tags.join(", "))
    };

    let description_text = if cmd.description.is_empty() {
        String::new()
    } else {
        format!("\n   {}", cmd.description)
    };

    let parameters_text = if cmd.parameters.is_empty() {
        String::new()
    } else {
        let mut text = format!("\n\n{}", "Parameters:".bold());
        for param in &cmd.parameters {
            text.push_str(&format!("\n   {}: {}", param.name, param.description));
        }
        text
    };

    format!(
        "{}\n   {} {}{}{}\n\n{}\n   {}{}",
        "Description:".bold(),
        cmd.collection,
        cmd.name.yellow(),
        description_text,
        tags_text,
        "Command:".bold(),
        cmd.command.cyan(),
        parameters_text
    )
}

/// Execute a shell command, connecting stdin/stdout/stderr to the parent process.
fn execute_command(command: &str) {
    let result = if cfg!(target_os = "windows") {
        process::Command::new("cmd")
            .args(["/C", command])
            .stdin(process::Stdio::inherit())
            .stdout(process::Stdio::inherit())
            .stderr(process::Stdio::inherit())
            .spawn()
    } else {
        let shell = std::env::var("SHELL").unwrap_or_else(|_| "sh".to_string());
        process::Command::new(&shell)
            .args(["-c", command])
            .stdin(process::Stdio::inherit())
            .stdout(process::Stdio::inherit())
            .stderr(process::Stdio::inherit())
            .spawn()
    };

    match result {
        Ok(mut child) => match child.wait() {
            Ok(status) => {
                if let Some(code) = status.code() {
                    if code != 0 {
                        process::exit(code);
                    }
                }
            }
            Err(e) => {
                eprintln!("Error executing command: {}", e.to_string().red());
                process::exit(1);
            }
        },
        Err(e) => {
            eprintln!("Error executing command: {}", e.to_string().red());
            process::exit(1);
        }
    }
}
