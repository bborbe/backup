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
var BackupStatusOverviewComponent_vue_1 = require("../components/BackupStatusOverviewComponent.vue");
var ActionPanelComponent_vue_1 = require("../components/ActionPanelComponent.vue");
var TargetListComponent_vue_1 = require("../components/TargetListComponent.vue");
var currentTime = (0, vue_1.ref)("");
(0, vue_1.onMounted)(function () {
    updateTime();
    setInterval(updateTime, 1000);
});
function updateTime() {
    currentTime.value = new Date().toISOString().split(".")[0].split("T")[1];
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
var __VLS_ctx = {};
var __VLS_elements;
var __VLS_components;
var __VLS_directives;
/** @type {__VLS_StyleScopedClasses['dashboard-header']} */ ;
/** @type {__VLS_StyleScopedClasses['dashboard-header']} */ ;
/** @type {__VLS_StyleScopedClasses['dashboard-content']} */ ;
// CSS variable injection 
// CSS variable injection end 
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "dashboard" }));
__VLS_asFunctionalElement(__VLS_elements.header, __VLS_elements.header)(__assign({ class: "dashboard-header" }));
__VLS_asFunctionalElement(__VLS_elements.h1, __VLS_elements.h1)({});
__VLS_asFunctionalElement(__VLS_elements.div, __VLS_elements.div)(__assign({ class: "time" }));
(__VLS_ctx.currentTime);
// @ts-ignore
[currentTime,];
__VLS_asFunctionalElement(__VLS_elements.main, __VLS_elements.main)(__assign({ class: "dashboard-content" }));
/** @type {[typeof ActionPanelComponent, ]} */ ;
// @ts-ignore
var __VLS_0 = __VLS_asFunctionalComponent(ActionPanelComponent_vue_1.default, new ActionPanelComponent_vue_1.default({}));
var __VLS_1 = __VLS_0.apply(void 0, __spreadArray([{}], __VLS_functionalComponentArgsRest(__VLS_0), false));
/** @type {[typeof BackupStatusOverviewComponent, ]} */ ;
// @ts-ignore
var __VLS_4 = __VLS_asFunctionalComponent(BackupStatusOverviewComponent_vue_1.default, new BackupStatusOverviewComponent_vue_1.default({}));
var __VLS_5 = __VLS_4.apply(void 0, __spreadArray([{}], __VLS_functionalComponentArgsRest(__VLS_4), false));
/** @type {[typeof TargetListComponent, ]} */ ;
// @ts-ignore
var __VLS_8 = __VLS_asFunctionalComponent(TargetListComponent_vue_1.default, new TargetListComponent_vue_1.default({}));
var __VLS_9 = __VLS_8.apply(void 0, __spreadArray([{}], __VLS_functionalComponentArgsRest(__VLS_8), false));
/** @type {__VLS_StyleScopedClasses['dashboard']} */ ;
/** @type {__VLS_StyleScopedClasses['dashboard-header']} */ ;
/** @type {__VLS_StyleScopedClasses['time']} */ ;
/** @type {__VLS_StyleScopedClasses['dashboard-content']} */ ;
var __VLS_dollars;
var __VLS_self = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
        return {
            BackupStatusOverviewComponent: BackupStatusOverviewComponent_vue_1.default,
            ActionPanelComponent: ActionPanelComponent_vue_1.default,
            TargetListComponent: TargetListComponent_vue_1.default,
            currentTime: currentTime,
        };
    },
});
exports.default = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
    },
});
; /* PartiallyEnd: #4569/main.vue */
