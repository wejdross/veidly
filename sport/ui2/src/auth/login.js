import React, { useState, useEffect } from 'react';
import {
    Button, Grid, TextField, Typography,
    CircularProgress, Avatar
} from '@mui/material';
import { Link, useHistory } from 'react-router-dom';
import { getOauthGoogleUrl, userLogin } from '../apicalls/user.api';
import { settoken } from '../helpers';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
import Divider from '@mui/material/Divider';
import StickyFooter from '../Footer';
import '../veidly-styles.css'
import CardWithBg from '../card/cardWithBg';

export const defaultRedirect = "/profile?fromlogin=true"

export default function Login(props) {

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [msg, setMsg] = useState("")
    const [busy, setBusy] = useState(false)

    let q = new URLSearchParams(window.location.search);
    let _email = q.get("email")
    if (_email && !email) {
        setEmail(_email)
    }

    useEffect(() => {
        window.scrollTo(0, 0)
      }, [])

    const history = useHistory()

    let login = async (e) => {
        e.preventDefault();
        setBusy(true)
        try {
            let res = await userLogin({
                email: email,
                password: password
            });
            settoken(res);
            await props.main.refresh()
            let q = new URLSearchParams(window.location.search);
            let returnUrl = q.get("return_url");
            if (returnUrl) {
                history.push(returnUrl)
                //window.location = returnUrl;
            } else {
                history.push(defaultRedirect)
                //window.location = "/";
            }
        } catch (ex) {
            setMsg(locale2.COULDNT_LOGIN[props.lang])
        } finally {
            setBusy(false)
        }
    }


    async function loginGoogle() {
        let q = new URLSearchParams(window.location.search);
        let returnUrl = q.get("return_url") || defaultRedirect;
        try {
            let url = await getOauthGoogleUrl(returnUrl);
            window.location.replace(url);
        } catch (ex) {
            console.log(ex)
        }
    }

    return (
        <CardWithBg
            header={locale2.LOGIN[props.lang]}
            subheader={locale2.WELCOME_AGAIN[props.lang]}
        >

                        <TextField
                            style={{marginBottom: 20}}
                            id="login" label="Email"
                            variant='outlined'
                            InputLabelProps={{
                                classes: {
                                    focused: "focused",
                                }
                            }}
                            value={email}
                            fullWidth
                            onChange={e => setEmail(e.target.value)} />
                        <TextField style={{ marginBottom: 30 }}
                            id="password" label={locale2.PASSWORD[props.lang]} type="password"
                            variant='outlined'
                            value={password}
                            fullWidth
                            InputLabelProps={{
                                classes: {
                                    focused: "focused",
                                }
                            }}
                            onChange={e => setPassword(e.target.value)} />
                        <Link to="/forgot_password" style={{ textDecoration: "none", color: MulwiColors.blueDark }}>
                            {locale2.FORGOT_PASSWORD[props.lang]}
                        </Link>

                        <Button
                            style={{ marginTop: 30, backgroundColor: MulwiColors.blueDark }}
                            onClick={login}
                            fullWidth
                            variant="contained" color="primary">
                            {(busy && (
                                <CircularProgress style={{ padding: 5, width: 30, height: 30 }} color="secondary" />
                                )) || (
                                    locale2.CONTINUE[props.lang]
                                    )}
                        </Button>
                        <Typography style={{
                            marginTop: 10,
                            color: MulwiColors.redError
                        }}>
                            {msg}
                        </Typography>
                        <div  style={{marginTop: 40, marginBottom: 10}}>
                            <Divider style={{fontSize: 12, fontWeight: 400}}>{locale2.OR[props.lang]}</Divider>
                        </div>
                        <Button
                            style={{ marginTop: 10, backgroundColor: "white" }}
                            variant="contained" color="secondary"
                            onClick={loginGoogle}
                            fullWidth
                            >
                                <Avatar src={"google.png"} style={{ height: 18, width: 18 }} />
                                <Typography style={{ marginLeft: 5, color: "Black", lineHeight: "1.2" }} align={'center'}>
                                    Google
                                </Typography>
                        </Button>
                        </CardWithBg>
    )
}