import moment from 'moment';
import fetchAPI from '@/utils/fetch';

const converter = {
  calAPIModelToModel(cal) {
    const calendar = {
      id: cal.id,
      user_id: cal.user_id,
      name: cal.name,
      active: true,
      color: cal.color,
      shares: cal.shares,
    };
    if (cal.plans) {
      calendar.plans = cal.plans.map((p) => this.planAPIModelToModel(p));
    } else {
      calendar.plans = [];
    }
    return calendar;
  },
  calModelToAPIModel(calendar) {
    return {
      id: calendar.id,
      user_id: calendar.user_id,
      name: calendar.name,
      color: calendar.color,
      shares: calendar.shares,
      plans: calendar.plans.map((p) => this.planModelToAPIModel(p)),
    };
  },
  planAPIModelToModel(plan) {
    const p = {
      id: plan.id,
      calendar_id: plan.calendar_id,
      user_id: plan.user_id,
      name: plan.name,
      memo: plan.memo,
      color: plan.color,
      private: plan.private,
      shares: plan.shares,
      start: moment.unix(plan.begin).toISOString(),
      end: moment.unix(plan.end).toISOString(),
      allday: plan.is_all_day,
    };
    if (plan.private && plan.name === '') {
      p.name = 'Private Plan';
    }
    return p;
  },
  planModelToAPIModel(plan) {
    return {
      id: plan.id,
      calendar_id: plan.calendar_id,
      user_id: plan.user_id,
      name: plan.name,
      memo: plan.memo,
      color: plan.color,
      private: plan.private,
      shares: plan.shares,
      begin: moment(plan.start).unix(),
      end: moment(plan.end).unix(),
      is_all_day: plan.allday,
    };
  },
};

// calendars: []calendar
// calendar: { id, user_id, name, color, []plan }
// plan: { id, name, memo, color, private, all_day, begin, end }
const state = () => ({
  focusDate: moment().format('YYYY-MM-DD'),
  calendars: [],
});

const getters = {
  getCalendarByID: ({ calendars }) => (id) => calendars.find((calendar) => calendar.id === id),
  // eslint-disable-next-line max-len
  getPlanByID: ({ calendars }) => (id) => calendars.map((calendar) => {
    const p = calendar.plans.find((plan) => id === plan.id);
    if (!p) {
      return null;
    }
    const plan = { ...p };
    plan.start = moment(p.start);
    plan.end = moment(p.end);
    return plan;
  }).filter((plan) => plan)[0],
  // eslint-disable-next-line max-len
  getMyCalendars: ({ calendars }, _, rootState) => calendars.filter((calendar) => calendar.user_id === rootState.user.user.id),
  // eslint-disable-next-line max-len
  getSharedCalendars: ({ calendars }, _, rootState) => calendars.filter((calendar) => calendar.user_id !== rootState.user.user.id),
  getActiveCalendars: ({ calendars }) => calendars.filter((calendar) => calendar.active),
  getActivePlans: ({ calendars }) => () => calendars.map((calendar) => {
    if (!calendar.active) {
      return [];
    }
    return calendar.plans.map((p) => {
      const plan = { ...p };
      plan.start = moment(p.start);
      plan.end = moment(p.end);
      return plan;
    });
  }).flat(),
};

const actions = {
  getCalendars({ state: st, commit }) {
    fetchAPI('/calendars')
      .then((obj) => {
        const calendars = obj.map((cal) => {
          const calendar = converter.calAPIModelToModel(cal);
          const old = st.calendars.filter((c) => c.id === cal.id);
          if (old.length) {
            calendar.active = old[0].active;
          }
          return calendar;
        });
        commit('setCalendars', calendars);
      })
      .catch(() => fetchAPI('/register', 'POST')
        .then(() => {
          const body = {
            name: 'calendar',
            color: 'red',
          };
          fetchAPI('/calendars', 'POST', JSON.stringify(body))
            .then((resCal) => commit('addCalendar', converter.calAPIModelToModel(resCal)))
            .then(() => fetchAPI('/calendars'))
            .then((obj) => {
              const calendars = obj.map((cal) => converter.calAPIModelToModel(cal));
              commit('setCalendars', calendars);
            });
        }));
  },
  addCalendar({ commit }, cal) {
    const body = {
      name: cal.name,
      color: cal.color,
    };
    fetchAPI('/calendars', 'POST', JSON.stringify(body))
      .then((resCal) => {
        commit('addCalendar', converter.calAPIModelToModel(resCal));
      });
  },
  removeCalendar({ commit }, { id }) {
    fetchAPI(`/calendars/${id}`, 'DELETE');
    commit('removeCalendar', id);
  },
  editCalendar({ commit }, { cal, noapi = false }) {
    const calendar = { ...cal };
    if (!noapi) {
      const body = converter.calModelToAPIModel(cal);
      fetchAPI(`/calendars/${cal.id}`, 'PATCH', JSON.stringify(body));
    }
    commit('setCalendar', calendar);
  },
  addPlan({ commit }, p) {
    const plan = { ...p };
    plan.start = p.start.toISOString();
    plan.end = p.end.toISOString();
    const body = converter.planModelToAPIModel(plan);
    fetchAPI('/plans', 'POST', JSON.stringify(body))
      .then((resPlan) => commit('addPlan', converter.planAPIModelToModel(resPlan)));
  },
  // eslint-disable-next-line camelcase
  removePlan({ commit }, { id, calendar_id }) {
    const body = {
      calendar_id,
    };
    fetchAPI(`/plans/${id}`, 'DELETE', JSON.stringify(body));
    commit('removePlan', id);
  },
  editPlan({ commit }, p) {
    const plan = { ...p };
    plan.start = p.start.toISOString();
    plan.end = p.end.toISOString();
    const body = converter.planModelToAPIModel(plan);
    fetchAPI(`/plans/${p.id}`, 'PATCH', JSON.stringify(body));
    commit('setPlan', plan);
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
    st.calendars = st.calendars.filter((calendar) => calendar.id !== id);
  },
  setPlan(st, plan) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      if (cal.id === plan.calendar_id) {
        // eslint-disable-next-line no-param-reassign
        cal.plans = cal.plans.map((p) => {
          if (p.id === plan.id) {
            return plan;
          }
          return p;
        });
      }
      return cal;
    });
  },
  addPlan(st, plan) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      if (cal.id === plan.calendar_id) {
        cal.plans.push(plan);
      }
      return cal;
    });
  },
  removePlan(st, id) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      // eslint-disable-next-line no-param-reassign
      cal.plans = cal.plans.filter((p) => p.id !== id);
      return cal;
    });
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
