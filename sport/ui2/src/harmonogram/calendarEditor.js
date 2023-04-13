import {
    Backdrop, Button, ButtonBase,
    CircularProgress, Dialog,
    DialogActions,
    DialogContent, FormControl, Grid,
    IconButton, InputLabel,
    MenuItem, Select, TextField,
    Typography, useMediaQuery, useTheme
} from '@mui/material';
import { AccessTime, ArrowDownward, ArrowUpward, PersonPinSharp, } from '@mui/icons-material';
import { DatePicker as KeyboardDatePicker } from '@mui/x-date-pickers/DatePicker';
import { TimePicker as KeyboardTimePicker} from '@mui/x-date-pickers/TimePicker';
import React, { useEffect, useRef, useState } from 'react';
import { putOcc } from '../apicalls/occ';
import {
    dateToEpoch, dayIndex, dfInHours, epochToDate,
    extendedDfInHours, randomString
} from '../helpers';
import { MulwiColors } from '../mulwiColors';
import {
    getErrorDialog, getInfoDialog,
    getNullDialog, StatusDialog
} from '../StatusDialog';
import { getWkFromMonth, MUISwitch, WeekSwitch } from './harmonogram';
import { HarmonogramMonth, sDays } from './month';
import { HarmonogramWeek } from './weekBigRes';
import { getTrainingByID } from '../apicalls/instructor.api';
import { ThemeProvider } from "@mui/styles";
import { DatePickerCustomTheme } from "../helpers"
import { locale2 } from '../locale';
import { HarmonogramDay } from './day';

const sessionColors = [
    MulwiColors.pinkDark,
    MulwiColors.pinkAction,
    MulwiColors.greenLight,
    MulwiColors.blueDark,
    MulwiColors.blueLight,
    "#69385c",
    "#9792e3",
    "#e9df00",
    "#db5aba",
    "#2c1320",
    "#73683b",
    "#ffe66d",
    "#955e42",
    "#40531b",
    "#ca7df9",
    "#896279",
    "#b9e6ff"
]

export function ColorEditor(props) {

    return (<FormControl fullWidth style={props.style}>
        <InputLabel id="cl">{locale2.EMPHASIZE_SESS[props.lang]}</InputLabel>
        <Select
            labelId="cl"
            // open={colorOpen}
            // onClose={() => setColorOpen(false)}
            // onOpen={() => setColorOpen(true)}
            value={props.color}
            fullWidth
            onChange={e => props.setColor(e.target.value)}>
            {/* <MenuItem value="">
                        <em>None</em>
                    </MenuItem> */}
            {sessionColors.map(c => (<MenuItem key={c} value={c}>
                <div style={{ backgroundColor: c, minHeight: 30, width: "100%" }} />
            </MenuItem>
            ))}
        </Select>
    </FormControl>)
}

export function ValueSwitch(props) {
    return (<Grid
        direction="row"
        style={{ cursor: "pointer", fontSize: 11 }}
        component="label" container alignItems="center" >
        <Grid item>{!props.selected ? <strong>{props.left}</strong> : <span>{props.left}</span>}</Grid>
        <Grid item>
            <MUISwitch disableRipple
                checked={props.selected}
                onChange={(e) => props.setSelected(e.target.checked)} />
        </Grid>
        <Grid item>{props.selected ? <strong>{props.right}</strong> : <span>{props.right}</span>}</Grid>
    </Grid>)
}

export function ModifyOcc2Modal(props) {
    return (<Dialog open={props.open} onClose={() => props.setOpen(false)}>
        <DialogContent>
            <OccEdit {...props} />
        </DialogContent>
    </Dialog>)
}

export function OpPicker(props) {

    return (<Dialog open={props.open} onClose={() => props.setOpen(false)}>
        <DialogContent>
            <Grid container direction="column" spacing={2}>
                <Grid item>
                    <Button fullWidth variant="contained" style={{
                        color: "white",
                        backgroundColor: MulwiColors.blueDark
                    }} onClick={() => {
                        props.setOp(1)
                    }}>
                        {locale2.MOVE_THIS_COURSE[props.lang]}
                    </Button>
                </Grid>
                <Grid item>
                    <Button fullWidth variant="contained" style={{
                        color: "white",
                        backgroundColor: MulwiColors.blueDark
                    }} onClick={() => {
                        props.setOp(3)
                    }}>
                        {locale2.CHANGE_SCHEDULE[props.lang]}
                    </Button>
                </Grid>
                <Grid item>
                    <Button style={{
                        color: "white",
                        backgroundColor: MulwiColors.redError
                    }} fullWidth variant="contained" onClick={() => {
                        props.setOp(2)
                    }}>
                        {locale2.DELETE_THIS_COURSE[props.lang]}
                    </Button>
                </Grid>
            </Grid>
        </DialogContent>
        <DialogActions>
            <Button onClick={() => {
                props.setOpen(false)
            }}>
                {locale2.CLOSE[props.lang]}
            </Button>
        </DialogActions>
    </Dialog>)
}

