import { Facebook, Instagram } from '@mui/icons-material';
import { Grid, useMediaQuery } from '@mui/material';
import Typography from '@mui/material/Typography';
import { makeStyles, useTheme } from '@mui/styles';
import React from 'react';
import { Link } from 'react-router-dom';
import { locale2 } from './locale';
import { MulwiColors } from './mulwiColors';
import { supportEmail } from './StatusDialog';
import "./veidly-styles.css"

function Copyright(props) {
  return (
    <Grid container direction='column' alignContent='center' justifyContent='center'>

      <Grid container direction={"row"} justifyContent="center" alignContent={"center"} style={{
        backgroundColor: MulwiColors.greenDark,
        paddingTop: 20,
        paddingBottom: 20,
      }}>
        <Typography variant="body2" color="textSecondary">
          {'Copyright © ' + new Date().getFullYear()}
          <Link to="/" color="inherit" href="https://veidly.com/" style={{
            textDecoration: "none",
            color: MulwiColors.blueDark,
            marginLeft: 5,
          }}>
            Veidly.com
          </Link>
        </Typography>
        <Grid container direction='row' justifyContent='center' alignContent='center' >
          {
            copyrightURLs(props.small)
          }
        </Grid>
      </Grid >
    </Grid>
  );
}

function copyrightURLs(small) {
  return (
    <>
      <a href={"https://www.facebook.com/veidly"} target="_blank" rel="noreferrer">
        <Facebook style={{ color: "white", marginLeft: small ? 0 : 20, }} />
      </a>
      <a href={"https://www.instagram.com/veidly_official/"} target="_blank" rel="noreferrer">
        <Instagram style={{ color: "white", marginLeft: 20 }} />
      </a>
    </>
  )
}

const useStyles = makeStyles((theme) => ({

  footer: {
    backgroundColor: MulwiColors.black,
    paddingBottom: 50,
  },
  footerDiv: {
    width: "100%",
    color: "white",
    margin: 0,
    padding: 0,
    [theme.breakpoints.down('md')]: {
      margin: 25,
    },

  },
  typographyBodyWhite: {
    color: "white",
    fontSize: "1.2em",
    fontWeight: 300
  }
}));

export default function StickyFooter(props) {
  const classes = useStyles();
  const theme = useTheme()
  const belowSMSize = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <>
      <Grid container direction={belowSMSize ? "column" : "row"} justify={"space-evenly"} alignItems={"center"} className={classes.footer}>
          <img src="/static/footer.png" style={{ width: "100%", height: "100%", backgroundColor: "#fbfbfb", margin: 0, padding: 0, objectFit: "fill", boxShadow: "none", border: "none", color: 'white', fill: "white"}} alt="footer graphic" />
        {
          // customer care
        }
          <Grid container direction={"column"} justify={"center"} alignItems={"center"}>
            <a href={"mailto: support@veidly.com"} style={{
              fontSize: "1.5em",
              textDecoration: "none",
              color: MulwiColors.greenDark
            }}>
              {supportEmail}
            </a>
            <Typography className='typographyBodyWhite'>ul. Lipowa 44</Typography>
            <Typography className='typographyBodyWhite'>43-190 Mikołów</Typography>
            <Typography variant='subtitle2' style={{color: MulwiColors.whiteBackground, marginTop: 50, width: "70%", marginBottom: 0}}>{locale2.SEO_FOOTER[props.lang]}</Typography>
          </Grid>
      </Grid>
      <Grid container direction={"column"} justify={"center"} alignItems={"center"}>
        <Copyright small={belowSMSize ? true : false} />
      </Grid>
    </>
  );
}