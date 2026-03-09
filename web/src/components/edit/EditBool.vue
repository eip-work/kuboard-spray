<template>
  <EditCommon v-model="modelValue">
    <template #edit>
      <el-switch v-model="modelValue" :disabled="disabled"></el-switch>
      <div class="desc">
        <slot name="edit_desc"></slot>
      </div>
    </template>
    <template #view>
      <el-switch v-model="modelValue" disabled></el-switch>
      <div class="desc">
        <slot name="view_desc"></slot>
      </div>
    </template>
  </EditCommon>
</template>

<script lang="ts" setup>
import EditCommon from "./EditCommon.vue"

withDefaults(defineProps<{
  disabled?: boolean;
}>(), {
  disabled: false,
})

const modelValue = defineModel<boolean>();

</script>

<style scoped lang="css">
.desc {
  color: var(--el-text-color-secondary);
  display: inline-block;
  vertical-align: top;
  margin-left: 10px;
}
</style>
