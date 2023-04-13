
import { Close } from '@mui/icons-material';
import PhotoCamera from '@mui/icons-material/PhotoCamera';
import {
    Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, Grid, IconButton, Typography
} from '@mui/material';
import React, { useState } from 'react';
import { putImg } from '../apicalls/user.api';
import { randomString } from '../helpers';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';

export default function AvatarEdit(props) {

    const [open, _setOpen] = useState(false)

    const [loading, setLoading] = useState(false)

    const [err, setErr] = useState("")

    const [fd, setFd] = useState(null)

    function setOpen(v) {
        if (v === false) setErr("")
        _setOpen(v)
    }

    function onchange(e) {
        if (!e.target.files || !e.target.files[0]) return
        setOpen(true)
        setLoading(true)
        let formData = new FormData()
        formData.append("image", e.target.files[0])
        setFd(formData)
        // load img to preview
        var fr = new FileReader()
        fr.onload = function () {
            document.getElementById("prevavatar").src = fr.result
            setLoading(false)
        }
        fr.readAsDataURL(e.target.files[0])
    }

    async function save() {
        if (!fd) {
            return
        }
        setLoading(true)
        try {
            if(props.putImg)
                await props.putImg(fd)
            else
                await putImg(fd)
            if(props.refreshInstr) 
                props.main.refreshInstructor()
            props.main.refreshUser()
            setOpen(false)
        } catch (ex) {
            setErr(ex)
        } finally {
            setLoading(false)
        }
    }

    const id = randomString(12)

    function dialog() {
        return (<Dialog aria-labelledby="form-dialog-title" open={open} onClose={() => setOpen(false)}>
            <DialogTitle id="form-dialog-title">
                {locale2.PREVIEW[props.lang]} {loading && <CircularProgress style={{ marginBottom: -10, marginLeft: 10 }} />}
            </DialogTitle>
            <img style={{
                padding: 20,
                maxWidth: 300,
                maxHeight: 300
            }} id="prevavatar" alt="" />
            {err && (<DialogContent>
                <Typography style={{ textAlign: "center" }}>
                    {locale2.SOMETHING_WENT_WRONG[props.lang]}
                </Typography>
                <Typography style={{ textAlign: "center" }}>
                    {
                        (err == 413) ? <center>
                            {locale2.PICTURE_TOO_LARGE[props.lang]}
                            <br />
                            {locale2.CLOSE_AND_TRY_AGAIN[props.lang]}
                        </center> : `${locale2.SOMETHING_WENT_WRONG[props.lang]}`
                    }
                </Typography>
            </DialogContent>)}
            <DialogActions>
                <Button color="secondary" onClick={() => setOpen(false)}>
                    {locale2.CANCEL[props.lang]} <Close />
                </Button>
                {(err && "err") || (<Button onClick={save} color="primary">
                    {locale2.SAVE[props.lang]}
                </Button>)}
            </DialogActions>
        </Dialog>)
    }

    return (
        props.fab ? (
            <React.Fragment>
                {dialog()}
                <input
                    onChange={onchange}
                    accept="image/*"
                    style={{ display: "none" }}
                    id={id} type="file" />
                <label htmlFor={id}>
                    {props.fab(id)}
                </label>
            </React.Fragment>
        ) : (
            <Grid container
                alignItems="center"
                spacing={3}>
                {dialog()}
                <Grid item xs={10}>
                    <span style={{ color: MulwiColors.subtitleTypography }}>Avatar</span>
                </Grid>
                <Grid item xs={2} >
                    <Grid item>
                        <input
                            onChange={onchange}
                            accept="image/*"
                            style={{ display: "none" }}
                            id={id} type="file" />
                        <label htmlFor={id}>
                            <IconButton
                                style={{
                                    marginTop: -10,
                                    marginBottom: -10,
                                    marginLeft: 10
                                }}
                                color="primary" aria-label="upload picture" component="span">
                                <PhotoCamera />
                            </IconButton>
                        </label>
                    </Grid>
                </Grid>
            </Grid>
        )
    )
}