// import DateFnsUtils from "@date-io/date-fns";
// import {
//   Button, Card, CircularProgress, Dialog, DialogActions,
//   DialogContent, DialogTitle, FormControl, Grid, InputAdornment, MenuItem, Slider, TextField, Typography, useMediaQuery, useTheme
// } from "@mui/material";
// import { Add } from "@mui/icons-material";
// import { KeyboardDatePicker, MuiPickersUtilsProvider } from "@mui/lab/";
// import React, { useEffect, useState } from 'react';
// import { apiCreateTraining } from "../apicalls/instructor.api";
// import { fetchGeo } from "../apicalls/user.api";
// import MulwiMap from "../gmaps/maps";
// import { MulwiColors } from "../mulwiColors";
// import { errToStr } from "../StatusDialog";
// import { ArrEdit } from "./arrEdit";
// import { DiffAtc } from "./DiffAtc";
// import { TagAtc } from "./TagAtc";

// export function AddTraining(props) {

//   const [open, _setOpen] = useState(false)
//   const [diff, setDiff] = useState([])
//   const [tags, setTags] = useState([])
//   const [name, setName] = useState("")
//   const [nameErr, setNameErr] = useState("")
//   const [desc, setDesc] = useState("")

//   const [req, setReq] = useState([])
//   const [reqGear, setReqGear] = useState([])
//   const [recGear, setRecGear] = useState([])
//   const [instrGear, setInstrGear] = useState([])
//   const [minAge, setMinAge] = useState(0)
//   const [maxAge, setMaxAge] = useState(0)

//   const theme = useTheme()
//   const isLowRes = useMediaQuery(theme.breakpoints.down('xs'))

//   const currencies = [
//     // {
//     //   value: 'USD',
//     //   label: '$',
//     // },
//     {
//       value: 'PLN',
//       label: 'zł',
//     },
//     // {
//     //   value: 'EUR',
//     //   label: '€',
//     // },
//     // :(
//     // {
//     //   value: 'BTC',
//     //   label: '฿',
//     // },
//     // {
//     //   value: 'JPY',
//     //   label: '¥',
//     // },
//   ]

//   const [price, setPrice] = useState(1)
//   const [priceErr, setPriceErr] = useState("")

//   const [currency, setCurrency] = useState("PLN")

//   const [capacity, setCapacity] = useState(1)

//   const [dateStart, setDateStart] = useState(null)
//   const [dateEnd, setDateEnd] = useState(null)

//   const [st, setSt] = useState(0)
//   const [msg, setMsg] = useState("")

//   function defaultLoc() {
//     return ({
//       LocationText: '',
//       LocationLat: null,
//       LocationLng: null,
//     })
//   }

//   // localisation purpose
//   const [unsufficientData, setUnsufficientData] = useState(false)
//   const [localizationUserInput, setLocalizationUserInput] = useState("")
//   const [localizationData, setLocalizationData] = useState(defaultLoc())
//   const [locErr, setLocErr] = useState("")

//   function setOpen(x) {
//     if (x) {
//       setMsg("")
//       setSt(0)
//     }
//     _setOpen(x)
//   }

//   async function save() {
//     let t = {
//       Title: name,
//       Description: desc,
//       DateStart: dateStart,
//       DateEnd: dateEnd,
//       Capacity: capacity,
//       //Requirements: '',
//       LocationText: localizationData.LocationText,
//       LocationLat: localizationData.LocationLat,
//       LocationLng: localizationData.LocationLng,
//       LocationCountry: '',
//       Price: Math.round(price * 100),
//       Currency: currency,
//       Tags: tags,
//       Diff: diff,
//       RequiredGear: reqGear,
//       RecommendedGear: recGear,
//       InstructorGear: instrGear,
//       MinAge: minAge,
//       MaxAge: maxAge,
//     }
//     setSt(1)
//     try {
//       let id = await apiCreateTraining({
//         Training: t,
//         ReturnID: true
//       })
//       if (!id)
//         throw new Error("nie można było zweryfikować id treningu")
//       t.ID = id
//       props.setDrawerData({
//         training: t,
//         openAddSession: true
//       })
//       props.onChange()
//       setOpen(false)
//       props.openDrawer()
//     } catch (ex) {
//       setMsg(ex)
//       setSt(2)
//     }
//   }

//   function validate() {
//     let e = false

//     if (name == "") {
//       if (!nameErr) {
//         setNameErr("nazwa treningu jest wymagana")
//       }
//       e = true
//     } else {
//       if (nameErr) {
//         setNameErr("")
//       }
//     }

//     if (localizationData.LocationLat == null ||
//       localizationData.LocationLng == null ||
//       !localizationData.LocationText) {
//       if (!locErr) {
//         setLocErr("Lokacja jest wymagana")
//       }
//       e = true
//     } else {
//       if (locErr) {
//         setLocErr("")
//       }
//     }

//     if (price < 1) {
//       if (!priceErr) {
//         setPriceErr("Cena musi być więszka niż 1zł")
//       }
//       e = true
//     } else {
//       if (priceErr) {
//         setPriceErr("")
//       }
//     }

