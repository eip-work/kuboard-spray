<i18n>
en:
  please: 'Please input '
  isRequiredField: ' is required.'
zh:
  please: 请输入
  isRequiredField: 为必填字段
</i18n>

<template>
  <el-form-item :rules="computedRules" :prop="prop" :label-width="compute_label_width">
    <template #label>
      {{ label }}
    </template>
    <template v-if="compute_edit_mode">
      <slot name="edit" :placeholder="compute_placeholder"></slot>
    </template>
    <div v-else class="app_text_mono">
      <slot name="view" :placeholder="compute_placeholder">
        <span v-if="modelValue !== undefined && modelValue !== ''">
          {{ compute_display_value }}
        </span>
        <span v-else class="field_placeholder">{{ compute_placeholder }}</span>
      </slot>
    </div>
    <div v-if="compute_edit_mode && (helpString || helpLink)" style="color: var(--el-text-color-secondary)">
      {{ helpString }}
      <a v-if="helpLink" :href="helpLink" target="_blank">
        <el-icon style="vertical-align: middle;"><el-icon-link /></el-icon>
        帮助
      </a>
    </div>
  </el-form-item>
</template>

<script setup lang="ts">
import { inject, computed, ComputedRef } from 'vue';
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n({
  inheritLocale: true,
  useScope: 'local'
})

const props = withDefaults(defineProps<{
  prop?: string;
  required?: boolean;
  rules?: any[];
  label?: string;
  placeholder?: string;
  antiFreeze?: boolean;
  readOnly?: boolean;
  labelWidth?: string;
  helpString?: string;
  helpLink?: string;
}>(), {
  labelWidth: "",
  required: false,
  readOnly: false,
  antiFreeze: false
})

const modelValue = defineModel();

const editMode = inject<ComputedRef<"view" | "edit" | "create" | "frozen">>('editMode')

const defaultLabelWidth = inject<string>("defaultLabelWidth", "")

const compute_edit_mode = computed(() => {
  if (props.readOnly) {
    return false
  }
  if (editMode?.value === 'view') {
    return false
  }
  if (props.antiFreeze) {
    return true
  }
  if (editMode?.value === 'frozen') {
    return false
  }
  return true
})

const compute_label_width = computed(() => {
  if (props.labelWidth != undefined && props.labelWidth != "") {
    return props.labelWidth;
  }
  return defaultLabelWidth || "";
})

const compute_display_value = computed(() => {
  if (modelValue.value !== undefined && modelValue.value !== '') {
    return modelValue.value
  }
  if (props.placeholder) {
    return props.placeholder
  }
  return t("please") + props.label
})

const compute_placeholder = computed(() => {
  if (props.placeholder) {
    return props.placeholder
  }
  return t("please") + props.label
})

const computedRules = computed(() => {
  let result = []
  if (editMode?.value === 'frozen' && !props.antiFreeze) {
    return []
  }
  if (props.required) {
    let fieldName = props.prop || props.label || '';
    fieldName = fieldName.slice((fieldName.lastIndexOf('.') || 0) + 1);
    let message = fieldName + t('isRequiredField');
    result.push({
      required: true,
      validator: (rule: any, value: any, callback: Function) => {
        console.log("Validating required field:", value);
        if (value !== undefined && value !== '' && value !== null) {
          return callback()
        }
        return callback(message)
      },
      message: message,
      trigger: 'blur'
    })
  }
  if (props.rules) {
    result.push(...props.rules)
  }
  return result
})
</script>

<style scoped lang="css">
.field_placeholder {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}
</style>