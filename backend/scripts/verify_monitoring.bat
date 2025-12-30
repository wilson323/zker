@echo off
REM 🔧 P0修复：监控系统验证脚本 (Windows版本)
REM 用于验证监控日志系统是否正确初始化和运行

setlocal enabledelayedexpansion

REM 配置
if "%API_BASE_URL%"=="" set API_BASE_URL=http://localhost:8888
if "%JAEGER_UI_URL%"=="" set JAEGER_UI_URL=http://localhost:16686
if "%PROMETHEUS_URL%"=="" set PROMETHEUS_URL=http://localhost:9090
if "%GRAFANA_URL%"=="" set GRAFANA_URL=http://localhost:3001

echo 🔍 开始验证监控系统...
echo.

REM 1. 检查Prometheus /metrics端点
echo 📊 检查 1: Prometheus /metrics 端点
curl -s "%API_BASE_URL%/metrics" > %TEMP%\metrics.txt 2>&1
findstr /C:"go_memstats" %TEMP%\metrics.txt >nul 2>&1
if %ERRORLEVEL%==0 (
    echo ✅ Prometheus /metrics 端点正常
    echo    找到指标示例:
    head -n 5 %TEMP%\metrics.txt
    echo.
) else (
    echo ❌ Prometheus /metrics 端点异常或无响应
    type %TEMP%\metrics.txt
    exit /b 1
)

REM 2. 检查/health端点
echo 💓 检查 2: 健康检查 /health 端点
curl -s "%API_BASE_URL%/health" > %TEMP%\health.txt 2>&1
findstr /C:"healthy" %TEMP%\health.txt >nul 2>&1
if %ERRORLEVEL%==0 (
    echo ✅ 健康检查端点正常
    echo    响应:
    type %TEMP%\health.txt
    echo.
) else (
    echo ⚠️  健康检查端点可能未配置（非必需）
    echo    响应:
    type %TEMP%\health.txt
    echo.
)

REM 3. 检查关键指标
echo 📈 检查 3: 关键业务指标
set METRICS_FOUND=0
for %%m in (http_requests_total http_request_duration_seconds tenant_total db_query_duration_seconds) do (
    findstr /C:"%%m" %TEMP%\metrics.txt >nul 2>&1
    if !ERRORLEVEL!==0 (
        echo ✅ 找到指标: %%m
        set /a METRICS_FOUND+=1
    ) else (
        echo ⚠️  未找到指标: %%m
    )
)
echo.

REM 4. 检查Jaeger
echo 🔍 检查 4: Jaeger分布式追踪
curl -s "%JAEGER_UI_URL%/api/status" > %TEMP%\jaeger.txt 2>&1
findstr /C:"jaeger" %TEMP%\jaeger.txt >nul 2>&1
if %ERRORLEVEL%==0 (
    echo ✅ Jaeger UI可访问: %JAEGER_UI_URL%
    echo.
) else (
    echo ⚠️  Jaeger UI可能未启动（非必需）
    echo    URL: %JAEGER_UI_URL%
    echo.
)

REM 5. 检查Prometheus
echo 🔍 检查 5: Prometheus服务
curl -s "%PROMETHEUS_URL%/-/healthy" > %TEMP%\prometheus.txt 2>&1
findstr /C:"Prometheus" %TEMP%\prometheus.txt >nul 2>&1
if %ERRORLEVEL%==0 (
    echo ✅ Prometheus服务正常
    echo    URL: %PROMETHEUS_URL%
    echo.
) else (
    echo ⚠️  Prometheus服务可能未启动（非必需）
    echo    URL: %PROMETHEUS_URL%
    echo.
)

REM 6. 检查Grafana
echo 🔍 检查 6: Grafana可视化服务
curl -s "%GRAFANA_URL%/api/health" > %TEMP%\grafana.txt 2>&1
findstr /C:"commit" %TEMP%\grafana.txt >nul 2>&1
if %ERRORLEVEL%==0 (
    echo ✅ Grafana服务正常
    echo    URL: %GRAFANA_URL%
    echo.
) else (
    echo ⚠️  Grafana服务可能未启动（非必需）
    echo    URL: %GRAFANA_URL%
    echo.
)

REM 7. 测试指标生成
echo 🧪 检查 7: 测试指标生成
echo    发送测试请求到 %API_BASE_URL%...
curl -s "%API_BASE_URL%" >nul 2>&1
timeout /t 2 /nobreak >nul

curl -s "%API_BASE_URL%/metrics" > %TEMP%\metrics_new.txt 2>&1
findstr /B /C:"http_requests_total" %TEMP%\metrics_new.txt > %TEMP%\requests_count.txt
set /a REQUEST_COUNT=0
for /f %%a in ('type %TEMP%\requests_count.txt ^| find /c /v ""') do set REQUEST_COUNT=%%a

if %REQUEST_COUNT% GTR 0 (
    echo ✅ 指标正在生成（找到 %REQUEST_COUNT% 个http_requests_total指标）
    echo.
) else (
    echo ⚠️  未检测到新的HTTP请求指标
    echo.
)

REM 总结
echo ==========================================
echo ✨ 监控系统验证完成！
echo.
echo 📋 访问地址：
echo    - API服务:     %API_BASE_URL%
echo    - Prometheus:  %PROMETHEUS_URL%
echo    - Grafana:     %GRAFANA_URL% (admin/admin)
echo    - Jaeger UI:   %JAEGER_UI_URL%
echo.
echo 💡 提示：
echo    1. 检查/metrics端点是否暴露所有预期指标
echo    2. 在Prometheus中查询: up{job=\"coze-api\"}
echo    3. 在Grafana中创建仪表盘可视化
echo    4. 在Jaeger中查看分布式追踪
echo ==========================================

REM 清理临时文件
del %TEMP%\metrics.txt %TEMP%\health.txt %TEMP%\jaeger.txt %TEMP%\prometheus.txt %TEMP%\grafana.txt %TEMP%\metrics_new.txt %TEMP%\requests_count.txt 2>nul

if %METRICS_FOUND% GEQ 3 (
    exit /b 0
) else (
    echo ⚠️  部分关键指标未找到，请检查配置
    exit /b 1
)
