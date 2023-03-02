<template>
  <div style="padding-top: 0; margin-top: 0px;">
    <v-container style="height: max-content;padding-bottom: 15px;">
      <v-row>
        <v-col mb-6 no-gutters>
          <v-autocomplete v-model="selectedRegionId" :items="regions" item-value="id" item-text="name" label="МО">
          </v-autocomplete>
          <v-switch
            v-model="GeneralInformation.isTown"
            label="Населенный пункт"
            color="#008B8B"
            hide-details
            style="padding-top: 5px;padding-bottom: 25px"
          ></v-switch>
          <v-text-field type="number" class="v-text-field" v-model.number="GeneralInformation.lng"
                        label="Долгота"></v-text-field>
          <v-text-field type="number" v-model.number="GeneralInformation.lat" label="Широта"></v-text-field>
          <v-text-field v-model="EnvironmentalAssessment.createDate" label="Дата образования"></v-text-field>
          <v-text-field v-model="GeneralInformation.oktmo" label="ОКТМО"></v-text-field>
        </v-col>
        <v-col mb-6 no-gutters>
          <v-select v-model="selectedYear" :items="years" label="Год"/>
          <v-text-field v-model="EnvironmentalAssessment.center" label="Административный центр"></v-text-field>
          <v-text-field v-model="EnvironmentalAssessment.population" label="Численность населения"></v-text-field>
          <v-text-field v-model="EnvironmentalAssessment.area" label="Площадь кв.км"></v-text-field>
          <v-text-field v-model="EnvironmentalAssessment.grossEmissions" label="Валовые выбросы, тонн"></v-text-field>
          <v-text-field v-model="EnvironmentalAssessment.formedWaste"
                        label="Масса образованных отходов, тонн"></v-text-field>
          <v-text-field v-model="EnvironmentalAssessment.withdrawnWater"
                        label="Объем забранной воды, млн куб.м"></v-text-field>
          <v-text-field v-model="EnvironmentalAssessment.dischargeVolume"
                        label="Объем водоотведения млн куб.м"></v-text-field>
          <v-row class="float-right">
            <v-btn style="margin-inline: 5px" @click="cancel">Отменить</v-btn>
            <v-btn style="margin-inline: 5px" @click="save">Сохранить</v-btn>
          </v-row>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import axios from "axios";
import {bus} from "@/main";

export default {
  data() {
    return {
      selectedYear: null,
      selectedRegionId: null,
      regions: [],
      years: [],
      GeneralInformation: {
        isTown: false,
        lng: '',
        lat: '',
      },
      EnvironmentalAssessment: {
        area: '',
        center: '',
        withdrawnWater: '',
        dischargeVolume: '',
        formedWaste: '',
        grossEmissions: '',
        population: '',
        createDate: '',
      },
      GeneralInformation2: {
        isTown: false,
        lng: '',
        lat: '',
      },
      EnvironmentalAssessment2: {
        center: '',
        withdrawnWater: '',
        dischargeVolume: '',
        formedWaste: '',
        grossEmissions: '',
        population: '',
        createDate: '',
      },
    }
  },
  watch: {
    selectedRegionId: function () {
      this.getYears()
      this.getGeneralInformation()
      if (this.EnvironmentalAssessment !== {}) {
        this.getEnvironmentalAssessment()
      }
    },
    years: function () {
      if (this.selectedYear == null) {
        this.selectedYear = this.years[0]
      }
    },
    selectedYear: function () {
      this.getEnvironmentalAssessment()
    }
  },
  methods: {
    save() {
      this.GeneralInformation.oktmo = this.GeneralInformation.oktmo.trim()
      axios
        .post("/api/saveRegionInfo", {
          regionId: this.selectedRegionId,
          year: this.selectedYear,
          EnvironmentalAssessment: this.EnvironmentalAssessment,
          GeneralInformation: this.GeneralInformation
        })
        .then(() => {
          bus.$emit('message', {
            message: "Информация обновлена",
            color: "green",
          })
          this.getEnvironmentalAssessment()
          this.getGeneralInformation()
        })
        .catch(e => alert(e.toString()));
    },
    cancel() {
      this.EnvironmentalAssessment = Object.assign({}, this.EnvironmentalAssessment2)
      this.GeneralInformation = Object.assign({}, this.GeneralInformation2)
    },
    getRegions() {
      this.regions = []
      axios
        .post("/api/regions")
        .then(response => (this.regions = response.data))
        .catch(e => alert(e.toString()));
    },
    getYears() {
      this.years = []
      axios
        .post("/api/yearsRegions", {regionId: this.selectedRegionId})
        .then(response => (this.years = response.data))
        .catch(e => alert(e.toString()));
    },
    getGeneralInformation() {
      this.GeneralInformation = {}
      axios
        .post("/api/getGeneralInformation", {regionId: this.selectedRegionId})
        .then(response => {
          this.GeneralInformation = response.data
          this.GeneralInformation2 = Object.assign({}, response.data)
        })
        .catch(e => alert(e.toString()));
    },
    getEnvironmentalAssessment() {
      this.EnvironmentalAssessment = {}
      axios
        .post("/api/getEnvironmentalAssessment", {regionId: this.selectedRegionId, year: this.selectedYear})
        .then(response => {
            this.EnvironmentalAssessment = response.data
            this.EnvironmentalAssessment2 = Object.assign({}, response.data)
          }
        )
        .catch(e => alert(e.toString()));
    }
  },
  created() {
    this.getRegions()
  }
}
</script>
<style scoped>
.v-text-field >>> input[type="number"]::-webkit-outer-spin-button,
.v-text-field >>> input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
}

.v-text-field >>> input[type='number'],
input[type="number"]:hover,
input[type="number"]:focus {
  appearance: none;
  -moz-appearance: textfield;
}
</style>
