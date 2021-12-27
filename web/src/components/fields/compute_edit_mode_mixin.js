const mixin = {
  props: {
    antiFreeze: { type: Boolean, required: false, default: false },
  },
  inject: ['editMode'],
  computed: {
    compute_edit_mode () {
      if (this.editMode === 'view') {
        return false
      }
      if (this.antiFreeze) {
        return true
      }
      if (this.editMode === 'frozen') {
        return false
      }
      return true
    },
  },
}

export default mixin
