import {
  Container, Fab,
  Slider, TextField, Button, Checkbox, FormControlLabel
} from "@mui/material";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import { DatePicker as KeyboardDatePicker } from '@mui/x-date-pickers/DatePicker';
import { TimePicker} from '@mui/x-date-pickers/TimePicker';
import React, { useEffect, useState } from "react";
import { MulwiColors } from "../mulwiColors";
import { dayIndex } from '../helpers'
import { SortContent } from "./SortContent";
import { DateBar } from "./Datebar";
import { locale2 } from "../locale";
import { makeStyles } from "@mui/styles";

const useStyles = makeStyles((theme) => ({
  // leaving it here as reference, to rapidly start developement
  paper: {
    height: "100%",
    minWidth: 300,
    maxWidth: 550,
    [theme.breakpoints.down("md")]: {
      minWidth: 100,
      maxWidth: 250,
      backgroundColor: MulwiColors.whiteSurface,
    },
  },
  arrowIcon: {
    height: 50,
    width: 50,
    color: MulwiColors.subtitleTypography,
  },
  accordion: {
    [theme.breakpoints.up("md")]: {
      // marginTop: 50,
    },
  },
  Slider: {
    color: MulwiColors.greenDark,
    maxWidth: 250,
    marginLeft: 20,
    [theme.breakpoints.down("md")]: {
      marginLeft: 20,
      marginTop: 15
    },
  },
  fabLevel: {
    color: MulwiColors.whiteSurface,
    boxShadow: "none",
  },
  filterButton: {
    marginTop: 20,
    backgroundColor: MulwiColors.greenDark,
  },
  topPaper: {
    marginLeft: 30,
    height: 300,
  },
  headingTop: {
    marginTop: 50,
    fontSize: 30,
  },
  heading: {
    fontSize: "1.2em",
  },
  typo: {
    marginBottom: 15,
    marginTop: 10,
  },
  separator: {
    marginBottom: 40,
  },
  fabFilterButton: {
    color: MulwiColors.whiteSurface,
    backgroundColor: MulwiColors.greenDark,
    height: 48,
    width: "100%",
    fontSize: "1.5em",
    "&:hover": {
      backgroundColor: MulwiColors.blueDark,
      color: MulwiColors.whiteSurface,
    },
    [theme.breakpoints.down("md")]: {

      marginTop: 20,
      color: MulwiColors.whiteSurface,
      backgroundColor: MulwiColors.greenDark,
      height: 50,
      width: 200,
      "&:hover": {
        backgroundColor: MulwiColors.blueDark,
        color: MulwiColors.whiteSurface,
      },
  
      fontSize: "1.2em",
    },
  },
  appBarsmall: {
    top: "auto",
    bottom: 0,
  },
  appBarButton: {
    width: "100%",
    backgroundColor: MulwiColors.whiteSurface,
    color: MulwiColors.blueDark,
  },
  ContainerOnBig: {
    backgroundColor: MulwiColors.whiteSurface,
    width: "100%",
    [theme.breakpoints.up("md")]: {
      // marginTop: 50,
      padding: 20,
      paddingRight: 30,
    },
  },
}));

const priceMax = 2400000
const capacityMax = 29979

