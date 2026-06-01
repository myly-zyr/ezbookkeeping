<template>
    <v-dialog width="1000" :persistent="loading || !!rules.length || !!newRule.targetId" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex flex-wrap align-center justify-center">
                    <div class="d-flex flex-wrap align-center">
                        <h4 class="text-h4 text-wrap">{{ tt('Batch Replace Categories / Accounts / Tags') }}</h4>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :disabled="loading"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                    <v-spacer/>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :icon="true" :disabled="loading">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <v-list-item :prepend-icon="mdiContentSaveOutline"
                                             :title="tt('Save Rule Set')"
                                             @click="showSaveDialog()"></v-list-item>
                                <v-divider/>
                                <v-list-subheader v-if="allRuleSetsList.length">{{ tt('Rule Sets') }}</v-list-subheader>
                                <v-list-item v-for="rs in allRuleSetsList" :key="rs.id"
                                             :prepend-icon="mdiFileDocumentOutline"
                                             :title="rs.name"
                                             :subtitle="rs.ruleCount + ' ' + tt('rules')"
                                             :active="rs.id === currentRuleSetId"
                                             @click="loadOtherRuleSet(rs.id)">
                                    <template #append>
                                        <v-btn v-if="!rs.isDefault" density="compact" variant="text" color="error"
                                               :icon="true" size="small"
                                               @click.stop="deleteRuleSet(rs.id)">
                                            <v-icon :icon="mdiDeleteOutline" size="18"/>
                                            <v-tooltip activator="parent">{{ tt('Delete') }}</v-tooltip>
                                        </v-btn>
                                    </template>
                                </v-list-item>
                                <v-list-item v-if="!allRuleSetsList.length"
                                             :disabled="true"
                                             :title="tt('No saved rule sets')"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <v-row>
                    <v-col cols="12">
                        <v-table density="comfortable" fixed-header fixed-footer height="400" striped="even">
                            <thead>
                                <tr>
                                    <th class="text-left">{{ tt('Type') }}</th>
                                    <th class="text-left">{{ tt('Source Value') }}</th>
                                    <th class="text-left">{{ tt('Target Value') }}</th>
                                    <th class="text-right">{{ tt('Operation') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(rule, index) in rules" :key="index">
                                    <td class="text-left">{{ getRuleTypeDisplayName(rule) }}</td>
                                    <td class="text-left">{{ rule.sourceValue || tt('(Empty)') }}</td>
                                    <td class="text-left">{{ getRuleTargetValueDisplayName(rule) }}</td>
                                    <td class="text-right">
                                        <v-btn density="comfortable" variant="tonal" color="error"
                                               :disabled="loading" @click="removeRule(index)">{{ tt('Delete') }}</v-btn>
                                    </td>
                                </tr>
                            </tbody>
                            <tfoot>
                                <tr style="background-color: rgb(var(--v-theme-surface))">
                                    <td>
                                        <v-select class="w-100" density="compact" variant="underlined"
                                                  item-title="name"
                                                  item-value="value"
                                                  :disabled="loading"
                                                  :items="[
                                                      {
                                                          value: 'expenseCategory',
                                                          name: tt('Expense Category')
                                                      },
                                                      {
                                                          value: 'incomeCategory',
                                                          name: tt('Income Category')
                                                      },
                                                      {
                                                          value: 'transferCategory',
                                                          name: tt('Transfer Category')
                                                      },
                                                      {
                                                          value: 'account',
                                                          name: tt('Account')
                                                      },
                                                      {
                                                          value: 'tag',
                                                          name: tt('Transaction Tag')
                                                      }
                                                  ]"
                                                  v-model="newRule.dataType"
                                                  @update:model-value="newRule.sourceValue = ''; newRule.targetId = ''"
                                        />
                                    </td>
                                    <td>
                                        <v-autocomplete class="w-100" density="compact" variant="underlined"
                                                        item-title="name" item-value="value" persistent-placeholder
                                                        :disabled="loading" :items="sourceItems"
                                                        :no-data-text="noSourceItemText"
                                                        v-model="newRule.sourceValue">
                                        </v-autocomplete>
                                    </td>
                                    <td>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                           secondary-hidden-field="hidden"
                                                           :disabled="loading || !hasVisibleExpenseCategories"
                                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                           :show-selection-primary-text="true"
                                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(newRule.targetId, allCategories[CategoryType.Expense])"
                                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(newRule.targetId, allCategories[CategoryType.Expense])"
                                                           :items="allCategories[CategoryType.Expense]"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'expenseCategory'">
                                        </two-column-select>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                           secondary-hidden-field="hidden"
                                                           :disabled="loading || !hasVisibleIncomeCategories"
                                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                           :show-selection-primary-text="true"
                                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(newRule.targetId, allCategories[CategoryType.Income])"
                                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(newRule.targetId, allCategories[CategoryType.Income])"
                                                           :items="allCategories[CategoryType.Income]"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'incomeCategory'">
                                        </two-column-select>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                           secondary-hidden-field="hidden"
                                                           :disabled="loading || !hasVisibleTransferCategories"
                                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                           :show-selection-primary-text="true"
                                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(newRule.targetId, allCategories[CategoryType.Transfer])"
                                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(newRule.targetId, allCategories[CategoryType.Transfer])"
                                                           :items="allCategories[CategoryType.Transfer]"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'transferCategory'">
                                        </two-column-select>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="category"
                                                           primary-title-field="name" primary-footer-field="displayBalance"
                                                           primary-icon-field="icon" primary-icon-type="account"
                                                           primary-sub-items-field="accounts"
                                                           :primary-title-i18n="true"
                                                           secondary-key-field="id" secondary-value-field="id"
                                                           secondary-title-field="name" secondary-footer-field="displayBalance"
                                                           secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                           :disabled="loading || !allVisibleAccounts.length"
                                                           :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                           :custom-selection-primary-text="getAccountDisplayName(newRule.targetId)"
                                                           :items="allVisibleCategorizedAccounts"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'account'">
                                        </two-column-select>
                                        <v-autocomplete density="compact" variant="underlined"
                                                        item-title="name" item-value="id"
                                                        persistent-placeholder chips
                                                        :disabled="loading" :items="allTagsWithGroupHeader"
                                                        :no-data-text="tt('No available tag')"
                                                        v-model="newRule.targetId"
                                                        v-if="newRule.dataType == 'tag'">
                                            <template #chip="{ props, item }">
                                                <v-chip :prepend-icon="mdiPound" :text="item.title" v-bind="props" v-if="newRule.targetId"/>
                                            </template>

                                            <template #subheader="{ props }">
                                                <v-list-subheader>{{ props['title'] }}</v-list-subheader>
                                            </template>

                                            <template #item="{ props, item }">
                                                <v-list-item :value="item.value" v-bind="props" v-if="item.raw instanceof TransactionTag && !item.raw.hidden">
                                                    <template #title>
                                                        <v-list-item-title>
                                                            <div class="d-flex align-center">
                                                                <v-icon size="20" start :icon="mdiPound"/>
                                                                <span>{{ item.title }}</span>
                                                            </div>
                                                        </v-list-item-title>
                                                    </template>
                                                </v-list-item>
                                            </template>
                                        </v-autocomplete>
                                    </td>
                                    <td class="text-right">
                                        <v-btn density="comfortable" variant="tonal" color="primary"
                                               :disabled="loading || !newRule.dataType || !newRule.targetId"
                                               @click="addNewRule()">{{ tt('Add Rule') }}</v-btn>
                                    </td>
                                </tr>
                            </tfoot>
                        </v-table>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="loading" @click="confirm">{{ tt('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="loading" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <v-dialog width="400" v-model="saveDialogState">
        <v-card>
            <v-card-title>{{ tt('Save Rule Set') }}</v-card-title>
            <v-card-text>
                <v-text-field :label="tt('Rule Set Name')"
                              :placeholder="tt('Enter rule set name')"
                              v-model="saveDialogName"
                              :rules="[v => !!v || tt('Name is required')]"
                              autofocus
                              @keyup.enter="confirmSave()" />
            </v-card-text>
            <v-card-actions>
                <v-spacer/>
                <v-btn @click="saveDialogState = false">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" :disabled="!saveDialogName" @click="confirmSave()">{{ tt('Save') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useTransactionTagSelectionBase } from '@/components/base/TransactionTagSelectionBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useImportReplaceRuleStore } from '@/stores/importReplaceRule.ts';

import type { NameValue } from '@/core/base.ts';
import type { ImportReplaceRuleItem } from '@/models/import_replace_rule.ts';
import { CategoryType } from '@/core/category.ts';
import { ImportTransactionReplaceRule } from '@/core/import_transaction.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';

import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.ts';

import logger from '@/lib/logger.ts';

import {
    mdiRefresh,
    mdiPound,
    mdiDotsVertical,
    mdiContentSaveOutline,
    mdiFileDocumentOutline,
    mdiDeleteOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

interface BatchReplaceAllTypesDialogOpenOptions {
    fileType?: string;
    expenseCategoryNames: NameValue[];
    incomeCategoryNames: NameValue[];
    transferCategoryNames: NameValue[];
    accountNames: NameValue[];
    tagNames: NameValue[];
}

interface BatchReplaceAllTypesDialogResponse {
    rules: ImportTransactionReplaceRule[]
}

const { tt, getCategorizedAccountsWithDisplayBalance } = useI18n();

const { allTagsWithGroupHeader } = useTransactionTagSelectionBase({ modelValue: [] }, false);

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const importReplaceRuleStore = useImportReplaceRuleStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const rules = ref<ImportTransactionReplaceRule[]>([]);
const newRule = ref<ImportTransactionReplaceRule>(ImportTransactionReplaceRule.of('expenseCategory', '', ''));

const sourceExpenseCategoryNames = ref<NameValue[]>([]);
const sourceIncomeCategoryNames = ref<NameValue[]>([]);
const sourceTransferCategoryNames = ref<NameValue[]>([]);
const sourceAccountNames = ref<NameValue[]>([]);
const sourceTagNames = ref<NameValue[]>([]);

const savedRuleSets = ref<{ id: string; name: string; ruleCount: number }[]>([]);
const currentRuleSetId = ref<string | null>(null);
const currentFileType = ref<string>('');

const saveDialogState = ref<boolean>(false);
const saveDialogName = ref<string>('');

let resolveFunc: ((response: BatchReplaceAllTypesDialogResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const defaultRuleSetName = computed(() => getDefaultRuleSetName(currentFileType.value));

const allRuleSetsList = computed(() => {
    return savedRuleSets.value.map(rs => ({
        ...rs,
        isDefault: rs.name === defaultRuleSetName.value
    }));
});

const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);
const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value, customAccountCategoryOrder.value));
const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
const allTagsMap = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);

const hasVisibleExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleExpenseCategories);
const hasVisibleIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleIncomeCategories);
const hasVisibleTransferCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleTransferCategories);

