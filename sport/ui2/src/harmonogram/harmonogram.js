import {
    Button,
    Grid,
    Typography,
    Switch as _MUISwitch,
    useMediaQuery, useTheme, FormControl, Container
} from '@mui/material'
import {
    Add,
    NavigateBefore, NavigateNext,
} from '@mui/icons-material'
import React, { useEffect, useState } from 'react'
import { HarmonogramMonth } from './month'
import { HarmonogramWeek } from './weekBigRes'
import { dateStartOfWeek, lMonths } from './../helpers'
import { HarmonogramDay } from './day'
import { MulwiColors } from '../mulwiColors'
import { DrawerResponsive } from '../card/DrawerResponsive'
import { TrainingDetailsSideContent } from './trainingDetails'
import { DeleteTraining } from './deleteTraining'
import { AddTraining } from './addTraining'
import { Route, Switch, useHistory, useLocation } from 'react-router-dom';
import { ListTrainings } from './trainingList'
import { getNullDialog, StatusDialog } from '../StatusDialog'
import PropTypes from 'prop-types';
import { TrainingSummary } from '../reservations/trainingSummary'
import { ListRsv } from './rsvList'
import TrainingAtc from './trainingAtc'
import { CalendarEditor, CalendarHeader } from './calendarEditor'
import { locale2 } from '../locale'
import { Checkbox, FormControlLabel } from '@mui/material'
import { makeStyles, withStyles } from "@mui/styles";
import { Schedule } from './schedule'

export function WeekSwitch(props) {

    let months = []

    for(let i = 0; i < 12; i++) {
        months.push(lMonths[i][props.lang])
    }

    return (
        <React.Fragment>
            <div style={{
                width: 460,
                margin: "auto",
                marginTop: 10,
                marginBottom: 10
            }}>
                <Button
                    onClick={() => {
                        let newStart = new Date(props.week.start)
                        newStart.setDate(newStart.getDate() - 7)
                        let newEnd = new Date(props.week.end)
                        newEnd.setDate(newEnd.getDate() - 7)
                        newEnd.setHours(23, 59, 59)
                        props.setWeek({ e: true, start: newStart, end: newEnd })
                        if (props.setMonthDate)
                            props.setMonthDate(newStart)
                    }}
                    variant="text" color="primary"><NavigateBefore /></Button>
                <div style={{ display: "inline-block", minWidth: 300, textAlign: "center" }}>
                    {
                        (props.week.start.getMonth() === props.week.end.getMonth() && (
                            <span>
                                {props.week.start.getDate()} - {props.week.end.getDate()} {months[props.week.start.getMonth()]} {props.week.start.getFullYear()}
                            </span>
                        )) || (
                            <span>
                                {props.week.start.getDate()} {months[props.week.start.getMonth()]} {props.week.start.getFullYear()} - {props.week.end.getDate()} {months[props.week.end.getMonth()]} {props.week.end.getFullYear()}
                            </span>
                        )
                    }
                </div>
                <Button
                    onClick={() => {
                        let newStart = new Date(props.week.start)
                        newStart.setDate(newStart.getDate() + 7)
                        let newEnd = new Date(props.week.end)
                        newEnd.setDate(newEnd.getDate() + 7)
                        newEnd.setHours(23, 59, 59)
                        props.setWeek({ e: true, start: newStart, end: newEnd })
                        if (props.setMonthDate)
                            props.setMonthDate(newStart)
                    }}
                    variant="text" color="primary"><NavigateNext /></Button>
            </div>
        </React.Fragment>
    )
}

WeekSwitch.propTypes = {
    week: PropTypes.object.isRequired,
    setWeek: PropTypes.func.isRequired,
    setMonthDate: PropTypes.func
}

export function getWkFromMonth(d) {
    let startOfWk = dateStartOfWeek(d)
    let endOfWk = new Date(startOfWk)
    endOfWk.setDate(endOfWk.getDate() + 6)
    endOfWk.setHours(23, 59, 59)
    return { e: true, start: startOfWk, end: endOfWk }
}

export const MUISwitch = withStyles({
    switchBase: {
        color: MulwiColors.blueDark,
        '&$checked': {
            color: MulwiColors.blueDark,
        },
        '&$checked + $track': {
            backgroundColor: "white",
        },
    },
    checked: {},
    track: {},
})(_MUISwitch);

