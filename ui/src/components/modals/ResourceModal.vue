<template>
  <dialog open class="modal">
    <div class="modal-header">
      <h3>Ressourcendetails</h3>
      <button class="close-btn" @click="$emit('close')">Schließen</button>
    </div>
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
/* Header mit Überschrift und Button */
.modal-header {
  display: flex; /* Flexbox für horizontale Anordnung */
  justify-content: space-between; /* Platz zwischen Überschrift und Button */
  align-items: center; /* Vertikale Zentrierung */
  margin-bottom: 20px; /* Abstand zum restlichen Inhalt */
}

h3 {
  margin: 0; /* Entferne Standardabstand der Überschrift */
  font-size: 20px;
  font-weight: bold;
}

/* Styling für Schließen-Button */
.close-btn {
  background: none;
  border: none;
  font-size: 14px;
  color: #007bff;
  cursor: pointer;
}

.close-btn:hover {
  color: red; /* Hover-Farbe für besseres Feedback */
}
</style>