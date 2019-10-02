<template>
  <div id="app">
    <NavbarComponent></NavbarComponent>
    <router-view></router-view>
  </div>
</template>

<script>
import NavbarComponent from "./components/NavbarComponent";
export default {
  name: 'app',
  components: {NavbarComponent},
}
</script>

<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    color: #2c3e50;
    padding: 1em;
    max-width: 1200px;
    margin-left: auto;
    margin-right: auto;
  }

  .clearfix {
    clear: both;
  }

  table {
    border-collapse: collapse;
  }

  table, th, td {
    border: 1px solid black;
  }

  th, td {
    padding: 10px;
    text-align: left;
  }

  th {
    text-align: center;
  }
</style>
