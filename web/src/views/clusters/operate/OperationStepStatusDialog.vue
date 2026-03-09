<i18n>
en:
  startTime: Check Time
  duration: Duration
  sequence: Sequence
  node: Node
zh:
  startTime: 状态检查执行时间
  duration: 耗时
  sequence: 序号
  node: 节点
</i18n>

<template>
  <el-dialog v-model="visible" top="5vh" width="90vw">
    <template #header>
      <div>
        <span class="dialog-title">
          {{ t("startTime") }}:
        </span>
        <span class="dialog-value app_text_mono">
          {{ formatTime(status?.data.start_time) }}
        </span>
        <span class="dialog-title" style="margin-left: 30px">
          {{ t("duration") }}:
        </span>
        <span class="dialog-value app_text_mono">
          {{ status?.data.duration }}ms
        </span>
        <span style="margin-left: 50px">
          <el-tag :type="computedStatus.success == computedStatus.total ? 'success' : 'danger'" size="">
            {{ computedStatus.success }}
          </el-tag>
          /
          <el-tag style="margin-right: 5px;" size="">
            {{ computedStatus.total }}
          </el-tag>
        </span>
      </div>
    </template>
    <div v-if="visible" class="check_step_status_table">
      <el-table :data="rows" style="width: 100%" size="">
        <el-table-column :label="t('sequence')" type="index" width="50">
        </el-table-column>
        <el-table-column :label="t('node')" width="120">
          <template #default="scope">
            {{ scope.row._nodeName }}
          </template>
        </el-table-column>
        <template v-for="column in columns" :key="column.name">
          <el-table-column width="150">
            <template #header>
              <el-popover width="800" trigger="click">
                <template #reference>
                  <el-link type="primary">
                    <el-icon>
                      <ElIconPosition />
                    </el-icon>
                    <span style="margin-left: 5px;">
                      {{ column.name }}
                    </span>
                  </el-link>
                </template>
                <el-descriptions border :column="1">
                  <el-descriptions-item label="名称"> {{ column.name }} </el-descriptions-item>
                  <el-descriptions-item label="描述"> {{ column.description[locale] }} </el-descriptions-item>
                  <el-descriptions-item label="cmd">
                    <div class="shell_response">
                      <div class="app_text_mono">
                        {{ column.cmd }}
                      </div>
                    </div>
                  </el-descriptions-item>
                </el-descriptions>

              </el-popover>
            </template>
            <template #default="scope">
              <el-tag v-if="scope.row[column.key].stdout == '0'" type="success" effect="dark" round size="small">
                <el-icon>
                  <ElIconCheck />
                </el-icon>
              </el-tag>
              <el-popover v-else width="75vw" effect="dark">
                <template #reference>
                  <el-tag type="danger" effect="dark" round size="small">
                    <el-icon>
                      <ElIconClose />
                    </el-icon>
                  </el-tag>
                </template>
                <div class="shell_response">
                  <div class="app_text_mono">
                    {{ scope.row[column.key].stdout || scope.row[column.key].msg }}
                  </div>
                </div>
              </el-popover>
            </template>
          </el-table-column>
        </template>
      </el-table>
    </div>
  </el-dialog>
</template>

<script>
import yaml from 'js-yaml';
import dayjs from 'dayjs';
import clone from "clone";

export default {
  data() {
    return {
      status: undefined,
      visible: false
    }
  },
  computed: {
    computedStatus() {
      if (this.status === undefined) {
        return { loading: true };
      }
      let total = 0;
      let success = 0;
      for (let i in this.status.data.result) {
        let node = this.status.data.result[i];
        for (let j in node) {
          let action = node[j];
          if (!action.skipped) {
            total++;
            if (action.stdout == "0") {
              success++;
            }
          }
        }
      }
      return {
        loading: false,
        success, total
      };
    },
    columns() {
      let columns = {};
      for (let i in this.status.data.result) {
        let node = this.status.data.result[i];
        for (let j in node) {
          let _metadata = yaml.load(j);
          if (node[j].cmd) {
            _metadata.cmd = node[j].cmd;
          }
          _metadata.key = j;

          let cmd = _metadata.cmd;
          if (columns[_metadata.name] && columns[_metadata.name].cmd) {
            cmd = columns[_metadata.name].cmd;
          }

          if (columns[_metadata.name] == undefined) {
            columns[_metadata.name] = _metadata;
          }
          columns[_metadata.name].cmd = cmd;
        }
      }
      let temp = [];
      for (let k in columns) {
        temp.push(columns[k])
      }
      return temp;
    },
    rows() {
      let tmp = [];
      for (let i in this.status.data.result) {
        let t = clone(this.status.data.result[i]);
        let skipped = true;
        for (let j in t) {
          if (!t[j].skipped) {
            skipped = false;
          }
        }
        if (!skipped) {
          t._nodeName = i;
          tmp.push(t);
        }
      }
      return tmp;
    }
  },
  methods: {
    formatTime(time) {
      return dayjs(time).format("YYYY-MM-DD HH:mm:ss")
    },
    show(status) {
      this.status = status;
      this.visible = true;
    }
  }
}
</script>

<style lang="scss">
.check_step_status_table {
  width: 100%;
  height: calc(90vh - 100px);

  .cell {
    text-align: center;
  }

  tr {
    height: 42px;
  }
}
</style>

<style lang="scss" scoped>
.dialog-title {
  color: var(--el-text-color-secondary);
  font-size: 14px
}

.dialog-value {
  font-weight: bold;
  margin-left: 5px;
}

.shell_response {
  display: flex;
  max-height: 60vh;
  overflow: hidden;
  overflow-y: auto;
  background-color: var(--el-color-black);
  color: var(--el-fill-color);
  padding: 5px;

  div {
    flex: 0 1 auto;
    overflow-x: auto;
    white-space: pre;
    word-break: keep-all;
    overflow-wrap: normal;
  }
}
</style>