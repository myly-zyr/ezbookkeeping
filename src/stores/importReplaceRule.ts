import { ref } from 'vue';
import { defineStore } from 'pinia';

import {
    type ImportReplaceRuleSetInfoResponse,
    type ImportReplaceRuleSetDetailResponse,
    type ImportReplaceRuleSetCreateRequest,
    type ImportReplaceRuleSetModifyRequest,
    type ImportReplaceRuleSetDeleteRequest
} from '@/models/import_replace_rule.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useImportReplaceRuleStore = defineStore('importReplaceRule', () => {
    const savedRuleSets = ref<ImportReplaceRuleSetInfoResponse[]>([]);

    function getAllRuleSets(): Promise<ImportReplaceRuleSetInfoResponse[]> {
        return new Promise((resolve, reject) => {
            services.getAllImportReplaceRuleSets().then(response => {
                if (!response.data || !response.data.success || !response.data.result) {
                    reject(new Error('failed to get rule sets'));
                    return;
                }

                savedRuleSets.value = response.data.result;
                resolve(response.data.result);
            }).catch(err => {
                logger.error('failed to get rule sets', err);
                reject(err);
            });
        });
    }

    function getRuleSet(id: string): Promise<ImportReplaceRuleSetDetailResponse> {
        return new Promise((resolve, reject) => {
            services.getImportReplaceRuleSet({ id }).then(response => {
                if (!response.data || !response.data.success || !response.data.result) {
                    reject(new Error('failed to get rule set'));
                    return;
                }

                resolve(response.data.result);
            }).catch(err => {
                logger.error('failed to get rule set', err);
                reject(err);
            });
        });
    }

    function getDefaultByFileType(fileType: string): Promise<ImportReplaceRuleSetDetailResponse | null> {
        return new Promise((resolve, reject) => {
            services.getDefaultImportReplaceRuleSetByFileType({ fileType }).then(response => {
                if (!response.data || !response.data.success) {
                    reject(new Error('failed to get default rule set'));
                    return;
                }

                resolve(response.data.result || null);
            }).catch(err => {
                logger.error('failed to get default rule set', err);
                reject(err);
            });
        });
    }

    function addRuleSet(req: ImportReplaceRuleSetCreateRequest): Promise<ImportReplaceRuleSetInfoResponse> {
        return new Promise((resolve, reject) => {
            services.addImportReplaceRuleSet(req).then(response => {
                if (!response.data || !response.data.success || !response.data.result) {
                    reject(new Error('failed to add rule set'));
                    return;
                }

                resolve(response.data.result);
            }).catch(err => {
                logger.error('failed to add rule set', err);
                reject(err);
            });
        });
    }

    function modifyRuleSet(req: ImportReplaceRuleSetModifyRequest): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.modifyImportReplaceRuleSet(req).then(response => {
                if (!response.data || !response.data.success) {
                    reject(new Error('failed to modify rule set'));
                    return;
                }

                resolve(true);
            }).catch(err => {
                logger.error('failed to modify rule set', err);
                reject(err);
            });
        });
    }

    function setDefaultRuleSet(id: string, fileType: string): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.setDefaultImportReplaceRuleSet({ id, fileType }).then(response => {
                if (!response.data || !response.data.success) {
                    reject(new Error('failed to set default rule set'));
                    return;
                }

                resolve(true);
            }).catch(err => {
                logger.error('failed to set default rule set', err);
                reject(err);
            });
        });
    }

    function deleteRuleSet(req: ImportReplaceRuleSetDeleteRequest): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteImportReplaceRuleSet(req).then(response => {
                if (!response.data || !response.data.success) {
                    reject(new Error('failed to delete rule set'));
                    return;
                }

                resolve(true);
            }).catch(err => {
                logger.error('failed to delete rule set', err);
                reject(err);
            });
        });
    }

    return {
        savedRuleSets,
        getAllRuleSets,
        getRuleSet,
        getDefaultByFileType,
        addRuleSet,
        modifyRuleSet,
        setDefaultRuleSet,
        deleteRuleSet
    };
});
