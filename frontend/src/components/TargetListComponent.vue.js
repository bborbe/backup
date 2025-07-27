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
var _a, _b, _c, _d, _e, _f, _g, _h, _j;
Object.defineProperty(exports, "__esModule", { value: true });
var vue_1 = require("vue");
var BackupApiClient_1 = require("../lib/BackupApiClient");
var ActionButtonComponent_vue_1 = require("./ActionButtonComponent.vue");
var targets = (0, vue_1.ref)([]);
var loadingState = (0, vue_1.ref)({
    isLoading: false,
    error: null,
});
var actionStates = (0, vue_1.ref)({});
(0, vue_1.onMounted)(function () { return __awaiter(void 0, void 0, void 0, function () {
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [4 /*yield*/, loadTargets()];
            case 1:
                _a.sent();
                return [2 /*return*/];
        }
    });
}); });
function loadTargets() {
    return __awaiter(this, void 0, void 0, function () {
        var _a, err_1;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    loadingState.value.isLoading = true;
                    loadingState.value.error = null;
                    _b.label = 1;
                case 1:
                    _b.trys.push([1, 3, 4, 5]);
                    _a = targets;
                    return [4 /*yield*/, BackupApiClient_1.backupApiClient.listTargets()];
                case 2:
                    _a.value = _b.sent();
                    return [3 /*break*/, 5];
                case 3:
                    err_1 = _b.sent();
                    loadingState.value.error = err_1 instanceof Error ? err_1.message : "Failed to load targets";
                    return [3 /*break*/, 5];
                case 4:
                    loadingState.value.isLoading = false;
                    return [7 /*endfinally*/];
                case 5: return [2 /*return*/];
            }
        });
    });
}
function triggerBackup(host) {
    return __awaiter(this, void 0, void 0, function () {
        var result, err_2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (!actionStates.value[host]) {
                        actionStates.value[host] = { isLoading: false, result: null };
                    }
                    actionStates.value[host].isLoading = true;
                    actionStates.value[host].result = null;
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 3, 4, 5]);
                    return [4 /*yield*/, BackupApiClient_1.backupApiClient.triggerBackup(host)];
                case 2:
                    result = _a.sent();
                    actionStates.value[host].result = result;
                    setTimeout(function () {
                        if (actionStates.value[host]) {
                            actionStates.value[host].result = null;
                        }
                    }, 5000);
                    return [3 /*break*/, 5];
                case 3:
                    err_2 = _a.sent();
                    actionStates.value[host].result = {
                        success: false,
                        message: err_2 instanceof Error ? err_2.message : "Backup failed",
                    };
                    return [3 /*break*/, 5];
                case 4:
                    actionStates.value[host].isLoading = false;
                    return [7 /*endfinally*/];
                case 5: return [2 /*return*/];
            }
        });
    });
}
function triggerCleanup(host) {
    return __awaiter(this, void 0, void 0, function () {
        var result, err_3;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (!actionStates.value[host]) {
                        actionStates.value[host] = { isLoading: false, result: null };
                    }
                    actionStates.value[host].isLoading = true;
                    actionStates.value[host].result = null;
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 3, 4, 5]);
                    return [4 /*yield*/, BackupApiClient_1.backupApiClient.triggerCleanup(host)];
                case 2:
                    result = _a.sent();
                    actionStates.value[host].result = result;
                    setTimeout(function () {
                        if (actionStates.value[host]) {
                            actionStates.value[host].result = null;
                        }
                    }, 5000);
                    return [3 /*break*/, 5];
                case 3:
                    err_3 = _a.sent();
                    actionStates.value[host].result = {
                        success: false,
                        message: err_3 instanceof Error ? err_3.message : "Cleanup failed",
                    };
                    return [3 /*break*/, 5];
                case 4:
                    actionStates.value[host].isLoading = false;
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
/** @type {__VLS_StyleScopedClasses['target-list']} */ ;
/** @type {__VLS_StyleScopedClasses['error']} */ ;
/** @type {__VLS_StyleScopedClasses['target-dirs']} */ ;
/** @type {__VLS_StyleScopedClasses['target-excludes']} */ ;
/** @type {__VLS_StyleScopedClasses['target-dirs']} */ ;
/** @type {__VLS_StyleScopedClasses['target-excludes']} */ ;
/** @type {__VLS_StyleScopedClasses['action-result']} */ ;
/** @type {__VLS_StyleScopedClasses['action-result']} */ ;
/** @type {__VLS_StyleScopedClasses['error']} */ ;
// CSS variable injection 
// CSS variable injection end 
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "target-list" }));
__VLS_asFunctionalElement(__VLS_elements.h2, __VLS_elements.h2)({});
if (__VLS_ctx.loadingState.isLoading) {
    // @ts-ignore
    [loadingState,];
    __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "loading" }));
}
else if (__VLS_ctx.loadingState.error) {
    // @ts-ignore
    [loadingState,];
    __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "error" }));
    (__VLS_ctx.loadingState.error);
    // @ts-ignore
    [loadingState,];
    __VLS_asFunctionalElement(__VLS_elements.button, __VLS_elements.button)(__assign({ onClick: (__VLS_ctx.loadTargets) }, { class: "retry-btn" }));
    // @ts-ignore
    [loadTargets,];
}
else if (__VLS_ctx.targets.length === 0) {
    // @ts-ignore
    [targets,];
    __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "empty" }));
}
else {
    __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "targets" }));
    var _loop_1 = function (target) {
        // @ts-ignore
        [targets,];
        __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ key: (target.host) }, { class: "target-card" }));
        __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "target-header" }));
        __VLS_asFunctionalElement(__VLS_elements.h3, __VLS_elements.h3)({});
        (target.host);
        __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "target-connection" }));
        (target.user);
        (target.host);
        (target.port);
        __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "target-dirs" }));
        __VLS_asFunctionalElement(__VLS_elements.h4, __VLS_elements.h4)({});
        __VLS_asFunctionalElement(__VLS_elements.ul, __VLS_elements.ul)({});
        for (var _l = 0, _m = __VLS_getVForSourceType((target.dirs)); _l < _m.length; _l++) {
            var dir = _m[_l][0];
            __VLS_asFunctionalElement(__VLS_elements.li, __VLS_elements.li)({
                key: (dir),
            });
            (dir);
        }
        if (target.excludes.length > 0) {
            __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "target-excludes" }));
            __VLS_asFunctionalElement(__VLS_elements.h4, __VLS_elements.h4)({});
            __VLS_asFunctionalElement(__VLS_elements.ul, __VLS_elements.ul)({});
            for (var _o = 0, _p = __VLS_getVForSourceType((target.excludes)); _o < _p.length; _o++) {
                var exclude = _p[_o][0];
                __VLS_asFunctionalElement(__VLS_elements.li, __VLS_elements.li)({
                    key: (exclude),
                });
                (exclude);
            }
        }
        __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "target-actions" }));
        /** @type {[typeof ActionButtonComponent, ]} */ ;
        // @ts-ignore
        var __VLS_0 = __VLS_asFunctionalComponent(ActionButtonComponent_vue_1.default, new ActionButtonComponent_vue_1.default(__assign({ 'onClick': {} }, { label: "Backup", disabled: (((_a = __VLS_ctx.actionStates[target.host]) === null || _a === void 0 ? void 0 : _a.isLoading) || false) })));
        var __VLS_1 = __VLS_0.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { label: "Backup", disabled: (((_b = __VLS_ctx.actionStates[target.host]) === null || _b === void 0 ? void 0 : _b.isLoading) || false) })], __VLS_functionalComponentArgsRest(__VLS_0), false));
        var __VLS_3 = void 0;
        var __VLS_4 = void 0;
        var __VLS_5 = ({ click: {} },
            { onClick: function () {
                    var _a = [];
                    for (var _i = 0; _i < arguments.length; _i++) {
                        _a[_i] = arguments[_i];
                    }
                    var $event = _a[0];
                    if (!!(__VLS_ctx.loadingState.isLoading))
                        return;
                    if (!!(__VLS_ctx.loadingState.error))
                        return;
                    if (!!(__VLS_ctx.targets.length === 0))
                        return;
                    __VLS_ctx.triggerBackup(target.host);
                    // @ts-ignore
                    [actionStates, triggerBackup,];
                } });
        /** @type {[typeof ActionButtonComponent, ]} */ ;
        // @ts-ignore
        var __VLS_7 = __VLS_asFunctionalComponent(ActionButtonComponent_vue_1.default, new ActionButtonComponent_vue_1.default(__assign({ 'onClick': {} }, { label: "Cleanup", variant: "secondary", disabled: (((_c = __VLS_ctx.actionStates[target.host]) === null || _c === void 0 ? void 0 : _c.isLoading) || false) })));
        var __VLS_8 = __VLS_7.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { label: "Cleanup", variant: "secondary", disabled: (((_d = __VLS_ctx.actionStates[target.host]) === null || _d === void 0 ? void 0 : _d.isLoading) || false) })], __VLS_functionalComponentArgsRest(__VLS_7), false));
        var __VLS_10 = void 0;
        var __VLS_11 = void 0;
        var __VLS_12 = ({ click: {} },
            { onClick: function () {
                    var _a = [];
                    for (var _i = 0; _i < arguments.length; _i++) {
                        _a[_i] = arguments[_i];
                    }
                    var $event = _a[0];
                    if (!!(__VLS_ctx.loadingState.isLoading))
                        return;
                    if (!!(__VLS_ctx.loadingState.error))
                        return;
                    if (!!(__VLS_ctx.targets.length === 0))
                        return;
                    __VLS_ctx.triggerCleanup(target.host);
                    // @ts-ignore
                    [actionStates, triggerCleanup,];
                } });
        if ((_e = __VLS_ctx.actionStates[target.host]) === null || _e === void 0 ? void 0 : _e.result) {
            // @ts-ignore
            [actionStates,];
            __VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: ([
                    'action-result',
                    ((_g = (_f = __VLS_ctx.actionStates[target.host]) === null || _f === void 0 ? void 0 : _f.result) === null || _g === void 0 ? void 0 : _g.success) ? 'success' : 'error'
                ]) }));
            // @ts-ignore
            [actionStates,];
            ((_j = (_h = __VLS_ctx.actionStates[target.host]) === null || _h === void 0 ? void 0 : _h.result) === null || _j === void 0 ? void 0 : _j.message);
            // @ts-ignore
            [actionStates,];
        }
    };
    var __VLS_2, __VLS_9;
    for (var _i = 0, _k = __VLS_getVForSourceType((__VLS_ctx.targets)); _i < _k.length; _i++) {
        var target = _k[_i][0];
        _loop_1(target);
    }
}
/** @type {__VLS_StyleScopedClasses['target-list']} */ ;
/** @type {__VLS_StyleScopedClasses['loading']} */ ;
/** @type {__VLS_StyleScopedClasses['error']} */ ;
/** @type {__VLS_StyleScopedClasses['retry-btn']} */ ;
/** @type {__VLS_StyleScopedClasses['empty']} */ ;
/** @type {__VLS_StyleScopedClasses['targets']} */ ;
/** @type {__VLS_StyleScopedClasses['target-card']} */ ;
/** @type {__VLS_StyleScopedClasses['target-header']} */ ;
/** @type {__VLS_StyleScopedClasses['target-connection']} */ ;
/** @type {__VLS_StyleScopedClasses['target-dirs']} */ ;
/** @type {__VLS_StyleScopedClasses['target-excludes']} */ ;
/** @type {__VLS_StyleScopedClasses['target-actions']} */ ;
/** @type {__VLS_StyleScopedClasses['action-result']} */ ;
var __VLS_dollars;
var __VLS_self = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
        return {
            ActionButtonComponent: ActionButtonComponent_vue_1.default,
            targets: targets,
            loadingState: loadingState,
            actionStates: actionStates,
            loadTargets: loadTargets,
            triggerBackup: triggerBackup,
            triggerCleanup: triggerCleanup,
        };
    },
});
exports.default = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
    },
});
; /* PartiallyEnd: #4569/main.vue */
