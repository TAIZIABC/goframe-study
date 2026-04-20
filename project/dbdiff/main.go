package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"dbdiff/diff"
	"dbdiff/schema"
	"dbdiff/sqlgen"
)

func main() {
	source := flag.String("source", "", "源库 DSN (期望的结构)，如: user:pass@tcp(host:3306)/db_dev")
	target := flag.String("target", "", "目标库 DSN (待同步)，如: user:pass@tcp(host:3306)/db_test")
	output := flag.String("out", "", "将 SQL 输出到文件 (不指定则打印到终端)")
	noDrop := flag.Bool("no-drop", false, "不生成 DROP TABLE 语句")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: dbdiff -source <DSN> -target <DSN> [选项]\n\n")
		fmt.Fprintf(os.Stderr, "对比两个 MySQL 数据库的表结构差异，生成同步 SQL。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  dbdiff -source 'root:123@tcp(localhost:3306)/dev_db' \\\n")
		fmt.Fprintf(os.Stderr, "         -target 'root:123@tcp(localhost:3306)/test_db'\n\n")
		fmt.Fprintf(os.Stderr, "  dbdiff -source '...' -target '...' -out sync.sql -no-drop\n")
	}
	flag.Parse()

	if *source == "" || *target == "" {
		fmt.Fprintln(os.Stderr, "✗ 请指定 -source 和 -target DSN")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println()

	// 加载源库
	fmt.Println("  🔍 连接源库...")
	srcDB, err := schema.LoadDatabase(*source, "源库")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ 源库 [%s]: %d 张表\n", srcDB.Name, len(srcDB.Tables))

	// 加载目标库
	fmt.Println("  🔍 连接目标库...")
	tgtDB, err := schema.LoadDatabase(*target, "目标库")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ 目标库 [%s]: %d 张表\n", tgtDB.Name, len(tgtDB.Tables))

	// 对比差异
	fmt.Println("\n  🔄 对比结构差异...")
	diffs := diff.Compare(srcDB, tgtDB)

	// 过滤 DROP
	if *noDrop {
		var filtered []diff.Difference
		for _, d := range diffs {
			if d.Type != diff.TableDropped {
				filtered = append(filtered, d)
			}
		}
		diffs = filtered
	}

	if len(diffs) == 0 {
		fmt.Println("\n  ✅ 两个数据库结构完全一致，无差异！")
		return
	}

	// 打印差异摘要
	fmt.Printf("\n  发现 %d 项差异:\n", len(diffs))
	fmt.Println("  " + strings.Repeat("─", 65))

	// 按类型统计
	typeCounts := make(map[diff.DiffType]int)
	for _, d := range diffs {
		typeCounts[d.Type]++
	}
	for _, dt := range []diff.DiffType{
		diff.TableAdded, diff.TableDropped,
		diff.ColumnAdded, diff.ColumnDropped, diff.ColumnChanged,
		diff.IndexAdded, diff.IndexDropped, diff.IndexChanged,
	} {
		if c, ok := typeCounts[dt]; ok {
			fmt.Printf("    %s: %d\n", dt, c)
		}
	}

	fmt.Println("  " + strings.Repeat("─", 65))

	// 打印详细差异
	for _, d := range diffs {
		switch d.Type {
		case diff.TableAdded:
			fmt.Printf("  \033[32m+ 新增表\033[0m  %s\n", d.Table)
		case diff.TableDropped:
			fmt.Printf("  \033[31m- 删除表\033[0m  %s\n", d.Table)
		case diff.ColumnAdded:
			fmt.Printf("  \033[32m+ 新增字段\033[0m  %s.%s  →  %s\n", d.Table, d.Name, d.SourceDef)
		case diff.ColumnDropped:
			fmt.Printf("  \033[31m- 删除字段\033[0m  %s.%s\n", d.Table, d.Name)
		case diff.ColumnChanged:
			fmt.Printf("  \033[33m~ 字段变更\033[0m  %s.%s\n", d.Table, d.Name)
			fmt.Printf("      源: %s\n", d.SourceDef)
			fmt.Printf("      目: %s\n", d.TargetDef)
		case diff.IndexAdded:
			fmt.Printf("  \033[32m+ 新增索引\033[0m  %s.%s  →  %s\n", d.Table, d.Name, d.SourceDef)
		case diff.IndexDropped:
			fmt.Printf("  \033[31m- 删除索引\033[0m  %s.%s\n", d.Table, d.Name)
		case diff.IndexChanged:
			fmt.Printf("  \033[33m~ 索引变更\033[0m  %s.%s\n", d.Table, d.Name)
			fmt.Printf("      源: %s\n", d.SourceDef)
			fmt.Printf("      目: %s\n", d.TargetDef)
		}
	}

	// 生成 SQL
	sqls := sqlgen.GenerateSQL(diffs, srcDB)

	fmt.Println("\n  " + strings.Repeat("─", 65))
	fmt.Printf("  生成 %d 条 SQL 语句:\n\n", len(sqls))

	sqlText := "-- =============================================\n"
	sqlText += fmt.Sprintf("-- 数据库同步脚本: %s → %s\n", srcDB.Name, tgtDB.Name)
	sqlText += "-- =============================================\n\n"

	for _, s := range sqls {
		sqlText += s + "\n\n"
		fmt.Printf("  %s\n\n", s)
	}

	// 输出到文件
	if *output != "" {
		if err := os.WriteFile(*output, []byte(sqlText), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ 写入文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("  📄 SQL 已保存到: %s\n\n", *output)
	}
}
