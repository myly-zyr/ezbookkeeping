package errs

import "net/http"

var (
	ErrImportReplaceRuleSetIdInvalid   = NewNormalError(NormalSubcategoryImportReplaceRule, 0, http.StatusBadRequest, "import replace rule set id is invalid")
	ErrImportReplaceRuleSetNotFound    = NewNormalError(NormalSubcategoryImportReplaceRule, 1, http.StatusBadRequest, "import replace rule set not found")
	ErrImportReplaceRuleSetDataInvalid = NewNormalError(NormalSubcategoryImportReplaceRule, 2, http.StatusBadRequest, "import replace rule set data is invalid")
)
