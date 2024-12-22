package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
)

// Tree 模型对应 ClickHouse 的 tree 表
type Tree struct {
	//ID     uint   `gorm:"column:id;primaryKey" json:"id"`
	Node1  string `gorm:"column:node1" json:"node1"`
	Node2  string `gorm:"column:node2" json:"node2"`
	Node3  string `gorm:"column:node3" json:"node3"`
	Node4  string `gorm:"column:node4" json:"node4"`
	Node5  string `gorm:"column:node5" json:"node5"`
	Node6  string `gorm:"column:node6" json:"node6"`
	Node7  string `gorm:"column:node7" json:"node7"`
	Node8  string `gorm:"column:node8" json:"node8"`
	Node9  string `gorm:"column:node9" json:"node9"`
	Node10 string `gorm:"column:node10" json:"node10"`
}

// TableName 指定表名
func (Tree) TableName() string {
	return "tree"
}

// TreeNode 定义前端树形结构的节点
type TreeNode struct {
	Name     string      `json:"name"`
	Children []*TreeNode `json:"children,omitempty"`
}

var db *gorm.DB

// 初始化数据库连接
func initDB() {
	var err error
	dsn := "clickhouse://default:36UdcuHtpZSfRj@192.168.6.41:30900/dosec?dial_timeout=10s&read_timeout=20s"

	// 使用 gorm-clickhouse 连接 ClickHouse
	db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 配置命名策略
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}

	// 自动迁移（创建表）
	err = db.AutoMigrate(&Tree{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Successfully connected to ClickHouse and migrated database.")
}

// 获取树形数据
func getTreeData() ([]*TreeNode, error) {
	var trees []Tree
	// 查询所有数据
	result := db.Find(&trees)
	if result.Error != nil {
		return nil, result.Error
	}

	// 构建树形结构
	treeMap := make(map[string]*TreeNode)

	for _, t := range trees {
		parent := &TreeNode{Name: t.Node1}
		if _, exists := treeMap[t.Node1]; !exists {
			treeMap[t.Node1] = parent
		}

		// 递归添加子节点
		current := treeMap[t.Node1]
		for i := 2; i <= 10; i++ {
			nodeName := ""
			switch i {
			case 2:
				nodeName = t.Node2
			case 3:
				nodeName = t.Node3
			case 4:
				nodeName = t.Node4
			case 5:
				nodeName = t.Node5
			case 6:
				nodeName = t.Node6
			case 7:
				nodeName = t.Node7
			case 8:
				nodeName = t.Node8
			case 9:
				nodeName = t.Node9
			case 10:
				nodeName = t.Node10
			}

			if nodeName == "" {
				break
			}

			// 检查子节点是否存在
			var child *TreeNode
			found := false
			for _, c := range current.Children {
				if c.Name == nodeName {
					child = c
					found = true
					break
				}
			}
			if !found {
				child = &TreeNode{Name: nodeName}
				current.Children = append(current.Children, child)
			}
			current = child
		}
	}

	// 转换为切片
	var tree []*TreeNode
	for _, node := range treeMap {
		tree = append(tree, node)
	}

	return tree, nil
}

// API 处理函数
func treeHandler(w http.ResponseWriter, r *http.Request) {
	tree, err := getTreeData()
	if err != nil {
		http.Error(w, "Failed to get tree data", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	// 编码为 JSON
	json.NewEncoder(w).Encode(tree)
}

func main() {
	// 初始化数据库
	initDB()

	// 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tree", treeHandler)

	// 配置 CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:8081"}, // 前端地址
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// 使用 CORS 中间件
	handler := c.Handler(mux)

	// 启动服务器
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
