import { PropType } from "vue";

export const StringProp = {
  type: String,
  required: true,
} as const;

export const OptionalStringProp = {
  type: String,
  required: false,
  default: "",
} as const;

export const BooleanProp = {
  type: Boolean,
  required: true,
} as const;

export const OptionalBooleanProp = {
  type: Boolean,
  required: false,
  default: false,
} as const;

export function TypedObjectProp<T>() {
  return {
    type: Object as PropType<T>,
    required: true,
  } as const;
}

export function OptionalTypedObjectProp<T>() {
  return {
    type: Object as PropType<T>,
    required: false,
  } as const;
}

export function TypedArrayProp<T>() {
  return {
    type: Array as PropType<T[]>,
    required: true,
  } as const;
}

export function OptionalTypedArrayProp<T>() {
  return {
    type: Array as PropType<T[]>,
    required: false,
    default: () => [],
  } as const;
}