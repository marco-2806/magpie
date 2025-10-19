@echo off
setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
pushd "%SCRIPT_DIR%\.."

where docker >nul 2>&1
if errorlevel 1 (
  echo Docker is required but was not found. Install Docker Desktop.
  popd
  exit /b 1
)

docker compose version >nul 2>&1
if not errorlevel 1 (
  set "COMPOSE_EXE=docker"
  set "COMPOSE_ARGS=compose"
) else (
  where docker-compose >nul 2>&1
  if errorlevel 1 (
    echo Docker Compose is required but was not found. Install Docker Desktop or docker-compose.
    popd
    exit /b 1
  )
  set "COMPOSE_EXE=docker-compose"
  set "COMPOSE_ARGS="
)

echo Pulling latest changes...
git pull --ff-only
if errorlevel 1 (
  popd
  exit /b %ERRORLEVEL%
)

echo Rebuilding frontend and backend containers...
if "%COMPOSE_ARGS%"=="" (
  call "%COMPOSE_EXE%" up -d --build frontend backend
) else (
  call "%COMPOSE_EXE%" %COMPOSE_ARGS% up -d --build frontend backend
)

if errorlevel 1 (
  popd
  exit /b %ERRORLEVEL%
)

echo Done. Frontend is available at http://localhost:8080 and backend API at http://localhost:8082/api
popd
exit /b 0
