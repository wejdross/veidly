import { Button, Grid, useMediaQuery, useTheme } from '@mui/material'
import React, { useEffect, useRef, useState } from 'react'
import { getRsvSchedule, getSchedule, getUserSchedule } from '../apicalls/instructor.api'
import { dayIndex, epochToDate, extendedDfInHours, rmtoken } from '../helpers'
import { MulwiColors } from '../mulwiColors'
import { locale2 } from '../locale'
import { makeStyles } from "@mui/styles";

export const sDays = [0,1,2,3,4,5,6]

export function HarmonogramMonth(props) {
    
    const useStyles = makeStyles(theme => ({
        hover: {
            cursor: "pointer",
            [theme.breakpoints.up("sm")]: {
                '&:hover': {
                    backgroundColor: MulwiColors.lightGreyAddedByLukasz,
                },
            },
        },
      }))

    let d = props.date
    let start = new Date(d.getFullYear(), d.getMonth(), 1, 0, 0, 0, 0)
    let end = new Date(d.getFullYear(), d.getMonth() + 1)
    end.setHours(0, 0, 0, -1)

    function displayDay(n) {
        if(n < 10)
            return "0" + String(n)
        return String(n)
    }
    
    const [schedule, setSchedule] = useState(null)

    const t = useTheme()

    const classes = useStyles(t)

    function setDaily(daily, day, val, hasCustomers) {
        let o = daily[day]
        if(o) {
            if(hasCustomers) o.confirmedCount+=val
            else o.count+=val
        } else {
            if(hasCustomers) daily[day] = { confirmedCount: val, count: 0 }
            else daily[day] = { count: val, confirmedCount: 0 }
        }
    }

    // function setNewDaily(daily, day, val) {
    //     if(daily[day]) {
    //         daily[day].newCount+=val
    //     } else {
    //         daily[day] = {
    //             count: 0,
    //             newCount: val
    //         }
    //     }
    // }

    // this will be set to false once unmounted
    const mountedRef = useRef(true)
    useEffect(() => {
        return () => { 
          mountedRef.current = false
        }
      }, [])

      function adornment(labelContent, style) {
        let s = {
            position:"absolute", 
            top: -3,
            right: -3,
            zIndex: 999,
            fontSize: 10,
            backgroundColor:MulwiColors.blueDark,
            borderRadius: "50%",
            minWidth: 18,
            textAlign: "center",
            paddingTop: 2,
            height: 18,
            color: t.palette.primary.contrastText,
        }
        if(style) for(let k in style) {
            s[k] = style[k]
        }
        return (
            <div style={s}>
                <strong>{labelContent}</strong>
            </div>
        )
      }

    function currentSessionsFragment(s) {
        if(!s.count) return
        return adornment(Math.round(s.count), null)
    }
    
    function confirmedSessionsFragment(s) {
        if(!s.confirmedCount) return
        return adornment(Math.round(s.confirmedCount), {
            top: 11,
            right: -10,
            backgroundColor: MulwiColors.greenDark
        })
    }

    function generateAdorment(n) {
        if(!schedule) return null
        let s = schedule[n]
        if(s) return (
            <React.Fragment>
                {currentSessionsFragment(s)}
                {confirmedSessionsFragment(s)}
            </React.Fragment>
        )
        return null
    }

    function switchToWeek(n) {
        let weekStart = new Date(
            start.getFullYear(), 
            start.getMonth(), 
            n,
            0, 0, 0, 0)
        
        // convert to 0-based
        let x = weekStart.getDay() - 1
        /* 
            a + p = a  (mod p) 
            im just making sure its >= 0 after sub
        */
        if(x < 0) x += 7
        // no need to mod here, it will always be < 7
        weekStart.setDate(weekStart.getDate() - x)
        let weekEnd = new Date(
            weekStart.getFullYear(), 
            weekStart.getMonth(),
            weekStart.getDate() + 6,
            23, 59, 59, 0)
        if(props.switchToWeek)
            props.switchToWeek(weekStart, weekEnd)
    }

    function generateHead() {
        let ret = []
        for(let i = 0; i < 7; i++) {
            ret.push(i)
        }
        return ret
    }


    function getDays() {
        let ret = []
        let i = 0
        let startDay = start.getDay() - 1 // because js returns 1-based
        if(startDay < 0) startDay += 7
        let df = Math.abs(startDay - i)
        for(let i = 0; i < df; i++) ret.push(-i-1)
        let startCpy = new Date(start)
        while(startCpy <= end) {
            ret.push(startCpy.getDate())
            startCpy.setDate(startCpy.getDate() + 1)
        }
        let endDay = end.getDay() - 1
        if(endDay < 0) endDay += 7
        df = 0
        if(endDay < 6) {
            df = 6-endDay
        }
        for(let i = 0; i < df; i++) ret.push(-1)

        // group ret in 7s
        let gret = []
        let ix = 0
        let page = 0
        for(let i = 0; i < ret.length; i++) {
            if(ix === 0) {
                gret.push([])
            }
            gret[page].push(ret[i])
            ix++
            if(ix === 7) {
                ix = 0
                page++
            }
        }

        // if(gret.length == 4) {
        //     gret.push([-1, -1, -1, -1, -1, -1, -1])
        // }

        // if(gret.length == 5) {
        //     gret.splice(0, 0, [-1, -1, -1, -1, -1, -1, -1])
        // }

        return gret
    }

    const [head, setHead] = useState(generateHead())
    const [days, setDays] = useState(getDays())

    async function setRemoteSchedule() {
        let daily = {}
        let s = [] //await getSchedule(start, end)

        try {
            if(props.usrRsv) {
                s = await getRsvSchedule(start, end, props.trainingID)
            } else {
                if(props.user) {
                    if(!props.instructorID) return null
                    s = await getUserSchedule(
                        start, 
                        end,
                        props.instructorID,
                        props.trainingID,
                        props.smID)
                } else {
                    s = await getSchedule(
                        start, 
                        end,
                        props.trainingID)
                }
            }
        } catch(ex) {
            props.setInfo && props.setInfo({
                open: true,
                hdr: locale2.ERROR[props.lang],
                msg: ex,
                buttons: (
                  <React.Fragment>
                    <Button onClick={() => {
                            rmtoken()
                            window.location = "/"
                        }} color="primary">
                            Zresetuj aplikacje
                    </Button>
                  </React.Fragment>
                )
              })
            return
        }
        // create daily index
        
        // for each training
        for(let i = 0; i < s.length; i++) {
            // for each training occurrence
            for(let j = 0; j < s[i].Schedule.length; j++) {
                let sc = s[i].Schedule[j]
                if (!sc.IsAvailable && !props.showUnavailable)
                    continue
                if(sc.Occ && sc.Occ.SecondaryOccs && sc.Occ.SecondaryOccs.length > 0) {
                    let soccs = sc.Occ.SecondaryOccs
                    for(let z = 0; z < soccs.length; z++) {
                        let s = new Date(sc.Start)
                        let e = new Date(sc.Start)
                        s.setMinutes(s.getMinutes() + soccs[z].OffsetStart)
                        e.setMinutes(e.getMinutes() + soccs[z].OffsetEnd)
                        let d = extendedDfInHours(s, e)
                        for(let z in d) {
                            let zd = epochToDate(z)
                            if(zd >= start && zd <= end)
                                setDaily(daily, zd.getDate(), d[z], sc.Count > 0)
                        }
                    }
                } else {
                    let d = extendedDfInHours(sc.Start, sc.End)
                    for(let z in d) {
                        let zd = epochToDate(z)
                        if(zd >= start && zd <= end)
                            setDaily(daily, zd.getDate(), d[z], sc.Count > 0)
                    }
                }
            }
        }

        if (!mountedRef.current) return null

        setSchedule(daily)
    } 

    function setEditSchedule() {
        let daily = {}
        let ed = props.editorData
        if(!ed) {
        setSchedule(daily)
        return
        }
        let occs = ed.occs
        if(!occs || occs.length === 0) {
            setSchedule(daily)
            return
            }
        for(let i = 0; i < occs.length; i++) {
            let occ = occs[i]
            if(occ.SecondaryOccs && occ.SecondaryOccs.length > 0) {
                for(let j = 0; j < occ.SecondaryOccs.length; j++) {
                    let occ2 = occ.SecondaryOccs[j]
                    let start = new Date(occ.DateStart)
                    let end = new Date(occ.DateStart)
                    start.setMinutes(start.getMinutes() + occ2.OffsetStart)
                    end.setMinutes(end.getMinutes() + occ2.OffsetEnd)
                    let d = extendedDfInHours(start, end)
                    for(let z in d) {
                        let zd = epochToDate(z)
                        if(zd >= start && zd <= end)
                            setDaily(daily, zd.getDate(), d[z])
                    }
                }
            } else {
                let d = extendedDfInHours(occ.DateStart, occ.DateEnd)
                for(let z in d) {
                    let zd = epochToDate(z)
                    if(zd >= start && zd <= end)
                        setDaily(daily, zd.getDate(), d[z])
                }
            }
            break
        }
        setSchedule(daily)
    }

    useEffect(() => {
        if(props.editorData) {
            setEditSchedule()
        } else {
            setRemoteSchedule()
        }
        setHead(generateHead())
        setDays(getDays())
    // eslint-disable-next-line
    }, [props.editorData, props.date, props.refreshToken, props.trainingID, props.showUnavailable])


    function getBgStyle(d) {
        if(!d) return null
        let matches = false
        if(props.week && props.week.start && props.week.end) {
            for(let i = 0; i < d.length; i++) {
                if((d[i] === props.week.start.getDate() && props.date.getMonth() === props.week.start.getMonth()) || 
                (d[i] === props.week.end.getDate() && props.date.getMonth() === props.week.end.getMonth())) {
                    matches = true
                    break
                }
            }
        }
        let s = {
            borderRadius: 8,
            cursor:"pointer",
            marginBottom: 4,
            padding: 4,
            marginLeft: -5
        }
        if(matches && !isLowRes) {
            // s.borderStyle = "solid"
            // s.borderColor = "gray"
            // s.borderWidth = 2

            s.backgroundColor = MulwiColors.pinkAction
        }
        return s
    }

    const today = new Date()

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

    function getStyleForDay(d) {
        let s = {
            marginLeft: 7, 
            marginRight: isLowRes ? 0 : 15, 
            padding: 10,
            position:"relative",
            borderRadius: "50%",
        }
        if(d < 0) {
            s.color = "transparent"
        }
        if(d === today.getDate() && 
                today.getMonth() === props.date.getMonth() && 
                    today.getFullYear() === props.date.getFullYear()) {
                    //s.backgroundColor = t.palette.primary.main
                    //s.color = "white"
                    // s.borderStyle = "solid"
                    // s.borderWidth = 2
                    // s.borderColor = t.palette.primary.main
                    s.backgroundColor = MulwiColors.pinkDark
                    s.color = "white"
        }
        if(props.day && isLowRes) {
            if(d === props.day.getDate() && 
                props.day.getMonth() === props.date.getMonth() &&
                    props.day.getFullYear() === props.date.getFullYear()) {
                s.backgroundColor = MulwiColors.greenDark
                s.color = "white"
            } 
        }
        return (s)
    }

    return (
        <React.Fragment>
            <div style={{minWidth: isLowRes ? 340 : 440, marginLeft: isLowRes ? 20 : 10}}>
                    <div style={{display:"inline-block", marginLeft: 10, marginBottom: 10}}>
                        {head.map(d => (
                            <span key={d} style={{marginLeft: isLowRes ? 3 : 7, marginRight: isLowRes ? 25 : 35}} >
                                    {/* head is 0-based days */}
                                    {dayIndex[d + 1][props.lang]}
                                </span>
                        ))}
                    </div>
                    {/* {shouldGeneratePrefix() && arrayRepeat(7, -1).map((d, i) => (
                        <React.Fragment key={i}>
                            {generateCalendarPaddingItem()}
                        </React.Fragment>
                    ))} */}
                    {days.map((gd, i) => (
                        <React.Fragment key={i}>
                            
                            <Grid container direction="row"  
                                        onClick={() => switchToWeek(i*7+1)}
                                        className={classes.hover}
                                        style={getBgStyle(gd)}>
                                    {gd.map((d, j) => (
                                        <Grid item key={j} 
                                            onClick={() => {
                                                if(props.setDay && isLowRes) {
                                                    props.setDay(new Date(
                                                        props.date.getFullYear(),
                                                        props.date.getMonth(),
                                                        d,
                                                        0,0,0,0
                                                    ))
                                                }
                                            }}>
                                            <div style={getStyleForDay(d)}>
                                                {(d < 0 && "00") || displayDay(d)} 
                                                {d > 0 && generateAdorment(d)}
                                            </div>
                                        </Grid>
                                    ))}
                                </Grid>

                        </React.Fragment>
                    ))}
                    {/* {shouldGenerateSuffix() && arrayRepeat(7, -1).map((d, i) => (
                        <React.Fragment key={i}>
                            {generateCalendarPaddingItem()}
                        </React.Fragment>
                    ))} */}
                    {/* {!shouldGeneratePrefix() && !shouldGenerateSuffix() && days.length < 6 && arrayRepeat(7, -1).map((d, i) => (
                        <React.Fragment key={i}>
                            {generateCalendarPaddingItem()}
                        </React.Fragment>
                    ))} */}
            </div>
        </React.Fragment>
    )
}
