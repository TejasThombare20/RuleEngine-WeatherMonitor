package repositories

import (
	"context"

	"github.com/TejasThombare20/rule-engine/config"
	"github.com/TejasThombare20/rule-engine/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type RuleRepository struct {
	collection *mongo.Collection
}

func NewRuleRepository() *RuleRepository {
	return &RuleRepository{
		collection: config.RuleCollection,
	}
}

func (ruleCollection *RuleRepository) Create(rule models.Rule) error {

	_, err := ruleCollection.collection.InsertOne(context.Background(), rule)
	return err
}
