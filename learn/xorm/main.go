package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type K8sPod struct {
	ID                   int               `xorm:"pk autoincr 'id'" json:"id,omitempty"`
	UID                  string            `xorm:"'uid'" json:"uid,omitempty"` // + uid
	PodName              string            `xorm:"'pod_name'" json:"pod_name,omitempty"`
	Type                 string            `xorm:"'type'" json:"type,omitempty"` // workloads
	Namespace            string            `xorm:"'namespace'" json:"namespace,omitempty"`
	Labels               map[string]string `xorm:"'labels'" json:"labels,omitempty"`
	PortNames            []string          `xorm:"'ports'" json:"port_names,omitempty"`
	NodeName             string            `xorm:"'node_name'" json:"node_name,omitempty"`
	ServiceAccount       string            `xorm:"'service_account'" json:"service_account,omitempty"`
	ServiceAccountID     int               `xorm:"'service_account_id'" json:"service_account_id,omitempty"`
	HostIP               string            `xorm:"'host_ip'" json:"host_ip,omitempty"`
	HostIPV6             string            `xorm:"'host_ipv6'" json:"host_ipv_6,omitempty"`
	PodIP                string            `xorm:"'pod_ip'" json:"pod_ip,omitempty"`
	PodIPInt             uint32            `xorm:"'pod_ip_int'" json:"pod_ip_int,omitempty"`
	PodIPs               []string          `xorm:"'pod_ips'" json:"pod_i_ps,omitempty"`
	Containers           []string          `xorm:"'containers'" json:"containers,omitempty"`
	ClusterID            int               `xorm:"'cluster_id'" json:"cluster_id,omitempty"`
	NamespaceUid         string            `xorm:"'namespace_uid'" json:"namespace_uid,omitempty"` // 命名空间 uid
	HostNetwork          bool              `xorm:"'host_network'" json:"host_network,omitempty"`
	HasPolicy            bool              `xorm:"default false 'has_policy'" json:"has_policy,omitempty"`
	CreationTimestamp    time.Time         `xorm:"TIMESTAMP 'creation_timestamp'" json:"creation_timestamp"`            // + k8s 资源创建时间
	Phase                string            `xorm:"'phase'" json:"phase,omitempty"`                                      // + 状态
	Ages                 string            `xorm:"'ages'" json:"ages,omitempty"`                                        // + 存活时间
	Restart              int               `xorm:"'restart'" json:"restart,omitempty"`                                  // + 重启次数
	Spec                 string            `xorm:"'spec'" json:"spec,omitempty"`                                        // + 详情信息 json
	IsHistory            bool              `xorm:"'is_history'" json:"is_history,omitempty"`                            // + 是否历史
	CreateTime           time.Time         `xorm:"created TIMESTAMP 'create_time'" json:"create_time"`                  // + 创建时间
	UpdateTime           time.Time         `xorm:"updated TIMESTAMP 'update_time'" json:"update_time"`                  // + 更新时间
	DeleteTime           time.Time         `xorm:"TIMESTAMP 'delete_time'" json:"delete_time"`                          // + 删除时间
	UniqueStr            string            `xorm:"'unique_str'" json:"unique_str,omitempty"`                            // 唯一字符串 格式:md5(cluster_id:namespace_uuid:asset_name)
	ContainerReadyNum    int               `xorm:"'container_ready_num'" json:"container_ready_num,omitempty"`          // 就绪的容器数量
	ContainerTotalNum    int               `xorm:"'container_total_num'" json:"container_total_num,omitempty"`          // 总的容器数量
	InitContainerImages  []string          `xorm:"'init_container_images'" json:"init_container_images,omitempty"`      // 初始化容器镜像
	ContainerImages      []string          `xorm:"'container_images'" json:"container_images,omitempty"`                // 容器镜像,不包含init容器
	InitContainerHashIDs []string          `xorm:"'init_container_hash_ids'" json:"init_container_hash_i_ds,omitempty"` // 初始化容器hash_id
	Annotations          map[string]string `xorm:"'annotations'" json:"annotations,omitempty"`                          // 注释
}

