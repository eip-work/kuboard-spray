<template>
  <div class="app_page_operations_wrapper">
    <div class="app_page_operations">
      <div class="titleWrapper">
        <span class="title">
          <slot name="title">
            {{title}}
          </slot>
        </span>
      </div>
      <div class="buttons">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    title: { type: String, required: false, default: '' }
  },
  data() {
    return {

    }
  },
  computed: {
  },
  components: { },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">
.app_page_operations_wrapper {
  z-index: 200;
  position: -webkit-sticky;
  position: sticky;
  top: 0px;
  display: flex;
  height: 52px;
  overflow: hidden;
  background-color: #f1f4fa;
}

.app_page_operations {
  box-shadow: 0 1px 8px 0 rgba(100, 100, 100, 0.8);
  padding-bottom: 10px;
  margin: auto;
  height: 33px;
  width: 100%;
  line-height: 32px;
  text-align: left;
  margin-bottom: 10px;
}

.titleWrapper {
  height: 32px;
  display: inline-block;
  margin-right: 20px;
}

.title {
  margin-left: 20px; 
  font-weight: 500;
}

.title::after {
  content: '';
  width: 0px;
  height: 12px;
  position: relative;
  border: 2px solid #234883;
  left: px;
  top: -24px;
  border-left-color: #234883;
  display: block;
  border-radius: 3px;
}

.buttons {
  vertical-align: top;
  display: inline-block;
}
</style>