export function MoveOccModal(props) {

    const [start, _setStart] = useState(null)
    const [end, setEnd] = useState(null)
    const df = useRef(null)

    const TextFieldComponent = (props) => {
        return <TextField {...props} disabled={true} />
    }

    function setStart(v) {
        if(!df)
            return
        let _end = new Date(v)
        _end.setMinutes(v.getMinutes() + df.current)
        _setStart(v)
        setEnd(_end)
    }

    useEffect(() => {
        if (!props.occs || props.occs.length == 0) return

        let ixs = props.editIndexes || []
        switch (ixs.length) {
            case 1:
            case 2:
                {
                    let occ = props.occs[ixs[0]]
                    if (!occ) return
                    let start = new Date(occ.DateStart)
                    let end = new Date(occ.DateEnd)
                    _setStart(start)
                    setEnd(end)
                    let dfInMinutes = Math.round(dfInHours(start, end) * 60)
                    df.current = dfInMinutes
                }
                break
        }


    }, [props.editIndexes, props.editorData])

    return (<Dialog open={props.open} onClose={() => props.setOpen(false)}>
        <DialogContent>
            {/* <ThemeProvider theme={DatePickerCustomTheme}>
                <KeyboardDatePicker
                    fullWidth
                    size="small"
                    ampm={false}
                    value={start}
                    onChange={v => {
                        if (v instanceof Date && isFinite(v)) {
                            setStart(v)
                        }
                    }}
                    label={locale2.DATE_OF_STARTING_COURSE[props.lang]}>
                </KeyboardDatePicker>
                <KeyboardTimePicker
                    fullWidth
                    size="small"
                    ampm={false}
                    keyboardIcon={<AccessTime />}
                    label={locale2.TIME_OF_STARTING_COURSE[props.lang]}
                    value={start}
                    onChange={v => {
                        if (v instanceof Date && isFinite(v)) {
                            let d = new Date(start)
                            d.setHours(v.getHours())
                            d.setMinutes(v.getMinutes())
                            d.setSeconds(0)
                            d.setMilliseconds(0)
                            setStart(d)
                        }
                    }} 
                    renderInput={(params) => <TextField {...params} />}
                    />
            </ThemeProvider>
            <ThemeProvider theme={DatePickerCustomTheme}>
                <KeyboardDatePicker
                    readOnly
                    fullWidth
                    size="small"
                    ampm={false}
                    TextFieldComponent={TextFieldComponent}
                    value={end}
                    onChange={v => {
                    }}
                    label={locale2.DATE_OF_ENDING_COURSE[props.lang]} />
                <KeyboardTimePicker
                    readOnly
                    TextFieldComponent={TextFieldComponent}
                    fullWidth
                    size="small"
                    ampm={false}
                    keyboardIcon={<AccessTime />}
                    label={locale2.TIME_OF_ENDING_COURSE[props.lang]}
                    value={end}
                    onChange={v => {
                    }} 
                    renderInput={(params) => <TextField {...params} />}
                    />
            </ThemeProvider> */}
        </DialogContent>
        <DialogActions>
            <Button onClick={() => {
                props.setOpen && props.setOpen(false)
            }}>
                {locale2.CLOSE[props.lang]}
            </Button>
            <Button onClick={() => {
                let eix = props.editIndexes || []
                if (!props.occs || props.occs.length == 0 || eix.length < 1) return

                let occ = props.occs[eix[0]]
                occ.DateStart = start
                occ.DateEnd = end
                props.setOccs([...props.occs])
            }} variant="contained" style={{
                color: "white",
                backgroundColor: MulwiColors.greenDark
            }}>
                {locale2.SAVE[props.lang]}
            </Button>
        </DialogActions>
    </Dialog>)
}

export function adjustOccDatesAfterEdit(occ) {
    let shift = 0
    occ.SecondaryOccs.sort((a, b) => a.OffsetStart - b.OffsetStart)
    let bso = null
    for (let i = 0; i < occ.SecondaryOccs.length; i++) {
        let so = occ.SecondaryOccs[i]
        if (i === 0) {
            shift = so.OffsetStart
            if (shift != 0) {
                occ.DateStart = new Date(occ.DateStart)
                occ.DateStart.setMinutes(occ.DateStart.getMinutes() + shift)
            }
        }
        if (shift != 0) {
            so.OffsetStart -= shift
            so.OffsetEnd -= shift
        }
        if (!bso || bso.OffsetEnd < so.OffsetEnd) bso = so
    }
    if (bso) {
        let x = new Date(occ.DateStart)
        x.setMinutes(x.getMinutes() + bso.OffsetEnd)
        occ.DateEnd = x
    }
}