const sourceItems = computed<NameValue[]>(() => {
    switch (newRule.value.dataType) {
        case 'expenseCategory':
            return sourceExpenseCategoryNames.value;
        case 'incomeCategory':
            return sourceIncomeCategoryNames.value;
        case 'transferCategory':
            return sourceTransferCategoryNames.value;
        case 'account':
            return sourceAccountNames.value;
        case 'tag':
            return sourceTagNames.value;
        default:
            return [];
    }
});

const noSourceItemText = computed<string>(() => {
    switch (newRule.value.dataType) {
        case 'expenseCategory':
            return tt('No available category');
        case 'incomeCategory':
            return tt('No available category');
        case 'transferCategory':
            return tt('No available category');
        case 'account':
            return tt('No available account');
        case 'tag':
            return tt('No available tag');
        default:
            return '';
    }
});

function getRuleTypeDisplayName(rule: ImportTransactionReplaceRule): string {
    switch (rule.dataType) {
        case 'expenseCategory':
            return tt('Expense Category');
        case 'incomeCategory':
            return tt('Income Category');
        case 'transferCategory':
            return tt('Transfer Category');
        case 'account':
            return tt('Account');
        case 'tag':
            return tt('Transaction Tag');
        default:
            return '';
    }
}

function getRuleTargetValueDisplayName(rule: ImportTransactionReplaceRule): string {
    switch (rule.dataType) {
        case 'expenseCategory':
            return getTransactionSecondaryCategoryName(rule.targetId, allCategories.value[CategoryType.Expense]) || '';
        case 'incomeCategory':
            return getTransactionSecondaryCategoryName(rule.targetId, allCategories.value[CategoryType.Income]) || '';
        case 'transferCategory':
            return getTransactionSecondaryCategoryName(rule.targetId, allCategories.value[CategoryType.Transfer]) || '';
        case 'account':
            return getAccountDisplayName(rule.targetId);
        case 'tag':
            return allTagsMap.value[rule.targetId]?.name ?? '';
        default:
            return '';
    }
}

