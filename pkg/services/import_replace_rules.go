package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

type ImportReplaceRuleService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

var ImportReplaceRules = &ImportReplaceRuleService{
	ServiceUsingDB: ServiceUsingDB{
		container: datastore.Container,
	},
	ServiceUsingUuid: ServiceUsingUuid{
		container: uuid.Container,
	},
}

func (s *ImportReplaceRuleService) GetAllRuleSets(c core.Context, uid int64) ([]*models.ImportReplaceRuleSet, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var ruleSets []*models.ImportReplaceRuleSet
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Find(&ruleSets)

	if err != nil {
		return nil, err
	}

	return ruleSets, nil
}

func (s *ImportReplaceRuleService) GetRuleSetById(c core.Context, uid int64, id int64) (*models.ImportReplaceRuleSet, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}
	if id <= 0 {
		return nil, errs.ErrImportReplaceRuleSetIdInvalid
	}

	ruleSet := &models.ImportReplaceRuleSet{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(id).Where("uid=? AND deleted=?", uid, false).Get(ruleSet)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrImportReplaceRuleSetNotFound
	}

	return ruleSet, nil
}

func (s *ImportReplaceRuleService) GetAllRulesByRuleSetId(c core.Context, uid int64, ruleSetId int64) ([]*models.ImportReplaceRule, error) {
	var rules []*models.ImportReplaceRule
	err := s.UserDataDB(uid).NewSession(c).Where("rule_set_id=?", ruleSetId).OrderBy("sort_order ASC").Find(&rules)

	if err != nil {
		return nil, err
	}

	return rules, nil
}

func (s *ImportReplaceRuleService) CountRulesByRuleSetId(c core.Context, uid int64, ruleSetId int64) (int64, error) {
	count, err := s.UserDataDB(uid).NewSession(c).Where("rule_set_id=?", ruleSetId).Count(&models.ImportReplaceRule{})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ImportReplaceRuleService) CreateRuleSet(c core.Context, uid int64, name string, fileType string, items []*models.ImportReplaceRuleItem) (*models.ImportReplaceRuleSet, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	ruleSet := &models.ImportReplaceRuleSet{
		RuleSetId:       s.GenerateUuid(uuid.UUID_TYPE_IMPORT_REPLACE_RULE_SET),
		Uid:             uid,
		Deleted:         false,
		Name:            name,
		CreatedUnixTime: now,
		UpdatedUnixTime: now,
	}

	err := s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(ruleSet)

		if err != nil {
			return err
		}

		for i, item := range items {
			rule := &models.ImportReplaceRule{
				RuleId:      s.GenerateUuid(uuid.UUID_TYPE_DEFAULT),
				RuleSetId:   ruleSet.RuleSetId,
				DataType:    item.DataType,
				SourceValue: item.SourceValue,
				TargetId:    item.TargetId,
				SortOrder:   int64(i),
			}

			_, err := sess.Insert(rule)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ruleSet, nil
}

func (s *ImportReplaceRuleService) ModifyRuleSet(c core.Context, uid int64, id int64, name string, fileType string, items []*models.ImportReplaceRuleItem) error {
	ruleSet, err := s.GetRuleSetById(c, uid, id)

	if err != nil {
		return err
	}

	now := time.Now().Unix()

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		ruleSet.Name = name
		ruleSet.UpdatedUnixTime = now

		_, err := sess.ID(ruleSet.RuleSetId).Cols("name", "updated_unix_time").Update(ruleSet)

		if err != nil {
			return err
		}

		_, err = sess.Where("rule_set_id=?", ruleSet.RuleSetId).Delete(&models.ImportReplaceRule{})

		if err != nil {
			return err
		}

		for i, item := range items {
			rule := &models.ImportReplaceRule{
				RuleId:      s.GenerateUuid(uuid.UUID_TYPE_DEFAULT),
				RuleSetId:   ruleSet.RuleSetId,
				DataType:    item.DataType,
				SourceValue: item.SourceValue,
				TargetId:    item.TargetId,
				SortOrder:   int64(i),
			}

			_, err := sess.Insert(rule)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ImportReplaceRuleService) GetDefaultByFileType(c core.Context, uid int64, fileType string) (*models.ImportReplaceRuleSet, error) {
        if uid <= 0 {
                return nil, errs.ErrUserIdInvalid
        }

        ruleSet := &models.ImportReplaceRuleSet{}
        has, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND file_type=? AND is_default=?", uid, false, fileType, true).Get(ruleSet)

        if err != nil {
                return nil, err
        } else if !has {
                return nil, nil
        }

        return ruleSet, nil
}

func (s *ImportReplaceRuleService) SetDefaultRuleSet(c core.Context, uid int64, id int64, fileType string) error {
        if uid <= 0 {
                return errs.ErrUserIdInvalid
        }

        return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
                // Clear all defaults for this file type
                _, err := sess.Where("uid=? AND deleted=? AND file_type=?", uid, false, fileType).Cols("is_default").Update(&models.ImportReplaceRuleSet{IsDefault: false})
                if err != nil {
                        return err
                }

                // Set new default
                _, err = sess.ID(id).Where("uid=? AND deleted=?", uid, false).Cols("is_default").Update(&models.ImportReplaceRuleSet{IsDefault: true})
                return err
        })
}

func (s *ImportReplaceRuleService) DeleteRuleSet(c core.Context, uid int64, id int64) error {
	ruleSet, err := s.GetRuleSetById(c, uid, id)

	if err != nil {
		return err
	}

	now := time.Now().Unix()

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Where("rule_set_id=?", ruleSet.RuleSetId).Delete(&models.ImportReplaceRule{})

		if err != nil {
			return err
		}

		_, err = sess.ID(ruleSet.RuleSetId).Cols("deleted", "deleted_unix_time").Update(&models.ImportReplaceRuleSet{
			Deleted:         true,
			DeletedUnixTime: now,
		})

		if err != nil {
			return err
		}

		return nil
	})
}
