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
        <td>{{ release ? statusIdToName(release.info.status.code) : "-" }}</td>
      </tr>
    </table>

    <div>
        <h2>Table of Contents</h2>
        <ol>
            <li><a href="#heading-revisions">Revisions</a></li>
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

    <div>
      <h2 id="heading-revisions">Revisions</h2>
      <table>
        <tr>
          <th>Revision</th>
        </tr>
        <tr v-for="version in revisionVersions" :key="version">
          <td><router-link :to="`/releases/${name}/${version} `">Version {{ version }}</router-link></td>
        </tr>
      </table>
    </div>

    <div v-if="release && release.manifest">
      <h2 id="heading-rendered-manifest">Rendered manifest</h2>
      <div>
        <prism language="yaml">{{release.manifest}}</prism>
      </div>
    </div>

    <h1 id="heading-chart-details">Chart details</h1>

    <div v-if="release && release.chart && release.chart.metadata">
        <table>
            <tr>
                <th>Property</th>
                <th>Value</th>
            </tr>
            <tr v-for="prop in Object.keys(release.chart.metadata)" :key="prop">
                <td>{{prop}}</td>
                <td><pre>{{JSON.stringify(release.chart.metadata[prop], null, 2)}}</pre></td>
            </tr>
        </table>
    </div>

    <div v-if="release && release.chart">
      <h2 id="heading-chart-templates">Chart templates</h2>
      <div v-for="template in release.chart.templates" :key="template.name">
        <b>{{template.name}}</b>
        <prism language="smarty">{{decodeBase64(template.data)}}</prism>
      </div>
    </div>

    <div v-if="release && release.chart">
      <h2 id="heading-chart-values">Chart values</h2>
        <prism language="yaml">{{release.chart.values.raw}}</prism>
    </div>

    <div v-if="release && release.chart">
      <h2 id="heading-chart-files">Chart files</h2>
      <div v-for="file in release.chart.files" :key="file.type_url">
        <b>{{file.type_url}}</b>
          <prism :language="file.type_url.endsWith('.md') ? 'markdown' : 'markup'">{{decodeBase64(file.value)}}</prism>
      </div>
    </div>

  </div>
</template>

<script>
  import axios from "axios"
  import { Base64 } from "js-base64"
  import { statusIdToName } from  "../helper"
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
        version: this.$route.params.version,
      }
    },
    watch: {
      '$route.params.version': function (version) {
        this.version = version;
        this.getRelease(this.name, version)
      }
    },
    methods: {
      statusIdToName(id) {
        return statusIdToName(id);
      },
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