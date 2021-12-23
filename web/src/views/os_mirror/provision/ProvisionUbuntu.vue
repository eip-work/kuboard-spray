<i18n>
en:
  special_time: Frequncy to sync
  KB: KiloBytes/Thread
  kb: about { kb } Mb/s per thread
  threadCount: Thread Count
  totalkb: about { kb } Mb/s in total
  releases: Releases included
zh:
  special_time: 从上游同步的频率
  KB: 千字节每秒每线程
  kb: 大约 { kb } Mb/s 每线程
  threadCount: 个下载线程
  totalkb: 总共大约 { kb } Mb/s
  release: 包含的版本
</i18n>

<template>
  <div>
    <FieldNumber :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_server_port" required></FieldNumber>
    <FieldString :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_dir" required></FieldString>
    <!-- <FieldSelect :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_ubuntu_mirror" required :loadOptions="loadOptions"></FieldSelect> -->
    <el-form-item :label="$t('field.apt_mirror_ubuntu_mirror')" prop="inventory.all.children.target.vars.apt_mirror_ubuntu_mirror" required>
      <el-select v-model="apt_mirror_ubuntu_mirror" style="width: 100%;" v-if="editMode !== 'view'" :placeholder="$t('field.apt_mirror_ubuntu_mirror')">
        <el-option v-for="(url, index) in mirrorOptions" :key="'url' + index" :value="url">
          {{ url }}
        </el-option>
      </el-select>
      <span v-else class="app_text_mono">{{apt_mirror_ubuntu_mirror}}</span>
    </el-form-item>
    <el-form-item :label="$t('release')" required message="Required" prop="inventory.all.children.target.vars.apt_mirror_repos">
      <el-checkbox-group v-model="apt_mirror_repos" :disabled="editMode === 'view'">
        <el-checkbox v-for="(value, key) in releases" :key="key + 'release'" :label="key">
          <span class="app_text_mono">
            Ubuntu {{value}} - {{key}}
          </span>
        </el-checkbox>
      </el-checkbox-group>
    </el-form-item>
    <FieldBool :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_schedule_updates"></FieldBool>
    <FieldRadio v-if="inventory.all.children.target.vars.apt_mirror_schedule"
      :disabled="!inventory.all.children.target.vars.apt_mirror_schedule_updates"
      :holder="inventory.all.children.target.vars.apt_mirror_schedule[0]"
      fieldName="special_time" :label="$t('special_time')"
      :options="['daily', 'weekly', 'monthly']"></FieldRadio>
    <FieldBool :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_populate_repos"></FieldBool>
    <FieldBool :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_enable_limit_rate"></FieldBool>
    <FieldNumber :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_limit_rate" 
      v-if="inventory.all.children.target.vars.apt_mirror_enable_limit_rate">
      <template #append>{{ $t('KB') }}</template>
      <div class="placeholder">{{ $t('kb', { kb: inventory.all.children.target.vars.apt_mirror_limit_rate * 8 / 1000 }) }}</div>
    </FieldNumber>
    <FieldNumber :holder="inventory.all.children.target.vars" prop="inventory.all.children.target.vars" fieldName="apt_mirror_nthreads"
      v-if="inventory.all.children.target.vars.apt_mirror_enable_limit_rate">
      <template #append>{{ $t('threadCount') }}</template>
      <div class="placeholder">{{ $t('totalkb', { kb: inventory.all.children.target.vars.apt_mirror_limit_rate * 8 * inventory.all.children.target.vars.apt_mirror_nthreads / 1000 }) }}</div>
    </FieldNumber>
  </div>
</template>