export function Harmonogram(props) {

    let months = []
    for(let i = 0; i < 12; i++) {
        months.push(lMonths[i][props.lang])
    }

    let basePath = props.basePath

    if (!basePath)
        basePath = "/harmonogram"

    const [monthDate, setMonthDate] = useState(new Date())
    const [week, setWeek] = useState(() => getWkFromMonth(new Date()))
    const [day, setDay] = useState(new Date())

    const [st, setSt] = useState(0)

    const [hideMultiday, setHideMultiday] = useState(false)
    const [showUnavailable, setShowUnavailable] = useState(false)

    function setWeekLevelWithDate(start, end) {
        setWeek({ e: true, start: start, end: end })
    }

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))
    //const isXLowRes = useMediaQuery(theme.breakpoints.down('xs'))

    function setWkFromMonth(d) {
        setWeek(getWkFromMonth(d))
    }

    function resetHarmonogram() {
        let d = new Date()
        setMonthDate(d)
        setWkFromMonth(d)
    }


    useEffect(() => {
        let query = new URLSearchParams(window.location.search)
        if (query.get("open_add_train")) {
            setSt(1)
            //setSelectedTraining(true)
        }
    }, [])

    function getMonthLabel() {
        return (
            <React.Fragment>
                <div style={{
                    width: 340,
                    margin: "auto",
                    marginTop: 10,
                    marginBottom: 10
                }}>
                    <Button
                        onClick={() => {
                            let c = new Date(monthDate)
                            c.setMonth(c.getMonth() - 1)
                            c = new Date(c)
                            c.setDate(1)
                            setMonthDate(c)
                            setWkFromMonth(c)
                        }}
                        variant="text" color="primary"><NavigateBefore /></Button>
                    <div style={{ display: "inline-block", minWidth: 120, textAlign: "center" }}>
                        {months[monthDate.getMonth()]} {monthDate.getFullYear()}
                    </div>
                    <Button
                        onClick={() => {
                            let c = new Date(monthDate)
                            c.setMonth(c.getMonth() + 1)
                            c = new Date(c)
                            c.setDate(1)
                            setMonthDate(c)
                            setWkFromMonth(c)
                        }}
                        variant="text" color="primary"><NavigateNext /></Button>
                    <Button variant="text"
                        onClick={resetHarmonogram}>Reset</Button>
                </div>
            </React.Fragment>
        )
    }

    // function getHeaderLabel() {
    //     switch(level) {
    //     case 0:
    //         return getMonthLabel()
    //     case 1:
    //         return getWeekLabel()
    //     }
    // }

    const [info, setInfo] = useState(getNullDialog())

    const [drawerOpen, setDrawerOpen] = useState(false)
    const [drawerData, setDrawerData] = useState({})
    const [refreshToken, setRefreshToken] = useState(false)

    const history = useHistory()
    const location = useLocation()
    const [selectedTraining, setSelectedTraining] = useState(null)

    // useEffect(() => {
    //     let query = new URLSearchParams(window.location.search)
    //     if (query.get("open_add_train")) {
    //         setSelectedTraining(true)
    //     }
    // }, [])

    useEffect(() => {
        setRefreshToken(!refreshToken)
        // on location change close drawer
        setDrawerOpen(false)
        // eslint-disable-next-line
    }, [location])

    function defaultHeaderContent() {
        return (
        <React.Fragment>
        <Grid item>
            <TrainingAtc
                lang={props.lang}
                forUsr={props.usrRsv}
                value={selectedTraining} setValue={setSelectedTraining} />
        </Grid>
        {!props.usrRsv && isLowRes &&  (
            <Grid item>
                <Button variant="contained"
                    onClick={() => {
                        if (st === 0) {
                            setSt(1)
                        } else {
                            setSt(0)
                        }
                    }}
                    style={{
                        backgroundColor: st === 0 ? MulwiColors.greenDark : "inherit",
                        color: st === 0 ? "white" : "black",
                        display: st ===2 && "none"
                    }}>{st !== 0 ? locale2.CLOSE_EDITOR[props.lang] : <React.Fragment>
                        <Add />
                        {locale2.ADD_TRAINING[props.lang]}</React.Fragment>}
                </Button>
            </Grid>
        )}
        <Grid item>
            <Grid
                style={{ cursor: "pointer" }}
                component="label" container alignItems="center" spacing={1}>
                <Grid item>{locale2.SCHEDULE[props.lang]}</Grid>
                <Grid item>
                    <MUISwitch
                        disableRipple
                        checked={location.pathname === basePath ? false : true}
                        onChange={(e) => {
                            setSt(0)
                            if (!e.target.checked && window.location.pathname !== basePath) {
                                history.push(basePath)
                            }
                            if (e.target.checked && window.location.pathname !== (basePath + "/list")) {
                                history.push(basePath + "/list")
                            }
                        }} />
                </Grid>
                <Grid item>{ locale2.LIST[props.lang] }</Grid>
            </Grid>
        </Grid>
    </React.Fragment>)
    }

    const [calendarSaveToken, setCalendarSaveToken] = useState(0)

    return (
        <React.Fragment>
            <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
            <DrawerResponsive
                padding={7}
                content={props.usrRsv ? <TrainingSummary lang={props.lang}
                    usrRsv={props.usrRsv}
                    setInfo={setInfo}
                    onChange={() => setRefreshToken(!refreshToken)}
                    sch={drawerData && drawerData.sch}
                    resize={drawerOpen}
                    training={drawerData && drawerData.training}
                    setDrawerOpen={setDrawerOpen} /> : <TrainingDetailsSideContent
                    onChange={() => setRefreshToken(!refreshToken)}
                    drawerData={drawerData}
                    setInfo={setInfo}
                    lang={props.lang}
                    setDrawerData={setDrawerData}
                    setDrawerOpen={setDrawerOpen} />}
                navContent={!props.usrRsv && drawerData && (!drawerData.sch || !drawerData.sch.IsOrphaned) &&
                    <React.Fragment>
                        <Grid container
                            direction="row"
                            style={{
                                marginTop: 4,
                            }}
                            spacing={1}
                            justify="flex-end"
                            alignItems="center">
                            <Grid item>
                                <DeleteTraining
                                    lang={props.lang}
                                    setDrawerOpen={setDrawerOpen}
                                    onChange={() => {
                                        setSt(0)
                                        setRefreshToken(!refreshToken)
                                    }}
                                    id={drawerData &&
                                        drawerData.training &&
                                        !(drawerData.session && drawerData.session.IsOrphaned) &&
                                        drawerData.training.ID}
                                />
                            </Grid>
                        </Grid>
                    </React.Fragment>
                }
                open={drawerOpen}
                width={isLowRes ? "100vw" : 500}
                onClose={() => setDrawerOpen(false)}
                onOpen={() => setDrawerOpen(true)}>
                    {!props.nohdr && (<Grid direction={"row"}
                        style={{
                            marginBottom: 10,
                            backgroundColor: "white",
                            paddingTop: isLowRes || (props.instructor && props.instructor.Config != 0) ? 0 : 8,
                            paddingLeft: isLowRes ? 0 : 10,
                        }}
                        spacing={2}
                        justify={isLowRes ? "center" : "flex-start"}
                        alignItems="center"
                        container>
                            <Grid item>
                                <Typography style={{
                                    paddingLeft: 5,
                                    color: MulwiColors.greenDark,
                                }} variant="h4">
                                    {props.usrRsv ? locale2.YOUR_RSVS[props.lang] : locale2.SCHEDULE[props.lang]}
                                </Typography>
                            </Grid>
                            {st === 2 ? <CalendarHeader 
                                    lang={props.lang}
                                    onClose={() => {
                                        setSt(0)
                                    }} 
                                    setSaveToken={setCalendarSaveToken} 
                                    isLowRes={isLowRes} 
                                    drawerData={drawerData} /> : defaultHeaderContent()}
                    </Grid>)}
                <div style={{
                    width: "100%",
                    height: "100%",
                    backgroundColor: "white"
                }}>
                    {(() => {
                        switch (st) {
                            case 1:
                                return (<AddTraining
                                    lang={props.lang}
                                    onChange={() => setRefreshToken(!refreshToken)}
                                    setDrawerData={setDrawerData}
                                    openDrawer={() => setDrawerOpen(true)}
                                    onClose={(success) => {
                                        if (success) setSt(2)
                                        else setSt(0)
                                    }} />)
                                return null
                            case 2:
                                return (<React.Fragment>
                                    {props.nohdr && (
                                        <Container>
                                            <Grid container direction="row" 
                                                    spacing={2}>
                                                <CalendarHeader 
                                                    lang={props.lang}
                                                    onClose={() => {
                                                        setSt(0)
                                                    }} 
                                                    setSaveToken={setCalendarSaveToken} 
                                                    isLowRes={isLowRes} 
                                                    drawerData={drawerData} />
                                            </Grid>
                                        </Container>
                                    )}
                                    <CalendarEditor
                                        lang={props.lang}
                                        setSaveToken={setCalendarSaveToken}
                                        saveToken={calendarSaveToken}
                                        drawerData={drawerData}
                                        setDrawerData={setDrawerData}
                                        close={() => {
                                            setSt(0)
                                        }}
                                        onChange={() => {
                                            setRefreshToken(!refreshToken)
                                        }} />
                                </React.Fragment>)
                            default:
                                return (<Switch>
                                    <Route path={basePath + "/list"}>
                                        {(props.usrRsv && (
                                            <ListRsv
                                                lang={props.lang}
                                                setDrawerData={setDrawerData}
                                                setDrawerOpen={setDrawerOpen}
                                                refreshToken={refreshToken}
                                                setRefreshToken={setRefreshToken}
                                            />
                                        )) ||
                                            <ListTrainings
                                                lang={props.lang}
                                                setDrawerData={setDrawerData}
                                                setDrawerOpen={setDrawerOpen}
                                                refreshToken={refreshToken}
                                                setRefreshToken={setRefreshToken}
                                            />}
                                    </Route>
                                    <Route path="/">
                                        <Grid justify={"space-evenly"}
                                            alignItems={"center"}
                                            direction="row"
                                            container
                                            spacing={0}>
                                            {(!isLowRes && (
                                                <React.Fragment>
                                                    <Grid item style={{
                                                        width: "100%"
                                                    }}>

                                                        <Grid container direction="row">
                                                            <Grid item>
                                                                <FormControl>
                                                                    <FormControlLabel control={<Checkbox checked={hideMultiday} onChange={(e) => {
                                                                            setHideMultiday(e.target.checked)
                                                                        }} />} label={locale2.HIDE_MULTIDAY[props.lang]} />
                                                                </FormControl>
                                                            </Grid>
                                                            <Grid item>
                                                                <FormControl>
                                                                    <FormControlLabel control={<Checkbox checked={showUnavailable} onChange={(e) => {
                                                                            setShowUnavailable(e.target.checked)
                                                                        }} />} label={"Pokaż niedostępne treningi"} />
                                                                </FormControl>
                                                            </Grid>
                                                        </Grid>
                                                        <WeekSwitch lang={props.lang} 
                                                            week={week} setWeek={setWeek} setMonthDate={setMonthDate} />
                                                        <Schedule

                                                            lang={props.lang}

                                                            trainingID={selectedTraining && selectedTraining.Training && selectedTraining.Training.ID}

                                                            usrRsv={props.usrRsv}

                                                            refreshToken={refreshToken}
                                                            setRefreshToken={setRefreshToken}
                                                            // setInfo={props.setInfo}
                                                            // resetInfo={props.resetInfo}
                                                            // short={true} 

                                                            // sessions={props.sessions} 
                                                            // tStart={props.tStart}
                                                            // tEnd={props.tEnd}

                                                            setInfo={setInfo}

                                                            hideMultiday={hideMultiday}
                                                            showUnavailable={showUnavailable}

                                                            setDrawerData={setDrawerData}
                                                            setDrawerOpen={setDrawerOpen}
                                                            week={week} />
                                                    </Grid>
                                                </React.Fragment>
                                            ))}
                                            {isLowRes && (
                                            <Grid item>
                                                <div style={{
                                                    display: "table",
                                                    margin: "0 auto",
                                                   // visibility: drawerOpen ? "hidden" : "inherit"
                                                }}>
                                                    {getMonthLabel()}
                                                    <HarmonogramMonth
                                                        lang={props.lang}

                                                        day={day}
                                                        setDay={setDay}

                                                        trainingID={selectedTraining && selectedTraining.Training && selectedTraining.Training.ID}

                                                        usrRsv={props.usrRsv}

                                                        refreshToken={refreshToken}

                                                        setInfo={setInfo}
                                                        showUnavailable={showUnavailable}

                                                        // sessions={props.sessions} 
                                                        // tStart={props.tStart}
                                                        // tEnd={props.tEnd}

                                                        week={week}
                                                        switchToWeek={setWeekLevelWithDate}
                                                        date={monthDate}

                                                    />
                                                </div>
                                            </Grid>
                                            )}
                                            {(isLowRes && (
                                                <React.Fragment>
                                                    <Grid item style={{
                                                        width: "95%"
                                                    }}>
                                                        <HarmonogramDay
                                                            lang={props.lang}

                                                            refreshToken={refreshToken}
                                                            //sessions={props.sessions} 
                                                            day={day}
                                                            trainingID={selectedTraining && selectedTraining.Training && selectedTraining.Training.ID}

                                                            usrRsv={props.usrRsv}
                                                            setDrawerData={setDrawerData}
                                                            setDrawerOpen={setDrawerOpen}
                                                            // tStart={props.tStart}
                                                            // tEnd={props.tEnd} 

                                                            setInfo={setInfo}

                                                        />
                                                    </Grid>
                                                </React.Fragment>
                                            ))}
                                        </Grid>
                                    </Route>
                                </Switch>)
                        }
                    })()}
                </div>
            </DrawerResponsive>
        </React.Fragment>
    )
}
