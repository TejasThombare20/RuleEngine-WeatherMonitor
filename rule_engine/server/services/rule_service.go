package services

import (
	"fmt"
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

func tokenize(rule string) []string {
	re := regexp.MustCompile(`\(|\)|AND|OR|[<>=!]+|\w+|'[^']*'`)
	return re.FindAllString(rule, -1)
}

func parseRule(tokens []string) (*models.Node, []string) {
	if tokens[0] == "(" {
		tokens = tokens[1:] // Remove opening parenthesis
		left, tokens := parseRule(tokens)
		op := tokens[0]
		tokens = tokens[1:]
		right, tokens := parseRule(tokens)
		tokens = tokens[1:] // Remove closing parenthesis
		return &models.Node{Type: models.OperatorNode, Value: &op, Left: left, Right: right}, tokens
	}
	left := tokens[0]
	op := tokens[1]
	right := tokens[2]
	return &models.Node{Type: models.ConditionNode, Value: &op, Left: &models.Node{Value: &left}, Right: &models.Node{Value: &right}}, tokens[3:]
}

func (s *RuleService) CreateRule(rule models.Rule) error {

	rules := tokenize(*rule.Description)
	rootNode, _ := parseRule(rules)

	finalRule := models.Rule{
		ID:          primitive.NewObjectID(),
		Name:        rule.Name,
		Description: rule.Description,
		RootNode:    rootNode,
	}

	return s.repo.Create(&finalRule)
}

func (s *RuleService) resolveReferences(node *models.Node) error {
	if node == nil {
		return nil
	}

	if node.Type == models.ReferenceNode {
		referencedRule, err := s.repo.FindByID(node.ReferenceID)
		if err != nil {
			return err
		}
		*node = *referencedRule.RootNode
	}

	if err := s.resolveReferences(node.Left); err != nil {
		return err
	}
	if err := s.resolveReferences(node.Right); err != nil {
		return err
	}

	return nil
}

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
	s.resolveReferences(rule.RootNode)

	return rule, nil
}

func (s *RuleService) CombineRules(name, description string, ruleIDs []primitive.ObjectID) (*models.Rule, error) {

	rules, err := s.repo.FindByIDs(ruleIDs)
	if err != nil {
		return nil, err
	}

	if len(rules) < 2 {
		return nil, fmt.Errorf("at least two rules are required for combination")
	}

	and := "AND"

	combinedNode := &models.Node{
		Type:  models.OperatorNode,
		Value: &and,
		Left:  &models.Node{Type: models.ReferenceNode, ReferenceID: rules[0].ID},
		Right: &models.Node{Type: models.ReferenceNode, ReferenceID: rules[1].ID},
	}

	for i := 2; i < len(rules); i++ {
		combinedNode = &models.Node{
			Type:  models.OperatorNode,
			Value: &and,
			Left:  combinedNode,
			Right: &models.Node{Type: models.ReferenceNode, ReferenceID: rules[i].ID},
		}
	}

	newRule := &models.Rule{
		Name:        &name,
		Description: &description,
		RootNode:    combinedNode,
	}

	err = s.repo.Create(newRule)
	if err != nil {
		return nil, err
	}

	return newRule, nil
}
