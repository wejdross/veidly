import { Button, CircularProgress, FormControl, TextField, Typography } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { resetPassword } from '../apicalls/user.api';
import CardWithBg from '../card/cardWithBg';
import { rmtoken } from '../helpers';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors'

export default function ResetPassword(props) {

    const [pass, setPass] = useState("")
    const [confirmPass, setConfirmPass] = useState("")
    const [passErr, setPassErr] = useState(false)
    const [workInProgress, setWorkInProgress] = useState(false)
    const [msg, setMsg] = useState("")
    const history = useHistory();

    let q = new URLSearchParams(window.location.search)
    let token = q.get("token")

    function validateForm() {
        let e = true
        if(!pass || pass !== confirmPass) {
            setPassErr(true)
            e = false
        } else {
            setPassErr(false)
        }
        return e
    }

    useEffect(validateForm)

    async function ResetPassword() {
        if(!validateForm()) return
        setWorkInProgress(true)
        try {
            await resetPassword(token, pass)
            rmtoken()
            history.push("/login")
        } catch(ex) {
            setMsg(locale2.COULDNT_RESET_PASSWOD[props.lang] + ex)
        } finally {
            setWorkInProgress(false)
        }
    }

    return (
        <CardWithBg>
            <FormControl>
                <Typography variant="h5">
                    {locale2.PASS_RESET[props.lang]}
                </Typography>
                <Typography style={{marginBottom: 15}}>
                    {locale2.CHANGE_PASSWORD[props.lang]}
                </Typography>
                <TextField 
                    error={passErr}
                    value={pass}
                    type="password"
                    label={locale2.PASSWORD[props.lang]}
                    onChange={(e) => setPass(e.target.value)}/>
                <TextField 
                    error={passErr}
                    value={confirmPass}
                    type="password"
                    label={locale2.CONFIRM_PASSWORD[props.lang]}
                    onChange={(e) => setConfirmPass(e.target.value)}/>
                <Button variant="contained" style={{
                    color: "white",
                    marginTop: 30,
                    backgroundColor: MulwiColors.blueDark
                }} onClick={ResetPassword}>
                        {(workInProgress && (
                            <CircularProgress style={{padding:5, width:30, height: 30}} color="secondary"/>
                        )) || (
                            locale2.RESET[props.lang]
                        )}
                </Button>
                <Typography color="error">
                    {msg}
                </Typography>
            </FormControl>
        </CardWithBg>
    )
}