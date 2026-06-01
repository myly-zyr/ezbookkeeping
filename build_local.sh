#!/bin/bash
# ============================================================
# ezbookkeeping 构建与部署脚本
# ============================================================
# 用法:
#   ./build.sh              # 完整构建（Go + 前端 + 部署）
#   ./build.sh backend      # 仅编译Go后端
#   ./build.sh frontend     # 仅编译前端
#   ./build.sh deploy       # 仅部署到容器
#   ./build.sh status       # 查看当前状态
# ============================================================

set -e

# ======================== 配置 ========================
# 项目根目录
PROJECT_DIR="/opt/ezbookkeeping"

# Go 配置
GO_BIN="/usr/local/go/bin/go"
GOPROXY="https://goproxy.cn,direct"
GO_BUILD_FLAGS="-trimpath -ldflags '-w -s -linkmode external -extldflags \"-static\"'"

# 前端配置
NODE_BUILD_CMD="npm run build"

# Docker 部署配置
CONTAINER_NAME="ezbookkeeping-edge-fix"
CONTAINER_BINARY_PATH="/ezbookkeeping/ezbookkeeping"
CONTAINER_PUBLIC_PATH="/ezbookkeeping/public/"

# 静态文件源目录（前端构建输出）
DIST_DIR="${PROJECT_DIR}/dist"

# ======================== 函数 ========================

log_info() {
    echo -e "\033[32m[INFO]\033[0m $1"
}

log_error() {
    echo -e "\033[31m[ERROR]\033[0m $1"
}

log_step() {
    echo -e "\033[34m[STEP]\033[0m $1"
}

# 检查环境
check_env() {
    log_step "检查构建环境..."

    # Go
    if [ ! -f "$GO_BIN" ]; then
        log_error "Go 未安装: $GO_BIN"
        exit 1
    fi
    GO_VERSION=$($GO_BIN version)
    log_info "Go: $GO_VERSION"

    # Node
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装"
        exit 1
    fi
    NODE_VERSION=$(node --version)
    log_info "Node.js: $NODE_VERSION"

    # GCC (CGO需要)
    if ! command -v gcc &> /dev/null; then
        log_error "GCC 未安装 (CGO需要)"
        exit 1
    fi
    GCC_VERSION=$(gcc --version | head -1)
    log_info "GCC: $GCC_VERSION"

    # Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装"
        exit 1
    fi
    log_info "Docker: $(docker --version)"

    # 容器
    if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_error "容器 ${CONTAINER_NAME} 未运行"
        exit 1
    fi
    log_info "容器: ${CONTAINER_NAME} 运行中"

    log_info "环境检查通过 ✓"
    echo ""
}

# 编译Go后端
build_backend() {
    log_step "编译Go后端..."

    cd "$PROJECT_DIR"

    # 设置Go代理
    export GOPROXY="$GOPROXY"

    # 静态编译（宿主机Ubuntu glibc -> 容器Alpine musl）
    log_info "静态编译中（CGO_ENABLED=1, 静态链接）..."
    CGO_ENABLED=1 CC=gcc $GO_BIN build \
        -trimpath \
        -ldflags '-w -s -linkmode external -extldflags "-static"' \
        -o "${PROJECT_DIR}/ezbookkeeping_new" \
        ezbookkeeping.go

    if [ ! -f "${PROJECT_DIR}/ezbookkeeping_new" ]; then
        log_error "Go编译失败"
        exit 1
    fi

    # 替换旧二进制
    mv "${PROJECT_DIR}/ezbookkeeping_new" "${PROJECT_DIR}/ezbookkeeping"
    chmod +x "${PROJECT_DIR}/ezbookkeeping"

    BINARY_SIZE=$(ls -lh "${PROJECT_DIR}/ezbookkeeping" | awk '{print $5}')
    BINARY_TYPE=$(file "${PROJECT_DIR}/ezbookkeeping" | grep -o 'statically linked\|dynamically linked')
    log_info "Go编译成功 ✓ (${BINARY_SIZE}, ${BINARY_TYPE})"
    echo ""
}

# 编译前端
build_frontend() {
    log_step "编译前端..."

    cd "$PROJECT_DIR"

    # 安装依赖（如果node_modules不存在）
    if [ ! -d "node_modules" ]; then
        log_info "安装前端依赖..."
        npm install
    fi

    # 编译
    log_info "编译前端..."
    $NODE_BUILD_CMD

    if [ ! -d "$DIST_DIR" ]; then
        log_error "前端编译失败: dist目录不存在"
        exit 1
    fi

    DIST_SIZE=$(du -sh "$DIST_DIR" | awk '{print $1}')
    log_info "前端编译成功 ✓ (${DIST_SIZE})"
    echo ""
}

