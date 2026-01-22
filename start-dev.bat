@echo off
setlocal
set ROOT=%~dp0


echo === Ensuring server deps (Bun) ===
pushd "%ROOT%Evolord-Server"
echo [server] bun install...
call bun install
popd

echo === Ensuring client dependencies (Go) ===
pushd "%ROOT%Evolord-Client"
if exist go.mod (
	echo [client] go mod tidy...
	go mod tidy
)
popd

echo === Launching windows ===
rem Bind server to all interfaces for remote access
set HOST=0.0.0.0
set PORT=5173
set EVOLORD_AGENT_TOKEN=dev-token-insecure-local-only

start "Evolord-Server" cmd /k "cd /d %ROOT%Evolord-Server && set "EVOLORD_AGENT_TOKEN=dev-token-insecure-local-only" && bun install && bun run dev"
start "Evolord-Client" cmd /k "cd /d %ROOT%Evolord-Client && set EVOLORD_SERVER=wss://localhost:5173 && set "EVOLORD_AGENT_TOKEN=dev-token-insecure-local-only" && set EVOLORD_TLS_INSECURE_SKIP_VERIFY=true && set EVOLORD_MODE=dev && set GOINSECURE=* && set GOSUMDB=off && set GOPROXY=https://proxy.golang.org,direct && go mod tidy && go run ./cmd/agent"

echo Done. Terminals stay open (/k) for logs.
endlocal