export function OccEdit(props) {

    const [start, setStart] = useState(null)
    const [endt, setEndt] = useState(false)
    const [color, setColor] = useState(MulwiColors.blueDark)
    const [remarks, setRemarks] = useState("")
    const [dur, setDur] = useState([0, 0])

    useEffect(() => {

        if (!props.occs || props.occs.length == 0) return

        let ixs = props.editIndexes || []
        switch (ixs.length) {
            case 2:
                {
                    let occ = props.occs[ixs[0]]
                    if (!occ) return
                    let occ2 = occ.SecondaryOccs[ixs[1]]
                    if (!occ2) return
                    let start = new Date(occ.DateStart)
                    let end = new Date(start)
                    start.setMinutes(start.getMinutes() + occ2.OffsetStart)
                    end.setMinutes(end.getMinutes() + occ2.OffsetEnd)
                    setStart(start)
                    setRemarks(occ2.Remarks)
                    setColor(occ2.Color)
                    if (start.getDate() != end.getDate()) {
                        setEndt(true)
                        let c = dfInHours(start, end)
                        setDur([c, Math.round((c * 60) % 60)])
                    } else {
                        setEndt(false)
                        let h = end.getHours()
                        let m = end.getMinutes()
                        setDur([h, m])
                    }
                }
                break
        }
    }, [props.editIndexes, props.occs])


    function dtf() {
        if (!dur[0] && !dur[1])
            return null
        let h = dur[0]
        if (h < 10) h = "0" + String(h)
        let m = dur[1]
        if (m < 10) m = "0" + String(m)
        return h + ":" + m
    }

    return (<React.Fragment>
        <Grid container direction="row">
            <Grid item sm={6}>
                {/* <ThemeProvider theme={DatePickerCustomTheme}>
                    <KeyboardTimePicker
                        fullWidth
                        dura
                        size="small"
                        ampm={false}
                        keyboardIcon={<AccessTime />}
                        label={locale2.FROM[props.lang]}
                        value={start}
                        onChange={v => {
                            if (v instanceof Date && isFinite(v)) {
                                let d = new Date(start)
                                d.setHours(v.getHours())
                                d.setMinutes(v.getMinutes())
                                d.setSeconds(0)
                                d.setMilliseconds(0)
                                setStart(d)
                            }
                        }} 
                        renderInput={(params) => <TextField {...params} />}
                        />
                </ThemeProvider> */}
            </Grid>
            <Grid item sm={6}>
                {/* <ThemeProvider theme={DatePickerCustomTheme}> */}
                    {/* <KeyboardTimePicker
                        fullWidth
                        ampm={false}
                        keyboardIcon={<AccessTime />}
                        label={endt ? locale2.DURATION[props.lang] : locale2.TO[props.lang]}
                        //value={dtf()}
                        inputValue={dtf()}
                        size="small"
                        error={false}
                        invalidDateMessage={""}
                        onChange={(_, e) => {
                            if (e) {
                                let splits = e.split(":")
                                if (splits.length === 2) {
                                    let h = Number(splits[0])
                                    let m = Number(splits[1])
                                    if (Number.isNaN(h)) h = 0
                                    if (Number.isNaN(m)) m = 0
                                    setDur([h, m])
                                } else {
                                    setDur([0, 0])
                                }
                            } else {
                                setDur([0, 0])
                            }
                        }} 
                        renderInput={(params) => <TextField {...params} />}
                        /> */}
                {/* </ThemeProvider> */}
            </Grid>
        </Grid>
        <Grid container direction="row">
            <Grid item md={6}>
                <ValueSwitch selected={endt} setSelected={setEndt}
                    left={locale2.TO[props.lang]}
                    right={locale2.DURATION[props.lang]} />
            </Grid>
        </Grid>
        <TextField fullWidth
            multiline
            rows={2}
            value={remarks}
            onChange={e => setRemarks(e.target.value)}
            label={locale2.CLIENT_REMARKS[props.lang]} variant="outlined" size="small" />
        <ColorEditor
            lang={props.lang}
            style={{
                marginBottom: 10
            }}
            color={color} setColor={setColor} />
        <Grid container direction="row" justify="space-between">
            <Grid item>
                <Button onClick={() => {
                    let ixs = props.editIndexes || []
                    if (ixs.length != 2) return
                    let occ = props.occs[ixs[0]]
                    if (!occ) return
                    occ.SecondaryOccs.splice(ixs[1], 1)
                    props.setLastAdded([])
                    if (occ.SecondaryOccs.length === 0) {
                        props.setOccs([])
                        return
                    }
                    for(let i = 0 ; i < props.occs.length; i++) {
                        props.occs[i].SecondaryOccs = occ.SecondaryOccs
                        adjustOccDatesAfterEdit(occ)
                    }
                    props.setOccs([...props.occs])
                }}>
                    {locale2.RM_ELEM[props.lang]}
                </Button>
            </Grid>
            <Grid item>
                <Button onClick={() => {
                    props.setOpen && props.setOpen(false)
                }}>
                    {locale2.CANCEL[props.lang]}
                </Button>
                <Button onClick={() => {
                    if (!start) return
                    if (!dur[0] && !dur[1]) return

                    let ixs = props.editIndexes || []
                    if (ixs.length != 2 || ixs[0].length === 0) return

                    let d = {
                        Color: color,
                        Remarks: remarks
                    }

                    let occ = props.occs[ixs[0]]
                    if (!occ) return
                    let offsetStart = Math.round(dfInHours(occ.DateStart, start) * 60)
                    let offsetEnd = 0

                    if (endt) {
                        offsetEnd = offsetStart + (dur[0] * 60) + dur[1]
                    } else {
                        let d = new Date(start)
                        d.setHours(dur[0])
                        d.setMinutes(dur[1])
                        if (+d <= +start) return
                        let df = dfInHours(start, d)
                        offsetEnd = offsetStart + (df * 60)
                    }

                    occ.DateStart = new Date(occ.DateStart)
                    occ.DateEnd = new Date(occ.DateEnd)

                    let occ2 = {
                        OffsetStart: offsetStart,
                        OffsetEnd: offsetEnd,
                        ...d,
                    }

                    if (!occ.SecondaryOccs) occ.SecondaryOccs = []

                    occ.SecondaryOccs[ixs[1]] = occ2

                    for (let i = 0; i < props.occs.length; i++) {
                        props.occs[i].SecondaryOccs = occ.SecondaryOccs
                        adjustOccDatesAfterEdit(props.occs[i])
                    }

                    props.setLastAdded([occ, occ2])
                    props.setOccs([...props.occs])

                }} variant="contained" style={{
                    color: "white",
                    backgroundColor: MulwiColors.greenDark
                }}>
                    {locale2.SAVE[props.lang]}
                </Button>
            </Grid>
        </Grid>
    </React.Fragment>)
}

