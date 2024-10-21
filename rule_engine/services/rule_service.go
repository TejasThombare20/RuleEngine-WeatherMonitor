package services

import (
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

	return s.repo.Create(finalRule)
}
