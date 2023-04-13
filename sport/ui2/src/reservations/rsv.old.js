// import React, { useEffect, useState } from 'react'
// import { useHistory } from 'react-router'
// import { HarmonogramWeek } from '../harmonogram/weekBigRes'
// import { getWkFromMonth, WeekSwitch } from '../harmonogram/harmonogram'
// import { epochToDate } from '../helpers'
// import {
//     AppBar, Avatar, Button, Card, CardContent, Chip, Divider, Grid,
//     InputAdornment, Tab, Tabs, TextField, Toolbar,
//     Typography, useMediaQuery, useTheme
// } from '@mui/material'
// import { TabPanel } from '../configure/configure'
// import { getInstructor, getTrainingByID, postReservation } from '../apicalls/instructor.api'
// import { CardGiftcardTwoTone, Instagram, Loupe, Star } from '@mui/icons-material'
// import InstagramIcon from "@mui/icons-material/Instagram";
// import FacebookIcon from "@mui/icons-material/Facebook";
// import TwitterIcon from "@mui/icons-material/Twitter";
// import { MulwiColors } from '../mulwiColors'
// import { getErrorDialog, getNullDialog, StatusDialog } from '../StatusDialog'
// import { TrainingDetailsSideContent } from '../harmonogram/trainingDetails'
// import { Summary } from './Summary'
// import { TrainingSummary } from './trainingSummary'
// import { WaitConfirm } from './WaitConfirm'
// import { DrawerResponsive } from '../card/DrawerResponsive'
// import { HarmonogramDay } from '../harmonogram/day'
// import { HarmonogramMonth } from '../harmonogram/month'
// import { InstructorInfo } from './instructorInfo'
// import DrawerSmall from '../card/DrawerSmall'
// import TrainingAtc from '../harmonogram/trainingAtc'

// export function Rsv(props) {

//     const [instructorID, setInstructorID] = useState(null)
//     const [instructor, setInstructor] = useState(null)
//     const [training, setTraining] = useState(null)
//     const [ready, setReady] = useState(false)

//     const [accessToken, setAccessToken] = useState("")

//     //const [registerData, setRegisterData] = useState(null)

//     const [wk, setWk] = useState(getWkFromMonth(new Date()))
//     const [day, setDay] = useState(new Date())

//     useEffect(() => {
//         setStateFromQuery(window.location)
//     }, [])

//     const history = useHistory()

//     useEffect(() => {
//         return history.listen((location) => {
//             setStateFromQuery(location)
//         })
//     }, [history])

//     const theme = useTheme()
//     const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))
//     const isXLowRes = useMediaQuery(theme.breakpoints.down('sm'))

//     const [value, setValue] = React.useState(0)
//     const handleChange = (event, newValue) => {
//         setValue(newValue);
//     }

//     function resetTraining() {
//         setTraining(null)
//     }

//     const [info, setInfo] = useState(getNullDialog())

//     const [drawerOpen, _setDrawerOpen] = useState(false)
//     const [drawerData, setDrawerData] = useState({})

//     function setDrawerOpen(x) {
//         if (x) {
//             siv("hidden")
//         } else {
//             siv("inherit")
//         }
//         _setDrawerOpen(x)
//     }

//     function onTrainingSelected() {
//         sd([0, 0, 1])
//         setValue(1)
//     }

//     const [iv, siv] = useState("inherit")

//     const [d, sd] = useState([0, 1, 1])

//     async function onConfirm(d, c) {
//         if(!training || !training.Training) return
//         sd([1, 1, 0])
//         let r = {
//             TrainingID: training.Training.ID,
//             Occurrence: drawerData && drawerData.sch.Start,
//             NoRedirect: true,
//             UserData: d,
//             ContactData: c,
//             UseSavedData: false
//         }
//         try {
//             // perform reservation
//             let res = await postReservation(r)
//             res = JSON.parse(res)
//             setAccessToken(res.AT)
//             window.open(res.Url)
//             //window.open("https://test.adyen.link/PLF522D23EB8E96DD3")
//             setValue(2)

//         } catch (ex) {
//             setInfo(getErrorDialog("Nie udało się zrobić rezerwacji", ex))
//         }
//     }

//     async function setInstructorFromApi(iid) {
//         try {
//             let i = await getInstructor(iid)
//             i = JSON.parse(i)
//             setInstructor(i)
//         } catch (ex) {
//             console.log(ex)
//             // TODO: log ex
//         }
//     }

