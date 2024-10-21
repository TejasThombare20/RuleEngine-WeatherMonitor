package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NodeType string

const (
	OperatorNode  NodeType = "operator"
	ConditionNode NodeType = "condition"
	ReferenceNode NodeType = "reference"
)

type Node struct {
	Type        NodeType           `bson:"type" json:"type"`
	Value       *string            `bson:"value" json:"value"`
	Left        *Node              `bson:"left,omitempty" json:"left,omitempty"`
	Right       *Node              `bson:"right,omitempty" json:"right,omitempty"`
	ReferenceID primitive.ObjectID `bson:"reference_id,omitempty" json:"reference_id,omitempty"`
}

type Rule struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Description *string            `json:"description" validate:"required"`
	Name        *string            ` json:"name" validate:"required"`
	RootNode    *Node              `json:"root_node"`
}
