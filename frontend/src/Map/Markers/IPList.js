import React, { useState } from 'react';
import PropTypes from 'prop-types';

import Link from '@mui/material/Link';
import Popover from '@mui/material/Popover';

function IPList({ ipList }) {
  const dm = ipList[0].domain;
  const addr = ipList[0].address;

  const [anchorEl, setAnchorEl] = useState(null);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'map-input-app-info-popover' : undefined;

  return (
    <>
      <div>{dm} {addr}</div>
      {
        ipList.length > 1 && (
          <Link href="#" onClick={handleClick}>More</Link> // eslint-disable-line jsx-a11y/anchor-is-valid
        )
      }

      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
      >
        <div className="custom-map-markers-ip-list">
          {
            ipList.map(({ domain, address }) => (
              <div key={[domain, address].join(':')}>{domain} {address}</div>
            ))
          }
        </div>
      </Popover>
    </>
  );
}

IPList.propTypes = {
  ipList: PropTypes.arrayOf(PropTypes.shape()).isRequired
};

export default IPList;