function getAccountDisplayName(accountId?: string): string {
    if (accountId) {
        return Account.findAccountNameById(allAccounts.value, accountId) || '';
    } else {
        return tt('None');
    }
}

function getDefaultRuleSetName(fileType: string): string {
    const defaultNames: Record<string, string> = {
        'ezbookkeeping': 'ezbookkeeping数据默认规则',
        'dsv': '分隔符分隔值(DSV)文件默认规则',
        'dsv_data': '分隔符分隔值(DSV)数据默认规则',
        'excel': 'Excel工作簿文件默认规则',
        'ofx': '开放式金融交换(OFX)文件默认规则',
        'qfx': 'Quicken Financial Exchange (QFX)文件默认规则',
        'qif': 'Quicken Interchange Format (QIF)文件默认规则',
        'iif': 'Intuit Interchange Format (IIF)文件默认规则',
        'camt052': 'Camt.052银行对账单文件默认规则',
        'camt053': 'Camt.053银行对账单文件默认规则',
        'mt940': 'MT940客户对账消息文件默认规则',
        'alipay_app_csv': '支付宝(App)交易流水文件默认规则',
        'alipay_web_csv': '支付宝(网页版)交易流水文件默认规则',
        'wechat_pay_app': '微信支付账单文件默认规则',
        'jdcom_finance_app_csv': '京东金融账单文件默认规则',
        'gnucash': 'GnuCash XML 数据库文件默认规则',
        'firefly_iii_csv': 'Firefly III 数据默认规则',
        'beancount': 'Beancount数据文件默认规则',
        'feidee_mymoney_csv': '随手记(App)数据默认规则',
        'feidee_mymoney_xls': '随手记(Web版)数据默认规则',
        'feidee_mymoney_elecloud_xlsx': '随手记(神象云账本)数据默认规则'
    };
    return defaultNames[fileType] || '默认规则';
}

