#!/usr/bin/env bash
# URLShine PATH Configuration Helper
# Usage: bash setup-path.sh
# This script automatically configures your shell to access URLShine and all tools from anywhere

set -euo pipefail

CYAN='\033[0;36m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BOLD='\033[1m'; NC='\033[0m'

info()    { echo -e "${CYAN}[INFO]${NC}  $*"; }
ok()      { echo -e "${GREEN}[ ✔ ]${NC}  $*"; }
warn()    { echo -e "${YELLOW}[ ! ]${NC}  $*"; }
section() { echo -e "\n${BOLD}${CYAN}── $* ──${NC}\n"; }

echo ""
echo -e "${CYAN}${BOLD}  ╔═══════════════════════════════════════════╗"
echo -e "  ║  URLShine — PATH Configuration Helper    ║"
echo -e "  ╚═══════════════════════════════════════════╝${NC}"
echo ""

# Detect shell
if [[ "$SHELL" == *"zsh"* ]]; then
  SHELL_RC="$HOME/.zshrc"
  SHELL_NAME="Zsh"
elif [[ "$SHELL" == *"bash"* ]]; then
  SHELL_RC="$HOME/.bashrc"
  SHELL_NAME="Bash"
elif [[ "$SHELL" == *"fish"* ]]; then
  SHELL_RC="$HOME/.config/fish/config.fish"
  SHELL_NAME="Fish"
else
  SHELL_RC="$HOME/.bashrc"
  SHELL_NAME="Bash (default)"
fi

info "Detected shell: $SHELL_NAME ($SHELL_RC)"

section "Configuring PATH for System-Wide Access"

GOPATH=${GOPATH:-$HOME/go}
GOBIN=$GOPATH/bin
PYBIN=$HOME/.local/bin

info "Go binaries location: $GOBIN"
info "Python binaries location: $PYBIN"

# Function to add PATH based on shell type
add_path_to_shell() {
  local rc_file="$1"
  local shell_type="$2"
  
  if [[ "$shell_type" == "fish" ]]; then
    # Fish shell syntax
    if ! grep -q "set -Ua fish_user_paths.*go/bin" "$rc_file" 2>/dev/null; then
      info "Adding Go binaries to Fish PATH..."
      echo "" >> "$rc_file"
      echo "# URLShine Tools PATH Configuration" >> "$rc_file"
      echo "set -Ua fish_user_paths $GOBIN" >> "$rc_file"
      echo "set -Ua fish_user_paths $PYBIN" >> "$rc_file"
      ok "Added to $rc_file"
    else
      info "PATH already configured in $rc_file"
    fi
  else
    # Bash/Zsh syntax
    if ! grep -q "export PATH.*go/bin" "$rc_file" 2>/dev/null; then
      info "Adding Go and Python binaries to PATH..."
      echo "" >> "$rc_file"
      echo "# URLShine Tools PATH Configuration" >> "$rc_file"
      echo "export PATH=\"\$PATH:$GOBIN:$PYBIN\"" >> "$rc_file"
      ok "Added to $rc_file"
    else
      info "PATH already configured in $rc_file"
    fi
  fi
}

# Add PATH based on shell type
if [[ "$SHELL_NAME" == "Fish" ]]; then
  add_path_to_shell "$SHELL_RC" "fish"
else
  add_path_to_shell "$SHELL_RC" "bash"
fi

section "Verifying Tool Installation"

verify_tool() {
  if command -v "$1" &>/dev/null; then
    ok "$1 ✓"
  else
    warn "$1 ✗ (not yet accessible, may need shell reload)"
  fi
}

echo "Checking installed tools..."
verify_tool "urlshine"
verify_tool "gau"
verify_tool "katana"
verify_tool "gospider"
verify_tool "httpx"
verify_tool "waymore"
verify_tool "xnlinkfinder"

section "Next Steps"

echo -e "${BOLD}1. Reload your shell configuration:${NC}"
if [[ "$SHELL_NAME" == "Fish" ]]; then
  echo "   exec fish"
else
  echo "   source $SHELL_RC"
fi

echo ""
echo -e "${BOLD}2. Verify everything is accessible:${NC}"
echo "   urlshine doctor"
echo "   gau --help"
echo "   katana --help"

echo ""
echo -e "${BOLD}3. Test a URL scan:${NC}"
echo "   urlshine -a -c google.com"

echo ""
echo -e "${CYAN}${BOLD}✨ URLShine is now accessible from anywhere!${NC}"
echo ""
