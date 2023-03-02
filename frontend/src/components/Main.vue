<template id="app">
  <v-container>
    <v-layout text-xs-center wrap>
      <v-flex xs6 pa-3>
        <v-autocomplete
          :filter="customFilter"
          v-model="selectedRegionName"
          :items="regions"
          item-text="name"
          label="МО или ОКТМО"
        />
      </v-flex>
      <v-flex xs6 pa-3>
        <v-autocomplete
          v-model="selectedTable"
          :items="tables"
          item-value="id"
          item-text="name"
          label="Таблица"
        />
      </v-flex>
      <v-row justify-start>
        <v-flex mx-3 sm1>
          <v-select solo v-model="selectedYear" :items="years" label="Год">
          </v-select>
        </v-flex>
        <v-btn small elevation="2" active-class='myBtn' :depressed="this.showTable" :inputValue=this.showForm
               v-on:click="showOrHideForm()" style="margin-left: 20px">
          Форма для заполнения
        </v-btn>
        <v-btn small elevation="2" active-class='myBtn' :depressed="this.showTable" :inputValue=this.showCalculatingForm
               v-on:click="showOrHideCalculatingForm()" style="margin-left: 20px">
          Посчитать "Итого"
        </v-btn>
        <v-btn small elevation="2" active-class='myBtn' :depressed="this.showTable" :inputValue=this.showTable
               v-on:click="showOrHideTable()" style="margin-left: 20px">
          Таблица со значениями
        </v-btn>
      </v-row>
      <template v-if="selectedTable && selectedRegionName !== ''">
        <!-- Начало создания формы для заполнения -->
        <template v-if="showForm">
          <v-flex xs12>
            <v-text-field v-for="(header, Index) in selectHeader" :key="header.value" :label='header.text'
                          v-model="valueForSave[Index]"/>
            <v-btn small color="primary" :disabled="isFormEmpty" @click="sendData()">
              Сохранить информацию
            </v-btn>
          </v-flex>
        </template>
        <!-- Начало создания формы для Итогового значения -->
        <template v-if="showCalculatingForm">
          <v-form style="margin: 20px 10px 20px 10px; min-width: 100%" v-model="valid">
            <v-layout row>
              <v-flex xs4>
                <v-text-field style="max-width: available; margin-right: 10px" v-model="desiredTextValue" autofocus
                              dense outlined label="ИТОГО от:"
                              hint="Введите значение по которому будет считаться 'Итого'"
                              :rules="[v => !!v || 'Введите значение по которому будет считаться \'Итого\'']"/>
                <v-text-field style="max-width: available; margin-right: 10px" v-model="inputText"
                              dense outlined label="Что будет записано в ячейку"
                              hint="Введите, что должно быть записанно в ячейку на замену Итого"
                              :rules="[v => !!v || 'Введите, что должно быть записанно в ячейку на замену \'Итого\'']"/>
              </v-flex>
              <v-flex xs8>
                <v-select prepend-icon="mdi-table-column" v-model="searchIn"
                          label="Выберите колонку в которой есть введённое значение" dense outlined
                          :rules="[v => !!v || 'Нужно выбрать значение поля']" :items="selectHeader" item-text="text"
                          item-value="value"/>
                <v-select prepend-icon="mdi-table-column" v-model="ignoreValue"
                          label="Выберите колонку в которую не учитывать при подсчёте(необязательно)" dense outlined
                          :items="selectHeader" item-text="text" item-value="value"/>
              </v-flex>
              <v-btn :disabled="!valid" small elevation="2" class="primary" v-on:click="calculateFinalValue()">
                Посчитать
              </v-btn>
            </v-layout>
          </v-form>
        </template>
        <!-- Начало создания таблицы -->
        <template v-if="showTable">
          <v-flex xs12>
            <v-text-field
              style="width: 40vh"
              v-model="search"
              append-icon="mdi-magnify"
              label="Поиск"
              single-line
              hide-details/>
            <v-data-table
              :loading="loadingData"
              :headers="tableHeader"
              @update:page="page=$event"
              @update:items-per-page="itemsPerPage=$event"
              :search="search"
              :items=rowsForTable
              dense
            >
              <template v-slot:item.action="{item,index}">
                <v-lazy>
                  <td class="justify-center layout px-0" v-if="elementForSort[0] == null">
                    <v-icon small class="mdi-file-multiple" color="#8A00FF" title="Копирование записи"
                            @click="(RowForChange = Object.assign({},searchRowByIndex(index))); (editOrCopy='Дублирование записи');
                      (windowsForEditOrCopy = true)"> mdi-file-multiple
                    </v-icon>
                    <v-icon small class="mdi-transfer" color="#008F00" title="Отправить в другой регион"
                            @click="windowForTransferRegion = true; valueForTransferRegion.idRowForTransfer = searchRowByIndex(index).id;">
                      mdi-transfer
                    </v-icon>
                    <v-icon small class="mdi-pencil" color="#00AAAA" title="Изменение записи"
                            @click="(RowForChange = Object.assign({},searchRowByIndex(index))); (editOrCopy='Редактирование записи');
                      (windowsForEditOrCopy = true)"> mdi-pencil
                    </v-icon>
                    <v-icon small class="mdi-compare-vertical" color="#A0A000" title="Перенести запись"
                            @click="changeSort(index);">mdi-compare-vertical
                    </v-icon>
                    <v-icon small class="mdi-delete" color="#FF0000" title="Удаление записи"
                            @click="(windowForDelete = true); deleteRow(searchRowByIndex(index))"> mdi-delete
                    </v-icon>
                  </td>
                  <td class="justify-center layout px-0" v-if="elementForSort[0] != null">
                    <v-icon small class="mdi-backspace" color="#990000" title="Отменить перенос"
                            @click="elementForSort = []"> mdi-cancel
                    </v-icon>
                    <v-icon small class="mdi-compare_arrows" color="#ff9900"
                            :disabled="elementForSort[0].index===(page-1)*itemsPerPage+index"
                            title="Поменять с этой записью"
                            @click="changeSort(index)">mdi-compare-vertical
                    </v-icon>
                  </td>
                </v-lazy>
              </template>
              <template v-slot:no-data>
                Нет данных
              </template>
            </v-data-table>
          </v-flex>
          <v-btn small color="primary" v-on:click="getRows(); getHeader()">
            Обновить информацию
          </v-btn>
        </template>
      </template>
    </v-layout>
    <!--Форма для редактирования-->
    <v-dialog v-if="windowsForEditOrCopy" v-model="windowsForEditOrCopy" persistent>
      <v-card>
        <v-card-title>
          <span class="headline">{{ editOrCopy }}</span>
        </v-card-title>
        <v-card-text>
          <v-container grid-list-md>
            <v-layout wrap align-end>
              <v-flex v-for="(header, indexPair) in selectHeader" :key="header.value" xs12 sm6 md4>
                {{ header.text }}
                <v-textarea outlined no-resize v-model="RowForChange.value[indexPair]"/>
              </v-flex>
            </v-layout>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn color="blue darken-1" text @click="windowsForEditOrCopy = false;RowForChange={}; getRows()">
            Отмена
          </v-btn>
          <v-btn color="blue darken-1" text @click="windowsForEditOrCopy = false; copyOrEditRow()">
            Сохранить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!--Окно выбора регионов для переноса строки-->
    <v-dialog v-if="windowForTransferRegion" v-model="windowForTransferRegion" scrollable max-width="300px">
      <v-card>
        <v-card-title>Выбрать регион для переноса</v-card-title>
        <v-divider/>
        <v-card-text style="height: 100px;">
          <v-autocomplete
            v-model="valueForTransferRegion.selectedTransferRegion"
            :items="regions"
            item-text="name"
            item-value="id"/>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="red darken-1"
            text
            @click="windowForTransferRegion = false; valueForTransferRegion.selectedTransferRegion = null; valueForTransferRegion.idRowForTransfer = null;"
          >Отмена
          </v-btn>
          <v-btn
            color="green darken-1"
            :disabled="!valueForTransferRegion.selectedTransferRegion"
            text
            @click="windowForTransferRegion = false; transferRegion()"
          >Сохранить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!--Оповещение об удалении -->
    <v-dialog v-if="windowForDelete" v-model="windowForDelete" width="400" persistent>
      <v-card>
        <v-alert type="warning" :value="true"> Внимание
        </v-alert>
        <v-card-text class="text-md-center">
          Вы действительно хотите удалить данную строку?
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn @click="permissionForAction = true; windowForDelete = false">
            Да
          </v-btn>
          <v-btn @click="windowForDelete = false">
            Нет
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!--Оповещение об ошибке -->
    <v-dialog v-if="windowsForError" v-model="windowsForError" width="400" persistent>
      <v-card>
        <v-alert type="error" :value="true"> Ошибка
        </v-alert>
        <v-card-text class="text-md-center"> {{ messageError }}
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn @click="windowsForError = false"> Закрыть
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
<script>

