// schema/schema.go
// MySQL 表结构读取模块
package schema

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Column 列信息
type Column struct {
	Name         string
	Type         string // 如 varchar(200), int unsigned
	Nullable     string // YES / NO
	Default      sql.NullString
	Extra        string // auto_increment 等
	Comment      string
	CharacterSet sql.NullString
	Collation    sql.NullString
	ColumnKey    string // PRI, UNI, MUL
}

// Index 索引信息
type Index struct {
	Name      string
	Columns   []string
	Unique    bool
	IndexType string // BTREE, FULLTEXT, HASH
}

// Table 表结构
type Table struct {
	Name    string
	Columns []Column
	Indexes []Index
	Engine  string
	Charset string
	Comment string
}

// Database 数据库结构
type Database struct {
	DSN    string
	Name   string
	Tables map[string]*Table
}

// LoadDatabase 连接数据库并读取所有表结构
func LoadDatabase(dsn, label string) (*Database, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接 %s 失败: %w", label, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("连接 %s 失败: %w", label, err)
	}

	// 获取数据库名
	var dbName string
	if err := db.QueryRow("SELECT DATABASE()").Scan(&dbName); err != nil {
		return nil, fmt.Errorf("获取数据库名失败: %w", err)
	}

	database := &Database{
		DSN:    dsn,
		Name:   dbName,
		Tables: make(map[string]*Table),
	}

	// 获取所有表
	tables, err := loadTables(db, dbName)
	if err != nil {
		return nil, err
	}

	for _, t := range tables {
		cols, err := loadColumns(db, dbName, t.Name)
		if err != nil {
			return nil, err
		}
		t.Columns = cols

		indexes, err := loadIndexes(db, dbName, t.Name)
		if err != nil {
			return nil, err
		}
		t.Indexes = indexes

		database.Tables[t.Name] = t
	}

	return database, nil
}

func loadTables(db *sql.DB, dbName string) ([]*Table, error) {
	rows, err := db.Query(`
		SELECT TABLE_NAME, IFNULL(ENGINE,''), IFNULL(TABLE_COLLATION,''), IFNULL(TABLE_COMMENT,'')
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_NAME`, dbName)
	if err != nil {
		return nil, fmt.Errorf("查询表列表失败: %w", err)
	}
	defer rows.Close()

	var tables []*Table
	for rows.Next() {
		t := &Table{}
		if err := rows.Scan(&t.Name, &t.Engine, &t.Charset, &t.Comment); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

func loadColumns(db *sql.DB, dbName, tableName string) ([]Column, error) {
	rows, err := db.Query(`
		SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT, EXTRA,
		       IFNULL(COLUMN_COMMENT,''), CHARACTER_SET_NAME, COLLATION_NAME, IFNULL(COLUMN_KEY,'')
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`, dbName, tableName)
	if err != nil {
		return nil, fmt.Errorf("查询列信息失败 (%s): %w", tableName, err)
	}
	defer rows.Close()

	var cols []Column
	for rows.Next() {
		c := Column{}
		if err := rows.Scan(&c.Name, &c.Type, &c.Nullable, &c.Default,
			&c.Extra, &c.Comment, &c.CharacterSet, &c.Collation, &c.ColumnKey); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	return cols, nil
}

func loadIndexes(db *sql.DB, dbName, tableName string) ([]Index, error) {
	rows, err := db.Query(`
		SELECT INDEX_NAME, COLUMN_NAME, NON_UNIQUE, IFNULL(INDEX_TYPE,'BTREE')
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`, dbName, tableName)
	if err != nil {
		return nil, fmt.Errorf("查询索引信息失败 (%s): %w", tableName, err)
	}
	defer rows.Close()

	indexMap := make(map[string]*Index)
	var indexOrder []string

	for rows.Next() {
		var name, col, idxType string
		var nonUnique int
		if err := rows.Scan(&name, &col, &nonUnique, &idxType); err != nil {
			return nil, err
		}
		if _, ok := indexMap[name]; !ok {
			indexMap[name] = &Index{
				Name:      name,
				Unique:    nonUnique == 0,
				IndexType: idxType,
			}
			indexOrder = append(indexOrder, name)
		}
		indexMap[name].Columns = append(indexMap[name].Columns, col)
	}

	var indexes []Index
	for _, name := range indexOrder {
		indexes = append(indexes, *indexMap[name])
	}
	return indexes, nil
}

// ColumnSignature 返回列的完整签名，用于对比
func (c *Column) Signature() string {
	parts := []string{c.Type}
	if c.Nullable == "NO" {
		parts = append(parts, "NOT NULL")
	} else {
		parts = append(parts, "NULL")
	}
	if c.Default.Valid {
		parts = append(parts, "DEFAULT '"+c.Default.String+"'")
	}
	if c.Extra != "" {
		parts = append(parts, c.Extra)
	}
	if c.Comment != "" {
		parts = append(parts, "COMMENT '"+c.Comment+"'")
	}
	return strings.Join(parts, " ")
}

// IndexSignature 返回索引签名
func (idx *Index) Signature() string {
	u := ""
	if idx.Unique && idx.Name != "PRIMARY" {
		u = "UNIQUE "
	}
	return fmt.Sprintf("%sINDEX(%s) %s", u, strings.Join(idx.Columns, ","), idx.IndexType)
}

// SortedTableNames 返回排序的表名列表
func (d *Database) SortedTableNames() []string {
	var names []string
	for n := range d.Tables {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
