<template>
  <v-container v-if="isLoading === false" style="margin-left: 10px; margin-top: 20px; margin-right: 10px; max-width: 100%">
    <v-snackbar
      v-model="alert"
      :timeout="3000"
      centered
      color="success">
      {{ this.message }}
    </v-snackbar>
    <v-snackbar
      v-model="alertError"
      :timeout="3000"
      centered
      color="red accent-2">
      {{ this.message }}
    </v-snackbar>
    <v-select style="margin-right: 10px"
      v-model="selectedYear"
      :items="years"
      label="Нажмите, чтобы выбрать год"
      dense
      outlined
      @input="getFile(selectedYear)"/>
    <v-layout row wrap>
      <v-flex lg3 md12 sm12 xs12>
        <div>
          <v-file-input ref="file" label="Выбранный файл" dense outlined counter truncate-length="60"
                        v-model="file" style="max-width: fit-content; padding-right: 10px; padding-top: 20px"/>
          <v-layout style="margin-top: 5px">
            <span style="padding-top: 5px; padding-left: 10px">С</span>
            <v-text-field class=".v-text-field input" style="max-width: 16%; margin-left: 10px; margin-top: 0; padding-top: 0;"
              v-model.number="since"
              type="number"
              label="Строка"
            ></v-text-field>
            <span style="padding-top: 5px; padding-left: 10px">До</span>
            <v-text-field class=".v-text-field input" style="max-width: 16%; margin-left: 10px; padding-top: 0; margin-top: 0"
              v-model.number="before"
              type="number"
              label="Строка"
            ></v-text-field>
          </v-layout>
          <v-layout row wrap style="padding-bottom: 20px; padding-top: 10px">
            <v-flex lg5>
              <v-btn style="margin-left: 20px;" v-on:click="submitFile(forms)">Отправить</v-btn>
            </v-flex>
            <v-flex lg12 style="margin-left: 20px">
              <v-textarea class="igor" v-if="responseError !== ''" readonly solo v-bind:value='responseError' style="margin-top: 10px"/>
            </v-flex>
          </v-layout>
        </div>
      </v-flex>
      <v-flex lg9 md12 sm12 xs12 >
        <div class="overflow-y-auto" style="padding-left:20px; max-height: 1000px">
          <div v-for="(form, i) in forms" :key="i" style="padding-right: 20px; padding-bottom: 10px;">
            <v-expansion-panels>
              <v-expansion-panel>
                <v-expansion-panel-header v-if="form.valid === false || form.linkedForms.findIndex(el => el.valid === false) !== -1" class="error">
                  <template>
                    <span style="color: #e7e7e7">{{form.finalValue.text}}</span>
                  </template>
                </v-expansion-panel-header>
                <v-expansion-panel-header v-else class="success">
                  <template>
                    <span style="color: #e7e7e7">{{form.finalValue.text}}</span>
                  </template>
                </v-expansion-panel-header>
                <v-expansion-panel-content style="padding-top: 10px">
                  <v-form v-model="form.valid">
                    <v-autocomplete :disabled="form.complete"
                      v-model="form.headerId"
                      :items="tables"
                      item-text="name"
                      item-value="id"
                      @input="getHeaders($event, i, null, form)"
                      label="Выберите таблицу"
                      prepend-icon="mdi-file-excel"
                      autofocus dense outlined persistent-hint
                      :rules="[v => !!v || 'Нужно выбрать значение поля']">
                      <template slot="no-data">
                        Такой таблицы не существует
                      </template>
                    </v-autocomplete>
                    <v-select :disabled="form.complete" v-model="form.finalValue.value"
                              label="Выберите колонку по которой будет считаться значение" dense outlined prepend-icon="mdi-table-column"
                              :rules="[v => !!v || 'Нужно выбрать значение поля']" :items="form.headers[0]" @input="updateFinalValue(form, $event)"/>
                    <v-select :disabled="form.complete" prepend-icon="mdi-table-column-remove"
                              v-model="form.ignoreValue" label="Выберите колонку в которой есть всего, для её исключения" dense outlined
                              :rules="[v => !!v || 'Нужно выбрать значение поля']" :items="form.headers[0]"/>
                    <v-select :disabled="form.complete"
                              v-model="form.orgNameCol" label="Выберите колонку в которой хранится название организаций" dense outlined
                              prepend-icon="mdi-table-column"
                              :rules="[v => !!v || 'Нужно выбрать значение поля']" :items="form.headers[0]"/>
                    <v-btn class="primary" :disabled="form.complete" @click="addLinkedForm(form)">Добавить связную таблицу</v-btn>
                    <div v-for="(linkedForm, y) in forms[i].linkedForms" :key="y" style="padding-top: 20px">
                      <v-form v-model="linkedForm.valid">
                        <v-autocomplete :disabled="form.complete"
                          v-model="linkedForm.headerId"
                          :items="tables"
                          item-text="name"
                          item-value="id"
                          @input="getHeaders($event, i, y, linkedForm)"
                          label="Выберите таблицу"
                          prepend-icon="mdi-file-excel"
                          autofocus dense outlined persistent-hint
                          :rules="[v => !!v || 'Нужно выбрать значение поля']">
                          <template slot="no-data">
                            Такой таблицы не существует
                          </template>
                        </v-autocomplete>
                        <v-select :disabled="form.complete" v-model="linkedForm.finalValue.value"
                                  label="Выберите колонку по которой будет считаться значение" dense outlined prepend-icon="mdi-table-column"
                                  :rules="[v => !!v || 'Нужно выбрать значение поля']" :items="linkedForm.headers[0]" @input="updateFinalValue(linkedForm, $event)"/>
                        <v-select :disabled="form.complete" prepend-icon="mdi-table-column-remove"
                                  v-model="linkedForm.ignoreValue" label="Выберите колонку в которой есть всего, для её исключения" dense outlined
                                  :rules="[v => !!v || 'Нужно выбрать значение поля']" :items="linkedForm.headers[0]"/>
                        <v-select :disabled="form.complete"
                                  v-model="linkedForm.orgNameCol" label="Выберите колонку в которой хранится название организаций"
                                  dense outlined prepend-icon="mdi-table-column"
                                  :rules="[v => !!v || 'Нужно выбрать значение поля']" :items="linkedForm.headers[0]"/>
                        <v-card-actions>
                          <v-btn :disabled="form.complete" class="primary" @click="remove(linkedForm.id, i, y)">-</v-btn>
                        </v-card-actions>
                      </v-form>
                    </div>
                    <v-spacer></v-spacer>
                    <v-spacer></v-spacer>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn v-if="form.complete" class="primary" @click="editing(form)">Редактировать</v-btn>
                      <v-btn class="primary" @click="remove(form.id, i, null)">Удалить</v-btn>
                      <v-btn :disabled="!form.valid || form.complete || form.linkedForms.findIndex(el => el.valid === false) !== -1"
                             color="primary" @click="save(form)">Сохранить</v-btn>
                    </v-card-actions>
                  </v-form>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>
          </div>
            <div style="display: flex; justify-content: space-between;">
              <v-btn @click="addForm" class="primary">Добавить форму</v-btn>
              <v-btn style="margin-right: 20px" v-if="this.forms.length !== 0" @click="createFile(JSON.stringify(forms), selectedYear)"
                     class="primary">Сохранить шаблон для этого года</v-btn>
            </div>
        </div>
      </v-flex>
    </v-layout>
  </v-container>
  <v-container v-else style="margin-top: 200px; display: flex; align-items: center; justify-content: center">
    <hollow-dots-spinner
      :animation-duration="1000"
      :dot-size="15"
      :dots-num="3"
      color="rgba(0,139,139,1)"
    />
  </v-container>
