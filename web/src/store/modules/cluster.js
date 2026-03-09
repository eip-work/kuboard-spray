const state = {}

// {
//   currentTab: "plan",
//   plan: {
//     currentTab: "",
//   },
//   operation: {
//     currentOpertion: 0,
//     currentStep: 0,
//   }
// }

const mutations = {
  CHANGE_TAB: (state, { cluster, currentTab }) => {
    state[cluster] = state[cluster] || {};
    state[cluster].currentTab = currentTab;
  },
  CHANGE_CLUSTER_STATE: (state, { cluster, tab, key, value }) => {
    state[cluster] = state[cluster] || {};
    state[cluster][tab] = state[cluster][tab] || {};
    state[cluster][tab][key] = value;
  }
}

const actions = {
  changeState({ commit }, data) {
    commit('CHANGE_STATE', data)
  }
}


export default {
  namespaced: true,
  state,
  mutations,
  actions,
}
