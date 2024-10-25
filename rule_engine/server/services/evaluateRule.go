package services

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/TejasThombare20/rule-engine/models"
)

func evaluateNode(node *models.Node, data map[string]interface{}) bool {
	if node == nil {
		return false
	}

	switch node.Type {
	case models.OperatorNode:
		return evaluateOperatorNode(node, data)
	case models.ConditionNode:
		return evaluateConditionNode(node, data)
	default:
		return false
	}
}

func evaluateOperatorNode(node *models.Node, data map[string]interface{}) bool {
	if node.Value == nil {
		return false
	}

	switch *node.Value {
	case "AND":
		return evaluateNode(node.Left, data) && evaluateNode(node.Right, data)
	case "OR":
		return evaluateNode(node.Left, data) || evaluateNode(node.Right, data)
	default:
		return false
	}
}

func evaluateConditionNode(node *models.Node, data map[string]interface{}) bool {
	if node.Value == nil {
		return false
	}

	// Parse condition string (e.g., "age > 30")
	condition := parseCondition(node)
	if condition == nil {
		return false
	}

	// Get actual value from data
	actualValue, exists := data[condition.Field]
	if !exists {
		return false
	}

	return evaluateConditionValue(condition, actualValue)
}

type parsedCondition struct {
	Field    string
	Operator string
	Value    string
}

func parseCondition(node *models.Node) *parsedCondition {

	if node.Left == nil || node.Right == nil || node.Value == nil || node.Left.Value == nil || node.Right.Value == nil {
		return nil
	}

	return &parsedCondition{
		Field:    *node.Left.Value,
		Operator: *node.Value,
		Value:    *node.Right.Value,
	}

}

func evaluateConditionValue(cond *parsedCondition, actualValue interface{}) bool {
	// Convert condition value based on actual value type
	var condValue interface{}
	var err error

	switch actualValue.(type) {
	case float64:
		condValue, err = strconv.ParseFloat(cond.Value, 64)
	case int:
		condValue, err = strconv.Atoi(cond.Value)
	case bool:
		condValue, err = strconv.ParseBool(cond.Value)
	case string:
		condValue = cond.Value
	default:
		return false
	}

	fmt.Printf("condition value  %v, %T\n", cond.Value, cond.Value)
	fmt.Printf("parsed condition value %v,  %T\n", condValue, condValue)
	fmt.Printf("actual value %v,  %T\n", actualValue, actualValue)
	if err != nil {
		fmt.Println("error parsing condition value: ", err)
		return false
	}

	fmt.Printf("operator%T\n", cond.Operator)

	// Compare values based on operator
	switch cond.Operator {
	case "=":
		return reflect.DeepEqual(actualValue, condValue)
	case "!=":
		return !reflect.DeepEqual(actualValue, condValue)
	case ">", "<", ">=", "<=":
		return compareValues(actualValue, condValue, cond.Operator)
	default:
		return false
	}
}

func compareValues(actual, condition interface{}, operator string) bool {
	// Convert to float64 for numerical comparison
	var actualFloat, condFloat float64

	switch v := actual.(type) {
	case int:
		actualFloat = float64(v)
	case float64:
		actualFloat = v
	default:
		return false
	}

	switch v := condition.(type) {
	case int:
		condFloat = float64(v)
	case float64:
		condFloat = v
	default:
		return false
	}

	switch operator {
	case ">":
		return actualFloat > condFloat
	case "<":
		return actualFloat < condFloat
	case ">=":
		return actualFloat >= condFloat
	case "<=":
		return actualFloat <= condFloat
	default:
		return false
	}
}

func (s *RuleService) EvaluateRule(req *models.EvaluationRequest) (*models.EvaluationResponse, error) {
	// Get rule from repository
	rule, err := s.repo.FindByID(req.RuleID)
	if err != nil {
		return nil, err
	}

	// Evaluate rule
	result := evaluateNode(rule.RootNode, req.Data)

	return &models.EvaluationResponse{
		RuleID: rule.ID,
		Result: result,
		Name:   rule.Name,
	}, nil
}
