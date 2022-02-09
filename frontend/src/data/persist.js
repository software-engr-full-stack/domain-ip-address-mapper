import { shapePersist } from './objShapes';

const persistKey = 'domain-ip-address-mapper';
const hydrate = () => {
  const saved = localStorage.getItem(persistKey);

  if (!saved) return {};

  return JSON.parse(saved);
};

const persist = (currentDomainData, listOfDomainData) => {
  if (!currentDomainData || !listOfDomainData) {
    console.error('... must pass defined params to be persisted'); // eslint-disable-line no-console
    throw Error('...');
  }

  localStorage.setItem(
    persistKey,
    JSON.stringify(shapePersist(currentDomainData, listOfDomainData))
  );
};

const persistClear = () => {
  localStorage.setItem(persistKey, null);
};

export {
  persist,
  hydrate,
  persistClear
};
