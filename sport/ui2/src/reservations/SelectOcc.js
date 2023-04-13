import { Button, Dialog, DialogActions, DialogContent, 
        DialogTitle, List, ListItem, 
        ListItemText, Typography } from "@mui/material";
import { Pagination } from "@mui/lab";
import React, { useState } from "react";
import { useHistory } from "react-router";
import { prettyPrintDate, prettyPrintDay } from "../harmonogram/trainingDetails";
import { avReasonToStr, dateToEpoch } from "../helpers";
import { MulwiColors } from "../mulwiColors";
import { SubCard } from "../sub/SubCard";
import { locale2 } from "../locale";

export default function SelectOcc(props) {

    const h = useHistory()
    const [page, setPage] = useState(0)
    const pageSize = 4

    function navigateRsv(d, smID) {
        if (smID) {
            let l = ("/sub_purch?instructorID=" +
                props.elem.Instructor.id +
                "&trainingID=" +
                props.elem.Training.ID +
                "&smID=" + smID)
            h.push(l)
            return
        }

        let l = ("/rsv?instructorID=" +
            props.elem.Instructor.id +
            "&trainingID=" +
            props.elem.Training.ID +
            "&dateStart=" + dateToEpoch(d))
        h.push(l)
    }


    if(!props.elem) {
        return null
    }
    
    let training = props.elem.Training
    let schedule = props.elem.Schedule
    let sms = props.elem.Sms

    if (!training || !schedule || !sms) return null

    function gotoschedule() {

        let d = new Date(props.dr.DateStart)
        for (let i = 0; i < schedule.length; i++) {
            {
                if (schedule[i].IsAvailable) {
                    d = new Date(schedule[i].Start)
                    break
                }
            }
        }

        // let l = ("/instr/sched?instructorID=" +
        //     training.InstructorID +
        //     "&trainingID=" +
        //     training.ID +
        //     "&dateStart=" + dateToEpoch(d))

        let l = ("/instr_profile?instructorID=" + training.InstructorID)

        h.push(l)
    }

    function pi(c) {

        if (!c.IsAvailable) {
            return <ListItem button disabled>
                <ListItemText primary={prettyPrintDate(new Date(c.Start), props.lang)} 
                    secondary={avReasonToStr(c.AvailabilityReason)} />
            </ListItem>
        }
        
        let lbl = c.Count + "/" + training.Capacity

        return <ListItem button onClick={() => {
            navigateRsv(new Date(c.Start))
        }}>
            <ListItemText 
                primary={prettyPrintDate(new Date(c.Start), props.lang)} 
                secondary={lbl} />
        </ListItem>
    }

    return (<Dialog open={props.open} onClose={() => props.setOpen(false)}>
            <DialogTitle>
                <Typography>{training.Title}</Typography>
                {locale2.WHEN_DO_YOU_WANT_TO_SIGN_UP[props.lang]}
                <Typography variant="body2">
                    {locale2.TERMS_FOR[props.lang]} 
                    {prettyPrintDay(new Date(props.dr.DateStart))}
                    {" - " + prettyPrintDay(new Date(props.dr.DateEnd))}
                </Typography>

            </DialogTitle>
            <DialogContent>
                <List>
                    {schedule && schedule.slice(page * pageSize, (page + 1) * pageSize).map((c, i) => (
                    <React.Fragment key={i}>
                        {pi(c)}
                    </React.Fragment>))}
                </List>
                <Pagination 
                    count={Math.ceil(schedule.length / pageSize)} 
                    page={page+1}
                    onChange={(e,v) => setPage(v-1)} />
                {sms && sms.length > 0 && (<React.Fragment>
                    <DialogTitle>{locale2.INSTR_ALSO_OFFERS_CARNETS[props.lang]}</DialogTitle>
                    {!props.user && (<center><Typography style={{
                        color: MulwiColors.redError
                    }}>{locale2.YOU_MUST_BE_LOGGED_IN_TO_BUY_CARNET[props.lang]}
                </Typography></center>)}
                    <List>
                        {sms.map((s,i) => <ListItem key={i}>
                            <SubCard sm={s} user={props.user} />
                        </ListItem>)}
                    </List>
                </React.Fragment>)}
            </DialogContent>
            <DialogActions>
                <Button
                    onClick={gotoschedule}
                    variant="contained" style={{
                        color: "white",
                        backgroundColor: MulwiColors.blueDark
                    }} fullWidth>{locale2.INSTRUCTORS_PROFILE[props.lang]}</Button>
            </DialogActions>
        </Dialog>
    )
}