import axios from "axios";
import {bus} from "@/main";

export default {
  name: "app",
  data: () => ({
    valid: false,
    desiredTextValue: '',                     // Значение которое нужно найти
    inputText: '',                            // Что будет записываться в ячейку после подсчётов
    search: '',
    searchIn: '',                             // где искать desiredTextValue
    ignoreValue: '',                          // исключаемый столбец
    selectHeader: [],                         // Вариант хедера для селекта "ИТОГО"

    selectedTable: null,
    selectedRegionName: '',                   // Выбранный регион
    selectedYear: null, // Выбранный год
    editOrCopy: "",                           // Выбранная функция, редактировать или копировать
    page: 1,
    itemsPerPage: 10,

    showForm: false,                          // Отображение формы, смена цвета кнопки
    showTable: true,                          // Отображение таблицы, смена цвета кнопки
    showCalculatingForm: false,               // Отображение формы для подсчёта "Итого", смена цвета кнопки

    windowForDelete: false,                   // Окно с предупреждением об удаление
    windowForTransferRegion: false,           // Окно для переноса информации в другой регион
    windowsForError: false,                   // Окно с информацией о возникновении ошибки
    windowsForEditOrCopy: false,              // Окно для редактирования/копирования информации
    permissionForAction: false,                // Разрешение на действие (удаление, редактирование, копирование)

    tables: [],                               // Список всех таблиц {id, name}
    regions: [],                              // Список всех регионов {id, name}
    years: [],                                // Список всех годов в базе

    tableHeader: [],                          // Вариант хедера для таблицы
    rows: [],                                 // Значения таблицы {id, sort, values[]}
    valueForSave: [],                         // Массив для хранения значений из формы (для сохраненя инфы)
    elementForSort: [],                       // Два элемента, менюящиеся местами
    RowForChange: {},                          // Выбранная строка для редактирования/копирования
    valueForTransferRegion: {                  // Для переноса строки в другой регион
      selectedTransferRegion: null,
      idRowForTransfer: null,
    },
    messageError: "", // текст ошибки?
    loadingData: false,
  }),

  computed: {
    // Выборка данных для отправки на бэкенд для получения значений таблицы
    sendId() {
      if (this.selectedTable === null || this.selectedRegionName === "") {
        return null;
      }
      return {
        id_header: this.selectedTable,
        id_region: this.regions.filter(item => item.name === this.selectedRegionName)[0]?.id
      };
    },
    isFormEmpty() {
      for (const argument of this.valueForSave) {
        if (argument != null && argument !== "") {
          return false
        }
      }
      return true
    },
    rowsForTable() {
      if (this.rows != null) {
        let rowForTable = []
        for (const row of this.rows) {
          rowForTable.push(Object.freeze(Object.assign({}, row.value)))
        }
        return rowForTable
      }
    },
  },
  watch: {
    selectHeader: function () {
      this.valueForSave = new Array(this.selectHeader.length)
    },
    sendId: function () {                      // Для постоянной проверки выбранного региона и таблицы
      if (this.sendId != null) {
        this.getRows();
        this.getHeader();
      }
    },
    selectedYear: function () {                      // Для постоянной проверки выбранного года
      if (this.tables.length > 0) {
        this.selectedTable = null
        this.rows = []
        this.tableHeader = []
        this.selectHeader = []
        this.getTables(this.selectedYear)
      }
    }
  },

  methods: {
    customFilter(item, queryText) {
      const textOne = item.name.toLowerCase().replaceAll(" ", "")
      const textTwo = item.oktmo.toLowerCase().replaceAll(" ", "")
      const searchText = queryText.toLowerCase().replaceAll(" ", "")
      const searchText2 = queryText.toLowerCase().replaceAll(" ", "").slice(0, 5)
      return textOne.indexOf(searchText) > -1 || textTwo.indexOf(searchText2) > -1
    },

    sendForTransfer() {
      if (this.valueForTransferRegion.idRowForTransfer) {
        return {
          id_header: this.selectedTable,
          id_region: this.regions.filter(item => item.name === this.selectedRegionName)[0].id,
          id_region_target: this.valueForTransferRegion.selectedTransferRegion,
          id_row: this.valueForTransferRegion.idRowForTransfer,
        };
      }
    },
    calculateFinalValue() {
      let toBack = {
        regionId: this.regions.filter(item => item.name === this.selectedRegionName)[0].id,
        headerId: this.selectedTable,
        desiredTextValue: this.desiredTextValue,
        inputText: this.inputText,
        searchIn: this.searchIn,
        ignoreValue: this.ignoreValue
      }
      this.sendForm(toBack)
    },
    showOrHideForm() {
      this.showForm = !this.showForm;
    },
    showOrHideTable() {
      this.showTable = !this.showTable;
    },
    showOrHideCalculatingForm() {
      this.showCalculatingForm = !this.showCalculatingForm;
    },
    createFormHeader(backendHeader, text) {
      if (backendHeader.sub.length !== 0) {
        for (let i = 0; i < backendHeader.sub.length; i++) {
          this.createFormHeader(backendHeader.sub[i], backendHeader?.text ? text + backendHeader.text + "|" : "")
        }
      } else {
        this.tableHeader.push({
          text: text + backendHeader.text,
          value: (this.tableHeader.length).toString(),
          align: 'center',
          divider: true
        })
        this.selectHeader.push({
          text: text + backendHeader.text,
          value: (this.selectHeader.length).toString()
        })
      }
    },
    searchRowByIndex: function (index) {
      return this.rows[(this.page - 1) * this.itemsPerPage + index]
    },
    getRowForSave(value) {
      return {
        id_header: this.sendId.id_header,
        id_region: this.sendId.id_region,
        value: this.getValueWithoutHoles(value, 0)
      };
    },
    getValueWithoutHoles: function (arrayWithHoles, start) {
      let valueArray = [];
      for (let i = start; i < this.tableHeader.length - (1 - start); i++) {
        if (arrayWithHoles[i] != null) {
          valueArray.push(arrayWithHoles[i])
        } else {
          valueArray.push("")
        }
      }
      return valueArray
    },
    message(message, options) {
      bus.$emit('message', {
        message: message,
        color: options.color,
        centered: options.centered
      })
    },
    //Запросы к бэку
    deleteRow(row) {
      if (this.permissionForAction) {
        axios
          .post("/api/deleteRow", row.id)
          .then(response => {
            if (response.status !== 200) {
              throw response.data;
            }
            this.permissionForAction = false;
            this.message("Информация удалена", {color: "green", centered: true})
          })
          .catch(e => {
            this.message("Ошибка удаления" + e.toString(), {color: "red", centered: true})
          });
        setTimeout(this.getRows, 300);
        return;
      }
      setTimeout(this.deleteRow, 0, row);
    },
    copyOrEditRow() {
      if (this.editOrCopy === "Редактирование записи") {
        axios
          .post("/api/editRow", this.RowForChange)
          .then(response => {
            if (response.status !== 200) {
              throw response.data;
            }
            this.message("Запись изменена", {color: "green", centered: true})
          })
          .catch(e => {
            this.message("Ошибка редактирования \n" + e.response.data, {color: "red", centered: true})
          });
      } else {
        axios
          .post("/api/copyRow", this.RowForChange)
          .then(response => {
            if (response.status !== 200) {
              throw response.data;
            }
            this.message("Запись продублирована", {color: "green", centered: true})
          })
          .catch(e => {
            this.message("Ошибка при копировании строки \n" + e.response.data, {color: "red", centered: true})
          });
      }
      this.RowForChange = {}
      setTimeout(this.getRows, 300);
    },
    changeSort(index) {
      let elementTableForSort = this.searchRowByIndex(index)
      this.elementForSort.push({
        id: elementTableForSort.id,
        sort: elementTableForSort.sort,
        index: (this.page - 1) * this.itemsPerPage + index
      })
      if (this.elementForSort.length === 2) {
        axios
          .post("/api/changeSortRows", this.elementForSort)
          .then(response => {
            if (response.status !== 200) {
              throw response.data;
            }
            this.message("Строки поменялись местами", {color: "green", centered: true})
            this.elementForSort = []
          })
          .catch(e => {
            this.message("Ошибка изменения сортировки \n" + e.toString(), {color: "red", centered: true})
            this.elementForSort = []
          });
        setTimeout(this.getRows, 200);
      }
    },
    transferRegion() {
      axios
        .post("/api/transferRowRegion", this.sendForTransfer())
        .then(response => {
          if (response.status !== 200) {
            throw response.data;
          }
          this.message("Строка отправлена в другой регион", {color: "green", centered: true})
        })
        .catch(e => {
          if (e.response.status === 400) {
            this.message(e.response.data, {color: "red", centered: true})
          } else {
            this.message("Ошибка переноса в другой регион", {color: "red", centered: true})
          }
        });
      setTimeout(this.getRows, 200);
      this.valueForTransferRegion.selectedTransferRegion = null
      this.valueForTransferRegion.idRowForTransfer = null
    },
    getTables(year) {
      return new Promise((resolve) => {
        axios
          .post("/api/table", {year: year})
          .then(response => {
            this.tables = response.data
            resolve();
          })
          .catch(e => alert(e.toString()));
      })
    },
    getRegions() {
      axios
        .post("/api/regions")
        .then(response => (this.regions = response.data))
        .catch(e => alert(e.toString()));
    },
    getYears() {
      return new Promise((resolve) => {
        axios
          .post("/api/years")
          .then((response) => {
            this.years = response.data
            resolve();
          })
          .catch(e => alert(e.toString()));
      })
    },
    getHeader() {
      this.tableHeader = []
      this.selectHeader = []
      axios
        .post("/api/getHeader", {id_header: this.sendId.id_header})
        .then(response => {
          if (response.data != null) {
            this.createFormHeader({sub: response.data}, "");
            this.tableHeader.unshift({
              text: "Действие",
              value: 'action',
              align: 'center',
              divider: true,
              sortable: false,
            })
          }
        })
        .catch(e => alert(e.toString()));
    },
    getRows() {
      this.loadingData = true
      this.rows = []
      axios
        .post("/api/getRows", this.sendId)
        .then(response => {
          this.rows = response.data
          this.loadingData = false
        })
        .catch(e => alert(e.toString()));
      this.valueForSave = []              // Сбрасываем введеную в форму информацию при смене региона или таблицы
      this.elementForSort = []            // Сбрасываем элемент сортировки
    },

    sendData() {
      axios
        .post("/api/saveRow", this.getRowForSave(this.valueForSave))
        .then(response => {
          if (response.status !== 200) {
            throw response.data;
          }
          if (response.data === "Запись уже существует") {
            this.message(response.data, {color: "red", centered: true})
          } else {
            this.message("Запись сохранена", {color: "green", centered: true})
            setTimeout(this.getRows, 100);
          }
        })
        .catch(e => alert(e.toString()))
    },
    sendForm(data) {
      axios
        .post("/api/calculateFinalValue", data)
        .then(response => {
          if (response.status === 200) {
            this.message("Запись сохранена", {color: "green", centered: true})
            this.ignoreValue = ''
            setTimeout(this.getRows, 100);
          }
        })
        .catch(e => {
          this.message(e.response.data, {color: "red", centered: true})
        })
    },
  },
  async created() {
    await this.getYears();
    this.selectedYear = this.years[0]
    await this.getTables(this.selectedYear);
    await this.getRegions();
  }
}
</script>
<style>
.v-data-table th {
  white-space: nowrap;
  background-color: hsl(0deg 0% 0% / 8%);
}

.v-data-table {
  margin-block: 10px;
}

.v-data-footer {
  background-color: hsl(0deg 0% 0% / 8%);
}

.v-data-table .v-data-table__wrapper {
  -webkit-scroll-snap-type: block;
  transform: rotateX(180deg);
}

.v-data-table table {
  transform: rotateX(180deg);
}

.myBtn {
  background-color: rgb(12 231 231) !important;
}
</style>
