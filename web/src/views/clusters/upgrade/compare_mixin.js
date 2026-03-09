const mixin = {
  computed: {
    tableData () {
      let result = []
      let res = this.cluster.resourcePackage.data.dependency

      for (let component of res) {
        result.push({name:component.name, version:component.version})
      }

      return result
    },
  },
}

export default mixin
