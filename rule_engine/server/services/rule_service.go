package services

import (
	"fmt"
	"strings"

	"regexp"

	"github.com/TejasThombare20/rule-engine/models"
	"github.com/TejasThombare20/rule-engine/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleService struct {
	repo *repositories.RuleRepository
}

func NewRuleService(repo *repositories.RuleRepository) *RuleService {
	return &RuleService{repo: repo}
}

// func tokenize(rule string) []string {
// 	re := regexp.MustCompile(`\(|\)|AND|OR|[<>=!]+|\w+|'[^']*'`)
// 	return re.FindAllString(rule, -1)
// }

func tokenize(rule string) []string {
	// First split by regular expression
	re := regexp.MustCompile(`\(|\)|AND|OR|[<>=!]+|\w+|'[^']*'`)
	tokens := re.FindAllString(rule, -1)

	// Process each token to remove quotes if present
	for i, token := range tokens {
		// Check if token starts and ends with single quotes
		if len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\'' {
			// Remove quotes and trim spaces
			tokens[i] = strings.TrimSpace(token[1 : len(token)-1])
		}
	}

	return tokens
}

func parseRule(tokens []string) (*models.Node, []string) {
	if len(tokens) == 0 {
		return nil, tokens
	}

	if tokens[0] == "(" {
		// Remove opening parenthesis
		tokens = tokens[1:]

		// Parse the left side of the expression
		left, remainingTokens := parseRule(tokens)

		// Check if we have more tokens to process
		if len(remainingTokens) == 0 {
			return left, remainingTokens
		}

		// If next token is a closing parenthesis and we have more tokens
		if remainingTokens[0] == ")" {
			tokens = remainingTokens[1:] // Remove closing parenthesis

			// If we have more tokens and the next is an operator (AND/OR)
			if len(tokens) > 0 && (tokens[0] == "AND" || tokens[0] == "OR") {
				op := tokens[0]
				tokens = tokens[1:]
				right, remainingTokens := parseRule(tokens)
				return &models.Node{
					Type:  models.OperatorNode,
					Value: &op,
					Left:  left,
					Right: right,
				}, remainingTokens
			}
			return left, tokens
		}

		// Process operator and right side within parentheses
		op := remainingTokens[0]
		remainingTokens = remainingTokens[1:]
		right, remainingTokens := parseRule(remainingTokens)

		// Remove closing parenthesis
		if len(remainingTokens) > 0 && remainingTokens[0] == ")" {
			remainingTokens = remainingTokens[1:]
		}

		node := &models.Node{
			Type:  models.OperatorNode,
			Value: &op,
			Left:  left,
			Right: right,
		}

		// If we have more tokens and the next is an operator (AND/OR)
		if len(remainingTokens) > 0 && (remainingTokens[0] == "AND" || remainingTokens[0] == "OR") {
			op := remainingTokens[0]
			remainingTokens = remainingTokens[1:]
			right, remainingTokens = parseRule(remainingTokens)
			return &models.Node{
				Type:  models.OperatorNode,
				Value: &op,
				Left:  node,
				Right: right,
			}, remainingTokens
		}

		return node, remainingTokens
	}

	// Handle simple conditions (e.g., "age > 30")
	if len(tokens) >= 3 {
		left := tokens[0]
		op := tokens[1]
		right := tokens[2]
		node := &models.Node{
			Type:  models.ConditionNode,
			Value: &op,
			Left:  &models.Node{Value: &left},
			Right: &models.Node{Value: &right},
		}

		tokens = tokens[3:]

		// If we have more tokens and the next is an operator (AND/OR)
		if len(tokens) > 0 && (tokens[0] == "AND" || tokens[0] == "OR") {
			op := tokens[0]
			tokens = tokens[1:]
			rightNode, remainingTokens := parseRule(tokens)
			return &models.Node{
				Type:  models.OperatorNode,
				Value: &op,
				Left:  node,
				Right: rightNode,
			}, remainingTokens
		}

		return node, tokens
	}

	return nil, tokens
}

func (s *RuleService) CreateRule(rule models.Rule) error {

	rules := tokenize(*rule.Description)

	fmt.Println("rules", rules, len(rules))

	rootNode, _ := parseRule(rules)

	finalRule := models.Rule{
		ID:          primitive.NewObjectID(),
		Name:        rule.Name,
		Description: rule.Description,
		RootNode:    rootNode,
	}

	return s.repo.Create(&finalRule)
}

