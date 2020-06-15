<template>
  <v-dialog v-model="open" persistent max-width="600px" @click:outside="open = false">
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
                v-model="newCalendar.name"
                autofocus
              />
            </v-card-title>
          </v-col>

          <v-col cols="2" class="pl-0">
            <ColorPicker
              v-model="picking"
              :color="newCalendar.color"
              @click="picking = !picking"
              @select="selectColor"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" class="pt-0">
            <v-combobox
              class="px-4 py-0"
              v-model="newCalendar.shares"
              label="Shares (user id)"
              multiple
              chips
              dense
              hide-details
              clearable
              deletable-chips
              append-icon
            ></v-combobox>
          </v-col>
        </v-row>
      </v-container>
      <v-card-actions>
        <v-spacer />
        <v-btn text color="red" @click="open = false">CANCEL</v-btn>
        <v-btn text color="green" :disabled="savable" @click="save">SAVE</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import ColorPicker from '@/components/ColorPicker.vue';

export default {
  name: 'CalendarEditor',
  components: {
    ColorPicker,
  },
  props: {
    value: Boolean,
    calendar: Object,
  },
  computed: {
    open: {
      get() {
        return this.value;
      },
      set(v) {
        this.$emit('input', v);
      },
    },
    savable() {
      if (!this.newCalendar.name || !this.newCalendar.color) {
        return true;
      }
      return false;
    },
  },
  data() {
    return {
      newCalendar: {},
      picking: false,
    };
  },
  methods: {
    selectColor(color) {
      this.newCalendar.color = color;
    },
    save() {
      this.picking = false;
      this.newCalendar.shares.unshift(this.$store.state.user.user.id);
      this.newCalendar.shares = this.newCalendar.shares
        .filter((elm, i, self) => self.indexOf(elm) === i);
      this.$emit('save', this.newCalendar);
      this.open = false;
    },
  },
  watch: {
    value(v) {
      if (v) {
        this.newCalendar = { ...this.calendar };
        this.newCalendar.shares = this.newCalendar.shares
          .filter((id) => id !== this.$store.state.user.user.id);
      }
    },
  },
};
</script>