<script>
export default {
  props: {
    os_mirror: { type: Object, required: true },
  },
  data() {
    return {
      releases: {
        impish: '21.10',
        hirsute: '21.04',
        focal: '20.04',
        bionic: '18.04',
        xenial: '16.04',
        trusty: '14.04',
      },
      mirrorOptions: [
        'https://repo.huaweicloud.com/ubuntu/',
        'https://mirrors.aliyun.com/ubuntu/',
        'https://mirrors.tuna.tsinghua.edu.cn/ubuntu/',
        'http://cn.archive.ubuntu.com/ubuntu/',
        'https://mirrors.cloud.tencent.com/ubuntu/',
        'http://mirror.aliyun.com/ubuntu/'
      ]
    }
  },
  inject: ['editMode'],
  computed: {
    inventory: {
      get () { return this.os_mirror.inventory },
      set () {},
    },
    apt_mirror_ubuntu_mirror: {
      get () {
        if (this.inventory.all.children.target.vars.apt_mirror_ubuntu_mirror_protocol) {
          return this.inventory.all.children.target.vars.apt_mirror_ubuntu_mirror_protocol + this.inventory.all.children.target.vars.apt_mirror_ubuntu_mirror
        }
        return ''
      },
      set (v) {
        if (v) {
          let splited = v.split("://")
          this.inventory.all.children.target.vars.apt_mirror_ubuntu_mirror_protocol = splited[0] + "://"
          this.inventory.all.children.target.vars.apt_mirror_ubuntu_mirror = splited[1]
        } else {
          this.inventory.all.children.target.vars.apt_mirror_ubuntu_mirror_protocol = undefined
          this.inventory.all.children.target.vars.apt_mirror_ubuntu_mirror = undefined
        }
      }
    },
    apt_mirror_repos: {
      get () {
        let temp = {}
        for (let i in this.inventory.all.children.target.vars.apt_mirror_repos) {
          let repo = this.inventory.all.children.target.vars.apt_mirror_repos[i]
          for (let key in this.releases) {
            if (repo.indexOf(key) > 0) {
              temp[key] = true
            }
          }
        }
        let result = []
        for (let key in temp) {
          result.push(key)
        }
        return result
      },
      set (v) {
        let temp = []
        for (let release of v) {
          temp.push(`deb-amd {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release} main restricted universe multiverse`)
          temp.push(`deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-backports main restricted universe multiverse`)
          temp.push(`deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-security main restricted universe multiverse`)
          temp.push(`deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-updates main restricted universe multiverse`)
          temp.push(`deb-amd64 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-proposed main restricted universe multiverse`)
          temp.push(`deb-i386 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release} main restricted universe multiverse`)
          temp.push(`deb-i386 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-backports main restricted universe multiverse`)
          temp.push(`deb-i386 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-security main restricted universe multiverse`)
          temp.push(`deb-i386 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-updates main restricted universe multiverse`)
          temp.push(`deb-i386 {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-proposed main restricted universe multiverse`)
          temp.push(`deb-src {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release} main restricted universe multiverse`)
          temp.push(`deb-src {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-backports main restricted universe multiverse`)
          temp.push(`deb-src {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-security main restricted universe multiverse`)
          temp.push(`deb-src {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-updates main restricted universe multiverse`)
          temp.push(`deb-src {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release}-proposed main restricted universe multiverse`)
          temp.push(`deb {{ apt_mirror_ubuntu_mirror_protocol }}{{ apt_mirror_ubuntu_mirror }} ${release} main/debian-installer multiverse/debian-installer restricted/debian-installer universe/debian-installer`)
        }
        this.inventory.all.children.target.vars.apt_mirror_repos = temp
      }
    }
  },
  watch: {
    'os_mirror': function() {
      this.initDefaults()
    },
  },
  components: { },
  mounted () {
    this.initDefaults()
  },
  methods: {
    initDefaults () {
      let vars = this.inventory.all.children.target.vars
      vars.apt_mirror_schedule = vars.apt_mirror_schedule || [ 
        {
          name: 'apt-mirror updates',
          special_time: 'daily',
          job: '/usr/bin/apt-mirror > {{ apt_mirror_dir }}/var/cron.log',
          cron_file: 'apt_mirror',
        } 
      ]
      vars.apt_mirror_repos = vars.apt_mirror_repos || []
      if (vars.apt_mirror_populate_repos === undefined) {
        vars.apt_mirror_populate_repos = true
      }
      if (vars.apt_mirror_enable_limit_rate === undefined) {
        vars.apt_mirror_enable_limit_rate = false
      }
      vars.apt_mirror_limit_rate = vars.apt_mirror_limit_rate || 125
      vars.apt_mirror_nthreads = vars.apt_mirror_nthreads || 10
      vars.apache2_virtual_hosts = vars.apache2_virtual_hosts || [
        {
          documentroot: 'var/www/html',
          default_site: true,
          port: '{{ apt_mirror_server_port }}',
          serveradmin: 'webmaster@localhost',
          servername: '',
        }
      ],
      vars.apt_mirror_server_port = vars.apt_mirror_server_port || 80
    },
    async loadOptions() {
      return [
        { label: 'https://repo.huaweicloud.com/ubuntu/', value: 'https://repo.huaweicloud.com/ubuntu/' },
        { label: 'https://mirrors.aliyun.com/ubuntu/', value: 'https://mirrors.aliyun.com/ubuntu/' },
        { label: 'https://mirrors.tuna.tsinghua.edu.cn/ubuntu/', value: 'https://mirrors.tuna.tsinghua.edu.cn/ubuntu/' },
        { label: 'http://cn.archive.ubuntu.com/ubuntu/ ', value: 'http://cn.archive.ubuntu.com/ubuntu/' },
        { label: 'https://mirrors.cloud.tencent.com/ubuntu/', value: 'https://mirrors.cloud.tencent.com/ubuntu/' },
      ]
    },
    async loadSpecialTimeOptions () {
      return [
        { label: this.$t('daily'), value: 'daily' },
        { label: this.$t('weekly'), value: 'weekly' },
        { label: this.$t('monthly'), value: 'monthly' },
      ]
    }
  }
}
</script>

<style scoped lang="scss">
.placeholder {
  font-size: 12px;
  color: $--color-text-secondary;
}
</style>
