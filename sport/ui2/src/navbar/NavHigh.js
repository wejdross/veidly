import React, { useEffect } from 'react';
import clsx from 'clsx';
import { useTheme } from '@mui/material/styles';
import Drawer from '@mui/material/Drawer';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import CssBaseline from '@mui/material/CssBaseline';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Link } from 'react-router-dom';
import { Button, Grid } from '@mui/material';
import { MulwiColors } from '../mulwiColors';
import { ConfigWarning } from './ConfigWarning';
import { NavList } from './NavList';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import { LangSelect } from './LangSelect';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { locale2 } from '../locale';
import { makeStyles } from "@mui/styles";

const drawerWidth = 220;

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginRight: 36,
  },
  hide: {
    display: 'none',
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
    whiteSpace: 'nowrap',
  },
  drawerOpen: {
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  drawerClose: {
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    overflowX: 'hidden',
    display: "none",
    width: theme.spacing(7) + 1,
    [theme.breakpoints.up('sm')]: {
      width: theme.spacing(8) + 1,
    },
  },
  toolbar: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
  },
  content: {
    flexGrow: 1
  },
  white: {
    color: "white"
  },
  activeLink: {
    backgroundColor: MulwiColors.black + " !important",
    color: "white"
  },
}));

export function NavHigh(props) {

  const classes = useStyles();
  const theme = useTheme();
  const [open, setOpen] = React.useState(false);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  function logout() {
    props.main.logout()
  }

  useEffect(() => {
    if (!props.user) {
      setOpen(false)
    }
  }, [props.user])

  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <div className={classes.root}>
      <CssBaseline />
      <AppBar
        id="#navbar"
        position="fixed"
        className={clsx(classes.appBar, {
          [classes.appBarShift]: open,
        })}>
        <Toolbar style={{ backgroundColor: MulwiColors.black }}>

          <Grid container direction="row" justifyContent={"space-between"} alignItems={"center"}>

            <Grid item>
              <Grid container direction={"row"} alignItems="center" justify="center" >

                {props.user && (
                  <Grid item>
                    <IconButton
                      color="inherit"
                      aria-label="open drawer"
                      onClick={handleDrawerOpen}
                      edge="start"
                      style={{ paddingLeft: "4px" }}
                      className={clsx(classes.menuButton, {
                        [classes.hide]: open,
                      })}>
                      <MenuIcon style={{ paddingRight: 0 }} />
                    </IconButton>
                  </Grid>
                )}

                <Grid item>
                  <Link style={{ color: 'inherit', textDecoration: 'inherit' }} to="/" >
                    <Button className={classes.white}  >
                      <img alt="V" src="/logo-light.png" height="35"></img>
                    </Button>
                  </Link>
                </Grid>


                <Grid item>
                  <Link style={{ color: 'inherit', textDecoration: 'inherit' }} to="/" >
                    <Button className={classes.white} style={{ fontWeight: 600, fontSize: 19, padding: "none" }}>
                      Veidly
                    </Button>
                  </Link>
                </Grid>

                <Grid item>
                  <Button aria-controls="about-company-menu" aria-haspopup="true"
                    onClick={handleClick} style={{ color: "white", lineHeight: "34px" }} >
                    {locale2.COMPANY[props.lang]}
                    {anchorEl ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                  </Button>
                  <Menu
                    id="about-company-menu"
                    anchorEl={anchorEl}
                    keepMounted
                    open={Boolean(anchorEl)}
                    onClose={handleClose}
                    style={{ marginTop: 10 }}>
                    <Link style={{ color: 'inherit', textDecoration: 'inherit' }} to="/benefits_instructor" >
                      <MenuItem onClick={handleClose}>{locale2.FOR_INSTRUCTOR[props.lang]}</MenuItem>
                    </Link>
                    <Link style={{ color: 'inherit', textDecoration: 'inherit' }} to="/benefits_user" >
                      <MenuItem onClick={handleClose}>{locale2.FOR_TRAINEE[props.lang]}</MenuItem>
                    </Link>
                    <Link style={{ color: 'inherit', textDecoration: 'inherit' }} to="/support" >
                      <MenuItem onClick={handleClose}>{locale2.DOCS[props.lang]}</MenuItem>
                    </Link>
                  </Menu>
                </Grid>
              </Grid>
            </Grid>
            {
              props.content &&
              <Grid item>
                {props.content}
              </Grid>}
            <Grid item>
              <Grid
                container
                direction={"row"}
                justifyContent={"flex-end"}
                alignContent={"flex-end"}
              >
                <LangSelect lang={props.lang} setLang={props.setLang} />
                {(props.user && (
                  <Button onClick={logout} color="inherit">{locale2.LOGOUT[props.lang]}</Button>
                )) || (
                    <>
                      <Link to="/login" style={{ textDecoration: "none" }}>
                        <Button className={classes.white}>{locale2.LOGIN[props.lang]}</Button>
                      </Link>
                      <Link to="/register" style={{ textDecoration: "none" }}>
                        <Button className={classes.white}>{locale2.REGISTER[props.lang]}</Button>
                      </Link>
                    </>
                  )}
              </Grid>
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      {props.user && (
        <Drawer
          variant="permanent"
          className={clsx(classes.drawer, {
            [classes.drawerOpen]: open,
            [classes.drawerClose]: !open,
          })}
          classes={{
            paper: clsx({
              [classes.drawerOpen]: open,
              [classes.drawerClose]: !open,
            }),
          }}>
          <div className={classes.toolbar}>
            <IconButton onClick={handleDrawerClose}>
              {theme.direction === 'rtl' ? <ChevronRightIcon /> : <ChevronLeftIcon />}
            </IconButton>
          </div>
          <Divider />

          <NavList location={props.location}
            lang={props.lang}
            nots={props.nots}
            instructor={props.instructor}
            user={props.user}
            logout={logout}
            state={open}
          />
        </Drawer>
      )}
      <main className={classes.content} style={{
        overflow: "hidden"
      }} >
        <div className={classes.toolbar} />
        <ConfigWarning lang={props.lang} instructor={props.instructor} />
        {props.children}
      </main>
    </div>
  );
}
