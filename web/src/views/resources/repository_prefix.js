export default function () {
  if (window.KuboardSpray.version.arch === 'arm64') {
    return "https://addons.kuboard.cn/v-kuboard-spray-arm64"
  }
  return "https://addons.kuboard.cn/v-kuboard-spray"
}