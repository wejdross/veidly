import React, { useState } from 'react'
import { Button, Divider, Grid, Typography } from '@mui/material'
import { prettyPrintCurrency } from '../helpers'
import { daySuffix as dl, prettyPrintDay } from '../harmonogram/trainingDetails'
import { MulwiColors } from '../mulwiColors'
import {Link} from 'react-router-dom'
import { getErrorDialog } from '../StatusDialog'
import { postSub } from '../apicalls/sm'
import { locale2 } from '../locale'

export function CreateSub(props) {

    async function buyPass() {
        if(!props.sm) return
        try {
            let res = await postSub(props.sm.ID)
            res = JSON.parse(res)
            window.location.href = res.Url
        } catch (ex) {
            props.setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    function smTerm() {
        let now = new Date()
        now.setDate(now.getDate() + props.sm.Duration)
        return now
    }

    if(!props.sm || !props.user) return null

    return (<React.Fragment>
        <Grid container direction="column" spacing={1}>
            <Grid item>
                <Typography variant="h6">{locale2.CARNET[props.lang]}</Typography>
            </Grid>
            <Grid item>
                <Grid container direction="row" spacing={2} 
                                alignItems="center" justify="space-between">
                    <Grid item>
                        <Typography>{locale2.PERIOD_OF_VALIDITY[props.lang]}</Typography> 
                    </Grid>
                    <Grid item>
                        <strong>{props.sm.Duration} {dl(props.sm.Duration, props.lang)}
                        </strong>
                    </Grid>
                </Grid>
            </Grid>
            <Grid item>
                <Grid container direction="row" spacing={2} alignItems="center" 
                            justify="space-between">
                    <Grid item>
                        <Typography>{locale2.CARNET_VALID_UNTIL[props.lang]}</Typography> 
                    </Grid>
                    <Grid item>
                        <strong>
                            {prettyPrintDay(smTerm())}
                        </strong>
                    </Grid>
                </Grid>
            </Grid>
            <Grid item>
                <Grid container direction="row" spacing={2} 
                            alignItems="center" justify="space-between">
                    <Grid item>
                        <Typography>{locale2.NUMBER_OF_ENTRIES[props.lang]}</Typography> 
                    </Grid>
                    <Grid item>
                        <strong>
                            {props.sm.MaxEntrances == -1 ? 
                                locale2.UNLIMITED[props.lang]
                                    : props.sm.MaxEntrances}
                        </strong>
                    </Grid>
                </Grid>
            </Grid>


            <Grid item>
                <Typography variant="h6">
                    {locale2.CARNET_OWNER_INFO[props.lang]}
                </Typography>
            </Grid>

            <Grid item>
                <Grid container direction="row" spacing={2} alignItems="center" justify="space-between">
                    <Grid item>
                        <Typography>{locale2.NAME[props.lang]}</Typography> 
                    </Grid>
                    <Grid item>
                        {props.user.Name ? <strong>
                            {props.user.Name}
                        </strong> : <span style={{
                            color: MulwiColors.subtitleTypography
                        }}> 
                            {locale2.ANON[props.lang]}
                        </span>}
                    </Grid>
                </Grid>
            </Grid>
            <Grid item>
                <Grid container direction="row" spacing={2} alignItems="center" justify="space-between">
                    <Grid item>
                        <Typography>{locale2.EMAIL[props.lang]}</Typography> 
                    </Grid>
                    <Grid item>
                        {props.user.ContactData.Email ? <strong>
                            {props.user.ContactData.Email}
                        </strong> : <span style={{
                            color: MulwiColors.subtitleTypography
                        }}> 
                            {locale2.NONE[props.lang]}
                        </span>}
                    </Grid>
                </Grid>
            </Grid>
            <Grid item>
                <Grid container direction="row" spacing={2} 
                            alignItems="center" justify="space-between">
                    <Grid item>
                        <Typography>{locale2.PHONE[props.lang]}</Typography> 
                    </Grid>
                    <Grid item>
                        {props.user.ContactData.Phone ? <strong>
                            {props.user.ContactData.Phone}
                        </strong> : <span style={{
                            color: MulwiColors.subtitleTypography
                        }}> 
                            {locale2.NONE[props.lang]}
                        </span>}
                    </Grid>
                </Grid>
            </Grid>

            <Grid item style={{marginBottom: 10}}>
                <Typography variant="body2" color="textSecondary" style={{whiteSpace:"pre-wrap"}}>
                    {locale2.CONTACT_DATA_DISCLAIMER[props.lang]}
                    <Link to="/" style={{
                        textDecoration:"none",
                        color: MulwiColors.blueDark
                    }}> {locale2.READ_FURTHER[props.lang]}</Link>
                </Typography>
            </Grid>

            <Divider/>

            <React.Fragment>
                <Grid item>
                    <Grid container direction="row" justify="space-between">
                        <Grid item>
                            {locale2.CARNET_PRICE[props.lang]}
                        </Grid>
                        <Grid item>
                            <Typography>{props.sm.Price/100} {prettyPrintCurrency(props.sm.Currency)}</Typography>
                        </Grid>
                    </Grid>
                    <Grid container direction="row" justify="space-between">
                        <Grid item>
                            {locale2.TRANSACTION_COST[props.lang]}
                        </Grid>
                        <Grid item>
                            <Typography>{props.sm.ProcessingFee/100} {prettyPrintCurrency(props.sm.Currency)}</Typography>
                        </Grid>
                    </Grid>
                    <Divider />
                    <Grid container direction="row" justify="space-between">
                        <Grid item>
                            {locale2.TOTAL_COST[props.lang]}
                        </Grid>
                        <Grid item>
                            <Typography><strong>{(props.sm.Price + props.sm.ProcessingFee)/100}</strong> {prettyPrintCurrency(props.sm.Currency)}</Typography>
                        </Grid>
                    </Grid>
                </Grid>
            </React.Fragment>

            <Grid item>
                <Button variant="contained" style={{
                    backgroundColor:MulwiColors.greenDark,
                    color: "white"
                }} onClick={buyPass} fullWidth>
                    {locale2.BUY_CARNET[props.lang]}
                </Button>
            </Grid>
        </Grid>
    </React.Fragment>)
}