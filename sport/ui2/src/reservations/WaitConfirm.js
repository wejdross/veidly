import { Button, CircularProgress, Grid, Typography } from '@mui/material'
import React, { useEffect, useRef, useState } from 'react'
import { readRsvByToken } from '../apicalls/instructor.api'
import { sleep } from '../helpers'
import {stateDisplay} from '../harmonogram/trainingDetails'
import { Check } from '@mui/icons-material'
import { MulwiColors } from '../mulwiColors'
import { Link } from 'react-router-dom'
import { locale2 } from '../locale';

export function WaitConfirm(props) {

    const mountedRef = useRef(true)
    useEffect(() => {
        return () => { 
          mountedRef.current = false
        }
      }, [])

    const [r, setr] = useState(null)
    const [state, setState] = useState(false)

    useEffect(async () => {
        while(1) {
            if (!mountedRef.current) return
            setState(1)
            try {
                let r = await readRsvByToken(props.accessToken)
                r = JSON.parse(r)
                r = r.Rsv[0]
                setr(r)
                if(r.State == "hold" || r.State == "capture") {
                    setState(0)
                    return
                }
            } catch(ex) {
                console.log(ex)
            } 
            await sleep(2000)
        }
    }, [])

    return (
        <Grid container direction="column"
                style={{backgroundColor:MulwiColors.whiteBackground}}
                justify="center"
                spacing={4}
                alignItems="center">
            {state === 1 && (<Grid item>
                <CircularProgress/>
            </Grid>)}
            {state === 0 && (<Grid item>
                <Check style={{
                    color: MulwiColors.greenDark,
                }} fontSize="large"/>
            </Grid>)}
            {r && (<Grid item>
                {stateDisplay(r.State, r, props.lang)}
            </Grid>)}
            <Grid item>
                <Typography>
                   {locale2.YOU_MAY_CLOSE_THIS_WINDOW[props.lang]}
                </Typography>
            </Grid>
            <Grid item>
                <Typography>
                    {locale2.AT_ANY_TIME_YOU_CAN_CHECK[props.lang]}
                </Typography>
            </Grid>
            <Grid item>
            <Link to={"/rsv_details?type=token&id=" + props.accessToken} 
                    style={{textDecoration:"none"}}>
                <Button variant="contained" style={{
                     backgroundColor: MulwiColors.blueDark,
                     color: "white"
                }}>
                    {locale2.HERE[props.lang]}
                </Button>
            </Link>
            </Grid>
            <Grid item>
                <Typography>
                {locale2.OR_IF_YOU_HAVE_ACCOUNT_YOU_CAN_CHECK[props.lang]}
                </Typography>
            </Grid>
            <Grid item>
            <Link to="/trainings" style={{textDecoration:"none"}}>
                <Button variant="contained" style={{
                     backgroundColor: MulwiColors.blueDark,
                     color: "white"
                }}>
                    {locale2.YOUR_RSV_TAB[props.lang]}
                </Button>
            </Link>
            </Grid>
        </Grid>
    )
}