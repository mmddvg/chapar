package postgres

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Initialize(nodeID int64) {
	var err error

	// snowflake.Epoch = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6

	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("Failed to initialize Snowflake node: %v", err)
	}
}

func GenerateId() uint64 {
	return uint64(node.Generate().Int64())
}
