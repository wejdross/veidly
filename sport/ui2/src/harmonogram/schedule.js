import { isOptionGroup } from '@mui/base'
import { LocalConvenienceStoreOutlined } from '@mui/icons-material'
import Edit from '@mui/icons-material/Edit'
import {
    Box, Button, Divider,
    Grid, Menu, MenuItem, Paper, Tooltip, Typography,
    useTheme
} from '@mui/material'
import { makeStyles, withStyles } from "@mui/styles"
import React, { useEffect, useRef, useState } from 'react'
import Draggable from 'react-draggable'
import {
    getRsvSchedule, getSchedule,
    getUserSchedule, getVacations
} from '../apicalls/instructor.api'
import {
    avReasonToStr, dayIndex, dfInHours, epochToDate,
    extendedDfInHours, randomString, rmtoken, trainingResToDrawerData
} from '../helpers'
import { locale2 } from '../locale'
import { MulwiColors } from '../mulwiColors'
import { AddTraining } from './addTraining'
import { addOcc2, adjustOccDatesAfterEdit, ModifyOcc2Modal, MoveOccModal, OpPicker } from './calendarEditor'
import { sDays } from './month'
import { ScheduleEditModal } from './scheduleEditModal'
import { prettyPrintDate } from './trainingDetails'

const b = { 0: [], 1: [], 2: [], 3: [], 4: [], 5: [], 6: [], 999: [] }

const HtmlTooltip = withStyles((theme) => ({
    tooltip: {
        backgroundColor: '#f5f5f9',
        color: 'rgba(0, 0, 0, 0.87)',
        maxWidth: 220,
        fontSize: theme.typography.pxToRem(12),
        border: '1px solid #dadde9',
    },
}))(Tooltip);


// ex. if hd = 12.5 ...
export function hrDecToLabel(hd) {
    //    = 12
    let h = Math.floor(hd)
    //    = 0.5
    let r = hd - h
    let ret = ""
    let min = Math.abs(String(Math.round(r * 60, 2)))
    if (min >= 60) {
        min = 0
        h++
    }
    if (h < 10) {
        ret += "0" + h
    } else {
        ret = h
    }
    ret += ":"
    if (min < 10) {
        ret += "0" + min
    } else {
        ret += min
    }
    // ... return 12:30
    return ret
}

const useStyles = makeStyles({
    cardItem: {
        "&:hover": {
            zIndex: 999
        }
    },
})