//     return e
//   }

//   useEffect(() => {
//     validate()
//   })

//   function wip() {
//     if (st !== 1) {
//       return null
//     }
//     return (
//       <React.Fragment> <DialogContent>
//         <CircularProgress style={{
//           marginLeft: 100, marginRight: 100, marginTop: 50, marginBottom: 50
//         }} />
//       </DialogContent></React.Fragment>)
//   }

//   function fin() {
//     if (st !== 2) {
//       return null
//     }
//     return (<React.Fragment>
//       <DialogContent>
//         <Typography>{errToStr(msg)}</Typography>
//       </DialogContent>
//       <DialogActions>
//         <Button onClick={() => {
//           setMsg("")
//           setSt(0)
//         }} color="primary">
//           Jeszcze raz
//         </Button>
//         <Button onClick={() => setOpen(false)} color="secondary">
//           Zamknij
//         </Button>
//       </DialogActions>
//     </React.Fragment>)
//   }

//   const [sm, setsm] = useState(false)

//   function editor() {
//     if (st !== 0) {
//       return null
//     }
//     return (<React.Fragment>
//       <form onSubmit={(e) => {
//         e.preventDefault()
//         save()
//       }}>
//         <DialogContent style={{
//           width: isLowRes ? "" : 500
//         }} >
//           <Grid container spacing={2}>
//             <Grid item style={{ width: "100%" }}>
//               <TextField
//                 label="nazwa"
//                 fullWidth
//                 inputProps={
//                   { maxLength: 128 }
//                 }

//                 error={Boolean(nameErr)}
//                 helperText={nameErr || (((name && name.length) || 0) + "/128")}
//                 value={name}
//                 onChange={(e) => setName(e.target.value)} />
//             </Grid>
//             <Grid item style={{ width: "100%" }}>
//               <TextField
//                 label="Gdzie odbędzie się trening?"
//                 fullWidth
//                 error={Boolean(locErr)}
//                 helperText={locErr}
//                 onBlur={() => setsm(false)}
//                 onFocus={() => setsm(true)}
//                 value={localizationUserInput}
//                 onChange={async (e) => {
//                   //setLocalText(e.target.value);
//                   let c = e.target.value
//                   setLocalizationUserInput(c)
//                   let resp = await fetchGeo(c)
//                   let tmp = { ...localizationData }
//                   tmp.LocationText = c
//                   let r = JSON.parse(resp)
//                   if (r.length === 0) {
//                     setUnsufficientData(true)
//                   } else {
//                     tmp.LocationLat = parseFloat(r[0].lat)
//                     tmp.LocationLng = parseFloat(r[0].lon)
//                     // tmp.LocationText = JSON.parse(resp)[0].display_name
//                     setUnsufficientData(false)
//                   }
//                   setLocalizationData(tmp)
//                 }} />
//               <div style={{
//                 visibility: sm ? "visible" : "hidden",
//                 zIndex: 9999,
//                 position: "relative"
//               }}>
//                 <Card style={{
//                   width: "100%",
//                   position: "absolute",
//                   minHeight: 300
//                 }}>
//                   {
//                     unsufficientData && <Typography>Musisz wprowadzić dokładniejszy adres!</Typography> || <center><MulwiMap center={localizationData.LocationText} /></center>
//                   }
//                 </Card>
//               </div>
//             </Grid>
//             <Grid item style={{ width: "100%" }}>
//               <TextField
//                 multiline
//                 rows={3}
//                 variant={"outlined"}
//                 helperText={desc.length + "/250"}
//                 fullWidth={true}
//                 value={desc}
//                 onChange={(event => (setDesc(event.target.value)))}
//                 label="opis"
//                 inputProps={
//                   { maxLength: 250 }
//                 } />
//             </Grid>
//             <Grid item style={{ width: "100%" }}>
//               <DiffAtc setDiff={setDiff} diff={diff} />
//             </Grid>
//             <Grid item style={{ width: "100%" }}>
//               <TagAtc tags={tags} setTags={setTags} />
//             </Grid>
//             <Grid item style={{ width: "100%" }}>
//               <Grid container spacing={2} >
//                 <Grid item>
//                   <FormControl fullWidth variant="outlined">
//                     <TextField
//                       type="number"
//                       value={String(price)}
//                       error={Boolean(priceErr)}
//                       helperText={priceErr}
//                       label={"Cena"}
//                       onChange={e => {
//                         let c = Number(e.target.value)
//                         if (!isNaN(c)) {
//                           setPrice(c)
//                         }
//                       }}
//                       variant="outlined"
//                       InputProps={{
//                         startAdornment: <InputAdornment position="start">{currency}</InputAdornment>
//                       }}
//                     />
//                   </FormControl>
//                   {/* <TextField label="Cena" fullWidth style={{marginBottom: 10}} /> */}
//                 </Grid>
//                 <Grid item>
//                   <TextField
//                     id="standard-select-currency"
//                     select
//                     label="Waluta"
//                     value={currency}
//                     onChange={(e) => setCurrency(e.target.value)}
//                   >
//                     {currencies.map((option) => (
//                       <MenuItem key={option.value} value={option.value}>
//                         {option.label} / {option.value}
//                       </MenuItem>
//                     ))}
//                   </TextField>
//                 </Grid>
//               </Grid>
//             </Grid>
//             {/* <Grid item style={{ width: "100%" }}>
//               <Grid container justify="center" alignItems="center" spacing={2}>
//                 <Grid item>
//                   <Typography id="continuous-slider" gutterBottom>
//                     Max ilość osób
//                 </Typography>
//                 </Grid>
//                 <Grid item>
//                   <TextField
//                     onChange={(e) => setCapacity(Number(e.target.value))}
//                     variant="outlined" type="number" value={String(capacity)}
//                     size="small"></TextField>
//                 </Grid>
//               </Grid>
//               <Slider
//                 min={1}
//                 label=""
//                 value={capacity} onChange={(e, v) => setCapacity(v)} />
//             </Grid> */}
//             <Grid item style={{ width: "100%" }}>
//               <TextField
//                 label="Max. ilość osób"
//                 onChange={(e) => setCapacity(Number(e.target.value))}
//                 variant="outlined" type="number" value={String(capacity)}></TextField>
//             </Grid>

