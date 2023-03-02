<template>
  <v-container style="height: max-content;padding-bottom: 15px;">
    <v-row>
      <v-text-field v-model="newName" label="Наименование" style="margin-inline: 15px"/>
      <v-text-field v-model="newCode" label="Код" style="margin-inline: 15px"/>
      <v-btn small class="mt-3" @click="addGroup">Добавить группу</v-btn>
    </v-row>
    <v-data-table
      :headers="headers"
      :items="groups">
      <template v-slot:item.name="{item}">
        <v-edit-dialog
          :return-value.sync="item.name"
          large
          save-text="Сохранить"
          cancel-text="Отменить"
          @save="updateGroup(item)"
        >
          <div>{{ item.name }}</div>
          <template v-slot:input>
            <div class="mt-4 text-h6">
              Изменить имя
            </div>
            <v-text-field
              v-model="item.name"
              label="Редактировать"
              single-line
              autofocus
            ></v-text-field>
          </template>
        </v-edit-dialog>
      </template>
      <template v-slot:item.code="{item}">
        <v-edit-dialog
          :return-value.sync="item.code"
          large
          save-text="Сохранить"
          cancel-text="Отменить"
          @save="updateGroup(item)"
        >
          {{ item.code }}
          <template v-slot:input>
            <div class="mt-4 text-h6">
              Изменить имя
            </div>
            <v-text-field
              v-model="item.code"
              label="Редактировать"
              single-line
              autofocus
            ></v-text-field>
          </template>
        </v-edit-dialog>
      </template>
      <template v-slot:item.delete="{ item }">
        <v-icon
          small
          @click="deleteItem(item)"
        >
          mdi-delete
        </v-icon>
      </template>
    </v-data-table>
    <!--Форма удаления таблицы -->
    <v-dialog v-model="dialogDelete" width="500" persistent>
      <v-card>
        <v-alert color="red" :value="true"> Внимание!
        </v-alert>
        <v-card-text class="text-md-center">
          <v-textarea
            readonly
            value="Вы уверены что хотите удалить группу?"
            rows="1"
            flat
            solo
            auto-grow
            dense>
          </v-textarea>
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn @click="dialogDelete = false; deleteGroup()">
            Да
          </v-btn>
          <v-btn
            @click="
              dialogDelete = false;">
            Нет
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      newName: '',
      newCode: '',
      deletingGroup: null,
      headers: [{text: 'Наименование', value: 'name'}, {text: 'Код', value: 'code'}, {
        text: 'Удаление',
        value: 'delete'
      }],
      groups: [],
      selectedCol: null,
      selectedRow: null,
      dialogDelete: false
    }
  },
  watch: {},
  methods: {
    deleteItem(item) {
      this.deletingGroup = item.id
      this.dialogDelete = true
    },
    deleteGroup() {
      axios
        .post("/api/permission-group/delete", {id: this.deletingGroup})
        .then(() => {
          this.getGroups()
          this.dialogDelete = false
          this.deletingGroup = null
        })
        .catch(e => alert(e.toString()));
    },
    addGroup() {
      axios
        .post("/api/permission-group/add", {name: this.newName, code: this.newCode})
        .then(() => {
          this.getGroups()
          this.newName = ''
          this.newCode = ''
        })
        .catch(e => alert(e.toString()));
    },
    updateGroup(group) {
      axios
        .post("/api/permission-group/edit", group)
        .then(() => this.getGroups())
        .catch(e => alert(e.toString()));
    },
    getGroups() {
      axios
        .post("/api/permission-group/get/all")
        .then(response => (this.groups = response.data))
        .catch(e => alert(e.toString()));
    },
  },
  created() {
    this.getGroups()
  }
}
</script>

