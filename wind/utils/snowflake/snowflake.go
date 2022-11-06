package utils

import (
	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
)

func InitSnowFlake(nodeId int64) {
	var err error
	node, err = snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}
}

func GetSnowflakeId() string {
	if node == nil {
		panic("snowflake not init")
	}
	return node.Generate().String()
}

func GetSnowflakeIdInt64() int64 {
	if node == nil {
		panic("snowflake not init")
	}
	return node.Generate().Int64()
}
