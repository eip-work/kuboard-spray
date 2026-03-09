<template>
  <EditCommon v-model="modelValue">
    <template #edit="scope">
      <el-input v-model.number="modelValue" :placeholder="scope.placeholder">
        <template #append v-if="$slots.append">
          <slot name="append"></slot>
        </template>
      </el-input>
      <div class="desc" v-if="$slots.edit_desc">
        <slot name="edit_desc"></slot>
      </div>
    </template>
    <template #view>
      {{ modelValue }} <span class="app_text_mono" style="font-size: 13px;">
        <slot name="append"></slot>
      </span>
      <div class="desc" v-if="$slots.view_desc">
        <slot name="view_desc"></slot>
      </div>
    </template>
  </EditCommon>
</template>

<script lang="ts" setup>
import EditCommon from "./EditCommon.vue"

const modelValue = defineModel<number>();
</script>

<style scoped lang="css">
.field_placeholder {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}

.desc {
  color: var(--el-text-color-secondary);
  display: inline-block;
  vertical-align: middle;
  margin-left: 10px;
}
</style>
