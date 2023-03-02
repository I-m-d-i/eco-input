<template>
  <div style="padding-top: 0; margin-top: 0px;">
    <v-container>
      <v-layout row wrap style="margin-top: 0px">
        <v-flex lg12 md12 sm12 xs12 class="ma-0">
          <v-select prepend-icon="mdi-calendar" v-model="selectedYear" solo :items="years" label="Год"
                    @change="status = 'received-new-year';changeHandler()">
            <!--Кнопка "Добавить таблицы для следующего года-->
            <template v-slot:prepend-item>
              <v-list-item ripple
                           dense
                           @mousedown.prevent
                           @click="status = 'addition-new-year';!isSaved?(delayedFunction=function (){this.windowAddNewYear=true},windowAlert=true):(windowAddNewYear = true)">
                <v-list-item-content style="padding: 0px">
                  <v-list-item-title>
                    <v-icon color="#2A5885">
                      mdi-plus
                    </v-icon>
                    <span style="color: #2A5885"> Добавить год</span>
                  </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
              <v-divider class="mt-2"></v-divider>
            </template>
          </v-select>
          <v-row v-if="selectedYear" class="ma-0">
            <v-select solo prepend-icon="mdi-table-large" v-model="selectedTableName"
                      :items="tablesListName" label="Таблица"
                      @change="status = 'received-new-table';changeHandler()">
              <!--Кнопка "Добавить новую таблицу"-->
              <template v-slot:prepend-item>
                <v-list-item ripple
                             dense
                             @mousedown.prevent
                             @click.stop="status = 'addition-new-table';!isSaved?(delayedFunction=function (){this.windowAddNewTable=true},windowAlert=true):windowAddNewTable=true">
                  <v-list-item-content style="padding: 0px">
                    <v-list-item-title>
                      <v-icon color="#2A5885">
                        mdi-plus
                      </v-icon>
                      <span style="color: #2A5885"> Добавить новую таблицу</span>
                    </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
                <v-divider class="mt-2"></v-divider>
              </template>
            </v-select>
            <v-menu
              offset-y
              close-on-content-click
            >
              <template v-slot:activator="{ on, attrs }">
                <v-btn style="margin: 5px"
                       v-bind="attrs"
                       v-on="on">
                  <v-icon
                    color="rgb(0 139 139)"
                    class="mdi-cog-outline">
                    mdi-cog-outline
                  </v-icon>
                  Настройки
                </v-btn>
              </template>
              <v-list>
                <v-list-item v-if="selectedTableName"
                             @click.stop="newNameTable=selectedTableName;windowEditNameTable=true">
                  <v-list-item-icon>
                    <v-icon
                      color="rgb(0, 139, 139)">
                      mdi-rename-box
                    </v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>Редактировать</v-list-item-title>
                </v-list-item>
                <v-list-item v-if="selectedTableName"
                             @click.stop="windowEditPermissions=true;getGroups(); selectedGroups = tables.filter(item => item.name === selectedTableName)[0]?.groups?.slice(0)">
                  <v-list-item-icon>
                    <v-icon
                      color="rgb(0, 139, 139)">
                      mdi-lock
                    </v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>Изменить права доступа</v-list-item-title>
                </v-list-item>
                <v-list-item @click.stop="windowSortTables=true; sortableTables=tables.slice(0)">
                  <v-list-item-icon>
                    <v-icon
                      color="rgb(0, 139, 139)">
                      mdi-sort
                    </v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>Сортировать</v-list-item-title>
                </v-list-item>
                <v-list-item v-if="selectedTableName" @click.stop="windowDeleteTable=true">
                  <v-list-item-icon>
                    <v-icon
                      color="rgb(0, 139, 139)">
                      mdi-delete
                    </v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>Удалить</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-row>
        </v-flex>
        <v-flex lg12 md12 sm12 xs12>
          <v-row v-if="selectedTableName" class="ma-1">
            <!--Блок кнопок-->
            <v-sheet color="#d3d3d3" class="ma-1" style="width: max-content" elevation="1">
              <v-btn min-height="41" class="pa-2" title="Добавить ячейку слева" tile depressed
                     @click="actionHandlerWithHeader('add-before')">
                <v-icon>mdi-table-column-plus-before</v-icon>
              </v-btn>
              <v-btn min-height="41" class="pa-2" title="Добавить ячейку справа" tile depressed
                     @click="actionHandlerWithHeader('add-after')">
                <v-icon>mdi-table-column-plus-after</v-icon>
              </v-btn>
              <v-btn min-height="41" class="pa-2" title="Добавить ячейку снизу" tile depressed
                     @click="actionHandlerWithHeader('add-level')">
                <v-icon>mdi-table-row-plus-after</v-icon>
              </v-btn>
              <v-btn min-height="41" class="pa-2" title="Удалить ячейку" tile depressed
                     @click="actionHandlerWithHeader('delete')">
                <v-icon>mdi-table-column-remove</v-icon>
              </v-btn>
              <v-btn min-height="41" class="pa-2 ml-1" title="Отменить действие" :disabled="historyChanges.length===1"
                     tile depressed
                     @click="undoAction()">
                <v-icon>mdi-arrow-u-left-top</v-icon>
              </v-btn>
              <v-btn min-height="41" class="pa-2" title="Вернуть действие" :disabled="historyChangesCanceled.length===0"
                     tile
                     depressed
                     @click="returnUndoAction()">
                <v-icon>mdi-arrow-u-right-top</v-icon>
              </v-btn>
            </v-sheet>
            <v-spacer/>
            <v-btn class="ma-1" :disabled="isSaved"
                   @click="status ='saving-header'; windowSaveHeader=true">
              Сохранить
            </v-btn>
          </v-row>
        </v-flex>
        <div v-if="selectedTableName" class="table-wrapper Flipped">
          <table id="main-table" class="data-table Content mt-2">
            <thead class="table-main-header subheading">
            <tr v-for="(headers,lvl) in headerLevels">
              <th v-for="(header, index) in headers" :key="header.id" :colspan="header.colspan"
                  :rowspan="header.rowspan"
                  class="cell" v-bind:class="{active: isActiveCell(header.id)}"
                  @click="selectedCellId!==header.id?(selectedCellId=header.id):'';$refs[header.id][0].focus()"
              >
                  <v-textarea
                    :ref="header.id"
                    class="textArea1"
                    v-model.lazy="header.text"
                    hide-details
                    no-resize
                    auto-grow
                    flat
                    solo
                    cols="14"
                    rows="1"
                    @input="updateTextInColumn(header.text,header.id)"
                    @focusin="selectedCellId=header.id"
                  />
              </th>
            </tr>
            </thead>
          </table>
        </div>
      </v-layout>
      <!--Форма изменения имени таблицы-->
      <v-dialog v-if="windowEditNameTable" v-model="windowEditNameTable"  width="500" persistent>
        <v-card>
          <v-alert :value="true"> Переименование
          </v-alert>
          <v-card-text class="text-md-center">
            <v-form v-model="valid">
              <v-textarea
                outlined
                label="Название"
                v-model="newNameTable"
                :rules="[rules.duplicate]"
                rows="1"
                auto-grow
                dense>
              </v-textarea>
            </v-form>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn :disabled="!valid"
                   @click=" windowEditNameTable = false; renameTable(); selectedTableName = newNameTable; newNameTable=''">
              Сохранить
            </v-btn>
            <v-btn @click="newNameTable='';windowEditNameTable = false"> Отменить
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!--Форма изменения прав доступа к таблице-->
      <v-dialog v-if="windowEditPermissions" v-model="windowEditPermissions" width="500" persistent>
        <v-card>
          <v-alert> Изменение права доступа</v-alert>
          <v-card-text style="min-height: 250px"  class="text-md-center">
            <v-select chips label="Права доступа" multiple solo :value="selectedGroups"
                      :items="userGroups"
                      item-text="name" item-value="code" hide-details @input="changeSelectedGroups($event)"/>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn
              @click="saveGroups(selectedGroups, sendIdHeader)">
              Сохранить
            </v-btn>
            <v-btn @click="windowEditPermissions = false; selectedGroups=[]"> Отменить
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!--Форма добавления новой таблицы-->
      <v-dialog v-if="windowAddNewTable" v-model="windowAddNewTable"  width="500" persistent>
        <v-card>
          <v-alert :value="true"> Добавление новой таблицы
          </v-alert>
          <v-card-text style="overflow-y: scroll" class="text-md-center">
            <v-form v-model="valid">
              <v-textarea
                outlined
                label="Название"
                v-model="newNameTable"
                :rules="[rules.duplicate, rules.required]"
                rows="1"
                auto-grow
                dense>
              </v-textarea>
              <v-textarea
                outlined
                label="Количество колонок"
                v-model="amountColInNewTable"
                :rules="[rules.required,rules.number,rules.counter]"
                rows="1"
                auto-grow
                dense>
              </v-textarea>
            </v-form>
            <v-select style="max-height: 50px" label="Право доступа" :value="selectedGroups" full-width multiple
                      outlined :items="userGroups" @input="changeSelectedGroups($event)"
                      item-value="code" chips item-text="name" hide-details/>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn :disabled="!valid"
                   @click="changeHandler(); windowAddNewTable = false">
              Добавить
            </v-btn>
            <v-btn @click="newNameTable=''; amountColInNewTable = null ;windowAddNewTable = false"> Отмена
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!--Форма изсенения сортировки таблицы-->
      <v-dialog v-if="windowSortTables" v-model="windowSortTables" width="600" scrollable persistent>
        <v-card>
          <v-alert :value="true"> Изменение сортировки таблиц
          </v-alert>
          <v-card-text class="text-md-center">
            <v-list>
              <template v-for="(item, index) in sortableTables">
                <v-list-item :key="item.sort">
                  <v-list-item-content class="pa-0">
                    <v-textarea
                      :value=item.name
                      no-resize
                      flat
                      solo
                      readonly
                      auto-grow
                      hide-details
                      class="pa-0"
                      rows="1"
                    ></v-textarea>
                  </v-list-item-content>
                  <v-icon
                    v-if="index!==0"
                    title="Переместить выше"
                    color="#990000"
                    @click="sortUp(index)">
                    mdi-arrow-up-thin
                  </v-icon>
                  <v-icon
                    v-if="index!==sortableTables.length-1"
                    title="Переместить ниже"
                    color="#990000"
                    @click="sortDown(index)">
                    mdi-arrow-down-thin
                  </v-icon>
                </v-list-item>
                <v-divider/>
              </template>
            </v-list>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn
              @click=" saveSort();windowSortTables = false">
              Применить
            </v-btn>
            <v-btn @click="windowSortTables = false"> Отмена
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!--Форма предупреждения о потере изменений-->
      <v-dialog v-if="windowAlert" v-model="windowAlert"  width="500" persistent>
        <v-card>
          <v-alert type="warning" :value="true"> Внимание!
          </v-alert>
          <v-card-text class="text-md-center">
            <v-textarea
              readonly
              value="Все не сохраненные изменения в шапке таблицы отменятся. Вы уверены что хотите продолжить?"
              rows="1"
              flat
              solo
              auto-grow
              dense>
            </v-textarea>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn @click="windowAlert = false; delayedFunction()">
              Да
            </v-btn>
            <v-btn
              @click="
              windowAlert = false;
              selectedTableName = oldSelectedNameTable;
              selectedYear = oldSelectedYear;
              newNameTable=''
              amountColInNewTable = null">
              Нет
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!--Форма удаления таблицы -->
      <v-dialog v-if="windowDeleteTable" v-model="windowDeleteTable"  width="500" persistent>
        <v-card>
          <v-alert color="red" :value="true"> Внимание!
          </v-alert>
          <v-card-text class="text-md-center">
            <v-textarea
              readonly
              value="Отменить удаление таблицы нельзя. Вы уверены что хотите продолжить?"
              rows="1"
              flat
              solo
              auto-grow
              dense>
            </v-textarea>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn @click="windowDeleteTable = false; deleteTable()">
              Да
            </v-btn>
            <v-btn
              @click="
              windowDeleteTable = false;">
              Нет
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!--Форма сохранения шапки-->
      <v-dialog v-if="windowSaveHeader" v-model="windowSaveHeader"  width="500" persistent>
        <v-card>
          <v-alert color="warning" :value="true"> Внимание!
          </v-alert>
          <v-card-text class="text-md-center">
            <v-textarea
              readonly
              value="Вы уверены, что хотите сохранить изменения?"
              rows="1"
              flat
              solo
              auto-grow
              dense>
            </v-textarea>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn @click="changeHandler();windowSaveHeader = false;">
              Да
            </v-btn>
            <v-btn
              @click="
              windowSaveHeader = false;">
              Нет
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!--Форма добавления нового года-->
      <v-dialog v-if="windowAddNewYear" v-model="windowAddNewYear"  max-width="600px" persistent>
        <v-card>
          <v-alert type="info" :value="true"> Добавление нового года
          </v-alert>
          <v-card-text>
            <v-form v-model="valid">
              <v-textarea
                outlined
                label="Укажите добавляемый год"
                v-model.number="addingYear"
                :rules="[rules.required,rules.duplicateYear,rules.number,rules.counterYear]"
                rows="1"
                auto-grow
                dense
              />
              <v-switch
                v-model="switchDuplicate"
                inset
                label='Использовать таблицы другого года'
              />
              <template v-if="switchDuplicate">
                Использовать таблицы какого года?
                <v-select
                  solo
                  v-model.number="selectedYearForCloning"
                  :items="years"
                  label="Год"
                  :rules="[rules.required]"
                >
                </v-select>
              </template>
            </v-form>
          </v-card-text>
          <v-card-actions class="text-lg-right">
            <v-spacer/>
            <v-btn :disabled="!valid"
                   @click="changeHandler();windowAddNewYear = false"> Добавить
            </v-btn>
            <v-btn
              @click="windowAddNewYear = false;switchDuplicate=false;selectedYearForCloning = null; addingYear=null">
              Отменить
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </div>
</template>
<script>
import axios from "axios";
import {bus} from "@/main";
import {debounce} from "@/components/FormEdit/debounce";

