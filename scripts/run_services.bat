@echo off
echo Starting KYC and AI services...

REM Create necessary directories if they don't exist
if not exist "uploads" mkdir uploads
if not exist "models" mkdir models

REM Check if Python is installed
python --version > nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Error: Python is not installed or not in PATH
    pause
    exit /b 1
)

REM Check if the AI server file exists
if not exist "services\ai-service\main.py" (
    echo Error: services\ai-service\main.py not found
    echo Current directory: %CD%
    dir services\ai-service
    pause
    exit /b 1
)

REM Start the AI service in a new window
echo Starting AI Service on port 8001...
start "AI Service" cmd /c "cd services\ai-service && pip install fastapi uvicorn python-multipart pydantic && python main.py || (echo AI Service failed to start && pause)"

REM Wait for the AI service to start
echo Waiting for AI service to start...
timeout /t 5 /nobreak > nul

REM Check if Go is installed
go version > nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed or not in PATH
    pause
    exit /b 1
)

REM Check if the KYC server file exists
if not exist "services\kyc-service\main.go" (
    echo Error: services\kyc-service\main.go not found
    echo Current directory: %CD%
    dir services\kyc-service
    pause
    exit /b 1
)

REM Start the KYC service in a new window
echo Starting KYC Service on port 8081...
start "KYC Service" cmd /c "cd services\kyc-service && go run main.go || (echo KYC Service failed to start && pause)"

REM Wait for the KYC service to start
echo Waiting for KYC service to start...
timeout /t 3 /nobreak > nul

REM Open browser to test the services
echo Opening browser to test the services...
start http://localhost:8001/docs
start http://localhost:8081/swagger-ui/

echo.
echo Services are running:
echo - AI Service: http://localhost:8001
echo - KYC Service: http://localhost:8081
echo.
echo API Key for AI service: your-api-key
echo.
echo Press any key to stop the services...
pause > nul

REM Kill the services
taskkill /FI "WINDOWTITLE eq AI Service*" /F
taskkill /FI "WINDOWTITLE eq KYC Service*" /F

echo Services stopped.
