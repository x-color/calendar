import moment from 'moment';

// calendars: []calendar
// calendar: { id, user_id, name, color, []plan }
// plan: { id, name, memo, color, private, all_day, begin, end }
const state = () => ({
  focusDate: moment().format('YYYY-MM-DD'),
  calendars: [
    {
      id: 'calendar id',
      user_id: 'dummmy id',
      active: true,
      name: 'dummy',
      color: 'red',
      plans: [
        {
          id: 'plan id',
          name: 'plan',
          memo: 'sample text',
          color: 'blue',
          private: false,
          start: moment([2020, 5 - 1, 18, 9, 0, 0]),
          end: moment([2020, 5 - 1, 19, 18, 0, 0]),
          allday: false,
        },
      ],
    },
    {
      id: 'calendar id2',
      user_id: 'other id',
      active: true,
      name: 'other',
      color: 'green',
      plans: [
        {
          id: 'plan id2',
          name: 'other\'s plan',
          memo: 'sample text',
          color: 'purple',
          private: false,
          start: moment([2020, 5 - 1, 25, 0, 0, 0]),
          end: moment([2020, 5 - 1, 25, 0, 0, 0]),
          allday: true,
        },
      ],
    },
  ],
});

const getters = {
  getCalendarByID: ({ calendars }) => (id) => calendars.find((calendar) => calendar.id === id),
  // eslint-disable-next-line max-len
  getPlanByID: ({ calendars }) => (id) => calendars.map((calendar) => calendar.plans.find((plan) => id === plan.id)).filter((plan) => plan)[0],
  // eslint-disable-next-line max-len
  getMyCalendars: ({ calendars }, _, rootState) => calendars.filter((calendar) => calendar.user_id === rootState.user.user.id),
  // eslint-disable-next-line max-len
  getSharedCalendars: ({ calendars }, _, rootState) => calendars.filter((calendar) => calendar.user_id !== rootState.user.user.id),
  getActiveCalendars: ({ calendars }) => calendars.filter((calendar) => calendar.active),
  getActivePlans: ({ calendars }) => calendars.map((calendar) => {
    if (calendar.active) {
      return calendar.plans;
    }
    return [];
  }).flat(),
};

const actions = {
  // getCalendars({ commit }) {
  //   Call API
  //   commit('setCalendars', calendars)
  // },
  addCalendar({ commit }, calendar) {
    commit('addCalendar', calendar);
  },
  removeCalendar({ commit }, id) {
    commit('removeCalendar', id);
  },
  editCalendar({ commit }, calendar) {
    commit('setCalendar', calendar);
  },
};

const mutations = {
  setFocusDate(st, date) {
    // eslint-disable-next-line no-param-reassign
    st.focusDate = date;
  },
  setCalendars(st, calendars) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = calendars;
  },
  setCalendar(st, calendar) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      if (calendar.id === cal.id) {
        return calendar;
      }
      return cal;
    });
  },
  addCalendar(st, calendar) {
    st.calendars.push(calendar);
  },
  removeCalendar(st, id) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.fileter((calendar) => calendar.id !== id);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
