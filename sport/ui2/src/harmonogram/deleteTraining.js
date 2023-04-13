import { Button, CircularProgress, Dialog, DialogActions, 
        DialogContent, DialogTitle, Typography } from '@mui/material'
import { Check } from '@mui/icons-material'
import React, { useState } from 'react'
import { apiDeleteTraining } from '../apicalls/instructor.api'
import { MulwiColors } from '../mulwiColors'
import { locale2 } from '../locale'

export function DeleteTraining(props) {

    async function deleteTraining() {
        setInfo({
            title: locale2.REMOVING_TRAINING[props.lang] +"...",
            msg: <CircularProgress/>
        })
        try {
            await apiDeleteTraining(props.id)
            setInfo({
                title: locale2.TRAINING_HAS_BEEN_REMOVED[props.lang],
                msg: <Check style={{color: "green"}} />,
            })
            props.onChange && props.onChange()
            props.setDrawerOpen(false)
        } catch(ex) {
            setInfo({
                title: locale2.ERROR[props.lang],
                msg: <React.Fragment>
                    <Typography>{ex}</Typography>
                </React.Fragment>,
            })
        }
    }

    const [open, setOpen] = useState(false)

    const [info, setInfo] = useState({
        title: "",
        msg: null
    })

    if(!props.id) {
        return null;
    }

    return (
        <React.Fragment>
            <Dialog  open={open} onClose={() => setOpen(false)}>
                <DialogTitle>
                    {info.title}
                </DialogTitle>
                <DialogContent>
                    <div style={{display:"table", margin: "0 auto"}}>
                        {info.msg}
                    </div>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpen(false)} color="secondary">
                        {locale2.CLOSE[props.lang]}
                    </Button>
                </DialogActions>
            </Dialog>
            <Button onClick={() => {
                setOpen(true)
                setInfo({
                    title: locale2.CONFIRM[props.lang],
                    msg: (
                        <React.Fragment>
                            <Typography>{locale2.CANT_UNDO[props.lang]}</Typography>
                            <Typography>
                                {locale2.EXISTING_RSV_WARN[props.lang]}
                            </Typography>
                            <Button variant="contained" style={{
                                color: "white",
                                backgroundColor: MulwiColors.redError
                            }}
                                onClick={deleteTraining}
                            >{locale2.DELETE_TRAINING[props.lang]}</Button>
                        </React.Fragment>
                    )
                })
            }} variant="contained" style={{
                backgroundColor:MulwiColors.redError,
                color:MulwiColors.whiteBackground
            }}>
                {locale2.DELETE_TRAINING[props.lang]}
            </Button>
        </React.Fragment>
    )
}