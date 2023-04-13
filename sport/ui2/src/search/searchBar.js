import { Search } from "@mui/icons-material";
import {
  Button,
  Card,
  CardContent,
  CircularProgress, createTheme, Fab,
  Grid, Typography, useMediaQuery
} from "@mui/material";
import { makeStyles } from "@mui/styles";
import React, { useEffect, useState } from "react";
import { useHistory } from 'react-router';
import { Link } from "react-router-dom";
import { getSupportedLanguage, locale2 } from "../locale";
import { MulwiColors } from "../mulwiColors";
import '../veidly-styles.css';
import GoogleMaps from "./googlemaps";
import TagEditor from "./tagEditor";

const useStyles = makeStyles((theme) => ({
  searchBar: {
    marginBottom: 2,
    [theme.breakpoints.down("sm")]: {
      marginTop: 50,
    },
    [theme.breakpoints.up("sm")]: {
      padding: 40,
    },
    minWidth: 450,
    backgroundColor: "inherit",
    width: "100%",
    boxShadow: "none"
  },
  navElements: {
    [theme.breakpoints.down("lg")]: {
      width: 120,
    },
    [theme.breakpoints.between("lg", "xl")]: {
      width: 200,
      //  marginLeft: 40,
    },
    [theme.breakpoints.up("xl")]: {
      width: 250,
    },
    backgroundColor: "white",
    borderRadius: 5,
  },
  // two fields for user input
  searchBarElements: {
    [theme.breakpoints.down("sm")]: {
      marginTop: 10,
      minWidth: 250,
    },
    // for column alignment
    // [theme.breakpoints.up("sm")]: {
    //   marginLeft: 40,
    // },
    // for row alignement
    [theme.breakpoints.up("sm")]: {
      width: 300,
      //marginTop: 80,
      //  marginLeft: 40,
    },
    borderRadius: 5,
    backgroundColor: "white",
  },
  // icon is icon
  searchBarIcon: {
    [theme.breakpoints.down("sm")]: {
      marginTop: 10,
    },
    // [theme.breakpoints.up("sm")]: {
    //   marginLeft: 40,
    // },
    backgroundColor: MulwiColors.greenDark,
    "&:hover": {
      backgroundColor: MulwiColors.pinkAction,
    },
  },
  // END of searchBar section
}));

