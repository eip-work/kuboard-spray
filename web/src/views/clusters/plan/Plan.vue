<i18n>
en:
  hi: kuboard
zh:
  hi: kuboard
</i18n>

<template>
  <div>
    <div class="plan">
      <div class="left">
        <div>
          <div class="localhost node">
            localhost
          </div>
        </div>
        <div>
          <div class="verticalConnection"></div>
          <div class="horizontalConnection"></div>
          <div class="verticalConnection"></div>
        </div>
        <div>
          <div class="bastion node">
            bastion
          </div>
          <div class="verticalConnection"></div>
          <div class="horizontalConnection"></div>
        </div>
      </div>
      <div class="right">
        <div class="masters">
          <div class="node" v-for="(item, index) in inventory.all.children.k8s_cluster.children.kube_control_plane.hosts" :key="'master' + index">
            {{index}}
          </div>
        </div>
        <div class="workers">
          <div class="node" v-for="(item, index) in inventory.all.children.k8s_cluster.children.kube_node.hosts" :key="'master' + index">
            {{index}}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true }
  },
  computed: {
    inventory: {
      get () {
        return this.cluster.inventory
      },
      set (v) {
        console.log(v)
      }
    }
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="scss">
.plan {
  display: flex;
  .left {
    width: 180px;
    .horizontalConnection {
      flex-grow: 1;
      margin-left: 60px;
      width: calc(100% - 60px);
      border-top: dashed 1.5px $--color-primary;
    }
    .verticalConnection {
      margin-left: 60px;
      border-left: dashed 1.5px $--color-primary;
      height: 20px;
    }
  }
  .right {
    flex-grow: 1;
    min-height: 300px;
    border: solid 1.5px $--color-primary-light-5;
    border-radius: 10px;
    margin-left: 5px;
    padding: 10px;
    .masters {
      display: flex;
      flex-wrap: wrap;
      padding-bottom: 10px;
      margin-bottom: 10px;
      border-bottom: dotted 1.5px $--color-primary-light-5;
    }
    .workers {
      display: flex;
      flex-wrap: wrap;
    }
  }
  .node {
    border: solid 1px $--border-color-light;
    border-radius: 5px;
    width: 120px;
    margin: 5px;
    padding: 15px;
    font-size: 13px;
  }
  .node:hover {
    border-color: $--color-primary;
  }
}
</style>
