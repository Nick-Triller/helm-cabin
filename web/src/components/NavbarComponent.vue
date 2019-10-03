<template>
    <div id="navbar">
        <div id="navitems">
            <router-link exact exact-active-class="active" to="/">Releases</router-link>
        </div>
        <div id="links">
            <a href="https://github.com/Nick-Triller/helmcabin">GitHub</a>
        </div>
        <div class="clearfix"></div>
    </div>
</template>

<script>
  export default {
    name: "NavbarComponent"
  }
</script>

<style scoped>
    #navbar {
        margin-bottom: 2em;
    }

    #navitems {
        float: left;
    }

    #links {
        float: right;
    }

    #navbar a {
        padding: 1em;
    }

    .active {
        background-color: lightgray;
    }
</style>
