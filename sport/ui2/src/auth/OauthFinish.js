import { Button, CircularProgress } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import { oauthActionWithCode } from '../apicalls/user.api';
import CardWithBg from '../card/cardWithBg';
import { settoken } from '../helpers';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';

export default function OauthFinish(props) {
    const [oauthErr, setOauthErr] = useState("")
    const history = useHistory();
    useEffect(() => {
        async function _() {
            let query = new URLSearchParams(window.location.search)
            let provider = query.get("oauth")
            let code = query.get("code")
            let state = query.get("state")
            if(provider && code) {
                try {
                    let token = await oauthActionWithCode(provider, code);
                    settoken(token)
                    await props.main.refresh()
                    if(state) {
                        history.push(state)
                    } else {
                        history.push("/")
                    }
                } catch(ex) {
                    setOauthErr(<React.Fragment>
                        <p>{locale2.FAILED_TO_LOGIN_VIA_OAUTH[props.lang] + provider}</p>
                        <p>{locale2.MAYBE_YOU_USE_PASS[props.lang]}</p>
                        <p>{locale2.ONCE_AGAIN[props.lang]}</p>
                        <Button><Link style={{textDecoration:"none", color: MulwiColors.blueDark}} to="/login">
                            {locale2.LOGIN[props.lang]}
                        </Link></Button>
                        <Button><Link style={{textDecoration:"none", color: MulwiColors.blueDark}} to="/register">
                            {locale2.REGISTER[props.lang]}
                        </Link></Button>
                    </React.Fragment>)
                }
            } else {
                setOauthErr(locale2.NO_OAUTH_CODE[props.lang])
            }
        }
        _()
    }, [])
    return (
        <CardWithBg img="/static/form-backgrounds/kosz.webp">
            <div style={{maxWidth: 300}}>
                {oauthErr || <CircularProgress/>}
            </div>
        </CardWithBg>
    )
}