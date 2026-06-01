package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
)

type ImportReplaceRuleApi struct {
	importReplaceRules *services.ImportReplaceRuleService
}

var ImportReplaceRules = &ImportReplaceRuleApi{
	importReplaceRules: services.ImportReplaceRules,
}

func (a *ImportReplaceRuleApi) ImportReplaceRuleSetListHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	ruleSets, err := a.importReplaceRules.GetAllRuleSets(c, uid)

	if err != nil {
		log.Errorf(c, "[import_replace_rules.ListHandler] failed to get rule sets for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	ruleSetInfos := make([]*models.ImportReplaceRuleSetInfoResponse, len(ruleSets))

	for i, ruleSet := range ruleSets {
		count, err := a.importReplaceRules.CountRulesByRuleSetId(c, uid, ruleSet.RuleSetId)

		if err != nil {
			log.Errorf(c, "[import_replace_rules.ListHandler] failed to count rules for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		ruleSetInfos[i] = ruleSet.ToInfoResponse(int(count))
	}

	return ruleSetInfos, nil
}

func (a *ImportReplaceRuleApi) ImportReplaceRuleSetGetHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	var request models.ImportReplaceRuleSetGetRequest

	err := c.ShouldBindQuery(&request)

	if err != nil {
		log.Warnf(c, "[import_replace_rules.GetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	ruleSet, err := a.importReplaceRules.GetRuleSetById(c, uid, request.Id)

	if err != nil {
		log.Errorf(c, "[import_replace_rules.GetHandler] failed to get rule set \"id:%d\" for user \"uid:%d\", because %s", request.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	rules, err := a.importReplaceRules.GetAllRulesByRuleSetId(c, uid, ruleSet.RuleSetId)

	if err != nil {
		log.Errorf(c, "[import_replace_rules.GetHandler] failed to get rules for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return ruleSet.ToDetailResponse(rules), nil
}

func (a *ImportReplaceRuleApi) ImportReplaceRuleSetCreateHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	var request models.ImportReplaceRuleSetCreateRequest

	err := c.ShouldBind(&request)

	if err != nil {
		log.Warnf(c, "[import_replace_rules.CreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	ruleSet, err := a.importReplaceRules.CreateRuleSet(c, uid, request.Name, request.FileType, request.Rules)

	if err != nil {
		log.Errorf(c, "[import_replace_rules.CreateHandler] failed to create rule set for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return ruleSet.ToInfoResponse(len(request.Rules)), nil
}

func (a *ImportReplaceRuleApi) ImportReplaceRuleSetModifyHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	var request models.ImportReplaceRuleSetModifyRequest

	err := c.ShouldBind(&request)

	if err != nil {
		log.Warnf(c, "[import_replace_rules.ModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	err = a.importReplaceRules.ModifyRuleSet(c, uid, request.Id, request.Name, request.FileType, request.Rules)

	if err != nil {
		log.Errorf(c, "[import_replace_rules.ModifyHandler] failed to modify rule set \"id:%d\" for user \"uid:%d\", because %s", request.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

func (a *ImportReplaceRuleApi) ImportReplaceRuleSetGetDefaultByFileTypeHandler(c *core.WebContext) (any, *errs.Error) {
        uid := c.GetCurrentUid()
        var request models.ImportReplaceRuleSetGetDefaultByFileTypeRequest

        err := c.ShouldBindQuery(&request)

        if err != nil {
                log.Warnf(c, "[import_replace_rules.GetDefaultByFileTypeHandler] parse request failed, because %s", err.Error())
                return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
        }

        ruleSet, err := a.importReplaceRules.GetDefaultByFileType(c, uid, request.FileType)

        if err != nil {
                log.Errorf(c, "[import_replace_rules.GetDefaultByFileTypeHandler] failed to get default rule set for user \"uid:%d\", because %s", uid, err.Error())
                return nil, errs.Or(err, errs.ErrOperationFailed)
        }

        if ruleSet == nil {
                return nil, nil
        }

        rules, err := a.importReplaceRules.GetAllRulesByRuleSetId(c, uid, ruleSet.RuleSetId)

        if err != nil {
                log.Errorf(c, "[import_replace_rules.GetDefaultByFileTypeHandler] failed to get rules for user \"uid:%d\", because %s", uid, err.Error())
                return nil, errs.Or(err, errs.ErrOperationFailed)
        }

        return ruleSet.ToDetailResponse(rules), nil
}

func (a *ImportReplaceRuleApi) ImportReplaceRuleSetSetDefaultHandler(c *core.WebContext) (any, *errs.Error) {
        uid := c.GetCurrentUid()
        var request models.ImportReplaceRuleSetSetDefaultRequest

        err := c.ShouldBind(&request)

        if err != nil {
                log.Warnf(c, "[import_replace_rules.SetDefaultHandler] parse request failed, because %s", err.Error())
                return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
        }

        err = a.importReplaceRules.SetDefaultRuleSet(c, uid, request.Id, request.FileType)

        if err != nil {
                log.Errorf(c, "[import_replace_rules.SetDefaultHandler] failed to set default for user \"uid:%d\", because %s", uid, err.Error())
                return nil, errs.Or(err, errs.ErrOperationFailed)
        }

        return true, nil
}

func (a *ImportReplaceRuleApi) ImportReplaceRuleSetDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	var request models.ImportReplaceRuleSetDeleteRequest

	err := c.ShouldBind(&request)

	if err != nil {
		log.Warnf(c, "[import_replace_rules.DeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	err = a.importReplaceRules.DeleteRuleSet(c, uid, request.Id)

	if err != nil {
		log.Errorf(c, "[import_replace_rules.DeleteHandler] failed to delete rule set \"id:%d\" for user \"uid:%d\", because %s", request.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}
