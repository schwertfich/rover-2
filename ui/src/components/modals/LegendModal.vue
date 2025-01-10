<template>
  <dialog open class="modal">
    <fieldset>
      <legend>Legende</legend>
      <b>Anweisungen</b>
      <hr />
      <p>
        Klicken oder bewegen Sie die Maus über einen Knoten, um die
        Verbindungen hervorzuheben. Zum Abwählen klicken Sie auf den
        Hintergrund.
      </p>
      <p>
        Alle Ressourcen, von denen der Knoten abhängt, sind durch eine
        durchgezogene Linie dargestellt. Ressourcen, die vom Knoten abhängen,
        sind durch eine gestrichelte Linie dargestellt.
      </p>
      <hr />
      <b>Ressourcentypen</b>
      <div class="node create">Resource - Create</div>
      <div class="node delete">Resource - Delete</div>
      <div class="node replace">Resource - Replace</div>
      <div class="node update">Resource - Update</div>
      <div class="node no-op">Resource - No Operation</div>
      <hr />
      <b>Sonstige Elemente</b>
      <div class="node variable">Variable</div>
      <div class="node output">Output</div>
      <div class="node data">Data</div>
      <div class="node module">Module</div>
      <div class="node locals">Local</div>
    </fieldset>
    <button @click="$emit('close')">Schließen</button>
  </dialog>
</template>

<script>
export default {
  name: "LegendModal",
};
</script>

<style scoped>
.modal {
  position: absolute;
  top: 10%;
  left: 50%;
  transform: translateX(-50%);
  background: white;
  border: 1px solid #ccc;
  padding: 15px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.2);
  z-index: 100;
}

.node {
  display: inline-block;
  margin: 5px;
  padding: 5px;
  border: 1px solid #ccc;
  font-size: 0.9em;
}

.node.create {
  border-color: green;
  color: green;
}
.node.delete {
  border-color: red;
  color: red;
}
.node.replace {
  border-color: orange;
  color: orange;
}
.node.update {
  border-color: blue;
  color: blue;
}
.node.no-op {
  border-color: grey;
  color: grey;
}
</style>