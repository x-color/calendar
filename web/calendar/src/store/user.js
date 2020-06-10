const state = () => ({
  user: {
    id: 'user01',
    name: 'Alice',
    signin: false,
  },
  sideMenu: false,
});

const getters = {
};

const actions = {
  signin({ commit }, { username, callback }) {
    commit('setUser', { id: 'user01', name: username, signin: true });
    callback(true);
  },
  signup(st, { callback }) {
    callback(true);
  },
  signout({ commit }) {
    commit('setUser', { id: '', name: '', signin: false });
  },
};

const mutations = {
  setUser(st, user) {
    // eslint-disable-next-line no-param-reassign
    st.user = user;
  },
  setSideMenu(st, open) {
    // eslint-disable-next-line no-param-reassign
    st.sideMenu = open;
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
