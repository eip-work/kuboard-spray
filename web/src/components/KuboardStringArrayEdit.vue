<i18n locale="en" lang="yaml">
  hi: kuboard
</i18n>

<i18n locale="cn" lang="yaml">
  hi: kuboard
</i18n>

<template>
  <div class="app_form_mini" style="margin: 0 -5px; display: flex; flex-wrap: wrap">
    <template v-if="holderData[fieldName0] && holderData[fieldName0][fieldName1]">
      <el-form-item class="arrayItem" :style="`width: ${width}; flex-grow: ${flexGrow};`" v-for="(item, index) in holderData[fieldName0][fieldName1]" :key="'item' + index"
        :prop="`${prop}.${index}`" :rules="rules || [{ required: true, message: message, trigger: 'blur' }]">
        <div style="display: flex;">
          <el-input style="flex-grow: 1; margin-right: 10px;" v-model.trim="holderData[fieldName0][fieldName1][index]" :placeholder="placeholder"></el-input>
          <el-button icon="el-icon-delete" type="primary" text @click="holderData[fieldName0][fieldName1].splice(index, 1)"></el-button>
        </div>
      </el-form-item>
    </template>
    <div class="arrayItem" :style="`width: ${width}; flex-grow: ${flexGrow}; cursor: pointer; border-style: dashed;`" @click="addItem">
      <el-button type="primary" text icon="el-icon-plus" size="default">{{addButtonText || '添 加'}}</el-button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    holder: { type: Object, required: true },
    fieldName0: { type: String, required: true },
    fieldName1: { type: String, required: true },
    prop: { type: String, required: false, default: undefined },
    placeholder: { type: String, required: false, default: undefined },
    message: { type: String, required: false, default: '不能为空' },
    width: { type: String, required: false, default: '160px' },
    flexGrow: { type: String, required: false, default: '0' },
    rules: { type: Array, required: false, default: undefined },
    addButtonText: { type: String, required: false, default: undefined },
  },
  data() {
    return {

    }
  },
  computed: {
    holderData: {
      get() {
        return this.holder
      },
      set () {

      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    addItem () {
      let fieldName0 = this.fieldName0
      let fieldName1 = this.fieldName1
      this.holderData[fieldName0] = this.holderData[fieldName0] || {}
      this.holderData[fieldName0][fieldName1] = this.holderData[fieldName0][fieldName1] || []
      this.holderData[fieldName0][fieldName1].push('')
    }
  }
}
</script>

<style scoped lang="css">
.arrayItem {
  margin: 5px;
  padding: 3px 12px 3px 20px;
  border-radius: 50px;
  border: solid 1px var(--el-color-primary-light-7);
  background-color: var(--el-color-primary-light-9);
  height: 28px;
}
</style>
