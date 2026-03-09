<template>
  <EditCommon>
    <template #edit="scope">
      <div style="display: flex; flex-grow: 1;">
        <el-select v-model="modelValue" style="flex-grow: 1;" :clearable="clearable" :disabled="disabled"
          :allowCreate="allowCreate" :filterable="filterable" :placeholder="scope.placeholder"
          @visible-change="load($event)" :loading="loading" :multiple="multiple">
          <el-option v-for="(item, index) in options" :key="'i' + index" :value="item.value" :label="item.label"
            :disabled="item.disabled">
            {{ item.label }}
          </el-option>
        </el-select>
        <slot name="edit"></slot>
      </div>
    </template>
    <template #view="scope">
      <slot name="view">
        <span v-if="modelValue !== undefined && modelValue !== ''">
          {{ compute_display_value }}
        </span>
        <span v-else class="field_placeholder">{{ scope.placeholder }}</span>
        <slot name="view_append"></slot>
      </slot>
    </template>
  </EditCommon>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue";
import EditCommon from "./EditCommon.vue"

type LoadOptionFunc = () => Promise<{ label: string; value: string, disabled?: boolean }[]>;

const props = withDefaults(defineProps<{
  loadOptions: LoadOptionFunc;
  disabled?: boolean;
  clearable?: boolean;
  allowCreate?: boolean;
  multiple?: boolean;
  filterable?: boolean;
}>(), {
  disabled: false,
  clearable: false,
  allowCreate: false,
  multiple: false,
  filterable: false,
})

const modelValue = defineModel<string>();

const options = ref<{ label: string; value: string, disabled?: boolean }[]>([])
const loading = ref(false)
const compute_display_value = computed(() => {
  for (let i in options.value) {
    if (options.value[i].value === modelValue.value) {
      return options.value[i].label
    }
  }
  return modelValue.value;
})

onMounted(() => {
  load(true)
})

function load(flag: boolean) {
  if (flag) {
    loading.value = true
    props.loadOptions.apply(modelValue).then((ops: any) => {
      options.value = ops
      loading.value = false
    }).catch((e: any) => {
      console.log(e)
      loading.value = false
    })
  }
}
</script>

<style scoped lang="css">
.field_placeholder {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}
</style>
