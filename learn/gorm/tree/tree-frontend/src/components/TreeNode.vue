<template>
  <li>
    <span @click="toggle" :class="{'expanded': expanded, 'collapsed': !expanded}">
      {{ node.name }}
    </span>
    <ul v-if="expanded && node.children && node.children.length">
      <TreeNode v-for="child in node.children" :key="child.name" :node="child"></TreeNode>
    </ul>
  </li>
</template>

<script>
export default {
  name: 'TreeNode',
  props: {
    node: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      expanded: false
    }
  },
  components: {
    TreeNode: () => import('./TreeNode.vue')
  },
  methods: {
    toggle() {
      this.expanded = !this.expanded
    }
  }
}
</script>

<style scoped>
li {
  position: relative;
}

span {
  user-select: none;
}
</style>
