@echo off
REM 集成测试执行脚本 (Windows)
REM 职责: 执行所有集成测试并生成覆盖率报告

setlocal enabledelayedexpansion

REM 获取脚本目录
set SCRIPT_DIR=%~dp0
set BACKEND_DIR=%SCRIPT_DIR%..\..
cd /d "%BACKEND_DIR%"

echo [INFO] 当前工作目录: %BACKEND_DIR%

REM 默认参数
set CLEANUP_ONLY=0
set VERBOSE=0
set SUITE_NAME=
set GENERATE_COVERAGE=1

REM 解析参数
:parse_args
if "%~1"=="" goto end_parse
if "%~1"=="-h" goto print_usage
if "%~1"=="--help" goto print_usage
if "%~1"=="-c" (
    set CLEANUP_ONLY=1
    shift
    goto parse_args
)
if "%~1"=="--cleanup-only" (
    set CLEANUP_ONLY=1
    shift
    goto parse_args
)
if "%~1"=="-v" (
    set VERBOSE=1
    shift
    goto parse_args
)
if "%~1"=="--verbose" (
    set VERBOSE=1
    shift
    goto parse_args
)
if "%~1"=="-s" (
    set SUITE_NAME=%~2
    shift
    shift
    goto parse_args
)
if "%~1"=="--suite" (
    set SUITE_NAME=%~2
    shift
    shift
    goto parse_args
)
if "%~1"=="--no-coverage" (
    set GENERATE_COVERAGE=0
    shift
    goto parse_args
)
echo [ERROR] 未知参数: %~1
goto print_usage

:end_parse

REM 仅清理模式
if %CLEANUP_ONLY%==1 (
    echo [INFO] 清理测试资源...
    docker ps -a --filter "label=org.testcontainers" --format "{{{{.ID}}}}" > temp_containers.txt
    for /f "tokens=*" %%i in (temp_containers.txt) do (
        if not "%%i"=="" (
            echo [INFO] 删除测试容器: %%i
            docker rm -f %%i >nul 2>&1
        )
    )
    del temp_containers.txt
    echo [SUCCESS] 测试资源清理完成
    goto end
)

echo [INFO] ==================== 集成测试开始 ====================
echo.

REM 检查依赖
echo [INFO] 检查依赖...
where go >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go未安装
    exit /b 1
)

where docker >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker未安装（testcontainers需要）
    exit /b 1
)

docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker未运行
    exit /b 1
)

echo [SUCCESS] 依赖检查通过

REM 设置测试环境
echo [INFO] 设置测试环境...
set GO111MODULE=on
set CGO_ENABLED=1
set TEST_ENV=integration
set LOG_LEVEL=debug
echo [SUCCESS] 测试环境设置完成

REM 运行集成测试
echo [INFO] 开始运行集成测试...

set GO_TEST_CMD=go test

if %VERBOSE%==1 (
    set GO_TEST_CMD=%GO_TEST_CMD% -v
)

set GO_TEST_CMD=%GO_TEST_CMD% -race
set GO_TEST_CMD=%GO_TEST_CMD% -coverprofile=coverage.out -covermode=atomic
set GO_TEST_CMD=%GO_TEST_CMD% -timeout 10m

if not "%SUITE_NAME%"=="" (
    goto run_suite
)

REM 运行所有测试
echo [INFO] 执行命令: %GO_TEST_CMD% ./tests/integration/...
%GO_TEST_CMD% ./tests/integration/...
set TEST_RESULT=%errorlevel%

if %TEST_RESULT%==0 (
    echo [SUCCESS] 集成测试全部通过
) else (
    echo [ERROR] 集成测试失败
)

echo.

REM 生成覆盖率报告
if %GENERATE_COVERAGE%==1 (
    if exist coverage.out (
        echo [INFO] 生成覆盖率报告...
        go tool cover -html=coverage.out -o coverage.html
        go tool cover -func=coverage.out -o coverage.txt

        echo [SUCCESS] 覆盖率报告生成完成
        echo [INFO] HTML报告: %BACKEND_DIR%\coverage.html
        echo [INFO] 文本报告: %BACKEND_DIR%\coverage.txt
    ) else (
        echo [WARNING] 未找到覆盖率文件 coverage.out
    )
)

REM 清理测试资源
echo [INFO] 清理测试资源...
docker ps -a --filter "label=org.testcontainers" --format "{{{{.ID}}}}" > temp_containers.txt
for /f "tokens=*" %%i in (temp_containers.txt) do (
    if not "%%i"=="" (
        docker rm -f %%i >nul 2>&1
    )
)
del temp_containers.txt
echo [SUCCESS] 测试资源清理完成

echo.
echo [INFO] ==================== 集成测试结束 ====================

if %TEST_RESULT%==0 (
    exit /b 0
) else (
    exit /b %TEST_RESULT%
)

:run_suite
echo [INFO] 运行测试套件: %SUITE_NAME%

if "%SUITE_NAME%"=="tenant-permission" (
    %GO_TEST_CMD% ./tests/integration/ -run TestTenantPermissionIntegrationSuite
    goto test_done
)

if "%SUITE_NAME%"=="saga-business" (
    %GO_TEST_CMD% ./tests/integration/ -run TestSagaBusinessIntegrationSuite
    goto test_done
)

if "%SUITE_NAME%"=="humaninloop-bot" (
    %GO_TEST_CMD% ./tests/integration/ -run TestHumanInLoopBotIntegrationSuite
    goto test_done
)

if "%SUITE_NAME%"=="memory-conversation" (
    %GO_TEST_CMD% ./tests/integration/ -run TestMemoryConversationIntegrationSuite
    goto test_done
)

echo [ERROR] 未知的测试套件: %SUITE_NAME%
echo [INFO] 可用的测试套件:
echo   - tenant-permission
echo   - saga-business
echo   - humaninloop-bot
echo   - memory-conversation
exit /b 1

:test_done
set TEST_RESULT=%errorlevel%
goto generate_coverage

:print_usage
cat <<EOF
集成测试执行脚本 (Windows)

用法: %~nx0 [选项] [测试套件]

选项:
    -h, --help          显示此帮助信息
    -c, --cleanup-only  仅清理测试资源
    -v, --verbose       详细输出
    -s, --suite <name>  运行特定测试套件
    --no-coverage       跳过覆盖率报告生成

测试套件:
    tenant-permission      租户+权限系统集成测试
    saga-business          Saga+业务系统集成测试
    humaninloop-bot        人机协同+Bot系统集成测试
    memory-conversation    记忆+对话系统集成测试

示例:
    # 运行所有集成测试
    %~nx0

    # 运行特定测试套件
    %~nx0 -s tenant-permission

    # 详细输出模式
    %~nx0 -v

    # 仅清理测试资源
    %~nx0 -c

EOF
exit /b 0

:end
endlocal
