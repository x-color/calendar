import moment from 'moment';
import fetchAPI from '@/utils/fetch';

const converter = {
  calAPIModelToModel(calendar) {
    return {
      id: calendar.id,
      user_id: calendar.user_id,
      active: true,
      name: calendar.name,
      color: calendar.color,
      shares: calendar.shares.filter((id) => id !== calendar.user_id),
      plans: calendar.plans.map((p) => this.planAPIModelToModel(p)),
    };
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
    return {
      id: plan.id,
      calendar_id: plan.calendar_id,
      user_id: plan.user_id,
      name: plan.name,
      memo: plan.memo,
      color: plan.color,
      private: plan.private,
      shares: plan.shares,
      start: moment().unix(plan.begin).toISOString(),
      end: moment().unix(plan.end).toISOString(),
      allday: plan.is_all_day,
    };
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
  calendars: [
    {
      id: 'calendar01',
      user_id: 'user01',
      active: true,
      name: 'my calendar',
      color: 'red',
      shares: [],
      plans: [
        {
          id: 'plan01',
          calendar_id: 'calendar01',
          user_id: 'user01',
          name: 'my plan',
          memo: 'sample text',
          color: 'blue',
          private: false,
          shares: [],
          start: moment([2020, 5 - 1, 18, 9, 0, 0]).toISOString(),
          end: moment([2020, 5 - 1, 19, 18, 0, 0]).toISOString(),
          allday: false,
        },
      ],
    },
    {
      id: 'calendar02',
      user_id: 'user02',
      active: true,
      name: 'other\'s calendar',
      color: 'green',
      shares: ['user01'],
      plans: [
        {
          id: 'plan02',
          calendar_id: 'calendar02',
          user_id: 'user02',
          name: 'other\'s plan',
          memo: 'sample text',
          color: 'purple',
          private: false,
          shares: ['calendar01'],
          start: moment([2020, 5 - 1, 25, 0, 0, 0]).toISOString(),
          end: moment([2020, 5 - 1, 25, 0, 0, 0]).toISOString(),
          allday: true,
        },
      ],
    },
  ],
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
  getCalendars({ commit }) {
    const obj = fetchAPI('/calendars');
    const calendars = obj.map((cal) => converter.calAPIModelToModel(cal));
    commit('setCalendars', calendars);
  },
  addCalendar({ commit }, cal) {
    const body = {
      name: cal.name,
      color: cal.color,
    };
    fetchAPI('/calendars', 'POST', JSON.stringify(body))
      .then((calendar) => commit('addCalendar', calendar));
  },
  removeCalendar({ commit }, { id }) {
    fetchAPI(`/calendars/${id}`, 'DELETE');
    commit('removeCalendar', id);
  },
  editCalendar({ commit }, cal) {
    const calendar = { ...cal };
    const body = converter.calModelToAPIModel(cal);
    fetchAPI(`/calendars/${cal.id}`, 'PATCH', JSON.stringify(body));
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
  removePlan({ commit }, { id }) {
    fetchAPI(`/plans/${id}`, 'DELETE');
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
