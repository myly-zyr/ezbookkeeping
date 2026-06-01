export interface ImportReplaceRuleItem {
    type: string;
    sourceValue: string;
    targetId: string;
}

export interface ImportReplaceRuleSetInfoResponse {
    id: string;
    name: string;
    ruleCount: number;
    createdTime: number;
    updatedTime: number;
}

export interface ImportReplaceRuleResponse {
    type: string;
    sourceValue: string;
    targetId: string;
}

export interface ImportReplaceRuleSetDetailResponse {
    id: string;
    name: string;
    rules: ImportReplaceRuleResponse[];
    createdTime: number;
    updatedTime: number;
}

export interface ImportReplaceRuleSetCreateRequest {
    name: string;
    fileType: string;
    rules: ImportReplaceRuleItem[];
}

export interface ImportReplaceRuleSetModifyRequest {
    id: string;
    name: string;
    fileType: string;
    rules: ImportReplaceRuleItem[];
}

export interface ImportReplaceRuleSetDeleteRequest {
    id: string;
}
