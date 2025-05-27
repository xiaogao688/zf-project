package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/syndtr/goleveldb/leveldb"
)

type Artifact struct {
	ID         int
	GroupID    string
	ArtifactID string
}

type Index struct {
	ArtifactID  int
	Version     string
	SHA1        []byte
	ArchiveType string
}

func main() {
	// 初始化 SQLite 连接
	sqliteDB, err := sql.Open("sqlite3", "./trivy-java.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDB.Close()

	InitTestDB()

	// 写入 LevelDB
	artifactsDB, indicesDB := initLevelDB()
	defer artifactsDB.Close()
	defer indicesDB.Close()

	// 测试写入性能
	//TestArtifactsWrite(sqliteDB, artifactsDB)
	//TestIndicesWrite(sqliteDB, indicesDB)

	// migrateToLevelDB(artifactsDB, indicesDB, artifacts, indices)

	// 性能测试
	TestReadArtifact(sqliteDB, artifactsDB)
	TestReadIndex(sqliteDB, indicesDB)
	// testReadPerformance(sqliteDB, artifactsDB, indicesDB, artifacts, indices)
}

func TestArtifactsWrite(db *sql.DB, artifactsDB *leveldb.DB) {
	// 写入耗时
	sqliteTime := int64(0)
	leveldbTime := int64(0)

	defer func() {
		fmt.Println("sqlite time: ", sqliteTime, "leveldb time: ", leveldbTime)
	}()

	// 读取 artifacts
	artifacts := make([]Artifact, 0, 1000)
	var lastID int
	for {
		rows, err := db.Query(`
            SELECT id, group_id, artifact_id
            FROM artifacts
            WHERE ID > ?
            ORDER BY ID
            LIMIT 1000`, lastID)
		if err != nil {
			panic(err)
		}

		count := 0
		for rows.Next() {
			var a Artifact
			if err := rows.Scan(&a.ID, &a.GroupID, &a.ArtifactID); err != nil {
				rows.Close()
				panic(err)
			}
			artifacts = append(artifacts, a)
			lastID = a.ID
			count++
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			panic(err)
		}

		if count == 0 {
			break
		}

		// 写sqlite
		t := time.Now().UnixNano()
		InsertArtifacts(artifacts)
		sqliteTime += time.Now().UnixNano() - t

		// leveldb
		t = time.Now().UnixNano()
		// 写入 artifacts
		for _, a := range artifacts {
			key := fmt.Sprintf("%s:%s", a.GroupID, a.ArtifactID)
			if err := artifactsDB.Put([]byte(key), []byte{}, nil); err != nil {
				log.Fatal(err)
			}
		}
		leveldbTime += time.Now().UnixNano() - t

		// 清空
		artifacts = artifacts[:0]
	}
}

func TestIndicesWrite(db *sql.DB, indicesDB *leveldb.DB) {
	// 写入耗时
	sqliteTime := int64(0)
	leveldbTime := int64(0)

	// 读取 artifacts
	indices := make([]Index, 0, 1000)
	var lastID int
	for {
		rows, err := db.Query(`
            SELECT artifact_id, version, sha1, archive_type FROM indices
            LIMIT 1000 OFFSET ?`, lastID)
		if err != nil {
			panic(err)
		}

		count := 0
		for rows.Next() {
			var a Index
			if err := rows.Scan(&a.ArtifactID, &a.Version, &a.SHA1, &a.ArchiveType); err != nil {
				rows.Close()
				panic(err)
			}
			indices = append(indices, a)
			count++
		}
		lastID += 1000
		rows.Close()
		if err := rows.Err(); err != nil {
			panic(err)
		}

		if count == 0 {
			break
		}

		// 写sqlite
		t := time.Now().UnixNano()
		InsertIndices(indices)
		sqliteTime += time.Now().UnixNano() - t

		// leveldb
		t = time.Now().UnixNano()
		// 写入 artifacts
		for _, idx := range indices {
			value := fmt.Sprintf("%d:%s:%s", idx.ArtifactID, idx.Version, idx.ArchiveType)
			if err := indicesDB.Put(idx.SHA1, []byte(value), nil); err != nil {
				log.Fatal(err)
			}
		}
		leveldbTime += time.Now().UnixNano() - t

		// 清空
		indices = indices[:0]
	}
}

func initLevelDB() (*leveldb.DB, *leveldb.DB) {
	artifactsDB, err := leveldb.OpenFile("./artifacts.db", nil)
	if err != nil {
		log.Fatal(err)
	}

	indicesDB, err := leveldb.OpenFile("./indices.db", nil)
	if err != nil {
		log.Fatal(err)
	}

	return artifactsDB, indicesDB
}

func TestReadIndex(sqliteDB *sql.DB, indicesDB *leveldb.DB) {
	indices := []Index{
		{ArtifactID: 1, Version: "0.3-3", ArchiveType: "jar"},
		{ArtifactID: 2, Version: "0.12.3", ArchiveType: "jar"},
		{ArtifactID: 2, Version: "0.13.0", ArchiveType: "jar"},
		{ArtifactID: 2, Version: "1.4.0", ArchiveType: "jar"},
		{ArtifactID: 3, Version: "1.4.0", ArchiveType: "jar"},
		{ArtifactID: 4, Version: "1.0", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.7.0", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.7.1", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.8.0", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.8.1", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.8.1.1", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.8.2", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.8.3", ArchiveType: "jar"},
		{ArtifactID: 5, Version: "0.9.0", ArchiveType: "jar"},
		{ArtifactID: 6, Version: "0.51", ArchiveType: "jar"},
		{ArtifactID: 6, Version: "0.6.1", ArchiveType: "jar"},
		{ArtifactID: 7, Version: "0.51", ArchiveType: "jar"},
		{ArtifactID: 7, Version: "0.6.1", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.7.0", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.7.1", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.8.0", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.8.1", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.8.1.1", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.8.2", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.8.3", ArchiveType: "jar"},
		{ArtifactID: 8, Version: "0.9.0", ArchiveType: "jar"},
	}

	// 测试 artifacts 读取
	var sqliteTime, leveldbTime int64
	defer func() {
		fmt.Println("Index reed sqlite time: ", sqliteTime, "leveldb time: ", leveldbTime)
	}()

	for _, ind := range indices {
		var id int
		t := time.Now().UnixNano()
		err := sqliteDB.QueryRow(
			"SELECT sha1 FROM indices WHERE artifact_id = ? AND version = ? AND archive_type = ?",
			ind.ArtifactID, ind.Version, ind.ArchiveType,
		).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		sqliteTime += time.Now().UnixNano() - t

		t = time.Now().UnixNano()
		// 测试 LevelDB
		key := fmt.Sprintf("%d:%s:%s", ind.ArtifactID, ind.Version, ind.ArchiveType)
		_, err = indicesDB.Get([]byte(key), nil)
		if err != nil {
			log.Fatal(err)
		}
		leveldbTime += time.Now().UnixNano() - t
	}

	fmt.Println("")
}

func TestReadArtifact(sqliteDB *sql.DB, artifactsDB *leveldb.DB) {
	artifacts := []Artifact{
		{GroupID: "HTTPClient", ArtifactID: "HTTPClient"},
		{GroupID: "abbot", ArtifactID: "abbot"},
		{GroupID: "abbot", ArtifactID: "costello"},
		{GroupID: "academy.alex", ArtifactID: "custommatcher"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-cas"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-catalina-common"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-catalina-server"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-catalina"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-jboss-lib"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-jboss"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-jetty-ext"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-jetty"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-resin-lib"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-resin"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-taglib"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security-tiger"},
		{GroupID: "acegisecurity", ArtifactID: "acegi-security"},
		{GroupID: "activecluster", ArtifactID: "activecluster"},
		{GroupID: "activeio", ArtifactID: "activeio"},
		{GroupID: "activemq", ArtifactID: "activemq-axis"},
		{GroupID: "activemq", ArtifactID: "activemq-container"},
		{GroupID: "activemq", ArtifactID: "activemq-core"},
		{GroupID: "activemq", ArtifactID: "activemq-gbean-management"},
		{GroupID: "activemq", ArtifactID: "activemq-gbean"},
	}

	// 测试 artifacts 读取
	var sqliteTime, leveldbTime int64
	defer func() {
		fmt.Println("Artifact reed sqlite time: ", sqliteTime, "leveldb time: ", leveldbTime)
	}()

	for _, artifact := range artifacts {
		var id int
		t := time.Now().UnixNano()
		err := sqliteDB.QueryRow(
			"SELECT id FROM artifacts WHERE group_id = ? AND artifact_id = ?",
			artifact.GroupID, artifact.ArtifactID,
		).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		sqliteTime += time.Now().UnixNano() - t

		t = time.Now().UnixNano()
		// 测试 LevelDB
		key := fmt.Sprintf("%s:%s", artifact.GroupID, artifact.ArtifactID)
		_, err = artifactsDB.Get([]byte(key), nil)
		if err != nil {
			log.Fatal(err)
		}
		leveldbTime += time.Now().UnixNano() - t
	}
}
