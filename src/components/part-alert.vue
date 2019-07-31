<template>
  <div id="message" class="alert alert-info alert-dismissible" role="alert">
    <button
      type="button"
      class="close"
      data-dismiss="alert"
      aria-label="閉じる"
    >
      <span aria-hidden="true">×</span>
    </button>
    <strong>Oh snap!</strong> Change a few things up and try submitting again.
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
@Component
export default class Alert extends Vue {}
</script>
<style scoped>
#message {
  display: none;
}
</style>
