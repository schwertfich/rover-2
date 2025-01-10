<template>
  <div id="app">
    <!-- Navigationsleiste -->
    <main-nav class="fixed-navbar" @saveGraph="saveGraph">
    </main-nav>

    <!-- Graph -->
    <div class="graph-container">
      <graph
          ref="filegraph"
          :displayGraph="displayGraph"
          v-on:getNode="openResourceDetails"
      />
    </div>

    <!-- Modals -->

    <resource-modal
        v-if="resourceID"
        :resourceID="resourceID"
        @close="closeResourceModal"
    />
  </div>
</template>

<script>
import MainNav from "@/components/MainNav.vue";
import Graph from "@/components/Graph/Graph.vue";
import ResourceModal from "@/components/modals/ResourceModal.vue";

export default {
  name: "App",
  metaInfo: {
    title: "Terraform Visualization",
  },
  components: {
    MainNav,
    Graph,
    ResourceModal,
  },
  data() {
    return {
      displayGraph: true, // Steuerung des Graphen
      resourceID: "", // ID der aktuell ausgewählten Ressource
    };
  },
  methods: {
    saveGraph() {
      // Speichere den aktuellen Graph
      this.$refs.filegraph.saveGraph();
    },
    openResourceDetails(resourceID) {
      // Öffne detaillierte Ansicht für eine Ressource
      this.resourceID = resourceID;
    },
    closeResourceModal() {
      // Schließe das Ressourcen-Detail-Modal
      this.resourceID = "";
    },
  },
};
</script>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  margin: 0;
  padding: 0;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Positionierung des Graphen */
.graph-container {
  flex: 1;
  height: calc(100% - 60px); /* Abzug der Höhe der MainNav */
  width: 100%;
  position: relative;
}
.fixed-navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;

  color: white;
  padding: 10px 20px;
  z-index: 10; /* Stellt sicher, dass die Navbar über allem liegt */
  display: flex;
  justify-content: space-between;
  align-items: center;
}

</style>