@echo off
REM ============================================================================
REM 集成测试执行脚本 (Windows)
REM
REM 功能：
REM 1. 执行后端集成测试（计费、多租户、API）
REM 2. 执行前端E2E测试（知识管理、计费管理）
REM 3. 生成测试覆盖率报告
REM 4. 生成测试执行报告
REM
REM 作者: 研发团队
REM 日期: 2025-01-04
REM ============================================================================

setlocal enabledelayedexpansion

REM 项目根目录
set PROJECT_ROOT=%~dp0..\
set BACKEND_DIR=%PROJECT_ROOT%backend
set FRONTEND_DIR=%PROJECT_ROOT%frontend\apps\coze-studio
set REPORT_DIR=%PROJECT_ROOT%test-reports

REM 时间戳
for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set TIMESTAMP=%dt:~0,8%_%dt:~8,6%

REM 创建报告目录
if not exist "%REPORT_DIR%" mkdir "%REPORT_DIR%"

echo.
echo ========================================
echo   集成测试执行开始
echo ========================================
echo.

REM ============================================================================
REM 运行后端集成测试
REM ============================================================================

echo.
echo [INFO] 运行后端集成测试...
echo.

cd /d "%BACKEND_DIR%"

REM 启动测试容器
echo [INFO] 启动测试依赖容器...
docker compose -f ..\..\docker\docker-compose.test.yml up -d
timeout /t 10 /nobreak >nul

REM 运行计费集成测试
echo [INFO] 运行计费集成测试...
go test -v ./tests/integration/billing/... -coverprofile="%REPORT_DIR%\billing_coverage.out" -timeout 30m > "%REPORT_DIR%\backend_billing_test.log" 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] ✅ 计费集成测试通过
) else (
    echo [ERROR] ❌ 计费集成测试失败
)

REM 运行多租户集成测试
echo [INFO] 运行多租户集成测试...
go test -v ./tests/integration/tenant/... -coverprofile="%REPORT_DIR%\tenant_coverage.out" -timeout 30m > "%REPORT_DIR%\backend_tenant_test.log" 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] ✅ 多租户集成测试通过
) else (
    echo [ERROR] ❌ 多租户集成测试失败
)

REM 运行API集成测试
echo [INFO] 运行API集成测试...
go test -v ./tests/integration/api/... -coverprofile="%REPORT_DIR%\api_coverage.out" -timeout 30m > "%REPORT_DIR%\backend_api_test.log" 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] ✅ API集成测试通过
) else (
    echo [ERROR] ❌ API集成测试失败
)

REM 生成覆盖率报告
echo [INFO] 生成覆盖率报告...
go tool cover -html="%REPORT_DIR%\billing_coverage.out" -o "%REPORT_DIR%\billing_coverage.html"
go tool cover -html="%REPORT_DIR%\tenant_coverage.out" -o "%REPORT_DIR%\tenant_coverage.html"
go tool cover -html="%REPORT_DIR%\api_coverage.out" -o "%REPORT_DIR%\api_coverage.html"

echo.
echo [SUCCESS] 后端集成测试完成
echo.

REM ============================================================================
REM 运行前端E2E测试
REM ============================================================================

echo.
echo [INFO] 运行前端E2E测试...
echo.

cd /d "%FRONTEND_DIR%"

REM 安装依赖
echo [INFO] 安装前端测试依赖...
call npm ci
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] 前端依赖安装失败
    goto :cleanup
)

REM 运行计费管理E2E测试
echo [INFO] 运行计费管理E2E测试...
call npx playwright test billing-management.spec.ts --reporter=json 2>&1 | tee "%REPORT_DIR%\frontend_billing_e2e_test.log"
if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] ✅ 计费管理E2E测试通过
) else (
    echo [ERROR] ❌ 计费管理E2E测试失败
)

REM 运行知识管理E2E测试
echo [INFO] 运行知识管理E2E测试...
call npx playwright test knowledge-management.spec.ts --reporter=json 2>&1 | tee "%REPORT_DIR%\frontend_knowledge_e2e_test.log"
if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] ✅ 知识管理E2E测试通过
) else (
    echo [ERROR] ❌ 知识管理E2E测试失败
)

echo.
echo [SUCCESS] 前端E2E测试完成
echo.

REM ============================================================================
REM 清理
REM ============================================================================

:cleanup
echo.
echo [INFO] 清理资源...
cd /d "%PROJECT_ROOT%"
docker compose -f docker\docker-compose-test.yml down -v 2>nul

echo.
echo ========================================
echo   集成测试执行完成
echo ========================================
echo.
echo [SUCCESS] 所有测试已完成！
echo [INFO] 测试报告位置: %REPORT_DIR%
echo.

pause
