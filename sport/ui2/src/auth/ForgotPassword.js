import { Avatar, Button, CircularProgress, Grid, TextField, Typography } from '@mui/material';
import Divider from '@mui/material/Divider';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import validator from 'validator';
import { forgotPassword } from '../apicalls/user.api';
import CardWithBg from '../card/cardWithBg';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
import '../veidly-styles.css';


export default function ForgotPassword(props) {
    const [emailErr, setEmailErr] = useState(false)
    const [email, setEmail] = useState("")
    const [workInProgress, setWorkInProgress] = useState(false)
    const [msg, setMsg] = useState("")
    const [hdr, setHdr] = useState("")

    function validate() {
        if (!email || !validator.isEmail(email)) {
            setEmailErr(true)
            return false
        } else {
            setEmailErr(false)
            return true
        }
    }

    useEffect(validate)

    async function sendLink() {
        if (!validate()) return
        setWorkInProgress(true)
        try {
            await forgotPassword(email)
            setMsg(locale2.SENT_EMAIL[props.lang])
            setHdr(locale2.SUCCESS[props.lang])
        } catch (ex) {
            setMsg(locale2.COULDNT_SEND_EMAIL[props.lang] + ex)
        } finally {
            setWorkInProgress(false)
        }
    }
    return (
        <CardWithBg
            header={locale2.PASSWORD_RECOVERY[props.lang]}
            subheader={locale2.PASSWORD_RECOVERY_SUB[props.lang]}
        >
            <TextField style={{ marginBottom: 10 }}
                error={emailErr}
                value={email}
                fullWidth
                variant="outlined"
                onChange={(e) => setEmail(e.target.value)}
                id="standard-basic" label="Email" type="email" />
            <Button
                style={{
                    marginTop: 30,
                    color: "white",
                    backgroundColor: emailErr ? "grey" : MulwiColors.blueDark
                }}
                fullWidth
                onClick={sendLink}
                disabled={emailErr}
                variant="contained">
                {
                    (workInProgress && (
                        <CircularProgress style={{ padding: 5, width: 30, height: 30 }} color="secondary" />
                    )) || (
                        locale2.SEND_EMAIL[props.lang]
                    )
                }
            </Button>
            <Typography variant="h6" color="primary">
                {hdr}
            </Typography>
            <Typography >
                {msg}
            </Typography>
            <Divider style={{ marginTop: 40, fontSize: 12, fontWeight: 400 }}>{locale2.OR[props.lang]}</Divider>

            {!hdr && (
                <>
                    <Link to="/login" style={{ textDecoration: "none" }}>
                        <Button
                            style={{ marginTop: 10, backgroundColor: "white", width: "100%", textDecoration: "none" }}
                            variant="contained" color="secondary" >
                            <Avatar src={"veidly-logo-mini-black.svg"} style={{ height: 18, width: 18 }} />
                            <Typography style={{ marginLeft: 5, color: "Black", lineHeight: "1.2", textDecoration: "none" }}>
                                {locale2.LOGIN[props.lang]} {locale2.WITH[props.lang]} Veidly
                            </Typography>
                        </Button>
                    </Link>

                </>
            )}

        </CardWithBg>
    )
}