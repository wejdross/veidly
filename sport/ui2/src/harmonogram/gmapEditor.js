import { Check, Close } from "@mui/icons-material";
import {
  Button, InputAdornment, TextField
} from "@mui/material";
import React, { useState } from 'react';
import { fetchGeo } from "../apicalls/user.api";
import MulwiMap from "../gmaps/maps";
import { locale2 } from "../locale";
import { MulwiColors } from "../mulwiColors";
import GoogleMaps from "../search/googlemaps";

 export function defaultLoc() {
    return ({
      LocationText: '',
      LocationLat: null,
      LocationLng: null,
    })
  }

var to = null

export function GmapEditor(props) {

  const [localizationUserInput, setLocalizationUserInput] = useState("")
  const [localizationData, setLocalizationData] = useState(defaultLoc())
  const [locOk, setLocOk] = useState(null)

  async function setGeo(c) {
    let resp = await fetchGeo(c)
    let tmp = { ...localizationData }
    tmp.LocationText = c
    let r = JSON.parse(resp)
    if (r.length === 0) {
      setLocOk(false)
      setLocalizationData(null)
    } else {
      tmp.LocationLat = parseFloat(r[0].lat)
      tmp.LocationLng = parseFloat(r[0].lon)
      setLocOk(true)
      setLocalizationData(tmp)
      if(!props.withConfirmButton && props.setLocalizationData)
        props.setLocalizationData(tmp)
    }
  }

  return (
    <React.Fragment>
      <GoogleMaps 
        setLocation={(e) => {
          setLocalizationData({
            LocationText: e.display_name,
            LocationLat: e.lat,
            LocationLng: e.lon,
          })
          props.setLocalizationData({
            LocationText: e.display_name,
            LocationLat: e.lat,
            LocationLng: e.lon,
          })
        }}
      
      />
      {/* <TextField
        variant="outlined"
        label={locale2.WHERE_WILL_BE_T[props.lang]}
        fullWidth
        error={Boolean(props.locErr)}
        helperText={props.locErr}
        value={localizationUserInput}
        InputProps={{
          endAdornment: locOk ? (
            <InputAdornment position="end">
              <Check style={{ color: MulwiColors.greenDark }} />
            </InputAdornment>
          ) : (
            <InputAdornment position="end">
              <Close style={{ color: MulwiColors.redError }} />
            </InputAdornment>
          )
        }}
        onChange={(e) => {
          let c = e.target.value
          setLocalizationUserInput(c)
          if (to) window.clearTimeout(to)
          to = setTimeout(() => {
            setGeo(c)
          }, 500)
        }} /> */}

      <div style={{ minHeight: 300, marginTop: 20, }}>
        <center>
          <MulwiMap center={localizationData.LocationText} height={300} />
        </center>
      </div>
      {
        props.withConfirmButton && (<Button 
            onClick={() => props.setLocalizationData(localizationData)}>
          {locale2.SAVE[props.lang]}
        </Button>)
      }
    </React.Fragment>
  )

}
