"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
var vue_1 = require("vue");
var BackupApiClient_1 = require("../lib/BackupApiClient");
var ActionButtonComponent_vue_1 = require("./ActionButtonComponent.vue");
var isBackupLoading = (0, vue_1.ref)(false);
var isCleanupLoading = (0, vue_1.ref)(false);
var backupResult = (0, vue_1.ref)(null);
var cleanupResult = (0, vue_1.ref)(null);
function triggerBackupAll() {
    return __awaiter(this, void 0, void 0, function () {
        var result, err_1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    isBackupLoading.value = true;
                    backupResult.value = null;
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 3, 4, 5]);
                    return [4 /*yield*/, BackupApiClient_1.backupApiClient.triggerBackupAll()];
                case 2:
                    result = _a.sent();
                    backupResult.value = result;
                    setTimeout(function () {
                        backupResult.value = null;
                    }, 10000);
                    return [3 /*break*/, 5];
                case 3:
                    err_1 = _a.sent();
                    backupResult.value = {
                        success: false,
                        message: err_1 instanceof Error ? err_1.message : "Backup failed",
                    };
                    return [3 /*break*/, 5];
                case 4:
                    isBackupLoading.value = false;
                    return [7 /*endfinally*/];
                case 5: return [2 /*return*/];
            }
        });
    });
}
function triggerCleanupAll() {
    return __awaiter(this, void 0, void 0, function () {
        var result, err_2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    isCleanupLoading.value = true;
                    cleanupResult.value = null;
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 3, 4, 5]);
                    return [4 /*yield*/, BackupApiClient_1.backupApiClient.triggerCleanupAll()];
                case 2:
                    result = _a.sent();
                    cleanupResult.value = result;
                    setTimeout(function () {
                        cleanupResult.value = null;
                    }, 10000);
                    return [3 /*break*/, 5];
                case 3:
                    err_2 = _a.sent();
                    cleanupResult.value = {
                        success: false,
                        message: err_2 instanceof Error ? err_2.message : "Cleanup failed",
                    };
                    return [3 /*break*/, 5];
                case 4:
                    isCleanupLoading.value = false;
                    return [7 /*endfinally*/];
                case 5: return [2 /*return*/];
            }
        });
    });
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
var __VLS_ctx = {};
var __VLS_elements;
var __VLS_components;
var __VLS_directives;
/** @type {__VLS_StyleScopedClasses['action-panel']} */ ;
/** @type {__VLS_StyleScopedClasses['action-result']} */ ;
/** @type {__VLS_StyleScopedClasses['action-result']} */ ;
/** @type {__VLS_StyleScopedClasses['action-info']} */ ;
/** @type {__VLS_StyleScopedClasses['actions']} */ ;
// CSS variable injection 
// CSS variable injection end 
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "action-panel" }));
__VLS_asFunctionalElement(__VLS_elements.h2, __VLS_elements.h2)({});
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "actions" }));
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "action-group" }));
/** @type {[typeof ActionButtonComponent, ]} */ ;
// @ts-ignore
var __VLS_0 = __VLS_asFunctionalComponent(ActionButtonComponent_vue_1.default, new ActionButtonComponent_vue_1.default(__assign({ 'onClick': {} }, { label: "Backup All Targets", disabled: (__VLS_ctx.isBackupLoading) })));
var __VLS_1 = __VLS_0.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { label: "Backup All Targets", disabled: (__VLS_ctx.isBackupLoading) })], __VLS_functionalComponentArgsRest(__VLS_0), false));
var __VLS_3;
var __VLS_4;
var __VLS_5 = ({ click: {} },
    { onClick: (__VLS_ctx.triggerBackupAll) });
// @ts-ignore
[isBackupLoading, triggerBackupAll,];
var __VLS_2;
if (__VLS_ctx.backupResult) {
    // @ts-ignore
    [backupResult,];
    __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: ([
            'action-result',
            __VLS_ctx.backupResult.success ? 'success' : 'error'
        ]) }));
    // @ts-ignore
    [backupResult,];
    (__VLS_ctx.backupResult.message);
    // @ts-ignore
    [backupResult,];
}
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "action-group" }));
/** @type {[typeof ActionButtonComponent, ]} */ ;
// @ts-ignore
var __VLS_7 = __VLS_asFunctionalComponent(ActionButtonComponent_vue_1.default, new ActionButtonComponent_vue_1.default(__assign({ 'onClick': {} }, { label: "Cleanup All Targets", variant: "secondary", disabled: (__VLS_ctx.isCleanupLoading) })));
var __VLS_8 = __VLS_7.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { label: "Cleanup All Targets", variant: "secondary", disabled: (__VLS_ctx.isCleanupLoading) })], __VLS_functionalComponentArgsRest(__VLS_7), false));
var __VLS_10;
var __VLS_11;
var __VLS_12 = ({ click: {} },
    { onClick: (__VLS_ctx.triggerCleanupAll) });
// @ts-ignore
[isCleanupLoading, triggerCleanupAll,];
var __VLS_9;
if (__VLS_ctx.cleanupResult) {
    // @ts-ignore
    [cleanupResult,];
    __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: ([
            'action-result',
            __VLS_ctx.cleanupResult.success ? 'success' : 'error'
        ]) }));
    // @ts-ignore
    [cleanupResult,];
    (__VLS_ctx.cleanupResult.message);
    // @ts-ignore
    [cleanupResult,];
}
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "action-info" }));
__VLS_asFunctionalElement(__VLS_elements.p, __VLS_elements.p)({});
__VLS_asFunctionalElement(__VLS_elements.strong, __VLS_elements.strong)({});
__VLS_asFunctionalElement(__VLS_elements.p, __VLS_elements.p)({});
__VLS_asFunctionalElement(__VLS_elements.strong, __VLS_elements.strong)({});
/** @type {__VLS_StyleScopedClasses['action-panel']} */ ;
/** @type {__VLS_StyleScopedClasses['actions']} */ ;
/** @type {__VLS_StyleScopedClasses['action-group']} */ ;
/** @type {__VLS_StyleScopedClasses['action-result']} */ ;
/** @type {__VLS_StyleScopedClasses['action-group']} */ ;
/** @type {__VLS_StyleScopedClasses['action-result']} */ ;
/** @type {__VLS_StyleScopedClasses['action-info']} */ ;
var __VLS_dollars;
var __VLS_self = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
        return {
            ActionButtonComponent: ActionButtonComponent_vue_1.default,
            isBackupLoading: isBackupLoading,
            isCleanupLoading: isCleanupLoading,
            backupResult: backupResult,
            cleanupResult: cleanupResult,
            triggerBackupAll: triggerBackupAll,
            triggerCleanupAll: triggerCleanupAll,
        };
    },
});
exports.default = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
    },
});
; /* PartiallyEnd: #4569/main.vue */
