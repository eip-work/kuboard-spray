import ControlBar from './page/ControlBar.vue'
import ConfigSection from './ConfigSection.vue'
import CopyToClipBoard from './clipboard/CopyToClipBoard.vue'
import ConfirmButton from './page/ConfirmButton.vue'
import PreviewYaml from './PreviewYaml.vue'
import PangeeClusterLink from './PangeeClusterLink.vue'
import PangeeClusterTime from './PangeeClusterTime.vue'
import EditCommon from "./edit/EditCommon.vue"
import EditString from "./edit/EditString.vue"
import EditBool from "./edit/EditBool.vue"
import EditNumber from './edit/EditNumber.vue'
import EditArray from "./edit/EditArray.vue"
import EditRadio from "./edit/EditRadio.vue"
import EditSelect from "./edit/EditSelect.vue"

export default {
  install(app) {
    app.component("EditCommon", EditCommon)
    app.component("EditString", EditString)
    app.component("EditBool", EditBool)
    app.component("EditNumber", EditNumber)
    app.component("EditArray", EditArray)
    app.component("EditRadio", EditRadio)
    app.component("EditSelect", EditSelect)
    app.component('ControlBar', ControlBar)
    app.component('ConfigSection', ConfigSection)
    app.component('CopyToClipBoard', CopyToClipBoard)
    app.component('ConfirmButton', ConfirmButton)
    app.component('PreviewYaml', PreviewYaml)
    app.component('PangeeClusterLink', PangeeClusterLink)
    app.component('PangeeClusterTime', PangeeClusterTime)
  } 
}