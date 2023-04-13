import {
    Button,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Grid,
    TextField
} from '@mui/material'
import React, { useEffect, useRef, useState } from 'react'
import EditIcon from '@mui/icons-material/Edit';
import { patchUserData } from '../apicalls/user.api';
import { Check, Close } from '@mui/icons-material';
import { errToStr } from '../StatusDialog';
import { MulwiColors } from '../mulwiColors';
import { locale2 } from '../locale';
import { makeStyles } from "@mui/styles";

export default function NameEdit(props) {

    const useStyles = makeStyles({
        padded: {
            marginLeft: 100, marginRight: 100, marginTop: 50, marginBottom: 50
        }
    })
    const classes = useStyles()

    useEffect(() => {
        if(props.open) 
            openNameEdit()
    }, [props.open])

    const [fnOpen, setFnOpen] = useState(false)
    function openNameEdit() {
        setNameEdit(c => ({
            firstName: (props.user && props.user.FirstName) || "",
            lastName: (props.user && props.user.LastName) || ""
        }))
        setNameEditSt({ st: "", msg: "" })
        setFnOpen(true)
    }

    const runOnce = useRef(false)

    useEffect(() => {
        if(!props.user)
            return
        if(runOnce.current)
            return
        runOnce.current = true
        let q = new URLSearchParams(window.location.search)
        let flag = q.get("fromlogin")
        if(flag && !props.user.Name)
            setFnOpen(true)
    }, [props.user])

    useEffect(() => {
        if(props.setOpen) {
            if(!fnOpen) props.setOpen(false)
        }
    }, [fnOpen])

    const [nameEdit, _setNameEdit] = useState({ firstName: "", lastName: "" })
    const [nameEditSt, setNameEditSt] = useState({
        st: "",
        msg: ""
    })

    const [err, setErr] = useState(null)

    function setNameEdit(v) {
        v = v(nameEdit)
        let we = false
        if ((v.firstName.length + v.lastName.length) >= 128) {
            setErr(locale2.MAX_ALLOWED_CHARS[props.lang])
            we = true
        }
        if (!we && err) {
            setErr(null)
        }
        _setNameEdit(v)
    }

    function nameEditForm() {
        return nameEditSt.st === "" && (
            <React.Fragment>
                <DialogContent>
                    <DialogContentText>
                        {locale2.SET_NAME_AND_SURNAME[props.lang]}
                    </DialogContentText>
                    <TextField value={nameEdit.firstName} error={Boolean(err)} 
                            helperText={err} 
                                onChange={(e) => setNameEdit(c => ({ ...c, firstName: e.target.value }))}
                        autoFocus margin="dense" 
                                id="name" 
                                label={locale2.FIRST_NAME[props.lang]} 
                                type="text" fullWidth />
                    <TextField value={nameEdit.lastName} error={Boolean(err)} 
                        helperText={err} 
                        onChange={(e) => setNameEdit(c => ({ ...c, lastName: e.target.value }))}
                        margin="dense" id="lastname" 
                        label={locale2.LAST_NAME[props.lang]} 
                        type="text" fullWidth />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setFnOpen(false)}>
                        {locale2.CANCEL[props.lang]}
                    </Button>
                    <Button onClick={saveName} style={{
                        backgroundColor: MulwiColors.greenDark,
                        color: "white"
                    }} variant="contained">
                        {locale2.SAVE[props.lang]}
                    </Button>
                </DialogActions>
            </React.Fragment>
        )
    }

    function nameEditWaiter() {
        return nameEditSt.st === "wip" && (
            <CircularProgress className={classes.padded} />
        )
    }

    function nameEditOK() {
        return nameEditSt.st === "ok" && (
            <Check color="primary" fontSize="large" className={classes.padded} />
        )
    }

    function nameEditEX() {
        return nameEditSt.st === "ex" && (
            <React.Fragment>
                <DialogContent>
                    <DialogContentText>
                        {errToStr(nameEditSt.msg)}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button color="primary" onClick={() => setNameEditSt({ st: "", msg: "" })}>
                        {locale2.ONCE_AGAIN[props.lang]}
                    </Button>
                    <Button onClick={() => setFnOpen(false)} color="secondary">
                        <Close />
                    </Button>
                </DialogActions>
            </React.Fragment>
        )
    }

    async function saveName() {
        if (err) return
        let u = props.user
        u.Name = nameEdit.firstName + " " + nameEdit.lastName
        try {
            setNameEditSt({ st: "wip", msg: "" })
            await patchUserData(u)
            props.main.refresh()
            setNameEditSt({ st: "ok" })
            setTimeout(async () => {
                setFnOpen(false)
            }, 100)
        } catch (ex) {
            setNameEditSt({
                st: "ex",
                msg: ex
            })
        }
    }

    function nameStr() {
        if (!props.user) return ""
        if (props.user.Name) {
            return (props.user.Name)
        }
        return ""
    }


    return (
        <React.Fragment>
            {!props.external && (<Grid
                container
                spacing={3}
                alignItems="center"
            >
                <Grid item xs={4} style={{ color: MulwiColors.subtitleTypography }}>
                    {locale2.NAME[props.lang]}
                </Grid>
                <Grid item xs={6}>
                    {nameStr() || ""}
                </Grid>
                <Grid item xs={2}>
                    <Button onClick={openNameEdit} color="primary" size='small' 
                                style={{ marginLeft: 5 }} aria-label="edit">
                        <EditIcon />
                    </Button>
                </Grid>
            </Grid>)}
            <Dialog open={fnOpen} onClose={() => setFnOpen(false)} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">{locale2.YOUR_NAME[props.lang]}</DialogTitle>
                {nameEditForm()}
                {nameEditWaiter()}
                {nameEditOK()}
                {nameEditEX()}
            </Dialog>
        </React.Fragment>
    )
}