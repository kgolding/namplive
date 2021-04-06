<template>
  <div class="container-fluid">
    <Results/>
  </div>
</template>

<script>
import Results from "./components/Results.vue";
import "../node_modules/bootstrap/dist/css/bootstrap.css"
// import "./assets/css/main.css";

export default {
  name: "app",
  components: {
    Results
  }
};
</script>