export function MonthView(props) {

    function getDays() {
        let ret = []
        let i = 0
        let startDay = start.getDay() - 1 // because js returns 1-based
        if (startDay < 0) startDay += 7
        let df = Math.abs(startDay - i)
        for (let i = 0; i < df; i++) ret.push(-i - 1)
        let startCpy = new Date(start)
        while (startCpy <= end) {
            ret.push(startCpy.getDate())
            startCpy.setDate(startCpy.getDate() + 1)
        }
        let endDay = end.getDay() - 1
        if (endDay < 0) endDay += 7
        df = 0
        if (endDay < 6) {
            df = 6 - endDay
        }
        for (let i = 0; i < df; i++) ret.push(-1)

        // group ret in 7s
        let gret = []
        let ix = 0
        let page = 0
        for (let i = 0; i < ret.length; i++) {
            if (ix === 0) {
                gret.push([])
            }
            gret[page].push(ret[i])
            ix++
            if (ix === 7) {
                ix = 0
                page++
            }
        }
        return gret
    }

    const [sel, setSel] = useState({})

    const [start, setStart] = useState(new Date())
    const [end, setEnd] = useState(new Date())

    useEffect(() => {
        let d = props.date
        if (!d) return
        let _start = new Date(d.getFullYear(), d.getMonth(), 1, 0, 0, 0, 0)
        let _end = new Date(d.getFullYear(), d.getMonth() + 1)
        _end.setHours(0, 0, 0, -1)
        setStart(_start)
        setEnd(_end)
    }, [props.date])

    useEffect(() => {
        let newSel = {}
        let ed = props.editorData
        if (!ed) return
        let occs = ed.occs
        if (!occs) return
        for (let i = 0; i < occs.length; i++) {
            let os = new Date(occs[i].DateStart)
            //let oe = new Date(occs[i].DateEnd)
            let soccs = occs[i].SecondaryOccs
            if (soccs && soccs.length > 0) {
                for (let j = 0; j < soccs.length; j++) {
                    let so = soccs[j]
                    let start2 = new Date(os)
                    let end2 = new Date(os)
                    start2.setMinutes(start2.getMinutes() + so.OffsetStart)
                    end2.setMinutes(end2.getMinutes() + so.OffsetEnd)
                    let _durs = extendedDfInHours(start2, end2)
                    for (let z in _durs) {
                        let tmpstart = epochToDate(z)
                        tmpstart.setHours(0, 0, 0, 0)
                        let key = dateToEpoch(tmpstart)
                        let pi = newSel[key]
                        if (pi) {
                            newSel[key] = {
                                grp: i,
                                occ2ix: Math.min(j, pi.occ2ix)
                            }
                        } else {
                            newSel[key] = {
                                grp: i,
                                occ2ix: j
                            }
                        }
                    }
                }
            } else {
                let occ = occs[i]
                let _durs = extendedDfInHours(new Date(occ.DateStart), new Date(occ.DateEnd))
                for (let z in _durs) {
                    let tmpstart = epochToDate(z)
                    tmpstart.setHours(0, 0, 0, 0)
                    let key = dateToEpoch(tmpstart)
                    newSel[key] = {
                        grp: i,
                        occ2ix: 0
                    }
                }
            }
        }
        setSel(newSel)
    }, [props.editorData])


    function displayDay(n) {
        if (n < 10)
            return "0" + String(n)
        return String(n)
    }

    const today = new Date()

    function shouldApplyBorder(d, grp) {
        let adj = new Date(start)
        adj.setDate(d)
        let adjSel = sel[dateToEpoch(adj)]
        if (!adjSel || adjSel.grp != grp) {
            return true
        }
        return false
    }

    function getPropsForDay(d, isLeft, isRight, isTop, isBottom) {
        let s = {
            padding: 8,
            position: "relative",
            width: 31
        }
        if (d < 0) {
            s.color = "transparent"
        }

        let onClick = null
        let isLesser = false

        if (d >= 0 && d < today.getDate() &&
            today.getMonth() === props.date.getMonth() &&
            today.getFullYear() === props.date.getFullYear()) {
            s.color = "gray"
            isLesser = true
        }

        if (!isLesser) {
            onClick = () => {
                let x = new Date(props.date)
                x.setDate(d)
                x.setHours(0, 0, 0, 0)
                props.onSelectedDate(x)
            }
        }

        if (d > 0) {
            let c = new Date(start)
            c.setDate(d)
            let csel = sel[dateToEpoch(c)]
            if (csel) {
                s.boxShadow = ""
                if (csel.occ2ix === 0) onClick = () => {
                    props.onSelectedDate(null, csel.grp)
                }
                if (isLeft || shouldApplyBorder(d - 1, csel.grp)) {
                    s.boxShadow += "-1px 0 " + MulwiColors.greenDark
                }
                if (isRight || shouldApplyBorder(d + 1, csel.grp)) {
                    if (s.boxShadow) s.boxShadow += ","
                    s.boxShadow += "1px 0 " + MulwiColors.greenDark
                }
                if (isTop || shouldApplyBorder(d - 7, csel.grp)) {
                    if (s.boxShadow) s.boxShadow += ","
                    s.boxShadow += "0 -1px " + MulwiColors.greenDark
                }
                if (isBottom || shouldApplyBorder(d + 7, csel.grp)) {
                    if (s.boxShadow) s.boxShadow += ","
                    s.boxShadow += "0 1px " + MulwiColors.greenDark
                }
            }
        }

        return {
            style: s,
            onClick: onClick
        }
    }

    function adornment(d) {
        if (d <= 0) return
        let c = new Date(start)
        c.setDate(d)
        let csel = sel[dateToEpoch(c)]
        if (!csel || csel.occ2ix != 0) return
        return (<ButtonBase style={{
            position: "absolute", color: "white", fontSize: 12,
            left: 0, top: 0, width: 10, height: 10,
            backgroundColor: MulwiColors.greenDark
        }}>
            x
        </ButtonBase>)
    }

    let days = getDays()

    return (<React.Fragment>

        <div style={{ maxWidth: 220 }}>

            {props.withTitle && (<div style={{ display: "inline-block", marginLeft: 10, marginBottom: 10 }}>
                {[0, 1, 2, 3, 4, 5, 6].map(d => (
                    <span key={d} style={{ marginRight: 14 }} >
                        <strong>{dayIndex[d + 1][props.lang]}</strong>
                    </span>
                ))}
            </div>)}
            <center>
                <Typography variant="body2">
                    <strong>{props.date.toLocaleString(navigator.language, { month: "long" }).toUpperCase()}</strong>
                </Typography>
            </center>

            {/* {shouldGeneratePrefix(days) && arrayRepeat(7, -1).map((d, i) => (
                <React.Fragment key={i}>
                    {generateCalendarPaddingItem()}
                </React.Fragment>
            ))} */}

            {days.map((gd, i) => {
                return (<React.Fragment key={i}>
                    <Grid container direction="row">
                        {gd.map((d, j) => {
                            let isLeft = true
                            let isRight = true
                            if (gd[j - 1] && gd[j - 1] > 0) isLeft = false
                            if (gd[j + 1] && gd[j + 1] > 0) isRight = false
                            let isTop = true
                            let isBottom = true
                            if (days[i - 1] && days[i - 1][j] && days[i - 1][j] > 0) isTop = false
                            if (days[i + 1] && days[i + 1][j] && days[i + 1][j] > 0) isBottom = false
                            return (<Grid item key={j}>
                                <ButtonBase disabled={d < 0} {...getPropsForDay(d, isLeft, isRight, isTop, isBottom)}>
                                    {(d < 0 && "00") || displayDay(d)}
                                    {adornment(d)}
                                </ButtonBase>
                            </Grid>
                            )
                        })}
                    </Grid>
                </React.Fragment>
                )
            })}
        </div>
    </React.Fragment>)
}

