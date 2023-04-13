import { Grid, Typography, useMediaQuery, useTheme } from '@mui/material';
import React from 'react';
import {
    getRsvStatus, getRsvStatusColor, prettyPrintDateRange,
    prettyPrintRsvDecision
} from '../harmonogram/trainingDetails';
import { MulwiColors } from '../mulwiColors';
import { locale2 } from '../locale';

export function RsvInfo(props) {

    const t = useTheme()
    const isLowRes = useMediaQuery(t.breakpoints.down('sm'))

    const contentMargin = 20

    if(!props.rsv) return null
    return (<React.Fragment>
        <Grid container direction="column" 
                alignItems={isLowRes ? "center" : "flex-start"}
                spacing={2}>
            <Grid item>
                <Typography variant="body2" 
                        style={{color: MulwiColors.subtitleTypography}}>
                    {locale2.TERM[props.lang]}
                </Typography>
            </Grid>
            <Grid item style={{marginLeft: contentMargin}}>
                <Typography>
                    {prettyPrintDateRange(
                        new Date(props.rsv.DateStart), 
                        new Date(props.rsv.DateEnd), 
                        0, 0, props.lang)}
                </Typography>
            </Grid>
            <Grid item>
                <Typography variant="body2" 
                    style={{color: MulwiColors.subtitleTypography}}>
                        {locale2.STATUS[props.lang]}
                    </Typography>
            </Grid>
            <Grid item style={{marginLeft: contentMargin}}>
                <Typography style={{
                    color: getRsvStatusColor(props.rsv)
                }}>{getRsvStatus(props.rsv, props.lang)}</Typography>
            </Grid>
            {/* <Grid item>
                <Typography variant="body2" 
                    style={{color: MulwiColors.subtitleTypography}}>
                        {locale2.INSTRUCTOR_DECISION[props.lang]}
                    </Typography>
            </Grid>
            
            <Grid item style={{marginLeft: contentMargin}}>
                  {prettyPrintRsvDecision(props.rsv, props.lang)}
              </Grid> */}
            {/* <Grid item style={{
                textAlign: isLowRes ? "center" : "left"
            }}>
                    {props.rsv.Training.ManualConfirm ? (
                        <Typography><strong>
                            {locale2.CONFIRM_REQUIRED[props.lang]}
                        </strong></Typography>
                    ) : (
                        <Typography><strong>
                            {locale2.CONFIRM_NOT_REQUIRED[props.lang]}
                        </strong></Typography>
                    )}
            </Grid> */}
        </Grid>
    </React.Fragment>)
}