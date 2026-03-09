<template>
  <EditCommon v-model="modelValue">
    <template #edit>
      <div style="flex-wrap: wrap; margin: -5px; width: 100%;" class="app_form_mini app_form_item_hide_label">
        <template v-for="(item, index) in modelValue" :key="index + 'item'">
          <el-form-item :rules="itemRules" :prop="`${prop}.${index}`" :class="isBlockItem ? 'block_item_in_array' : 'item_in_array'">
            <template v-if="isBlockItem">
              <div style="display: flex; gap: 0">
                <div style="width: calc(100% - 20px);">
                  <slot name="editItem" :index="index" :item="item"></slot>
                </div>
                <div style="width: 20px;">
                  <el-button style="margin-left: 0px;" icon="el-icon-delete" type="primary" text @click="modelValue!.splice(index, 1)">s</el-button>
                </div>
              </div>
            </template>
            <template v-else>
              <div style="display: flex; gap: 0; width: 100%;">
                <slot name="editItem" :index="index" :item="item" style="flex-grow: 1"></slot>
                <el-button style="margin-left: 5px;" icon="el-icon-delete" type="primary" text @click="modelValue!.splice(index, 1)"></el-button>
              </div>
            </template>
          </el-form-item>
        </template>
        <div class="item_in_array" @click="addNewItem" style="cursor: pointer; border-style: dashed; flex-grow: 1;">
          <el-button type="primary" text icon="el-icon-plus">{{ $t('msg.add') }}</el-button>
        </div>
        <slot name="help"></slot>
      </div>
    </template>
    <template #view>
      <div style="flex-wrap: wrap; margin: -5px;" class="app_form_mini">
        <template v-for="(item, index) in modelValue" :key="index + 'item'">
          <div class="view_item_in_array">
            <slot name="viewItem" :index="index" :item="item">{{ item }}</slot>
          </div>
        </template>
        <slot name="helpView"></slot>
      </div>
    </template>
  </EditCommon>
</template>

<script lang="ts" setup>
import clone from "clone";
import EditCommon from "./EditCommon.vue";

const props = withDefaults(
  defineProps<{
    prop: string;
    itemRules?: any[];
    newItemTemplate?: any;
    isBlockItem?: boolean;
  }>(),
  {
    itemRules: () => [],
    newItemTemplate: "",
    isBlockItem: false
  }
);

const modelValue = defineModel<any[]>();

function addNewItem() {
  modelValue.value = modelValue.value || [];
  modelValue.value.push(clone(props.newItemTemplate));
}
</script>

<style scoped lang="css">
.item_in_array {
  display: flex;
  margin: 5px;
  padding: 3px 12px 3px 20px;
  border-radius: 50px;
  border: solid 1px var(--el-color-primary-light-7);
  background-color: var(--el-color-primary-light-9);
}

.block_item_in_array {
  display: flex;
  margin: 5px;
  padding: 3px 12px 3px 20px;
  border-radius: 5px;
  border: solid 1px var(--el-color-primary-light-7);
  background-color: var(--el-color-white);
  height: auto;
}

.view_item_in_array {
  display: flex;
  align-items: center;
  margin: 3px;
  padding: 0px 6px 0px 12px;
  border-radius: 4px;
  /* border: solid 1px var(--el-color-primary-light-8); */
  background-color: white;
  height: 24px;
  position: relative;
  
}

.view_item_in_array::before {
  content: '-';
  position: absolute;
  left: 0;
  font-weight: 500;
}
</style>