func BatchUpdateData(engine *xorm.Engine, items []K8sPod, updateFields []string, uniqueField string) error {
	if len(items) == 0 || len(updateFields) == 0 || uniqueField == "" {
		return nil
	}

	vType := reflect.TypeOf(items[0])
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}
	if vType.Kind() != reflect.Struct {
		return fmt.Errorf("T must be a struct")
	}

	caseMap := make(map[string][]string)
	var idList []string

	for _, item := range items {
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		val, err := GetFieldValueByTag(&item, "xorm", uniqueField)
		if err != nil {
			return err
		}
		uid := fmt.Sprintf("'%v'", val)
		idList = append(idList, uid)

		// 构建每个更新字段的 CASE WHEN 子句
		for _, field := range updateFields {
			val, err = GetFieldValueByTag(item, "xorm", field)
			if err != nil {
				return err
			}
			caseLine := fmt.Sprintf("WHEN %s THEN %v", uid, val)
			caseMap[field] = append(caseMap[field], caseLine)
		}
	}

	// 拼接 SET 子句
	var setClauses []string
	for _, field := range updateFields {
		caseLines := strings.Join(caseMap[field], " ")
		setClause := fmt.Sprintf("%s = CASE %s %s END", field, uniqueField, caseLines)
		setClauses = append(setClauses, setClause)
	}

	// 拼接 WHERE 子句
	whereClause := fmt.Sprintf("%s IN (%s)", uniqueField, strings.Join(idList, ", "))

	// 默认使用结构体名为表名（小写）
	tableName := strings.ToLower(vType.Name())

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName,
		strings.Join(setClauses, ", "),
		whereClause,
	)

	_, err := engine.Exec(sql)
	return err
}

func GetFieldValueByTag(obj interface{}, tagName, tagValue string) (reflect.Value, error) {
	tagValue = fmt.Sprintf("'%s'", tagValue)
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	if val.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("expected struct, got %s", val.Kind())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == tagValue {
			return val.Field(i), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("no field with tag %s=%s found", tagName, tagValue)
}

func BatchUpdateData2(engine *xorm.Engine, items []builder.Eq, updateFields []string, uniqueField, tableName string) error {
	// If no items or update fields or uniqueField/tableName is missing, nothing to do.
	if len(items) == 0 || len(updateFields) == 0 || uniqueField == "" || tableName == "" {
		return nil
	}

	// 用于存储每个 updateField 对应的 CASE WHEN 子句
	caseMap := make(map[string][]string)
	// 存储所有 uniqueField 的值，用于构建 WHERE 子句
	var idList []string

	// 遍历每个项，提取 uniqueField 以及更新字段的值
	for _, item := range items {
		// 取 uniqueField 的值
		uidValue, ok := item[uniqueField]
		if !ok {
			return fmt.Errorf("uniqueField %s not found in item", uniqueField)
		}
		// 对 uniqueField 的值进行格式化（如果是字符串则加引号，否则直接格式化）
		uid := ""
		switch v := uidValue.(type) {
		case string:
			uid = fmt.Sprintf("'%s'", v)
		default:
			uid = fmt.Sprintf("'%v'", v)
		}
		idList = append(idList, uid)

		// 遍历每个需要更新的字段
		for _, field := range updateFields {
			fieldValue, ok := item[field]
			if !ok {
				return fmt.Errorf("field %s not found in item", field)
			}
			// 根据字段值类型决定是否需要加引号
			var fieldValueStr string
			switch v := fieldValue.(type) {
			case string:
				fieldValueStr = fmt.Sprintf("'%s'", v)
			default:
				fieldValueStr = fmt.Sprintf("%v", v)
			}
			// 构造 CASE WHEN 子句：WHEN <uniqueField值> THEN <字段新值>
			caseLine := fmt.Sprintf("WHEN %s THEN %s", uid, fieldValueStr)
			caseMap[field] = append(caseMap[field], caseLine)
		}
	}

	// 构建 SET 子句，每个更新字段使用 CASE 语句
	var setClauses []string
	for _, field := range updateFields {
		// 拼接所有 CASE WHEN 子句
		caseLines := strings.Join(caseMap[field], " ")
		setClause := fmt.Sprintf("%s = CASE %s %s END", field, uniqueField, caseLines)
		setClauses = append(setClauses, setClause)
	}

	// 构建 WHERE 子句
	whereClause := fmt.Sprintf("%s IN (%s)", uniqueField, strings.Join(idList, ", "))

	// 最终的 SQL 语句
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(setClauses, ", "), whereClause)

	_, err := engine.Exec(sql)
	return err
}
