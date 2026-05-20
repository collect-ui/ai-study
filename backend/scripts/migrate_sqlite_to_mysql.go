package main

import (
	"bufio"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

type columnInfo struct {
	Name    string
	Type    string
	NotNull bool
	PKOrder int
}

type indexInfo struct {
	Name    string
	Unique  bool
	Origin  string
	Partial bool
	Columns []string
}

type logger struct {
	path  string
	lines []string
}

func (l *logger) add(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	fmt.Println(line)
	l.lines = append(l.lines, line)
	_ = os.WriteFile(l.path, []byte(strings.Join(l.lines, "\n")+"\n"), 0644)
}

func main() {
	sourceSQLite := flag.String("sqlite", "database/ai_study_admin.db", "source SQLite database")
	autoCheckConf := flag.String("auto-check-conf", "/data/project/auto-check/conf/application.properties", "auto-check application.properties")
	targetDB := flag.String("target-db", "ai_study", "target MySQL database")
	reportDir := flag.String("report-dir", "../reports/mysql_migration", "report directory")
	flag.Parse()

	if err := os.MkdirAll(*reportDir, 0755); err != nil {
		fatalf("create report dir: %v", err)
	}
	stamp := time.Now().Format("20060102-150405")
	logPath := filepath.Join(*reportDir, "migration_"+stamp+".md")
	log := &logger{path: logPath}

	log.add("# AI Study SQLite -> MySQL migration")
	log.add("")
	log.add("- started_at: %s", time.Now().Format(time.RFC3339))
	log.add("- source_sqlite: %s", *sourceSQLite)
	log.add("- target_mysql_db: %s", *targetDB)
	log.add("- auto_check_conf: %s", *autoCheckConf)

	props, err := readProperties(*autoCheckConf)
	if err != nil {
		fatalf("read auto-check config: %v", err)
	}
	sourceDSN := props["dataSourceName"]
	if sourceDSN == "" {
		fatalf("dataSourceName not found in %s", *autoCheckConf)
	}
	cfg, err := mysql.ParseDSN(sourceDSN)
	if err != nil {
		fatalf("parse auto-check MySQL DSN: %v", err)
	}
	configureMySQLConfig(cfg, "")
	log.add("- mysql_source: %s", sanitizeDSN(cfg))

	sqliteFile, err := os.Open(*sourceSQLite)
	if err != nil {
		fatalf("open sqlite file: %v", err)
	}
	sqliteHash, err := sha256File(sqliteFile)
	_ = sqliteFile.Close()
	if err != nil {
		fatalf("hash sqlite file: %v", err)
	}
	sqliteBackup := filepath.Join(*reportDir, "ai_study_admin.sqlite.before_mysql_"+stamp+".db")
	if err := copyFile(*sourceSQLite, sqliteBackup); err != nil {
		fatalf("backup sqlite: %v", err)
	}
	log.add("- sqlite_backup: %s", sqliteBackup)
	log.add("- sqlite_sha256: %s", sqliteHash)

	srcDB, err := sql.Open("sqlite", *sourceSQLite)
	if err != nil {
		fatalf("open sqlite: %v", err)
	}
	defer srcDB.Close()

	var integrity string
	if err := srcDB.QueryRow("PRAGMA integrity_check;").Scan(&integrity); err != nil {
		fatalf("sqlite integrity check: %v", err)
	}
	log.add("- sqlite_integrity_check: %s", integrity)
	if integrity != "ok" {
		fatalf("sqlite integrity check failed: %s", integrity)
	}

	tables, err := sqliteTables(srcDB)
	if err != nil {
		fatalf("list sqlite tables: %v", err)
	}
	log.add("- sqlite_table_count: %d", len(tables))

	adminDB, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fatalf("open mysql admin connection: %v", err)
	}
	defer adminDB.Close()
	if err := adminDB.Ping(); err != nil {
		fatalf("ping mysql admin connection: %v", err)
	}

	exists, err := mysqlDatabaseExists(adminDB, *targetDB)
	if err != nil {
		fatalf("check target database: %v", err)
	}
	if exists {
		tableCount, err := mysqlTableCount(adminDB, *targetDB)
		if err != nil {
			fatalf("check existing mysql tables: %v", err)
		}
		if tableCount > 0 {
			backupDB := fmt.Sprintf("%s_backup_%s", *targetDB, stamp)
			if err := cloneMySQLDatabase(adminDB, *targetDB, backupDB); err != nil {
				fatalf("backup existing mysql database: %v", err)
			}
			log.add("- mysql_existing_backup: %s", backupDB)
		} else {
			log.add("- mysql_existing_backup: skipped_existing_db_empty")
		}
	} else {
		log.add("- mysql_existing_backup: skipped_target_db_not_found")
	}

	if _, err := adminDB.Exec("DROP DATABASE IF EXISTS " + quoteMySQLIdent(*targetDB)); err != nil {
		fatalf("drop target database: %v", err)
	}
	if _, err := adminDB.Exec("CREATE DATABASE " + quoteMySQLIdent(*targetDB) + " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		fatalf("create target database: %v", err)
	}
	log.add("- mysql_target_recreated: true")

	targetCfg := cfg.Clone()
	configureMySQLConfig(targetCfg, *targetDB)
	target, err := sql.Open("mysql", targetCfg.FormatDSN())
	if err != nil {
		fatalf("open target mysql: %v", err)
	}
	defer target.Close()
	if err := target.Ping(); err != nil {
		fatalf("ping target mysql: %v", err)
	}
	_, _ = target.Exec("SET FOREIGN_KEY_CHECKS=0")

	log.add("")
	log.add("## Table Migration")
	log.add("")
	log.add("| table | rows | indexes |")
	log.add("| --- | ---: | ---: |")

	sourceCounts := make(map[string]int64)
	targetCounts := make(map[string]int64)
	totalRows := int64(0)
	for _, table := range tables {
		cols, err := sqliteColumns(srcDB, table)
		if err != nil {
			fatalf("read columns for %s: %v", table, err)
		}
		colTypes := map[string]string{}
		for _, col := range cols {
			colTypes[col.Name] = mysqlColumnType(col)
		}
		if err := createMySQLTable(target, table, cols); err != nil {
			fatalf("create table %s: %v", table, err)
		}
		rows, err := copyTableData(srcDB, target, table, cols)
		if err != nil {
			fatalf("copy table %s: %v", table, err)
		}
		indexes, err := sqliteIndexes(srcDB, table)
		if err != nil {
			fatalf("read indexes for %s: %v", table, err)
		}
		indexCount := 0
		for _, idx := range indexes {
			if idx.Origin == "pk" || len(idx.Columns) == 0 {
				continue
			}
			if err := createMySQLIndex(target, table, idx, colTypes); err != nil {
				fatalf("create index %s.%s: %v", table, idx.Name, err)
			}
			indexCount++
		}
		sourceCounts[table] = rows
		targetCount, err := countRows(target, table)
		if err != nil {
			fatalf("count target %s: %v", table, err)
		}
		targetCounts[table] = targetCount
		if targetCount != rows {
			fatalf("row count mismatch for %s: sqlite=%d mysql=%d", table, rows, targetCount)
		}
		totalRows += rows
		log.add("| %s | %d | %d |", table, rows, indexCount)
	}
	_, _ = target.Exec("SET FOREIGN_KEY_CHECKS=1")

	log.add("")
	log.add("## Count Check")
	log.add("")
	log.add("| table | sqlite | mysql |")
	log.add("| --- | ---: | ---: |")
	for _, table := range tables {
		log.add("| %s | %d | %d |", table, sourceCounts[table], targetCounts[table])
	}
	log.add("")
	log.add("- total_rows: %d", totalRows)
	log.add("- completed_at: %s", time.Now().Format(time.RFC3339))
	log.add("- result: PASS")
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}

