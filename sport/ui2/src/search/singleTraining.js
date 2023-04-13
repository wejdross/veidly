import RoomIcon from '@mui/icons-material/Room';
import { Rating } from "@mui/lab";
import {
  Avatar, Dialog, Grid,
  Typography,
  useMediaQuery,
  useTheme
} from "@mui/material";
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import CardMedia from '@mui/material/CardMedia';
import { red } from '@mui/material/colors';
import IconButton from '@mui/material/IconButton';
import { makeStyles } from "@mui/styles";
import React, { useEffect, useState } from "react";
import MulwiMap from "../gmaps/maps";
import { prettyPrintDate } from "../harmonogram/trainingDetails";
import { defaultLogoPath, prettyPrintCurrency, randomString } from "../helpers";
import { locale2 } from "../locale";
import { MulwiColors } from "../mulwiColors";
import { DisabledSupport } from "../reservations/disabled";
import { ImgOverlay } from "../training/ImgOverlay";

const useStyles = makeStyles((theme) => ({
  media: {
    height: 0,
    paddingTop: '56.25%', // 16:9
    objectFit: 'scale-down'
  },
  expand: {
    transform: 'rotate(0deg)',
    marginLeft: 'auto',
    transition: theme.transitions.create('transform', {
      duration: theme.transitions.duration.shortest,
    }),
  },
  expandOpen: {
    transform: 'rotate(180deg)',
  },
  avatar: {
    backgroundColor: red[500],
    width: 60,
    height: 60
  },
}));

