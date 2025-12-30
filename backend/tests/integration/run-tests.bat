@echo off
REM 组织中心集成测试运行脚本 (Windows)
REM 用途: 简化测试执行，支持多种运行模式

setlocal enabledelayedexpansion

REM 切换到脚本所在目录
cd /d "%~dp0"
cd ..\..\

REM 默认参数
set VERBOSE=
set COVER=
set RUN_TEST=
set PARALLEL=-parallel 4
set SHORT=
set CLEAN=

REM 解析参数
:parse_args
if "%~1"=="" goto end_parse
if "%~1"=="-h" goto show_help
if "%~1"=="--help" goto show_help
if "%~1"=="-v" (
    set VERBOSE=-v
    shift
    goto parse_args
)
if "%~1"=="-c" (
    set COVER=-cover
    shift
    goto parse_args
)
if "%~1"=="-r" (
    set RUN_TEST=-run %~2
    shift
    shift
    goto parse_args
)
if "%~1"=="-p" (
    set PARALLEL=-parallel %~2
    shift
    shift
    goto parse_args
)
if "%~1"=="-s" (
    set SHORT=-short
    shift
    goto parse_args
)
if "%~1"=="--clean" (
    set CLEAN=1
    shift
    goto parse_args
)
echo 未知选项: %~1
goto show_help

:end_parse

REM 显示配置
echo ================================================
echo 组织中心集成测试
echo ================================================
echo.
echo 测试配置:
echo   详细输出: %VERBOSE%
echo   覆盖率:   %COVER%
echo   并行数:   %PARALLEL%
echo   特定测试: %RUN_TEST%
echo   跳过测试: %SHORT%
echo.

REM 清理缓存
if defined CLEAN (
    echo 清理测试缓存...
    go clean -testcache
    echo ✓ 缓存已清理
    echo.
)

REM 检查Docker
echo 检查Docker环境...
docker --version >nul 2>&1
if errorlevel 1 (
    echo ✗ Docker未安装
    exit /b 1
)

docker ps >nul 2>&1
if errorlevel 1 (
    echo ✗ Docker未运行
    exit /b 1
)

echo ✓ Docker环境正常
echo.

REM 构建测试命令
set TEST_CMD=go test %VERBOSE% %COVER% %PARALLEL% %SHORT% -tags=integration -timeout 10m ./tests/integration/...

REM 添加特定测试
if defined RUN_TEST (
    set TEST_CMD=!TEST_CMD! %RUN_TEST%
)

REM 运行测试
echo 运行测试...
echo.
echo 命令: %TEST_CMD%
echo.

REM 执行测试
%TEST_CMD%
if errorlevel 1 (
    echo.
    echo ================================================
    echo ✗ 测试失败
    echo ================================================
    exit /b 1
) else (
    echo.
    echo ================================================
    echo ✓ 所有测试通过！
    echo ================================================
    exit /b 0
)

:show_help
echo 组织中心集成测试运行脚本
echo.
echo 用法: run-tests.bat [选项]
echo.
echo 选项:
echo   -h, --help              显示帮助信息
echo   -v, --verbose           详细输出模式
echo   -c, --cover             生成覆盖率报告
echo   -r, --run ^<name^>       运行特定测试
echo   -p, --parallel ^<num^>   并行运行测试（默认4）
echo   -s, --short             跳过集成测试
echo   --clean                 清理测试缓存
echo.
echo 示例:
echo   run-tests.bat                    # 运行所有测试
echo   run-tests.bat -v                 # 详细输出
echo   run-tests.bat -c                 # 生成覆盖率
echo   run-tests.bat -r OrgCRUD         # 运行特定测试
echo   run-tests.bat -p 8               # 8个并行
echo.
exit /b 0
