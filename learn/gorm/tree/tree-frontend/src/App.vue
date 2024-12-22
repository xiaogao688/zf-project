<template>
  <div id="app">
    <h1>树形展示</h1>
    <ul>
      <TreeNode v-for="node in treeData" :key="node.name" :node="node"></TreeNode>
    </ul>
  </div>
</template>

<script>
import axios from 'axios'
import TreeNode from './components/TreeNode.vue'

export default {
  name: 'App',
  components: {
    TreeNode
  },
  data() {
    return {
      treeData: []
    }
  },
  created() {
    this.fetchTreeData()
  },
  methods: {
    fetchTreeData() {
      axios.get('http://localhost:8080/api/tree')
        .then(response => {
          this.treeData = response.data
        })
        .catch(error => {
          console.error("Error fetching tree data:", error)
        })
    }
  }
}
</script>

<style>
ul {
  list-style-type: none;
  padding-left: 20px;
}

li {
  cursor: pointer;
  margin: 5px 0;
}

.expanded::before {
  content: "▼ ";
}

.collapsed::before {
  content: "▶ ";
}
</style>
