// https://stackoverflow.com/questions/26093545/how-to-validate-domain-name-using-regex/38578855#38578855
function isDomainValid(domain) {
  const re = /^(?:(?:(?:[a-zA-z-]+):\/{1,3})?(?:[a-zA-Z0-9])(?:[a-zA-Z0-9\-.]){1,61}(?:\.[a-zA-Z]{2,})+|\[(?:(?:(?:[a-fA-F0-9]){1,4})(?::(?:[a-fA-F0-9]){1,4}){7}|::1|::)\]|(?:(?:[0-9]{1,3})(?:\.[0-9]{1,3}){3}))(?::[0-9]{1,5})?$/;
  return domain.trim().match(re);
}

export {
  isDomainValid // eslint-disable-line import/prefer-default-export
};
