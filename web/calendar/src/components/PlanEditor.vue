<template>
  <v-dialog v-model="open" persistent max-width="600px" @click:outside="close">
    <v-card>
      <v-container>
        <v-row>
          <v-col cols="10">
            <v-card-title class="py-0">
              <v-text-field
                class="title pt-1 pb-0"
                dense
                hide-details
                placeholder="Name..."
                v-model="newPlan.name"
                autofocus
              />
            </v-card-title>
          </v-col>

          <v-col cols="2">
            <v-btn
              icon
              @click="newPlan.private = !newPlan.private"
            >
              <v-icon v-if="newPlan.private" color="yellow" medium>mdi-lock</v-icon>
              <v-icon v-else medium>mdi-lock-open</v-icon>
            </v-btn>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="5" class="pr-0 ml-1 py-0">
            <DatePicker v-model="newPlan.start" label="Start Date"  :disabled="newPlan.allday" />
          </v-col>
          <v-col cols="5" class="pr-0 py-0">
            <DatePicker v-model="newPlan.end" label="End Date" :disabled="newPlan.allday" />
          </v-col>

          <v-col cols="1" class="py-0 my-auto">
            <v-checkbox
              v-model="newPlan.allday"
              :label="allDayLabel"
              dense
              hide-details
            ></v-checkbox>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="10" class="pt-0">
            <v-select
              v-model="calendar"
              :items="calendars"
              :rules="[v => !!v || 'calendar is required']"
              label="Calendar"
              required
              class="pl-4"
            ></v-select>
          </v-col>

          <v-col cols="2">
            <v-menu
              v-model="colorPicker"
              offset-y
              :close-on-content-click="false"
            >
              <template v-slot:activator="{ on }">
                <v-icon
                  :color="newPlan.color"
                  class="ma-1"
                  medium
                  v-on="on"
                >mdi-circle</v-icon>
              </template>
              <v-card>
                <v-list dense subheader max-width="400">
                  <v-subheader>Color</v-subheader>
                  <v-container>
                    <v-row>
                      <v-col
                        v-for="(color, index) in ['red', 'green', 'blue', 'purple']"
                        :key="index"
                        cols="auto"
                        class="pa-0"
                      >
                        <v-list-item class="px-0">
                          <v-icon
                            large
                            :color="color"
                            @click.stop="selectColor(color)"
                          >mdi-circle</v-icon>
                        </v-list-item>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-list>
              </v-card>
            </v-menu>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" class="pt-0">
            <v-card-text class="pt-0 pb-2">
              <v-textarea
                class="body-2"
                v-model="newPlan.memo"
                placeholder="Add desription..."
                auto-grow
                dense
                hide-details
              />
            </v-card-text>
          </v-col>
        </v-row>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text color="red" @click="$emit('input', false)">CANCEL</v-btn>
          <v-btn
            text
            color="green"
            :disabled="!newPlan.name"
            @click="save"
          >SAVE</v-btn>
        </v-card-actions>
      </v-container>
    </v-card>
  </v-dialog>
</template>

<script>
import moment from 'moment';
import DatePicker from '@/components/DatePicker.vue';

export default {
  name: 'PlanEditor',
  components: {
    DatePicker,
  },
  filters: {
    replaceToHintText(text) {
      if (!text) {
        return 'Add description...';
      }
      return text;
    },
    replaceToHintTitle(text) {
      if (!text) {
        return 'No title...';
      }
      return text;
    },
  },
  props: {
    value: Boolean, // open flag
    start: moment,
    plan: Object,
  },
  computed: {
    calendar: {
      get() {
        return this.cal;
      },
      set(v) {
        this.newPlan.calendar_id = v.id;
        this.newPlan.color = v.color;
      },
    },
    calendars() {
      return this.$store.state.calendars.calendars.map((calendar) => ({
        text: calendar.name,
        value: calendar,
      }));
    },
    open: {
      get() {
        return this.value;
      },
      set(v) {
        this.$emit('input', v);
      },
    },
    allDayLabel() {
      if (!this.$vuetify.breakpoint.xs) {
        return 'allday';
      }
      return '';
    },
  },
  data() {
    return {
      newPlan: {
        id: '',
        calendar_id: '',
        name: '',
        memo: '',
        color: '',
        private: false,
        start: null,
        end: null,
        allday: false,
      },
      editTitleMode: false,
      editTextMode: false,
      colorPicker: false,
      cal: null,
    };
  },
  methods: {
    selectColor(color) {
      this.newPlan.color = color;
      this.colorPicker = false;
    },
    save() {
      if (!this.colorPicker) {
        this.$emit('save');
      }
    },
  },
  watch: {
    value(v) {
      if (v) {
        if (this.plan) {
          this.newPlan = { ...this.plan };
        } else {
          this.newPlan = {
            id: '',
            calendar_id: '',
            name: '',
            memo: '',
            color: 'red',
            private: false,
            start: this.timeList[0],
            end: this.timeList[4],
            allday: false,
          };
        }
      }
    },
  },
};
</script>
