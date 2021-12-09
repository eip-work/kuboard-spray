import { createStore } from 'vuex'
import header from './modules/header.js'

let modules = {}

modules.header = header

const store = createStore({
  modules,
})

export default store
