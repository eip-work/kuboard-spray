const mixin = {
  computed: {
    compute_placeholder () {
      if (this.placeholder) {
        return this.placeholder
      }
      let temp = this.$t('field.' + this.fieldName + '_placeholder')
      if (temp == ('field.' + this.fieldName + '_placeholder')) {
        return this.$t('field.please_input') + this.$t(' ' +this.fieldName)
      } else {
        return temp
      }
    },
  },
}

export default mixin
