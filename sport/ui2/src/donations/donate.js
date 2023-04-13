import React, { useEffect, useState } from "react";
import Slider from '@mui/material/Slider';
import { Alert, Button, Card, CardActions, CardContent, CardHeader, Dialog, DialogActions, DialogContent, DialogTitle, Grid, Modal, Paper, TextField, Typography } from "@mui/material";
import CardWithBg from "../card/cardWithBg";
import { MulwiColors } from "../mulwiColors";
import { getErrorDialog, getNullDialog, StatusDialog } from '../StatusDialog'
import { donate } from "../apicalls/donate";
import { locale2 } from "../locale";
import isEmail from "validator/lib/isEmail";
import { validateDate } from "@mui/x-date-pickers/internals";


export function DonateForm(props) {
    const [donation, setDonation] = useState(50)
    const [email, setEmail] = useState("")
    const [err, setErr] = useState(null)

    const [info, setInfo] = useState(getNullDialog())

    async function handleDonate() {
        if (err || !validate(email))
            return
        try {
            let rawres = await donate(donation * 100, email)
            let res = JSON.parse(rawres)
            window.location.href = res.Url
        } catch (ex) {
            setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG["pl"], ex))
        }
    }

    function validate(v) {
        if (!(v && isEmail(v))) {
            setErr(locale2.BROKEN_EMAIL[props.lang])
            return false
        }
        setErr(null)
        return true
    }

    return (
        <React.Fragment>
            <StatusDialog lang={"pl"} info={info} setInfo={setInfo} />
            <Grid style={{
                marginTop: 50,
                padding: 10
            }} container direction="row" justifyContent={"center"}>
                <Grid item>

                    <CardWithBg style={{
                        maxWidth: 400
                    }}
                        header={locale2.SUPPORT_VEIDLY[props.lang]}
                        subheader={locale2.FREE_TO_USE_APP[props.lang]}
                    >
                        <CardContent>
                            <Grid
                                container
                                spacing={3}>
                                <Grid item xs={9}>
                                    <Slider
                                        aria-label="Money"
                                        defaultValue={50}
                                        valueLabelDisplay="auto"
                                        step={1}
                                        min={5}
                                        max={500}
                                        value={donation}
                                        onChange={(e) => setDonation(e.target.value)}
                                        style={{
                                            color: MulwiColors.blueDark,
                                            marginLeft: 10,
                                        }}
                                    />

                                </Grid>
                                <Grid item xs={3}>
                                    <Typography variant="h6" align="right">
                                        {donation} zł
                                    </Typography>
                                </Grid>
                            </Grid>
                            <TextField
                                size="small"
                                label={locale2.NOTIFY_EMAIL[props.lang]}
                                variant="outlined"
                                value={email}
                                fullWidth
                                style={{
                                    marginLeft: 0,
                                    marginRight: 0,
                                    paddingLeft: 0,
                                    paddingRight: 0,
                                }}
                                error={err}
                                onChange={e => {
                                    let v = e.target.value
                                    setEmail(v)
                                    validate(v)
                                }} />
                        </CardContent>
                        <CardActions>
                            <Button variant="contained" fullWidth disabled={err} style={{
                                backgroundColor: MulwiColors.greenDark
                            }}
                                onClick={() => handleDonate()} >
                                {locale2.SUPPORT_VEIDLY[props.lang]}
                            </Button>
                        </CardActions>
                    </CardWithBg>
                </Grid>
            </Grid>
        </React.Fragment>
    )
}
export default function Donate(props) {
    const [donation, setDonation] = useState(50)
    const [hovered, setHovered] = useState(false)
    const [email, setEmail] = useState("")
    const [err, setErr] = useState(null)

    const [info, setInfo] = useState(getNullDialog())
    const [open, setOpen] = useState(false)

    async function handleDonate() {
        if (err || !validate(email))
            return
        try {
            let rawres = await donate(donation * 100, email)
            let res = JSON.parse(rawres)
            window.location.href = res.Url
        } catch (ex) {
            setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG["pl"], ex))
        }
    }

    useEffect(() => {
        setOpen(props.open)
    }, [props.open])

    function validate(v) {
        if (!(v && isEmail(v))) {
            setErr(locale2.BROKEN_EMAIL[props.lang])
            return false
        }
        setErr(null)
        return true
    }

    return (
        <React.Fragment>
            <StatusDialog lang={"pl"} info={info} setInfo={setInfo} />
            <Dialog open={open} onClose={() => setOpen(false)}>
                <DialogTitle>
                {locale2.SUPPORT_VEIDLY[props.lang]}
                </DialogTitle>
                <DialogContent>
                    <Typography>
                    {locale2.FREE_TO_USE_APP[props.lang]}
                    </Typography>
                    <Grid
                        container
                        spacing={3}>
                        <Grid item xs={9}>
                            <Slider
                                aria-label="Money"
                                defaultValue={500}
                                valueLabelDisplay="auto"
                                step={1}
                                min={5}
                                max={500}
                                value={donation}
                                onChange={(e) => setDonation(e.target.value)}
                                style={{
                                    color: MulwiColors.blueDark,
                                    marginTop: 20,
                                }}
                            />

                        </Grid>
                        <Grid item xs={3}>
                            <Typography variant="h6" style={{
                                marginTop: 15,
                            }}>
                                {donation} zł
                            </Typography>
                        </Grid>
                    <Grid item xs={12}>

                    <TextField
                        size="small"
                        label={locale2.NOTIFY_EMAIL[props.lang]}
                        variant="outlined"
                        value={email}
                        fullWidth
                        error={err}
                        onChange={e => {
                            let v = e.target.value
                            setEmail(v)
                            validate(v)
                        }} />
                        </Grid>
                    </Grid>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpen(false)}>
                        {locale2.CLOSE[props.lang]}
                    </Button>
                    <Button variant="contained" disabled={err} style={{
                        backgroundColor: MulwiColors.greenDark
                    }}
                        onClick={() => handleDonate()} >
                        {locale2.SUPPORT[props.lang]}
                    </Button>
                </DialogActions>
            </Dialog>
        </React.Fragment>
    )
}
