package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	getRuleByID      string = `SELECT id, user_id, title, trigger, condition, rule_action, in_effect FROM rules WHERE id = $1;`
	getRulesByUserID string = `SELECT id, user_id, title, trigger, condition, rule_action, in_effect FROM rules WHERE user_id = $1;`
)

type Rule struct {
	ID         uuid.UUID
	User       User
	Title      string
	Trigger    string
	Condition  Condition
	RuleAction Action
	InEffect   bool
}

func (r *Rule) unmarshalRow(row pgx.Row) error {
	return row.Scan(&r.ID, &r.User.ID, &r.Title, &r.Trigger, &r.Condition, &r.RuleAction, &r.InEffect)
}

func (db *DB) GetRuleByID(id uuid.UUID) (*Rule, error) {
	row := db.conn.QueryRow(context.Background(), getRuleByID, id)
	var rule Rule
	err := row.Scan(&rule.ID, &rule.User.ID, &rule.Title, &rule.Trigger, &rule.Condition, &rule.RuleAction)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}
