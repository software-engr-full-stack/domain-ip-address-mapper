import React, { useEffect } from 'react';
import PropTypes from 'prop-types';

import {
  useMap,
  Marker,
  Popup,
  Tooltip
} from 'react-leaflet';

import { latLngBounds, Icon } from 'leaflet';

import markerIconPng from 'leaflet/dist/images/marker-icon.png';

import Contents from './Contents';
import IPList from './IPList';

function Markers({ markers }) {
  const map = useMap();

  const markerBounds = latLngBounds([]);
  markers.forEach(({ location }) => {
    const { latitude, longitude } = location;
    markerBounds.extend([latitude, longitude]);
  });

  useEffect(() => {
    map.fitBounds(markerBounds, { padding: [1, 1] });
  }, [map, markerBounds]);

  const validUrbanArea = (urbanArea) => urbanArea.iD > 0;

  return (
    <>
      {
        markers.map(({ location, ipList }) => {
          const {
            latitude, longitude, fullName, viewFullName,
            urbanArea
          } = location;

          return (
            <Marker
              key={fullName}
              riseOnHover
              position={[latitude, longitude]}
              icon={new Icon({ iconUrl: markerIconPng, iconSize: [25, 41], iconAnchor: [12, 41] })}
            >
              <Tooltip
                direction="top"
                offset={[0, -45]}
                opacity={1}
                permanent
                className="custom-map-marker-tooltip"
              >
                {viewFullName}
                <IPList ipList={ipList} />
                {
                  validUrbanArea(urbanArea) || <div>Not near an urban area</div>
                }
              </Tooltip>

              {
                validUrbanArea(urbanArea) && (
                  <Popup offset={[0, 380]}><Contents urbanArea={urbanArea} /></Popup>
                )
              }
            </Marker>
          );
        })
      }
    </>
  );
}

Markers.propTypes = {
  markers: PropTypes.arrayOf(PropTypes.any).isRequired
};

export default Markers;
