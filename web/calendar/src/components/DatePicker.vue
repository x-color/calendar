<template>
  <v-container class="pb-0 pr-0">
    <v-row>
      <v-col cols="8" sm="7" class="pr-0 pb-0">
        <v-menu
          v-model="open"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
        >
          <template v-slot:activator="{ on }">
            <v-text-field
              v-model="dateFormatted"
              :rules="[(v) => !!v || 'Required']"
              :label="label"
              @blur="date = parseDate(dateFormatted)"
              v-on="on"
              dense
              hide-details
              required
              :class="{ 'sm-text': $vuetify.breakpoint.xs }"
            ></v-text-field>
          </template>
          <v-date-picker
            :value="date"
            no-title
            dense
            @input="open = false; date = parseDate($event)"
          ></v-date-picker>
        </v-menu>
      </v-col>

      <v-col cols="4" sm="5" class="pl-0">
        <!--
        NOTE: v-combobox has bug...?
              A change of value(set to v-model) does not reflect view if input value is invallid.
              e.g. '', ' ', 'a' etc...
        -->
        <v-combobox
          v-model="time"
          :rules="[(v) => !!v || 'Required']"
          :items="timeList"
          @blur="date = parseTime(time)"
          validate-on-blur
          dense
          hide-details
          append-icon
          :class="{ 'sm-text': $vuetify.breakpoint.xs }"
          :disabled="disabled"
        ></v-combobox>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import moment from 'moment';

export default {
  name: 'DatePicker',
  props: {
    value: moment,
    label: String,
    disabled: Boolean,
  },
  data(vm) {
    return {
      open: false,
      time: vm.value.format('HH:mm'),
      dateFormatted: vm.value.format('YYYY/MM/DD'),
    };
  },
  computed: {
    date: {
      get() {
        return this.value.format('YYYY-MM-DD');
      },
      set(v) {
        this.$emit('input', moment(v));
      },
    },
    timeList() {
      const offset = moment().hour(0).minute(0).subtract(15, 'minutes');
      return [...Array(96)].map(() => offset.add(15, 'minutes').clone()).map((m) => m.format('HH:mm'));
    },
  },
  methods: {
    parseDate(s) {
      const date = moment(s, 'YYYY/MM/DD');
      if (!date.isValid()) {
        return this.value.format();
      }
      return date.hour(this.value.hour()).minute(this.value.minute()).format();
    },
    parseTime(s) {
      const date = moment(s, 'HH:mm');
      if (!date.isValid()) {
        return this.value.format();
      }
      // eslint-disable-next-line max-len
      return date.year(this.value.year()).month(this.value.month()).date(this.value.date()).format();
    },
  },
  watch: {
    value(v) {
      this.dateFormatted = v.format('YYYY/MM/DD');
      this.time = v.format('HH:mm');
    },
  },
};
</script>

<style scoped>
.sm-text {
  font-size: 12px;
}
</style>
