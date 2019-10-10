<template>
  <div>
    <fieldset>
      <legend>Filter</legend>

      <fieldset>
        <legend>Status</legend>
        <template v-for="statusId in Object.keys(statusIdToNameMap)">
          <input :key="`${statusId}-input`" :id="`${statusIdToNameMap[statusId]}-status`" type="checkbox" checked/>
          <label :key="`${statusId}-label`" :for="`${statusIdToNameMap[statusId]}-status`">
            {{statusIdToNameMap[statusId]}}
          </label>
          <br :key="`${statusId}-br`"/>
        </template>
        <button>Enable all</button>
        <button>Disable all</button>
      </fieldset>

    </fieldset>
    <button :onclick="handleRefresh">Refresh</button>
  </div>
</template>

<script>
  import {statusIdToNameMap} from '../helper'

  export default {
    name: "ReleasesFilterComponent",
    props: {
      releasesList: Object,
      onRefresh: Function,
    },
    methods: {
      handleRefresh(event) {
        this.onRefresh(event);
      }
    },
    data() {
      return {
        statusIdToNameMap
      }
    }
  }
</script>

<style scoped>
button {
  margin: 20px;
}
</style>