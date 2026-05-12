#!/usr/bin/env bash
# URLShine — Professional Environment Installer
# Usage: bash install.sh
set -euo pipefail

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'

info()    { echo -e "${CYAN}[INFO]${NC}  $*"; }
ok()      { echo -e "${GREEN}[ ✔ ]${NC}  $*"; }
warn()    { echo -e "${YELLOW}[ ! ]${NC}  $*"; }
skip()    { echo -e "${YELLOW}[---]${NC}  $* (already installed)"; }
section() { echo -e "\n${BOLD}${CYAN}── $* ──${NC}\n"; }

check_go() {
  command -v go &>/dev/null || { echo -e "${RED}[✘] Go not found. Install: https://go.dev/dl${NC}"; exit 1; }
  info "Go $(go version | awk '{print $3}')"
  
  # Ensure GOPATH/bin is in PATH
  GOPATH=${GOPATH:-$HOME/go}
  GOBIN=$GOPATH/bin
  if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
    warn "GOPATH/bin ($GOBIN) not in PATH"
    info "Add this to your ~/.bashrc or ~/.zshrc:"
    echo "    export PATH=\"\$PATH:$GOBIN\""
  fi
}

go_install() {
  local bin="$1" pkg="$2"
  if command -v "$bin" &>/dev/null; then skip "$bin"; return; fi
  info "Installing $bin ..."
  go install "$pkg" && ok "$bin" || warn "Failed: $bin"
}

pip_install() {
  local bin="$1" pkg="$2"
  if command -v "$bin" &>/dev/null; then skip "$bin"; return; fi
  command -v pip3 &>/dev/null || { warn "pip3 not found, skip $bin"; return; }
  info "Installing $bin via pip ..."
  pip3 install "$pkg" --quiet && ok "$bin" || warn "Failed: $bin"
}

echo ""
echo -e "${CYAN}${BOLD}  ╔══════════════════════════════════════════╗"
echo -e "  ║     URLShine — Tool Installer  v2.0.0    ║"
echo -e "  ╚══════════════════════════════════════════╝${NC}"
echo ""

check_go

section "Go-based tools"
go_install "gau"         "github.com/lc/gau/v2/cmd/gau@latest"
go_install "gospider"    "github.com/jaeles-project/gospider@latest"
go_install "katana"      "github.com/projectdiscovery/katana/cmd/katana@latest"
go_install "waybackurls" "github.com/tomnomnom/waybackurls@latest"
go_install "hakrawler"   "github.com/hakluke/hakrawler@latest"
go_install "gobuster"    "github.com/OJ/gobuster/v3@latest"
go_install "httpx"       "github.com/projectdiscovery/httpx/cmd/httpx@latest"

section "Python-based tools"
pip_install "waymore"       "waymore"
pip_install "xnLinkFinder"  "xnlinkfinder"

section "System tools (dirb)"
install_dirb() {
  if command -v dirb &>/dev/null; then skip "dirb"; return; fi
  info "Installing dirb ..."
  if command -v apt-get &>/dev/null; then
    sudo apt-get update && sudo apt-get install -y dirb && ok "dirb" || warn "Failed: dirb"
  elif command -v brew &>/dev/null; then
    brew install dirb && ok "dirb" || warn "Failed: dirb"
  else
    warn "dirb requires apt or brew. Visit: https://sourceforge.net/projects/dirb/"
  fi
}
install_dirb

section "Building URLShine Binary"
info "Compiling URLShine ..."
go mod tidy
go build -ldflags "-X main.version=2.0.0 -s -w" -o urlshine .
ok "urlshine binary compiled"

section "Installing URLShine to System PATH"
if sudo cp urlshine /usr/local/bin/urlshine; then
  sudo chmod +x /usr/local/bin/urlshine
  ok "urlshine installed to /usr/local/bin"
  ok "Usage: urlshine --help"
  ok "        urlshine doctor"
else
  warn "Failed to install to /usr/local/bin (requires sudo)"
  info "You can still run: ./urlshine --help"
fi

section "Verifying Installation & PATH Configuration"

# Check if tools are accessible
verify_tool() {
  if command -v "$1" &>/dev/null; then
    ok "$1 ✓ accessible in PATH"
  else
    warn "$1 ✗ not found in PATH"
  fi
}

info "Verifying installed tools..."
verify_tool "urlshine"
verify_tool "gau"
verify_tool "katana"
verify_tool "gospider"
verify_tool "httpx"
verify_tool "waymore"
verify_tool "xnlinkfinder"

section "PATH Configuration"

# Get current shell
SHELL_RC=""
if [[ "$SHELL" == *"zsh"* ]]; then
  SHELL_RC="$HOME/.zshrc"
elif [[ "$SHELL" == *"bash"* ]]; then
  SHELL_RC="$HOME/.bashrc"
fi

if [ -n "$SHELL_RC" ]; then
  GOPATH=${GOPATH:-$HOME/go}
  GOBIN=$GOPATH/bin
  
  # Check if GOBIN is already in PATH
  if ! grep -q "export PATH.*$GOBIN" "$SHELL_RC" 2>/dev/null; then
    info "Adding $GOBIN to PATH in $SHELL_RC..."
    echo "" >> "$SHELL_RC"
    echo "# URLShine Go tools PATH" >> "$SHELL_RC"
    echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
    ok "PATH configuration added"
    info "Run: source $SHELL_RC"
  else
    ok "PATH already configured"
  fi
else
  warn "Could not determine shell configuration file"
  info "Manually add to ~/.bashrc or ~/.zshrc:"
  echo "    export PATH=\"\$PATH:\$HOME/go/bin\""
fi

section "Post-Installation Setup"

echo ""
echo -e "${BOLD}${CYAN}1. Reload shell configuration:${NC}"
echo "   source $SHELL_RC"
echo ""
echo -e "${BOLD}${CYAN}2. Verify installation:${NC}"
echo "   urlshine doctor"
echo ""
echo -e "${BOLD}${CYAN}3. Run your first scan:${NC}"
echo "   urlshine -a -c google.com"
echo ""
echo -e "${BOLD}${CYAN}4. For detailed usage:${NC}"
echo "   urlshine --help"
echo ""

echo ""
