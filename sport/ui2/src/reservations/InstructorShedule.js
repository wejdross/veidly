import { Button, Card, CardContent, Container, Dialog, 
        DialogActions, DialogContent, Grid, 
        Typography, useMediaQuery, useTheme } from '@mui/material'
import React, { useEffect, useState } from 'react'
import { useHistory } from 'react-router'
import { getInstructor, getTrainingByID } from '../apicalls/instructor.api'
import { getSmForInstr } from '../apicalls/sm'
import DrawerSmall from '../card/DrawerSmall'
import { HarmonogramDay } from '../harmonogram/day'
import { getWkFromMonth, WeekSwitch } from '../harmonogram/harmonogram'
import { HarmonogramMonth } from '../harmonogram/month'
import { MonthLabel } from '../harmonogram/MonthLabel'
import TrainingAtc from '../harmonogram/trainingAtc'
import { HarmonogramWeek } from '../harmonogram/weekBigRes'
import { dateToEpoch, epochToDate } from '../helpers'
import { MulwiColors } from '../mulwiColors'
import { getNullDialog, StatusDialog } from '../StatusDialog'
import { SubCard } from '../sub/SubCard'
import { InstructorInfo } from './instructorInfo'
import { TrainingSummary } from './trainingSummary'
import { getSupportedLanguage, locale2 } from '../locale'
import { fitDateToOcc, prepOccs } from '../harmonogram/calendarEditor'
import { getErrorDialog } from '../StatusDialog'
import { Schedule } from '../harmonogram/schedule'

