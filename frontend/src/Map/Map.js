import React from 'react';
import PropTypes from 'prop-types';

import 'leaflet/dist/leaflet.css';

import {
  MapContainer,
  TileLayer
} from 'react-leaflet';

import './Map.css';

import Markers from './Markers/Markers';

function Map({ markers }) {
  if (markers.length < 1) return null;

  return (
    <MapContainer scrollWheelZoom>
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {
        markers.length > 0 && <Markers markers={markers} />
      }
    </MapContainer>
  );
}

Map.propTypes = {
  markers: PropTypes.arrayOf(PropTypes.any).isRequired
};
export default Map;
