const shapeDD = (domain, domainData) => ({ domain, data: domainData });

const shapePersist = (currentDomainData, listOfDomainData) => ({
  currentDomainData, listOfDomainData
});

const duplicate = (newDomain, listOfDomainData) => {
  const hash = listOfDomainData.reduce((memo, dd) => ({
    ...memo,
    [dd.domain]: dd
  }), {});

  return (hash[newDomain]);
};

export {
  shapeDD,
  shapePersist,
  duplicate
};
