package models

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UUIDArray []uuid.UUID

func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	parts := make([]string, 0, len(a))
	for _, id := range a {
		parts = append(parts, id.String())
	}

	return "{" + strings.Join(parts, ",") + "}", nil
}

func (UUIDArray) GormDataType() string {
	return "uuid[]"
}

func (a UUIDArray) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	v, err := a.Value()
	if err != nil {
		panic(fmt.Sprintf("UUIDArray.Value() failed in GormValue: %v", err))
	}
	return clause.Expr{SQL: "?::uuid[]", Vars: []interface{}{v}}
}

func (a *UUIDArray) Scan(src interface{}) error {
	if src == nil {
		*a = UUIDArray{}
		return nil
	}

	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("unsupported UUID array source type: %T", src)
	}

	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")

	if raw == "" {
		*a = UUIDArray{}
		return nil
	}

	strs := strings.Split(raw, ",")

	res := make([]uuid.UUID, 0, len(strs))
	for _, s := range strs {
		s = strings.Trim(s, `"`)
		id, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid UUID in array: %w", err)
		}
		res = append(res, id)
	}

	*a = UUIDArray(res)
	return nil
}

type ProjectSlot struct {
	ID                uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID         uuid.UUID  `gorm:"column:project_id;type:uuid;not null;uniqueIndex:idx_project_user,where:project_id IS NOT NULL" json:"project_id"`
	PrimarySkillsID   UUIDArray  `gorm:"column:primary_skills_id;type:uuid[];not null" json:"primary_skills_id"`
	SecondarySkillsID UUIDArray  `gorm:"column:secondary_skills_id;type:uuid[];not null" json:"secondary_skills_id"`
	UserID            *uuid.UUID `gorm:"column:user_id;type:uuid;uniqueIndex:idx_project_user,where:user_id IS NOT NULL" json:"user_id"`
	Title             string     `gorm:"column:title; not null" json:"title"`
	Description       *string    `gorm:"column:description" json:"description"`
	Status            string     `gorm:"column:status; not null" json:"status"`
	CreatedAt         time.Time  `gorm:"column:created_at; not null" json:"created_at"`

	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
	User    User    `gorm:"foreignKey:UserID" json:"-"`
}

func (ps *ProjectSlot) TableName() string { return "Project_Slot" }
