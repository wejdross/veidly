import React from "react";
import './Spinner.css';

import useMediaQuery from '@mui/material/useMediaQuery';
import { Typography } from "@mui/material";
//import { useTheme } from "@material-ui/core/styles";


export default function Spinner() {
  const smallScreen = useMediaQuery('(min-width:600px)');

  return (
      <div className="lds-container" id="overlay">
        <Typography variant={smallScreen ? "h3" : "h5"}>
          Veidly
        </Typography>
        <div className="lds-dual-ring"></div>
      </div>
  );
}
