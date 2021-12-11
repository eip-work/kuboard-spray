<template>
  <el-form-item :rules="computedRules" :prop="prop ? prop + '.' + fieldName : undefined">
    <template #label>
      {{ $t('field.' + fieldName) }}
    </template>
    <div style="display: flex;" v-if="editMode !== 'view'">
      <el-select v-model.trim="obj[fieldName]" style="flex-grow: 1;" clearable
        :placeholder="$t('field.' + fieldName + '_placeholder')" @visible-change="load">
        <el-option v-for="(item, index) in options" :key="'i' + index" :value="item.value">
          {{item.label}}
        </el-option>
      </el-select>
      <slot></slot>
    </div>
    <div v-else class="app_text_mono">{{ obj[fieldName] }}</div>
  </el-form-item>
</template>

<script>

export default {
  props: {
    holder: { type: Object, required: true },
    fieldName: { type: String, required: true },
    prop: { type: String, required: false },
    required: { type: Boolean, required: false, default: false },
    rules: { type: Array, required: false, default: () => ([])},
    loadOptions: { type: Function, required: true },
  },
  data () {
    return {
      options: []
    }
  },
  inject: ['editMode'],
  computed: {
    obj: {
      get () {
        return this.holder
      },
      set (v) {
        console.log(v)
      }
    },
    computedRules () {
      let result = []
      if (this.required) {
        let message = this.$t('field.' + this.fieldName) + ' ' + this.$t('isRequiredField')
        result.push({
          validator: (rule, value, callback) => {
            if (value !== undefined && value !== '') {
              return callback()
            }
            return callback(message)
          },
          trigger: 'blur'
        })
      }
      result.push(... this.rules)
      return result
    }
  },
  components: { },
  mounted () {
  },
  methods: {
    load () {
      this.loadOptions.call(this.$parent).then(ops => {
        this.options = ops
      }).catch(e => {
        console.log(e)
      })
    }
  }
}
</script>

<style scoped lang="scss">

</style>
