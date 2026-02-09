# CmdVault

A CLI tool to store, search, and execute your frequently used shell commands.

Stop digging through terminal history or scattered notes. CmdVault lets you save commands with descriptions, tags, and parameters, then find and run them instantly with fuzzy search.

![CI](https://github.com/sakshamgupta05/cmdvault/actions/workflows/release.yml/badge.svg)
![Crates.io](https://img.shields.io/crates/v/cmdvault)
![Downloads](https://img.shields.io/crates/d/cmdvault)
![License](https://img.shields.io/github/license/sakshamgupta05/cmdvault)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20linux-blue)
![Arch](https://img.shields.io/badge/arch-arm64%20%7C%20amd64-green)
![Homebrew](https://img.shields.io/badge/Homebrew-available-orange)

## Features

- **Fuzzy search** - Find commands by typing any part of the name, collection, or tag
- **Parameterized commands** - Save command templates with `{{placeholders}}` that prompt you for values
- **Optional parameters** - Wrap sections in `{? ?}` to include them only when a value is provided
- **Collections** - Organize commands into YAML files by topic (docker, kubernetes, git, etc.)
- **Execute or copy** - Run commands directly or copy them to your clipboard
- **Preview** - See full command details in a side panel while searching
- **Cross-platform** - Works on macOS and Linux

## Installation

### Homebrew (macOS and Linux)

```bash
brew tap sakshamgupta05/tap
brew install cmdvault
```

### Cargo (Rust)

```bash
cargo install cmdvault
```

### Binary download

Download the latest binary for your platform from the [releases page](https://github.com/sakshamgupta05/cmdvault/releases).

**macOS users:** After downloading, you may need to remove the quarantine attribute:

```bash
tar -xzf cmdvault-darwin-arm64.tar.gz
xattr -d com.apple.quarantine ./cmdvault-darwin-arm64
./cmdvault-darwin-arm64
```

### Build from source

If you want to build and install from source locally:

```bash
# Clone the repository
git clone https://github.com/sakshamgupta05/cmdvault.git
cd cmdvault

# Build and install with cargo
cargo install --path .
```

Or to just build without installing:

```bash
# Build in release mode
cargo build --release

# The binary will be in target/release/cmdvault
./target/release/cmdvault
```

## Quick Start

### 1. Create a collection file

CmdVault reads commands from YAML files in `~/.config/cmdvault/collections/`. Create your first collection:

```bash
mkdir -p ~/.config/cmdvault/collections
```

Create `~/.config/cmdvault/collections/docker.yaml`:

```yaml
commands:
  - name: Run container
    command: docker run --name {{name}} -p {{host_port}}:{{container_port}} {{image}}:{{tag}}
    description: Run a Docker container with port mapping
    tags:
      - docker
      - run
    parameters:
      - name: name
        description: Container name
      - name: host_port
        description: Host port
      - name: container_port
        description: Container port
      - name: image
        description: Image name
      - name: tag
        description: Image tag
        defaultValue: latest

  - name: Stop all containers
    command: docker stop $(docker ps -q)
    description: Stop all running containers
    tags:
      - docker
      - stop
```

### 2. Search and run

```bash
cmdvault
```

This opens the fuzzy finder. Type to filter, press Enter to select, then choose to execute, copy, or view details.

## Usage

### Interactive search (default)

```bash
cmdvault              # browse all commands
cmdvault search git   # pre-filter by "git"
```

### List commands in a collection

```bash
cmdvault list -c docker
```

### List all collections

```bash
cmdvault collections
```

## Command YAML Format

Each collection is a YAML file with a `commands` list:

```yaml
commands:
  - name: Command name
    command: the actual shell command
    description: What this command does # optional
    tags: # optional
      - tag1
      - tag2
    parameters: # optional
      - name: param_name
        description: What this param is
        optional: false # default: false
        defaultValue: some_default # optional
```

### Parameters

Use `{{param_name}}` as placeholders in your command string. When you execute or copy a parameterized command, CmdVault prompts you for each value.

```yaml
- name: SSH into server
  command: ssh {{user}}@{{host}} -p {{port}}
  parameters:
    - name: user
      description: Username
      defaultValue: root
    - name: host
      description: Server hostname
    - name: port
      description: SSH port
      defaultValue: 22
```

### Optional parameters

Wrap sections in `{? ?}` to include them only when a value is provided. If left empty, the entire block is removed from the command.

```yaml
- name: Helm template
  command: helm template {{repo}}/{{chart}}{? --version {{version}} ?}{? -f {{values}}.yaml?} > {{output}}.yaml
  parameters:
    - name: repo
      description: Helm repo name
    - name: chart
      description: Chart name
    - name: version
      description: Chart version
      optional: true
    - name: values
      description: Values file name
      optional: true
    - name: output
      description: Output file name
```

In this example, if you leave `version` empty, `--version ...` is omitted entirely from the final command.

## Configuration

CmdVault works out of the box with zero configuration. By default, collections are stored in `~/.config/cmdvault/collections/`.

You can optionally create `~/.config/cmdvault/config.yaml` for advanced configuration. The config directory location can be overridden with the `CMDVAULT_CONFIG` environment variable.

## Advanced Configuration

Most users don't need a config file. However, you can create `~/.config/cmdvault/config.yaml` for advanced customization.

### Custom Collection Directories

You can configure CmdVault to load collections from multiple directories. This is useful for:

- **Importing public collections** - Clone community-maintained command repositories and reference them directly
- **Team collaboration** - Share collections via Git across your team
- **Organized separation** - Keep personal and work collections separate

Create `~/.config/cmdvault/config.yaml`:

```yaml
collectionDirs:
  - ~/.config/cmdvault-awesome-public/collections/
  - ~/.config/cmdvault/collections/
```

For example, to use a public collection repository:

```bash
# Clone a public collection repository
git clone https://github.com/someone/cmdvault-awesome-commands ~/.config/cmdvault-awesome-public

# Add it to your config
cat > ~/.config/cmdvault/config.yaml <<EOF
collectionDirs:
  - ~/.config/cmdvault-awesome-public/collections/
  - ~/.config/cmdvault/collections/
EOF
```

#### Collection Name Conflicts

Collections are identified by their filename (e.g., `docker.yaml` defines the "docker" collection). If multiple directories contain files with the same name, **the directory listed later in `collectionDirs` takes precedence** and overrides earlier ones.

In the example above, if both directories have a `docker.yaml` file, the one in `~/.config/cmdvault/collections/docker.yaml` will be used (since it's listed last). This allows you to selectively override commands from public collections with your own customizations.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