</template>

<script>
import axios from 'axios'
import { HollowDotsSpinner } from 'epic-spinners'

export default {
  name: 'UploadDataFromExcel',
  components: {
    HollowDotsSpinner,
  },
  data () {
    return {
      file: [],
      backendHeader: [],
      formHeader: [],
      tables: [],
      forms: [],
      since: 3,
      responseError: '',
      before: 0,
      isLoading: false,
      count: 0,
      years: [],
      selectedYear: '',
      message: '',
      alertError: false,
      alert: false
    }
  },
  methods: {
    addForm () {
      if (this.forms.findIndex(el => el.complete === false) === -1) {
        this.forms.push({
          id: this.count,
          headerId: '',
          finalValue: {
           text: '',
           value: ''
          },
          ignoreAddition: true,
          orgNameCol: '',
          ignoreValue: '',
          valid: false,
          complete: false,
          linkedForms: [],
          headers: []
        })
        this.count++
      } else {
        alert('Заполните предыдущую форму')
      }
    },
    addLinkedForm (form) {
      form.linkedForms.push({
        id: this.count,
        formName: this.count,
        headerId: '',
        finalValue: {
          text: '',
          value: ''
        },
        ignoreAddition: false,
        orgNameCol: '',
        valid: false,
        ignoreValue: '',
        headers: []
      })

      this.count++
    },
    editing (form) {
      form.complete = false
    },
    remove (id, indexI, indexY) {
      if (indexY === null) {
        this.forms.splice(this.forms.findIndex(el => el.id === id), 1)
      } else {
        this.forms[indexI].linkedForms.splice(this.forms[indexI].linkedForms.findIndex(el => el.id === id), 1)
      }
    },
    save (form) {
      form.complete = true
    },
    updateFinalValue (form, headerValue) {
      form.finalValue.text = form.headers[0][form.headers[0].findIndex(el => el.value === headerValue)].text
      form.finalValue.value = headerValue
    },
    getHeaders (id_header, i, y, form) {
      form.finalValue = {
        text: '',
        value: ''
      }
      form.orgNameCol = ''
      form.ignoreValue = ''

      this.backendHeader = []
      axios
        .post('/api/getHeader', {id_header})
        .then(response => {
          (this.backendHeader = response.data)
          if (this.backendHeader != null) {
            this.createFormHeader(i, y)
          }
        })
        .catch(e => alert(e.toString()))
    },
    createFormHeader (indexI, indexY) {
      let count = 1
      this.formHeader = []
      for (let i = 0; i < this.backendHeader.length; i++) {
        if (this.backendHeader[i].sub != '') {
          for (let j = 0; j < this.backendHeader[i].sub.length; j++) {
            if (this.backendHeader[i].sub[j].sub != '') {
              for (let q = 0; q < this.backendHeader[i].sub[j].sub.length; q++) {
                if (this.backendHeader[i].sub[j].sub[q].sub != '') {
                  for (let v = 0; v < this.backendHeader[i].sub[j].sub[q].sub.length; v++) {
                    this.formHeader.push({
                      text: this.backendHeader[i].text + ' | ' + this.backendHeader[i].sub[j].text +
                        ' | ' + this.backendHeader[i].sub[j].sub[q].text + ' | ' +
                        this.backendHeader[i].sub[j].sub[q].sub[v].text,
                      value: count.toString()
                    })
                    count++
                  }
                } else {
                  this.formHeader.push({
                    text: this.backendHeader[i].text + ' | ' + this.backendHeader[i].sub[j].text +
                      ' | ' + this.backendHeader[i].sub[j].sub[q].text,
                    value: count.toString()
                  })
                  count++
                }
              }
            } else {
              this.formHeader.push({ text: this.backendHeader[i].text + ' | ' + this.backendHeader[i].sub[j].text, value: count.toString() })
              count++
            }
          }
        } else {
          this.formHeader.push({
            text: this.backendHeader[i].text,
            value: count.toString()
          })
          count++
        }
      }
      if (indexY === null) {
        this.forms[indexI].headers.splice(0, 1)
        this.forms[indexI].headers.splice(0, 0, this.formHeader)
      } else {
        this.forms[indexI].linkedForms[indexY].headers.splice(0, 1)
        this.forms[indexI].linkedForms[indexY].headers.splice(0, 0, this.formHeader)
      }
    },
    submitFile (unFilteredForms) {
      if (this.since === 0 || this.before === 0) {
        this.alertError = true
        this.message = "Выберите с какой и до какой строки будет выполняться программа"
      } else {
        this.responseError = ''
        if (this.forms.length !== 0) {
          if (this.forms.findIndex(el => el.complete === false) === -1) {
            if (this.file.length !== 0) {
              const formData = new FormData()
              formData.append('file', this.file)
              this.isLoading = true
              axios
                .post('/api/uploadFile', formData, {headers: {'Content-Type': 'multipart/form-data'}})
                .then(() => {
                  let year = this.selectedYear
                  let since = this.since
                  let before = this.before
                  let forms = []
                  unFilteredForms.forEach(form => {
                    forms.push({
                      id: form.id,
                      headerId: form.headerId,
                      finalValue: form.finalValue,
                      ignoreAddition: form.ignoreAddition,
                      orgNameCol: form.orgNameCol,
                      ignoreValue: form.ignoreValue,
                      linkedForms: form.linkedForms,
                    })
                  })
                  axios
                    .post('/api/uploadFile', {forms, since, before, year})
                    .then(() => {
                      this.isLoading = false
                      this.alert = true
                      this.message = "Данные успешно загружены"
                    })
                    .catch(error => {
                      if (error.response) {
                        this.isLoading = false
                        this.responseError += "Ошибки:\n" + error.response.data
                      }
                    })
                })
                .catch(() => {
                  this.isLoading = false
                  this.alertError = true
                  this.message = "Не удалось загрузить файл, возможно он был изменён"
                })
            } else {
              this.alertError = true
              this.message = "Выберите файл"
            }
          } else {
            this.alertError = true
            this.message = "Заполните форму"
          }
        } else {
          this.alertError = true
          this.message = "Нужно добавить формы"
        }
      }
    },
    getTables(year) {
      axios
        .post('/api/table', { year })
        .then(response => {
          this.tables = response.data
        })
        .catch(e => alert(e.toString()));
    },
    createFile(data, year) {
      if (this.forms.length !== 0) {
        if (this.forms.findIndex(el => el.complete === false) === -1) {
          axios
            .post('/api/createFile', {data, year})
            .then(() => {
              this.message = "Шаблон успешно создан"
              this.alert = true
            })
            .catch(e => alert(e.toString()));
        } else {
          this.alertError = true
          this.message = "Для сохранения шаблона, сохраните его форму!"
        }
      } else {
        this.alertError = true
        this.message = "Для сохранения шаблона, сначала создайте форму!"
      }
    },
    getFile(year) {
      this.getTables(this.selectedYear)
      axios
        .get('/api/getTemplate', { params: { year: year }})
        .then(response => {
          this.forms = []
          this.forms = JSON.parse(response.data)
          this.count = this.forms.length + 1
        })
        .catch(e => {
          if (e.response.data === 4452) {
            this.forms = []
            this.count = this.forms.length
            this.alertError = true
            this.message = "Шаблон для этого года отсутсвует, создайте его!"
          } else {
            alert(e.toString())
          }
        })
    },
    getYears(){
      axios
        .post("/api/years")
        .then(response => {
          this.years = response.data
          this.selectedYear = this.years[0]
          this.getTables(this.selectedYear)
          this.getFile(this.selectedYear)
        })
        .catch(e => alert(e.toString()));
    }
  },
  created () {
    this.getYears()
  }
}
</script>

<style scoped>
.igor >>> textarea {
  color: rgba(180, 22, 27, 0.7);
}
.v-text-field input {
  padding: 6px 0 0;
}
</style>
