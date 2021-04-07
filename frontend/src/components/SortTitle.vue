<template>
  <div @click="handleClick()">
    <slot></slot>
    <fa v-if="icon" :icon="icon" class="ml-2"/>
  </div>
</template>

<script>
export default {
  props: {
    field: String,
    value: Object,
  },
  computed: {
    icon() {
      if (this.field === this.value.field) {
        return this.value.asc ? "sort-up" : "sort-down"
      }
      return undefined
    }
  },
  methods: {
    handleClick() {
      let asc = this.value.asc;
      if (this.value.field == this.field) {
        asc = !asc;
      }
      this.$emit("input", { field: this.field, asc: asc });
    },
  },
};
</script>