func readProperties(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	props := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		props[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return props, scanner.Err()
}

func configureMySQLConfig(cfg *mysql.Config, dbName string) {
	cfg.DBName = dbName
	cfg.ParseTime = true
	cfg.Loc = time.Local
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}
	cfg.Params["charset"] = "utf8mb4"
}

func sanitizeDSN(cfg *mysql.Config) string {
	dbName := cfg.DBName
	if dbName == "" {
		dbName = "<server>"
	}
	return fmt.Sprintf("%s@%s(%s)/%s", cfg.User, cfg.Net, cfg.Addr, dbName)
}

func sha256File(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func sqliteTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func sqliteColumns(db *sql.DB, table string) ([]columnInfo, error) {
	rows, err := db.Query("PRAGMA table_info(" + quoteSQLiteIdent(table) + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cols []columnInfo
	for rows.Next() {
		var cid int
		var name string
		var typ sql.NullString
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return nil, err
		}
		cols = append(cols, columnInfo{
			Name:    name,
			Type:    typ.String,
			NotNull: notNull != 0,
			PKOrder: pk,
		})
	}
	return cols, rows.Err()
}

func sqliteIndexes(db *sql.DB, table string) ([]indexInfo, error) {
	rows, err := db.Query("PRAGMA index_list(" + quoteSQLiteIdent(table) + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var indexes []indexInfo
	for rows.Next() {
		var seq int
		var name string
		var unique int
		var origin string
		var partial int
		if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			return nil, err
		}
		cols, err := sqliteIndexColumns(db, name)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, indexInfo{
			Name:    name,
			Unique:  unique != 0,
			Origin:  origin,
			Partial: partial != 0,
			Columns: cols,
		})
	}
	return indexes, rows.Err()
}

func sqliteIndexColumns(db *sql.DB, indexName string) ([]string, error) {
	rows, err := db.Query("PRAGMA index_info(" + quoteSQLiteIdent(indexName) + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type item struct {
		seq  int
		name string
	}
	var items []item
	for rows.Next() {
		var seqno int
		var cid int
		var name string
		if err := rows.Scan(&seqno, &cid, &name); err != nil {
			return nil, err
		}
		items = append(items, item{seq: seqno, name: name})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].seq < items[j].seq })
	cols := make([]string, 0, len(items))
	for _, item := range items {
		cols = append(cols, item.name)
	}
	return cols, rows.Err()
}

