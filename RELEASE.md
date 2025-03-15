## Building binary

```bash
# Initialize Go module (run once)
go mod init github.com/sakshamgupta05/cmdvault

# Install dependencies
go get github.com/spf13/cobra
go get github.com/AlecAivazis/survey/v2
go get github.com/atotto/clipboard
go get github.com/fatih/color
go get gopkg.in/yaml.v3

# Build the binary
go build -o cmdvault .

# Install locally
go install
```

## Homebrew Formula

```ruby
class Cmdvault < Formula
  desc "Store and retrieve shell commands"
  homepage "https://github.com/sakshamgupta05/cmdvault"
  url "https://github.com/sakshamgupta05/cmdvault/archive/v1.0.0.tar.gz"
  sha256 "YOUR_TARBALL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args
  end

  test do
    system "#{bin}/cmdvault", "--version"
  end
end
```

## Building for Different Platforms

```bash
# Create a build directory
mkdir -p build

# Build for macOS (amd64)
GOOS=darwin GOARCH=amd64 go build -o build/cmdvault-darwin-amd64

# Build for macOS (arm64)
GOOS=darwin GOARCH=arm64 go build -o build/cmdvault-darwin-arm64

# Build for Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o build/cmdvault-linux-amd64

# Build for Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o build/cmdvault-windows-amd64.exe
```

## Creating Debian/RPM Packages

```bash
# Install fpm
gem install fpm

# Create a deb package
fpm -s dir -t deb -n cmdvault -v 1.0.0 -C . \
  --prefix /usr/local/bin \
  build/cmdvault-linux-amd64=/usr/local/bin/cmdvault

# Create an rpm package
fpm -s dir -t rpm -n cmdvault -v 1.0.0 -C . \
  --prefix /usr/local/bin \
  build/cmdvault-linux-amd64=/usr/local/bin/cmdvault
```