import FieldCommon from './fields/FieldCommon.vue'
import FieldArray from './fields/FieldArray.vue'
import FieldString from './fields/FieldString.vue'
import FieldBool from './fields/FieldBool.vue'
import FieldNumber from './fields/FieldNumber.vue'
import FieldRadio from './fields/FieldRadio.vue'
import FieldSelect from './fields/FieldSelect.vue'
import ControlBar from './page/ControlBar.vue'
import ConfigSection from './ConfigSection.vue'
import CopyToClipBoard from './clipboard/CopyToClipBoard.vue'
import ConfirmButton from './page/ConfirmButton.vue'
import PreviewYaml from './PreviewYaml.vue'
import KuboardSprayLink from './KuboardSprayLink.vue'
import KuboardSprayTime from './KuboardSprayTime.vue'

export default {
  install(app) {
    app.component('FieldCommon', FieldCommon)
    app.component('FieldArray', FieldArray)
    app.component('FieldString', FieldString)
    app.component('FieldBool', FieldBool)
    app.component('FieldNumber', FieldNumber)
    app.component('FieldRadio', FieldRadio)
    app.component('FieldSelect', FieldSelect)
    app.component('ControlBar', ControlBar)
    app.component('ConfigSection', ConfigSection)
    app.component('CopyToClipBoard', CopyToClipBoard)
    app.component('ConfirmButton', ConfirmButton)
    app.component('PreviewYaml', PreviewYaml)
    app.component('KuboardSprayLink', KuboardSprayLink)
    app.component('KuboardSprayTime', KuboardSprayTime)
  } 
}