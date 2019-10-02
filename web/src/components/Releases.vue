<template>
  <div>
    <table>
      <tr>
        <th>Name</th>
        <th>Namespace</th>
        <th>Status</th>
        <th>Revision</th>
        <th>Chart</th>
      </tr>
      <tr v-for="release in releases" v-bind:key="release.id">
        <td><router-link :to="'/releases/' + release.Name">{{ release.Name }}</router-link></td>
        <td>{{ release.Namespace }}</td>
        <td>{{ release.Info.Status.StatusId }}</td>
        <td>{{ release.Version }}</td>
        <td>
          <a :href="release.Chart.Home">{{ release.Chart.Name }}</a>
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Releases',
  data() {
    return {
      releases: null,
    }
  },
  props: {},
  mounted () {
    const getReleases = () => {
      axios
              .get('/api/releases?status=0&status=1&status=2&status=3&status=4&status=5&status=6&status=7&status=8')
              .then(response => (this.releases = response.data))
    };
    setInterval(getReleases, 5000);
    getReleases();
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
 table {
   width: 100%;
 }
</style>
