package main

import (
	"testing"
	"xorm.io/builder"
	"xorm.io/xorm"
)

func TestBatchUpdateData(t *testing.T) {
	type args struct {
		engine       *xorm.Engine
		items        []K8sPod
		updateFields []string
		uniqueField  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "t1",
			args: args{
				engine: nil,
				items: []K8sPod{
					{
						UniqueStr:         "a1",
						ContainerReadyNum: 11,
						ContainerTotalNum: 12,
					},
					{
						UniqueStr:         "a2",
						ContainerReadyNum: 21,
						ContainerTotalNum: 22,
					},
					{
						UniqueStr:         "a3",
						ContainerReadyNum: 31,
						ContainerTotalNum: 32,
					},
					{
						UniqueStr:         "a4",
						ContainerReadyNum: 41,
						ContainerTotalNum: 42,
					},
				},
				updateFields: []string{"container_ready_num", "container_total_num"},
				uniqueField:  "unique_str",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchUpdateData(tt.args.engine, tt.args.items, tt.args.updateFields, tt.args.uniqueField); (err != nil) != tt.wantErr {
				t.Errorf("BatchUpdateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBatchUpdateData2(t *testing.T) {
	type args struct {
		engine       *xorm.Engine
		items        []builder.Eq
		updateFields []string
		uniqueField  string
		tableName    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				engine: nil,
				items: []builder.Eq{
					{
						"container_total_num": 10011,
						"container_ready_num": 10012,
						"unique_str":          "a1",
					},
					{
						"container_total_num": 10021,
						"container_ready_num": 10022,
						"unique_str":          "a2",
					},
					{
						"container_total_num": 10031,
						"container_ready_num": 10032,
						"unique_str":          "a3",
					},
					{
						"container_total_num": 10041,
						"container_ready_num": 10042,
						"unique_str":          "a4",
					},
				},
				updateFields: []string{"container_total_num", "container_ready_num"},
				uniqueField:  "unique_str",
				tableName:    "k8s_pod",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchUpdateData2(tt.args.engine, tt.args.items, tt.args.updateFields, tt.args.uniqueField, tt.args.tableName); (err != nil) != tt.wantErr {
				t.Errorf("BatchUpdateData2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
