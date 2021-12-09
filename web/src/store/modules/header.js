const state = {
  percentage: 0,
  breadcrumb: [],
  refresh: undefined,
  errorMsg: undefined,
}

const mutations = {
  CHANGE_HEADER: (state, { key, value }) => {
    state[key] = value
  }
}

const actions = {
  changeHeader({ commit }, data) {
    commit('CHANGE_HEADER', data)
  }
}

const getters = {
  breadcrumb: state => {
    return state.breadcrumb
  },
  percentage: state => {
    return state.percentage
  },
  refresh: state => {
    return state.refresh
  },
  errorMsg: state => {
    return state.errorMsg
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}