export function FilterContent(props) {
  /*
    below var diffs is bad hack to keep code everything running
    this must be completely refactored before release, I'm not deleting
    diffs.js to make it simpler to read what and why is taken from
  */
  var diffs = {
    1: locale2.BASIC_LEVEL[props.lang],
    2: locale2.INTERMEDIATE_LEVEL[props.lang],
    3: locale2.ADVANCED_LEVEL[props.lang]
  }
  const classes = useStyles();
  const [price, setPrice] = useState([0, 0]);
  const [age, setAge] = useState(0);

  //const [sdiffs, setsdiffs] = React.useState([])
  const [level, _setLevel] = useState([]);
  function setLevel(x) {
    _setLevel(x || [])
  }
  const [days, _setDays] = useState([]);
  function setDays(x) {
    _setDays(x || [])
  }
  const [allowedDistance, setAllowedDistance] = useState(10);
  const [capacity, setCapacity] = useState([0, 0]);
  const [minHour, setMinHour] = useState(null);
  const [maxHour, setMaxHour] = useState(null);
  const [minHourToReturn, setMinHourToReturn] = useState(null);
  const [maxHourToReturn, setMaxHourToReturn] = useState(null);

  const [trainingSupportsDisabled, setTrainingSupportsDisabled] = useState(false)
  const [placeSupportsDisabled, setPlaceSupportsDisabled] = useState(false)

  function restoreDefaults() {
    setPrice([0, 0])
    setLevel([])
    setDays([])
    setAllowedDistance(10)
    setCapacity([0, 0])
    setMinHour(null)
    setMaxHour(null)
    setMinHourToReturn(null)
    setMaxHourToReturn(null)
    setAge(0)
    setTrainingSupportsDisabled(false)
    setPlaceSupportsDisabled(false)
  }

  function reverseMergeRequest() {

    let r = props.searchRequest || {}

    if (r.DistKm) {
      setAllowedDistance(r.DistKm)
    }

    setPrice([
      r.PriceMin ? Math.round(r.PriceMin / 100) : 0,
      r.PriceMax === priceMax ? 0 : (Math.round(r.PriceMax / 100) || 0)
    ])

    setAge(r.Age || 0)

    setLevel(r.Diffs || [])

    setDays(r.Days || [])

    if (r.CapacityMin || r.CapacityMax) {
      setCapacity([r.CapacityMin, r.CapacityMax])
    }

    setMinHourToReturn(r.HrStart || null)
    setMaxHourToReturn(r.HrEnd || null)

    setPlaceSupportsDisabled(r.placeSupportsDisabled || false)
    setTrainingSupportsDisabled(r.trainingSupportsDisabled || false)

  }

  // merge existing api search request with filter
  function mergeRequest() {

    let r = props.searchRequest

    if (!price[1] || price[1] > priceMax) {
      price[1] = priceMax
    }

    if (price[0] <= price[1]) {
      r.PriceMin = Math.round(price[0] * 100)
      r.PriceMax = Math.round(price[1] * 100)
    }

    r.Age = age || 0

    r.Diffs = level

    r.Days = days

    if (r.DistKm != allowedDistance) {
      props.setUdist && props.setUdist(allowedDistance)
    }

    r.DistKm = allowedDistance

    if (!capacity[1] || capacity[1] > capacityMax) {
      capacity[1] = capacityMax
    }

    if (capacity[0] <= capacity[1]) {
      r.CapacityMin = capacity[0]
      r.CapacityMax = capacity[1]
    }

    if (minHourToReturn && maxHourToReturn && minHour <= maxHour) {
      r.HrStart = minHourToReturn
      r.HrEnd = maxHourToReturn
    }

    r.placeSupportsDisabled = placeSupportsDisabled
    r.trainingSupportsDisabled = trainingSupportsDisabled

    return r
  }

  const [isMount, setIsMount] = useState(true)

  useEffect(() => {
    setIsMount(false)
  }, [])

  useEffect(() => {
    if (isMount)
      return
    if (props.updateToken) {
      props.onChange(mergeRequest())
    }
  }, [props.updateToken])

  useEffect(() => {
    reverseMergeRequest()
  }, [props.searchRequest])

  return (
    <>
      <Container
        className={classes.ContainerOnBig}>
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.SORTING[props.lang]}
          </Typography>
          <Grid container spacing={3} justify="flex-end" alignItems="stretch">
            <Grid item>
              <SortContent lang={props.lang}
                searchRequest={props.searchRequest} onChange={props.onChange} />
            </Grid>
          </Grid>
        </div>
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.WHICH_DATES_DO_YOU_PREFER[props.lang]}
          </Typography>
          <Grid container spacing={3} justify="flex-end" alignItems="stretch">
            <Grid item>
              <DateBar lang={props.lang}
                searchRequest={props.searchRequest} onChange={props.onChange} />
            </Grid>
          </Grid>
        </div>
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.ADJUST_PRICE[props.lang]}
          </Typography>
          <Grid container direction="row"
            justify="flex-end"
            spacing={2}
            alignItems="center">
            <Grid item>
              <Typography>{locale2.FROM[props.lang]}</Typography>
            </Grid>
            <Grid item style={{
              width: 100
            }}>
              <TextField
                value={price[0] === 0 ? "" : String(Number(price[0]))}
                onChange={(e) => {
                  let v = e.target.value
                  if (v === "") {
                    v = 0
                  } else {
                    v = Number(v) || price[0]
                  }
                  setPrice([v, price[1]]);
                }}
                variant="outlined" size="small" />
            </Grid>
            <Grid item>
              <Typography>{locale2.TO[props.lang]}</Typography>
            </Grid>
            <Grid item style={{
              width: 100
            }}>
              <TextField
                value={price[1] === 0 ? "" :
                  (price[1] >= priceMax ? "" : String(Number(price[1])))}
                onChange={(e) => {
                  let v = e.target.value
                  if (v === "") {
                    v = 0
                  } else {
                    v = Number(v) || price[1]
                  }
                  setPrice([price[0], v]);
                }}
                variant="outlined" size="small" />
            </Grid>
          </Grid>
        </div>

        <div className={classes.separator}>
          <Typography style={{ marginBottom: 25 }} className={classes.typo}>
            {locale2.ADJUST_DISTANCE[props.lang]}
          </Typography>
          <Grid container spacing={3} justify="flex-end" alignItems="stretch">
            <Grid item>
              <Slider
                style={{
                  width: 300,
                }}
                value={allowedDistance}
                aria-labelledby="discrete-slider"
                valueLabelDisplay="on"
                valueLabelFormat={(v) => <span style={{
                  fontSize: 9
                }}>{v} km</span>}
                step={5}
                marks
                min={5}
                max={100}
                className={classes.Slider}
                onChange={(e, n) => {
                  setAllowedDistance(n);
                }}
              />
            </Grid>
          </Grid>
        </div>
        {
          // LEVEL
        }
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.ADJUST_LEVEL[props.lang]}
          </Typography>
          <Grid
            container
            direction="row"
            spacing={2}
            justify="flex-end"
            alignItems="center"
          >
            {[1, 2, 3].map((elem, id) => {
              return (
                <Grid item key={id}>
                  <Fab
                    size="small"
                    variant={"extended"}
                    key={id}
                    className={classes.fabLevel}
                    style={
                      level.indexOf(elem) >= 0
                        ? {
                          backgroundColor: MulwiColors.greenDark,
                        }
                        : {
                          color: MulwiColors.blueLight,
                        }
                    }
                    onClick={(event) => {
                      if (level.indexOf(elem) === -1) {
                        let x = [...level]
                        x.push(elem)
                        setLevel(x)
                      } else {
                        let x = []
                        for (let i = 0; i < level.length; i++) {
                          if (level[i] !== elem) x.push(level[i])
                        }
                        setLevel(x)
                      }
                    }}
                  >
                    {diffs[elem] ? diffs[elem] : elem}
                  </Fab>
                </Grid>
              );
            })}
          </Grid>
        </div>

        {
          // CAPACITY
        }
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.ADJUST_GROUP_SIZE[props.lang]}
          </Typography>
          <Grid container direction="row"
            justify="flex-end"
            spacing={2}
            alignItems="center">
            <Grid item>
              <Typography>{locale2.FROM[props.lang]}</Typography>
            </Grid>
            <Grid item style={{
              width: 100
            }}>
              <TextField
                value={capacity[0] === 0 ? "" : String(Number(capacity[0]))}
                onChange={(e) => {
                  let v = e.target.value
                  if (v === "") {
                    v = 0
                  } else {
                    v = Number(v) || capacity[0]
                  }
                  setCapacity([v, capacity[1]]);
                }}
                variant="outlined" size="small" />
            </Grid>
            <Grid item>
              <Typography>{locale2.TO[props.lang]}</Typography>
            </Grid>
            <Grid item style={{
              width: 100
            }}>
              <TextField
                value={capacity[1] === 0 ? "" :
                  (capacity[1] >= capacityMax ? "" : String(Number(capacity[1])))}
                onChange={(e) => {
                  let v = e.target.value
                  if (v === "") {
                    v = 0
                  } else {
                    v = Number(v) || capacity[1]
                  }
                  setCapacity([capacity[0], v]);
                }}
                variant="outlined" size="small" />
            </Grid>
          </Grid>
        </div>

        {
          // DAYS
        }
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.WHICH_DAYS_DO_YOU_PREFER[props.lang]}
          </Typography>
          <Grid
            container
            direction="row"
            justify="flex-end"
            alignItems="center"
            spacing={2}
          >
            {[1, 2, 3, 4, 5, 6, 7].map((day, id) => {
              return (
                <Grid item key={id}>
                  <Fab
                    size="small"
                    variant={"extended"}
                    className={classes.fabLevel}
                    style={
                      days.indexOf(day) >= 0
                        ? {
                          backgroundColor: MulwiColors.greenDark,
                        }
                        : {
                          color: MulwiColors.blueLight,
                        }
                    }
                    onClick={(event) => {
                      if (days.indexOf(day) === -1) {
                        let x = [...days]
                        x.push(day)
                        setDays(x)
                      } else {
                        let x = []
                        for (let i = 0; i < days.length; i++) {
                          if (days[i] !== day) x.push(days[i])
                        }
                        setDays(x)
                      }
                    }}>
                    {dayIndex[day] ? dayIndex[day][props.lang] : day}
                  </Fab>
                </Grid>
              );
            })}
          </Grid>
        </div>
        {
          // training time
        }
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.ADJUST_HOURS[props.lang]}
          </Typography>
          <Grid
            container
            direction="row"
            justify="flex-end"
            alignItems="center"
            spacing={2}
          >
            <Grid item>
              <TimePicker
                label={locale2.FROM[props.lang]}
                size="small"
                value={minHour}
                ampm={false}
                ampmInClock={false}
                style={{
                  width: 120
                }}
                onChange={(date, timeString) => {
                  setMinHourToReturn(timeString);
                  setMinHour(date);
                }}
                KeyboardButtonProps={{
                  "aria-label": "change time",
                }}
                renderInput={(params) => <TextField {...params} />}
              />
            </Grid>
            <Grid item>
              <TimePicker
                label={locale2.TO[props.lang]}
                size="small"
                ampmInClock={false}
                style={{
                  width: 120
                }}
                value={maxHour}
                ampm={false}
                onChange={(date, timeString) => {
                  setMaxHourToReturn(timeString);
                  setMaxHour(date);
                }}
                KeyboardButtonProps={{
                  "aria-label": "change time",
                }}
                renderInput={(params) => <TextField {...params} />}
              />
            </Grid>

          </Grid>
        </div>
        <div className={classes.separator}>
          <Typography className={classes.typo}>
            {locale2.ADJUST_AGE[props.lang]}
          </Typography>
          <Grid container justify="flex-end">
            <TextField
              type="number"
              variant="outlined"
              size="small"
              value={age || ""}
              onChange={(v) => {
                let x = Number(v.target.value || 0)
                if (x > 120) return
                setAge(x)
              }}
              style={{
                width: 100
              }}
            />
          </Grid>
        </div>

        <div className="classes.separator">

          <FormControlLabel
            control={<Checkbox
              checked={trainingSupportsDisabled}
              onChange={(c, v) => {
                setTrainingSupportsDisabled(v)
              }}
            />}
            label={
              <React.Fragment>
              <Grid container direction="row">
                <img height={25} style={{
                  marginBottom: -5, marginRight: 5
                }} src="/1.svg" alt="" />
                <Typography>
                  {locale2.TRAINING_SUPPORTS_DISABLED[props.lang]}
                </Typography>
              </Grid>
              </React.Fragment>}
          />
        </div>

        <div className="classes.separator">
          <FormControlLabel
            control={<Checkbox
              checked={placeSupportsDisabled}
              onChange={(c, v) => {
                setPlaceSupportsDisabled(v)
              }}
            />}
            label={<React.Fragment>
              <Grid container direction="row">
                <img height={25} style={{
                  marginBottom: -5, marginRight: 5
                }} src="/2.svg" alt="" />
                <Typography>
                  {locale2.PLACE_SUPPORTS_DISABLED[props.lang]}
                </Typography>
              </Grid>
            </React.Fragment>}
          />
        </div>

        <br/>

        <Button onClick={restoreDefaults}>
          {locale2.RESET_FILTERS[props.lang]}
        </Button>
      </Container>
    </>
  )
}
