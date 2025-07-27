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
Object.defineProperty(exports, "__esModule", { value: true });
var ComponentTypes_1 = require("../lib/ComponentTypes");
var props = defineProps({
    label: ComponentTypes_1.StringProp,
    disabled: __assign(__assign({}, ComponentTypes_1.BooleanProp), { required: false, default: false }),
    variant: {
        type: String,
        required: false,
        default: "primary",
        validator: function (value) { return ["primary", "secondary", "danger"].includes(value); },
    },
});
var emit = defineEmits();
function handleClick() {
    if (!props.disabled) {
        emit("click");
    }
}
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
var __VLS_ctx = {};
var __VLS_elements;
var __VLS_components;
var __VLS_directives;
/** @type {__VLS_StyleScopedClasses['btn']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-primary']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-secondary']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-danger']} */ ;
// CSS variable injection 
// CSS variable injection end 
__VLS_asFunctionalElement(__VLS_elements.button, __VLS_elements.button)(__assign(__assign({ onClick: (__VLS_ctx.handleClick) }, { class: ([
        'btn',
        "btn-".concat(props.variant),
        { 'btn-disabled': props.disabled }
    ]) }), { disabled: (props.disabled) }));
// @ts-ignore
[handleClick,];
(props.label);
/** @type {__VLS_StyleScopedClasses['btn']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-disabled']} */ ;
var __VLS_dollars;
var __VLS_self = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
        return {
            handleClick: handleClick,
        };
    },
    __typeEmits: {},
    props: {
        label: ComponentTypes_1.StringProp,
        disabled: __assign(__assign({}, ComponentTypes_1.BooleanProp), { required: false, default: false }),
        variant: {
            type: String,
            required: false,
            default: "primary",
            validator: function (value) { return ["primary", "secondary", "danger"].includes(value); },
        },
    },
});
exports.default = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
    },
    __typeEmits: {},
    props: {
        label: ComponentTypes_1.StringProp,
        disabled: __assign(__assign({}, ComponentTypes_1.BooleanProp), { required: false, default: false }),
        variant: {
            type: String,
            required: false,
            default: "primary",
            validator: function (value) { return ["primary", "secondary", "danger"].includes(value); },
        },
    },
});
; /* PartiallyEnd: #4569/main.vue */
