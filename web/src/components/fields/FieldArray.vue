<template>
  <FieldCommon :fieldName="fieldName" :holder="holder" :prop="prop" :rules="rules" :required="required" :label="label" :placeholder="placeholder" :helpString="helpString" :helpLink="helpLink">
    <template #edit>
      <div style="flex-wrap: wrap; margin: -5px;" class="app_form_mini">
        <template v-for="(item, index) in value" :key="index + 'item'">
          <el-form-item :rules="itemRules" class="item_in_array">
            <div style="display: flex;">
              <slot name="editItem" :index="index" :item="item"></slot>
              <el-button style="margin-left: 10px;" icon="el-icon-delete" type="text" @click="obj[fieldName].splice(index, 1)"></el-button>
            </div>
          </el-form-item>
        </template>
        <div class="item_in_array" @click="addNewItem" style="cursor: pointer; border-style: dashed; flex-grow: 1;">
          <el-button type="text" icon="el-icon-plus">{{ $t('msg.add') }}</el-button>
        </div>
      </div>
    </template>
    <template #view>
      <div style="flex-wrap: wrap; margin: -5px;" class="app_form_mini">
        <template v-for="(item, index) in value" :key="index + 'item'">
          <div class="view_item_in_array">
            <slot name="viewItem" :index="index" :item="item">{{item}}</slot>
          </div>
        </template>
      </div>
    </template>
  </FieldCommon>
</template>

<script>
import compute_placeholder_mixin from './compute_placeholder_mixin.js'
import clone from 'clone'

export default {
  props: {
    holder: { type: Object, required: true },
    fieldName: { type: String, required: true },
    prop: { type: String, required: false },
    required: { type: Boolean, required: false, default: false },
    rules: { type: Array, required: false, default: () => ([])},
    itemRules: { type: Array, required: false, default: ()=> ([]) },
    placeholder: { type: String, required: false, default: undefined },
    label: { type: String, required: false, default: undefined },
    disabled: { type: Boolean, required: false, default: false },
    newItemTemplate: { required: false, default: '' },
    helpString: { type: String, required: false, default: undefined },
    helpLink: { type: String, required: false, default: undefined },
  },
  data () {
    return {
      hidePassword: true,
    }
  },
  mixins: [ compute_placeholder_mixin ],
  computed: {
    obj: {
      get () {
        return this.holder
      },
      set (v) {
        console.log(v)
      }
    },

    value: {
      get () {
        return this.holder[this.fieldName]
      },
      set (v) {
        if (v) {
          this.obj[this.fieldName] = v
        } else {
          delete this.obj[this.fieldName]
        }
      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    addNewItem() {
      this.obj[this.fieldName] = this.obj[this.fieldName] || []
      this.obj[this.fieldName].push(clone(this.newItemTemplate))
    }
  }
}
</script>

<style scoped lang="css">
.item_in_array {
  display: flex;
  margin: 5px;
  padding: 3px 12px 3px 20px;
  border-radius: 50px;
  border: solid 1px var(--el-color-primary-light-6);
  background-color: var(--el-color-primary-light-9);
  height: 28px;
}
.view_item_in_array {
  display: flex;
  margin: 5px;
  padding: 0px 12px 0px 12px;
  border-radius: 50px;
  border: solid 1px var(--el-color-primary-light-8);
  background-color: white;
  height: 28px;
}
</style>