// func (s *RuleService) resolveReferences(node *models.Node) error {
// 	if node == nil {
// 		return nil
// 	}

// 	if node.Type == models.ReferenceNode {
// 		referencedRule, err := s.repo.FindByID(node.ReferenceID)
// 		if err != nil {
// 			return err
// 		}
// 		*node = *referencedRule.RootNode
// 	}

// 	if err := s.resolveReferences(node.Left); err != nil {
// 		return err
// 	}
// 	if err := s.resolveReferences(node.Right); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *RuleService) GetRule(id string) (*models.Rule, error) {

// 	return s.repo.FindByID(id)
// }

func (s *RuleService) GetRule(id string) (*models.Rule, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)

	rule, err := s.repo.FindByID(objectID)
	if err != nil {
		return nil, err
	}

	// Resolve references
	// s.resolveReferences(rule.RootNode)

	return rule, nil
}

func (s *RuleService) GetRules() ([]*models.Rule, error) {

	rules, err := s.repo.Find()

	if err != nil {
		return nil, err
	}

	return rules, nil
}

// func (s *RuleService) CombineRules(name, description string, ruleIDs []primitive.ObjectID) (*models.Rule, error) {

// 	rules, err := s.repo.FindByIDs(ruleIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(rules) < 2 {
// 		return nil, fmt.Errorf("at least two rules are required for combination")
// 	}

// 	and := "AND"

// 	combinedNode := &models.Node{
// 		Type:  models.OperatorNode,
// 		Value: &and,
// 		Left:  &models.Node{Type: models.ReferenceNode, ReferenceID: rules[0].ID},
// 		Right: &models.Node{Type: models.ReferenceNode, ReferenceID: rules[1].ID},
// 	}

// 	for i := 2; i < len(rules); i++ {
// 		combinedNode = &models.Node{
// 			Type:  models.OperatorNode,
// 			Value: &and,
// 			Left:  combinedNode,
// 			Right: &models.Node{Type: models.ReferenceNode, ReferenceID: rules[i].ID},
// 		}
// 	}

// 	newRule := &models.Rule{
// 		Name:        &name,
// 		Description: &description,
// 		RootNode:    combinedNode,
// 	}

// 	err = s.repo.Create(newRule)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return newRule, nil
// }

// func traverseAndCount(node *models.Node, frequency map[string]int, seen map[string]bool) {
// 	if node == nil {
// 		return
// 	}

// 	if node.Type == models.ConditionNode && node.Value != nil {
// 		if !seen[*node.Value] {
// 			frequency[*node.Value]++
// 			seen[*node.Value] = true
// 		}
// 	}

// 	traverseAndCount(node.Left, frequency, seen)
// 	traverseAndCount(node.Right, frequency, seen)
// }

func getCompleteCondition(node *models.Node) string {
	if node.Left.Value == nil || node.Value == nil || node.Right.Value == nil {
		return ""
	}

	conditionString := *node.Left.Value + " " + *node.Value + " " + *node.Right.Value

	return conditionString
}

func traverseAndCount(node *models.Node, frequency map[string]int, seen map[string]bool) {
	if node == nil {
		return
	}

	if node.Type == models.ConditionNode {
		condition := getCompleteCondition(node)

		fmt.Println("condition", condition)

		if condition != "" && !seen[condition] {
			frequency[condition]++
			seen[condition] = true
		}
	}

	traverseAndCount(node.Left, frequency, seen)
	traverseAndCount(node.Right, frequency, seen)
}

// func findCommonNode(rules []*models.Rule) map[string]int {

// 	frequency := make(map[string]int)

// 	for _, rule := range rules {
// 		seen := make(map[string]bool)
// 		traverseAndCount(rule.RootNode, frequency, seen)
// 	}

// 	for _, rule := range rules {
// 		fmt.Println("Rule:", *rule) // Dereference the pointer to get the actual value
// 	}

// 	fmt.Println("F", frequency)

// 	return frequency
// }

