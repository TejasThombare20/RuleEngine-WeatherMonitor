package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NodeType string

const (
	OperatorNode  NodeType = "operator"
	ConditionNode NodeType = "condition"
)

type Node struct {
	Type  NodeType `bson:"type" json:"type"`
	Value *string  `bson:"value" json:"value"`
	Left  *Node    `bson:"left" json:"left"`
	Right *Node    `bson:"right" json:"right"`
}

type Rule struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Description *string            `json:"description" validate:"required"`
	Name        *string            ` json:"name" validate:"required"`
	RootNode    *Node              `json:"root_node"`
}
