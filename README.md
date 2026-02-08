# CmdVault

## The Problem CmdVault Solves

As developers, we constantly find ourselves:

- Scrolling through terminal history to find that command we ran last week
- Keeping scattered notes with commands across different tools and projects
- Forgetting the exact flags and options for complex commands
- Re-creating the same commands with small variations again and again
- Losing track of useful commands when switching between machines

**CmdVault** is your personal command-line Swiss Army knife - a modern, cross-platform CLI tool that helps you store, search, and execute frequently used shell commands with an intuitive interface.

## Key Features

- 🔍 **Fuzzy Search**: Quickly find commands with partial search terms
- 🧩 **Parameterized Commands**: Store command templates with placeholders for dynamic values
- 🗂️ **Collections**: Organize commands into different collections (e.g., Docker, Kubernetes, etc.)
- 📋 **Clipboard Integration**: Copy commands with a single keystroke
- ⌨️ **TUI Interface**: Beautiful terminal UI powered by Bubble Tea
- 🔄 **Import/Export**: Share command collections across machines via a GitHub repo
- 🌍 **Cross-Platform**: Works on macOS, Linux, and Windows
- 🚀 **Fast**: Written in Go for lightning-fast performance

## Installation

### Via Homebrew (macOS and Linux)

```bash
brew tap yourusername/cmdvault
brew install cmdvault
```

### Using Go

```bash
go install github.com/sakshamgupta05/cmdvault@latest
```

### Direct Download

Download the latest binary for your platform from the [releases page](https://github.com/sakshamgupta05/cmdvault/releases).

## Usage

### Interactive Search

Simply run `cmdvault` without arguments to open the interactive search interface:

```bash
cmdvault
```

### Add a New Command

```bash
cmdvault add
```

### Search for Commands

```bash
cmdvault search docker
```

### Manage Collections

```bash
# Create a new collection
cmdvault create-collection kubernetes

# Set the default collection
cmdvault set-default kubernetes

# List all collections
cmdvault collections
```

### Import/Export Commands

```bash
# Export all collections to a directory (e.g., to commit to Git)
cmdvault export ~/github/my-cmdvault

# Import commands from a directory
cmdvault import ~/github/my-cmdvault
```

## Parameterized Commands

CmdVault supports command templates with dynamic parameters, similar to function arguments.

### Parameter Syntax

- **Mandatory parameters**: `<param_name>`
- **Optional parameters**: Enclosed in square brackets `[ --flag <param_name>]`

### Example

```
helm template <repo>/<chart>[ --version <version>][ -f <values>.yaml] > <output>.yaml
```

In this example:
- `repo`, `chart`, and `output` are mandatory parameters
- `version` and `values` are optional parameters

When using this command, CmdVault will:
1. Prompt you to enter values for each parameter
2. For optional parameters, if you leave them empty, that entire section is omitted
3. Generate the final command with your values substituted

### Parameter Input

When executing a parameterized command, CmdVault displays a form to enter values:

<p align="center">
  <img src="assets/params-demo.png" alt="Parameter Input" width="600" />
</p>

You can preview the final command before execution with `Ctrl+P`.

## Examples

Here are some examples of how CmdVault can simplify your workflow:

### Kubernetes Context Switching

```
kubectl config use-context <context>
```

### Complex Docker Commands

```
docker run --name <container_name> -p <host_port>:<container_port>[ -v <volume>][ -e <env_var>=<value>] <image>:<tag>
```

### Git Workflow

```
git checkout -b <branch_name> && git add . && git commit -m "<commit_message>" && git push -u origin <branch_name>
```

### AWS CLI Commands

```
aws ec2 describe-instances --filters "Name=tag:<tag_key>,Values=<tag_value>"[ --region <region>]
```

## Keyboard Shortcuts

In the interactive search interface:

- **Tab**: Navigate between UI elements
- **Up/Down**: Navigate through commands
- **/** or just type: Fuzzy search
- **c**: Copy command to clipboard
- **e**: Execute command
- **v**: View command details
- **ESC**: Go back or clear search
- **q**: Quit

When entering parameters:

- **Tab/Shift+Tab**: Navigate between fields
- **Ctrl+P**: Toggle command preview
- **Enter**: Execute (when on the last field)
- **Ctrl+C/q**: Cancel

## Why CmdVault?

Unlike shell history or aliases, CmdVault:

- **Provides context**: Stores descriptions and tags with commands
- **Organizes**: Group commands into collections
- **Cross-platform**: Works across different operating systems
- **Git-ready**: Easy to back up and sync between machines
- **Parameter support**: Turns static commands into flexible templates

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/yourusername">yourusername</a>
</p>