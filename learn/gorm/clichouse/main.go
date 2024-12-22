package main

import (
	"fmt"
	"log"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

/*
CREATE TABLE tree
(
    node1 String,
    node2 String,
    node3 String,
    node4 String,
    node5 String,
    node6 String,
    node7 String,
    node8 String,
    node9 String,
    node10 String
)
ENGINE = MergeTree()
ORDER BY (node1, node2, node3, node4, node5, node6, node7, node8, node9, node10);
*/

type Tree struct {
	Node1  string `gorm:"column:node1"`
	Node2  string `gorm:"column:node2"`
	Node3  string `gorm:"column:node3"`
	Node4  string `gorm:"column:node4"`
	Node5  string `gorm:"column:node5"`
	Node6  string `gorm:"column:node6"`
	Node7  string `gorm:"column:node7"`
	Node8  string `gorm:"column:node8"`
	Node9  string `gorm:"column:node9"`
	Node10 string `gorm:"column:node10"`
}

func main() {
	// 连接 ClickHouse 数据库
	dsn := "clickhouse://default:36UdcuHtpZSfRj@192.168.6.41:30900/dosec?dial_timeout=10s&read_timeout=20s"
	//dsn := "tcp://192.168.6.41:30900/?username=default&password=36UdcuHtpZSfRj&database=dosec"
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	// 自动创建表
	//if err := db.AutoMigrate(&Tree{}); err != nil {
	//	log.Fatalf("Failed to migrate database: %v", err)
	//}
	//fmt.Println("Table `tree` created successfully!")

	//// 插入数据
	//tree := Tree{
	//	Node1:  "value1",
	//	Node2:  "value2",
	//	Node3:  "value3",
	//	Node4:  "value4",
	//	Node5:  "value5",
	//	Node6:  "value6",
	//	Node7:  "value7",
	//	Node8:  "value8",
	//	Node9:  "value9",
	//	Node10: "value10",
	//}
	//if err := db.Create(&tree).Error; err != nil {
	//	log.Fatalf("Failed to insert data: %v", err)
	//}
	//fmt.Println("Data inserted successfully into `tree`!")

	// 查询数据
	var results []Tree
	if err := db.Table("tree").Find(&results).Error; err != nil {
		log.Fatalf("Failed to query data: %v", err)
	}
	fmt.Println("Query result:")
	for _, result := range results {
		fmt.Printf("%+v\n", result)
	}
}