//     async function setTrainingFromApi(id) {
//         try {
//             let i = await getTrainingByID(id)
//             setTraining(i[0])
//             setReady(true)
//         } catch (ex) {
//             console.log(ex)
//             // TODO: log ex
//         }
//     }

//     function setStateFromQuery(l) {
//         let query = new URLSearchParams(l.search)

//         let _instructorID = query.get("instructorID")
//         if (_instructorID) {
//             setInstructorID(_instructorID)
//             setInstructorFromApi(_instructorID)
//         }

//         let _trainingID = query.get("trainingID")
//         if (_trainingID) {
//             setTrainingFromApi(_trainingID)
//         }

//         let _dateStart = query.get("dateStart")
//         if (_dateStart) {
//             let d = epochToDate(_dateStart)
//             if (d && isFinite(d)) {
//                 setWk(getWkFromMonth(d))
//                 setDay(d)
//             }
//         }

//     }


//     return (
//         <div style={{ backgroundColor: "white" }}>
//             <StatusDialog info={info} setInfo={setInfo} />
//             <DrawerSmall
//                     padding={7}
//                     navContent={value === 0 ?
//                         <Button
//                             style={{
//                                 color: "white",
//                                 backgroundColor: MulwiColors.greenDark,
//                             }}
//                             onClick={() => {
//                                 if (isLowRes)
//                                     setDrawerOpen(false)
//                                 onTrainingSelected()
//                             }}
//                             fullWidth variant="contained">
//                             Zapisz się na trening
//                         </Button> : null
//                     }
//                     content={
//                         <React.Fragment>
//                             <TrainingSummary
//                                 sch={drawerData && drawerData.sch}
//                                 training={drawerData && drawerData.training}
//                                 setDrawerOpen={setDrawerOpen} />
//                             {/* <Button onClick={() => {
//                                 setDrawerOpen(false)
//                                 onTrainingSelected()
//                             }} 
//                             style={{color: "white", backgroundColor:MulwiColors.greenDark}}>
//                                 Zarezerwuj
//                             </Button> */}
//                         </React.Fragment>
//                         //     <TrainingDetailsSideContent 
//                         //         onRsv={onTrainingSelected}
//                         //     //onChange={() => setRefreshToken(!refreshToken) }
//                         //     drawerData={drawerData}
//                         //     setInfo={setInfo}
//                         //     //setDrawerData={setDrawerData}
//                         //     setDrawerOpen={setDrawerOpen}
//                         //   />
//                     }
//                     open={drawerOpen}
//                     width={isXLowRes ? "100%" : 850}
//                     onClose={() => setDrawerOpen(false)}
//                     onOpen={() => setDrawerOpen(true)} >
//                 <React.Fragment>
//                     <Grid direction={"column"}
//                         style={{
//                             marginBottom: 10,
//                             paddingTop: isLowRes ? 60 : 10,
//                             paddingLeft: 10,
//                         }}
//                         spacing={2}
//                         container >
//                         <Grid item>
//                             <Tabs value={value} onChange={handleChange}
//                                 indicatorColor="primary"
//                                 textColor="primary">

//                                 <Tab label="Wybierz trening"
//                                     disabled={Boolean(d[0])}
//                                     id="tt1" aria-controls="stt1" />
//                                 <Tab label="Twoje dane"
//                                     disabled={Boolean(d[1])}
//                                     id="tt2" aria-controls="stt2" />
//                                 <Tab label="Płatność"
//                                     disabled={Boolean(d[2])}
//                                     id="tt3" aria-controls="stt3" />

//                             </Tabs>

//                         </Grid>
//                     </Grid>
//                     <div style={{ backgroundColor: MulwiColors.whiteBackground }}>
//                         <TabPanel value={value} index={0}>
//                             <Grid direction={"row"}
//                                 style={{
//                                     marginBottom: 10,
//                                 }}
//                                 spacing={2}
//                                 justify={isLowRes ? "center" : "flex-start"}
//                                 alignItems="center"
//                                 container>

//                                 <Grid item>
//                                         <Typography style={{
//                                             paddingLeft: 5,
//                                             color: MulwiColors.greenDark,
//                                         }} variant="h4">
//                                             Wybierz trening
//                                         </Typography>
//                                 </Grid>
//                                 <Grid item>
//                                     {/* <TextField
//                                         id="input-with-icon-textfield"
//                                         label="albo znajdź inny trening"
//                                         size="small"
//                                         variant="outlined"
//                                         InputProps={{
//                                             endAdornment: (
//                                                 <InputAdornment position="end">
//                                                     <Loupe style={{
//                                                         color: "gray"
//                                                     }} />
//                                                 </InputAdornment>
//                                             )
//                                         }}
//                                     /> */}
//                                     <TrainingAtc forUsr 
//                                         instructorID={instructorID} 
//                                         setValue={setTraining} 
//                                         value={training} />
//                                 </Grid>
//                                 {/* {training && (
//                                     <Chip
//                                         color="primary"
//                                         label={training.Title}
//                                         onDelete={resetTraining} />
//                                 )} */}