export default {
  data() {
    return {
      windowEditPermissions: false,
      windowEditNameTable: false,
      windowAddNewTable: false,
      windowAddNewYear: false,
      windowAlert: false,
      windowSaveHeader: false,
      windowDeleteTable: false,
      windowSortTables: false,

      switchDuplicate: false,
      selectedYearForCloning: null,

      status: '',
      valid: false,
      savedVersion: null,
      version: null,

      header: [],
      years: [],
      tables: [{groups: [], id: null, name: '', sort: null}],
      userGroups: [],
      selectedGroups: [],

      headerLevels: [],
      listUUID: [],
      listUUIDHistory: [],

      newNameTable: '',
      selectedTableName: '',
      selectedYear: null,
      selectedCellId: null,
      amountColInNewTable: null,
      oldSelectedYear: null,
      oldSelectedNameTable: null,
      sortableTables: null,
      addingYear: null,

      historyChanges: [],
      historyChangesCanceled: [],

      rules: {
        duplicate: value => {
          for (const tableName of this.tablesListName) {
            if (tableName.toLocaleLowerCase() === value.toLocaleLowerCase()) {
              return "Таблица с таким именем уже существует"
            }
          }
          return true
        },
        number: value => {
          const pattern = /^\d+$/
          if (pattern.test(value)) {
            return true
          } else {
            return 'Можно использовать только числами'
          }
        },
        required: value => !!value || 'Поле обязательно для заполнения',
        counter: value => value > 0 && value <= 100 || 'Колонок может быть быть от 1 до 100',
        counterYear: value => value >= 2017 && value <= Math.max.apply(null, this.years) + 1 || `Год может быть быть от 2017 до ${String(Math.max.apply(null, this.years) + 1)}`,
        duplicateYear: value => {
          for (const year of this.years) {
            if (year === Number(value)) {
              return "Такой год уже есть"
            }
          }
          return true
        },
      },
    }
  },
  computed: {
    tablesListName() {
      return this.tables.map(item => item.name);
    },
    sendIdHeader() {
      return this.tables.filter(item => item.name === this.selectedTableName)[0].id
    },
    isSaved() {
      if (this.savedVersion == null && this.historyChanges.length ===1) {
        return true
      }
      return this.savedVersion === this.version
    },
  },
  methods: {
    changeSelectedGroups(groups) {
      if (groups.length > 0) {
        if (!groups.includes('60260cf42e8f13c8d8ec5994', 0)) {
          groups.push('60260cf42e8f13c8d8ec5994')
        } else {
          if (groups.length + 1 === this.selectedGroups.length && groups.length === 1) {
            groups = []
          }
        }
      }
      this.selectedGroups = groups
    },
    saveGroups(groups, idHeader) {
      return new Promise(resolve => {
        if (groups.length > 0 && !groups.includes('60260cf42e8f13c8d8ec5994', 0)) {
          groups.push('60260cf42e8f13c8d8ec5994')
        }
        axios
          .post("/api/table/saveGroups", {groups: groups, idHeader: idHeader})
          .then(response => {
            this.selectedGroups = []
            this.windowEditPermissions = false
            this.getTables(this.selectedYear)
            resolve('resolved')
          })
          .catch(e => alert(e.toString()));
      })
    },
    addNewYear(year, cloningYear) {
      axios
        .post("/api/addYear", {year: year, cloningYear: cloningYear})
        .then(() => {
            bus.$emit('message', {
              message: "Новый год добавлен",
              color: "green",
            })
            this.years.unshift(this.addingYear)
            this.years.sort(function (a, b) {
              return a - b;
            })
            this.selectedYear = year;
            this.getTables(this.selectedYear);
          }
        )
        .catch(e => {
          bus.$emit('message', {
            message: "Ошибка добавления нового года \n" + e.toString(),
            color: "red",
          })
          this.selectedYearForCloning = null;
          this.addingYear = null;
        });
    },

    sortUp(index) {
      this.sortableTables[index].sort--
      this.sortableTables[index - 1].sort++
      const tempTable1 = this.sortableTables[index - 1]
      const tempTable2 = this.sortableTables[index]
      this.sortableTables[index] = tempTable1
      this.sortableTables[index - 1] = tempTable2
    },
    sortDown(index) {
      this.sortableTables[index].sort++
      this.sortableTables[index + 1].sort--
      const tempTable1 = this.sortableTables[index + 1]
      const tempTable2 = this.sortableTables[index]
      this.sortableTables[index] = tempTable1
      this.sortableTables[index + 1] = tempTable2
    },
    saveSort() {
      axios
        .post("/api/saveSortTables", {tables: this.sortableTables})
        .then(() => {
          bus.$emit('message', {
            message: "Сортировка таблиц успешно изменена",
            color: "green",
          })
          this.getTables(this.selectedYear)
          this.sortableTables = []
        })
        .catch(e => {
          bus.$emit('message', {
            message: "Ошибка при изменении сортировки таблиц \n" + e.toString(),
            color: "red",
          })
        });
    },
    saveHeader() {
      axios
        .post("/api/saveHeader", {id: this.sendIdHeader, header: this.header})
        .then(response => {
          bus.$emit('message', {
            message: "Сохранено",
            color: "green",
          })
        })
        .catch(e => bus.$emit('message', {
          message: "Ошибка сохранения \n" + e.toString(),
          color: "red",
        }));
    },
    deleteTable() {
      axios
        .post("/api/deleteTable", {id: this.sendIdHeader})
        .then(response => {
          if (response.data === "В таблицу внесены данные, удаление невозможно") {
            bus.$emit('message', {
              message: response.data,
              color: "red",
            })
          } else {
            bus.$emit('message', {
              message: "Успешно удалено",
              color: "green",
            })
            this.selectedTableName = '';
            this.clearHistory()
            this.getTables(this.selectedYear);
          }
        })
        .catch(e => bus.$emit('message', {
          message: "Ошибка удаления \n" + e.toString(),
          color: "red",
        }));
    },
    addNewTable() {
      axios
        .post("/api/addTable", {
          year: this.selectedYear,
          name: this.newNameTable,
          header: this.header,
          groups: this.selectedGroups
        })
        .then(async response => {
          if (response.data === "Таблица с таким именем уже есть") {
            bus.$emit('message', {
              message: response.data,
              color: "red",
            })
          } else {
            bus.$emit('message', {
              message: "Таблица " + this.newNameTable + " добавленна",
              color: "green",
            })
            this.selectedTableName = this.newNameTable
            this.oldSelectedNameTable = this.selectedTableName
            this.newNameTable = ''
            this.selectedGroups = []
            await this.getTables(this.selectedYear);
            this.getHeader(this.sendIdHeader);
          }
        })
        .catch(e => bus.$emit('message', {
          message: "Ошибка добавления таблицы \n" + e.toString(),
          color: "red",
        }));
    },
    renameTable() {
      axios
        .post("/api/renameTable", {id: this.sendIdHeader, newName: this.newNameTable})
        .then(() => {
          bus.$emit('message', {
            message: "Таблица переименованна",
            color: "green",
          })
          this.getTables(this.selectedYear);
        })
        .catch(e => bus.$emit('message', {
          message: "Ошибка переименования таблицы \n" + e.toString(),
          color: "red",
        }));
    },
    changeHandler() {
      switch (this.status) {
        case 'received-new-year':
          this.delayedFunction = () => {
            this.oldSelectedYear = this.selectedYear
            this.clearHistory()
            this.header = []
            this.headerLevels = []
            this.selectedTableName = ''
            this.getTables(this.selectedYear)
          }
          break
        case 'addition-new-year':
          this.delayedFunction = () => {
            this.clearHistory()
            this.header = []
            this.headerLevels = []
            this.selectedTableName = ''
            if (this.switchDuplicate) {
              this.addNewYear(this.addingYear, this.selectedYearForCloning)
            } else {
              this.selectedYear = this.addingYear
              this.addingYear = null
              this.tables = []
            }
          }
          break
        case 'received-new-table':
          this.delayedFunction = () => {
            this.oldSelectedNameTable = this.selectedTableName
            this.clearHistory()
            this.getHeader(this.sendIdHeader)
          }
          break
        case 'addition-new-table':
          this.delayedFunction = () => {
            this.header = []
            for (let i = 0; i < this.amountColInNewTable; i++) {
              this.header.push({sub: [], text: ''})
            }
            this.clearHistory()
            this.addNewTable()
            this.amountColInNewTable = null
          }
          break
        case 'saving-header':
          this.delayedFunction = () => {
            this.saveHeader()
            this.savedVersion = this.version
          }
          break
        default:
          break
      }
      if (!this.isSaved&&this.tables) {
        this.windowAlert = true
      } else {
        this.delayedFunction()
        this.status = ''
      }

    },

    getState(){
      return {
        id: this.version,
        header: structuredClone(this.header),
        listUUID: this.listUUID.slice(0)
      }
    },
    actionHandlerWithHeader(action) {
      if (this.selectedCellId != null) {
        let id = null
        switch (action) {
          case 'delete':
            if (this.header.length > 1) {
              id = this.deleteColumn()
            } else {
              return
            }
            break
          case 'add-after':
            id = this.addColumn(1)
            break
          case 'add-before':
            id = this.addColumn(0)
            break
          case 'add-level':
            this.addLevelColumn(this.selectedCellId)
            break
          default:
            break
        }
        this.addInHistory(this.getState())
        this.genLevelHeader(this.header)
        if (id!=null){
           this.$nextTick(() => {
            this.$refs[id][0].focus()
          })
        }
      }
    },
    updateTextInColumn(text, id){
      this.AddTextChangeToHistory()
      const column = this.findColumn(id, this.header)
      if (column.text !== text) {
        column.text = text
      }
    },
//TODO переписать в общий обработчик действий
    AddTextChangeToHistory: debounce(function (){
      this.addInHistory(this.getState())
    }, 500,false),
    addLevelColumn(id) {
      const column = this.findColumn(id, this.header)
      if (column != null) {
        column['sub'].push({id: this.genUUID(this.listUUID), sub: [], text: ''})
      }
    },
    addColumn(shift) {
      const siblings = this.getSiblings(this.header, this.selectedCellId, 0)
      for (let i = 0; i < siblings.length; i++) {
        if (siblings[i].id === this.selectedCellId) {
          const newId = this.genUUID(this.listUUID)
          siblings.splice(i + shift, 0, {id: newId, sub: [], text: ''})
          return newId
        }
      }
    },
    deleteColumn() {
      const siblings = this.getSiblings(this.header, this.selectedCellId, 0)
      for (let i = 0; i < siblings.length; i++) {
        if (siblings[i].id === this.selectedCellId) {
          siblings.splice(i, 1)
          if (siblings.length === 0) {
            this.selectedCellId = null
            return null
          } else if (siblings.length === i) {
            //TODO Можно переделать на отслеживание изменения selectCol и последующее изменение фокуса из него
            return siblings[i - 1].id
          } else {
            return  siblings[i].id
          }
        }
      }
      this.listUUID.splice(this.listUUID.indexOf(this.selectedCellId), 1)
    },

    //Получение братьев и сестер
    getSiblings(header, childId, lvl) {
      for (const column of header) {
        if (column.id === childId) {
          if (lvl === 0) {
            return header
          }
          return true
        } else {
          if (column['sub'].length > 0) {
            const req = this.getSiblings(column['sub'], childId, lvl + 1)
            if (req === true) {
              return column['sub']
            }
            if (req != null) {
              return req
            }
          }
        }
      }
      return null
    },
    //Получение ссылки на колонку
    findColumn(id, header) {
      for (const headerElement of header) {
        if (headerElement.id === id) {
          return headerElement
        } else {
          if (headerElement['sub'].length > 0) {
            const col = this.findColumn(id, headerElement['sub'])
            if (col != null) {
              return col
            }
          }
        }
      }
      return null
    },

    //Отменить действие
    undoAction() {
      this.historyChangesCanceled.push(this.historyChanges.at(-1))
      this.header = this.historyChanges[this.historyChanges.length - 2].header
      this.listUUID = this.historyChanges[this.historyChanges.length - 2].listUUID
      this.version = this.historyChanges[this.historyChanges.length - 2].id
      this.selectedCellId = null
      this.genLevelHeader(this.header)
      this.historyChanges.pop()
    },
    //Вернуть действие
    returnUndoAction() {
      this.historyChanges.push(this.historyChangesCanceled.at(-1))
      this.version = this.historyChangesCanceled.at(-1).id
      this.header = this.historyChangesCanceled.at(-1).header
      this.listUUID = this.historyChangesCanceled.at(-1).listUUID
      this.selectedCellId = null
      this.genLevelHeader(this.header)
      this.historyChangesCanceled.pop()
    },
    //Очистка истории
    clearHistory() {
      this.historyChanges = []
      this.historyChangesCanceled = []
      this.listUUIDHistory = []
    },
    addInHistory(state) {
      this.historyChangesCanceled = []
      this.historyChanges.push(state)
      this.version = this.genUUID(this.listUUIDHistory)
    },

    isActiveCell: function (id) {
      if (this.selectedCol !== null) {
        return this.selectedCol.id === id
      } else {
        return false
      }
    },

    message(message, options) {
      bus.$emit('message', {
        message: message,
        color: options.color,
        centered: options.centered
      })
    },
    //----------------------------------------
    //Генерация массива уровней из дерева
    genLevelHeader(header) {
      this.headerLevels = []
      let maxLvl = 1
      for (let i = 0; i < header.length; i++) {
        if (header[i]['sub'].length > 0) {
          let numberOfLvl = this.rowLength(header[i]['sub'], 1)
          if (maxLvl < this.rowLength(header[i]['sub'], 1)) {
            maxLvl = numberOfLvl
          }
        }
      }
      this.createLevelHeader(header, maxLvl)
    },
    createLevelHeader(header, rowspanLen) {
      let headerLevel = []
      let newChild = []
      for (const headerElement of header) {
        if (headerElement['sub'].length > 0) {
          newChild = newChild.concat(headerElement['sub'])
          headerLevel.push({
            id: headerElement.id,
            text: headerElement.text,
            colspan: this.colLength(headerElement['sub']),
            rowspan: 1
          })
        } else {
          headerLevel.push({id: headerElement.id, text: headerElement.text, rowspan: rowspanLen, colspan: 1})
        }
      }
      this.headerLevels.push(headerLevel)
      if (newChild.length > 0) {

        this.createLevelHeader(newChild, rowspanLen--)
      }
    },
    //Подсчет colspan
    colLength(header) {
      let colspanLen = 0
      for (let i = 0; i < header.length; i++) {
        if (header[i]['sub'].length > 0) {
          colspanLen += this.colLength(header[i]['sub'])
        } else {
          colspanLen++
        }
      }
      return colspanLen
    },
    //Подсчет rowspan
    rowLength(header, Lvl) {
      Lvl++
      let maxLevel = Lvl
      for (const headElement of header) {
        if (headElement['sub'].length > 0) {
          const sublevelLength = this.rowLength(headElement['sub'], Lvl)
          if (sublevelLength > maxLevel) {
            maxLevel = sublevelLength
          }
        }
      }
      return maxLevel
    },
    //Установка уникальных идентификаторов для каждой колонки
    addIdForHeader(header) {
      for (let i = 0; i < header.length + 0; i++) {
        header[i].id = this.genUUID(this.listUUID)
        if (header[i]['sub'].length > 0) {
          this.addIdForHeader(header[i]['sub'])
        }
      }
    },
    //Генерация уникального идентификатора
    genUUID(listUUID) {
      let UUID = Math.floor(Math.random() * 10000 + 1000)
      if (listUUID.includes(UUID)) {
        UUID = this.genUUID(listUUID)
      } else {
        listUUID.push(UUID)
      }
      return UUID
    },
    //---------------------------------------
    getHeader(id_header) {
      this.header = []
      this.headerLevels = []
      this.listUUID = []
      this.listUUIDHistory = []
      this.savedVersion = null
      axios
        .post("/api/getHeader", {id_header: id_header})
        .then(response => {
          this.header = response.data;
          this.version = this.genUUID(this.listUUIDHistory)
          this.addIdForHeader(this.header);
          this.genLevelHeader(this.header)
          this.addInHistory(this.getState())
        })
        .catch(e => alert(e.toString()));
    },
    getTables(sendYear) {
      this.tables = []
      return new Promise(resolve => {
        axios
          .post("/api/table", {year: sendYear})
          .then(response => {
            this.tables = response.data;
            resolve('resolved')
          })
          .catch(e => alert(e.toString()));
      })
    },
    getYears() {
      axios
        .post("/api/years")
        .then(response => (this.years = response.data))
        .catch(e => alert(e.toString()));
    },
    getGroups() {
      axios
        .post("/api/permission-group/get/all")
        .then(response => (this.userGroups = response.data))
        .catch(e => alert(e.toString()));
    },
  },
  created() {
    this.getYears()
  }
}
</script>
<style scoped>
.Flipped,
.Flipped .Content {
  transform: rotateX(180deg);
  -ms-transform: rotateX(180deg); /* IE 9 */
  -webkit-transform: rotateX(180deg); /* Safari and Chrome */
}

.table-wrapper {
  overflow-x: scroll;
}

table {
  width: 100%;
  border-color: #ffffff;
  border-collapse: collapse;
  margin-right: 10px;
}

.table-main-header >>> th {
  border: 1px solid #bdbdbd;
  text-align: center;
  background-color: #ffffff;
  color: #424242;
  padding: 2px;
  cursor: text;
}

.cell:hover {
  background-color: #e0e0e0;
}

.textArea1 >>> .v-input__slot {
  background-color: inherit !important;
  padding: 0 !important;
}

.cell.active {
  background-color: rgba(0, 194, 194, 0.75);
}

.textArea1 >>> textarea {
  text-align: center;
  vertical-align: center;
  font-size: 14px;
  width: max-content;
}
</style>
