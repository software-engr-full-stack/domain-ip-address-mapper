import React, { useState, useEffect } from 'react';

import AlertTitle from '@mui/material/AlertTitle';
import Alert from '@mui/material/Alert';

import * as data from './data';

import Map from './Map/Map';
import Input from './Map/Input';
import RecentSearches from './Reports/RecentSearches';
import UI from './UI/UI';

import { isDomainValid } from './lib';

function App() {
  const [apiError, setAPIError] = useState();

  const [currentDomainData, setCurrentDomainData] = useState({});

  const [listOfDomainData, setListOfDomainData] = useState([]);

  const [inputError, setInputError] = useState(false);

  const defaultInputLabel = 'Enter domain name';
  const [inputLabel, setInputLabel] = useState(defaultInputLabel);

  const [isFetching, setIsFetching] = useState(false);

  // Probably get rid of this in the future. Back up if no inputs.
  const defaultDomain = 'hired.com';
  const [domain, setDomain] = useState(defaultDomain);

  const handleSubmit = (domainFromInput) => {
    if (isDomainValid(domainFromInput)) {
      setIsFetching(true);
      data.retrieve({ domain: domainFromInput }).then((resp) => resp.json()).then((json) => {
        const jsonErr = json.error;
        if (jsonErr) {
          setAPIError(jsonErr);
          setIsFetching(false);
          return;
        }
        const domainData = data.uniqueMarkers(json);

        const cdd = data.shapeDD(domainFromInput, domainData);
        setCurrentDomainData(cdd);
        setDomain(domainFromInput);

        if (!data.duplicate(domainFromInput, listOfDomainData)) {
          const newList = [...listOfDomainData, cdd];
          setListOfDomainData(newList);
          data.persist(cdd, newList);
        } else {
          data.persist(cdd, listOfDomainData);
        }

        setIsFetching(false);
      });
      return;
    }

    setInputError(true);
    setInputLabel('Invalid domain');
  };

  const handleClickDomainName = (domainData) => {
    setCurrentDomainData(domainData);
    setDomain(domainData.domain);
    data.persist(domainData, listOfDomainData);
  };

  const handleClearSearches = () => {
    setListOfDomainData([]);
    data.persistClear();
  };

  useEffect(() => {
    const persistedData = data.hydrate();
    let pcd;
    if (persistedData) {
      pcd = persistedData.currentDomain;
      if (pcd) setDomain(pcd);
      const persListOfMarkers = persistedData.listOfDomainData;
      if (persListOfMarkers && persListOfMarkers.length > 0) {
        const cdd = persistedData.currentDomainData;
        if (currentDomainData) {
          setCurrentDomainData(cdd);
          setListOfDomainData(persListOfMarkers);
          // Apparently, the states aren't set inside useEffect
          if (!pcd) setDomain(currentDomainData.domain);
          return;
        }
      }
    }

    const actualDm = pcd || domain;
    if (!actualDm) return;
    setIsFetching(true);
    data.retrieve({ domain: actualDm }).then((resp) => resp.json()).then((json) => {
      const jsonErr = json.error;
      if (jsonErr) {
        setAPIError(jsonErr);
        setIsFetching(false);
        return;
      }
      setCurrentDomainData(data.shapeDD(actualDm, data.uniqueMarkers(json)));
      setIsFetching(false);
    });
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const markers = currentDomainData.data;

  if (!markers || markers.length < 1) {
    if (apiError) {
      return (
        // TODO: Make front end error reporting more robust.
        //   Probably add a on close (onClose) handler.
        //   See src/Map/Input.js Alert.
        <UI>
          <Alert severity="error">
            <AlertTitle>Error</AlertTitle>
            {apiError}
          </Alert>
        </UI>
      );
    }
    return null;
  }

  return (
    <UI
      Input={(
        <Input
          error={apiError}
          setAPIError={setAPIError}
          onSubmit={handleSubmit}
          inputError={inputError}
          setInputError={setInputError}
          defaultLabel={defaultInputLabel}
          setLabel={setInputLabel}
          label={inputLabel}
          isFetching={isFetching}
        />
      )}
      LeftSideBar={(
        <RecentSearches
          onClickDomainName={handleClickDomainName}
          onClearSearches={handleClearSearches}
          listOfDomainData={listOfDomainData}
          currentDomainData={currentDomainData}
        />
      )}
    >
      {
        isFetching ? (
          <div
            style={{
              opacity: 0.25,
              background: '#000000',
              height: '92vh'
            }}
          />
        ) : (
          <Map markers={markers} />
        )
      }
    </UI>
  );
}

export default App;
