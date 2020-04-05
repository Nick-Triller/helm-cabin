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
        <td>{{ name }}</td>
      </tr>
      <tr>
        <td>Revision</td>
        <td>{{ version }}</td>
      </tr>
      <tr>
        <td>Status</td>
        <td>{{ release && release.Info && release.Info.Status ? release.Info.Status.StatusID : "-" }}</td>
      </tr>
      <tr>
        <td>Revisions</td>
        <td>
          <div v-for="version in revisionVersions" :key="version">
            <router-link :to="`/releases/${name}/${version} `">Version {{ version }}</router-link>
          </div>
        </td>
      </tr>
    </table>

    <div>
        <h2>Table of Contents</h2>
        <ol>
            <li><a href="#heading-rendered-manifest">Rendered manifest</a></li>
            <li><a href="#heading-chart-details">Chart details</a>
                <ol>
                    <li><a href="#heading-chart-templates">Chart templates</a></li>
                    <li><a href="#heading-chart-values">Chart values</a></li>
                    <li><a href="#heading-chart-files">Chart files</a></li>
                </ol>
            </li>
        </ol>
    </div>

    <div v-if="release && release.Manifest">
      <h2 id="heading-rendered-manifest">Rendered manifest</h2>
      <div>
        <prism language="yaml">{{release.Manifest}}</prism>
      </div>
    </div>

    <h1 id="heading-chart-details">Chart details</h1>

    <div v-if="release && release.Chart">
        <table>
            <tr>
                <th>Property</th>
                <th>Value</th>
            </tr>
            <tr v-for="prop in Object.keys(release.Chart).filter(x => release.Chart[x])" :key="prop">
                <td>{{prop}}</td>
                <td><pre>{{JSON.stringify(release.Chart[prop], null, 2)}}</pre></td>
            </tr>
        </table>
    </div>

    <div v-if="release && release.Templates">
      <h2 id="heading-chart-templates">Chart templates</h2>
      <div v-for="template in release.Templates" :key="template.name">
        <b>{{template.Name}}</b>
        <prism language="smarty">{{decodeBase64(template.Data)}}</prism>
      </div>
    </div>

    <div v-if="release && release.Values">
      <h2 id="heading-chart-values">Chart values</h2>
        <prism language="yaml">{{release.Values}}</prism>
    </div>

    <div v-if="release && release.Files">
      <h2 id="heading-chart-files">Chart files</h2>
      <div v-for="file in release.Files" :key="file.TypeURL">
        <b>{{file.TypeURL}}</b>
          <prism :language="file.TypeURL.endsWith('.md') ? 'markdown' : 'markup'">{{decodeBase64(file.Value)}}</prism>
      </div>
    </div>

  </div>
</template>

<script>
  import axios from "axios"
  import { Base64 } from "js-base64"
  import 'prismjs'
  import 'prismjs/themes/prism.css'
  import 'prismjs/components/prism-smarty'
  import 'prismjs/components/prism-yaml'
  import 'prismjs/components/prism-markdown'
  import Prism from 'vue-prism-component'

  export default {
    name: "ReleaseComponent",
    data() {
      return {
        revisionVersions: [],
        release: null,
        name: this.$route.params.name,
        version: this.$route.params.version
      }
    },
    watch: {
      '$route.params.version': function (version) {
        this.version = version;
        this.getRelease(this.name, version)
      }
    },
    methods: {
      decodeBase64(data) {
        return Base64.decode(data)
      },
      getRevisions(name) {
        axios
                .get(`/api/releases/${name}/revisions/`)
                .then(response => (this.revisionVersions = response.data));
      },
      getRelease(name, version) {
        axios
                .get(`/api/releases/${name}/versions/${version}`)
                .then(response => (this.release = response.data));
      }
    },
    mounted () {
      this.getRelease(this.name, this.version);
      this.getRevisions(this.name);
    },
    components: {
      Prism
    }
  }
</script>

<style scoped>
.collapsed {
  display: none;
}
</style>