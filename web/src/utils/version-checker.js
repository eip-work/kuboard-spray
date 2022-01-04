import axios from 'axios'

async function init () {
    axios.get(`https://addons.kuboard.cn/kuboard-spray-latest-version.json?r=${Math.random()}`, {
      timeout: 2000,
      withCredentials: false,
    }).then(resp => {
      let installed = window.KuboardSpray.version.version
      let arch = installed.slice(installed.lastIndexOf('-') + 1)
      let channel = 'latest'
      if (installed.indexOf('beta') > 0) channel = 'beta'
      if (installed.indexOf('alpha') > 0) channel = 'alpha'
      let currentVersion = resp.data[arch].channel[channel]
      
      let newVersion = currentVersion.version
      console.log('installed version: ' + installed)
      console.log('      new version: ' + newVersion)
      if (installed !== newVersion) {
        if (localStorage.getItem('IgnoreThisVersion') === newVersion) {
          console.log('  Ignored version: ' + newVersion)
          return
        }
        setTimeout(() => {
            window.$kuboardVersionCheckNotification = window.VueAppComponent.$notify({
              title: '检测到新版本',
              dangerouslyUseHTMLString: true,
              message: `<a href="https://github.com/eip-work/kuboard-spray/releases" target=_blank class="app_notify" style="width: 100%;" onclick="window.$kuboardVersionCheckNotification.close()">
                        <div class="app_text_mono">已经安装的版本是 <span style="color: var(--el-color-success); font-weight: bold;">${installed} </span></div>
                        <div class="app_text_mono">当前最新的版本是 <span style="color: var(--el-color-danger); font-weight: bold;">${newVersion} </span></div>
                        <div style="text-align: right; margin-top: 10px; color: var(--el-color-primary); font-weight: bold;">
                          <i class="el-icon-thumb" style="margin-right: 5px;"></i>点击此处前往升级
                        </div>
                        </a>
                        <div onclick="ignoreThisVersion('${newVersion}')"
                          style="margin-left: -30px; padding: 10px 20px 0px 0;" class="app_text_description app_notify">
                          ${window.VueAppComponent.$t('upgrade.ignoreThisVersion')}
                        </div>`,
              position: 'bottom-left',
              type: 'info',
              duration: 0,
              'custom-class': 'kuboard-spray-upgrade-notification',
            });
        }, 5000)
      }
    }).catch(e => {
      console.log('查询新版本失败: ', e)
    })

}

export default {init}

window.ignoreThisVersion = function(version) {
  localStorage.setItem('IgnoreThisVersion', version)
  window.VueAppComponent.$message.success(window.VueAppComponent.$t('upgrade.ignored', {version}))
  window.$kuboardVersionCheckNotification.close()
}

