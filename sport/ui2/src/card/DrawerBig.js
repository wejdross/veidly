import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import {
  Drawer, Grid,
  IconButton,
  useTheme
} from '@mui/material';
import { makeStyles } from "@mui/styles";
import clsx from 'clsx';
import React from 'react';
import { MulwiColors } from '../mulwiColors';

const drawerWidth = 500

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  appBar: {
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginRight: props => props.width || drawerWidth,
  },
  title: {
    flexGrow: 1,
  },
  hide: {
    display: 'none',
  },
  drawer: {
    width: props => props.width || drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: props => props.width || drawerWidth,
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: 'flex-start',
  },
  content: {
    flexGrow: 1,
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    //marginRight: -drawerWidth,
  },
  contentShift: {
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginRight: props => props.width || drawerWidth,
  }
}))  

export function DrawerBig(props) {
    const theme = useTheme();
    const classes = useStyles(props);

    return (<React.Fragment>
        <Drawer
            transitionDuration={0}
            
            className={classes.drawer}
            anchor="right" 
            variant="persistent"
            classes={{
                paper: classes.drawerPaper,
            }}
            open={props.open} onClose={props.onClose}>

            <div style={{height:"100%", backgroundColor:MulwiColors.whiteBackground}}>
                <Grid container direction="row" 
                      style={{
                        marginTop: 70, 
                        paddingTop: props.padding || 17, 
                        paddingBottom: props.padding || 17
                      }}
                      justify="space-between"
                      alignItems="center">
                  <Grid xs={1} item>
                    <IconButton onClick={props.onClose}>
                        {theme.direction === 'rtl' ? <ChevronLeftIcon /> : <ChevronRightIcon />}
                    </IconButton>
                  </Grid>
                  <Grid xs={10} item style={{marginLeft: 10, marginRight: 20}}>
                    {props.navContent}
                  </Grid>
                </Grid>
              {props.content}
            </div>

        </Drawer>

        <main
            className={clsx(classes.content, {
                [classes.contentShift]: props.open,
            })}>
            {props.children}
        </main>
    </React.Fragment>)
}
