<template>
  <v-app id="keep">
    <v-tabs v-model="active" background-color="rgba(0,139,139,1)" slider-color="black">
      <v-tab class="edit_form_style white--text" ripple>
        {{ "Форма заполнения" }}
      </v-tab>
      <v-tab-item>
        <main-tab/>
      </v-tab-item>
      <v-spacer></v-spacer>
      <v-tab class="edit_form_style white--text" ripple>
        {{ "Форма редактирования" }}
      </v-tab>
      <v-tab-item>
        <edit-form/>
      </v-tab-item>
      <v-tab class="edit_form_style white--text" ripple>
        {{ "Загрузка данных из excel" }}
      </v-tab>
      <v-tab-item>
        <upload-data-from-excel/>
      </v-tab-item>
    </v-tabs>
    <message></message>
  </v-app>
</template>

<script>

import Main from '@/components/Main'
import message from '@/components/Message'
import {bus} from "@/main";


export default {
  components: {
    'message':message,
    'main-tab': Main,
    'edit-form': () => import('@/components/FormEdit/EditForm'),
    'upload-data-from-excel': () => import('@/components/UploadDataFromExcel')
  },
  data: () => ({
    active: 0,
  }),
  watch:{
    active:function (){
      bus.$emit('message', {
        show:false
      })
    }
  }
}

</script>
<style>

</style>
