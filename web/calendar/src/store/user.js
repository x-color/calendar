const state = () => ({
  user: {
    id: 'dummmy id',
    name: 'Alice',
  },
});

const getters = {
};

const actions = {
};

const mutations = {
  setUser(st, user) {
    // eslint-disable-next-line no-param-reassign
    st.user = user;
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
