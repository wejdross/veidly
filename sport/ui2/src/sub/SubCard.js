import { Button, Grid, Typography } from "@mui/material";
import React from "react";
import { useHistory } from "react-router";
import { daySuffix as dl } from "../harmonogram/trainingDetails";
import { prettyPrintCurrency } from "../helpers";
import { locale2 } from "../locale";

export function KeyVal(props) {
    return (<Grid container direction="row" spacing={2} 
                alignItems="center" justify="space-between">
        <Grid item>
            <Typography>{props.k}</Typography>
        </Grid>
        <Grid item>
            {!props.raw ? 
                (<Typography><strong>{props.v}</strong></Typography>) : props.v}
        </Grid>
    </Grid>)
}

export function SmContent(props) {

    let s = props.sm
    if(!s) return null

    return (<React.Fragment>
        {!props.noname && (<Typography variant="body2">
            <strong>{s.Name}</strong>
        </Typography>)}
        <KeyVal k={locale2.PRICE[props.lang]} 
            v={s.Price/100 + " " + prettyPrintCurrency(s.Currency)} />
        <KeyVal k={locale2.NUMBER_OF_ENTRIES[props.lang]} 
            v={s.MaxEntrances === -1 ? locale2.UNLIMITED[props.lang] : s.MaxEntrances} />
        <KeyVal k={locale2.PERIOD_OF_VALIDITY[props.lang]} 
            v={s.Duration + " " + dl(s.Duration, props.lang)} />
    </React.Fragment>)
}

export function SubCard(props) {

    const h = useHistory()
    function navigateToSmDetails() {
        let l = ("/sub_purch?instructorID=" 
            + s.InstructorID + "&smID=" + s.ID)
        h.push(l)
        return
    }

    let s = props.sm
    if(!s) return null

    return (<React.Fragment>
        <Button disabled={!props.user} style={{width: "100%"}} 
                variant="outlined" onClick={navigateToSmDetails}>
            <div style={{width: "100%"}}>
                <SmContent lang={props.lang} sm={s} />
            </div>
        </Button>
    </React.Fragment>)
}
