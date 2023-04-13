import { Drawer, Grid, IconButton } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import React from 'react';
import { MulwiColors } from '../mulwiColors';
import { makeStyles } from "@mui/styles";
import useMediaQuery from '@mui/material/useMediaQuery';

export default function DrawerSmall(props) {

  const theme = useTheme()
  const sm = useMediaQuery(theme.breakpoints.down('xs'))

  const useStyles = makeStyles(theme => ({
    list: {
      width: props => props.width || 360,
      height:"100%",
      backgroundColor: MulwiColors.whiteBackground,
    },
    toolbar: {
      marginTop: sm ? 0 : 75,
      marginBottom: 7,
      paddingBottom: 7,
      //backgroundColor: "white"
    },
  }));

  const classes = useStyles(props)

  return (
    <React.Fragment>
      <Drawer
        variant={"persistent"}
        anchor={"right"}
        open={props.open}
        onClose={props.onClose}>
        <div
          style={{
            overflowX: "hidden",
          }}
          className={classes.list}
          role="presentation" >
          <div className={classes.toolbar}>
              <Grid container direction="row" 
                  style={{paddingRight: 5, paddingLeft: 5}}
                  justify="space-between"
                  spacing={1}
                  alignItems="center">
                <Grid item>
                    <IconButton onClick={props.onClose}>
                    <ChevronRightIcon />
                  </IconButton>
                </Grid>
                <Grid xs item>
                  {props.navContent}
                </Grid>
              </Grid>
            </div>
          {/* <div style={{ height:"80vh", overflowY:"auto", 
                    overflowX: "hidden",
              paddingRight: 20, paddingLeft: 20}}>
            {props.content}
          </div> */}
            {props.content}
        </div>
    </Drawer>
    {props.children && (
      <div style={{
        marginTop: sm ? 70 : 15,
        overflow: "hidden",
      }}>
        {props.children}
      </div>
    )}
    </React.Fragment>
  );
}