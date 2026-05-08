#!/bin/bash
# Go 并发练习项目 - 运行脚本
# 用法: ./run.sh [示例编号|test|bench|race|all]

set -e
cd "$(dirname "$0")"

GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_header() {
    echo ""
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}========================================${NC}"
}

run_example() {
    local dir=$1
    local name=$2
    print_header "$name"
    go run ./$dir/
}

case "${1:-all}" in
    1|goroutine)
        run_example "01_goroutine" "Goroutine 基础"
        ;;
    2|channel)
        run_example "02_channel" "Channel 通信"
        ;;
    3|select)
        run_example "03_select" "Select 语句"
        ;;
    4|sync)
        run_example "04_sync" "Sync 并发原语"
        ;;
    5|context)
        run_example "05_context" "Context 控制"
        ;;
    6|patterns)
        run_example "06_patterns" "并发模式"
        ;;
    test)
        print_header "运行单元测试"
        go test -v -count=1 .
        ;;
    race)
        print_header "竞态条件检测"
        go test -race -v -count=1 .
        ;;
    bench)
        print_header "性能基准测试"
        go test -bench=. -benchmem -count=1 .
        ;;
    all)
        run_example "01_goroutine" "01 - Goroutine 基础"
        run_example "02_channel" "02 - Channel 通信"
        run_example "03_select" "03 - Select 语句"
        run_example "04_sync" "04 - Sync 并发原语"
        run_example "05_context" "05 - Context 控制"
        run_example "06_patterns" "06 - 并发模式"

        echo ""
        print_header "运行测试 + 竞态检测"
        go test -race -v -count=1 .

        print_header "性能基准测试"
        go test -bench=. -benchmem -count=1 .

        echo ""
        echo -e "${GREEN}✅ 全部运行完毕！${NC}"
        ;;
    *)
        echo "用法: $0 [选项]"
        echo ""
        echo "选项:"
        echo "  1 | goroutine   运行 Goroutine 示例"
        echo "  2 | channel     运行 Channel 示例"
        echo "  3 | select      运行 Select 示例"
        echo "  4 | sync        运行 Sync 示例"
        echo "  5 | context     运行 Context 示例"
        echo "  6 | patterns    运行并发模式示例"
        echo "  test            运行单元测试"
        echo "  race            运行竞态检测"
        echo "  bench           运行基准测试"
        echo "  all             运行全部（默认）"
        ;;
esac
