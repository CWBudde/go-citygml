import { initMap, updateMap, highlightFeature, destroyMap } from "./modules/map.js";
import { initScene, updateScene, highlightObject as highlightSceneObject, destroyScene, setWireframe } from "./modules/scene.js";
import { updateSidebar, setSelectionCallback } from "./modules/sidebar.js";

let wasmReady = false;
let currentData = null;
let activeView = "map";

// --- WASM initialization ---

async function initWasm() {
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch("citygml.wasm"),
    go.importObject
  );
  go.run(result.instance);
  wasmReady = true;
}

initWasm().catch((err) => {
  console.error("WASM init failed:", err);
  showError("Failed to initialize WebAssembly: " + err.message);
});

// --- DOM refs ---

const dropZoneView = document.getElementById("drop-zone-view");
const vizView = document.getElementById("viz-view");
const dropZone = document.getElementById("drop-zone");
const fileInput = document.getElementById("file-input");
const loading = document.getElementById("loading");
const errorDiv = document.getElementById("error");
const errorMsg = document.getElementById("error-message");
const fileName = document.getElementById("file-name");
const fileMeta = document.getElementById("file-meta");
const newFileBtn = document.getElementById("new-file-btn");

// --- Drop zone ---

dropZone.addEventListener("click", () => fileInput.click());

dropZone.addEventListener("dragover", (e) => {
  e.preventDefault();
  dropZone.classList.add("drag-over");
});

dropZone.addEventListener("dragleave", () => {
  dropZone.classList.remove("drag-over");
});

dropZone.addEventListener("drop", (e) => {
  e.preventDefault();
  dropZone.classList.remove("drag-over");
  const file = e.dataTransfer.files[0];
  if (file) processFile(file);
});

fileInput.addEventListener("change", () => {
  const file = fileInput.files[0];
  if (file) processFile(file);
  fileInput.value = "";
});

// --- View switching ---

newFileBtn.addEventListener("click", () => {
  vizView.classList.remove("active");
  dropZoneView.classList.add("active");
  destroyMap();
  destroyScene();
  currentData = null;
});

// View toggle buttons
document.querySelectorAll(".viz-toggle-btn").forEach((btn) => {
  btn.addEventListener("click", () => {
    const view = btn.dataset.view;
    if (view === activeView) return;

    document.querySelectorAll(".viz-toggle-btn").forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");

    document.querySelectorAll(".viz-canvas").forEach((c) => c.classList.remove("active"));
    document.getElementById(view === "map" ? "map-container" : "scene-container").classList.add("active");

    activeView = view;

    if (view === "map" && currentData) {
      initMap("map", currentData);
    } else if (view === "scene" && currentData) {
      initScene("scene-canvas", currentData);
    }
  });
});

// Sidebar tabs
document.querySelectorAll(".sidebar-tab").forEach((tab) => {
  tab.addEventListener("click", () => {
    document.querySelectorAll(".sidebar-tab").forEach((t) => t.classList.remove("active"));
    tab.classList.add("active");

    document.querySelectorAll(".sidebar-content").forEach((c) => c.classList.remove("active"));
    document.getElementById("tab-" + tab.dataset.tab).classList.add("active");
  });
});

// Wireframe toggle
document.getElementById("wireframe-toggle").addEventListener("change", (e) => {
  setWireframe(e.target.checked);
});

// --- Sidebar selection callback ---

setSelectionCallback((objectId) => {
  if (activeView === "map") {
    highlightFeature(objectId);
  } else {
    highlightSceneObject(objectId);
  }
});

// --- File processing ---

async function processFile(file) {
  if (!wasmReady) {
    showError("WebAssembly is still loading. Please wait and try again.");
    return;
  }

  loading.classList.remove("hidden");
  errorDiv.classList.add("hidden");
  dropZone.classList.add("hidden");

  try {
    const arrayBuffer = await file.arrayBuffer();
    const uint8Array = new Uint8Array(arrayBuffer);

    // Give the UI a frame to show the spinner
    await new Promise((r) => requestAnimationFrame(r));

    const result = parseCityGML(uint8Array);

    if (!result.success) {
      showError(result.error || "Unknown parse error");
      return;
    }

    currentData = result;

    // Update header
    fileName.textContent = file.name;
    const parts = [];
    if (result.meta.version) parts.push("v" + result.meta.version);
    if (result.meta.buildingCount) parts.push(result.meta.buildingCount + " buildings");
    if (result.meta.terrainCount) parts.push(result.meta.terrainCount + " terrains");
    fileMeta.textContent = parts.join(" | ");

    // Switch to viz view
    dropZoneView.classList.remove("active");
    vizView.classList.add("active");

    // Initialize default view
    activeView = "map";
    document.querySelectorAll(".viz-toggle-btn").forEach((b) => b.classList.remove("active"));
    document.querySelector('[data-view="map"]').classList.add("active");
    document.querySelectorAll(".viz-canvas").forEach((c) => c.classList.remove("active"));
    document.getElementById("map-container").classList.add("active");

    updateSidebar(result);
    initMap("map", result);
  } catch (err) {
    showError("Error processing file: " + err.message);
    console.error(err);
  } finally {
    loading.classList.add("hidden");
  }
}

function showError(msg) {
  loading.classList.add("hidden");
  dropZone.classList.remove("hidden");
  errorDiv.classList.remove("hidden");
  errorMsg.textContent = msg;
}