export function Schedule(props) {

    const cls = useStyles()

    // height multiplier. 
    // effective content height will be that * 24
    let _hm = props.sm ? 20 : 32
    const hm = _hm

    // content will be shifted to bottom by this offset
    // if set to 0 content will be overlapping title bar (days and divider)
    const off = 35
    const [soff, setSoff] = useState(off)

    // column width
    // it no longer will represent pixel width,
    // instead its maximum allowed % width
    let columnItemWidth = 90 // props.sm ? 80 : 200
    if(props.user)
        columnItemWidth = 100 // no need to shift items for readonly view

    const [mday, setMday] = useState(false)

    function getDays() {
        let d = new Date(props.week.start)
        let ret = []
        let ix = 0
        while (d <= props.week.end) {
            let dx = d.getDay() - 1
            if (dx < 0) dx += 7
            ix = sDays[dx]
            ret.push(dayIndex[ix + 1][props.lang])
            d.setDate(d.getDate() + 1)
        }
        return ret
    }

    const theme = useTheme()
    const gray = theme.palette.grey[200]
    const darkgray = theme.palette.grey[400]

    function getHours() {
        let r = []
        for (let i = 0; i <= 24; i += 2) {
            if (i < 10) {
                r.push({ h: i, s: "0" + String(i) + ":00" })
            } else {
                r.push({ h: i, s: String(i) + ":00" })
            }
        }
        return r
    }

    function bcpy(b) {
        return {
            999: [...b[999]],
            0: [...b[0]],
            1: [...b[1]],
            2: [...b[2]],
            3: [...b[3]],
            4: [...b[4]],
            5: [...b[5]],
            6: [...b[6]],
        }
    }

    const [items, setItems] = useState(bcpy(b))

    function percentToDigit(input) {
        if(!input)
            return null
        if(!input.endsWith)
            return null
        if(!input.endsWith("%"))
            return null
        input = input.replace("%", "")
        return Number(input)
    }

    /* this s*** needs to be reworked */
    /* edit: this is perfectly fine */
    function setOverlaps(o) {

        const marginLeftOffset = 3
        for (let k in o) {
            let d = o[k]
            if (!d) continue

            let removedEls = []

            for (let i = d.length - 1; i >= 0; i--) {
                let el = d[i]
                if (!el || el.backgroundItem) {
                    removedEls.push(el)
                    d.splice(i, 1)
                    continue
                }
                if (!el.style) continue
                let p = percentToDigit(el.style.width)
                if(p) {
                    el.style.width = p
                    el.__convertedFromPercent = true
                }
            }

            for (let i = 0; i < d.length; i++) {
                let el = d[i]
                if (!el.style || el.mday || !el.t) continue
                let yStart = el.style.marginTop
                let yEnd = yStart + el.style.height
                let xStart = el.style.marginLeft || 0
                let xEnd = xStart + el.style.width
                for (let z = 0; z < d.length; z++) {
                    if (z === i) continue
                    let zel = d[z]
                    if (!zel.style || zel.mday || !zel.t) continue
                    let yStart2 = zel.style.marginTop
                    let yEnd2 = yStart2 + zel.style.height
                    let xStart2 = zel.style.marginLeft || 0
                    let xEnd2 = xStart2 + zel.style.width
                    let yOverlap = yEnd > yStart2 && yStart < yEnd2
                    let xOverlap = xEnd > xStart2 && xStart < xEnd2
                    if (yOverlap) {
                        if (xOverlap) {
                            zel.style.marginLeft = xEnd + marginLeftOffset
                            if (!zel._shiftRight) zel._shiftRight = 1
                            else zel._shiftRight++
                        }
                        if (!el._overlaps) el._overlaps = []
                        el._overlaps.push(z)
                    }
                }
            }
            let ovps = {}
            for (let i = 0; i < d.length; i++) {
                let el = d[i]
                if (!el._overlaps) continue
                let shifts = el._shiftRight || 0
                if (!ovps[shifts]) ovps[shifts] = [i]
                else ovps[shifts].push(i)
            }
            let keys = Object.keys(ovps)
            keys.sort((a, b) => b - a)
            for (let i = 0; i < keys.length; i++) {
                let elems = ovps[keys[i]]
                for (let j = 0; j < elems.length; j++) {
                    let target = d[elems[j]]
                    if (!target._overlaps) continue
                    let shift = target._shiftRight
                    if (!target._dv || target._dv < shift) target._dv = shift + 1
                    for (let z = 0; z < target._overlaps.length; z++) {
                        let ovElem = d[target._overlaps[z]]
                        if (!ovElem._dv || ovElem._dv < target._dv) ovElem._dv = target._dv
                    }
                }
            }
            for (let i = 0; i < d.length; i++) {
                let el = d[i]
                if (!el._overlaps) continue
                el.style.marginLeft /= el._dv
                el.style.width /= el._dv
                el.style.marginLeft = Math.round(el.style.marginLeft)
                el.style.width = Math.round(el.style.width)
            }
            for (let i = 0; i < d.length; i++) {
                let el = d[i]
                if (!el._overlaps) continue
                let minMl = 99999
                for (let j = 0; j < el._overlaps.length; j++) {
                    let ml = d[el._overlaps[j]].style.marginLeft - marginLeftOffset
                    if (ml < minMl) minMl = ml
                }
                if (minMl != 99999 && minMl > el.style.width) {
                    el.style.width = minMl
                }
            }
            for (let i = 0; i < d.length; i++) {
                let el = d[i]
                if (!el._overlaps) continue
                let f = false
                let right = el.style.marginLeft + el.style.width
                for (let j = 0; j < el._overlaps.length; j++) {
                    let ml = d[el._overlaps[j]].style.marginLeft
                    if (ml > right) {
                        f = true
                        break
                    }
                }
                if (f) continue
                el.style.width = columnItemWidth - el.style.marginLeft - marginLeftOffset
            }


            for (let i = 0; i < d.length; i++) {
                let el = d[i]
                if (!el.style || !el.__convertedFromPercent) continue
                let w = el.style.width
                el.style.width = w + "%"
            }

            for (let i = 0; i < d.length; i++) {
                let el = d[i]
                let left = el.style.marginLeft
                if(left)
                    el.style.marginLeft = left + "%"
            }

            for (let i = 0; i < removedEls.length; i++) {
                d.push(removedEls[i])
            }

        }

        return o
    }

    // this will be set to false once unmounted
    const mountedRef = useRef(true)
    useEffect(() => {
        return () => {
            mountedRef.current = false
        }
    }, [])

    function dayBlockStyle(hr, hhr, _off) {
        return {
            position: "absolute",
            height: Math.round(hhr * hm),
            width: columnItemWidth + "%",
            marginLeft: 0,
            color: "white",
            textAlign: "left",
            cursor: "pointer",
            marginTop: Math.round((_off || off) + (hr * hm) - 10),
            overflow: "hidden",
            textOverflow: "ellipsis",
            wordBreak: "break-word",
            zIndex: 1,
            backgroundColor: MulwiColors.blueDark
        }
    }

    function makeStyleMultiDay(style, start, end, numberDays) {
        if (start < props.week.start) {
            start = props.week.start
        }
        // adjust left margin for start hr
        // and set element to span multiple days
        let _hr = start.getHours() + (start.getMinutes() / 60)

        style.marginLeft = (_hr / 24) * 100
        if (end > props.week.end) {
            end = props.week.end
        }
        _hr = end.getHours() + (end.getMinutes() / 60)
        let mright = (1 - (_hr / 24)) * 100

        // correction for borders
        //mright -= (numberDays - 1) * 3

        style.width = 100 * numberDays - style.marginLeft - mright

        style.width += "%"
    }

    function schDayBlockStyle(baseStyle, sch, durs) {
        baseStyle.backgroundColor = (sch.Occ && sch.Occ.Color) || MulwiColors.blueDark

        if (durs) {
            makeStyleMultiDay(baseStyle, sch.Start, sch.End, durs.length)
            if (sch.IsAvailable) {
                baseStyle.border = "solid 1px " + MulwiColors.blueDark
            } else {
                baseStyle.border = "solid 1px " + gray
            }
        }

        baseStyle.borderRadius = 3

        return baseStyle
    }

    function setEditorElements(a) {
        if (!props.editorData || !props.editorData.occs) return
        let occs = props.editorData.occs

        let _mday = false

        if (!occs) {
            setMday(_mday)
            return
        }

        for (let i = 0; i < occs.length; i++) {
            let o = occs[i]

            if (props.editIndexes && props.editIndexes.length > 1 && i != props.editIndexes[0])
                continue

            if (o.SecondaryOccs.length > 1) {
                _mday = true

                for (let j = 0; j < o.SecondaryOccs.length; j++) {
                    let so = o.SecondaryOccs[j]
                    let start2 = new Date(o.DateStart)
                    let end2 = new Date(o.DateStart)
                    start2.setMinutes(start2.getMinutes() + so.OffsetStart)
                    end2.setMinutes(end2.getMinutes() + so.OffsetEnd)

                    if (start2 < props.week.start || end2 > props.week.end)
                        continue
                    if (start2 < props.week.start)
                        start2 = props.week.start
                    if (end2 > props.week.end)
                        end2 = props.week.end

                    let _durs = extendedDfInHours(start2, end2)
                    for (let z in _durs) {
                        let tmpstart = epochToDate(z)
                        let hr = tmpstart.getHours() + (tmpstart.getMinutes() / 60)
                        let day = tmpstart.getDay() - 1
                        if (day < 0) day += 7
                        let s = dayBlockStyle(hr, _durs[z])
                        s.height -= 2
                        s.borderRadius = 3
                        s.backgroundColor = (o.Color) || MulwiColors.blueDark
                        a[day].push({
                            style: s,
                            noEdit: true,
                            editIndexes: props.mdayEdit ? [i, j] : null,
                            occ2: so
                        })
                    }
                }


                let start = new Date(o.DateStart)
                let end = new Date(o.DateEnd)

                if (start < props.week.start)
                    start = props.week.start
                if (end > props.week.end)
                    end = props.week.end
                if (start >= end)
                    continue

                if (!props.mdayEdit) {

                    let _durs = extendedDfInHours(start, end)
                    for (let z in _durs) {
                        let tmpstart = epochToDate(z)
                        let hr = tmpstart.getHours() + (tmpstart.getMinutes() / 60)
                        let day = tmpstart.getDay() - 1
                        if (day < 0) day += 7
                        let s = dayBlockStyle(hr, _durs[z])
                        s.borderRadius = 3
                        s.backgroundColor = "transparent"
                        a[day].push({
                            style: s,
                            backgroundItem: true,
                            noBoxShadow: true,
                            t: props.editorData.training,
                            editIndexes: [i, 0]
                        })
                    }
                }

            }

            let start = new Date(o.DateStart)
            let end = new Date(o.DateEnd)

            if (start < props.week.start)
                start = props.week.start
            if (end > props.week.end)
                end = props.week.end
            if (start >= end)
                continue

            let _durs = extendedDfInHours(start, end)
            for (let z in _durs) {
                let tmpstart = epochToDate(z)
                let hr = tmpstart.getHours() + (tmpstart.getMinutes() / 60)
                let day = tmpstart.getDay() - 1
                if (day < 0) day += 7
                let s = dayBlockStyle(hr, _durs[z])
                s.height -= 2
                s.borderRadius = 3
                let backgroundItem = false
                if (o.SecondaryOccs.length > 1) {
                    s.backgroundColor = gray
                    backgroundItem = true
                    s.zIndex = 0
                } else {
                    s.backgroundColor = (o.Color) || MulwiColors.blueDark
                }
                a[day].push({
                    style: s,
                    backgroundItem: backgroundItem,
                    noBoxShadow: backgroundItem ? true : false,
                    noEdit: true,
                    editIndexes: backgroundItem ? null : [i, 0],
                    occ2: o.SecondaryOccs && o.SecondaryOccs.length === 1 && o.SecondaryOccs[0]
                })
            }

        }

        setMday(_mday)
    }

    function verticalBarObject(primary, offset) {
        return ({
            style: {
                position: "absolute",
                backgroundColor: primary ? MulwiColors.redError : MulwiColors.subtitleTypography,
                height: 2,
                width: 'calc(100% + 4px)',
                marginLeft: -2,
                marginTop: offset,
            },
            backgroundItem: true
        })
    }

    async function refreshHarmonogram() {
        let a = bcpy(b)

        setEditorElements(a)

        if (props.editorData) {
            if (!mountedRef.current) return null
            a = setOverlaps(a)
            setItems(a)
            return
        }

        let s = []
        let vacs = []

        try {
            if (!mountedRef.current) return null
            if (props.usrRsv) {
                s = await getRsvSchedule(
                    props.week.start,
                    props.week.end,
                    props.trainingID)
            } else {
                if (props.user) {
                    if (!props.instructorID) return null
                    s = await getUserSchedule(
                        props.week.start,
                        props.week.end,
                        props.instructorID,
                        props.trainingID,
                        props.smID)
                } else {
                    s = await getSchedule(
                        props.week.start,
                        props.week.end,
                        props.trainingID)
                }
            }
            if (!s) return null
            if (!mountedRef.current) return null
            vacs = await getVacations(props.instructorID)
        } catch (ex) {
            if (!mountedRef.current) return null
            props.setInfo && props.setInfo({
                open: true,
                hdr: locale2.SOMETHING_WENT_WRONG[props.lang],
                msg: ex,
                buttons: (
                    <React.Fragment>
                        <Button onClick={() => {
                            rmtoken()
                            window.location = "/"
                        }} color="primary">
                            {locale2.RESET[props.lang]}
                        </Button>
                    </React.Fragment>
                )
            })
            return
        }

        let now = new Date()
        let nowday = now.getDay() - 1
        if (nowday < 0) nowday += 7
        let nowEnd = new Date(now)
        nowEnd.setMinutes(nowEnd.getMinutes() + 30)

        let _soff = off

        const shift = hm / 2

        let hr = 0
        for (let i = 0; i < s.length; i++) {
            let t = s[i].Training

            //let found = false
            for (let j = 0; j < s[i].Schedule.length; j++) {
                let sch = s[i].Schedule[j]
                let start = sch.Start
                if (start < props.week.start) {
                    start = props.week.start
                }
                let end = sch.End
                if (end > props.week.end) {
                    end = props.week.end
                }
                let _durs = extendedDfInHours(start, end)
                let durs = []
                for (let z in _durs) {
                    durs.push({
                        epoch: z,
                        val: _durs[z]
                    })
                }
                let day = start.getDay() - 1
                if (day < 0) day += 7

                // multiple days only
                if (durs.length < 2) {
                    continue
                }

                if (!sch.IsAvailable && !props.showUnavailable)
                    continue

                if (props.hideMultiday) continue

                // this layout is like an ogre
                // it has layers
                // ha ha

                // let carpetStyle = schDayBlockStyle(dayBlockStyle(hr, 1, off+shift), sch, durs)
                // carpetStyle.borderRadius = 3
                // carpetStyle.filter = "brightness(1.3)"
                // //carpetStyle.opacity = "50%"
                // carpetStyle.zIndex = 18
                // a[day].push({
                //     style: carpetStyle,
                //     mday: true
                // })

                if (!sch.Occ.SecondaryOccs || sch.Occ.SecondaryOccs.length == 0) {
                    let so = { ...sch.Occ }
                    so.OffsetStart = 0
                    so.OffsetEnd = dfInHours(start, end) * 60
                    sch.Occ.SecondaryOccs = [so]
                }

                for (let j = 0; j < sch.Occ.SecondaryOccs.length; j++) {

                    let so = sch.Occ.SecondaryOccs[j]
                    let start2 = new Date(sch.Start)
                    let end2 = new Date(sch.Start)
                    start2.setMinutes(start2.getMinutes() + so.OffsetStart)
                    end2.setMinutes(end2.getMinutes() + so.OffsetEnd)
                    if (start2 < props.week.start || end2 > props.week.end)
                        continue
                    if (start2 < props.week.start)
                        start2 = props.week.start
                    if (end2 > props.week.end)
                        end2 = props.week.end
                    let _durs = extendedDfInHours(start2, end2)
                    let len = 0
                    for (let _ in _durs) len++

                    let day2 = start2.getDay() - 1
                    if (day2 < 0) day2 += 7

                    let secondaryOccStyle = dayBlockStyle(hr, 1, off)
                    makeStyleMultiDay(secondaryOccStyle, start2, end2, len)
                    secondaryOccStyle.zIndex = 19
                    if (j == 0) {
                        secondaryOccStyle.borderTopLeftRadius = 3
                        secondaryOccStyle.borderBottomLeftRadius = 3
                    }
                    if (j == (sch.Occ.SecondaryOccs.length - 1)) {
                        secondaryOccStyle.borderTopRightRadius = 3
                        secondaryOccStyle.borderBottomRightRadius = 3
                    }
                    //secondaryOccStyle.borderRadius = 5

                    let noBoxShadow = false

                    if (sch.IsAvailable) {
                        secondaryOccStyle.backgroundColor
                            = (so.Color) || (sch.Occ && sch.Occ.Color) || MulwiColors.blueDark
                    } else {
                        secondaryOccStyle.backgroundColor = darkgray
                        noBoxShadow = true
                    }

                    //const mul = cw / 24 / 60 
                    //secondaryOccStyle.marginLeft += mul * so.OffsetStart

                    //a[day].push({ style: secondaryOccStyle, noBoxShadow: true })
                    a[day2].push({ mday: true, t:t, style: secondaryOccStyle, noBoxShadow: noBoxShadow })
                }

                // add day marker
                if (+now >= +props.week.start && +now <= +props.week.end) {
                    let dayStyle = dayBlockStyle(hr, 1, off)
                    dayStyle.backgroundColor = MulwiColors.redError
                    dayStyle.zIndex = 20
                    dayStyle.height += 2
                    makeStyleMultiDay(dayStyle, now, nowEnd, 1)
                    a[nowday].push({ mday: true, style: dayStyle, noBoxShadow: true })
                }

                let overlayStyle = dayBlockStyle(hr, 1, off)
                makeStyleMultiDay(overlayStyle, sch.Start, sch.End, 7) // schDayBlockStyle(ds, sch, durs)
                overlayStyle.backgroundColor = "transparent"
                if (!overlayStyle.marginRight)
                    overlayStyle.marginRight = 0
                overlayStyle.marginRight += columnItemWidth
                overlayStyle.marginLeft = 0
                overlayStyle.zIndex = 22
                a[999].push({
                    style: overlayStyle,
                    sch: sch,
                    t: t,
                    dark: true,
                    record: s[i],
                    mday: true,
                    noBoxShadow: true,
                    altLabel: sch.Count + " / " + t.Capacity
                })

                let borderStyle = dayBlockStyle(hr, 1, off)
                borderStyle = schDayBlockStyle(borderStyle, sch, durs)
                borderStyle.backgroundColor = "transparent"
                borderStyle.zIndex = 23
                // borderStyle.borderWidth = 2
                a[day].push({
                    style: borderStyle,
                    sch: sch,
                    t: t,
                    record: s[i],
                    mday: true,
                    noBoxShadow: true,
                    notext: true
                })

                _soff += hm

                hr += 1.05
            }
        }

        if (_soff != off) {
            _soff += shift
        }

        for (let i = 0; i < s.length; i++) {
            let t = s[i].Training

            //let found = false
            for (let j = 0; j < s[i].Schedule.length; j++) {
                let sch = s[i].Schedule[j]
                let start = sch.Start
                if (start < props.week.start) {
                    start = props.week.start
                }
                let end = sch.End
                if (end > props.week.end) {
                    end = props.week.end
                }

                if (!sch.IsAvailable && !props.showUnavailable)
                    continue

                let _durs = extendedDfInHours(start, end)
                let durs = []
                for (let z in _durs) {
                    durs.push({ epoch: z, val: _durs[z] })
                }
                let day = start.getDay() - 1
                if (day < 0) day += 7


                // single day only
                if (durs.length === 1) {
                    //let tmpstart = epochToDate(durs[0].epoch)
                    //let hr = 0
                    //if (start.getDate() === tmpstart.getDate()) {
                    let hr = start.getHours() + (start.getMinutes() / 60)
                    //}
                    let style = schDayBlockStyle(dayBlockStyle(hr, durs[0].val, _soff), sch)
                    style.height -= 2

                    if (!sch.IsAvailable) {
                        style.backgroundColor = darkgray
                    }

                    let hasMultipleDays = sch.Occ.SecondaryOccs && sch.Occ.SecondaryOccs.length > 0
                    if (hasMultipleDays) {

                        for (let j = 0; j < sch.Occ.SecondaryOccs.length; j++) {

                            let so = sch.Occ.SecondaryOccs[j]
                            let start2 = new Date(sch.Start)
                            let end2 = new Date(sch.Start)
                            start2.setMinutes(start2.getMinutes() + so.OffsetStart)
                            end2.setMinutes(end2.getMinutes() + so.OffsetEnd)
                            if (start2 < props.week.start || end2 > props.week.end)
                                continue
                            if (start2 < props.week.start)
                                start2 = props.week.start
                            if (end2 > props.week.end)
                                end2 = props.week.end
                            let _durs = extendedDfInHours(start2, end2)

                            let day2 = start2.getDay() - 1
                            if (day2 < 0) day2 += 7

                            let durs = []
                            for (let z in _durs) {
                                durs.push({ epoch: z, val: _durs[z] })
                            }

                            for (let i = 0; i < durs.length; i++) {
                                let x = epochToDate(durs[i].epoch)
                                let hr = x.getHours() + (x.getMinutes() / 60)
                                let style = schDayBlockStyle(dayBlockStyle(hr, durs[i].val, _soff), sch)
                                style.height -= 2
                                if (!sch.IsAvailable) {
                                    style.backgroundColor = darkgray
                                }
                                style.borderRadius = 3
                                a[day2].push({
                                    style: style,
                                    sch: sch,
                                    t: t,
                                    record: s[i]
                                })
                            }

                        }

                    } else {
                        style.borderRadius = 3
                        a[day].push({
                            style: style,
                            sch: sch,
                            t: t,
                            record: s[i]
                        })
                    }

                } else if (durs.length > 1) {


                    if (!sch.Occ.SecondaryOccs || sch.Occ.SecondaryOccs.length == 0) {
                        let so = { ...sch.Occ }
                        so.OffsetStart = 0
                        so.OffsetEnd = dfInHours(start, end) * 60
                        sch.Occ.SecondaryOccs = [so]
                    }

                    for (let j = 0; j < sch.Occ.SecondaryOccs.length; j++) {


                        let so = sch.Occ.SecondaryOccs[j]
                        let start2 = new Date(sch.Start)
                        let end2 = new Date(sch.Start)
                        start2.setMinutes(start2.getMinutes() + so.OffsetStart)
                        end2.setMinutes(end2.getMinutes() + so.OffsetEnd)
                        if (start2 < props.week.start || end2 > props.week.end)
                            continue
                        if (start2 < props.week.start)
                            start2 = props.week.start
                        if (end2 > props.week.end)
                            end2 = props.week.end
                        let _durs = extendedDfInHours(start2, end2)

                        let day2 = start2.getDay() - 1
                        if (day2 < 0) day2 += 7

                        let durs = []
                        for (let z in _durs) {
                            durs.push({ epoch: z, val: _durs[z] })
                        }

                        /* adds grayed out normal vertical item */
                        for (let i = 0; i < durs.length; i++) {
                            let x = epochToDate(durs[i].epoch)
                            let hr = x.getHours() + (x.getMinutes() / 60)
                            let style = schDayBlockStyle(dayBlockStyle(hr, durs[i].val, _soff), sch)
                            style.borderRadius = 3
                            style.backgroundColor = gray
                            style.border = null
                            style.opacity = 1
                            style.zIndex = 0
                            style.borderRadius = 0
                            style.width += 2 * style.marginLeft
                            style.marginLeft = 0
                            a[day2].push({
                                style: style,
                                sch: sch,
                                // t: t,
                                record: s[i],
                                noBoxShadow: true,
                                backgroundItem: true
                            })
                        }
                    }
                }
            }
        }

        setSoff(_soff)

        // vacations

        let vhr = 23
        if (vacs)
            for (let i = 0; i < vacs.length; i++) {
                let v = vacs[i]
                if (!(v.DateEnd >= props.week.start && v.DateStart <= props.week.end)) {
                    continue
                }

                let start = v.DateStart
                if (start < props.week.start) {
                    start = props.week.start
                }
                let end = v.DateEnd
                if (end > props.week.end) {
                    end = props.week.end
                }


                let _durs = extendedDfInHours(start, end)
                let l = 0
                for (let z in _durs) {
                    l++
                }

                let s = dayBlockStyle(vhr, 1)
                s.backgroundColor = "gray"
                s.color = "black"
                makeStyleMultiDay(s, v.DateStart, v.DateEnd, l)
                vhr -= 1.05
                let day = start.getDay() - 1
                if (day < 0) day += 7
                a[day].push({
                    style: s,
                    vacation: {
                        title: locale2.VACATION[props.lang]
                    }
                })
            }

        // indicator for current hour

        if (+now >= +props.week.start && +now <= +props.week.end) {
            for (let i = 0; i < 7; i++) {
                let offset = off + ((now.getHours() + (now.getMinutes() / 60)) * hm)
                if (i == nowday) {
                    a[i].push(verticalBarObject(true, offset))
                } else {
                    a[i].push(verticalBarObject(false, offset))
                }
            }
        }

        for (let i = 0; i < 7; i++) {
            for (let j = 0; j < 24; j += 1) {
                let start = new Date(props.week.start)
                start.setHours(j)
                start.setDate(start.getDate() + i)
                let end = new Date(start)
                end.setHours(end.getHours() + 1)
                let style = dayBlockStyle(j, 1, _soff)
                style.backgroundColor = 'white'
                style.border = null
                style.opacity = 1
                style.zIndex = 0
                style.borderRadius = 0
                style.width = "100%"
                style.marginLeft = 0
                style.borderBottomStyle = j % 2 != 0 ? 'solid' : 'none'
                style.borderTopStyle = j % 2 != 0 ? 'none' : 'solid'
                style.borderLeftStyle = 'solid'
                style.borderRightStyle = 'solid'
                style.borderWidth = '1px'
                style.borderColor = gray
                a[i].push({
                    style: style,
                    noBoxShadow: true,
                    backgroundItem: true,
                    addTrainingBackground: true,
                    addTrainingData: {
                        start: start,
                        end: end
                    }
                })
            }
        }

        a = setOverlaps(a)

        if (!mountedRef.current) return null
        setItems(a)
    }

    useEffect(() => {
        refreshHarmonogram()
    }, [props.hideMultiday, props.week, props.refreshToken,
    props.trainingID, props.editorData, props.showUnavailable])

    const [forceModalEdit, setForceModalEdit] = useState(false)
    const [modalOpen, setModalOpen] = useState(false)
    const [editIndexes, setEditIndexes] = useState([])

    const [opPickerOpen, setOpPickerOpen] = useState(false)

    const [forceOccEdit, setForceOccEdit] = useState(false)

    function setOp(op) {
        setOpPickerOpen(false)
        switch (op) {
            case 1:
                setForceOccEdit(true)
                setForceModalEdit(false)
                setModalOpen(true)
                break
            case 2:
                let i = editIndexes[0]
                props.editorData.occs.splice(i, 1)
                props.setEditorData({ ...props.editorData })
                break
            case 3:
                setForceOccEdit(false)
                setForceModalEdit(true)
                setModalOpen(true)
                break
        }
    }

    const [lastAdded, setLastAdded] = useState([])

    function handleOnClickColumn(e, dd) {
        // NEG props.mdayEdit 
        //    || props.repeating === -1 
        //    || (!props.editorData.occs || props.editorData.occs.length === 0)
        if (!props.mdayEdit && props.repeating !== -1 &&
            !(!props.editorData || !props.editorData.occs || props.editorData.occs.length === 0))
            return
        if (!props.editorData)
            return
        let br = e.target.getBoundingClientRect()
        let clicked = e.pageY
        let allowedStart = br.top + soff - 10
        if (clicked < allowedStart)
            return
        let bottom = allowedStart + br.height - soff
        let top = allowedStart
        bottom -= top
        clicked -= top
        let norm = clicked / bottom
        if (norm < 0)
            norm = 0
        if (norm > 1)
            norm = 1
        let hr = (norm * 24)
        hr = Math.round(hr)
        dd.setHours(hr)
        // new element
        setForceOccEdit(false)
        if (!props.editorData.occs || props.editorData.occs.length === 0) {

            let [occs, indexes] = addOcc2(props.editorData.occs, null, dd, props.repeating)
            setEditIndexes(indexes)
            props.editorData.occs = occs
            props.setEditorData({ ...props.editorData })
            e.stopPropagation()

        } else if (props.mdayEdit) {

            let eix = props.editIndexes
            if (eix && eix.length > 0) {
                let i = eix[0]
                let occ = props.editorData.occs[i]
                let df = dfInHours(occ.DateStart, dd) * 60
                let so = {
                    OffsetStart: df,
                    OffsetEnd: df + 120,
                    Color: "",
                    Remarks: ""
                }
                let fix = randomString(24)
                so.fix = fix
                occ.SecondaryOccs.push(so)
                adjustOccDatesAfterEdit(occ)
                let ix = occ.SecondaryOccs.findIndex(c => c.fix === fix)
                delete so.fix
                for (let j = 0; j < props.editorData.occs.length; j++) {
                    props.editorData.occs[j].SecondaryOccs = occ.SecondaryOccs
                    adjustOccDatesAfterEdit(props.editorData.occs[j])
                }
                let editIndexes = [i, ix]
                props.setEditorData({ ...props.editorData })
                setEditIndexes(editIndexes)
            }

        } else if (props.editorData.occs && props.editorData.occs.length >= 1) {

            let len = props.editorData.occs.length
            let occ = { ...props.editorData.occs[len - 1] }
            // dd = new Date(asUTC(dd))
            occ.DateStart = dd
            adjustOccDatesAfterEdit(occ)
            props.editorData.occs.push(occ)
            let ed = { ...props.editorData }
            setEditIndexes([props.editorData.occs.length - 1])
            props.setEditorData(ed)
        }

        setForceModalEdit(false)
        setModalOpen(true)
        e.stopPropagation()
    }

    function DD() {
        return (<Divider style={{
            height: 2,
            borderColor: theme.palette.grey[300],
            marginTop: -10 + soff,
            position: "absolute",
            width: "100%"
        }} />)
    }

    function mdayEditModal() {
        // i dont even know any more...
        if (!forceOccEdit && ((props.repeating != -1 && mday && !props.mdayEdit) || (!props.mdayEdit && forceModalEdit))) {
            return (<ScheduleEditModal
                open={modalOpen}
                lang={props.lang}
                editIndexes={editIndexes}
                editorData={props.editorData}
                setEditorData={props.setEditorData}
                setOpen={setModalOpen}
                week={props.week}
            />)
        } else if (!props.readonly && props.editorData) {
            if (props.repeating == -1 || forceOccEdit) {
                return (<MoveOccModal
                    editIndexes={editIndexes}
                    setLastAdded={setLastAdded}
                    occs={props.editorData.occs}
                    setOccs={occs => {
                        let ed = { ...props.editorData }
                        ed.occs = occs
                        setEditIndexes([])
                        setModalOpen(false)
                        props.setEditorData(ed)
                    }}
                    open={modalOpen}
                    lang={props.lang}
                    setOpen={setModalOpen}
                />)
            } else {
                return (<ModifyOcc2Modal
                    editIndexes={editIndexes}
                    setLastAdded={setLastAdded}
                    occs={props.editorData.occs}
                    setOccs={occs => {
                        let ed = { ...props.editorData }
                        ed.occs = occs
                        setEditIndexes([])
                        setModalOpen(false)
                        props.setEditorData(ed)
                    }}
                    open={modalOpen}
                    lang={props.lang}
                    setOpen={setModalOpen}
                />)
            }
        }
    }

    function getColumnForDay(day, index) {
        let dis = items[index]
        let dd = new Date(props.week.start)
        dd.setDate(dd.getDate() + index)
        return (
            <Grid container direction="column"
                onClick={e => {
                    handleOnClickColumn(e, dd)
                }}
                style={{
                    borderLeft: day == -1 ? "none" : "solid",
                    borderWidth: 2,
                    height: 24 * hm + soff - 11,
                    width: "13%",
                    position: "relative",
                    borderColor: theme.palette.grey[300],
                    cursor: (props.mdayEdit || props.repeating === -1 || (props.editorData && (!props.editorData.occs || props.editorData.occs.length === 0))) ? "pointer" : null
                }}>
                {soff != off && day != -1 && <React.Fragment>
                    {DD()}
                </React.Fragment>}
                {day != -1 && (
                    <Typography align="center"><span style={{ color: "gray" }}>
                        {dd.getDate()} / </span>{day}
                    </Typography>
                )}
                <Divider />
                {dis && dis.map((di, i) => {
                    let fs = 12
                    if (fs > di.style.height) {
                        fs = di.style.height - 2
                        if (fs < 0) fs = 0
                    }
                    if (di.style.width <= (columnItemWidth / 2)) {
                        fs = 9
                    }
                    if (di.style.width <= (columnItemWidth / 3)) {
                        fs = 7
                    }
                    let c = di.dark ? "black" : "white"
                    return (<React.Fragment key={i}>
                            <HtmlTooltip title={(di.sch && di.t && (<React.Fragment>
                                <p><strong>{di.t.Title}</strong></p>
                                <p>{locale2.START[props.lang]}: <strong>{prettyPrintDate(di.sch.Start, props.lang)}</strong></p>
                                <p>{locale2.END[props.lang]}: <strong>{prettyPrintDate(di.sch.End, props.lang)}</strong></p>
                                <p>{di.sch.Count} / {di.t.Capacity}</p>
                                <p>{avReasonToStr(di.sch.AvailabilityReason)}</p>
                            </React.Fragment>)) || ""}>
                                <Box
                                    className={di.addTrainingBackground ? "hoverable-bg-item" : cls.cardItem}
                                    boxShadow={di.noBoxShadow ? null : 3} style={di.style} onClick={(e) => {

                                        if (!props.user && di.addTrainingBackground) {
                                            setTrainingOpen(true)
                                            setOccData(di.addTrainingData)
                                            return
                                        }

                                        if (di.record && di.sch && (!props.user || di.sch.IsAvailable)) {
                                            if (props.setDrawerData)
                                                props.setDrawerData(trainingResToDrawerData(di.record, props.user, di.sch))
                                            if (props.setDrawerOpen)
                                                props.setDrawerOpen(true)
                                        }

                                        if (props.editorData && di.editIndexes) {
                                            if (props.mdayEdit) {
                                                setEditIndexes(di.editIndexes)
                                                setForceModalEdit(true)
                                                setModalOpen(true)
                                            } else {
                                                setEditIndexes(di.editIndexes)
                                                setOpPickerOpen(true)
                                            }
                                            e.stopPropagation()
                                        }
                                    }}>
                                    {!di.notext && !di.altLabel && di.sch && di.t && <div style={{
                                        position: "absolute",
                                        right: 1,
                                        bottom: -1,
                                        color: c,
                                        fontSize: 7,
                                    }}><strong>{di.sch.Count + " / " + di.t.Capacity}</strong></div>}

                                    {di.no && (<div style={{
                                        position: "absolute",
                                        left: 1,
                                        top: 1,
                                        fontSize: fs,
                                        color: c
                                    }}><strong>{locale2.ACTIVITY[props.lang]} {di.no}</strong></div>)}

                                    {!props.readonly && props.editorData && !di.noEdit && (<Edit size="small" style={{
                                        position: "absolute",
                                        right: 1,
                                        top: 1,
                                        width: Math.min(25, di.style.height),
                                        height: Math.min(25, di.style.height),
                                        color: c
                                    }} />)}

                                    {di.occ2 && (<div style={{
                                        position: "absolute",
                                        left: 4,
                                        bottom: 0,
                                        fontSize: fs,
                                        color: c
                                    }}><strong>{di.occ2.Remarks}</strong></div>)}

                                    {!di.notext && !di.altLabel && (di.t || di.vacation) && (<div style={{
                                        position: "absolute",
                                        left: 2,
                                        top: 1,
                                        fontSize: fs,
                                        color: c
                                    }}><strong>{(di.t && di.t.Title) || di.vacation.title}</strong></div>)}


                                    {di.altLabel && (di.t || di.vacation) && (<Grid container direction="row" justify="space-between" style={{
                                        position: "absolute",
                                        left: 2,
                                        top: 3,
                                        width: columnItemWidth - 5,
                                        fontSize: fs,
                                        color: c
                                    }}>
                                        <Grid item style={{
                                            maxWidth: columnItemWidth - (columnItemWidth / 3),
                                            textOverflow: "ellipsis",
                                            whiteSpace: "nowrap",
                                            overflow: "hidden"
                                        }}>
                                            <strong>{(di.t && di.t.Title) || di.vacation.title}</strong>
                                        </Grid>
                                        <Grid item>
                                            {di.altLabel}
                                        </Grid>
                                    </Grid>)}

                                </Box>
                            </HtmlTooltip>
                    </React.Fragment>)
                })}
            </Grid>)
    }

    function marginTopForHour(h) {
        return -10 + soff + (hm * h * 0.97)
    }

    const [trainingOpen, setTrainingOpen] = useState(false)
    const [occData, setOccData] = useState(null)

    return (<React.Fragment>
        <AddTraining
            lang={props.lang}
            setDrawerData={props.setDrawerData}
            openDrawer={() => props.setDrawerOpen(true)}
            onChange={() => props.setRefreshToken(!props.refreshToken)}
            modal
            open={trainingOpen}
            setOpen={setTrainingOpen}
            occData={occData}
            onClose={(success) => { }} />
        <Grid container
            direction="row"
            style={{
                borderColor: theme.palette.grey[300],
                borderStyle: "solid",
                borderWidth: 2,
                width: "100%",
                position: "relative"
            }}
            justify="space-between">
            <Grid container direction="column" style={{
                width: "9%",
            }}>
                <Grid item>
                    <div style={{ position: "absolute" }}>{getColumnForDay(-1, 999)}</div>
                    <Grid container direction="column" style={{
                        width: "100%",
                        position: "relative"
                    }}>
                        <Typography style={{
                            backgroundColor: theme.palette.grey[300]
                        }} align="center">&nbsp;</Typography>
                        <Divider />

                        {soff != off && (<React.Fragment>
                            {DD()}
                        </React.Fragment>)}
                        {getHours().map((h) => (
                            <Typography key={h.h}
                                style={{
                                    position: "absolute",
                                    marginLeft: "calc(50% - 25px)",
                                    marginTop: marginTopForHour(h.h)
                                }}
                            >{h.s}</Typography>
                        ))}
                    </Grid>
                </Grid>

            </Grid>
            {getDays().map((d, i) => <React.Fragment key={i}>
                {getColumnForDay(d, i)}
            </React.Fragment>)}
        </Grid>
    </React.Fragment>)
}