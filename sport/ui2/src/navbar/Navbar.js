import React, { useState } from 'react';
import { NavHigh } from "./NavHigh"
import { withRouter } from 'react-router-dom';
import { useMediaQuery, useTheme } from "@mui/material";
import NavLow from "./NavLow";

function Navbar(props) {
  const theme = useTheme()
  const matches = useMediaQuery(theme.breakpoints.down('sm'))

  const [location, setLocation]
    = useState(window.location.pathname.toLowerCase())

  props.history.listen((location, action) => {
    setLocation(location.pathname.toLowerCase())
  });

  return (
    (matches && (
      <React.Fragment>
        <NavLow location={location} 
                main={props.main}
                nots={props.nots}
                lang={props.lang} setLang={props.setLang}
                history={props.history} 
                user={props.user} 
                instructor={props.instructor}/>
                <main style={{
          overflow:"hidden"
        }}>
        {props.children}
          </main>
      </React.Fragment>
    )) || (
      <NavHigh  location={location} 
                content={props.content}
                history={props.history} 
                user={props.user} 
                lang={props.lang} setLang={props.setLang}
                main={props.main}
                nots={props.nots}
                instructor={props.instructor}>
        {props.children}
      </NavHigh>
    )
  )
}

export default withRouter(Navbar)