export default function SingleTraining(props) {
  const classes = useStyles();

  const expanded = false
  const [openDialog, setOpenDialog] = React.useState(false)

  // const handleExpandClick = (e) => {
  //   e.stopPropagation()
  //   setExpanded(!expanded);
  // };

  const theme = useTheme()
  const isSmall = useMediaQuery(theme.breakpoints.down('xs'))
  // if(isSmall) return null

  const imgid = "i" + randomString(12)
  const cid = "c" + randomString(12)

  // useEffect(() => {
  //   let e = document.getElementById(imgid)
  //   if (e) {
  //     let c = document.getElementById(cid)
  //     let dh = (e.clientHeight - 106)
  //     if (caches) {
  //       c.style.height = dh + "px"
  //     }
  //   }
  // }, [])

  useEffect(() => {
    if (props.hover)
      setShadow(1)
    else
      setShadow(0)
  }, [props.hover])

  const [shadow, setShadow] = useState(0)

  const onMouseEnter = () => {
    setShadow(1)
    props.onMouseEnter && props.onMouseEnter()
  }

  const onMouseLeave = () => {
    setShadow(0)
    props.onMouseLeave && props.onMouseLeave()
  }

  function getFirstFreeTraining(s) {
    if (s) {
      for (let i = 0; i < s.length; i++) {
        let x = s[i]
        if (x.IsAvailable)
          return <React.Fragment>
            <Grid item>
              <Typography variant="body2" color="textSecondary" component="span">
                {locale2.NEAREST_DATE[props.lang]}  </Typography>
            </Grid>
            <Grid item>
              <Typography style={{
                marginLeft: props.list ? 5 : 0,
              }} variant="body2" component="span">
                <strong>{prettyPrintDate(new Date(x.Start), props.lang)}</strong> </Typography>
            </Grid>
          </React.Fragment>
      }
    }
    return ` ${locale2.NEAREST_DATE[props.lang]}`
  }

  const [ovOpen, setOvOpen] = useState(false)

  return (
    <React.Fragment>
      {/* <SelectOcc
        user={props.user}
        dr={props.dr}
        open={occOpen}
        schedule={props.d.Schedule}
        onSelect={navigateRsv}
        training={props.d.Training}
        d={props.d}
        setOpen={setOccOpen} /> */}
      <ImgOverlay
        open={ovOpen}
        setOpen={setOvOpen}
        MainImgUrl={props.d.Training.MainImgUrl}
        SecondaryImgUrls={props.d.Training.SecondaryImgUrls}
      />
      <Card id={'listt' + props.d.Training.ID} onClick={props.onClick} onMouseEnter={onMouseEnter} style={{
        cursor: "pointer",
        height: "100%",
        minWidth: 160,
        borderRadius: 0,
        width: "100%"
      }}
        onMouseLeave={onMouseLeave} raised={Boolean(shadow)}>
        <Grid container direction="row" style={{ margin: "none" }}>
          <Grid item lg>
            <Grid container direction="column" style={{ height: "100%" }} justifyContent="space-between">

              <Grid item>
                {
                  !props.list && props.d.Training.MainImgUrl && (
                    <CardMedia
                      style={{
                        minWidth: 160,
                      }}
                      className={classes.media}
                      image={props.d.Training.MainImgUrl || (defaultLogoPath)}
                      title="training"
                    />
                  )
                }

                <CardHeader
                  avatar={props.d.UserInfo.AvatarUrl &&
                    <Avatar aria-label="recipe" className={classes.avatar}
                      src={props.d.UserInfo.AvatarUrl}>
                      R
                    </Avatar>
                  }
                  title={<Typography variant="h6"><strong>{props.d.Training.Title}</strong></Typography>}
                  subheader={<React.Fragment>
                    <DisabledSupport small training={props.d.Training} lang={props.lang} />
                    <strong>{props.d["UserInfo"]["Name"]}</strong>
                    {" " + locale2.INSTRUCTOR_SINCE[props.lang]} {(() => {
                      let x = new Date(
                        props.d["Instructor"]["CreatedOn"]
                      );
                      return x.getFullYear();
                    })()}
                  </React.Fragment>}
                />

                <CardContent id={cid} style={{ paddingTop: 0 }}>
                  <Grid container direction="column"
                    alignItems={props.list ? "flex-start" : "center"}
                    justify={"center"} spacing={1}>
                    {!props.suggestion && (
                      <Grid item style={{ height: "100%" }}>
                        <Grid style={{
                          marginLeft: props.list ? 20 : 0,
                          width: "100%"
                        }} container direction="row" alignItems="center"
                          justify={!props.list ? "center" : null}>
                          {getFirstFreeTraining(props.d.Schedule)}
                        </Grid>
                      </Grid>
                    )}
                    {!props.list && (
                      <React.Fragment>
                        <Grid item style={{
                          width: "100%"
                        }}>
                          <Grid container direction="row" justifyContent="space-between">
                            <Grid item>
                              <Grid container direction="row">
                                <Rating readOnly value={props.d.Training.AvgMark} max={6} />
                                <Typography>
                                  ({props.d.Training.NumberReviews})
                                </Typography>
                              </Grid>
                            </Grid>
                            <Grid item style={{
                              marginLeft: 5,
                            }}>
                              <Typography>
                                <strong>{props.d.Training.Price / 100}</strong>
                                {" " + prettyPrintCurrency(props.d.Training.Currency)}
                              </Typography>
                            </Grid>
                          </Grid>
                        </Grid>

                        <Grid item>
                          <Grid container direction="row" alignItems="center">
                            <Grid item>
                              <IconButton size="small" aria-label="share"
                                style={{ color: MulwiColors.blueLight }}
                                onClick={(e) => {
                                  e.stopPropagation()
                                  setOpenDialog(!openDialog)
                                }}>
                                <Dialog open={openDialog} >
                                  <div style={{ padding: 30 }}>
                                    <center>
                                      <Typography>{props.d.Training.Title}</Typography>
                                    </center>
                                    <MulwiMap center={props.d.Training.LocationText} />
                                  </div>
                                </Dialog>
                                <RoomIcon />
                              </IconButton>
                            </Grid>
                            <Grid item>
                              <Typography variant="body2" color="textSecondary">
                                {props.d.Training.LocationText}
                              </Typography>
                            </Grid>
                          </Grid>
                        </Grid>
                      </React.Fragment>
                    )}
                  </Grid>
                </CardContent>
              </Grid>

              {props.list && (<Grid item><CardActions>
                <Rating readOnly value={props.d.Training.AvgMark} max={6} />
                <Typography style={{
                  marginLeft: 10
                }} variant="body2" component="span">
                  ({props.d.Training.NumberReviews})
                </Typography>
                <Typography style={{ marginLeft: 10 }}>
                  <strong>{props.d.Training.Price / 100}</strong>
                  {" " + prettyPrintCurrency(props.d.Training.Currency)}
                </Typography>

                <Grid item>
                  <Grid style={{
                    marginLeft: 5
                  }} container direction="row" alignItems="center">
                    <Grid item>
                      <IconButton size="small" aria-label="share"
                        style={{ color: MulwiColors.blueLight }}
                        onClick={(e) => {
                          e.stopPropagation()
                          setOpenDialog(!openDialog)
                        }}>
                        <Dialog open={openDialog} >
                          <div style={{ padding: 30 }}>
                            <center>
                              <Typography>{props.d.Training.Title}</Typography>
                            </center>
                            <MulwiMap center={props.d.Training.LocationText} />
                          </div>
                        </Dialog>
                        <RoomIcon />
                      </IconButton>
                    </Grid>
                    <Grid item>
                      <Typography variant="body2" color="textSecondary">
                        {props.d.Training.LocationText}
                      </Typography>
                    </Grid>
                  </Grid>
                </Grid>
              </CardActions></Grid>)}
            </Grid>
          </Grid>

          {props.list && props.d.Training.MainImgUrl && (
            <Grid item xs={4}>
              <img
                alt="main-training"
                onClick={(e) => {
                  e.stopPropagation()
                  setOvOpen(true)
                }}
                id={imgid}
                style={{
                  width: isSmall ? "100%" : "auto",
                  height: isSmall ? "100%" : "100%",
                  margin: "auto",
                  maxHeight: 220,
                  display: "block",
                  marginBottom: -6,
                  borderBottomLeftRadius: expanded ? 15 : 0
                }}
                src={props.d.Training.MainImgUrl}
              />
            </Grid>)}

        </Grid>
      </Card>
    </React.Fragment >
  );
}
