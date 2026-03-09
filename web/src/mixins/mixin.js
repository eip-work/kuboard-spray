const mixin = {
  computed: {
    $breadcrumb () {
      if (this.$options.breadcrumb) {
        return this.$options.breadcrumb.call(this)
      } else {
        return undefined
      }
    },
    $percentage () {
      if (this.$options.percentage) {
        return this.$options.percentage.call(this)
      } else {
        return 0
      }
    },
  },
  mounted () {
    if (this.$options.breadcrumb) {
      this.$store.dispatch('header/changeHeader', { key: 'breadcrumb', value: this.$breadcrumb})
      this.$store.dispatch('header/changeHeader', { key: 'percentage', value: this.$percentage})

      if (this.$options.switchToNamespaceHome) {
        this.$store.dispatch('header/changeHeader', { key: 'switchToNamespaceHome', value: true })
      } else {
        this.$store.dispatch('header/changeHeader', { key: 'switchToNamespaceHome', value: false })
      }
    }
    if (this.$options.refresh) {
      this.$store.dispatch('header/changeHeader', { key: 'refresh', value: this.$options.refresh ? { ref: this, refresh: this.$options.refresh } : undefined })
    }
  },
  beforeUnmount () {
    if (this.$options.refresh) {
      this.$store.dispatch('header/changeHeader', { key: 'refresh', value: undefined })
    }
  },
  watch: {
    '$breadcrumb': function () {
      if (this.$options.breadcrumb) {
        this.$store.dispatch('header/changeHeader', { key: 'breadcrumb', value: this.$breadcrumb})
      }
    },
    '$percentage': function () {
      if (this.$options.percentage) {
        this.$store.dispatch('header/changeHeader', { key: 'percentage', value: this.$percentage})
      }
    }
  },
}

export default mixin