func findCommonNodes(rules []*models.Rule) map[string]int {
	frequency := make(map[string]int)

	for _, rule := range rules {
		seen := make(map[string]bool)
		traverseAndCount(rule.RootNode, frequency, seen)
	}

	return frequency
}

// func hasCondition(node *models.Node, condition string) bool {
// 	if node == nil {
// 		return false
// 	}

// 	if node.Type == models.ConditionNode && node.Value != nil && *node.Value == condition {
// 		return true
// 	}

// 	return hasCondition(node.Left, condition) || hasCondition(node.Right, condition)
// }

func hasCondition(node *models.Node, condition string) bool {
	if node == nil || node.Left.Value == nil || node.Right.Value == nil || node.Value == nil {
		return false
	}

	if node.Type == models.ConditionNode {
		// Compare complete conditions
		conditionString := getCompleteCondition(node)
		return conditionString == condition
	}

	return hasCondition(node.Left, condition) || hasCondition(node.Right, condition)
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func buildCommonSubtree(rules []models.Rule, commonCond string, maxFreq int) *models.Node {
	commonNode := &models.Node{
		Type:  models.OperatorNode,
		Value: strPtr("AND"),
	}
	components := strings.Split(commonCond, " ")

	// Add the common condition
	conditionNodeLeftNode := &models.Node{
		Value: &components[0],
	}
	conditionNodeRightNode := &models.Node{
		Value: &components[2],
	}
	conditionNode := &models.Node{
		Type:  models.ConditionNode,
		Value: &components[1],
		Left:  conditionNodeLeftNode,
		Right: conditionNodeRightNode,
	}
	commonNode.Left = conditionNode

	// Create OR node for remaining conditions
	remainingNode := &models.Node{
		Type:  models.OperatorNode,
		Value: strPtr("OR"),
	}

	// Add remaining conditions for each rule
	for _, rule := range rules {
		var ruleBranch *models.Node
		if maxFreq > 1 {

			ruleBranch = removeCommonCondition(copyNode(rule.RootNode), commonCond)
		} else {

			ruleBranch = copyNode(rule.RootNode)
		}

		if remainingNode.Left == nil {
			remainingNode.Left = ruleBranch
		} else if remainingNode.Right == nil {
			remainingNode.Right = ruleBranch
		} else {
			// Create new OR node if needed
			newOR := &models.Node{
				Type:  models.OperatorNode,
				Value: strPtr("OR"),
				Left:  remainingNode.Right,
				Right: ruleBranch,
			}
			remainingNode.Right = newOR
		}

		// Track rule IDs for visualization
		addRuleID(ruleBranch, rule.ID)
	}

	commonNode.Right = remainingNode
	return commonNode
}

func copyNode(node *models.Node) *models.Node {
	if node == nil {
		return nil
	}

	newNode := &models.Node{
		Type:     node.Type,
		Value:    node.Value,
		Left:     copyNode(node.Left),
		Right:    copyNode(node.Right),
		RuleIDs:  append([]primitive.ObjectID{}, node.RuleIDs...),
		Metadata: make(map[string]interface{}),
	}

	for k, v := range node.Metadata {
		newNode.Metadata[k] = v
	}

	return newNode
}

func removeCommonCondition(node *models.Node, condition string) *models.Node {
	if node == nil {
		return nil
	}

	conditionStringComponent := strings.Split(condition, " ")

	if node.Type == models.ConditionNode &&
		node.Value != nil && *node.Value == conditionStringComponent[1] &&
		*node.Left.Value == conditionStringComponent[0] &&
		*node.Right.Value == conditionStringComponent[2] {
		fmt.Println("Hello")
		return nil
	}

	if node.Type == models.OperatorNode {
		node.Left = removeCommonCondition(node.Left, condition)
		node.Right = removeCommonCondition(node.Right, condition)

		// If one child is nil after removal, return the other child
		if node.Left == nil {
			return node.Right
		}
		if node.Right == nil {
			return node.Left
		}
	}

	return node
}
func addRuleID(node *models.Node, ruleID primitive.ObjectID) {
	if node == nil {
		return
	}

	if node.RuleIDs == nil {
		node.RuleIDs = []primitive.ObjectID{}
	}
	node.RuleIDs = append(node.RuleIDs, ruleID)

	addRuleID(node.Left, ruleID)
	addRuleID(node.Right, ruleID)
}

func buildIndependentSubtree(rules []models.Rule) *models.Node {
	if len(rules) == 0 {
		return nil
	}

	if len(rules) == 1 {
		node := copyNode(rules[0].RootNode)
		addRuleID(node, rules[0].ID)
		return node
	}

	// Create OR node for multiple independent rules
	orNode := &models.Node{
		Type:  models.OperatorNode,
		Value: strPtr("OR"),
	}

	for i, rule := range rules {
		node := copyNode(rule.RootNode)
		addRuleID(node, rule.ID)

		if i == 0 {
			orNode.Left = node
		} else if i == 1 {
			orNode.Right = node
		} else {
			// Create new OR node if needed
			newOR := &models.Node{
				Type:  models.OperatorNode,
				Value: strPtr("OR"),
				Left:  orNode.Right,
				Right: node,
			}
			orNode.Right = newOR
		}
	}

	return orNode
}

func (s *RuleService) CombinedRules(name, description string, ruleIDs []primitive.ObjectID) error {

	rules, err := s.repo.FindByIDs(ruleIDs)
	if err != nil {
		return err
	}

	fmt.Println("rules:", rules)

	if len(rules) < 2 {
		return fmt.Errorf("at least two rules are required for combination")
	}

	// Find common conditions
	frequency := findCommonNodes(rules)

	fmt.Printf("%p\n", frequency)

	// Find most common conditionpackage controllers

	var mostCommon string
	maxFreq := 0
	for cond, freq := range frequency {
		if freq > maxFreq {
			maxFreq = freq
			mostCommon = cond
		}
	}

	fmt.Println("maxFreq and cond", mostCommon, maxFreq)

	var rulesWithCommon []models.Rule
	var rulesWithout []models.Rule

	for _, rule := range rules {
		if hasCondition(rule.RootNode, mostCommon) {
			rulesWithCommon = append(rulesWithCommon, *rule)
		} else {
			rulesWithout = append(rulesWithout, *rule)
		}
	}

	root := &models.Node{
		Type:  models.OperatorNode,
		Value: strPtr("OR"),
	}

	// Add metadata for visualization
	root.Metadata = map[string]interface{}{
		"combined":    true,
		"rules_count": len(rules),
	}

	// Build common condition subtree
	if len(rulesWithCommon) > 0 {
		commonSubtree := buildCommonSubtree(rulesWithCommon, mostCommon, maxFreq)
		root.Left = commonSubtree
	}

	if len(rulesWithout) > 0 && len(rulesWithCommon) > 0 {
		independentSubtree := buildIndependentSubtree(rulesWithout)
		if root.Left == nil {
			root.Left = independentSubtree
		} else {
			root.Right = independentSubtree
		}
	}

	if len(rulesWithCommon) == 0 {
		combinedNode := &models.Node{
			Type:  models.OperatorNode,
			Value: strPtr("OR"),
			Left:  &models.Node{Type: models.ReferenceNode, ReferenceID: rules[0].ID},
			Right: &models.Node{Type: models.ReferenceNode, ReferenceID: rules[1].ID},
		}

		// Create a root node and attach rules based on odd/even index
		root := combinedNode // Start with the combined node of the first two rules

		for i := 2; i < len(rules); i++ {
			newNode := &models.Node{
				Type:        models.ReferenceNode,
				ReferenceID: rules[i].ID,
			}

			if i%2 == 0 {
				// Even indexed rule: attach to the right of root
				root = &models.Node{
					Type:  models.OperatorNode,
					Value: strPtr("OR"),
					Left:  root,
					Right: newNode,
				}
			} else {
				// Odd indexed rule: attach to the left of root
				root = &models.Node{
					Type:  models.OperatorNode,
					Value: strPtr("OR"),
					Left:  newNode,
					Right: root,
				}
			}
		}
	}

	combinedRule := &models.Rule{
		ID:          primitive.NewObjectID(),
		Name:        &name,
		Description: &description,
		RootNode:    root,
	}

	//if all rules have the same condition
	if len(rulesWithCommon) == len(rules) {
		combinedRule.RootNode = root.Left
	}

	fmt.Println("combinedRule", combinedRule)

	err = s.repo.Create(combinedRule)
	return err

}
