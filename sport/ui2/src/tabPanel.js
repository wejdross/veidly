import { Box } from '@mui/material'
import React from 'react'

export function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`stt-${index}`}
      aria-labelledby={`tt${index}`}
      {...other}
    >
      {value === index && (
        <Box p={props.p}>
          {children}
        </Box>
      )}
    </div>
  );
}
