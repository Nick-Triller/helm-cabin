<template>
  <div>
    <h1>Helm Cabin</h1>

    <div v-if="releasesList">

      <!--releases-filter :releases-list="releasesList"/-->

      <table>
        <tr>
          <th>Name</th>
          <th>Namespace</th>
          <th>Status</th>
          <th>Revision</th>
          <th>Chart</th>
        </tr>
        <tr v-for="release in releasesList.releases" v-bind:key="release.id">
          <td><router-link :to="`/releases/${release.name}/${release.version} `">{{ release.name }}</router-link></td>
          <td>{{ release.namespace }}</td>
          <td>{{ statusIdToNameMap[release.info.status.code] }}</td>
          <td>{{ release.version }}</td>
          <td>
            <a :href="release.chart.metadata.home">{{ release.chart.metadata.name }}</a>
          </td>
        </tr>
      </table>

    </div>
  </div>
</template>

<script>
import axios from 'axios'
import ReleasesFilterComponent from "./ReleasesFilterComponent";
import {statusIdToNameMap} from "../helper";

export default {
  name: 'ReleasesComponent',
  components: {
    releasesFilter: ReleasesFilterComponent,
  },
  data() {
    return {
      releasesList: null,
      statusIdToNameMap
    }
  },
  methods: {
    getReleases() {
      axios
        .get('/api/releases?status=0&status=1&status=2&status=3&status=4&status=5&status=6&status=7&status=8&offset=')
        .then(response => (this.releasesList = response.data))
    }
  },
  props: {},
  mounted () {
    setInterval(this.getReleases, 5000);
    this.getReleases();
  }
}
</script>

<style scoped>
   table {
     width: 100%;
   }
</style>
