import fetchAPI from '@/utils/fetch';

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
  signin({ commit }, { username, password, callback }) {
    const body = {
      name: username,
      password,
    };
    fetchAPI('/auth/signin', 'POST', JSON.stringify(body), false)
      .then((obj) => {
        commit('setUser', { id: obj.id, name: username, signin: true });
        callback(true);
      })
      .catch(() => callback(false));
  },
  signup(st, { username, password, callback }) {
    const body = {
      name: username,
      password,
    };
    fetchAPI('/auth/signup', 'POST', JSON.stringify(body), false)
      .then(() => callback(true))
      .catch(() => callback(false));
  },
  signout({ commit }) {
    fetchAPI('/auth/signout', 'POST');
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
