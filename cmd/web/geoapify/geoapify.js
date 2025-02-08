let currentMap = null;
let marker = null;

function initMap() {
  const mapContainer = document.querySelector("#map");

  const myAPIKey = "e1a483fb384e43d28d5a8eb3698df683";

  const markerIcon = L.icon({
    iconUrl: `https://api.geoapify.com/v1/icon/?type=awesome&color=%232ea2ff&size=large&scaleFactor=2&apiKey=${myAPIKey}`,
    iconSize: [38, 56],
    iconAnchor: [19, 51],
    popupAnchor: [0, -60],
  });

  if (mapContainer) {
    if (currentMap) {
      currentMap.remove();
      currentMap = null;
    }

    const latInput = document.querySelector("#geolocation-lat");
    const lonInput = document.querySelector("#geolocation-lon");

    let initialLat = 38.908838755401035; 
    let initialLon = -77.02346458179596; 

    if (latInput && lonInput && latInput.value && lonInput.value) {
      initialLat = parseFloat(latInput.value);
      initialLon = parseFloat(lonInput.value);
    }

    currentMap = L.map("map", { zoomControl: false }).setView(
      [initialLat, initialLon],
      12,
    );

    const isRetina = L.Browser.retina;
    const baseUrl = `https://maps.geoapify.com/v1/tile/osm-bright/{z}/{x}/{y}.png?apiKey=${myAPIKey}`;
    const retinaUrl = `https://maps.geoapify.com/v1/tile/osm-bright/{z}/{x}/{y}@2x.png?apiKey=${myAPIKey}`;

    L.tileLayer(isRetina ? retinaUrl : baseUrl, {
      attribution:
        'Powered by <a href="https://www.geoapify.com/" target="_blank">Geoapify</a> | <a href="https://openmaptiles.org/" rel="nofollow" target="_blank">© OpenMapTiles</a> <a href="https://www.openstreetmap.org/copyright" rel="nofollow" target="_blank">© OpenStreetMap</a> contributors',
      apiKey: myAPIKey,
      maxZoom: 20,
      id: "osm-bright",
    }).addTo(currentMap);

    L.control
      .zoom({
        position: "bottomright",
      })
      .addTo(currentMap);

    if (initialLat && initialLon) {
      marker = L.marker([initialLat, initialLon], {
        icon: markerIcon,
      }).addTo(currentMap);
      currentMap.panTo([initialLat, initialLon]);
    }
  }

  window.handleLocationSelect = function (location) {
    if (marker) {
      marker.remove();
    }

    marker = L.marker([location.lat, location.lon], {
      icon: markerIcon,
    }).addTo(currentMap);

    currentMap.panTo([location.lat, location.lon]);

    marker.bindPopup(location.formatted).openPopup();

    const resultsContainer = document.getElementById("geoloc-results");
    if (resultsContainer) {
      resultsContainer.style.display = "none";
    }

    const latInput = document.querySelector("#geolocation-lat");
    const lonInput = document.querySelector("#geolocation-lon");
    const locName = document.querySelector("#geolocation-name");

    if (latInput && lonInput) {
      latInput.value = location.lat;
      lonInput.value = location.lon;
      locName.value = location.formatted;
    } else {
      console.log("Lat and/or Lon input fields not found.");
    }
  };

  const geolocInput = document.getElementById("geoloc");
  if (geolocInput) {
    geolocInput.addEventListener("input", function () {
      const resultsContainer = document.getElementById("geoloc-results");
      if (resultsContainer) {
        resultsContainer.style.display = "block";
      }
    });
  }
}
