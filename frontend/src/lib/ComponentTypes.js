export const StringProp = {
  type: String,
  required: true,
};

export const OptionalStringProp = {
  type: String,
  required: false,
  default: "",
};

export const BooleanProp = {
  type: Boolean,
  required: true,
};

export const OptionalBooleanProp = {
  type: Boolean,
  required: false,
  default: false,
};

export function TypedObjectProp() {
  return {
    type: Object,
    required: true,
  };
}

export function OptionalTypedObjectProp() {
  return {
    type: Object,
    required: false,
  };
}

export function TypedArrayProp() {
  return {
    type: Array,
    required: true,
  };
}

export function OptionalTypedArrayProp() {
  return {
    type: Array,
    required: false,
    default: () => [],
  };
}