export default function SearchBar(props) {
  //const [shouldRedirect, setShouldRedirect] = useState(false);
  const styles = useStyles();
  const theme = createTheme();
  const downSm = useMediaQuery(theme.breakpoints.down("sm"));
  const downXs = useMediaQuery(theme.breakpoints.down("xs"));
  const downXl = useMediaQuery(theme.breakpoints.down("1185"));

  // const [askedForPermission, setAskedForPermission] = useState(false)
  const [location, setLocation] = useState(null)

  const [searchRequest, _setSearchRequest] = useState({ Query: "" })
  function setSearchRequest(x) {
    x = x(searchRequest)
    if (!x.Pagination) {
      x.Pagination = {
        Page: 0,
        Size: 100
      }
    }
    if (!x.Lang)
      x.Lang = getSupportedLanguage()
    if (!x.DistKm)
      x.DistKm = 35
    if (!x.DateStart && !x.DateEnd) {
      let d = new Date()
      d.setHours(0, 0, 0, 0)
      //d.setDate(1)
      x.DateStart = d
      let e = new Date(d)
      //e.setHours(23, 59, 59, 0)
      e.setMonth(e.getMonth() + 4)
      e.setDate(e.getDate() - 1)
      x.DateEnd = e
    }
    _setSearchRequest(x)
    if (props.setSearchRequest)
      props.setSearchRequest(x)
  }

  const history = useHistory()

  function getResultUrl() {
    return "/search/l?q=" + encodeURIComponent(JSON.stringify(searchRequest))
  }

  useEffect(() => {
    if (props.searchRequest) {
      _setSearchRequest(props.searchRequest)
      if (!location || (location.lat !== props.searchRequest.Lat &&
        location.lon !== props.searchRequest.Lng))
        setLocation({
          lat: props.searchRequest.Lat,
          lon: props.searchRequest.Lng,
          display_name: props.searchRequest.display_name
        })
    }
  }, [props.searchRequest])

  return (
    (props.simple && (
      <React.Fragment >
        <Grid
          spacing={props.nav ? 1 : 3}
          container
          direction={downSm ? "column" : "row"}
          justifyContent="center"
          style={{
            marginTop: (downXl && props.belowNav && !downXs) ? 60 : null,
          }}
          alignItems="center">
          <Grid item style={{display: props.nav ? "none" : ""}}>
            <Typography variant={"h4"} style={{ fontWeight: "300", color: MulwiColors.whiteBackground, paddingLeft: 40 }}>
              {locale2.SKIING_DANCING_BOX[props.lang]}
            </Typography>
          </Grid>
          <Grid item>
            <TagEditor lang={props.lang}
              noshrink
              val={searchRequest.Query}
              class={props.nav ? "tagLocationNavbar" : "tagLocation"}
              setVal={(e) => {
                setSearchRequest(c => ({ ...c, Query: e }))
              }} />
          </Grid>
          <Grid item >
            <GoogleMaps
              lang={props.lang}
              class={props.nav ? "tagLocationNavbar" : "tagLocation"}
              location={location}
              setLocation={e => {
                let lat = Number(e.lat)
                let lng = Number(e.lon)
                if (lat && lng) {
                  setSearchRequest(c => ({ ...c, Lat: lat, Lng: lng, display_name: e.display_name }))
                  setLocation(e)
                }
              }}
            />
            {/* <LocationAtc lang={props.lang}
              size="small"
              noshrink
              class={props.nav ? styles.navElements : styles.searchBarElements}
              setLocation={e => {
                let lat = Number(e.lat)
                let lng = Number(e.lon)
                if (lat && lng) {
                  setSearchRequest(c => ({ ...c, Lat: lat, Lng: lng, display_name: e.display_name }))
                  setLocation(e)
                }
              }}
              onConfirm={() => {
                history.push(getResultUrl())
              }}
              location={location} /> */}
          </Grid>
          <Grid item lg={props.column ? 8 : null}>

            <Link
              to={getResultUrl()}>
              {props.fullSearchBtn ? (
                <Button style={{
                  backgroundColor: MulwiColors.greenDark,
                  color: "white",
                }}>
                  {locale2.SEARCH[props.lang]}
                </Button>
              ) : (
                <Fab
                  size="small"
                  color="primary"
                  aria-label="add"
                  className={styles.searchBarIcon}>

                  {props.loading ?
                    (<CircularProgress style={{
                      color: "white"
                    }} />) :
                    (<Search />)}
                </Fab>
              )}
            </Link>
          </Grid>
        </Grid>
      </React.Fragment>
    )) || (
      <React.Fragment>
        <Card className={styles.searchBar}>
          <CardContent>
            <Grid
              spacing={2}
              container

              direction={downXl ? "column" : "row"}
              justifyContent="center"
              alignItems="center"
            >
              <Grid item xs={12}>
                <Typography variant={"h4"} style={{ fontWeight: "300", color: MulwiColors.whiteBackground, marginBottom: 30 }} align="center">
                  {locale2.SKIING_DANCING_BOX[props.lang]}
                </Typography>
              </Grid>
              <Grid item style={!downXl ? { paddingTop: 0, paddingBottom: 0, marginTop: 0, marginBottom: 0 } : {}}>
                <TagEditor lang={props.lang}
                  noshrink
                  val={searchRequest.Query}
                  class={"tagLocation"}
                  setVal={(e) => {
                    setSearchRequest(c => ({ ...c, Query: e }))
                  }} />
              </Grid>
              <Grid item style={!downXl ? { paddingTop: 0, paddingBottom: 0, marginTop: 0, marginBottom: 0 } : {}}>
                {/* <LocationAtc lang={props.lang}
                  noshrink
                  class={styles.searchBarElements}
                  setLocation={e => {
                    let lat = Number(e.lat)
                    let lng = Number(e.lon)
                    if (lat && lng) {
                      setSearchRequest(c => (
                        { ...c, Lat: lat, Lng: lng, display_name: e.display_name }))
                      setLocation(e)
                    }
                  }}
                  onConfirm={() => {
                    history.push(getResultUrl())
                  }}
                  location={location} /> */}
                <GoogleMaps
                  lang={props.lang}
                  class={"tagLocation"}
                  location={location}
                  setLocation={e => {
                    let lat = Number(e.lat)
                    let lng = Number(e.lon)
                    if (lat && lng) {
                      setSearchRequest(c => ({ ...c, Lat: lat, Lng: lng, display_name: e.display_name }))
                      setLocation(e)
                    }
                  }}
                />
              </Grid>
              <Grid item style={{ paddingTop: 0, paddingBottom: 0, marginTop: 0, marginBottom: 0 }}>
                <Link
                  to={getResultUrl()}>
                  <Fab
                    style={{
                      backgroundColor: MulwiColors.greenDark,
                      color: "white"
                    }}
                    aria-label="add"
                    className={styles.searchBarIcon}>
                    <Search />
                  </Fab>
                </Link>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      </React.Fragment>
    )
  );
}
