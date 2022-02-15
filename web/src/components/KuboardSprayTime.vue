<i18n>
en:
  ago: " ago"
zh:
  ago: 以前
</i18n>

<template>
  <el-tooltip class="item" effect="dark" placement="top-start">
    <template #content>
      <div class="app_text_mono">{{formattedTime}}</div>
    </template>
    <span class="text" :style="isFuture > 0 ? 'color: var(--el-color-primary); font-weight: bold' : ''">{{distance}}{{$t('ago')}}</span>
  </el-tooltip>
</template>

<script>
import { differenceInMinutes, format, formatDistanceToNow, parseISO } from 'date-fns'
import { zhCN } from 'date-fns/locale'

export default {
  props: {
    time: { type: [String, Date], required: false, default: () => { return undefined } },
    compareTo: { type: [String, Date], required: false, default: undefined }
  },
  data() {
    return {
      interval: undefined,
    }
  },
  computed: {
    formattedTime () {
      if (this.time) {
        return format(this.timeComputed, 'yyyy-MM-dd HH:mm:ss')
      } else {
        return ''
      }
    },
    timeComputed () {
      if (typeof(this.time) === 'string') {
        return parseISO(this.time)
      }
      return this.time
    },
    distance () {
      if (this.time) {
        return formatDistanceToNow(this.timeComputed, { locale: zhCN })
      } else {
        return 'N/A'
      }
    },
    isFuture () {
      return this.differenceInMinutes(this.timeComputed, new Date())
    },
  },
  components: { },
  mounted () {
    this.interval = setInterval(() => { this.$forceUpdate() }, 10000)
  },
  beforeDestory () {
    clearInterval(this.interval)
  },
  methods: {
    differenceInMinutes,
    format,
    formatDistanceToNow,
  }
}
</script>

<style scoped lang="scss">
.text {
  border-bottom: dashed 1px var(--el-text-color-secondary);
  border-radius: 0px;
}
</style>
