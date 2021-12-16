<i18n>
en:
  "findInTerminal": "Find in terminal"
  "find": "Find"
  "string": "String"
  "backward": "Find backward"
  "forward": "Find forward"
  "pattern": "Pattern"
  "caseSensitive": "Case sensitive"
  "wholeWords": "Whole worlds"
  "isHit": "Hit"
  "notFound": "Not found"
zh:
  "findInTerminal": "在终端中查找"
  "find": "查 找"
  "string": "字符串"
  "backward": "向后查找"
  "forward": "向前查找"
  "pattern": "正则表达式"
  "caseSensitive": "大小写敏感"
  "wholeWords": "完整单词"
  "isHit": "是否命中"
  "notFound": "未找到"
</i18n>

<template>
  <div style="margin-top: 1px;">
    <div v-if="visible" class="finder-form">
      <el-form style="display: inline-block;" label-width="80px" label-position="left" @submit.prevent>
        <el-form-item :label="$t('string')">
          <el-input ref="text" v-model="textToFind" clearable
            @keyup.enter="find(); $refs.text.focus()"
            style="width: calc(100vw - 1000px); max-width: 600px;"/>
        </el-form-item>
      </el-form>
      <div style="line-height: 28px; display: inline-block;">
        <el-tooltip class="tip" effect="light" :content="$t('backward')" placement="bottom-start">
          <div class="icon-wrapper" @click="directionToFind = 0; find()">
            <img :class="directionToFind === 0 ? 'icon selected' : 'icon'" :src="iconLeft" />
          </div>
        </el-tooltip>
        <el-tooltip class="tip" effect="light" :content="$t('forward')" placement="bottom-start">
          <div class="icon-wrapper" @click="directionToFind = 1; find()">
            <img :class="directionToFind === 1 ? 'icon selected' : 'icon'" :src="iconRight" />
          </div>
        </el-tooltip>
        <el-tooltip class="tip" effect="light" :content="$t('pattern')" placement="bottom-start">
          <div class="icon-wrapper" @click="isRegex = !isRegex">
            <img :class="isRegex ? 'icon selected' : 'icon'" :src="iconRegex" />
          </div>
        </el-tooltip>
        <el-tooltip class="tip" effect="light" :content="$t('caseSensitive')" placement="bottom-start">
          <div class="icon-wrapper" @click="isCaseSensitive = !isCaseSensitive">
            <img :class="isCaseSensitive ? 'icon selected' : 'icon'" :src="iconCaseSensitive" />
          </div>
        </el-tooltip>
        <el-tooltip class="tip" effect="light" :content="$t('wholeWords')" placement="bottom-start">
          <div class="icon-wrapper" @click="isWholeWord = !isWholeWord">
            <img :class="isWholeWord ? 'icon selected' : 'icon'" :src="iconWholeWord" />
          </div>
        </el-tooltip>
        <el-tooltip class="tip" effect="light" :content="$t('isHit')" placement="bottom-start" style="margin-left: 32px;">
          <div class="icon-wrapper" style="cursor: auto;">
            <img :class="found ? 'icon selected-hit' : 'icon'" :src="iconHit" />
          </div>
        </el-tooltip>
        <div style="display: inline-block; margin-left: 50px;">
          <el-button type="primary" icon="el-icon-search" @click="find" :disabled="!textToFind">{{$t('find')}}</el-button>
          <el-button type="info" icon="el-icon-close" @click="visible = false">{{$t('msg.cancel')}}</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { SearchAddon } from 'xterm-addon-search'
import iconCaseSensitive from './svg/icon-case-sensitive.svg'
import iconRegex from './svg/icon-regex.svg'
import iconWholeWord from './svg/icon-whole-word.svg'
import iconLeft from './svg/icon-left.svg'
import iconRight from './svg/icon-right.svg'
import iconHit from './svg/icon-hit.svg'

export default {
  props: {
    terminal: { type: Object, required: false, default: undefined },
  },
  data () {
    return {
      directionToFind: 1,
      textToFind: undefined,
      iconLeft,
      iconRight,
      iconCaseSensitive,
      iconRegex,
      iconWholeWord,
      iconHit,
      isCaseSensitive: false,
      isRegex: false,
      isWholeWord: false,
      found: false,
      searchAddon: undefined,
      visible: false,
    }
  },
  watch: {
    textToFind () {
      this.found = false
      this.terminal.clearSelection()
    },
    isCaseSensitive () {
      this.found = false
      this.terminal.clearSelection()
    },
    isRegex () {
      this.found = false
      this.terminal.clearSelection()
    },
    isWholeWord () {
      this.found = false
      this.terminal.clearSelection()
    },
  },
  methods: {
    show () {
      if (this.visible) {
        this.visible = false
        return
      }
      if (this.searchAddon === undefined) {
        this.searchAddon = new SearchAddon();
        this.terminal.loadAddon(this.searchAddon);
      }
      this.found = false
      this.visible = true
      this.terminal.select(0,0,0)
      setTimeout(() => {
        this.$refs.text.focus()
      }, 100)
    },
    find () {
      this.found = false
      this.terminal.focus()
      if (this.directionToFind === 1) {
        this.found = this.searchAddon.findNext(this.textToFind, {regex: this.isRegex, wholeWord: this.isWholeWord, caseSensitive: this.isCaseSensitive, incremental: false })
      } else {
        this.found = this.searchAddon.findPrevious(this.textToFind, {regex: this.isRegex, wholeWord: this.isWholeWord, caseSensitive: this.isCaseSensitive, incremental: false })
      }
      if (this.found === false) {
        this.$message.warning(`${this.$t('notFound')}: ${this.textToFind}`)
      }
    },
  }
}
</script>

<style lang="scss" scoped>
.tip {
  vertical-align: middle;
  display: inline-block;
  margin-top: 3px;
}
.icon-wrapper {
  overflow: hidden;
  display: inline-block;
  margin-left: 6px;
  width: 22px;
  height: 28px;
  cursor: pointer;
}
.icon {
  width: 16px;
  height: 16px;
  color: $--color-info;
  border: 1px solid $--color-info;
  padding: 2px;
}
.icon.selected {
  filter: drop-shadow($--color-primary 80px 0);
  transform: translateX(-80px);
}
.icon.selected-hit {
  filter: drop-shadow($--color-success 80px 0);
  transform: translateX(-80px);
}

.finder-form {
  position: fixed;
  right: 0;
  top: 0;
  background-color: #ccc;
  color: white;
  width: calc(100vw - 448px);
  text-align: left;
  padding: 3px 10px;
  height: 32px;
  user-select: none;
}
</style>