//             <Grid item style={{ width: "100%" }}>
//               <Grid container spacing={2} direction="row" >

//                 <Grid item xs={6}>
//                   <TextField
//                     label="Min. wiek"
//                     onChange={(e) => setMinAge(Number(e.target.value))}
//                     variant="outlined" type="number" value={minAge ? String(minAge) : ""} />
//                 </Grid>

//                 <Grid item xs={6}>
//                   <TextField
//                     label="Max. wiek"
//                     onChange={(e) => setMaxAge(Number(e.target.value))}
//                     variant="outlined" type="number" value={maxAge ? String(maxAge) : ""} />
//                 </Grid>
//               </Grid>

//             </Grid>

//             <Grid item style={{ width: "100%" }}>
//               <ArrEdit value={reqGear} setValue={setReqGear} label="Sprzęt który klient musi mieć na treningu" />
//             </Grid>
//             <Grid item style={{ width: "100%" }}>
//               <ArrEdit value={recGear} setValue={setRecGear} label="Sprzęt który rekomendujesz by klient miał" />
//             </Grid>
//             <Grid item style={{ width: "100%" }}>
//               <ArrEdit value={instrGear} setValue={setInstrGear} label="Sprzęt który ty posiadasz" />
//             </Grid>
//             {/* <Grid item style={{ width: "100%" }}>
//               <Grid container spacing={2} style={{ marginBottom: 20 }}>
//                 <Grid item sm={12} md={6}>
//                   <KeyboardDatePicker fullWidth
//                     margin="normal"
//                     id="date-picker-dialog"
//                     label="Data rozpoczęcia"
//                     format="MM/dd/yyyy"
//                     value={dateStart}
//                     minDate={new Date()}
//                     onChange={e => {
//                       let x = new Date(e)
//                       x.setHours(0, 0, 0, 0)
//                       setDateStart(x)
//                     }}
//                     KeyboardButtonProps={{
//                       'aria-label': 'change date',
//                     }} />
//                 </Grid>
//                 <Grid item sm={12} md={6}>
//                   <KeyboardDatePicker fullWidth
//                     margin="normal"
//                     id="date-picker-dialog"
//                     label="Data zakończenia"
//                     format="MM/dd/yyyy"
//                     value={dateEnd}
//                     minDate={dateStart || new Date()}
//                     onChange={e => {
//                       let x = new Date(e)
//                       x.setHours(23, 59, 59, 0)
//                       setDateEnd(x)
//                     }}
//                     KeyboardButtonProps={{
//                       'aria-label': 'change date',
//                     }} />
//                 </Grid>
//               </Grid>
//             </Grid> */}
//           </Grid>
//         </DialogContent>
//         <DialogActions>
//           <Button onClick={() => setOpen(false)} color="secondary">
//             Anuluj
//         </Button>
//           <Button type="submit" color="primary">
//             Dalej
//         </Button>
//         </DialogActions>
//       </form>
//     </React.Fragment>)
//   }

//   useEffect(() => {
//     let query = new URLSearchParams(window.location.search)
//     if (query.get("open_add_train")) {
//       setOpen(true)
//     }
//   }, [])

//   return (<React.Fragment>
//     <Button variant="contained"
//       onClick={() => setOpen(true)}
//       style={{
//         backgroundColor: MulwiColors.greenDark,
//         color: "white"
//       }}><Add /> Dodaj nowy trening
//     </Button>
//     <Dialog open={open} onClose={() => setOpen(false)}
//       aria-labelledby="form-dialog-title">
//       <DialogTitle id="form-dialog-title">Dodaj nowy trening</DialogTitle>
//       <React.Fragment>
//         {editor()}
//         {wip()}
//         {fin()}
//       </React.Fragment>
//     </Dialog>
//   </React.Fragment>)
// }