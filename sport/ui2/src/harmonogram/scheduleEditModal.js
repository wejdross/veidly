import { Button, Dialog, DialogContent } from '@mui/material'
import { DialogActions } from '@mui/material'
import React, { useEffect, useRef, useState } from 'react'
import { locale2 } from '../locale'
import { MulwiColors } from '../mulwiColors'
import { getWkFromMonth, WeekSwitch } from './harmonogram'
import { prettyPrintDate, prettyPrintDateRange } from './trainingDetails'
import { HarmonogramWeek } from './weekBigRes'

export function ScheduleEditModal(props) {

    const [week, setWeek] = useState(() => getWkFromMonth(new Date()))

    useEffect(() => {
        setWeek(props.week)
    }, [props.week])

    let start, end

    if(props.editIndexes.length >= 1) {
        let occ = props.editorData.occs[props.editIndexes[0]]
        if(!occ)
            return null
        start = occ.DateStart
        end = occ.DateEnd
    }

    return (<Dialog fullWidth maxWidth="md" open={props.open} onClose={() => props.setOpen(false)}>
        <DialogContent>
            <center>
                <h1>{locale2.TRAINING_SCHEDULE[props.lang]}</h1>
                <p>{start && end && prettyPrintDateRange(start, end, null, false, props.lang)} </p>
                
                {/* <Button
                    size='small'
                    variant="contained" style={{
                        color: "white",
                        marginLeft: 5,
                        marginBottom: 2,
                        backgroundColor:MulwiColors.redError
                }} onClick={() => {
                    if(props.editIndexes.length < 1)
                        return
                    let i = props.editIndexes[0]
                    props.editorData.occs.splice(i, 1)
                    props.setEditorData({...props.editorData})
                    props.setOpen(false)
                }}>Usuń te występowanie</Button></p> */}
                
            </center>
            <WeekSwitch lang={props.lang} 
                week={week} setWeek={setWeek} setMonthDate={null} />
            <HarmonogramWeek 
                lang={props.lang}
                setInfo={props.setInfo}
                mdayEdit
                editIndexes={props.editIndexes}
                editorData={props.editorData}
                setEditorData={props.setEditorData}
                week={week} 
            />
        </DialogContent>
        <DialogActions>
            <Button onClick={() => props.setOpen(false)}>
                {locale2.CLOSE[props.lang]}
            </Button>
        </DialogActions>
    </Dialog>)
}