import React from 'react';
import PropTypes from 'prop-types';

import List from '@mui/material/List';
import ListSubheader from '@mui/material/ListSubheader';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';

import Tooltip from '@mui/material/Tooltip';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';

export default function RecentSearches({
  onClickDomainName,
  onClearSearches,
  listOfDomainData,
  currentDomainData
}) {
  return (
    <List>
      <ListSubheader component="div" color="primary">
        My recent searches
        <Tooltip title="Clear searches">
          <IconButton aria-label="clear-searches" onClick={onClearSearches}>
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </ListSubheader>

      {listOfDomainData.map((dd) => (
        <ListItem
          button
          key={dd.domain}
          selected={currentDomainData.domain === dd.domain}
          onClick={() => { onClickDomainName(dd); }}
        >
          <ListItemText primary={dd.domain} />
        </ListItem>
      ))}
    </List>
  );
}

RecentSearches.propTypes = {
  onClickDomainName: PropTypes.func.isRequired,
  onClearSearches: PropTypes.func.isRequired,
  currentDomainData: PropTypes.shape().isRequired,
  listOfDomainData: PropTypes.arrayOf(PropTypes.any).isRequired
};
