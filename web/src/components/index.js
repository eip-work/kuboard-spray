import FieldString from './fields/FieldString.vue'
import FieldBool from './fields/FieldBool.vue'
import FieldNumber from './fields/FieldNumber.vue'
import FieldRadio from './fields/FieldRadio.vue'
import FieldSelect from './fields/FieldSelect.vue'

export default {
  install(app) {
    app.component('FieldString', FieldString)
    app.component('FieldBool', FieldBool)
    app.component('FieldNumber', FieldNumber)
    app.component('FieldRadio', FieldRadio)
    app.component('FieldSelect', FieldSelect)
  } 
}