export function CalendarHeader(props) {
    if (!props.drawerData || !props.drawerData.training) return null
    return (
        <React.Fragment>
            <Grid item>
                {/* <Typography variant="h5">{returnLocaleString(['harmonogram', 'calendarEditor'])[9]} <strong>{props.drawerData.training.Title}</strong></Typography> */}
                <Typography variant="h5"><strong>{props.drawerData.training.Title}</strong></Typography>
            </Grid>
            <Grid item>
                <Button variant="contained" style={{
                    color: "white",
                    backgroundColor: MulwiColors.greenDark
                }} onClick={() => {
                    props.setSaveToken && props.setSaveToken(1)
                }}>
                    {locale2.SAVE[props.lang]}
                </Button>
            </Grid>
            <Grid item>
                <Button variant="contained" style={{
                    backgroundColor: "inherit",
                    color: "black",
                }} onClick={props.onClose}>
                    {locale2.CLOSE[props.lang]}
                </Button>
            </Grid>
        </React.Fragment>)
}


function RptView(props) {

    // const [months, setMonths] = useState(() => {
    //     let now = new Date()
    //     let next = new Date(now)
    //     next.setMonth(next.getMonth() + 1)
    //     let next2 = new Date(next)
    //     next2.setMonth(next.getMonth() + 1)
    //     return [now, next, next2]
    // })

    // function getOccs() {
    //     if (!props.editorData) return []
    //     let occs = props.editorData.occs
    //     if (!occs || occs.length === 0) return []
    //     return occs
    // }

    const started = useRef(false)

    useEffect(() => {

        if (!props.editorData) return
        let occs = props.editorData.occs
        if (!occs || occs.length === 0) return

        if (started.current)
            return
        started.current = true


        if (occs.length > 1) {
            for (let i = 0; i < occs.length; i++) {
                if (occs[i].RepeatDays) {
                    props.setInfo(getInfoDialog(locale2.WARNING[props.lang], locale2.EDITOR_COMPAT_WARN[props.lang]))
                }
            }
            props.setRepeating(-1)
            return
        }
        let occ = occs[0]
        switch (occ.RepeatDays) {
            case 7:
                props.setRepeating(7)
                break
            case 1:
                props.setRepeating(1)
                break
        }
    }, [props.editorData])

    // function rotate(down) {
    //     if (down) {
    //         let x = new Date(months[2])
    //         x.setMonth(x.getMonth() + 1)
    //         setMonths([months[1], months[2], x])
    //     } else {
    //         let x = new Date(months[0])
    //         x.setMonth(x.getMonth() - 1)
    //         let today = new Date()
    //         // if (x < today && x.getMonth() != today.getMonth()) {
    //         //     return
    //         // }
    //         setMonths([x, months[0], months[1]])
    //     }
    // }

    // const [fallbackOcc, setFallbackOcc] = useState(null)

    // function copyOccToDay(occ, d) {
    //     let df = dfInHours(occ.DateStart, occ.DateEnd)
    //     let ds = new Date(occ.DateStart)
    //     ds.setYear(d.getFullYear())
    //     ds.setMonth(d.getMonth())
    //     ds.setDate(d.getDate())
    //     let de = new Date(ds)
    //     de.setHours(de.getHours() + df)
    //     let occNew = { ...occ }
    //     occNew.DateStart = ds
    //     occNew.DateEnd = de
    //     return occNew
    // }

    // function setCustomRpt(pv) {
    //     if (pv === -1) return
    //     let occs = getOccs()
    //     if (occs.length > 1) {
    //         props.setInfo(getInfoDialog(locale2.WARNING[props.lang], locale2.EDITOR_COMPAT_WARN[props.lang]))
    //         return
    //     }
    //     if (occs.length == 0) {
    //         props.setRepeating(-1)
    //         return
    //     }
    //     let occ = occs[0]
    //     occ.RepeatDays = 0
    //     occs = [occ]
    //     let ed = { ...props.editorData }
    //     ed.occs = occs
    //     props.setEditorData(ed) 
    //     props.setRepeating(-1)
    // }

    // function setAutoRpt(_, v) {
    //     let occs = getOccs()
    //     if (occs.length < 1) return
    //     let occ = occs[0]
    //     occ.RepeatDays = v
    //     occs = [occ]
    //     let ed = { ...props.editorData }
    //     ed.occs = occs
    //     props.setEditorData(ed)
    // }

    return (
        <React.Fragment>
            <Grid container direction="row" spacing={2} style={{ marginTop: 4 }}>
                <Grid item>
                    <Typography variant="h5">
                        {locale2.REPEATING[props.lang]}:
                    </Typography>
                </Grid>
                <Grid item>
                    <Select value={props.repeating} onChange={e => {
                        let v = e.target.value
                        // let pv = props.repeating
                        // switch (v) {
                        //     case -1:
                        //         setCustomRpt(pv)
                        //         break
                        //     default:
                        //         if (v < 0) return
                        //         setAutoRpt(pv, v)
                        //         break
                        // }
                        props.setRepeating(v)
                    }}>
                        <MenuItem value={0}>{locale2.DONT_REPEAT[props.lang]}</MenuItem>
                        <MenuItem value={1}>{locale2.DAILY[props.lang]}</MenuItem>
                        <MenuItem value={7}>{locale2.WEEKLY[props.lang]}</MenuItem>
                        <MenuItem value={-1}>{locale2.NONSTANDARD[props.lang]}</MenuItem>
                    </Select>
                </Grid>
            </Grid>
        </React.Fragment>)
}

