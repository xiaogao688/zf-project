package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var testDB *sql.DB

func InitTestDB() {
	// 打开（或创建）SQLite 数据库文件
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	//defer db.Close()

	// 执行建表和建索引语句
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS artifacts(
            id INTEGER PRIMARY KEY,
            group_id TEXT,
            artifact_id TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS indices(
            artifact_id INTEGER,
            version TEXT,
            sha1 BLOB,
            archive_type TEXT,
            FOREIGN KEY (artifact_id) REFERENCES artifacts(id)
        );`,
		`CREATE UNIQUE INDEX IF NOT EXISTS artifacts_idx 
            ON artifacts(artifact_id, group_id);`,
		`CREATE INDEX IF NOT EXISTS indices_artifact_idx 
            ON indices(artifact_id);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS indices_sha1_idx 
            ON indices(sha1);`,
	}

	for _, sqlStmt := range stmts {
		if _, err := db.Exec(sqlStmt); err != nil {
			log.Fatalf("执行语句失败: %v\nSQL: %s", err, sqlStmt)
		}
	}
	fmt.Println("✅ 表和索引创建完毕")
	testDB = db
}

func InsertArtifacts(artifacts []Artifact) {
	// 批量插入 artifacts
	if len(artifacts) > 0 {
		var placeholders []string
		var args []interface{}
		for _, a := range artifacts {
			placeholders = append(placeholders, "(?,?,?)")
			args = append(args, a.ID, a.GroupID, a.ArtifactID)
		}
		stmt := fmt.Sprintf("INSERT OR IGNORE INTO artifacts(id, group_id, artifact_id) VALUES %s", strings.Join(placeholders, ","))
		if _, err := testDB.Exec(stmt, args...); err != nil {
			log.Fatalf("批量插入 artifacts 失败: %v", err)
		}
		//fmt.Println("✅ artifacts 数据批量插入完毕")
	}
}

func InsertIndices(indices []Index) {
	// 批量插入 indices
	if len(indices) > 0 {
		var placeholders []string
		var args []interface{}
		for _, idx := range indices {
			placeholders = append(placeholders, "(?,?,?,?)")
			args = append(args, idx.ArtifactID, idx.Version, idx.SHA1, idx.ArchiveType)
		}
		stmt := fmt.Sprintf("INSERT OR IGNORE INTO indices(artifact_id, version, sha1, archive_type) VALUES %s", strings.Join(placeholders, ","))
		if _, err := testDB.Exec(stmt, args...); err != nil {
			log.Fatalf("批量插入 indices 失败: %v", err)
		}
		//fmt.Println("✅ indices 数据批量插入完毕")
	}
}
