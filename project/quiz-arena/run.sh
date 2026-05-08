#!/bin/bash
# quiz-arena 一键启动脚本
# 同时启动 Go 后端和 Vite 前端开发服务器

set -e
cd "$(dirname "$0")"

GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

cleanup() {
    echo -e "\n${YELLOW}🛑 正在停止服务...${NC}"
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
    exit 0
}
trap cleanup INT TERM

# 启动后端
echo -e "${CYAN}🚀 启动 Go 后端 (端口 8090)...${NC}"
cd backend
go run main.go &
BACKEND_PID=$!
cd ..

# 检查前端依赖
if [ ! -d "frontend/node_modules" ]; then
    echo -e "${YELLOW}📦 安装前端依赖...${NC}"
    cd frontend && npm install && cd ..
fi

# 启动前端
echo -e "${CYAN}🎨 启动 Vite 前端 (端口 5173)...${NC}"
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

sleep 2
echo ""
echo -e "${GREEN}✅ 服务启动完成！${NC}"
echo -e "   后端 API: ${CYAN}http://localhost:8090${NC}"
echo -e "   前端页面: ${CYAN}http://localhost:5173${NC}"
echo -e "${YELLOW}按 Ctrl+C 停止服务${NC}"

wait