func mysqlDatabaseExists(db *sql.DB, name string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM information_schema.schemata WHERE schema_name = ?", name).Scan(&count)
	return count > 0, err
}

func mysqlTableCount(db *sql.DB, schema string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE'", schema).Scan(&count)
	return count, err
}

func cloneMySQLDatabase(db *sql.DB, source, backup string) error {
	if _, err := db.Exec("CREATE DATABASE " + quoteMySQLIdent(backup) + " DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		return err
	}
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE' ORDER BY table_name", source)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		src := quoteMySQLIdent(source) + "." + quoteMySQLIdent(table)
		dst := quoteMySQLIdent(backup) + "." + quoteMySQLIdent(table)
		if _, err := db.Exec("CREATE TABLE " + dst + " LIKE " + src); err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO " + dst + " SELECT * FROM " + src); err != nil {
			return err
		}
	}
	return rows.Err()
}

func createMySQLTable(db *sql.DB, table string, cols []columnInfo) error {
	var defs []string
	var pkCols []columnInfo
	for _, col := range cols {
		nullSQL := "NULL"
		if col.PKOrder > 0 {
			nullSQL = "NOT NULL"
			pkCols = append(pkCols, col)
		}
		defs = append(defs, fmt.Sprintf("%s %s %s", quoteMySQLIdent(col.Name), mysqlColumnType(col), nullSQL))
	}
	sort.Slice(pkCols, func(i, j int) bool { return pkCols[i].PKOrder < pkCols[j].PKOrder })
	if len(pkCols) > 0 {
		var pkNames []string
		for _, col := range pkCols {
			pkNames = append(pkNames, quoteMySQLIdent(col.Name))
		}
		defs = append(defs, "PRIMARY KEY ("+strings.Join(pkNames, ", ")+")")
	}
	createSQL := "CREATE TABLE " + quoteMySQLIdent(table) + " (\n  " + strings.Join(defs, ",\n  ") + "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"
	_, err := db.Exec(createSQL)
	return err
}

func mysqlColumnType(col columnInfo) string {
	if col.PKOrder > 0 {
		return "varchar(191)"
	}
	typ := strings.ToUpper(strings.TrimSpace(col.Type))
	switch {
	case strings.Contains(typ, "INT"):
		return "bigint"
	case strings.Contains(typ, "REAL") || strings.Contains(typ, "FLOA") || strings.Contains(typ, "DOUB"):
		return "double"
	case strings.Contains(typ, "BLOB"):
		return "longblob"
	case strings.Contains(typ, "NUM") || strings.Contains(typ, "DEC"):
		return "decimal(20,6)"
	default:
		return "longtext"
	}
}

func copyTableData(src, dst *sql.DB, table string, cols []columnInfo) (int64, error) {
	if len(cols) == 0 {
		return 0, nil
	}
	var quotedCols []string
	var placeholders []string
	for _, col := range cols {
		quotedCols = append(quotedCols, quoteMySQLIdent(col.Name))
		placeholders = append(placeholders, "?")
	}
	var sqliteCols []string
	for _, col := range cols {
		sqliteCols = append(sqliteCols, quoteSQLiteIdent(col.Name))
	}

	rows, err := src.Query("SELECT " + strings.Join(sqliteCols, ", ") + " FROM " + quoteSQLiteIdent(table))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	tx, err := dst.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO " + quoteMySQLIdent(table) + " (" + strings.Join(quotedCols, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")")
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(cols))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var count int64
	for rows.Next() {
		for i := range values {
			values[i] = nil
		}
		if err := rows.Scan(scanArgs...); err != nil {
			_ = tx.Rollback()
			return 0, err
		}
		if _, err := stmt.Exec(values...); err != nil {
			_ = tx.Rollback()
			return 0, err
		}
		count++
	}
	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	return count, tx.Commit()
}

func createMySQLIndex(db *sql.DB, table string, idx indexInfo, colTypes map[string]string) error {
	prefix := "CREATE "
	if idx.Unique {
		prefix += "UNIQUE "
	}
	var cols []string
	for _, col := range idx.Columns {
		colSQL := quoteMySQLIdent(col)
		if isTextIndexType(colTypes[col]) {
			colSQL += "(64)"
		}
		cols = append(cols, colSQL)
	}
	_, err := db.Exec(prefix + "INDEX " + quoteMySQLIdent(idx.Name) + " ON " + quoteMySQLIdent(table) + " (" + strings.Join(cols, ", ") + ")")
	return err
}

func isTextIndexType(typ string) bool {
	typ = strings.ToLower(typ)
	return strings.Contains(typ, "text") || strings.Contains(typ, "char") || strings.Contains(typ, "blob")
}

func countRows(db *sql.DB, table string) (int64, error) {
	var count int64
	err := db.QueryRow("SELECT COUNT(1) FROM " + quoteMySQLIdent(table)).Scan(&count)
	return count, err
}

func quoteMySQLIdent(s string) string {
	return "`" + strings.ReplaceAll(s, "`", "``") + "`"
}

func quoteSQLiteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
