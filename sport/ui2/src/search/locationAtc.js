import { CircularProgress, TextField, Autocomplete } from "@mui/material";
import { useEffect, useState } from "react";
import { fetchGeo, reverseFetchGeo } from "../apicalls/user.api";
import { locale2 } from "../locale"
import React from "react";
import { makeStyles } from "@mui/styles";


const useStyles = makeStyles((t) => (
  {
    option: {
      minHeight: 'auto',
      alignItems: 'flex-start',
      padding: 8,
      '&[aria-selected="true"]': {
        backgroundColor: 'transparent',
      },
      '&[data-focus="true"]': {
        backgroundColor: 'rgba(0, 0, 0, 0.15)',
      },
    },
    typographyLineHeight: {
      lineHeight: 2,
    },
    searchBarElements: {
      [t.breakpoints.down("sm")]: {
        marginTop: 20,
        minWidth: 250,
      },
      [t.breakpoints.up("lg")]: {
        minWidth: 350,
        marginLeft: 40,
      },
      backgroundColor: "white",
    }
  }
))

export default function LocationAtc(props) {

  const classes = useStyles()

  const [cv, setCv] = useState("")
  const [th, setTh] = useState(null)

  const [options, setOptions] = useState([])

  useEffect(() => {
    if (props.location && props.location.display_name) {
      if (props.location.display_name === locale2.YOUR_LOCATION[props.lang]) {

        (async () => {
          try {
            let v = { ...props.location }
            let x = await reverseFetchGeo(v.lat, v.lon)
            x = JSON.parse(x)
            console.log(x)
            // v.display_name = x.display_name
            v.display_name = x.address.city || x.address.town || x.address.village
            //_cv = v.display_name
            setCv(v.display_name)
            props.setLocation(v)
          } catch (ex) {
            console.log(ex)
          }
        })()

      } else {
        //_cv = props.location.display_name
        setCv(props.location.display_name)
      }
    }
  }, [props.location])

  function updateOptions(value) {
    if (th) window.clearTimeout(th)
    setTh(setTimeout(async () => {
      try {
        let c = await fetchGeo(value, 1)
        c = JSON.parse(c)
        setOptions(c)
      } catch (ex) {
        console.log(ex)
      } finally {
        setTh(null)
      }
    }, 300))
  }

  return (
    <Autocomplete
      value={props.location}
      onChange={(event, newValue) => {
        if (!newValue) return
        if (newValue.sp === 1) {
          navigator.geolocation.getCurrentPosition(
            async p => {
              let v = {
                lat: p.coords.latitude,
                lon: p.coords.longitude,
                display_name: newValue.display_name
              }
              props.setLocation(v)
            },
            () => console.log("failed to getCurrentPosition"),
            {
              // keep geo in cache one hour
              maximumAge: 3600000,
              enableHighAccuracy: true,
            })
        } else {
          props.setLocation(newValue)
        }
      }}
      filterOptions={(options, params) => {
        options.unshift({
          display_name: locale2.YOUR_LOCATION[props.lang],
          sp: 1
        })
        return options;
      }}
      classes={{
        option: classes.option,
      }}
      fullWidth
      selectOnFocus
      clearOnBlur
      size={props.size}
      handleHomeEndKeys
      filterSelectedOptions
      options={options}
      getOptionLabel={o => o.display_name || ""}
      renderOption={(o) => o.display_name || ""}
      renderInput={(params) => {
        if (th)
          params.InputProps.endAdornment = (
            <CircularProgress style={{ width: 30, height: 30 }} />
          )
        return (
          <TextField
            className={props.class}
            style={{
              marginLeft: `${props.noMargin && 0}`
            }}
            value={cv}
            onChange={async e => {
              setCv(e.target.value)
              //_cv = e.target.value
              if (e.target.value) {
                updateOptions(e.target.value)
              }
            }}
            onKeyPress={props.onConfirm ? (e) => {
              if(e.key === "Enter") {
                props.onConfirm && props.onConfirm()
              }
            } : null}
            variant="outlined"
            {...params}
            InputLabelProps={props.noshrink ? { shrink: false } : null}
            label={props.noshrink ? ((cv || (props.location && props.location.display_name)) 
                  ? null : locale2.WHERE[props.lang]) : 
                  locale2.WHERE[props.lang]} />
        )
      }}
      freeSolo
    />)
}