import React from 'react';
import PropTypes from 'prop-types';

function Contents({ urbanArea }) {
  const { urbanAreaName, urbanAreaScores } = urbanArea;
  return (
    <div>
      <div>Urban area: {urbanAreaName}</div>
      <div>Scores...</div>
      <div>
        {
          urbanAreaScores.map(({ name, scoreOutOf10, color }) => (
            <div key={name} style={{ color }} className="custom-urban-area-score">
              {name}: {Math.round((scoreOutOf10 + Number.EPSILON) * 100) / 100}
            </div>
          ))
        }
      </div>
    </div>
  );
}

Contents.propTypes = {
  urbanArea: PropTypes.shape().isRequired
};

export default Contents;
