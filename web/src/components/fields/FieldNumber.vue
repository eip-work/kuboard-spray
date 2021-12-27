<template>
  <FieldCommon :fieldName="fieldName" :holder="holder" :prop="prop" :rules="rules" :required="required" :label="label" :placeholder="placeholder">
    <template #edit>
      <el-input v-if="editMode !== 'view'" v-model.number="obj[fieldName]" :placeholder="compute_placeholder">
        <template #append v-if="$slots.append"><slot name="append"></slot></template>
      </el-input>
    </template>
    <template #view>
      {{ obj[fieldName] }} <span class="app_text_mono" style="font-size: 13px;"><slot name="append"></slot></span>
    </template>
  </FieldCommon>
</template>

<script>
import compute_placeholder_mixin from './compute_placeholder_mixin.js'

export default {
  props: {
    holder: { type: Object, required: true },
    fieldName: { type: String, required: true },
    prop: { type: String, required: false },
    required: { type: Boolean, required: false, default: false },
    rules: { type: Array, required: false, default: () => ([])},
    placeholder: { type: String, required: false, default: undefined },
    label: { type: String, required: false, default: undefined },
  },
  data () {
    return {

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
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.field_placeholder {
  color: $--color-text-placeholder;
  font-size: 12px;
}
</style>