export function prepOccs(occs) {
    if (!occs) {
        return
    }
    for (let i = 0; i < occs.length; i++) {
        occs[i].DateStart = new Date(occs[i].DateStart)
        occs[i].DateEnd = new Date(occs[i].DateEnd)
        let soccs = occs[i].SecondaryOccs
        if (soccs && soccs.length > 1) {
            soccs.sort((a, b) => a.OffsetStart - b.OffsetStart)
        }
        if (!soccs || soccs.length == 0) {
            let so = { ...occs[i] }
            so.OffsetStart = 0
            so.OffsetEnd = dfInHours(occs[i].DateStart, occs[i].DateEnd) * 60
            occs[i].SecondaryOccs = [so]
        }
    }
    occs.sort((a, b) => +a.DateStart - +b.DateStart)
}

export function fitDateToOcc(occs, setWeek) {
    if (occs && occs.length > 0) {
        setWeek(getWkFromMonth(new Date(occs[0].DateStart)))
    }
}


export function CalendarEditor(props) {

    const [waiting, setWaiting] = useState(false)
    const [info, setInfo] = useState(getNullDialog())
    const [week, setWeek] = useState(() => getWkFromMonth(new Date()))

    const [editorData, _setEditorData] = useState({
        occs: null,
        training: null
    })

    function _fitDateToOcc(occs) {
        fitDateToOcc(occs, setWeek)
    }

    function setEditorDataWithAlign(ed) {
        // let occs = ed.occs
        // // fix alignment for copies
        // if (occs.length > 1) {
        //     let occ = occs[0]
        //     occ.DateStart = new Date(occ.DateStart)
        //     occ.DateEnd = new Date(occ.DateEnd)
        //     let hr0 = occ.DateStart.getHours()
        //     let m0 = occ.DateStart.getMinutes()
        //     let hre0 = occ.DateEnd.getHours()
        //     let me0 = occ.DateEnd.getMinutes()
        //     for (let i = 1; i < occs.length; i++) {
        //         let x = occs[i]
        //         occs[i].SecondaryOccs = occ.SecondaryOccs
        //         x.DateEnd = new Date(x.DateEnd)
        //         x.DateStart = new Date(x.DateStart)
        //         x.Remarks = occ.Remarks
        //         x.Color = occ.Color
        //         let hr = x.DateStart.getHours()
        //         let hre = x.DateEnd.getHours()
        //         let m = x.DateStart.getMinutes()
        //         let me = x.DateEnd.getMinutes()
        //         if (hr != hr0) {
        //             x.DateStart.setHours(hr0)
        //         }
        //         if (m != m0) {
        //             x.DateStart.setMinutes(m0)
        //         }
        //         if (hre != hre0) {
        //             x.DateEnd.setHours(hre0)
        //         }
        //         if (me != me0) {
        //             x.DateEnd.setMinutes(me0)
        //         }
        //     }
        //     // _fitDateToOcc(occs)
        // }
        _setEditorData(ed)
    }

    const [repeating, _setRepeating] = useState(0)

    function setAutoRepeating(occs, rpt) {
        if (occs.length == 0)
            return
        occs = [editorData.occs[0]]
        occs[0].RepeatDays = rpt
        editorData.occs = occs
    }

    function setRepeating(v) {
        switch (v) {
            case 0:
                setAutoRepeating(editorData.occs, 0)
                break
            case 1:
                setAutoRepeating(editorData.occs, 1)
                break
            case 7:
                setAutoRepeating(editorData.occs, 7)
                break
            case -1:
                for (let i = 0; i < editorData.occs.length; i++) {
                    let occ = editorData.occs[i]
                    occ.RepeatDays = 0
                }
        }
        _setRepeating(v)
        _setEditorData({ ...editorData })
    }


    useEffect(() => {
        if (!props.drawerData || !props.drawerData.occs) return
        let occs = JSON.parse(JSON.stringify(props.drawerData.occs))
        prepOccs(occs)
        _fitDateToOcc(occs)
        _setEditorData({
            occs: occs,
            training: props.drawerData.training
        })
    }, [])

    async function refresh() {
        if (!props.drawerData || !props.drawerData.training) return
        let tid = props.drawerData.training.ID
        try {
            let t = await getTrainingByID(tid)
            let dd = { ...props.drawerData }
            dd.training = t[0].Training
            let occs = t[0].Occurrences
            prepOccs(occs)
            _fitDateToOcc(occs)
            dd.occs = occs
            props.setDrawerData(dd)
            _setEditorData({
                occs: occs,
                training: dd.training
            })
        } catch (ex) {
            setInfo(getErrorDialog(ex))
        }
    }

    async function save() {
        if (!props.drawerData || !props.drawerData.training) return
        try {
            setWaiting(1)
            let occs = editorData.occs || []
            for (let i = 0; i < occs.length; i++) {
                if (!occs[i].SecondaryOccs) occs[i].SecondaryOccs = []
                if (occs[i].SecondaryOccs.length == 1) {
                    let so = occs[i].SecondaryOccs[0]
                    occs[i].Remarks = so.Remarks
                    occs[i].Color = so.Color
                    occs[i].SecondaryOccs = []
                }
            }
            let req = {
                TrainingID: props.drawerData.training.ID,
                Occurrences: occs
            }
            await putOcc(req)
            refresh()
            setWaiting(0)
            props.close()
        } catch (ex) {
            setInfo(getErrorDialog(locale2.COULDNT_SAVE_OCC[props.lang], ex))
            setWaiting(0)
        }
    }

    useEffect(() => {
        if (!props.saveToken) return
        props.setSaveToken(0)
        save()
    }, [props.saveToken])

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))
    const [monthDate] = useState(new Date())
    const [day, setDay] = useState(new Date())

    if (!props.drawerData) return null

    return (<React.Fragment>
        <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
        <Backdrop
            sx={{ color: '#fff', zIndex: 9999999 }}
            style={{
                zIndex: 9999999,
                opacity: 0.4
            }}
            open={waiting === 1}>
            <CircularProgress style={{
                color: MulwiColors.blueDark.greenDark
            }} />
        </Backdrop>
        <Grid container direction="column" alignItems="center" style={{
            marginLeft: 10,
        }}>
            <Grid item>
                <Grid container direction="row" spacing={2} style={{ marginBottom: 10 }}>
                    {isLowRes ? (
                        <React.Fragment>
                            <HarmonogramMonth
                                day={day} setDay={setDay}
                                date={monthDate}
                                lang={props.lang}
                                editorData={editorData} />
                            <Grid item style={{
                                width: "95%"
                            }}>
                                <HarmonogramDay
                                    day={day}
                                    lang={props.lang}
                                    repeating={repeating}
                                    setInfo={setInfo}
                                    setEditorData={setEditorDataWithAlign}
                                    editorData={editorData} />
                            </Grid>
                        </React.Fragment>
                    ) : (<React.Fragment>
                        <Grid item>
                            <Grid container direction="row">
                                <Grid item>
                                    <WeekSwitch
                                        lang={props.lang}
                                        week={week} setWeek={setWeek} />
                                </Grid>
                                <Grid item>
                                    <RptView
                                        lang={props.lang}
                                        setInfo={setInfo}
                                        repeating={repeating}
                                        setRepeating={setRepeating}
                                        editorData={editorData} />
                                </Grid>
                            </Grid>
                            <HarmonogramWeek
                                lang={props.lang}
                                setInfo={setInfo}
                                setEditorData={setEditorDataWithAlign}
                                repeating={repeating}
                                editorData={editorData}
                                week={week} />
                        </Grid>
                    </React.Fragment>)}
                </Grid>
            </Grid>
        </Grid>
    </React.Fragment>)
}

