import {
  Button, Checkbox, Container, Dialog, DialogContent, DialogTitle, FormControl, FormControlLabel, Grid, InputAdornment, MenuItem,
  TextField, Typography
} from "@mui/material";
import React, { useEffect, useState } from 'react';
import { apiCreateTraining } from "../apicalls/instructor.api";
import { dateToEpoch, prettyPrintCurrency } from "../helpers";
import { locale2 } from "../locale";
import { MulwiColors } from "../mulwiColors";
import { getErrorDialog, getNullDialog, StatusDialog } from "../StatusDialog";
import { defaultLoc, GmapEditor } from "./gmapEditor";
import GoogleMaps from '../search/googlemaps';
import CardWithBg from "../card/cardWithBg";
import { putOcc } from "../apicalls/occ";
import { DateTimePicker, LocalizationProvider } from "@mui/x-date-pickers";

export const minPrice = 5

export function AddTraining(props) {

  //const [open, _setOpen] = useState(false)
  const [name, setName] = useState("")
  const [nameErr, setNameErr] = useState("")

  const currencies = [
    {
      value: "PLN",
      label: prettyPrintCurrency("PLN"),
    },
  ]

  const [price, setPrice] = useState(minPrice)
  const [priceErr, setPriceErr] = useState("")

  const [currency, setCurrency] = useState(currencies[0].value)

  const [capacity, setCapacity] = useState(1)

  const [locErr, setLocErr] = useState("")

  const [info, setInfo] = useState(getNullDialog())

  async function save() {
    if (!location.lat || !location.lon) return
    let t = {
      Title: name,
      Capacity: capacity,
      LocationText: location.display_name,
      LocationLat: location.lat,
      LocationLng: location.lon,
      LocationCountry: '',
      Price: Math.round(price * 100),
      Currency: currency,
    }
    try {
      let id = await apiCreateTraining({
        Training: t,
        ReturnID: true
      })
      if (!id)
        throw new Error(locale2.COULDNT_VERIFY_TID[props.lang])
      let occs
      if (props.occData) {
        occs = [
          {
            DateStart: dateStart,
            DateEnd: dateEnd,
            RepeatDays: repeatDays
          }
        ]
        
        let req = {
          TrainingID: id,
          Occurrences: occs 
        }
        await putOcc(req)
      }
      t.ID = id
      props.setDrawerData({
        training: t,
        occs: occs
        //openAddSession: true
      })
      if (props.onChange)
        props.onChange()
      //setOpen(false)
      props.openDrawer()
      close(true)
    } catch (ex) {
      setInfo(getErrorDialog(ex))
    }
  }

  function close(success) {
    if (props.onClose)
      props.onClose(success)
    if (props.setOpen)
      props.setOpen(false)
  }

  function validate() {
    let e = false

    if (name == "") {
      if (!nameErr) {
        setNameErr(locale2.T_NAME_REQUIRED[props.lang])
      }
      e = true
    } else {
      if (nameErr) {
        setNameErr("")
      }
    }

    if (location.lat == null || location.lon == null || !location.display_name) {
      if (!locErr) {
        setLocErr(locale2.LOC_REQUIRED[props.lang])
      }
      e = true
    } else {
      if (locErr) {
        setLocErr("")
      }
    }

    if (price < minPrice) {
      if (!priceErr) {
        setPriceErr(locale2.PRICE_VAL[props.lang] + ' ' + minPrice
          + ' ' + prettyPrintCurrency(currency))
      }
      e = true
    } else {
      if (priceErr) {
        setPriceErr("")
      }
    }

    return e
  }

  const [repeatTraining, setRepeatTraining] = useState(false)
  const [repeatDays, setRepeatDays] = useState(0)
  const [dateStart, setDateStart] = useState(null)
  const [dateEnd, setDateEnd] = useState(null)
  const [location, setLocation] = useState({})

  useEffect(() => {
    let d = props.occData
    if (!d)
      return
    setDateStart(d.start)
    setDateEnd(d.end)
  }, [props.occData])

  useEffect(validate)

  function form() {
    return (<React.Fragment>
      <Grid item style={{ width: "100%" }}>
        <Grid container direction="row" spacing={2}>
          <Grid item xs={6}>
            <TextField
              label={locale2.NAME[props.lang]}
              variant="outlined"
              fullWidth
              inputProps={
                { maxLength: 128 }
              }
              error={Boolean(nameErr)}
              helperText={nameErr || (((name && name.length) || 0) + "/128")}
              value={name}
              style={{
                marginBottom: 10
              }}
              onChange={(e) => setName(e.target.value)} />
          </Grid>
          <Grid item xs={6}>
            <GoogleMaps
              lang={props.lang}
              class={"tagLocation"}
              location={location}
              fullWidth
              error={Boolean(locErr)}
              errorText={locErr}
              label={locale2.LOCATION[props.lang]}
              setLocation={e => {
                let lat = Number(e.lat)
                let lng = Number(e.lon)
                if (lat && lng) {
                  setLocation(e)
                }
              }}
            />
          </Grid>
        </Grid>
      </Grid>
      {/* <Grid item style={{ width: "100%" }}>
        <GmapEditor lang={props.lang} setLocalizationData={setLocalizationData} />
      </Grid> */}

      <Grid item style={{ width: "100%", marginTop: 15 }}>
        <Grid container direction="row" spacing={2}>

          <Grid item xs={2}>
            <TextField
              id="standard-select-currency"
              select fullWidth
              label={locale2.CURRENCY[props.lang]}
              value={currency}
              variant="outlined"
              onChange={(e) => setCurrency(e.target.value)}>
              {currencies.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.value}
                </MenuItem>
              ))}
            </TextField>
          </Grid>

          <Grid item xs={6}>
            <FormControl fullWidth variant="outlined">
              <TextField
                type="number"
                value={String(price)}
                error={Boolean(priceErr)}
                helperText={priceErr}
                label={locale2.PRICE[props.lang]}
                onChange={e => {
                  let c = Number(e.target.value)
                  if (!isNaN(c)) {
                    setPrice(c)
                  }
                }}
                variant="outlined"
                InputProps={{
                  startAdornment: <InputAdornment position="start">{currency}</InputAdornment>
                }}
              />
            </FormControl>
          </Grid>

          <Grid item xs={4}>
            <TextField
              label={locale2.MAX_CAPACITY[props.lang]}
              onChange={(e) => setCapacity(Number(e.target.value))}
              variant="outlined" type="number" value={String(capacity)}></TextField>
          </Grid>

        </Grid>
      </Grid>

      {props.occData && (<React.Fragment>
        <Grid item xs={12} style={{ marginTop: 15 }}>
          <Grid container direction="row" spacing={2}>
            <Grid item xs={6}>
              <DateTimePicker
                value={dateStart}
                onChange={e => {
                  if (e instanceof Date && isFinite(e))
                    setDateStart(e)
                }}
                ampm={false}
                renderInput={(params) => <TextField {...params} />}
                label={locale2.START_DATE[props.lang]} />
            </Grid>
            <Grid item xs={6}>
              <DateTimePicker
                value={dateEnd}
                onChange={e => {
                  if (e instanceof Date && isFinite(e))
                    setDateEnd(e)
                }}
                ampm={false}
                renderInput={(params) => <TextField {...params} />}
                label={locale2.END_DATE[props.lang]} />
            </Grid>
          </Grid>
        </Grid>

        <FormControlLabel
          control={<Checkbox checked={repeatTraining}
            onChange={e => setRepeatTraining(e.target.checked)} />}
          label="Powtarzaj trening" />

      </React.Fragment>)}


      {repeatTraining && (<React.Fragment>
        <Grid item xs={12}>
          <TextField label="powtarzaj trening co"
            value={repeatDays} type="number"
            onChange={e => setRepeatDays(Math.abs(Number(e.target.value)))}
            InputProps={{
              endAdornment: <InputAdornment position="end">dni</InputAdornment>
            }} />
        </Grid>
      </React.Fragment>)}

      <Grid item style={{ marginTop: 15 }} xs={12}>
        <Grid container spacing={2} direction={'row'} alignContent={'space-between'} justifyContent={'space-between'}>
          <Grid item>
            <Button style={{
              color: "white",
              backgroundColor: MulwiColors.redError
            }} onClick={() => {
              close(false)
            }}>{locale2.CANCEL[props.lang]}</Button>
          </Grid>
          <Grid item>
            <Button style={{
              color: "white",
              backgroundColor: MulwiColors.greenDark
            }} variant="contained" onClick={save}>{locale2.NEXT[props.lang]}</Button>
          </Grid>
        </Grid>
      </Grid>
    </React.Fragment>)
  }

  return (props.modal ? (<React.Fragment>
    <Dialog open={props.open} onClose={() => props.setOpen(false)}>
      <DialogTitle>
        {locale2.NEW_TRAINING[props.lang]}
      </DialogTitle>
      <DialogContent>
        {form()}
      </DialogContent>
    </Dialog>
  </React.Fragment>) : (<CardWithBg header={locale2.NEW_TRAINING[props.lang]} >
    <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
    {form()}
  </CardWithBg>)
  )
}