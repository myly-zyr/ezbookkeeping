package models

// ImportReplaceRuleSet represents a set of import replace rules
type ImportReplaceRuleSet struct {
	RuleSetId       int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_import_replace_rule_set_uid_deleted) NOT NULL"`
	Deleted         bool   `xorm:"INDEX(IDX_import_replace_rule_set_uid_deleted) NOT NULL"`
	Name            string `xorm:"VARCHAR(128) NOT NULL"`
	FileType        string `xorm:"VARCHAR(32) NOT NULL"`
	IsDefault       bool   `xorm:"NOT NULL DEFAULT false"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// ImportReplaceRule represents a single import replace rule
type ImportReplaceRule struct {
	RuleId      int64  `xorm:"PK"`
	RuleSetId   int64  `xorm:"INDEX(IDX_import_replace_rule_rule_set_id) NOT NULL"`
	DataType    string `xorm:"VARCHAR(32) NOT NULL"`
	SourceValue string `xorm:"VARCHAR(255) NOT NULL"`
	TargetId    string `xorm:"VARCHAR(255) NOT NULL"`
	SortOrder   int64  `xorm:"NOT NULL"`
}

// ImportReplaceRuleSetCreateRequest represents all parameters of rule set creation request
type ImportReplaceRuleSetCreateRequest struct {
	Name     string                   `json:"name" binding:"required,notBlank,max=128"`
	FileType string                   `json:"fileType" binding:"required,max=32"`
	Rules    []*ImportReplaceRuleItem `json:"rules" binding:"required,min=1"`
}

// ImportReplaceRuleSetModifyRequest represents all parameters of rule set modification request
type ImportReplaceRuleSetModifyRequest struct {
	Id       int64                    `json:"id,string" binding:"required,min=1"`
	Name     string                   `json:"name" binding:"required,notBlank,max=128"`
	FileType string                   `json:"fileType" binding:"required,max=32"`
	Rules    []*ImportReplaceRuleItem `json:"rules" binding:"required,min=1"`
}

// ImportReplaceRuleSetGetDefaultByFileTypeRequest represents request to get default rule set by file type
type ImportReplaceRuleSetGetDefaultByFileTypeRequest struct {
        FileType string `form:"fileType" binding:"required,max=32"`
}

// ImportReplaceRuleSetGetRequest represents all parameters of rule set getting request
type ImportReplaceRuleSetGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// ImportReplaceRuleSetSetDefaultRequest represents request to set default rule set
type ImportReplaceRuleSetSetDefaultRequest struct {
        Id       int64  `json:"id,string" binding:"required,min=1"`
        FileType string `json:"fileType" binding:"required,max=32"`
}

// ImportReplaceRuleSetDeleteRequest represents all parameters of rule set deleting request
type ImportReplaceRuleSetDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// ImportReplaceRuleItem represents a single rule item in request
type ImportReplaceRuleItem struct {
	DataType    string `json:"type" binding:"required"`
	SourceValue string `json:"sourceValue" binding:"required"`
	TargetId    string `json:"targetId" binding:"required"`
}

// ImportReplaceRuleSetInfoResponse represents a view-object of rule set info (without rules)
type ImportReplaceRuleSetInfoResponse struct {
	Id              int64  `json:"id,string"`
	Name            string `json:"name"`
	RuleCount       int    `json:"ruleCount"`
	CreatedUnixTime int64  `json:"createdTime"`
	UpdatedUnixTime int64  `json:"updatedTime"`
}

// ImportReplaceRuleSetDetailResponse represents a view-object of rule set with all rules
type ImportReplaceRuleSetDetailResponse struct {
	Id              int64                        `json:"id,string"`
	Name            string                       `json:"name"`
	Rules           []*ImportReplaceRuleResponse `json:"rules"`
	CreatedUnixTime int64                        `json:"createdTime"`
	UpdatedUnixTime int64                        `json:"updatedTime"`
}

// ImportReplaceRuleResponse represents a view-object of a single rule
type ImportReplaceRuleResponse struct {
	DataType    string `json:"type"`
	SourceValue string `json:"sourceValue"`
	TargetId    string `json:"targetId"`
}

// ToInfoResponse converts a rule set model to info response
func (rs *ImportReplaceRuleSet) ToInfoResponse(ruleCount int) *ImportReplaceRuleSetInfoResponse {
	return &ImportReplaceRuleSetInfoResponse{
		Id:              rs.RuleSetId,
		Name:            rs.Name,
		RuleCount:       ruleCount,
		CreatedUnixTime: rs.CreatedUnixTime,
		UpdatedUnixTime: rs.UpdatedUnixTime,
	}
}

// ToDetailResponse converts a rule set model to detail response
func (rs *ImportReplaceRuleSet) ToDetailResponse(rules []*ImportReplaceRule) *ImportReplaceRuleSetDetailResponse {
	ruleResps := make([]*ImportReplaceRuleResponse, len(rules))

	for i, rule := range rules {
		ruleResps[i] = &ImportReplaceRuleResponse{
			DataType:    rule.DataType,
			SourceValue: rule.SourceValue,
			TargetId:    rule.TargetId,
		}
	}

	return &ImportReplaceRuleSetDetailResponse{
		Id:              rs.RuleSetId,
		Name:            rs.Name,
		Rules:           ruleResps,
		CreatedUnixTime: rs.CreatedUnixTime,
		UpdatedUnixTime: rs.UpdatedUnixTime,
	}
}
