<template>
  <EditCommon v-model="modelValue">
    <template #edit>
      <el-radio-group v-model="modelValue" :disabled="disabled">
        <template v-for="(option, index) in options" :key="index">
          <el-radio-button :label="option.label" :value="option.value" :disabled="option.disabled">
            {{ option.label }}
          </el-radio-button>
        </template>
      </el-radio-group>
      <slot name="edit"></slot>
    </template>
    <template #view>
      <template v-for="(option, index) in options" :key="index">
        <el-button v-if="modelValue === option.value" :label="option.label" type="primary" disabled>
          {{ option.label }}
        </el-button>
      </template>
      <slot name="view"></slot>
    </template>
  </EditCommon>
</template>

<script lang="ts" setup>
import EditCommon from "./EditCommon.vue"

withDefaults(defineProps<{
  options: {
    label: string;
    value: any;
    disabled?: boolean;
  }[];
  disabled?: boolean;
}>(), {
  disabled: false,
})

const modelValue = defineModel<string>();

</script>

<style scoped lang="css"></style>
