import { Button, Card, Grid, Typography, useTheme } from '@mui/material'
import React, { useEffect, useRef, useState } from 'react'
import { getRsvSchedule, getSchedule, getUserSchedule } from '../apicalls/instructor.api'
import { dfInHours, epochToDate, extendedDfInHours, 
        rmtoken, trainingResToDrawerData } from '../helpers'
import { MulwiColors } from '../mulwiColors'
import { prettyPrintDateRange } from './trainingDetails'
import { hrDecToLabel } from './weekBigRes'
import { locale2 } from '../locale'
import { Add, Edit } from '@mui/icons-material'
import { addOcc2, ModifyOcc2Modal } from './calendarEditor'

export function HarmonogramDay(props) {

    const [data, setData] = useState([])

    // this will be set to false once unmounted
    const mountedRef = useRef(true)
    useEffect(() => {
        return () => { 
          mountedRef.current = false
        }
      }, [])

    async function setRemoteSchedule() {
        let start = new Date(props.day)
        start.setHours(0, 0, 0, 0)
        let end = new Date(props.day)
        end.setHours(23, 59, 59, 0)
        //let d = await getSchedule(start, end)
        //d = JSON.parse(d)

        let d = []

        try { 
            if (!mountedRef.current) return null
            if(props.usrRsv) {
                d = await getRsvSchedule(
                    start, 
                    end, 
                    props.trainingID)
            } else {
                if(props.user) {
                    if(!props.instructorID) return null
                    d = await getUserSchedule(
                        start, 
                        end,
                        props.instructorID,
                        props.trainingID,
                        props.smID)
                } else {
                    d = await getSchedule(
                        start, 
                        end, 
                        props.trainingID)
                }
            }
        } catch(ex) {

            props.setInfo && props.setInfo({
                open: true,
                hdr: locale2.COULDNT_GET_DATA[props.lang],
                msg: locale2.ERROR[props.lang] + ': ' + ex,
                buttons: (
                  <React.Fragment>
                    <Button onClick={() => {
                            rmtoken()
                            window.location = "/"
                        }} color="primary">
                            {locale2.RESET_APP[props.lang]}
                    </Button>
                  </React.Fragment>
                )
              })
            return
        }

        let _data = []

        for(let i = 0; i < d.length; i++) {
            let s = d[i].Schedule
            let title = d[i].Training.Title
            for(let j = 0; j < s.length; j++) {
                let df = extendedDfInHours(s[j].Start, s[j].End)
                for(let z in df) {
                    let x = epochToDate(z)
                    if(x.getFullYear() === props.day.getFullYear() &&
                        x.getMonth() === props.day.getMonth() &&
                            x.getDate() === props.day.getDate()) {
                        let hr = x.getHours() + (x.getMinutes() / 60)
                        _data.push({
                            title: title,
                            decLabel: hr,
                            hrStart: hrDecToLabel(hr),
                            hrEnd: hrDecToLabel(hr + df[z]),
                            duration: df[z],
                            trainingStart: s[j].Start,
                            trainingEnd: s[j].End,
                            t: d[i],
                            session: s[j],
                            isMultiDay: Object.keys(df).length === 1 ? false : true,
                            color: (s[j].Occ && s[j].Occ.Color) || MulwiColors.blueDark
                        })
                    }
                }
            }
        }

        _data.sort((a,b) => a.decLabel - b.decLabel)

        if (!mountedRef.current) return null
        setData(_data)
    }

    const [modalOpen, setModalOpen] = useState(false)
    const [lastAdded, setLastAdded] = useState([])
    const [editIndexes, setEditIndexes] = useState([])


    function setEditorSchedule() {
        let _data = []
        let ed = props.editorData
        if(!ed) {
            setData([])
            return
        }
        let occs = ed.occs
        if(!occs || occs.length === 0) {
            setData([])
            return
        }
        for(let i = 0; i < occs.length; i++) {
            let occ = occs[i]
            if(!occ.SecondaryOccs || occ.SecondaryOccs.length == 0) {
                let df = dfInHours(occ.DateStart, occ.DateEnd)
                let so = {...occ}
                so.OffsetStart = 0
                so.OffsetEnd = df * 60
                occ.SecondaryOccs = [so]
                let ed = {...props.editorData}
                props.setEditorData(ed)
                // setEditorData will trigger refresh so no need to continue
                return
            }
            if(occ.SecondaryOccs && occ.SecondaryOccs.length > 0) {
                for(let j = 0; j < occ.SecondaryOccs.length; j++) {
                    let occ2 = occ.SecondaryOccs[j]
                    let start = new Date(occ.DateStart)
                    let end = new Date(occ.DateStart)
                    start.setMinutes(start.getMinutes() + occ2.OffsetStart)
                    end.setMinutes(end.getMinutes() + occ2.OffsetEnd)
                    let df = extendedDfInHours(start, end)
                    for(let z in df) {
                        let x = epochToDate(z)
                        if(x.getFullYear() === props.day.getFullYear() &&
                            x.getMonth() === props.day.getMonth() &&
                                x.getDate() === props.day.getDate()) {
                            let hr = x.getHours() + (x.getMinutes() / 60)
                            _data.push({
                               title: ed.training.Title,
                               decLabel: hr,
                               hrStart: hrDecToLabel(hr),
                               hrEnd: hrDecToLabel(hr + df[z]),
                               duration: df[z],
                               trainingStart: occ.DateStart,
                               trainingEnd: occ.DateEnd,
                               t: ed.training,
                               isMultiDay: Object.keys(df).length === 1 ? false : true,
                               color: occ2.Color || occ.Color || MulwiColors.blueDark,
                               editIndexes: [i, j]
                            })
                        }
                    }
                }
            }
        }
        if (!mountedRef.current) return null
        setData(_data)
    }

    useEffect(() => {

        if(props.editorData) {
            setEditorSchedule()
        } else {
            setRemoteSchedule()
        }
        
        // eslint-disable-next-line
    }, [props.day, props.editorData, props.refreshToken, props.trainingID])

    return (<React.Fragment>

        {props.editorData && <ModifyOcc2Modal
                editIndexes={editIndexes} 
                setLastAdded={setLastAdded}
                occs={props.editorData.occs}
                setOccs={occs => {
                    let ed = {...props.editorData}
                    ed.occs = occs
                    setEditIndexes([])
                    setModalOpen(false)
                    props.setEditorData(ed)
                }}
                open={modalOpen} 
                lang={props.lang}
                setOpen={setModalOpen} />}
        {props.editorData && (<Button onClick={() => {
            let [occs, indexes] = addOcc2(props.editorData.occs, lastAdded, props.day, props.repeating)
            let ed = {...props.editorData}
            ed.occs = occs
            props.setEditorData(ed)
            setEditIndexes(indexes)
            setModalOpen(true)
        }} style={{
            backgroundColor: MulwiColors.blueDark,
            color: "white",
            marginBottom: 10,
        }} fullWidth><Add/></Button>)}
        {data.map((d,i) => (<React.Fragment key={i} >
            <Card
                onClick={() => {
                    if(d.t && d.session && d.t.Training && d.session) {
                        props.setDrawerOpen(true)
                        props.setDrawerData(trainingResToDrawerData(d.t, props.user, d.session))
                    }
                }}
                style={{
                    padding: 5,
                    backgroundColor: d.color,
                    marginBottom: 5,
                    cursor:"pointer",
                    color: "white",
                    position:"relative",
                    borderRadius: 5,
                }}>
                    {props.editorData && d.editIndexes && (
                        <Button onClick={() => {
                            setEditIndexes(d.editIndexes)
                            setModalOpen(true)
                        }} style={{
                            position:"absolute",
                            right: 0,
                            top: 0,
                            color:"white"
                        }}><Edit/></Button>
                    )}
                    {d.isMultiDay ? (
                        <Typography variant="body2" style={{fontSize: 12}}>{prettyPrintDateRange(d.trainingStart, d.trainingEnd, null, null, props.lang)}</Typography>
                    ) : (
                        <Typography>{d.hrStart} - {d.hrEnd}</Typography>
                    )}
                    <Grid container justify="space-between" direction="row">
                        <Grid item>
                            <Typography><strong>{d.title}</strong></Typography>
                        </Grid>
                        <Grid item>
                            {d.t && d.session && d.t.Training && d.session && (
                                <Typography>{d.session.Count} / {d.t.Training.Capacity}</Typography>
                            )}
                        </Grid>
                    </Grid>
                    {d.t2 && (<React.Fragment>
                        <Typography variant="body2">
                            {d.t2}
                        </Typography>
                        <Typography variant="body2">
                            {d.t3}
                        </Typography>
                    </React.Fragment>)}
            </Card> 
        </React.Fragment>))}
    </React.Fragment>)
}