function open(options: BatchReplaceAllTypesDialogOpenOptions): Promise<BatchReplaceAllTypesDialogResponse> {
    rules.value = [];
    newRule.value = ImportTransactionReplaceRule.of('expenseCategory', '', '');
    sourceExpenseCategoryNames.value = options.expenseCategoryNames;
    sourceIncomeCategoryNames.value = options.incomeCategoryNames;
    sourceTransferCategoryNames.value = options.transferCategoryNames;
    sourceAccountNames.value = options.accountNames;
    sourceTagNames.value = options.tagNames;
    currentFileType.value = options.fileType || '';
    currentRuleSetId.value = null;
    savedRuleSets.value = [];
    showState.value = true;

    // Load rule sets and auto-load matching one
    loadRuleSetsFromDatabase();

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
    loading.value = true;

    Promise.allSettled([
        accountsStore.loadAllAccounts({ force: true }),
        transactionCategoriesStore.loadAllCategories({ force: true }),
        transactionTagsStore.loadAllTags({ force: true })
    ]).then(results => {
        loading.value = false;

        const isAllUpToDate = results.length === 3
            && results[0].status === 'rejected' && results[0].reason?.isUpToDate
            && results[1].status === 'rejected' && results[1].reason?.isUpToDate
            && results[2].status === 'rejected' && results[2].reason?.isUpToDate;

        // show info if all up to date
        if (isAllUpToDate) {
            snackbar.value?.showMessage('Data is up to date');
            return;
        }

        // show error if any
        for (const result of results) {
            if (result.status === 'rejected' && !result.reason?.isUpToDate) {
                snackbar.value?.showError(result.reason);
                return;
            }
        }

        // show info if one of them updated
        for (const result of results) {
            if (result.status === 'fulfilled') {
                snackbar.value?.showMessage('Data has been updated');
                return;
            }
        }
    });
}

function loadRuleSetsFromDatabase(): void {
    loading.value = true;
    importReplaceRuleStore.getAllRuleSets().then(result => {
        savedRuleSets.value = result.map(rs => ({ id: rs.id, name: rs.name, ruleCount: rs.ruleCount }));

        // Try to get default rule set for this file type from database
        importReplaceRuleStore.getDefaultByFileType(currentFileType.value).then(defaultResult => {
            if (defaultResult && defaultResult.id) {
                currentRuleSetId.value = defaultResult.id;
                const parsedRules: ImportTransactionReplaceRule[] = [];
                for (const rule of defaultResult.rules) {
                    parsedRules.push(ImportTransactionReplaceRule.of(rule.type as ImportTransactionReplaceRule['dataType'], rule.sourceValue, rule.targetId));
                }
                rules.value = parsedRules;
            } else {
                // No default for this file type - use default name matching
                const defaultName = getDefaultRuleSetName(currentFileType.value);
                const matched = savedRuleSets.value.find(rs => rs.name === defaultName);
                if (matched) {
                    currentRuleSetId.value = matched.id;
                    loadRuleSetDetail(matched.id);
                }
            }
        }).catch(() => {
            // Fallback to default name matching
            const defaultName = getDefaultRuleSetName(currentFileType.value);
            const matched = savedRuleSets.value.find(rs => rs.name === defaultName);
            if (matched) {
                currentRuleSetId.value = matched.id;
                loadRuleSetDetail(matched.id);
            }
        });
    }).catch(err => {
        logger.error('Failed to load rule sets', err);
    }).finally(() => {
        loading.value = false;
    });
}