export function cpOcc(occ) {
    let cpy = { ...occ }
    cpy.SecondaryOccs = []
    if (occ.SecondaryOccs) for (let i = 0; i < occ.SecondaryOccs.length; i++) {
        cpy.SecondaryOccs.push({ ...occ.SecondaryOccs[i] })
    }
    return cpy
}

// adds new occ/occ2 into occs and returns new array along with edited indexes
export function addOcc2(occs, lastAdded, relDate, repeating) {
    let dateStart = new Date(relDate)
    // dateStart.setHours(0, 0, 0, 0)
    if (occs && occs.length > 0) {
        let occ = occs[0]
        let df = dfInHours(occ.DateStart, dateStart) * 60
        let so = null
        let fix = randomString(12)
        if (lastAdded && lastAdded.length === 2) {
            let o = lastAdded[0]
            so = { ...lastAdded[1] }
            let pd = new Date(o.DateStart)
            pd.setMinutes(pd.getMinutes() + so.OffsetStart)
            if (dateStart.toDateString() === pd.toDateString()) {
                so.OffsetStart = df
                so.OffsetEnd = df + 120
            } else {
                let startCpy = new Date(dateStart)
                startCpy.setHours(pd.getHours(), pd.getMinutes(), pd.getSeconds())
                let diffInDays = Math.round(dfInHours(pd, startCpy) / 24)
                so.OffsetStart += diffInDays * 1440
                so.OffsetEnd += diffInDays * 1440
            }
        } else 
        {
            so = {
                OffsetStart: df,
                OffsetEnd: df + 120,
                Color: occ.Color,
                Remarks: ""
            }
        }
        so.fix = fix
        occ.SecondaryOccs.push(so)
        adjustOccDatesAfterEdit(occ)
        let ix = occ.SecondaryOccs.findIndex(c => c.fix === fix)
        delete so.fix
        return [occs, [0, ix]]
    } else {
        let dateEnd = new Date(dateStart)
        dateEnd.setMinutes(120)
        let _occs = [{
            DateStart: dateStart,
            DateEnd: dateEnd,
            Color: MulwiColors.blueDark,
            Remarks: "",
            Repeating: repeating || 0,
            SecondaryOccs: [{
                OffsetStart: 0,
                OffsetEnd: 120,
                Color: "",
                Remarks: ""
            }]
        }]
        return [_occs, [0, 0]]
    }
}
