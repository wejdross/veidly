import React, { useEffect, useState } from 'react'
import { useHistory } from 'react-router'
import {
    Button, Checkbox, Divider, FormControlLabel,
    Grid, IconButton, TextField, Typography
} from '@mui/material'
import { dateToEpoch, dfInHours, epochToDate, prettyPrintCurrency, sprintf } from '../helpers'
import { prettyPrintDate } from '../harmonogram/trainingDetails'
import { Edit } from '@mui/icons-material'
import { MulwiColors } from '../mulwiColors'
import { getSupportedLocale, locale2 } from '../locale'
import { Link } from 'react-router-dom'
import { postReservation } from '../apicalls/instructor.api'
import { getPricingInfo } from '../apicalls/rsv'
import { getErrorDialog } from '../StatusDialog'
import { DcSelectModal } from '../dc/dcSelect'

export function CreateRsv(props) {
    const [day, setDay] = useState(new Date())
    const history = useHistory()
    const [pricingInfo, setPrcicingInfo] = useState(null)

    const [express, setExpress] = useState(false)

    const [name, setName] = useState("")
    const [email, setEmail] = useState("")
    const [phoneNumber, setPhoneNumber] = useState("")
    const [dc, setDc] = useState(null)

    function setStateFromQuery(l) {
        let query = new URLSearchParams(l.search)

        let _dateStart = query.get("dateStart")
        if (_dateStart) {
            let d = epochToDate(_dateStart)
            if (d && isFinite(d)) {
                //setWk(getWkFromMonth(d))
                setDay(d)
            }
        }

        if (dfInHours(new Date(), _dateStart) <= 24 && props.training.AllowExpress) {
            setExpress(true)
        }
    }

    useEffect(() => {
        setStateFromQuery(window.location)
    }, [props.training])


    useEffect(() => {
        return history.listen((location) => {
            setStateFromQuery(location)
        })
    }, [history])

    const [useSavedData, setUseSavedData] = useState(false)

    function getRsvRequest(ll, cc) {
        let dcID = null
        if (dc) dcID = dc.ID
        return {
            TrainingID: props.training.ID,
            Occurrence: day,
            NoRedirect: true,
            UserData: {
                Name: name,
                Language: ll,
                Country: cc
            },
            ContactData: {
                Email: email,
                Phone: phoneNumber
            },
            DcID: dcID,
            UseSavedData: useSavedData
        }
    }

    async function refreshPricingInfo() {
        try {
            let pi = await getPricingInfo(getRsvRequest())
            pi = JSON.parse(pi)
            setPrcicingInfo(pi)
        } catch (ex) {
            props.setInfo(getErrorDialog(locale2.FAILED_TO_DOWNLOAD_PRICING[props.lang], ex))
        }
    }

    useEffect(() => {
        refreshPricingInfo()
    }, [dc])

    const [errs, setErrs] = useState({ Name: "" })

    async function makeRsv() {

        if (!name) {
            setErrs(c => ({ ...c, Name: locale2.FIELD_IS_REQUIRED[props.lang] }))
            return
        } else {
            setErrs(c => ({ ...c, Name: "" }))
        }

        const [ll, cc] = getSupportedLocale()

        if (!props.training)
            return
        let r = getRsvRequest(ll, cc)
        try {
            let res = await postReservation(r)
            res = JSON.parse(res)
            window.location.href = res.Url
        } catch (ex) {
            props.setInfo(
                getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }


    let df = Math.floor(dfInHours(new Date(), day) / 24)

    return (<React.Fragment>
        <Grid container direction="column" spacing={1}>
            <Grid item>
                <Typography variant="h6">{locale2.TERM[props.lang]}</Typography>
            </Grid>
            <Grid item>
                <Grid container direction="row" spacing={2} alignItems="center">
                    <Grid item>
                        <Typography>{prettyPrintDate(day, props.lang)}</Typography>
                    </Grid>
                    <Grid item>
                        <IconButton onClick={() => {
                            let l = ("/instr/sched?instructorID=" +
                                props.training.InstructorID +
                                "&trainingID=" +
                                props.training.ID +
                                "&dateStart=" + dateToEpoch(day))
                            history.push(l)
                        }}>
                            <Edit />
                        </IconButton>
                    </Grid>
                </Grid>
            </Grid>

            <Grid item>
                <Typography variant="h6">
                    {locale2.YOUR_INFO[props.lang]}
                </Typography>
            </Grid>
            {props.user && props.user.Name && (<Grid item>
                <FormControlLabel
                    control={<Checkbox checked={useSavedData} onChange={(e) => {
                        if (e.target.checked) {
                            setName(props.user.Name)
                            setPhoneNumber(props.user.ContactData.Phone)
                            setEmail(props.user.ContactData.Email || props.user.Email)
                        }
                        setUseSavedData(e.target.checked)
                    }} />}
                    label={locale2.USE_ACCOUNT_DATA[props.lang]}
                />
            </Grid>)}
            <Grid item>
                <Grid container direction="column" spacing={3}>
                    <Grid item>
                        <TextField type="text" size="small"
                            disabled={useSavedData}
                            value={name}
                            fullWidth
                            helperText={errs.Name}
                            error={Boolean(errs.Name)}
                            onChange={(e) => setName(e.target.value)}
                            variant="outlined"
                            label={locale2.NAME_OR_NICK[props.lang]} />
                    </Grid>
                    {/* <Grid item>
                        <Grid container direction="row" spacing={3}>
                            <Grid item xs={6}>
                                <TextField type="text" size="small"
                                    value={email}
                                    disabled={useSavedData}
                                    onChange={(e) => setEmail(e.target.value)}
                                    variant="outlined" label="Email"/>
                                {!email && <Typography style={{fontSize: 11}} 
                                            variant="body2" color="textSecondary">
                                    <strong>{locale2.WE_ENCOURAGE_EMAIL[props.lang]}</strong>   
                                </Typography>}
                            </Grid>
                            <Grid item xs={6}>
                                <TextField type="text" size="small"
                                    value={phoneNumber}
                                    disabled={useSavedData} 
                                    onChange={(e) => setPhoneNumber(e.target.value)}
                                    variant="outlined" 
                                    label={locale2.PHONE[props.lang]}/>
                            </Grid>
                        </Grid>
                    </Grid> */}
                    <Grid item>
                        <TextField type="text" size="small"
                            value={email} fullWidth
                            disabled={useSavedData}
                            onChange={(e) => setEmail(e.target.value)}
                            variant="outlined" label="Email" />
                        {!email && <Typography style={{ fontSize: 11 }}
                            variant="body2" color="textSecondary">
                            <strong>{locale2.WE_ENCOURAGE_EMAIL[props.lang]}</strong>
                        </Typography>}
                    </Grid>
                    <Grid item style={{ marginBottom: 10 }}>
                        <Typography variant="body2" color="textSecondary" style={{ whiteSpace: "pre-wrap" }}>
                            {locale2.CONTACT_DATA_DISCLAIMER[props.lang]}
                            <Link to="/" style={{
                                textDecoration: "none",
                                color: MulwiColors.blueDark
                            }}>{locale2.READ_FURTHER[props.lang]}</Link>
                        </Typography>
                    </Grid>
                </Grid>
            </Grid>

            <Divider />
            {/* {pricingInfo && (<React.Fragment>
                <Grid item>
                    <Grid container direction="row" justify="space-between">
                        <Grid item>
                            {locale2.TRAINING_PRICE[props.lang]}
                        </Grid>
                        <Grid item>
                            <Typography>{props.training.Price / 100} {prettyPrintCurrency(props.training.Currency)}</Typography>
                        </Grid>
                    </Grid>
                    <Grid container direction="row" justify="space-between">
                        <Grid item>
                            {locale2.TRANSACTION_COST[props.lang]}
                        </Grid>
                        <Grid item>
                            <Typography>{pricingInfo.ProcessingFee / 100} {prettyPrintCurrency(props.training.Currency)}</Typography>
                        </Grid>
                    </Grid>
                    {pricingInfo.Dc ? (
                        <Grid container direction="row" justify="space-between">
                            <Grid item>
                                {locale2.DISCOUNT[props.lang]}
                            </Grid>
                            <Grid item>
                                <Typography>{pricingInfo.Dc.Discount}%</Typography>
                            </Grid>
                        </Grid>
                    ) : (
                        <DcSelectModal lang={props.lang} onChange={e => {
                            setDc(e)
                        }} trainingID={props.training.ID} />
                    )}
                    <Divider />
                    <Grid container direction="row" justify="space-between">
                        <Grid item>
                            {locale2.TOTAL_COST[props.lang]}
                        </Grid>
                        <Grid item>
                            <Typography><strong>{pricingInfo.TotalPrice / 100}</strong> {prettyPrintCurrency(props.training.Currency)}</Typography>
                        </Grid>
                    </Grid>
                </Grid>
            </React.Fragment>)} */}
            <Grid item>
                <Button variant="contained" style={{
                    backgroundColor: MulwiColors.greenDark,
                    color: "white"
                }} onClick={makeRsv} fullWidth>{locale2.RSV_SIGN_UP[props.lang]}</Button>
            </Grid>

            {/* {express && (
                <Grid item>
                    {props.training.ManualConfirm && (
                        <Typography><strong>
                            {locale2.CONFIRM_REQUIRED[props.lang]}
                        </strong></Typography>
                    ) || (
                            <Typography><strong>
                                {locale2.CONFIRM_NOT_REQUIRED[props.lang]}
                            </strong></Typography>
                        )}
                </Grid>
            )} */}

            {/* <Grid item>
                {express ? (
                    <React.Fragment></React.Fragment>
                ) : (
                    <Typography variant="body2" color="textSecondary">
                        {sprintf(locale2.YOU_HAVE_DAYS_UNTIL_RSV_FMT[props.lang], df)}
                        <br />{locale2.WE_WILL_CHARGE_AFTER[props.lang]}
                        {props.training.ManualConfirm && (
                            <React.Fragment>
                                <br />{locale2.IF_INSTR_AGRESS[props.lang]}
                            </React.Fragment>
                        )}</Typography>
                )}
            </Grid> */}
        </Grid>
    </React.Fragment>)
}