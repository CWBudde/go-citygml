let map = null;
let popup = null;

export function initMap(containerId, data) {
  if (map) {
    updateMap(data);
    map.resize();
    return;
  }

  const geojson = data.geojson;
  if (!geojson || !geojson.features || geojson.features.length === 0) {
    return;
  }

  map = new maplibregl.Map({
    container: containerId,
    style: {
      version: 8,
      sources: {},
      layers: [
        {
          id: "background",
          type: "background",
          paint: { "background-color": getBackgroundColor() },
        },
      ],
    },
    center: [0, 0],
    zoom: 1,
  });

  popup = new maplibregl.Popup({ closeButton: true, closeOnClick: false });

  map.on("load", () => {
    addDataToMap(geojson, data);
    fitBounds(geojson);
  });
}

export function updateMap(data) {
  if (!map) return;
  const source = map.getSource("citygml");
  if (source) {
    source.setData(data.geojson);
  }
}

export function highlightFeature(objectId) {
  if (!map) return;
  map.setFilter("buildings-highlight", ["==", ["id"], objectId]);
}

export function destroyMap() {
  if (map) {
    map.remove();
    map = null;
    popup = null;
  }
}

function getBackgroundColor() {
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "#1e293b"
    : "#f0f0f0";
}

function addDataToMap(geojson, data) {
  map.addSource("citygml", {
    type: "geojson",
    data: geojson,
    promoteId: "id",
  });

  // Building fill
  map.addLayer({
    id: "buildings-fill",
    type: "fill",
    source: "citygml",
    filter: ["==", ["get", "type"], "Building"],
    paint: {
      "fill-color": buildingColorExpression(data),
      "fill-opacity": 0.6,
    },
  });

  // Building outline
  map.addLayer({
    id: "buildings-outline",
    type: "line",
    source: "citygml",
    filter: ["==", ["get", "type"], "Building"],
    paint: {
      "line-color": "#1e40af",
      "line-width": 1,
    },
  });

  // Building highlight
  map.addLayer({
    id: "buildings-highlight",
    type: "line",
    source: "citygml",
    filter: ["==", ["id"], ""],
    paint: {
      "line-color": "#f59e0b",
      "line-width": 3,
    },
  });

  // Terrain fill
  map.addLayer({
    id: "terrain-fill",
    type: "fill",
    source: "citygml",
    filter: ["==", ["get", "type"], "Terrain"],
    paint: {
      "fill-color": "#22c55e",
      "fill-opacity": 0.3,
    },
  });

  // Terrain outline
  map.addLayer({
    id: "terrain-outline",
    type: "line",
    source: "citygml",
    filter: ["==", ["get", "type"], "Terrain"],
    paint: {
      "line-color": "#16a34a",
      "line-width": 1,
    },
  });

  // Click handler
  map.on("click", "buildings-fill", (e) => {
    if (e.features.length === 0) return;
    const f = e.features[0];
    const props = f.properties;

    let html = `<strong>${f.id || "Building"}</strong><br>`;
    if (props.class) html += `Class: ${props.class}<br>`;
    if (props.function) html += `Function: ${props.function}<br>`;
    if (props.measuredHeight) html += `Height: ${props.measuredHeight}m (measured)<br>`;
    else if (props.derivedHeight) html += `Height: ${props.derivedHeight}m (derived)<br>`;
    if (props.lod) html += `LoD: ${props.lod}<br>`;

    popup.setLngLat(e.lngLat).setHTML(html).addTo(map);

    map.setFilter("buildings-highlight", ["==", ["id"], f.id || ""]);
  });

  map.on("click", "terrain-fill", (e) => {
    if (e.features.length === 0) return;
    const f = e.features[0];
    popup
      .setLngLat(e.lngLat)
      .setHTML(`<strong>${f.id || "Terrain"}</strong>`)
      .addTo(map);
  });

  // Cursor
  map.on("mouseenter", "buildings-fill", () => {
    map.getCanvas().style.cursor = "pointer";
  });
  map.on("mouseleave", "buildings-fill", () => {
    map.getCanvas().style.cursor = "";
  });
}

function buildingColorExpression(data) {
  // Color by height if available
  const heights = (data.objects || [])
    .filter((o) => o.type === "Building" && o.height > 0)
    .map((o) => o.height);

  if (heights.length === 0) {
    return "#3b82f6";
  }

  const minH = Math.min(...heights);
  const maxH = Math.max(...heights);

  if (minH === maxH) {
    return "#3b82f6";
  }

  return [
    "interpolate",
    ["linear"],
    ["coalesce", ["get", "measuredHeight"], ["get", "derivedHeight"], 0],
    minH,
    "#3b82f6",
    (minH + maxH) / 2,
    "#f59e0b",
    maxH,
    "#ef4444",
  ];
}

function fitBounds(geojson) {
  if (!map || !geojson.features.length) return;

  const bounds = new maplibregl.LngLatBounds();
  let hasCoords = false;

  for (const feature of geojson.features) {
    if (!feature.geometry) continue;
    visitCoords(feature.geometry.coordinates, (lng, lat) => {
      // Only add reasonable coordinates (WGS84 range)
      if (Math.abs(lng) <= 180 && Math.abs(lat) <= 90) {
        bounds.extend([lng, lat]);
        hasCoords = true;
      }
    });
  }

  if (hasCoords) {
    map.fitBounds(bounds, { padding: 40, maxZoom: 18 });
  }
}

function visitCoords(coords, fn) {
  if (!Array.isArray(coords)) return;
  if (typeof coords[0] === "number") {
    fn(coords[0], coords[1]);
    return;
  }
  for (const c of coords) {
    visitCoords(c, fn);
  }
}
