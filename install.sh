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

echo ""
