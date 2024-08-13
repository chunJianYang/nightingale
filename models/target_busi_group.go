package models

import (
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm/clause"
)

type TargetBusiGroup struct {
	Id          int64  `json:"id" gorm:"primaryKey;type:bigint;autoIncrement"`
	TargetIdent string `json:"target_ident" gorm:"type:varchar(191);not null;index:idx_target_group,unique,priority:1"`
	GroupId     int64  `json:"group_id" gorm:"type:bigint;not null;index:idx_target_group,unique,priority:2"`
	UpdateAt    int64  `json:"update_at" gorm:"type:bigint;not null"`
}

func (t *TargetBusiGroup) TableName() string {
	return "target_busi_group"
}

func TargetBusiGroupsGetAll(ctx *ctx.Context) (map[string][]*TargetBusiGroup, error) {
	var lst []*TargetBusiGroup
	err := DB(ctx).Find(&lst).Error
	if err != nil {
		return nil, err
	}
	tgs := make(map[string][]*TargetBusiGroup)
	for _, tg := range lst {
		tgs[tg.TargetIdent] = append(tgs[tg.TargetIdent], tg)
	}
	return tgs, nil
}

func TargetBindBgids(ctx *ctx.Context, idents []string, bgids []int64) error {
	lst := make([]TargetBusiGroup, 0, len(bgids)*len(idents))
	updateAt := time.Now().Unix()
	for _, bgid := range bgids {
		for _, ident := range idents {
			cur := TargetBusiGroup{
				TargetIdent: ident,
				GroupId:     bgid,
				UpdateAt:    updateAt,
			}
			lst = append(lst, cur)
		}
	}

	return DB(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&lst).Error
}

func TargetUnbindBgids(ctx *ctx.Context, idents []string, bgids []int64) error {
	return DB(ctx).Where("target_ident in ? and group_id in ?",
		idents, bgids).Delete(&TargetBusiGroup{}).Error
}
