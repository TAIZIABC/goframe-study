// sqlgen/sqlgen.go
// 根据差异生成同步 SQL 语句
package sqlgen

import (
	"fmt"
	"strings"

	"dbdiff/diff"
	"dbdiff/schema"
)

// GenerateSQL 根据差异列表生成 SQL 语句
func GenerateSQL(diffs []diff.Difference, source *schema.Database) []string {
	var sqls []string

	for _, d := range diffs {
		switch d.Type {
		case diff.TableAdded:
			sqls = append(sqls, generateCreateTable(source.Tables[d.Table]))
		case diff.TableDropped:
			sqls = append(sqls, fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", d.Table))
		case diff.ColumnAdded:
			sqls = append(sqls, generateAddColumn(d, source))
		case diff.ColumnDropped:
			sqls = append(sqls, fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", d.Table, d.Name))
		case diff.ColumnChanged:
			sqls = append(sqls, generateModifyColumn(d, source))
		case diff.IndexAdded:
			sqls = append(sqls, generateAddIndex(d, source))
		case diff.IndexDropped:
			if d.Name == "PRIMARY" {
				sqls = append(sqls, fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY;", d.Table))
			} else {
				sqls = append(sqls, fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`;", d.Table, d.Name))
			}
		case diff.IndexChanged:
			// 先删后加
			if d.Name == "PRIMARY" {
				sqls = append(sqls, fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY;", d.Table))
			} else {
				sqls = append(sqls, fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`;", d.Table, d.Name))
			}
			sqls = append(sqls, generateAddIndex(d, source))
		}
	}

	return sqls
}

func generateCreateTable(t *schema.Table) string {
	if t == nil {
		return ""
	}

	var lines []string
	for _, c := range t.Columns {
		line := fmt.Sprintf("  `%s` %s", c.Name, c.Type)
		if c.Nullable == "NO" {
			line += " NOT NULL"
		}
		if c.Default.Valid {
			if c.Default.String == "CURRENT_TIMESTAMP" {
				line += " DEFAULT CURRENT_TIMESTAMP"
			} else {
				line += fmt.Sprintf(" DEFAULT '%s'", c.Default.String)
			}
		}
		if c.Extra != "" {
			line += " " + strings.ToUpper(c.Extra)
		}
		if c.Comment != "" {
			line += fmt.Sprintf(" COMMENT '%s'", c.Comment)
		}
		lines = append(lines, line)
	}

	// 索引
	for _, idx := range t.Indexes {
		cols := backtickJoin(idx.Columns)
		if idx.Name == "PRIMARY" {
			lines = append(lines, fmt.Sprintf("  PRIMARY KEY (%s)", cols))
		} else if idx.Unique {
			lines = append(lines, fmt.Sprintf("  UNIQUE KEY `%s` (%s)", idx.Name, cols))
		} else {
			lines = append(lines, fmt.Sprintf("  KEY `%s` (%s)", idx.Name, cols))
		}
	}

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n%s\n)", t.Name, strings.Join(lines, ",\n"))
	if t.Engine != "" {
		sql += " ENGINE=" + t.Engine
	}
	sql += ";"
	return sql
}

func generateAddColumn(d diff.Difference, source *schema.Database) string {
	t := source.Tables[d.Table]
	if t == nil {
		return ""
	}

	// 找到该列的完整定义和位置
	for i, c := range t.Columns {
		if c.Name == d.Name {
			def := columnDefinition(&c)
			after := ""
			if i > 0 {
				after = fmt.Sprintf(" AFTER `%s`", t.Columns[i-1].Name)
			} else {
				after = " FIRST"
			}
			return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s%s;", d.Table, def, after)
		}
	}
	return ""
}

func generateModifyColumn(d diff.Difference, source *schema.Database) string {
	t := source.Tables[d.Table]
	if t == nil {
		return ""
	}

	for _, c := range t.Columns {
		if c.Name == d.Name {
			def := columnDefinition(&c)
			return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s;", d.Table, def)
		}
	}
	return ""
}

func generateAddIndex(d diff.Difference, source *schema.Database) string {
	t := source.Tables[d.Table]
	if t == nil {
		return ""
	}

	for _, idx := range t.Indexes {
		if idx.Name == d.Name {
			cols := backtickJoin(idx.Columns)
			if idx.Name == "PRIMARY" {
				return fmt.Sprintf("ALTER TABLE `%s` ADD PRIMARY KEY (%s);", d.Table, cols)
			} else if idx.Unique {
				return fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE INDEX `%s` (%s);", d.Table, idx.Name, cols)
			} else {
				return fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `%s` (%s);", d.Table, idx.Name, cols)
			}
		}
	}
	return ""
}

func columnDefinition(c *schema.Column) string {
	def := fmt.Sprintf("`%s` %s", c.Name, c.Type)
	if c.Nullable == "NO" {
		def += " NOT NULL"
	}
	if c.Default.Valid {
		if c.Default.String == "CURRENT_TIMESTAMP" {
			def += " DEFAULT CURRENT_TIMESTAMP"
		} else {
			def += fmt.Sprintf(" DEFAULT '%s'", c.Default.String)
		}
	}
	if c.Extra != "" {
		def += " " + strings.ToUpper(c.Extra)
	}
	if c.Comment != "" {
		def += fmt.Sprintf(" COMMENT '%s'", c.Comment)
	}
	return def
}

func backtickJoin(cols []string) string {
	var parts []string
	for _, c := range cols {
		parts = append(parts, "`"+c+"`")
	}
	return strings.Join(parts, ", ")
}
