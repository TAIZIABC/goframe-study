// diff/diff.go
// 结构对比引擎：检测表级、字段级、索引级差异
package diff

import (
	"dbdiff/schema"
)

// DiffType 差异类型
type DiffType int

const (
	TableAdded   DiffType = iota // 源有目标无 → 需要 CREATE
	TableDropped                 // 源无目标有 → 需要 DROP
	ColumnAdded                  // 新增字段
	ColumnDropped                // 删除字段
	ColumnChanged                // 字段变更
	IndexAdded                   // 新增索引
	IndexDropped                 // 删除索引
	IndexChanged                 // 索引变更
)

func (d DiffType) String() string {
	switch d {
	case TableAdded:
		return "新增表"
	case TableDropped:
		return "删除表"
	case ColumnAdded:
		return "新增字段"
	case ColumnDropped:
		return "删除字段"
	case ColumnChanged:
		return "字段变更"
	case IndexAdded:
		return "新增索引"
	case IndexDropped:
		return "删除索引"
	case IndexChanged:
		return "索引变更"
	default:
		return "未知"
	}
}

// Difference 一个具体的差异项
type Difference struct {
	Type      DiffType
	Table     string
	Name      string // 字段名或索引名
	SourceDef string // 源端定义
	TargetDef string // 目标端定义
}

// Compare 对比两个数据库的结构差异
// source = 源库（期望的结构），target = 目标库（待同步的）
func Compare(source, target *schema.Database) []Difference {
	var diffs []Difference

	srcNames := make(map[string]bool)
	tgtNames := make(map[string]bool)

	for _, n := range source.SortedTableNames() {
		srcNames[n] = true
	}
	for _, n := range target.SortedTableNames() {
		tgtNames[n] = true
	}

	// 1. 源有目标无 → 新增表
	for _, name := range source.SortedTableNames() {
		if !tgtNames[name] {
			diffs = append(diffs, Difference{
				Type:  TableAdded,
				Table: name,
			})
		}
	}

	// 2. 源无目标有 → 删除表
	for _, name := range target.SortedTableNames() {
		if !srcNames[name] {
			diffs = append(diffs, Difference{
				Type:  TableDropped,
				Table: name,
			})
		}
	}

	// 3. 两端都有 → 对比字段和索引
	for _, name := range source.SortedTableNames() {
		if !tgtNames[name] {
			continue
		}
		srcTable := source.Tables[name]
		tgtTable := target.Tables[name]

		colDiffs := compareColumns(name, srcTable.Columns, tgtTable.Columns)
		diffs = append(diffs, colDiffs...)

		idxDiffs := compareIndexes(name, srcTable.Indexes, tgtTable.Indexes)
		diffs = append(diffs, idxDiffs...)
	}

	return diffs
}

func compareColumns(table string, srcCols, tgtCols []schema.Column) []Difference {
	var diffs []Difference

	srcMap := make(map[string]*schema.Column)
	tgtMap := make(map[string]*schema.Column)

	for i := range srcCols {
		srcMap[srcCols[i].Name] = &srcCols[i]
	}
	for i := range tgtCols {
		tgtMap[tgtCols[i].Name] = &tgtCols[i]
	}

	// 源有目标无 → 新增字段
	for _, c := range srcCols {
		if _, ok := tgtMap[c.Name]; !ok {
			diffs = append(diffs, Difference{
				Type:      ColumnAdded,
				Table:     table,
				Name:      c.Name,
				SourceDef: c.Signature(),
			})
		}
	}

	// 源无目标有 → 删除字段
	for _, c := range tgtCols {
		if _, ok := srcMap[c.Name]; !ok {
			diffs = append(diffs, Difference{
				Type:      ColumnDropped,
				Table:     table,
				Name:      c.Name,
				TargetDef: c.Signature(),
			})
		}
	}

	// 两端都有 → 对比签名
	for _, sc := range srcCols {
		tc, ok := tgtMap[sc.Name]
		if !ok {
			continue
		}
		srcSig := sc.Signature()
		tgtSig := tc.Signature()
		if srcSig != tgtSig {
			diffs = append(diffs, Difference{
				Type:      ColumnChanged,
				Table:     table,
				Name:      sc.Name,
				SourceDef: srcSig,
				TargetDef: tgtSig,
			})
		}
	}

	return diffs
}

func compareIndexes(table string, srcIdx, tgtIdx []schema.Index) []Difference {
	var diffs []Difference

	srcMap := make(map[string]*schema.Index)
	tgtMap := make(map[string]*schema.Index)

	for i := range srcIdx {
		srcMap[srcIdx[i].Name] = &srcIdx[i]
	}
	for i := range tgtIdx {
		tgtMap[tgtIdx[i].Name] = &tgtIdx[i]
	}

	// 源有目标无
	for _, idx := range srcIdx {
		if _, ok := tgtMap[idx.Name]; !ok {
			diffs = append(diffs, Difference{
				Type:      IndexAdded,
				Table:     table,
				Name:      idx.Name,
				SourceDef: idx.Signature(),
			})
		}
	}

	// 源无目标有
	for _, idx := range tgtIdx {
		if _, ok := srcMap[idx.Name]; !ok {
			diffs = append(diffs, Difference{
				Type:      IndexDropped,
				Table:     table,
				Name:      idx.Name,
				TargetDef: idx.Signature(),
			})
		}
	}

	// 两端都有 → 对比签名
	for _, si := range srcIdx {
		ti, ok := tgtMap[si.Name]
		if !ok {
			continue
		}
		srcSig := si.Signature()
		tgtSig := ti.Signature()
		if srcSig != tgtSig {
			diffs = append(diffs, Difference{
				Type:      IndexChanged,
				Table:     si.Name,
				Name:      si.Name,
				SourceDef: srcSig,
				TargetDef: tgtSig,
			})
		}
	}

	return diffs
}
