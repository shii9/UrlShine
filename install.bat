@echo off
REM URLShine — Windows Environment Installer
REM Usage: install.bat

setlocal enabledelayedexpansion
set VERSION=2.0.0

color 0B
echo.
echo   ╔════════════════════════════════════════╗
echo   ║  URLShine - Tool Installer v%VERSION%  ║
echo   ╚════════════════════════════════════════╝
echo.

REM Check Go installation
where go >nul 2>&1
if errorlevel 1 (
    color 0C
    echo [✘] Go not found. Install from: https://go.dev/dl
    color 0B
    pause
    exit /b 1
)
for /f "tokens=3" %%g in ('go version') do set GO_VERSION=%%g
echo [INFO] Go %GO_VERSION%

REM Install Go-based tools
echo.
echo ── Go-based Tools ──
echo.

setlocal enabledelayedexpansion
for %%T in ("gau:github.com/lc/gau/v2/cmd/gau@latest" ^
           "gospider:github.com/jaeles-project/gospider@latest" ^
           "katana:github.com/projectdiscovery/katana/cmd/katana@latest" ^
           "waybackurls:github.com/tomnomnom/waybackurls@latest" ^
           "hakrawler:github.com/hakluke/hakrawler@latest" ^
           "gobuster:github.com/OJ/gobuster/v3@latest" ^
           "httpx:github.com/projectdiscovery/httpx/cmd/httpx@latest") do (
    for /f "tokens=1,2 delims=:" %%A in ("%%T") do (
        where %%A >nul 2>&1
        if not errorlevel 1 (
            echo [---] %%A ^(already installed^)
        ) else (
            echo [INFO] Installing %%A...
            call go install %%B
            if errorlevel 1 (
                color 0C
                echo [!] Failed: %%A
                color 0B
            ) else (
                color 0A
                echo [ ✔ ] %%A installed
                color 0B
            )
        )
    )
)

REM Install Python-based tools
echo.
echo ── Python-based Tools ──
echo.

where pip3 >nul 2>&1
if errorlevel 1 (
    where pip >nul 2>&1
    if errorlevel 1 (
        color 0C
        echo [!] Python/pip not found - skipping Python tools
        color 0B
    ) else (
        set PIP=pip
    )
) else (
    set PIP=pip3
)

if defined PIP (
    for %%P in (waymore xnlinkfinder) do (
        %PIP% show %%P >nul 2>&1
        if not errorlevel 1 (
            echo [---] %%P ^(already installed^)
        ) else (
            echo [INFO] Installing %%P...
            call %PIP% install %%P -q
            if errorlevel 1 (
                color 0C
                echo [!] Failed: %%P
                color 0B
            ) else (
                color 0A
                echo [ ✔ ] %%P installed
                color 0B
            )
        )
    )
)

REM Build URLShine
echo.
echo ── Building URLShine ──
echo.
go mod tidy
go build -ldflags "-s -w" -o urlshine.exe .
if errorlevel 1 (
    color 0C
    echo [✘] Build failed
    color 0B
    pause
    exit /b 1
)

color 0A
echo [ ✔ ] urlshine.exe built successfully
color 0B
echo.
echo Usage: urlshine.exe --help
echo.
pause