# 部署到容器
deploy() {
    log_step "部署到容器..."

    # 检查文件是否存在
    if [ ! -f "${PROJECT_DIR}/ezbookkeeping" ]; then
        log_error "二进制文件不存在，请先编译后端"
        exit 1
    fi

    if [ ! -d "$DIST_DIR" ]; then
        log_error "前端dist目录不存在，请先编译前端"
        exit 1
    fi

    # 复制二进制
    log_info "复制后端二进制..."
    docker cp "${PROJECT_DIR}/ezbookkeeping" "${CONTAINER_NAME}:${CONTAINER_BINARY_PATH}"

    # 复制前端静态文件
    log_info "复制前端静态文件..."
    docker cp "${DIST_DIR}/." "${CONTAINER_NAME}:${CONTAINER_PUBLIC_PATH}"

    # 重启容器
    log_info "重启容器..."
    docker restart "$CONTAINER_NAME"

    # 等待健康检查
    log_info "等待容器启动..."
    sleep 8

    STATUS=$(docker ps --format '{{.Status}}' --filter "name=${CONTAINER_NAME}")
    if echo "$STATUS" | grep -q "healthy"; then
        log_info "部署成功 ✓ 容器状态: ${STATUS}"
    else
        log_error "容器可能未正常启动: ${STATUS}"
        log_error "查看日志: docker logs ${CONTAINER_NAME}"
    fi
    echo ""
}

# 查看状态
show_status() {
    echo "=============================="
    echo "  ezbookkeeping 状态"
    echo "=============================="
    echo ""

    # Go
    echo "Go:     $($GO_BIN version 2>/dev/null || echo '未安装')"
    echo "Node:   $(node --version 2>/dev/null || echo '未安装')"
    echo "Docker: $(docker --version 2>/dev/null || echo '未安装')"
    echo ""

    # 二进制
    if [ -f "${PROJECT_DIR}/ezbookkeeping" ]; then
        BINARY_TIME=$(ls -la "${PROJECT_DIR}/ezbookkeeping" | awk '{print $6, $7, $8}')
        BINARY_SIZE=$(ls -lh "${PROJECT_DIR}/ezbookkeeping" | awk '{print $5}')
        BINARY_TYPE=$(file "${PROJECT_DIR}/ezbookkeeping" | grep -o 'statically linked\|dynamically linked')
        echo "后端二进制: ${BINARY_SIZE} (${BINARY_TYPE}) 编译于 ${BINARY_TIME}"
    else
        echo "后端二进制: 不存在"
    fi

    # 前端
    if [ -d "$DIST_DIR" ]; then
        DIST_TIME=$(ls -la "$DIST_DIR" | head -2 | tail -1 | awk '{print $6, $7, $8}')
        DIST_SIZE=$(du -sh "$DIST_DIR" | awk '{print $1}')
        echo "前端构建: ${DIST_SIZE} 构建于 ${DIST_TIME}"
    else
        echo "前端构建: 不存在"
    fi

    # 容器
    echo ""
    echo "容器状态:"
    docker ps --format '  {{.Names}}  {{.Status}}  {{.Ports}}' --filter "name=${CONTAINER_NAME}" 2>/dev/null || echo "  容器未运行"
    echo ""

    # 数据库
    if [ -f "${PROJECT_DIR}/db/ezbookkeeping.db" ]; then
        DB_SIZE=$(ls -lh "${PROJECT_DIR}/db/ezbookkeeping.db" | awk '{print $5}')
        echo "数据库: ${DB_SIZE}"
    else
        echo "数据库: 不存在"
    fi
    echo ""
}

# ======================== 主流程 ========================

cd "$PROJECT_DIR"

case "${1:-all}" in
    backend)
        check_env
        build_backend
        ;;
    frontend)
        build_frontend
        ;;
    deploy)
        deploy
        ;;
    status)
        show_status
        ;;
    all)
        check_env
        build_backend
        build_frontend
        deploy
        show_status
        ;;
    *)
        echo "用法: $0 [backend|frontend|deploy|status|all]"
        exit 1
        ;;
esac
