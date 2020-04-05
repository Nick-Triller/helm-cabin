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
        <tr v-for="release in releasesList" v-bind:key="release.Name">
          <td><router-link :to="`/releases/${release.Name}/${release.Version} `">{{ release.Name }}</router-link></td>
          <td>{{ release.Namespace }}</td>
          <td>{{release.Info.Status.StatusID}}</td>
          <td>{{ release.Version }}</td>
          <td>
            <a v-if="release.Chart.Home" :href="release.Chart.Home">{{ release.Chart.Name }}</a>
            <span v-if="!release.Chart.Home">{{ release.Chart.Name }}</span>
          </td>
        </tr>
      </table>

    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'ReleasesComponent',
  components: {

  },
  data() {
    return {
      releasesList: null
    }
  },
  methods: {
    getReleases() {
      axios
        .get('/api/releases')
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