function loadRuleSetDetail(id: string): void {
    loading.value = true;
    importReplaceRuleStore.getRuleSet(id).then(result => {
        const parsedRules: ImportTransactionReplaceRule[] = [];
        for (const rule of result.rules) {
            parsedRules.push(ImportTransactionReplaceRule.of(rule.type as ImportTransactionReplaceRule['dataType'], rule.sourceValue, rule.targetId));
        }
        rules.value = parsedRules;
        currentRuleSetId.value = id;
    }).catch(err => {
        logger.error('Failed to load rule set detail', err);
        snackbar.value?.showError(tt('Failed to load rule set'));
    }).finally(() => {
        loading.value = false;
    });
}

function loadOtherRuleSet(id: string): void {
    importReplaceRuleStore.setDefaultRuleSet(id, currentFileType.value).catch(() => { /* ignore errors */ });
    loadRuleSetDetail(id);
}

function deleteRuleSet(id: string): void {
    loading.value = true;
    importReplaceRuleStore.deleteRuleSet({ id }).then(() => {
        snackbar.value?.showMessage(tt('Rule set deleted successfully'));
        // Clear current so loadRuleSetsFromDatabase will find the default
        currentRuleSetId.value = null;
        rules.value = [];
        // loadRuleSetsFromDatabase manages its own loading state
        loadRuleSetsFromDatabase();
    }).catch(err => {
        logger.error('Failed to delete rule set', err);
        snackbar.value?.showError(tt('Failed to delete rule set'));
        loading.value = false;
    });
}

function showSaveDialog(): void {
    if (!rules.value.length) {
        snackbar.value?.showError(tt('No rules to save'));
        return;
    }
    // If a rule set was loaded, use its name; otherwise use default name
    if (currentRuleSetId.value) {
        const matched = savedRuleSets.value.find(rs => rs.id === currentRuleSetId.value);
        saveDialogName.value = matched ? matched.name : getDefaultRuleSetName(currentFileType.value);
    } else {
        saveDialogName.value = getDefaultRuleSetName(currentFileType.value);
    }
    saveDialogState.value = true;
}

function confirmSave(): void {
    if (!saveDialogName.value) {
        return;
    }

    const ruleItems: ImportReplaceRuleItem[] = rules.value.map(r => ({
        type: r.dataType as ImportReplaceRuleItem['type'],
        sourceValue: r.sourceValue,
        targetId: r.targetId
    }));

    loading.value = true;
    saveDialogState.value = false;

    // Check if name matches existing rule set -> update; otherwise create new
    const matchedByName = savedRuleSets.value.find(rs => rs.name === saveDialogName.value);

    if (matchedByName) {
        importReplaceRuleStore.modifyRuleSet({
            id: matchedByName.id,
            name: saveDialogName.value,
            fileType: currentFileType.value,
            rules: ruleItems
        }).then(() => {
            currentRuleSetId.value = matchedByName.id;
            snackbar.value?.showMessage(tt('Rule set saved successfully'));
            loadRuleSetsFromDatabase();
        }).catch(err => {
            logger.error('Failed to modify rule set', err);
            snackbar.value?.showError(tt('Failed to save rule set'));
        }).finally(() => {
            loading.value = false;
        });
    } else {
        importReplaceRuleStore.addRuleSet({
            name: saveDialogName.value,
            fileType: currentFileType.value,
            rules: ruleItems
        }).then(result => {
            currentRuleSetId.value = result.id;
            snackbar.value?.showMessage(tt('Rule set created successfully'));
            loadRuleSetsFromDatabase();
        }).catch(err => {
            logger.error('Failed to create rule set', err);
            snackbar.value?.showError(tt('Failed to save rule set'));
        }).finally(() => {
            loading.value = false;
        });
    }
}

function removeRule(index: number): void {
    rules.value.splice(index, 1);
}

function addNewRule(): void {
    if (!newRule.value.dataType || !newRule.value.targetId) {
        return;
    }

    rules.value.push(newRule.value);
    newRule.value = ImportTransactionReplaceRule.of('expenseCategory', '', '');
}

function confirm(): void {
    resolveFunc?.({
        rules: rules.value
    });
    showState.value = false;
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
