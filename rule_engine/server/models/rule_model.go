package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NodeType string

const (
	OperatorNode  NodeType = "operator"
	ConditionNode NodeType = "condition"
	ReferenceNode NodeType = "reference"
)

type Value struct {
	StringValue *string
	IntValue    *int
}

type Node struct {
	Type        NodeType               `bson:"type" json:"type"`
	Value       *string                `bson:"value" json:"value"`
	Left        *Node                  `bson:"left,omitempty" json:"left,omitempty"`
	Right       *Node                  `bson:"right,omitempty" json:"right,omitempty"`
	RuleIDs     []primitive.ObjectID   `bson:"rule_ids,omitempty" json:"rule_ids,omitempty"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	ReferenceID primitive.ObjectID     `bson:"reference_id,omitempty" json:"reference_id,omitempty"`
}

type Rule struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Description *string            `json:"description" validate:"required"`
	Name        *string            ` json:"name" validate:"required"`
	RootNode    *Node              `json:"root_node"`
}

type EvaluationRequest struct {
	RuleID primitive.ObjectID     `json:"rule_id" binding:"required"`
	Data   map[string]interface{} `json:"data" binding:"required"`
}

type EvaluationResponse struct {
	RuleID primitive.ObjectID `json:"rule_id"`
	Result bool               `json:"result"`
	Name   *string            `json:"name,omitempty"`
}
