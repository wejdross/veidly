import MenuIcon from '@mui/icons-material/Menu';
import { AppBar, Button } from '@mui/material';
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import Toolbar from '@mui/material/Toolbar';
import { makeStyles } from "@mui/styles";
import clsx from 'clsx';
import React from 'react';
import { Link } from 'react-router-dom';
import { locale2 } from '../locale';
import { LangSelect } from './LangSelect';
import { NavList } from './NavList';

const useStyles = makeStyles((theme) => ({
  text: {
    padding: theme.spacing(2, 2, 0),
  },
  paper: {
    paddingBottom: 50,
  },
  list: {
    marginBottom: theme.spacing(2),
  },
  subheader: {
    backgroundColor: theme.palette.background.paper,
  },
  appBar: {
    top: 'on',
  },
  grow: {
    flexGrow: 1,
  },
  fabButton: {
    position: 'absolute',
    zIndex: 1,
    top: -30,
    left: 0,
    right: 0,
    margin: '0 auto',
  },
  white: {
    color: "white",
    fontSize: "0.75em"
  },
  itemlist: {
    width: 250,
  },
  itemfullList: {
    width: 'auto',
  },
}));

export default function NavLow(props) {
  const classes = useStyles();
  const [open, setOpen] = React.useState(false);

  const toggleDrawer = (open) => (event) => {
    if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
      return;
    }
    setOpen(open)
  };

  function logout() {
    props.main.logout()
  }

  const list = (anchor) => (
    <div
      className={clsx(classes.itemlist, {
        [classes.itemfullList]: anchor === 'top' || anchor === 'bottom',
      })}
      role="presentation"
      onClick={toggleDrawer(false)}
      onKeyDown={toggleDrawer(false)}
    >
      <NavList
        nots={props.nots}
        lang={props.lang} location={props.location}
        instructor={props.instructor} state={open} 
        user={props.user} logout={logout}/>
    </div>
  );
  return (
    <div>
      <React.Fragment>
        {(<React.Fragment>
          <Drawer anchor="left" open={open} onClose={toggleDrawer(false)}>
            {list("left")}
          </Drawer>
        </React.Fragment>)}
        <AppBar position="fixed" color="primary" className={classes.appBar} id="#navbar">
          <Toolbar style={{ backgroundColor: "#0e486e" }}>
            {(<IconButton onClick={toggleDrawer(true)} edge="start" color="inherit" aria-label="open drawer">
              <MenuIcon />
            </IconButton>)}
            <Link style={{ color: 'inherit', textDecoration: 'inherit' }} to="/" >

                      <img alt="V" src="/logo-light.png" height="25" style={{marginLeft: 10, marginRight: 10}}></img>
                  </Link>
            <div className={classes.grow} />
            <LangSelect lang={props.lang} setLang={props.setLang} />
            {(props.user && (
              <React.Fragment>
                <Button onClick={logout} color="inherit">{locale2.LOGOUT[props.lang]}</Button>
              </React.Fragment>
            )) || (
                <React.Fragment>
                  <Link to="/login">
                    <Button className={classes.white}>{locale2.LOGIN[props.lang]}</Button>
                  </Link>
                  <Link to="/register">
                    <Button className={classes.white}>{locale2.REGISTER[props.lang]}</Button>
                  </Link>
                </React.Fragment>
              )}
          </Toolbar>
        </AppBar>
      </React.Fragment>
    </div>
  );
}