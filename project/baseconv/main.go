package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║        进制转换工具 (baseconv)       ║")
	fmt.Println("╠══════════════════════════════════════╣")
	fmt.Println("║  用法: 输入 <数字> <进制(2/8/10/16)> ║")
	fmt.Println("║  示例: FF 16                         ║")
	fmt.Println("║        255 10                        ║")
	fmt.Println("║        11111111 2                    ║")
	fmt.Println("║        377 8                         ║")
	fmt.Println("║  输入 q 退出                         ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "q" || line == "quit" || line == "exit" {
			fmt.Println("再见！")
			break
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Println("  ✗ 格式错误，请输入: <数字> <进制(2/8/10/16)>")
			continue
		}

		numStr := strings.ToUpper(parts[0])
		base, err := strconv.Atoi(parts[1])
		if err != nil || (base != 2 && base != 8 && base != 10 && base != 16) {
			fmt.Println("  ✗ 进制只支持 2、8、10、16")
			continue
		}

		// 支持负数：提取符号
		negative := false
		raw := numStr
		if strings.HasPrefix(raw, "-") {
			negative = true
			raw = raw[1:]
		}

		// 去掉常见前缀
		raw = strings.TrimPrefix(raw, "0X")
		raw = strings.TrimPrefix(raw, "0O")
		raw = strings.TrimPrefix(raw, "0B")

		value, err := strconv.ParseUint(raw, base, 64)
		if err != nil {
			fmt.Printf("  ✗ \"%s\" 不是有效的 %d 进制数\n", parts[0], base)
			continue
		}

		sign := ""
		if negative {
			sign = "-"
		}

		fmt.Println("  ┌──────────────────────────────────")
		fmt.Printf("  │  输入: %s (base %d)\n", parts[0], base)
		fmt.Println("  ├──────────────────────────────────")
		fmt.Printf("  │  二进制 (BIN):  %s%s\n", sign, formatBin(value))
		fmt.Printf("  │  八进制 (OCT):  %s%s\n", sign, strconv.FormatUint(value, 8))
		fmt.Printf("  │  十进制 (DEC):  %s%d\n", sign, value)
		fmt.Printf("  │  十六进制(HEX):  %s%s\n", sign, strings.ToUpper(strconv.FormatUint(value, 16)))
		fmt.Println("  └──────────────────────────────────")
	}
}

// formatBin 格式化二进制输出，每 4 位加空格便于阅读
func formatBin(v uint64) string {
	s := strconv.FormatUint(v, 2)
	// 补齐到 4 的倍数
	if r := len(s) % 4; r != 0 {
		s = strings.Repeat("0", 4-r) + s
	}
	var parts []string
	for i := 0; i < len(s); i += 4 {
		parts = append(parts, s[i:i+4])
	}
	return strings.Join(parts, " ")
}
