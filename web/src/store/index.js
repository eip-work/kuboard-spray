import { createStore } from 'vuex'
import header from './modules/header.js'
import cluster from "./modules/cluster.js"

let modules = {
  header: header,
  cluster: cluster
}

const store = createStore({
  modules,
})

export default store
