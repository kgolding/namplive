<template>
  <span @click="handleClick()">
    {{ value }}
  </span>
</template>

<script>
export default {
  props: {
    value: Boolean,
  },
  methods: {
    handleClick() {
      this.$emit("input", !this.value);
    },
  },
};
</script>