//                             </Grid>

//                             <Grid container direction="row" justify="space-between"
//                                 alignItems="stretch" spacing={3}>
//                                 <Grid item>
//                                     <Grid container
//                                         justify="flex-start"
//                                         alignItems="center">
//                                         <Grid item style={{paddingLeft: 20}}>
//                                             {ready && (
//                                                 <React.Fragment>
//                                                     {(isLowRes && (
//                                                         <React.Fragment>
//                                                             <Grid item style={{
//                                                                 width: "95%"
//                                                             }}>
//                                                                 <HarmonogramMonth
//                                                                     user
//                                                                     trainingID={training && training.Training && training.Training.ID}
//                                                                     instructorID={instructorID}

//                                                                     day={day}
//                                                                     setDay={setDay}

//                                                                     refreshToken={null}

//                                                                     setInfo={setInfo}

//                                                                     // sessions={props.sessions} 
//                                                                     // tStart={props.tStart}
//                                                                     // tEnd={props.tEnd}

//                                                                     date={day} />
//                                                                 <HarmonogramDay
//                                                                     user
//                                                                     trainingID={training  && training.Training && training.Training.ID}
//                                                                     instructorID={instructorID}

//                                                                     refreshToken={null}
//                                                                     //sessions={props.sessions} 
//                                                                     day={day}
//                                                                     setDrawerData={setDrawerData}
//                                                                     setDrawerOpen={setDrawerOpen}
//                                                                     // tStart={props.tStart}
//                                                                     // tEnd={props.tEnd} 

//                                                                     setInfo={setInfo}

//                                                                 />
//                                                             </Grid>
//                                                         </React.Fragment>
//                                                     )) || (
//                                                             <React.Fragment>
//                                                                 <Grid item>
//                                                                     <WeekSwitch week={wk} setWeek={setWk} />
//                                                                     <HarmonogramWeek
//                                                                         spclick={() => {
//                                                                             onTrainingSelected()
//                                                                         }}
//                                                                         user
//                                                                         trainingID={training  && training.Training && training.Training.ID}
//                                                                         instructorID={instructorID}
//                                                                         refreshToken={null}
//                                                                         setInfo={setInfo}
//                                                                         setDrawerData={setDrawerData}
//                                                                         setDrawerOpen={setDrawerOpen}
//                                                                         week={wk} />
//                                                                 </Grid>
//                                                             </React.Fragment>
//                                                         )}
//                                                 </React.Fragment>
//                                             )}
//                                         </Grid>
//                                     </Grid>
//                                 </Grid>
//                                 {instructor && (<Grid item lg style={{
//                                                         visibility: iv
//                                                     }}>
//                                     <Grid container justify="center" alignItems="center">
//                                         <Grid item>
//                                             <Card>
//                                                 <CardContent>
//                                                     <InstructorInfo 
//                                                         user={props.user}
//                                                         setInfo={setInfo} 
//                                                         instructor={instructor} />
//                                                 </CardContent>
//                                             </Card>
//                                         </Grid>
//                                     </Grid>
//                                 </Grid>)}
//                             </Grid>
//                         </TabPanel>
//                     </div>
//                     <TabPanel value={value} index={1}>
//                         {/* <Grid container direction="row" justify="space-evenly"
//                         alignItems="stretch" spacing={3}>
//                     <Grid item>
//                         <Summary/>
//                     </Grid>
//                     <Grid item>
//                         <TrainingSummary 
//                             drawerData={drawerData} 
//                             setDrawerOpen={setDrawerOpen}
//                             onRsv={onTrainingSelected}/>
//                     </Grid>
//                 </Grid> */}
//                         <Grid container>
//                             <Grid item lg={6}>
//                                 <Summary onConfirm={onConfirm} />
//                             </Grid>
//                         </Grid>
//                     </TabPanel>
//                     <TabPanel value={value} index={2}>
//                         <Grid container>
//                             <Grid item lg={6}>
//                                 <WaitConfirm accessToken={accessToken} />
//                             </Grid>
//                         </Grid>
//                     </TabPanel>
//                 </React.Fragment>
//             </DrawerSmall>
//         </div>
//     )
// }