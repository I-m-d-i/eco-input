<template>
  <v-snackbar
    max-height="600"
    :centered="messagesDisplay.centered"
    multi-line
    :timeout="-1"
    v-model="messagesDisplay.show"
    :color="messagesDisplay.color"
    transition="scale-transition"
  >
<span v-bind:style="{color:messagesDisplay.textColor }"
      style="max-height:600px; white-space: pre-line; font-size: 0.85rem;">
         {{ messagesDisplay.message }}
       </span>
    <template v-slot:action="{ attrs }">
      <v-icon
        :color="messagesDisplay.textColor"
        v-bind="attrs"
        @click="messagesDisplay.show = false;timeoutID=null">
        mdi-close-circle-outline
      </v-icon>
    </template>
  </v-snackbar>
</template>
<script>
import {bus} from "@/main";

export default {
  name: "message",
  data: () => ({
    active: 0,
    messagesDisplay: {
      message: "",
      color: null,
      show: false,
      textColor: null,
      centered: false,
      timeout: 0,
    },
    timeoutID: null,
  }),
  mounted() {
    bus.$on('message', (message) => {
      this.messagesDisplay.show = (message.show === undefined || message.show)
      if (!this.messagesDisplay.show) {
        if (this.timeoutID != null) {
          clearTimeout(this.timeoutID)
          this.close()
        }
        return;
      }
      this.messagesDisplay.message = (message.message ?? "")
      this.messagesDisplay.color = message.color ?? "#4caf50"
      this.messagesDisplay.textColor = message.textColor ?? "FFFFFFFF"
      this.messagesDisplay.centered = message.centered ?? false
      this.messagesDisplay.timeout = message.timeout ?? 4000
      if (this.timeoutID) {
        this.messagesDisplay.show = false
        setTimeout(() => {
          this.messagesDisplay.show = true
        }, 100)
        clearTimeout(this.timeoutID)
      }
      if (this.messagesDisplay.timeout === -1) return
      this.timeoutID = setTimeout(() => {
        this.close()
      }, this.messagesDisplay.timeout)
    })
  },
  methods: {
    close() {
      this.messagesDisplay.show = false
      setTimeout(() => {
        this.messagesDisplay.message = ""
        this.messagesDisplay.color = null
        this.messagesDisplay.textColor = null
        this.messagesDisplay.timeout = null
        this.messagesDisplay.centered = false
        this.timeoutID = null
      }, 200)
    }
  },
}
</script>

