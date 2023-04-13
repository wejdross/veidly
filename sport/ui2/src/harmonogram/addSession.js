// import DateFnsUtils from "@date-io/date-fns";
// import {
//     AppBar,
//     ButtonBase,
//     Checkbox, Chip, DialogContent, FormControl, FormControlLabel,
//     Grid, IconButton, InputAdornment, InputLabel,
//     MenuItem, Select, Tab, Tabs, TextField, Typography
// } from "@mui/material";
// import { Add, Close } from "@mui/icons-material";
// import { KeyboardDatePicker, KeyboardTimePicker, MuiPickersUtilsProvider } from "@mui/lab/";
// import React, { useState } from "react";
// import { postTrainingSession, postTrainingSessionInBatch } from "../apicalls/instructor.api";
// import { MulwiColors } from "../mulwiColors";
// import ModalEdit from "../profile/ModalEdit";
// import { TabPanel } from "../tabPanel";

// export function AddSession(props) {

//     const [occDate, setOccDate] = useState(new Date())
//     const [dur, setdur] = useState(1)
//     const [weekly, setWeekly] = useState(false)
//     const [remarks, setRemarks] = useState("")
//     const [c, setc] = useState(MulwiColors.pinkDark)
//     const [colorOpen, setColorOpen] = React.useState(false)

//     const [days, setDays] = useState([false, false, false, false, false, false, false])
//     const labels = ["PN", "WT", "ŚR", "CZ", "PT", "SO", "ND"]
//     const [hrs, setHrs] = useState([new Date()])

//     function setHrsIx(ix, v) {
//         let cpy = [...hrs]
//         cpy[ix] = v
//         setHrs(cpy)
//     }

//     function setDaysIx(ix) {
//         let cpy = [...days]
//         cpy[ix] = !cpy[ix]
//         setDays(cpy)
//     }

//     const [value, setValue] = React.useState(0);
//     const handleChange = (event, newValue) => {
//         setValue(newValue);
//     };

//     return (
//         <React.Fragment>
//             <ModalEdit
//                 nocontent
//                 open={props.open}
//                 onlyButton
//                 buttonProps={{
//                     style: {
//                         backgroundColor: MulwiColors.greenDark,
//                         color: "white"
//                     }
//                 }}
//                 title="Występowanie"
//                 label="Dodaj występowanie treningu"
//                 onSave={async () => {

//                     if (!props.trainingID) {
//                         throw new Error("Nieprawidłowe parametry")
//                     }

//                     switch (value) {
//                         case 0:
//                             let now = new Date()
//                             let today = (now.getDay() + 6) % 7
//                             if (hrs.length == 0) {
//                                 throw new Error("Nie podałeś godzin występowania")
//                             }
//                             // if(days.length == 0) {
//                             //     throw new Error("Nie dni występowania")
//                             // }

//                             let cc = {
//                                 Color: c,
//                                 Remarks: remarks,
//                                 RepeatDays: 7
//                             }

//                             let req = []

//                             for (let i = 0; i < days.length; i++) {
//                                 if (!days[i]) continue
//                                 for (let j = 0; j < hrs.length; j++) {
//                                     let ccp = { ...cc }
//                                     let start = new Date()
//                                     start.setHours(hrs[j].getHours(), hrs[j].getMinutes(), 0, 0)
//                                     let df = i - today
//                                     if (df < 0) df += 7
//                                     start.setDate(start.getDate() + df)
//                                     let end = new Date(start)
//                                     let min = Math.round((dur % 1) * 60)
//                                     end.setHours(
//                                         end.getHours() + Math.floor(dur),
//                                         end.getMinutes() + min)
//                                     ccp.DateEnd = end
//                                     ccp.DateStart = start
//                                     req.push(ccp)
//                                 }
//                             }

//                             if (req.length == 0) {
//                                 throw new Error("Nie ma występowań")
//                             }

