import { Typography } from '@mui/material';
import React  from 'react';
import { useHistory } from 'react-router-dom';
import CardWithBg from '../card/cardWithBg';
import { locale2 } from '../locale';

export default function RegisterFinish(props) {
    
    let query = new URLSearchParams(window.location.search)
    let status = query.get("state")
    let x = query.get("return_url")
    if(x) {
        x = escape(x);
    } else {
        x = "";
    }
    let err = ""
    const history = useHistory();
    switch(status) {
        case "200":
            history.push("/login?email=" + query.get("email") + "&return_url=" + x)
            return null
        case "401":
            err = locale2.TOKEN_HAS_EXPIRED[props.lang]
            break
        case "501":
            err = locale2.TWOFA_NOT_SUPPORTED[props.lang]
            break
        default:
            err = locale2.UNEXPECTED_ERROR[props.lang]
            break
    }
    return (
        <CardWithBg img="/static/form-backgrounds/tenis.webp">
            <Typography color="error">
                {err}
            </Typography>
        </CardWithBg>
    )
} 