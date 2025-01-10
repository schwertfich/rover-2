<template>
  <dialog open class="modal">
    <h3>Ressourcendetails</h3><button @click="$emit('close')">Schließen</button>
    <resource-detail :resourceID="resourceID" />
  </dialog>
</template>

<script>
import ResourceDetail from "@/components/ResourceDetail.vue";

export default {
  name: "ResourceModal",
  components: {
    ResourceDetail,
  },
  props: {
    resourceID: {
      type: String,
      required: true,
    },
  },
};
</script>

<style scoped>
.modal {
  /* Modal bündig rechts und volle Bildschirmhöhe */
  position: fixed;
  top: 0;
  right: 0;
  width: 800px; /* Breite des Modals */
  height: 100vh; /* 100% der Bildschirmhöhe */
  background: white; /* Hintergrundfarbe des Modals */
  box-shadow: -4px 0 10px rgba(0, 0, 0, 0.2); /* Schattierung links */
  overflow-y: auto; /* Ermöglicht Scrollen bei größeren Inhalten */
  z-index: 100; /* Stellt sicher, dass das Modal darüber liegt */
  padding: 20px; /* Abstand zum Rand */
}

</style>