//                             await postTrainingSessionInBatch({
//                                 Occs: req,
//                                 TrainingID: props.trainingID,
//                             })
//                             props.onChange()
//                             break
//                         case 1:
//                             let end = new Date(occDate)
//                             end.setHours(end.getHours() + dur)
//                             let newSess = {
//                                 DateStart: occDate,
//                                 DateEnd: end,
//                                 RepeatDays: weekly ? 7 : 0,
//                                 Color: c,
//                                 Remarks: remarks,
//                                 TrainingID: props.trainingID
//                             }
//                             await postTrainingSession(newSess)
//                             props.onChange()
//                             break
//                         default:
//                             throw new Error("Invalid mode")
//                     }

//                 }}
//                 content={
//                     <React.Fragment>
//                         <AppBar position="static" style={{
//                             backgroundColor: MulwiColors.whiteBackground,
//                             color: "black"
//                         }}>
//                             <Tabs TabIndicatorProps={{ style: { backgroundColor: MulwiColors.blueDark } }} value={value} onChange={handleChange}>
//                                 <Tab label="Treningi tygodniowe" id="stt-1" />
//                                 <Tab label="Inne" id="stt-2" />
//                             </Tabs>
//                         </AppBar>
//                         <DialogContent>
//                             <TabPanel value={value} index={0}>
//                                 <br />
//                                 <Typography style={{
//                                     marginBottom: 10
//                                 }}>W jakie dni ma trening występować?</Typography>
//                                 <Grid container direction="row" spacing={1}>
//                                     {labels.map((l, i) => (<Grid key={i} item>
//                                         <ButtonBase style={{
//                                             borderRadius: 15,
//                                         }} onClick={() => setDaysIx(i)} >
//                                             <Chip key={l} label={l} style={{
//                                                 backgroundColor: days[i] ? MulwiColors.blueLight : null,
//                                                 color: days[i] ? "white" : null,
//                                                 cursor: "pointer"
//                                             }} />
//                                         </ButtonBase>
//                                     </Grid>))}
//                                 </Grid>
//                                 <br />
//                                 <TextField
//                                     style={{
//                                         marginBottom: 15
//                                     }}
//                                     fullWidth
//                                     size="small"
//                                     variant="outlined"
//                                     label="długość trwania" type="number"
//                                     InputProps={{
//                                         startAdornment: (
//                                             <InputAdornment style={{ marginRight: 10 }} position="start">hr</InputAdornment>
//                                         )
//                                     }}
//                                     value={String(Number(dur))}
//                                     onChange={e => setdur(Number(e.target.value))} />
//                                 <Typography>
//                                     Godziny występowania treningu w podane dni
//                             </Typography>
//                                 <Grid container direction="row" spacing={2} alignItems="center">
//                                     {hrs.map((h, i) => (<Grid key={i} item><KeyboardTimePicker
//                                         ampm={false}
//                                         margin="normal"
//                                         style={{
//                                             marginBottom: 20,
//                                             width: 80
//                                         }}
//                                         value={h}
//                                         onChange={v => {
//                                             if (v instanceof Date && isFinite(v)) {
//                                                 let d = new Date(h)
//                                                 d.setHours(v.getHours())
//                                                 d.setMinutes(v.getMinutes())
//                                                 d.setSeconds(0)
//                                                 d.setMilliseconds(0)
//                                                 setHrsIx(i, d)
//                                                 //setOccDate(d)
//                                             }
//                                         }}
//                                         KeyboardButtonProps={{
//                                             'aria-label': 'change time',
//                                         }}
//                                         keyboardIcon={null}
//                                         KeyboardButtonProps={{ size: "small" }}
//                                         InputProps={{
//                                             startAdornment: <IconButton style={{
//                                                 width: 20,
//                                                 height: 20,
//                                                 marginRight: 5,
//                                             }} size="small" onClick={() => {
//                                                 let cpy = []
//                                                 for (let j = 0; j < hrs.length; j++) {
//                                                     if (j == i) continue
//                                                     cpy.push(hrs[j])
//                                                     setHrs(cpy)
//                                                 }
//                                             }} >
//                                                 <Close />
//                                             </IconButton>
//                                         }}
//                                     /></Grid>))}
//                                     <IconButton onClick={() => {
//                                         let cpy = [...hrs]
//                                         cpy.push(new Date())
//                                         setHrs(cpy)
//                                     }}>
//                                         <Add />
//                                     </IconButton>
//                                 </Grid>
//                                 <br />
//                                 <TextField
//                                     multiline
//                                     inputProps={
//                                         { maxLength: 250 }
//                                     }
//                                     style={{ marginTop: 10 }}
//                                     rows={5}
//                                     variant="outlined"
//                                     helperText={remarks.length + "/250"}
//                                     fullWidth={true}
//                                     value={remarks}
//                                     onChange={e => setRemarks(e.target.value)}
//                                     label="uwagi do tych sesji"
//                                 />
//                                 <FormControl fullWidth style={{ marginTop: 5 }}>
//                                     <InputLabel id="demo-controlled-open-select-label">Wyróżnij sesje kolorem</InputLabel>
//                                     <Select
//                                         labelId="demo-controlled-open-select-label"
//                                         id="demo-controlled-open-select"
//                                         open={colorOpen}
//                                         onClose={() => setColorOpen(false)}
//                                         onOpen={() => setColorOpen(true)}
//                                         value={c}
//                                         fullWidth
//                                         onChange={e => setc(e.target.value)}>
//                                         {/* <MenuItem value="">
//                                     <em>None</em>
//                                 </MenuItem> */}
//                                         <MenuItem value={MulwiColors.pinkDark}><div style={{ backgroundColor: MulwiColors.pinkDark, minHeight: 30, width: "100%" }} /></MenuItem>
//                                         <MenuItem value={MulwiColors.pinkAction}><div style={{ backgroundColor: MulwiColors.pinkAction, minHeight: 30, width: "100%" }} /></MenuItem>
//                                         <MenuItem value={MulwiColors.greenDark}><div style={{ backgroundColor: MulwiColors.greenDark, minHeight: 30, width: "100%" }} /></MenuItem>
//                                         <MenuItem value={MulwiColors.greenLight}><div style={{ backgroundColor: MulwiColors.greenLight, minHeight: 30, width: "100%" }} /></MenuItem>
//                                         <MenuItem value={MulwiColors.blueDark}><div style={{ backgroundColor: MulwiColors.blueDark, minHeight: 30, width: "100%" }} /></MenuItem>
//                                         <MenuItem value={MulwiColors.blueLight}><div style={{ backgroundColor: MulwiColors.blueLight, minHeight: 30, width: "100%" }} /></MenuItem>
//                                     </Select>
//                                 </FormControl>
//                             </TabPanel>
//                             <TabPanel value={value} index={1}>
//                                 <React.Fragment>
//                                     <Grid container direction="column" spacing={1}>
//                                         <Grid item>
//                                             <KeyboardDatePicker fullWidth
//                                                 margin="normal"
//                                                 id="date-picker-dialog"
//                                                 label="Data rozpoczęcia"
//                                                 format="MM/dd/yyyy"
//                                                 value={occDate}
//                                                 onChange={v => {
//                                                     if (v instanceof Date && isFinite(v)) {
//                                                         let d = new Date(occDate)
//                                                         d.setMonth(v.getMonth())
//                                                         d.setDate(v.getDate())
//                                                         d.setFullYear(v.getFullYear())
//                                                         d.setSeconds(0)
//                                                         d.setMilliseconds(0)
//                                                         setOccDate(d)
//                                                     }
//                                                 }}
//                                                 KeyboardButtonProps={{
//                                                     'aria-label': 'change date',
//                                                 }} />
//                                         </Grid>
//                                         <Grid item>
//                                             <TextField
//                                                 fullWidth
//                                                 size="small"
//                                                 variant="outlined"
//                                                 label="długość trwania" type="number"
//                                                 InputProps={{
//                                                     startAdornment: (
//                                                         <InputAdornment style={{ marginRight: 10 }} position="start">hr</InputAdornment>
//                                                     )
//                                                 }}
//                                                 value={String(Number(dur))}
//                                                 onChange={e => setdur(Number(e.target.value))} />
//                                         </Grid>
//                                         <Grid item>
//                                             <KeyboardTimePicker
//                                                 fullWidth
//                                                 ampm={false}
//                                                 margin="normal"
//                                                 id="time-picker"
//                                                 label="Godzina występowania"
//                                                 value={occDate}
//                                                 onChange={v => {
//                                                     if (v instanceof Date && isFinite(v)) {
//                                                         let d = new Date(occDate)
//                                                         d.setHours(v.getHours())
//                                                         d.setMinutes(v.getMinutes())
//                                                         d.setSeconds(0)
//                                                         d.setMilliseconds(0)
//                                                         setOccDate(d)
//                                                     }
//                                                 }}
//                                                 KeyboardButtonProps={{
//                                                     'aria-label': 'change time',
//                                                 }}
//                                             />
//                                         </Grid>
//                                         <Grid item>
//                                             <FormControl>
//                                                 <FormControlLabel
//                                                     control={<Checkbox checked={weekly}
//                                                         onChange={e => setWeekly(e.target.checked)} />}
//                                                     label="Powtarzaj co tydzień"
//                                                 />
//                                             </FormControl>
//                                         </Grid>
//                                     </Grid>
//                                     <TextField
//                                         multiline
//                                         inputProps={
//                                             { maxLength: 250 }
//                                         }
//                                         style={{ marginTop: 10 }}
//                                         rows={5}
//                                         variant="outlined"
//                                         helperText={remarks.length + "/250"}
//                                         fullWidth={true}
//                                         value={remarks}
//                                         onChange={e => setRemarks(e.target.value)}
//                                         label="uwagi do treningu"
//                                     />
//                                     <FormControl fullWidth style={{ marginTop: 5 }}>
//                                         <InputLabel id="demo-controlled-open-select-label">Wyróżnij sesję kolorem</InputLabel>
//                                         <Select
//                                             labelId="demo-controlled-open-select-label"
//                                             id="demo-controlled-open-select"
//                                             open={colorOpen}
//                                             onClose={() => setColorOpen(false)}
//                                             onOpen={() => setColorOpen(true)}
//                                             value={c}
//                                             fullWidth
//                                             onChange={e => setc(e.target.value)}>
//                                             <MenuItem value={MulwiColors.pinkDark}><div style={{ backgroundColor: MulwiColors.pinkDark, minHeight: 30, width: "100%" }} /></MenuItem>
//                                             <MenuItem value={MulwiColors.pinkAction}><div style={{ backgroundColor: MulwiColors.pinkAction, minHeight: 30, width: "100%" }} /></MenuItem>
//                                             <MenuItem value={MulwiColors.greenDark}><div style={{ backgroundColor: MulwiColors.greenDark, minHeight: 30, width: "100%" }} /></MenuItem>
//                                             <MenuItem value={MulwiColors.greenLight}><div style={{ backgroundColor: MulwiColors.greenLight, minHeight: 30, width: "100%" }} /></MenuItem>
//                                             <MenuItem value={MulwiColors.blueDark}><div style={{ backgroundColor: MulwiColors.blueDark, minHeight: 30, width: "100%" }} /></MenuItem>
//                                             <MenuItem value={MulwiColors.blueLight}><div style={{ backgroundColor: MulwiColors.blueLight, minHeight: 30, width: "100%" }} /></MenuItem>
//                                         </Select>
//                                     </FormControl>
//                                 </React.Fragment>
//                             </TabPanel>
//                         </DialogContent>
//                     </React.Fragment>
//                 }
//             />
//         </React.Fragment>
//     )
// }