<template>
  <div style="padding-top: 0; margin: 0; ">
    <v-tabs vertical >
      <v-tab class="edit_form_style" ripple>
        <v-icon>
          mdi-table-large
        </v-icon>
        {{"Таблицы"}}
      </v-tab>
      <v-tab-item>
        <edit-table/>
      </v-tab-item>
      <v-tab class="edit_form_style" ripple>
        <v-icon>
          mdi-file-tree
        </v-icon>
        {{"Cтруктура"}}
      </v-tab>
      <v-tab-item>
        <edit-tree/>
      </v-tab-item>
      <v-tab class="edit_form_style" ripple>
        <v-img width="26" height="26" contain src="@/assets/countryicon.png"></v-img>
        {{"Регионы"}}
      </v-tab>
      <v-tab-item>
        <region-info/>
      </v-tab-item>
      <v-tab class="edit_form_style" ripple>
        <v-icon>
          mdi-account-group
        </v-icon>
        {{"Группы пользователей"}}
      </v-tab>
      <v-tab-item>
        <user-groups/>
      </v-tab-item>
    </v-tabs>
  </div>
</template>

<script>

export default {
  components: {
    'region-info' : () => import('@/components/FormEdit/RegionInfo'),
    'edit-tree' : () => import('@/components/FormEdit/EditTree'),
    'edit-table':() => import('@/components/FormEdit/EditTable'),
    'user-groups':() => import('@/components/FormEdit/UserGroups'),
  },
  name: 'EditForm'
}
</script>

<style scoped>
.edit_form_style {
  height: 36px;
  text-align: center;
  align-self: baseline;
  color: #000000;
  font-size: 11px;
}
</style>