export function InstructorShedule(props) {

    const [instructorID, setInstructorID] = useState(null)
    const [instructor, setInstructor] = useState(null)
    const [training, setTraining] = useState(null)
    const [sm, setSm] = useState(null)
    const [wk, setWk] = useState(getWkFromMonth(new Date()))
    const [monthDate, _setMonthDate] = useState(new Date())

    function setMonthDate(v) {
        _setMonthDate(v)
        setWk(getWkFromMonth(v))
    }

    const [day, setDay] = useState(new Date())

    useEffect(() => {
        setStateFromQuery(window.location)
    }, [])

    const history = useHistory()

    useEffect(() => {
        return history.listen((location) => {
            setStateFromQuery(location)
        })
    }, [history])

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))
    const downMd = useMediaQuery(theme.breakpoints.down('md'))
    const h = useHistory()
    
    const [info, setInfo] = useState(getNullDialog())

    const [drawerOpen, _setDrawerOpen] = useState(false)
    const [drawerData, setDrawerData] = useState({})

    function setDrawerOpen(x) {
        _setDrawerOpen(x)
    }

    function onTrainingSelected() {
        let l = ("/rsv?instructorID=" +
            drawerData.training.InstructorID +
            "&trainingID=" +
            drawerData.training.ID +
            "&dateStart=" + dateToEpoch(new Date(drawerData.sch.Start)))
        h.push(l)
    }

    async function setInstrSm(iid) {
        try {
            let i = await getSmForInstr(iid)
            i = JSON.parse(i)
            setSm(i)
        } catch (ex) {
            setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    async function setInstructorFromApi(iid) {
        try {
            let i = await getInstructor(iid)
            i = JSON.parse(i)
            setInstructor(i)
        } catch (ex) {
            setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    async function setTrainingFromApi(id) {
        try {
            let i = await getTrainingByID(id)
            setTraining(i[0])
        } catch (ex) {
            setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    useEffect(() => {
        if(!props.instructor)
            return
        setStateFromQuery(window.location)
        setInstructor(props.instructor)
        let iid = props.instructor.id
        setInstructorID(iid)
        setInstrSm(iid)
    }, [props.instructor])

    function setStateFromQuery(l) {
        let query = new URLSearchParams(l.search)

        let _instructorID = query.get("instructorID")
        if (_instructorID && !props.instructor) {
            setInstructorID(_instructorID)
            setInstructorFromApi(_instructorID)
            setInstrSm(_instructorID)
        }

        let _trainingID = query.get("trainingID")
        if (_trainingID) {
            setTrainingFromApi(_trainingID)
        }

        let _dateStart = query.get("dateStart")
        if (_dateStart) {
            let d = epochToDate(_dateStart)
            if (d && isFinite(d)) {
                setWk(getWkFromMonth(d))
                setMonthDate(d)
            }
        } else {
            let d = new Date()
            setWk(getWkFromMonth(d))
            setMonthDate(d)
        }
    }

    function setWeekLevelWithDate(start, end) {
        setWk({ e: true, start: start, end: end })
    }

    const [schedOpen, setSchedOpen] = useState(false)
    const [schedWk, setSchedWk] = useState(getWkFromMonth(new Date()))
    const [schedData, setSchedData] = useState([])

    let lang = getSupportedLanguage()

    return (<React.Fragment>
		
        <Dialog fullWidth maxWidth="md" open={schedOpen}
            onClose={() => setSchedOpen(false)}>
            <DialogContent>
                <Grid container direction="row" justifyContent="center">
                    <Grid item>
                        <WeekSwitch week={schedWk} setWeek={setSchedWk} />
                        <HarmonogramWeek
                            setInfo={setInfo}
                            readonly
                            lang={props.lang}
                            editorData={schedData} week={schedWk} />
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setSchedOpen(false)}>
                    {locale2.CLOSE[lang]}
                </Button>
            </DialogActions>
        </Dialog>

        {drawerOpen && (
            (<Grid container direction="row" style={{
                position: "fixed",
                top: isLowRes ? 6 : 73,
                right: isLowRes ? 6 : 12,
                maxWidth: isLowRes ? 260 : 600,
                zIndex: 9999,
            }} justifyContent="flex-end" spacing={2}>
                {/* {!downMd && (<Grid item>
                    <Button
                        style={{
                            color: "white",
                            backgroundColor: MulwiColors.blueDark,
                            maxWidth: isLowRes ? 260 : 270,
                        }}
                        onClick={() => {
                            if (!drawerData || !drawerData.sch) return
                            let occs = [JSON.parse(JSON.stringify(drawerData.sch.Occ))]
                            occs[0].DateStart = new Date(drawerData.sch.Start)
                            occs[0].DateEnd = new Date(drawerData.sch.End)
                            prepOccs(occs)
                            fitDateToOcc(occs, setSchedWk)
                            setSchedData({
                                occs: occs,
                                training: drawerData.training
                            })
                            setSchedOpen(true)
                        }}
                        fullWidth variant="contained">
                        {locale2.SCHEDULE[props.lang]}
                    </Button>
                </Grid>)} */}
                <Grid item>
                    <Button
                        style={{
                            color: "white",
                            backgroundColor: MulwiColors.greenDark,
                            maxWidth: isLowRes ? 260 : 270,
                        }}
                        onClick={() => {
                            onTrainingSelected()
                        }}
                        fullWidth variant="contained">
                    {locale2.SIGN_IN_FOR_THIS_TRAINING[props.lang]}
                    </Button>
                </Grid>
            </Grid>)
        )}
        <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
        <DrawerSmall
            padding={7}
            content={
                <React.Fragment>
                    <TrainingSummary lang={props.lang}
                        sch={drawerData && drawerData.sch}
                        training={drawerData && drawerData.training}
                        setDrawerOpen={setDrawerOpen}
                        dateStart={drawerData && drawerData.sch && drawerData.sch.Start}
                        resize={drawerOpen} />
                </React.Fragment>
            }
            open={drawerOpen}
            width={isLowRes ? "100vw" : 700}
            onClose={() => setDrawerOpen(false)}
            onOpen={() => setDrawerOpen(true)} >
            {!props.onlySchedule && <Grid direction={"row"}
                style={{
                    marginBottom: isLowRes ? 0 : 10,
                    backgroundColor: "white",
                    paddingTop: isLowRes ? 0 : 10,
                    paddingLeft: 10,
                }}
                spacing={2}
                justify="flex-start"
                alignItems="center"
                container>
                <Grid item>
                    <Typography style={{
                        paddingLeft: 5,
                        color: MulwiColors.greenDark,
                    }} variant="h4">
                        {locale2.INSTRUCTOR_SCHEDULE[props.lang]}
                    </Typography>

                </Grid>
                <Grid item>
                    <TrainingAtc forUsr 
                        lang={props.lang}
                        instructorID={instructorID}
                        setValue={setTraining}
                        value={training} />
                </Grid>
            </Grid>}
            {isLowRes ? (
                <Grid spacing={2}
                    container
                    direction="row"
                    justify="center"
                    alignItems="stretch">


                    <Grid item>
                        <Grid container direction="column">
                            <Grid item>
                                {instructor && (<Card style={{ marginBottom: 10 }}>
                                    <CardContent>
                                        <InstructorInfo lang={props.lang}
                                            user={props.user}
                                            setInfo={setInfo} instructor={instructor} />
                                    </CardContent>
                                </Card>)}
                            </Grid>
                        </Grid>
                    </Grid>

                    <Grid item>
                        <Grid>
                            <MonthLabel 
                                lang={props.lang}
                                monthDate={monthDate} setMonthDate={setMonthDate} />
                            <HarmonogramMonth
                                user
                                trainingID={training && training.Training && training.Training.ID}
                                instructorID={instructorID}

                                //day={day}
                                //setDay={setDay}

                                refreshToken={null}
                                setInfo={setInfo}
                                lang={props.lang}

                                date={monthDate}
                                week={wk}
                                switchToWeek={setWeekLevelWithDate}

                                day={day}
                                setDay={setDay}

                                lang={props.lang}
                            />
                        </Grid>
                    </Grid>

                    <Grid item style={{ width: "95%" }}>
                        <HarmonogramDay
                            user

                            day={day}
                            lang={props.lang}
                                                
                            trainingID={training && training.Training && training.Training.ID}
                            instructorID={instructorID}

                            setDrawerData={setDrawerData}
                            setDrawerOpen={setDrawerOpen}

                            refreshToken={null}
                            setInfo={setInfo}

                        />
                    </Grid>

                </Grid>
            ) : (<Container style={{
                marginBottom: 60
            }}>
                <Grid container direction="column"
                    justify="flex-start"
                    spacing={2}
                    style={{
                        marginBottom: 30,
                    }}
                    alignItems="stretch">
                    <Grid item>
                        <Grid spacing={2}
                            container
                            direction="row"
                            justify="center"
                            alignItems="stretch">
                            <Grid item xs={12} >
                                <Grid container direction="column" spacing={2}>
                                    <Grid item>
                                        <WeekSwitch week={wk} setWeek={setWk} 
                                                    lang={props.lang}
                                                    setMonthDate={setMonthDate} />
                                        <Schedule
                                            spclick={() => {
                                                onTrainingSelected()
                                            }}
                                            user
                                            trainingID={training && training.Training && training.Training.ID}
                                            instructorID={instructorID}
                                            refreshToken={null}
                                            setInfo={setInfo}
                                            lang={props.lang}
                                            setDrawerData={setDrawerData}
                                            setDrawerOpen={setDrawerOpen}
                                            week={wk} />
                                    </Grid>
                                    <Grid item>
                                        {sm && sm.length > 0 && (
                                            <Typography variant="h5" style={{
                                                marginBottom: 10
                                            }}>
                                                {locale2.INSTRUCTOR_CARNETS[props.lang]}
                                            </Typography>
                                        )}
                                        <Grid container direction="row" spacing={2}>
                                            {sm && sm.map((s, i) => (<Grid item sm={12} md={6}>
                                                <SubCard lang={props.lang} sm={s} user={props.user} />
                                            </Grid>))}
                                        </Grid>
                                    </Grid>
                                </Grid>
                            </Grid>
                            <Grid item lg={4} style={{
                                position: "relative"
                            }}>
                                <Grid
                                    style={{
                                        position: "sticky",
                                        top: 100
                                    }}
                                    container
                                    direction="column">

                                    {!props.onlySchedule && (<Grid item>
                                        {instructor && (<Card style={{ marginBottom: 10 }}>
                                            <CardContent>
                                                <InstructorInfo lang={props.lang}
                                                    user={props.user}
                                                    setInfo={setInfo}
                                                    instructor={instructor} />
                                            </CardContent>
                                        </Card>)}
                                    </Grid>)}
                                    {/* <Grid item>
                                        <Card>
                                            <CardContent>
                                                <Grid>
                                                    <MonthLabel 
                                                        lang={props.lang}
                                                        monthDate={monthDate} setMonthDate={setMonthDate} />
                                                    <HarmonogramMonth
                                                        user
                                                        trainingID={training && training.Training && training.Training.ID}
                                                        instructorID={instructorID}

                                                        //day={day}
                                                        //setDay={setDay}

                                                        refreshToken={null}
                                                        setInfo={setInfo}

                                                        date={monthDate}
                                                        lang={props.lang}
                                                        week={wk}
                                                        switchToWeek={setWeekLevelWithDate}
                                                    />
                                                </Grid>
                                            </CardContent>
                                        </Card>
                                    </Grid> */}
                                </Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </Container>)}
        </DrawerSmall>
    </React.Fragment>)
}