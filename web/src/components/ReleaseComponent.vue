<template>
  <div>
    <h1>Release details</h1>
    <table>
      <tr>
        <th>Property</th>
        <th>Value</th>
      </tr>
      <tr>
        <td>Name</td>
        <td>{{ $route.params.name }}</td>
      </tr>
    </table>
  </div>
</template>

<script>
  export default {
    name: "ReleaseComponent"
  }
</script>

<style scoped>

</style>