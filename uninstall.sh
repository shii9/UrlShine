#!/usr/bin/env bash
# URLShine — Complete Uninstaller
# Usage: bash uninstall.sh
set -uo pipefail

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'

info() { echo -e "${CYAN}[INFO]${NC}  $*"; }
ok()   { echo -e "${GREEN}[ ✔ ]${NC}  $*"; }
warn() { echo -e "${YELLOW}[ ! ]${NC}  $*"; }

echo ""
echo -e "${CYAN}${BOLD}  ╔══════════════════════════════════════════╗"
echo -e "  ║         URLShine — Uninstaller          ║"
echo -e "  ╚══════════════════════════════════════════╝${NC}"
echo ""

# 1. Remove main urlshine binary
info "Removing URLShine binary..."
if [ -f "/usr/local/bin/urlshine" ]; then
    if sudo rm -f "/usr/local/bin/urlshine"; then
        ok "Removed /usr/local/bin/urlshine"
    else
        warn "Failed to remove /usr/local/bin/urlshine (requires sudo)"
    fi
else
    ok "URLShine binary already removed from /usr/local/bin"
fi

# 2. Ask to remove sub-tools installed by installer
echo -e "\n${BOLD}${YELLOW}Do you also want to remove the sub-tools installed by URLShine?${NC}"
echo "These include: gau, katana, gobuster, waymore, xnLinkFinder, httpx, gospider, waybackurls"
read -p "Remove these tools too? [y/N]: " remove_subs

if [[ "$remove_subs" =~ ^[Yy]$ ]]; then
    info "Removing sub-tools from /usr/local/bin..."
    tools=("gau" "gospider" "katana" "waybackurls" "gobuster" "httpx" "waymore" "xnLinkFinder")
    for t in "${tools[@]}"; do
        if [ -f "/usr/local/bin/$t" ]; then
            if sudo rm -f "/usr/local/bin/$t"; then
                ok "Removed /usr/local/bin/$t"
            else
                warn "Failed to remove /usr/local/bin/$t"
            fi
        else
            info "$t not found in /usr/local/bin, skipping"
        fi
    done
    
    # Optional python tools cleanup
    info "Uninstalling python packages via pip3..."
    if command -v pip3 &>/dev/null; then
        pip3 uninstall -y waymore xnlinkfinder --quiet > /dev/null 2>&1 || true
        ok "Removed waymore and xnlinkfinder pip packages"
    fi
else
    info "Keeping sub-tools."
fi

# 3. Clean up local repository directory instructions
echo -e "\n${BOLD}${GREEN}✓ System binaries cleaned up!${NC}"
echo -e "To completely remove the source folder, run the following command in your terminal:"
echo -e "${CYAN}cd .. && rm -rf UrlShine${NC}\n"
