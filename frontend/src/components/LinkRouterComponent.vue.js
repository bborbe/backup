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
var props = defineProps({
    to: {
        type: Object,
        required: true,
    },
    title: {
        type: String,
        required: true,
    },
});
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
var __VLS_ctx = {};
var __VLS_elements;
var __VLS_components;
var __VLS_directives;
/** @type {__VLS_StyleScopedClasses['link']} */ ;
/** @type {__VLS_StyleScopedClasses['link']} */ ;
// CSS variable injection 
// CSS variable injection end 
var __VLS_0 = {}.RouterLink;
/** @type {[typeof __VLS_components.RouterLink, typeof __VLS_components.routerLink, typeof __VLS_components.RouterLink, typeof __VLS_components.routerLink, ]} */ ;
// @ts-ignore
RouterLink;
// @ts-ignore
var __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0(__assign({ class: "link" }, { to: (props.to), title: (props.title) })));
var __VLS_2 = __VLS_1.apply(void 0, __spreadArray([__assign({ class: "link" }, { to: (props.to), title: (props.title) })], __VLS_functionalComponentArgsRest(__VLS_1), false));
var __VLS_4 = {};
var __VLS_5 = __VLS_3.slots.default;
var __VLS_6 = {};
var __VLS_3;
/** @type {__VLS_StyleScopedClasses['link']} */ ;
// @ts-ignore
var __VLS_7 = __VLS_6;
var __VLS_dollars;
var __VLS_self = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
        return {};
    },
    props: {
        to: {
            type: Object,
            required: true,
        },
        title: {
            type: String,
            required: true,
        },
    },
});
var __VLS_component = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
    },
    props: {
        to: {
            type: Object,
            required: true,
        },
        title: {
            type: String,
            required: true,
        },
    },
});
exports.default = {};
; /* PartiallyEnd: #4569/main.vue */
