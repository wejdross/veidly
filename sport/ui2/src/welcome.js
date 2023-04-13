import {
  Button,
  Card,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Grid,
  Typography,
  useMediaQuery,
} from "@mui/material";
import React, { useEffect, useState } from "react";
import { Parallax } from "react-parallax";
import StickyFooter from "./Footer";
import { makeStyles } from "@mui/styles";
import SearchBar from "./search/searchBar";
import { MulwiColors } from "./mulwiColors";
import { Link } from "react-router-dom";
import { locale2 } from "./locale";
import { JoinUs } from "./JoinUs";
import { createTheme } from "@mui/material";
import Divider from '@mui/material/Divider';
import { Container, width } from "@mui/system";
import { removeFromQs } from "./helpers";
import { DonateForm } from "./donations/donate";


const mS = makeStyles((theme) => ({
  paralax: {
    height: "80vh",
    minHeight: 600,
    [theme.breakpoints.down("sm")]: {
      height: "85vh",
      marginTop: 50,
    },
  },
  //  dancing, boxing? blah, blah, blah, this is it
  advertiseText: {
    [theme.breakpoints.down("lg")]: {
      marginTop: 40,
    },
    [theme.breakpoints.up("lg")]: {
      marginTop: 12,
    },
    // width: "100%",
    backgroundColor: "rgba(0,0,0,0)",
    boxShadow: "0 0 0 0",
  },
  // END of searchBar section

  // START of "How it works section

  parallaxContent: {
    marginTop: "15vh",
  },
  primarylinkbutton: {
    transition: "all 0.3s ease",
    "&:hover": {
      backgroundColor: MulwiColors.blueDark,
      color: MulwiColors.whiteBackground,
    },
    //border: "1px solid grey",
    borderRadius: 20,
    textDecoration: "none",
    color: MulwiColors.blackText,
    paddingLeft: 75,
    paddingRight: 75,
    paddingBottom: 75,
    maxWidth: "95vw",
    paddingTop: 20,
    marginTop: 40,
    width: 390,
    height: 300,
    [theme.breakpoints.down("sm")]: {
      marginTop: 20,
      width: 280,
      paddingLeft: 5,
      paddingRight: 5,
      paddingBottom: 5,
      height: 200,
    },
  }
}));

function PrimaryLinkButton(props) {
  const [hovered, setHovered] = useState(false)
  const myStyles = mS();
  return (<React.Fragment>
    <Link to={props.href} style={{ textDecoration: "none" }}>
      <Card className={myStyles.primarylinkbutton} elevation={5} onMouseEnter={(e) => { setHovered(true) }} onMouseLeave={(e) => { setHovered(false) }}>
        <Grid container direction="column" justifyContent="center" alignItems="center">
          <Typography align="center" variant={"h5"} style={{ marginTop: 10 }}>
            {props.text}
          </Typography>
          <Divider style={{
            width: "100%",
            marginBottom: 20,
            backgroundColor: hovered ? MulwiColors.greenDark : ""
          }} />
          <Typography align="center" variant="h6">
            {props.subtext}
          </Typography>
        </Grid>
      </Card>
    </Link>
  </React.Fragment>)
}

export default function Welcome(props) {
  const myStyles = mS();
  const theme = createTheme({
    breakpoints: {
      values: {
        xs: 0,
        sm: 600,
        md: 960,
        lg: 1280,
        xl: 1921,
      },
    },
  });
  const belowSMSize = useMediaQuery(theme.breakpoints.down("sm"));
  const [donated, setDonated] = useState(false)

  useEffect(() => {
    window.scrollTo(0, 0)
    let query = new URLSearchParams(window.location.search)
    let donated = query.get("donated")
    if(donated) {
      setDonated(true)
      removeFromQs("donated")
    }
  }, [])

  const l = props.lang
  if (!l)
    return null

  return (
    <>
      <Dialog open={donated} onClose={() => setDonated(false)}>
        <DialogTitle>
          {locale2.THANKS[props.lang]}
        </DialogTitle>
        <DialogContent>
          Twoje wsparcie pozwala nam na dalszą pracę nad Veidly!
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDonated(false)}>
            OK!
          </Button>
        </DialogActions>
      </Dialog>
      <Parallax
        blur={0}
        strength={-500}
        bgImage={belowSMSize ? "/static/veidly-spring-small.webp" : "/static/veidly-spring-big.webp"}
        bgImageAlt="main paralax"
        className={myStyles.paralax}
      >
        {belowSMSize && (
          <Card className={myStyles.advertiseText}>
            <div style={{
              backgroundColor: "rgba(0,0,0,0.7)",
              marginLeft: "10px",
              marginRight: "10px",
              marginBottom: 15,
              padding: "20px",
              borderTopRightRadius: 50,
              borderBottomLeftRadius: 50,
              boxShadow: "0px 4px rgba(36,192,132,0.9)",
              minWidth: 250,
            }}>
              <SearchBar lang={l} simple />
            </div>
          </Card>
        )}
        {
          // I don't want to use && and || in one brackets, as react complains to not mix those operators...
          !belowSMSize && (
            <center>
              <div className={myStyles.parallaxContent} style={{
                backgroundColor: "rgba(0,0,0,0.7)",
                width: "60%",
                borderTopRightRadius: 120,
                borderBottomLeftRadius: 120,
                boxShadow: "7px 7px rgba(36,192,132,0.9)",
                minWidth: 550
              }}>
                <Grid
                  container
                  direction={"column"}
                  justifyContent={"center"}
                  alignItems={"center"}
                >
                  <Card className={myStyles.advertiseText}>
                    <Grid
                      container
                      direction={"column"}
                      justifyContent={"center"}
                      alignItems={"flex-start"}
                    >
                      <SearchBar lang={l} />
                    </Grid>
                  </Card>
                </Grid>
              </div>
            </center>
          )
        }
      </Parallax>
      {
        // END of Search Bar section
        // START of "How it works section
      }
      <Grid
        container
        style={{
          marginTop: 50
        }}
        direction="column"
        justifyContent="center"
        alignItems="center"
      >

        <Typography variant="h3" align={"center"}>
          {locale2.MEET_PLATFORM[l]}
        </Typography>
        <Typography variant="h4" style={{ color: MulwiColors.blueDark }}>
          {locale2.HOW_IT_WORKS[l]}
        </Typography>
      </Grid>
      <Container>

        <Grid direction={"column"} textAlign={"center"} justifyContent={"center"} alignItem={"center"} container style={{
          marginTop: 20,
        }}>
          {
            // this func just split text over \n character
            String(locale2.MARKETING_DESCRIPTION_MAIN_PAGE[props.lang]).split('\n').map((i, v) => {
              return (

                <Typography variant="body1">
                    {i}
                  </Typography>
                    )
            })
          }
        </Grid>
      </Container>
      <Grid
        container
        direction={belowSMSize ? "column" : "row"}
        justifyContent="space-around"
        alignItems="center"
        style={{
          marginTop: 50
        }}
      >
        <PrimaryLinkButton href={"/benefits_instructor"} alt={"trainings"}
          text={locale2.FOR_INSTRUCTOR[props.lang].toUpperCase()} subtext={locale2.CHECK_AMOUNT_OF_BENEFITS[props.lang]} />
        <PrimaryLinkButton href={"/benefits_user"} alt={"trainings"}
          text={locale2.FOR_TRAINEE[props.lang].toUpperCase()} subtext={locale2.CHECK_HOW_EASY_VEIDLY_IS[props.lang]} />
      </Grid>

      <JoinUs lang={l} />

      <DonateForm lang={props.lang}/>
    </>
  );
}
