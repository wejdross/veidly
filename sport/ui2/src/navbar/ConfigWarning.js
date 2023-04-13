import { Button } from '@mui/material';
import { Alert, AlertTitle } from '@mui/lab';
import React from 'react';
import { Link } from 'react-router-dom';
import { MulwiColors } from '../mulwiColors';
import { locale2 } from '../locale';

export function ConfigWarning(props) {

    if (!window.localStorage.getItem("mulwii_token")) {
        return null
    }

    if(!props.instructor) return null

    let c = props.instructor.Config

    if(!c) return null

    if(c & 4) {
        return (<Alert severity="warning">
            <AlertTitle>
                {locale2.TEMP_INACTIVE[props.lang]}
            </AlertTitle>
        </Alert>)
    }

    if(window.location.pathname.startsWith("/configure")) 
        return null

    return (<Alert severity="error" 
        action={<Link style={{
            textDecoration:"none"
        }} to="/configure">
            <Button style={{
                marginLeft: 10,
                color: "white",
                backgroundColor: MulwiColors.blueDark
            }}>{locale2.GOTO_CONFIG[props.lang]}</Button>
        </Link>}>
        <AlertTitle>
            {locale2.MISSING_INFO[props.lang]}
        </AlertTitle>
        
    </Alert>)
}