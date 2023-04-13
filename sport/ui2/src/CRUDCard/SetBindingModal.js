import { Button, Dialog, DialogContent, DialogTitle, 
        IconButton, Table, TableBody, 
        TableCell, TableHead, TableRow, Typography } from '@mui/material';
import { Remove, Repeat } from '@mui/icons-material';
import React, { useEffect, useState } from 'react';
import TrainingAtc from '../harmonogram/trainingAtc';
import { MulwiColors } from '../mulwiColors';
import { errToStr, getErrorDialog } from '../StatusDialog';
import { locale2 } from '../locale';
import { DialogActions } from '@mui/material';
import { sprintf } from '../helpers';

export function SetBindingModal(props) {

    const [open, _setOpen] = useState(false)
    const [trainings, setTrainings] = useState([])
    const [err, setErr] = useState("")

    let lang = props.lang

    function setOpen(o) {
        if(o) setErr("")
        _setOpen(o)
    } 

    async function refresh() {
        try {
            let t = await props.getTrainings(props.record)
            if (t) {
                setTrainings(t)
            }
        } catch (ex) {
            props.setInfo(getErrorDialog("Couldnt download trainings", ex))
        }
    }

    useEffect(() => {
        if(!props.newKey || !props.record) return
        if(props.record.ID === props.newKey) {
            props.setNewKey(null)
            setOpen(true)
        }
    }, [props.newKey, props.record])

    useEffect(() => {
        if (!open || !props.record) return
        refresh()
    }, [open])

    useEffect(() => {
        if(props.open && props.record)
            setOpen(true)
    }, [props.open, props.record])

    if (!props.record) return null

    return (<React.Fragment>
        <Dialog open={open} onClose={() => setOpen(false)}>
            <DialogTitle>{sprintf(locale2.ASSIGNED_TRAININGS_FMT[lang], props.nameSelector(props.record))}</DialogTitle>
            <DialogContent>
                <Typography style={{marginBottom: 5}}>{locale2.ASSIGN_TRAINING[props.lang]}</Typography>
                <TrainingAtc setValue={async v => {
                    if(!v || !v.Training) return
                    try {
                        await props.createBinding(props.record, v.Training.ID)
                        await refresh()
                    } catch(ex) {
                        setErr(errToStr(ex))
                    }
                }} />
                <Typography 
                    variant="body2"
                    style={{color: MulwiColors.redError}}>{err}&nbsp;</Typography>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell align="center">
                                {locale2.TRAINING[props.lang]}
                            </TableCell>
                            <TableCell align="center">
                                {locale2.DELETE[props.lang]}
                            </TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {trainings && trainings.map((t,i) => (
                            <TableRow key={i}>
                                <TableCell align="center">
                                    {t.Training.Title}
                                </TableCell>
                                <TableCell align="center">
                                    <IconButton onClick={async () => {
                                        if (!props.record) return
                                        try {
                                            await props.deleteBinding(props.record, t.Training.ID)
                                            await refresh()
                                        } catch(ex) {
                                            setErr(errToStr(ex))
                                        }
                                    }} style={{
                                        color: MulwiColors.redError
                                    }}>
                                        <Remove/>
                                    </IconButton>
                                </TableCell>
                            </TableRow>))}
                    </TableBody>
                </Table>
            </DialogContent>
            <DialogActions>
                <Button onClick={() => {
                    setOpen(false)
                }}>
                    {locale2.DONE[props.lang]}
                </Button>
            </DialogActions>
        </Dialog>
        <Button style={{
            color: MulwiColors.blueDark
        }} onClick={() => setOpen(true)}>
            <Repeat/>
        </Button>
    </React.Fragment>)
}