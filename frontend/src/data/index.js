import humps from 'humps';

import { persist, hydrate, persistClear } from './persist';
import { shapeDD, duplicate } from './objShapes';

const retrieve = ({ domain }) => (
  fetch(
    process.env.REACT_APP_API_END_POINT,
    {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'same-origin',
      headers: { 'Content-Type': 'application/json' },
      redirect: 'follow',
      referrerPolicy: 'no-referrer',
      body: JSON.stringify({ domain })
    }
  )
);

const buildIPList = (ip, fullName, acc) => {
  const hashValue = acc[fullName];
  if (!hashValue) return [ip];

  const { ipList } = hashValue;

  if (ipList) return [...ipList, ip];

  return [ip];
};

const uniqueMarkers = (json) => {
  const results = humps.camelizeKeys(json.results);

  const markerTable = results.ips.reduce((acc, ip) => {
    const { fullName } = ip.location;

    if (!fullName) {
      console.error('No full name found for IP...', ip); // eslint-disable-line no-console
      throw Error('...');
    }

    return {
      ...acc,
      [fullName]: {
        ...ip,
        ip,
        ipList: buildIPList(ip, fullName, acc)
      }
    };
  }, {});

  results.subDomains.forEach((sdm) => {
    sdm.ips.forEach((ip) => {
      const { location } = ip;
      if (!location) {
        console.error('No location found for IP...', ip); // eslint-disable-line no-console
        throw Error('...');
      }

      const { fullName } = location;
      if (!fullName) {
        console.error('No full name found for IP...', ip); // eslint-disable-line no-console
        throw Error('...');
      }

      markerTable[fullName] = {
        location,
        ipList: buildIPList(ip, fullName, markerTable)
      };
    });
  });

  return Object.keys(markerTable).map((key) => markerTable[key]);
};

const markersToLonLat = (markers) => markers.map((ip) => {
  const { longitude, latitude } = ip.location;
  return [longitude, latitude];
});

export {
  retrieve,
  uniqueMarkers,
  markersToLonLat,
  hydrate,
  persist,
  shapeDD,
  duplicate,
  persistClear
};
