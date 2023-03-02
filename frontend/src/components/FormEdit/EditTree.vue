<template>
  <v-container>
    <v-layout row wrap style="margin-top: 0px">
      <v-flex lg7 md12 sm12 xs12 style="padding-inline: 10px;">
        <v-select hide-details :disabled="!activeForm" :readonly="isAddedNewStruct" solo v-model="selectedYear"
                  :items="yearsTree"
                  label="Год" @click="isAddedNewStruct ? messageAboutCopy():null"
                  @input="getTree($event);getTables($event)">
          <template v-slot:prepend-item>
            <v-list-item ripple
                         dense
                         @mousedown.prevent
                         @click.stop="windowForAddNewStruct=true;">
              <v-list-item-content style="padding: 0px">
                <v-list-item-title>
                  <v-icon color="#2A5885">
                    mdi-plus
                  </v-icon>
                  <span style="color: #2A5885"> Добавить структуру для следующего года</span>
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-divider class="mt-2"></v-divider>
          </template>
        </v-select>
        <!--Форма дерева-->
        <v-treeview
          dense
          :items='treeForForm'
          :open.sync="openIds"
          open-all
          item-key="id"
          style="overflow-y: scroll; height: 80vh; margin-block:5px "
        >
          <template v-slot:prepend="{item,open}">
            <!--Выравнивание таблиц в зависимости от выбранного действия-->
            <template v-if="item.lvl!==-1&&item.idHeader">
              <i v-bind:style="{marginLeft: sort?'67px':transfer?'36px':'50px'}"></i>
            </template>
            <template v-if="!sort&&!transfer&&item.lvl!==-1">
              <!--Кнопки только для папок-->
              <template v-if="item.idHeader === null">
                <v-icon :disabled="!activeForm" small class="mdi-plus-thick" color="#8A00FF" title="Добавить"
                        @click="(windowForAdd = true); (selectedElement = item)">mdi-plus-thick
                </v-icon>
                <v-icon :disabled="!activeForm" small class="mdi-pencil" color="#00AAAA" title="Переименовать"
                        @click="(windowForEdit=true); (selectedElement = item);newNameItem=item.name"> mdi-pencil
                </v-icon>
                <v-icon small class="mdi-compare-vertical" color="#ff9900" title="Сортировать"
                        :disabled="!activeForm"
                        @click="(selectedElement = item); sort = true">mdi-compare-vertical
                </v-icon>
              </template>
              <v-icon :disabled="!activeForm" small class="mdi-transfer" color="#008F00"
                      title="Отправить в другую папку"
                      @click="transfer=true; (selectedElement = item)">mdi-transfer
              </v-icon>
              <v-icon :disabled="!activeForm" small class="mdi-delete" color="#FF0000" title="Удалить"
                      @click.stop="(windowForDelete = true); (selectedElement = item)"> mdi-delete
              </v-icon>
            </template>
            <!--Кнопки при перемещении-->
            <template v-if="transfer&&item.lvl!==-1&&!item.idHeader">
              <v-icon small class="mdi-backspace" color="#990000" title="Отменить перенос"
                      @click="transfer=false; selectedElement= null"> mdi-cancel
              </v-icon>
              <v-icon small class="mdi-compare_arrows" color="#ff9900" title="Поместить в эту папку"
                      :disabled="isElementInTree(item.id,selectedElement)"
                      @click="to = item; transferElement(); transfer = false">mdi-arrow-right-thin
              </v-icon>
            </template>
            <!--Кнопки при сортировки-->
            <template v-if="sort&&item.lvl!==-1&&item.idHeader === null">
              <v-icon small class="mdi-backspace" color="#990000" title="Отменить сортировку"
                      @click="sort=false; selectedElement = null"> mdi-cancel
              </v-icon>
              <v-icon
                :disabled="(!isElementInTree(item.id,selectedParent)||(item.lvl!==selectedElement.lvl)||(getIndexChildren(selectedParent,item.id)===0))"
                color="#990000"
                title="Переместить выше"
                @click="sortUp(item)">
                mdi-arrow-up-thin
              </v-icon>
              <v-icon
                :disabled="(!isElementInTree(item.id,selectedParent)||(item.lvl!==selectedElement.lvl)||(getIndexChildren(selectedParent,item.id)===selectedParent.children.length-1))"
                title="Переместить ниже"
                color="#990000"
                @click="sortDown(item)">
                mdi-arrow-down-thin
              </v-icon>
            </template>
            <!--Для корневой папки ставим только кнопку добавить           -->
            <template v-if="item.lvl===-1">
              <v-icon :disabled="!activeForm" small class="mdi-plus-thick" color="#8A00FF" title="Добавить"
                      @click="(windowForAdd = true); (selectedElement = item)">mdi-plus-thick
              </v-icon>
            </template>
            <v-icon :disabled="!activeForm" v-if="item.idHeader===null" class="">
              {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
            </v-icon>
            <v-icon :disabled="!activeForm" v-else>
              mdi-table-large
            </v-icon>
          </template>
          <template v-slot:label={item}>
            <v-tooltip bottom>
              <template v-slot:activator="{on, attrs }">
                <span slot="activator" v-bind="attrs" v-on="on"> {{ item.name }}</span>
              </template>
              <span>{{ item.name }}</span>
            </v-tooltip>
          </template>
        </v-treeview>
      </v-flex>
      <v-flex lg5 md12 sm12 xs12 style="padding-inline: 10px;">
        <v-row style="margin-block: 3px; margin-left: 0px">
          <v-btn :disabled="!activeForm" large @click="windowForAddInDB=true; getPin()">
            Добавить в базу данных
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn :disabled="!activeForm" large v-if="isAddedNewStruct"
                 @click="isAddedNewStruct = false; messagesDisplay.show = false; yearsTree.shift();selectedYear=maxYearTree; getTree(maxYearTree)">
            Отменить добавление
          </v-btn>
        </v-row>
        <!--Форма добавления элементов-->
        <div style="margin-block: 15px" v-if="windowForAdd">
          <v-card class="pa-2">
            <v-card-title>Добавление в папку "{{ this.selectedElement.name }}"</v-card-title>
            <v-divider></v-divider>
            <v-form v-model="valid">
              <v-select
                v-model="selectedType"
                :items="['Папка','Таблица']"
                :rules="[rules.required]"
                label="Тип добавляемого элемента"
                @change="addedElement={id:genUUID(listUUID) ,name: null,idHeader: null, sort: null}"/>
              <v-textarea
                v-if="selectedType==='Папка'"
                outlined
                label="Название папки"
                :rules="[rules.required,rules.duplicate(addedElement.id)]"
                v-model="addedElement.name"
                rows="1"
                auto-grow
                dense/>
              <v-select
                v-if="selectedType==='Таблица'"
                @input="addedElement.idHeader = [$event]; addedElement.name = getNameTable(addedElement); addedElement.sort = getSortTable(addedElement)"
                no-data-text="Таблиц нет"
                :rules="[rules.required]"
                :items="unUsedTables"
                item-text="name"
                item-value="id"
                label="Таблица"
                dense/>
            </v-form>
            <v-btn
              color="red darken-1"
              text
              @click="(selectedType = null); (addedElement={name: null,idHeader: null, sort: null});(windowForAdd = false)"
            >Отмена
            </v-btn>
            <v-btn
              :disabled="!valid"
              color="green darken-1"
              text
              @click="(windowForAdd = false);(selectedType= null); addElement(); (addedElement={name:null, idHeader: null, sort: null})"
            >Сохранить
            </v-btn>
          </v-card>
        </div>
        <!--Форма редактирования имени-->
        <div style="margin-block: 15px" v-if="windowForEdit">
          <v-card class="pa-2">
            <v-card-title>Редактирование</v-card-title>
            <v-divider></v-divider>
            <v-form v-model="valid">
              <v-textarea
                outlined
                label="Название"
                v-model="newNameItem"
                :rules="[rules.required,rules.duplicate]"
                rows="1"
                auto-grow
                dense>
              </v-textarea>
            </v-form>
            <v-btn
              color="red darken-1"
              text
              @click="(windowForEdit = false); newNameItem=''"
            >Отмена
            </v-btn>
            <v-btn
              color="green darken-1"
              :disabled="!valid"
              text
              @click="(windowForEdit = false); selectedElement.name = newNameItem; newNameItem=''"
            >Сохранить
            </v-btn>
          </v-card>
        </div>
      </v-flex>
    </v-layout>
    <!--Оповещение об удалении-->
    <v-dialog v-model="windowForDelete" width="400" persistent>
      <v-card>
        <v-alert type="warning" :value="true"> Внимание
        </v-alert>
        <v-card-text class="text-md-center">
          Вы действительно хотите удалить данный элемент?
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn @click="deleteElement(); windowForDelete = false"> Да
          </v-btn>
          <v-btn @click="windowForDelete = false"> Нет
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!--Диалог добавления новой структуры-->
    <v-dialog v-model="windowForAddNewStruct" max-width="600px" persistent>
      <v-card>
        <v-alert type="info" :value="true"> Добавление новой структуры
        </v-alert>
        <v-card-text>
          <v-form v-model="valid">
            <v-textarea
              outlined
              label="Укажите добавляемый год"
              v-model.number.lazy="addingYear"
              :rules="[rules.required,rules.number,rules.counterYear,rules.duplicateYear]"
              rows="1"
              auto-grow
              dense/>
            <v-switch
              v-model="switchDuplicate"
              inset
              label='Использовать готовую структуру'
            ></v-switch>
            <template v-if="switchDuplicate">
              Дублировать структуру какого года?
              <v-select solo v-model="selectedYearForCloning" :items="yearsTree" label="Год">
              </v-select>
            </template>
          </v-form>
        </v-card-text>
        <v-card-actions class="text-lg-right">
          <v-spacer/>
          <v-btn :disabled="!valid" @click="addStruct()"> Добавить
          </v-btn>
          <v-btn
            @click="windowForAddNewStruct = false;switchDuplicate=false;selectedYearForCloning = null;addingYear=null">
            Отменить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!--    -->
    <v-dialog v-model="windowForAddInDB" max-width="600px" persistent>
      <v-card>
        <v-alert type="info" :value="true"> Подтверждение добавления
        </v-alert>
        <v-card-text>
          Введите код для подтверждения. Код: {{ pin }}
        </v-card-text>
        <v-card-actions>
          <v-otp-input
            v-model="otp"
            :disabled="loading"
            @finish="checkPass"
          ></v-otp-input>
          <v-overlay absolute :value="loading">
            <v-progress-circular
              indeterminate
              color="primary"
            ></v-progress-circular>
          </v-overlay>
        </v-card-actions>
        <v-card-actions style="justify-content: flex-end;">
          <v-btn @click="windowForAddInDB = false; otp= ''"> Отменить
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <!--Оповещение-->
    <v-snackbar
      v-model="messagesDisplay.show"
      :timeout="messagesDisplay.timeout"
      :color="messagesDisplay.color"
    >
      <span v-bind:style="{color:messagesDisplay.textColor }" style="font-size: 1.15rem;">
        {{ messagesDisplay.message }}
      </span>
      <template v-slot:action="{ attrs }">
        <v-icon
          :color="messagesDisplay.textColor"
          v-bind="attrs"
          @click="messagesDisplay.show = false">
          mdi-close-circle-outline
        </v-icon>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      selectedYear: null,
      nowYear: new Date().getFullYear(),
      selectedYearForCloning: null,
      selectedType: null,
      selectedElement: {}, // Выбранный элемент
      addedElement: {},
      to: {},
      openIds: [],
      newNameItem: null,
      pin: null,
      otp: '',
      listUUID: [],

      yearsTree: [],
      yearsHeader: [],
      tables: [],
      tree: {
        id: null,
        idHeader: [],
        children: [],
        name: null,
        sort: null
      },

      windowForAdd: false, // Определяет видимость окна добавления
      windowForDelete: false, // Определяет видимость окна удаления
      windowForEdit: false, // Определяет видимость окна редактирования
      windowForAddNewStruct: false,
      windowForAddInDB: false,
      errorDuplicate: false, // ошибка наличия дуиблика
      isErrorInForm: false,
      transfer: false,
      sort: false,
      switchDuplicate: false,
      isAddedNewStruct: false,
      loading: null,
      valid: false,
      addingYear: null,
      messagesDisplay: {                        // Оповещение о выполненом действие (сохранение, удаление, изменение)
        message: "",
        color: "#4caf50",
        show: false,
        textColor: null,
        timeout: 3000,
      },
      rules: {
        number: value => {
          if (/^\d+$/.test(value)) {
            return true
          } else {
            return 'Можно использовать только числами'
          }
        },
        duplicate: value => {
          return !this.isElementInTree(value, this.tree) || 'Файл с таким именем уже существует'
        },
        required: value => !!value || 'Поле обязательно для заполнения',
        counterYear: value => value >= 2017 || `Год может быть быть от 2017 `,
        duplicateYear: value => this.yearsTree.includes(Number(value), 0) ? "Такой год уже есть" : true,
      }
    }
  },
  computed: {
    unUsedTables() {
      let tables = []
      for (const table of this.tables) {
        if (!this.isTableIsTree(table.id, this.tree)) {
          tables.push(table)
        }
      }
      return tables
    },
    selectedParent() {
      if (this.selectedElement != null) {
        return this.getParent(this.tree, this.selectedElement.id)
      }
    },
    treeForForm() {
      return [this.tree]
    },
    activeForm() {
      return !(this.windowForEdit || this.windowForAdd || this.windowForDelete || this.transfer || this.sort);
    },
    maxYearTree() {
      return Math.max.apply(null, this.yearsTree);
    },
  },
  methods: {
    checkPass() {
      this.loading = true
      setTimeout(() => {
        this.loading = false
        if (this.otp === this.pin.toString()) {
          this.messagesDisplay.color = '#4caf50'
          this.messagesDisplay.textColor = "black"
          this.messagesDisplay.show = true
          this.messagesDisplay.message = "Код верен"
          this.messagesDisplay.timeout = 4000
          this.windowForAddInDB = false;
          this.addTree()
          this.otp = ''
        } else {
          this.messagesDisplay.color = '#d70000'
          this.messagesDisplay.textColor = "black"
          this.messagesDisplay.show = true
          this.messagesDisplay.message = "Код не верен"
          this.messagesDisplay.timeout = 4000
          this.otp = ''
        }
      }, 1000)
    },
    getPin() {
      this.pin = Math.round(Math.random() * 899999 + 100000)
    },

    sortUp(item) {
      const parent = this.selectedParent
      const indexItem = this.getIndexChildren(parent, item.id)
      const item1 = structuredClone(parent.children[indexItem])
      const item2 = structuredClone(parent.children[indexItem - 1])
      item2.sort++
      item1.sort--
      parent.children.splice(indexItem, 1, item2)
      parent.children.splice(indexItem - 1, 1, item1)
    },
    sortDown(item) {
      const parent = this.selectedParent
      const indexItem = this.getIndexChildren(parent, item.id)
      const item1 = structuredClone(parent.children[indexItem])
      const item2 = structuredClone(parent.children[indexItem + 1])
      item2.sort--
      item1.sort++
      parent.children.splice(indexItem, 1, item2)
      parent.children.splice(indexItem + 1, 1, item1)
    },
    addElement() {
      if (this.addedElement.idHeader) {
        //TODO Сделать добавление в порядке сортировки
        this.selectedElement.children.push({
          id: this.genUUID(this.listUUID),
          name: this.addedElement.name,
          lvl: this.selectedElement.lvl + 1,
          children: [],
          idHeader: this.addedElement.idHeader,
          sort: this.addedElement.sort
        })
      } else {
        this.selectedElement.children.unshift({
          id: this.genUUID(this.listUUID),
          name: this.addedElement.name,
          lvl: this.selectedElement.lvl + 1,
          children: [],
          idHeader: this.addedElement.idHeader,
          sort: this.selectedElement.children.length + 1
        })
      }
      this.openIds.push(this.selectedElement.id)
      this.selectedElement = null
    },
    deleteElement() {
      const parent = this.selectedParent
      const index = this.getIndexChildren(parent, this.selectedElement.id)
      parent.children.splice(index, 1)
      this.selectedElement = null
    },
    transferElement() {
      const copySelectedElement = this.selectedElement
      const parent = this.selectedParent
      const index = this.getIndexChildren(parent, copySelectedElement.id)
      parent.children.splice(index, 1)
      if (this.selectedElement.idHeader) {
        this.to.children.push(
          {
            id: copySelectedElement.id,
            name: copySelectedElement.name,
            lvl: this.to.lvl + 1,
            children: copySelectedElement.children,
            idHeader: copySelectedElement.idHeader,
            sort: this.to.children.length
          })
      } else {
        this.to.children.unshift(
          {
            id: copySelectedElement.id,
            name: copySelectedElement.name,
            lvl: this.to.lvl + 1,
            children: copySelectedElement.children,
            idHeader: copySelectedElement.idHeader,
            sort: this.to.children.length
          })
      }
      this.openIds.push(this.to.id)
      this.selectedElement = null
      this.to = null
    },

    getNameTable(table) {
      if (table.idHeader != null) {
        return this.tables.filter(item => item.id === table.idHeader[0])[0].name
      }
    },
    getSortTable(table) {
      if (table.idHeader != null) {
        return this.tables.filter(item => item.id === table.idHeader[0])[0].sort
      }
    },
    getParent(tree, childId) {
      if (tree.children != null && tree.children.length > 0) {
        for (let i = 0; i < tree.children.length; i++) {
          if (tree.children[i].id === childId) {
            return tree
          }
          const req = this.getParent(tree.children[i], childId)
          if (req != null) {
            return req
          }
        }
        return null
      }
    },
    getIndexChildren(tree, childId) {
      if (tree.children != null) {
        for (let i = 0; i < tree.children.length; i++) {
          if (tree.children[i].id === childId) {
            return i
          }
        }
      }
    },
    isElementInTree: function (elementId, tree) {
      if (elementId != null) {
        if (elementId !== tree.id) {
          if (tree.children != null) {
            for (const child of tree.children) {
              const request = this.isElementInTree(elementId, child)
              if (request === true) {
                return true
              }
            }
          }
          return false
        }
        return true
      }
    },
    isTableIsTree(tableId, tree) {
      if (tree.idHeader?.length > 0) {
        if (tree.idHeader[0] === tableId) {
          return true
        }
      }
      if (tree.children != null) {
        for (let i = 0; i < tree.children.length; i++) {
          let ok = this.isTableIsTree(tableId, tree.children[i])
          if (ok) {
            return true
          }
        }
      }
      return false
    },

    messageAboutCopy() {
      this.messagesDisplay.color = '#ffae23'
      this.messagesDisplay.textColor = "#000000"
      this.messagesDisplay.show = true
      this.messagesDisplay.message = "В режиме добавления новой структуры запрещено выбирать другие года. Для разблокировки такой возможности отмените добавление или добавьте структуру в базу данных"
      this.messagesDisplay.timeout = 50000000
    },
    messageAboutEditNewStruct() {
      this.messagesDisplay.color = '#1c93ff'
      this.messagesDisplay.textColor = "white"
      this.messagesDisplay.show = true
      this.messagesDisplay.message = "Режим добавления новой структуры"
      this.messagesDisplay.timeout = 50000000
    },
    async checkHeaders(year) {
      await this.getYearsHeader()
      if (!this.yearsHeader.includes(year, 0)) {
        this.messagesDisplay.color = '#ffae23'
        this.messagesDisplay.textColor = "#000000"
        this.messagesDisplay.show = true
        this.messagesDisplay.message = 'Для ' + year + ' года отсутствуют таблицы. Перейдите в раздел "таблицы" для их добавления'
        return false
      } else {
        return true
      }
    },
    async addStruct() {
      if (await this.checkHeaders(this.addingYear)) {
        this.isAddedNewStruct = true //Ставим состояние добавления новой структуры
        this.yearsTree.unshift(this.addingYear) // Добавляем новый год
        this.selectedYear = this.yearsTree[0];  // Ставим его как выбранный
        await this.getTables(this.addingYear); // Получаем таблицы для выбранного года
        if (this.switchDuplicate) {
          await this.cloningTree(this.selectedYearForCloning)
        } else {
          this.tree = {id: -1, name: 'Структура экопаспортов', children: [], lvl: -1, idHeader: null, sort: 0}
        }
        this.windowForAddNewStruct = false
        this.switchDuplicate = false
        this.addingYear = null
        this.messageAboutEditNewStruct()
      }
    },
    async cloningTree(year) {
      await this.getTree(year)
      this.cleaningTree(this.tree); // Убираем все id из клонируемого дерева и генерируем новые. Также меняем idHeader на новые
    },
    cleaningTree: function (tree) {
      if (tree.children != null) {
        for (let i = 0; i < tree.children.length; i++) {
          const req = this.cleaningTree(tree.children[i])
          if (req === false) {
            tree.children.splice(i, 1) // Удаляем таблицу, если для нее не нашлось idHeader
            i-- // Меням смещаям счетчик назад после удаления, иначе пропускает один элемент
          }
        }
      }
      if (tree.idHeader != null) {
        let newTable = this.tables.filter(item => item.name === tree.name)[0]
        if (newTable == null) {
          return false
        }
        tree.idHeader = [newTable.id]
      }
      // Для корневой папки не генерируем id
      if (tree.id !== -1) {
        this.genUUID(this.listUUID)
      }
    },
    deleteId(tree) {
      if (tree.children != null) {
        for (let i = 0; i < tree.children.length; i++) {
          this.deleteId(tree.children[i])
        }
      }
      if (tree.id !== -1) {
        tree.id = null
      }
    },
    getYearsTree() {
      return new Promise(resolve => {
        axios
          .post("/api/yearsTree")
          .then(response => {
            this.yearsTree = response.data;
            resolve('resolved');
          })
          .catch(e => {
            alert(e.toString());
            resolve('resolved');
          });
      });
    },
    getYearsHeader() {
      return new Promise(resolve => {
        axios
          .post("/api/years")
          .then(response => {
            this.yearsHeader = response.data;
            resolve('resolved');
          })
          .catch(e => {
            alert(e.toString());
            resolve('resolved');
          });
      });
    },
    getTables(year) {
      this.tables = []
      return new Promise(resolve => {
        axios
          .post("/api/table", {year})
          .then(response => {
            this.tables = response.data;
            resolve('resolved');
          })
          .catch(e => {
            alert(e.toString());
            resolve('resolved');
          });
      })
    },
    getTree(year) {
      this.tree = {}
      this.openIds = []
      return new Promise(resolve => {
        axios
          .post("/api/getTree", {year: year})
          .then(response => {
            this.tree = {
              id: -1,
              name: "Структура экопаспортов",
              children: response.data,
              lvl: -1,
              idHeader: null,
              sort: 0
            };
            this.openIds.push(-1);
            resolve('resolved');
          })
          .catch(e => {
            alert("Шкипер, у нас проблемы: " + e.toString());
            resolve('resolved');
          });
      });

    },
    addTree() {
      this.deleteId(this.tree)
      axios
        .post("/api/addTree", {tree: this.tree, year: this.selectedYear})
        .then(() => {
          this.updateData()
          this.isAddedNewStruct =false
    })
        .catch(e => alert(e.toString()));
    },
    async updateData() {
      await this.getYearsTree();
      await this.getTree(this.selectedYear);
      await this.getTables(this.selectedYear);
    },
    genUUID(listUUID) {
      let UUID = "#" + Math.floor(Math.random() * 10000 + 1000)
      if (listUUID.includes(UUID)) {
        UUID = this.genUUID(listUUID)
      } else {
        listUUID.push(UUID)
      }
      return UUID
    },
  },
  async created() {
    await this.getYearsTree();
    this.selectedYear = this.maxYearTree
    await this.getTree(this.maxYearTree);
    await this.getTables(this.maxYearTree);
  }
  ,
}
</script>
