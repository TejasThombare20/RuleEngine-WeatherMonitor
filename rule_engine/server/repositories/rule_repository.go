package repositories

import (
	"context"
	"fmt"

	"github.com/TejasThombare20/rule-engine/config"
	"github.com/TejasThombare20/rule-engine/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (ruleCollection *RuleRepository) Create(rule *models.Rule) error {

	_, err := ruleCollection.collection.InsertOne(context.Background(), rule)
	return err
}

func (r *RuleRepository) FindByID(id primitive.ObjectID) (*models.Rule, error) {
	// objectID, _ := primitive.ObjectIDFromHex(id)
	var rule models.Rule
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&rule)

	fmt.Println("rule", &rule, err)
	return &rule, err
}

func (r *RuleRepository) FindByIDs(ids []primitive.ObjectID) ([]*models.Rule, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rules []*models.Rule
	if err = cursor.All(context.Background(), &rules); err != nil {
		return nil, err
	}
	return rules, nil
}
