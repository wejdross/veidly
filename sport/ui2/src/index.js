import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import './index.css';
import Main from './Main';
// import { MuiPickersUtilsProvider } from '@mui/lab/';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { MulwiColors } from './mulwiColors';

const theme = createTheme({
  typography: {
    fontFamily: [
      'Nunito Sans',
      'sans-serif'
    ].join(','),
    color: MulwiColors.blackText,
  }
});

export default
  ReactDOM.render(
    <React.StrictMode>
      <ThemeProvider theme={theme}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <BrowserRouter>
            <Main />
          </BrowserRouter>
        </LocalizationProvider>
        </ThemeProvider>
    </React.StrictMode>,
    document.getElementById('root')
  );
