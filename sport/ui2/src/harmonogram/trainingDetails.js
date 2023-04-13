import EditIcon from '@mui/icons-material/Edit';
import {
    Avatar, Button, Checkbox, Dialog, FormControl, FormControlLabel, Grid, InputAdornment, InputLabel, List, ListItem, MenuItem, OutlinedInput, TextField, Typography
} from '@mui/material';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import { DateTimePicker, MobileDateTimePicker } from '@mui/x-date-pickers';
import React, { useEffect, useState } from 'react';
import {
    getSchedule, getTrainings, patchTraining
} from '../apicalls/instructor.api';
import { putOcc } from '../apicalls/occ';
import { explainTags, fetchGeo, getTags } from '../apicalls/user.api';
import { diffs, getDiffsByLang } from '../diffs';
import {
    eachFdayIndex, getSupportedLanguage,
    getSupportedLocale, locale2
} from '../locale';
import { MulwiColors } from '../mulwiColors';
import AutocompleteEdit from '../profile/AutocompleteEditor';
import DynAtcEdit from '../profile/DynAtcEditor';
import GenericInlineEdit from '../profile/GenericInlineEdit';
import ModalEdit from '../profile/ModalEdit';
import TextAreaModalEdit from '../profile/TextAreaModalEdit';
import { TrainingSummary } from '../reservations/trainingSummary';
import { getErrorDialog } from '../StatusDialog';
import { DcEdit } from './DcEdit';
import { GmapEditor } from './gmapEditor';
import { GroupEdit } from './GroupEdit';
import TrainingImageEditor from './imageEditor';
import { OccDisplay } from './occDisplay';
import { RsvInteractMenu } from './rsvInteract';
import { SmEdit } from './SmEdit';
import { SubInteractMenu } from './subInteract';

export function prettyPrintHr(d) {
    let h = d.getHours()
    let str = ""
    if (h < 10) {
        str += "0"
    }
    str += String(h) + ":"
    let m = d.getMinutes()
    if (m < 10) {
        str += "0"
    }
    str += String(m)
    return str
}

export function daySuffix(d, lang) {
    d = Math.abs(d)
    if (d === 1) {
        return locale2.DAY[lang]
    }
    return locale2.DAYS[lang]
}

export function prettyPrintDateRange(start, end, rpt, noday, lang) {
    if (!start.toDateString) start = new Date(start)
    if (!end.toDateString) end = new Date(end)
    if (rpt) {
        if (rpt === 7) {
            if (start.toDateString() === end.toDateString()) {
                let day = ""
                switch (start.getDay()) {
                    case 1:
                        day = eachFdayIndex.MO[lang]
                        break
                    case 2:
                        day = eachFdayIndex.TU[lang]
                        break
                    case 3:
                        day = eachFdayIndex.WE[lang]
                        break
                    case 4:
                        day = eachFdayIndex.TH[lang]
                        break
                    case 5:
                        day = eachFdayIndex.FR[lang]
                        break
                    case 6:
                        day = eachFdayIndex.SA[lang]
                        break
                    case 7:
                        day = eachFdayIndex.SU[lang]
                        break
                }
                return locale2.EACH[lang] + day + locale2.IN_HRS[lang]
                    + prettyPrintHr(start) + " - " + prettyPrintHr(end)
            } else {
                return prettyPrintDateRange(start, end, 0, 0, lang)
                    + locale2.EVERY_WEEK[lang]
            }
        }
        return prettyPrintDateRange(start, end, 0, 0, lang)
            + locale2.REPEATED_EVERY[lang]
            + " " + rpt + " " + daySuffix(rpt, lang)
    }
    if (!start || !end) return null
    if (start.toDateString() === end.toDateString()) {
        if (noday)
            return prettyPrintHr(start) + " - " + prettyPrintHr(end)
        return prettyPrintHr(start) + " - " + prettyPrintHr(end) + " " + prettyPrintDay(start)
    }
    if (noday) {
        return prettyPrintHr(start) + " - " + prettyPrintHr(end) + " " + prettyPrintShort(end)
    }
    return prettyPrintHr(start) + " " + prettyPrintDay(start) + " - "
        + prettyPrintHr(end) + " " + prettyPrintDay(end)
}

export function prettyPrintDay(d) {
    return d.toLocaleDateString(getSupportedLocale().join("-"), {
        weekday: "long",
        year: "numeric",
        month: "2-digit",
        day: "numeric"
    })
}

export function prettyPrintShort(d) {
    return d.toLocaleDateString(getSupportedLocale().join("-"), {
        weekday: "short",
        month: "2-digit",
        day: "numeric"
    })
}

export function prettyPrintDate(d) {
    return d.toLocaleDateString(getSupportedLocale().join("-"), {
        weekday: "long",
        year: "numeric",
        month: "2-digit",
        day: "numeric"
    }) + " " + prettyPrintHr(d)
}

function AgeEditor(props) {

    const [minAge, setMinAge] = useState(0)
    const [maxAge, setMaxAge] = useState(0)

    return (
        <ModalEdit
            lang={props.lang}
            title={props.title + ' - ' + locale2.AGE[props.lang]}
            value={(!props.minAge && !props.maxAge && locale2.ANY[props.lang]) ||
                (props.minAge || "") + " - " + (props.maxAge || "")}
            label={<Typography style={{ color: "gray" }}>{locale2.AGE[props.lang]}</Typography>}
            onSave={() => {
                props.onChange && props.onChange(Number(minAge), Number(maxAge))
            }} content={<React.Fragment>
                <Grid container spacing={2} style={{ marginTop: 20, marginBottom: 20 }}>
                    <Grid item>
                        <TextField
                            type="number"
                            label={locale2.MIN_AGE[props.lang]}
                            value={minAge ? String(minAge) : ""}
                            onChange={(e) => setMinAge(e.target.value)}
                        />
                    </Grid>
                    <Grid item>
                        <TextField
                            type="number"
                            label={locale2.MAX_AGE[props.lang]}
                            value={maxAge ? String(maxAge) : ""}
                            onChange={(e) => setMaxAge(e.target.value)}
                        />
                    </Grid>
                </Grid>
            </React.Fragment>}
            onOpen={() => {
                setMinAge(props.minAge)
                setMaxAge(props.maxAge)
            }}
        />
    )
}

function PriceEditor(props) {

    const currencies = [
        {
            value: 'PLN',
            label: 'zł',
        },
        // {
        //     value: 'USD',
        //     label: '$',
        // },
        // {
        //     value: 'EUR',
        //     label: '€',
        // },
        // :(
        // {
        //   value: 'BTC',
        //   label: '฿',
        // },
        // {
        //     value: 'JPY',
        //     label: '¥',
        // },
    ]

    const [price, setPrice] = useState(0)
    const [currency, setCurrency] = useState("")

    return (
        <ModalEdit
            lang={props.lang}
            title={props.title + ' - ' + locale2.PRICE[props.lang]}
            value={props.price / 100 + " " + props.currency}
            label={<Typography style={{ color: "gray" }}>{locale2.PRICE[props.lang]}</Typography>}
            onSave={() => {
                props.onChange && props.onChange(Math.round(price * 100), currency)
            }} content={<React.Fragment>
                <Grid container spacing={2} style={{ marginTop: 20, marginBottom: 20 }}>
                    <Grid item>
                        <FormControl fullWidth variant="outlined">
                            <InputLabel htmlFor="outlined-adornment-amount">
                                {locale2.PRICE[props.lang]}
                            </InputLabel>
                            <OutlinedInput
                                type="number"
                                id="outlined-adornment-amount"
                                value={String(price)}
                                onChange={e => {
                                    let c = Number(e.target.value)
                                    if (!isNaN(c)) {
                                        setPrice(c)
                                    }
                                }}
                                startAdornment={
                                    <InputAdornment position="start">{props.currency}</InputAdornment>
                                }
                                labelWidth={60}
                            />
                        </FormControl>
                    </Grid>
                    <Grid item>
                        <TextField
                            id="standard-select-currency"
                            select
                            label={locale2.CURRENCY[props.lang]}
                            value={currency}
                            onChange={(e) => setCurrency(e.target.value)}>
                            {currencies.map((option) => (
                                <MenuItem key={option.value} value={option.value}>
                                    {option.label} / {option.value}
                                </MenuItem>
                            ))}
                        </TextField>
                    </Grid>
                </Grid>
            </React.Fragment>}
            onOpen={() => {
                setPrice(props.price / 100)
                setCurrency(props.currency)
            }}
        />
    )
}

export function getUntilDateLabel(d, x, lang) {
    let now = new Date()
    switch (d.getDay()) {
        case now.getDay():
            if (x === 1) {
                return locale2.TODAY[lang] + prettyPrintHr(d)
            }
            return locale2.TODAY_AT[lang] + prettyPrintHr(d)
        case now.getDay() + 1:
            if (x === 1) {
                return locale2.TOMORROW[lang] + prettyPrintHr(d)
            }
            return locale2.TOMORROW_AT[lang] + prettyPrintHr(d)
        default:
            return prettyPrintDate(d, lang)
    }
}

export function getRsvStatus(rsv, lang) {
    let s = rsv.State
    console.log(s)
    switch (s) {
        case "link_express":
        case "link":
        case "link_expire":
        case "retry_cancel_or_refund":
        case "wait_cancel_or_refund":
        case "wait_capture":
        case "wait_payout":
        case "error":
        case "wait_refund":
        case "payout":
            return locale2.SOMETHING_WENT_WRONG_CONTACT[lang]
        case "hold":
            return "zlożono rezerwacje, oczekiwanie na decyzję instruktora"
        case "capture":
            return "instruktor potwierdził rezerwację"
        case "dispute":
            return locale2.ISSUE_REPORTED[lang]
        case "refund":
        case "cancel_or_refund":
            return locale2.CANCELLED[lang]
    }
    // let state = rsv.State
    // let text = ""
    // let d = new Date(rsv.SmTimeout)
    // switch (state) {
    //     case "link_express":
    //     case "link":
    //         return locale2.WAITING_FOR_PAYMENT[lang]
    //     case "link_expire":
    //         return locale2.LINK_EXPIRED[lang]
    //     case "hold":
    //         return locale2.PAYMENT[lang] + getUntilDateLabel(d, 0, lang)
    //     case "retry_cancel_or_refund":
    //     case "wait_cancel_or_refund":
    //         return locale2.CANCELLATION[lang]
    //     case "cancel_or_refund":
    //         text = locale2.CANCELLED[lang]
    //         return text
    //     case "wait_capture":
    //         return locale2.FETCHING_PAYMENTS[lang]
    //     case "capture":
    //         text = locale2.AWAITING_PAYOUT[lang] + getUntilDateLabel(d, 0, lang)
    //         return text
    //     case "wait_payout":
    //         return locale2.PAYOUT_IN_PROGRESS[lang]
    //     case "error":
    //         return locale2.SOMETHING_WENT_WRONG_CONTACT[lang]
    //     case "wait_refund":
    //         return locale2.REFUNDING[lang]
    //     case "refund":
    //         return locale2.REFUNDED[lang]
    //     case "dispute":
    //         return locale2.ISSUE_REPORTED[lang]
    //     case "payout":
    //         return locale2.DONE_PAYOUT[lang]
    //     default:
    //         return state
    // }
}

export function stateDisplay(state, r, lang) {
    // let text = ""
    // let d = new Date(r.SmTimeout)
    // let now = new Date()
    return <Typography>
        {getRsvStatus(r, lang)}
    </Typography>
}

export function OccEditor(props) {

    const [repeatTraining, setRepeatTraining] = useState(false)
    const [repeatDays, setRepeatDays] = useState(0)
    const [dateStart, setDateStart] = useState(null)
    const [dateEnd, setDateEnd] = useState(null)

    function labelValue() {
        let occs = props.occs

        if (!occs || occs.length === 0)
            return (<Typography>Brak</Typography>)

        let leadingOcc = occs[0]

        return (<Typography>
            {prettyPrintDateRange(
                new Date(leadingOcc.DateStart), new Date(leadingOcc.DateEnd),
                leadingOcc.RepeatDays, null, props.lang)}
        </Typography>)
    }

    useEffect(() => {

        if (!props.occs || props.occs.length === 0)
            return null

        let leadingOcc = props.occs[0]

        setDateStart(new Date(leadingOcc.DateStart))
        setDateEnd(new Date(leadingOcc.DateEnd))
        if (leadingOcc.RepeatDays > 0) {
            setRepeatDays(leadingOcc.RepeatDays)
            setRepeatTraining(true)
        } else {
            setRepeatDays(0)
            setRepeatTraining(false)
        }

    }, [props.occs])

    return (<React.Fragment>
        <ModalEdit
            lang={props.lang}
            title={locale2.OCCURRENCE[props.lang]}
            value={labelValue()}
            custom
            onSave={async () => {
                if (!props.occs || props.occs.length === 0)
                    return null
                let leadingOcc = props.occs[0]
                let occs = {
                    TrainingID: leadingOcc.TrainingID,
                    Occurrences: [
                        {
                            DateStart: dateStart,
                            DateEnd: dateEnd,
                            RepeatDays: repeatDays
                        }
                    ]
                }
                await putOcc(occs)
                if (props.onChange)
                    props.onChange()
            }} content={<React.Fragment>
                <Grid container direction="column" spacing={2} style={{ padding: 10 }}>
                    <Grid item>
                        <DateTimePicker
                            value={dateStart}
                            onChange={e => {
                                if (e instanceof Date && isFinite(e))
                                    setDateStart(e)
                            }}
                            ampm={false}
                            renderInput={(params) => <TextField
                                size='small'
                                {...params} />}
                            label={locale2.START_DATE[props.lang]} />
                    </Grid>
                    <Grid item>
                        <DateTimePicker
                            value={dateEnd}
                            onChange={e => {
                                if (e instanceof Date && isFinite(e))
                                    setDateEnd(e)
                            }}
                            ampm={false}
                            renderInput={(params) => <TextField
                                size='small'
                                {...params} />}
                            label={locale2.END_DATE[props.lang]} />
                    </Grid>
                    <Grid item>
                        <FormControlLabel
                            control={<Checkbox size='small' checked={repeatTraining}
                                onChange={e => {
                                    let checked = e.target.checked
                                    if (!checked)
                                        setRepeatDays(0)
                                    setRepeatTraining(checked)
                                }} />}
                            label="Powtarzaj trening" />
                    </Grid>
                    {repeatTraining && (<React.Fragment>
                        <Grid item xs={12}>
                            <TextField label="powtarzaj trening co"
                                value={repeatDays}
                                type="number"
                                size='small'
                                onChange={e => setRepeatDays(Number(e.target.value))}
                                InputProps={{
                                    endAdornment: <InputAdornment position="end">dni</InputAdornment>
                                }} />
                        </Grid>
                    </React.Fragment>)}
                </Grid>
            </React.Fragment>}
            onOpen={() => { }}
        />
    </React.Fragment>)

}

export function instrDecisionDisplay(d, rsv, lang) {
    switch (d) {
        case "unset":
            if (rsv) {
                if (rsv.Training.ManualConfirm && rsv.IsActive) {
                    return locale2.AWAITING_DECISION[lang]
                        + Math.ceil((new Date(rsv.SmTimeout) - new Date()) / 1000 / 60 / 60)
                        + ' ' + locale2.HOURS_TO_MAKE_DECISION[lang]
                } else {
                    return locale2.NONE[lang]
                }
            } else {
                return locale2.NONE[lang]
            }
        case "approve":
            return locale2.ACCEPTED[lang]
        case "reject":
            return locale2.REJECTED[lang]
    }
}


export function prettyPrintRsvDecision(rsv, lang) {
    return instrDecisionDisplay(rsv.InstructorDecision, rsv, lang)
}

export function getRsvStatusColor(r) {
    return r.IsConfirmed ? MulwiColors.greenDark : (r.IsActive ? MulwiColors.blueLight : MulwiColors.redError)
}

export function isRsvCancelled(r) {
    return !r.IsConfirmed && !r.IsActive
}

export function RsvInfoListItem(props) {
    let r = props.rsv

    if (!r || !r.UserInfo) {
        return null
    }

    return (
        <Grid container style={{
            borderLeft: "solid",
            paddingLeft: 5,
            borderWidth: 6,
            borderColor: r.IsConfirmed ? MulwiColors.greenDark :
                (r.IsActive ? MulwiColors.blueLight : MulwiColors.redError)
        }}
            direction="row"
            justify="space-between"
            alignItems="center">
            <Grid item>
                <Grid container orientation="row" alignItems="center" justify="center">
                    {r.UserInfo.AvatarUrl && <Avatar src={r.UserInfo.AvatarUrl} />}
                    <div style={{ marginLeft: 10 }}>
                        {r.UserInfo.Name}
                    </div>
                </Grid>
            </Grid>
            <Grid item>
                <Grid container direction="column">
                    {r.InstructorDecision !== "unset" && <Grid item>
                        <span style={{ color: (r.InstructorDecision === "approve" && MulwiColors.greenDark) || MulwiColors.redError }}>
                            {instrDecisionDisplay(r.InstructorDecision, null, props.lang)}
                        </span>
                    </Grid>}
                    <Grid item>
                        {stateDisplay(r.State, r, props.lang)}
                    </Grid>
                </Grid>
            </Grid>
            <Grid item>
                <RsvInteractMenu
                    lang={props.lang}
                    instructor={props.instructor}
                    onChange={props.onChange}
                    setInfo={props.setInfo}
                    rsv={r} />
            </Grid>
        </Grid>)
}

export function getTagLabel(t) {
    if (!t || !t.Tag) return null
    return (t.Tag.Translations && t.Tag.Translations[getSupportedLanguage()]) || t.Tag.Name || ""
}
export function getCategoryLabel(t) {
    if (!t || !t.Category) return null
    return (t.Category.Translations && t.Category.Translations[getSupportedLanguage()]) || t.Category.Name
}


export function TrainingDetailsSideContent(props) {
    const [openMapDialog, setOpenMapDialog] = useState(false)
    const [newLocalization, setNewLocalization] = useState('')
    const [unsufficientData, setUnsufficientData] = useState(true)
    const [checkLocalization, setCheckLocalization] = useState(false)
    const [addrReceivedFromApi, setAddrReceivedFromApi] = useState({
        LocationText: '',
        LocationLat: '',
        LocationLng: '',
    })
    useEffect(() => {
        if (checkLocalization === false) {
            return
        }
        async function grabData(input) {
            try {
                let resp = await fetchGeo(input);
                if (JSON.parse(resp).length === 0) {
                    setUnsufficientData(true)
                } else {
                    setUnsufficientData(false)
                    let tmp = { ...addrReceivedFromApi }
                    tmp.LocationText = JSON.parse(resp)[0].display_name
                    tmp.LocationLat = JSON.parse(resp)[0].lat
                    tmp.LocationLng = JSON.parse(resp)[0].lon
                    setAddrReceivedFromApi(tmp)
                }
            } catch (e) {
                console.log(e)
            }
            setCheckLocalization(false)
        }
        grabData(newLocalization)
        // eslint-disable-next-line
    }, [checkLocalization])

    async function updateAll() {
        await updateComponent()
        props.onChange()
    }

    async function updateComponent() {
        let sch = props.drawerData.sch
        if (sch && sch.Occ) {
            let schRes = await getSchedule(
                sch.Start,
                sch.End,
                props.drawerData.training && props.drawerData.training.ID)
            for (let i = 0; i < schRes.length; i++) {
                let s2 = schRes[i]
                if (s2.Training.ID === props.drawerData.training.ID) {
                    for (let j = 0; j < s2.Schedule.length; j++) {
                        let sch = s2.Schedule[j]
                        /*
                        i use + before dates to convert them into milisecond timestamps
                        */
                        //if(+sch.Session.DateStart === +sess.Session.DateStart && +sch.Session.DateEnd === +sess.Session.DateEnd && sch.Session.ID === sess.Session.ID) {
                        if (sch.Occ.ID === sch.Occ.ID) {
                            let nd = { ...props.drawerData }
                            nd.training = s2.Training
                            nd.occs = s2.Occurrences
                            nd.groups = s2.Groups
                            nd.dc = s2.Dcs
                            nd.sch = sch
                            nd.openAddSession = false
                            props.setDrawerData(nd)
                            return
                        }
                    }
                }
            }
            props.setDrawerOpen(false)
            props.setDrawerData(null)
        } else {
            let ts = await getTrainings(props.drawerData.training.ID)
            if (!ts || ts.length === 0) {
                // training may not exist no more - close this drawer and return
                props.setDrawerOpen(false)
                props.setDrawerData(null)
                return
            }
            let t = ts[0]
            let nd = { ...props.drawerData }
            nd.training = t.Training
            nd.occs = t.Occurrences
            nd.groups = t.Groups
            nd.dc = t.Dcs
            nd.openAddSession = false
            nd.sch = null
            props.setDrawerData(nd)
        }
    }

    async function editTrainingWithFields(training, fields, values) {
        if (!fields || !values || !fields.length ||
            !values.length || values.length !== fields.length) {
            return
        }
        let cpy = { ...training }
        for (let i = 0; i < fields.length; i++) {
            cpy[fields[i]] = values[i]
        }
        try {
            await patchTraining(cpy)
            // notify interested parties that something changed about training
            props.onChange()
            // update urself
            await updateComponent()
        } catch (ex) {
            props.setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang],
                ex))
        }
    }

    const [options, setOptions] = useState([])

    async function updateTagOptions(i) {
        let c = await getTags(i)
        setOptions(JSON.parse(c))
    }

    const [tags, setTags] = useState([])

    async function refreshAtcTags() {
        if (!props.drawerData || !props.drawerData.training) return
        let t = props.drawerData.training.Tags;
        if (!t) return
        try {
            let et = JSON.parse(await explainTags(t))
            setTags(et || [])
        } catch (ex) {
            console.log(ex)
        }
    }

    useEffect(() => {
        refreshAtcTags()
    }, [props.drawerData])

    // you can prefill options with instructor tags 
    //     (if you manage to pass instructor to this component)
    // useEffect(() => {
    //     let c = []
    //     if(props.instructor && props.instructor.Tags) {
    //       let s = props.instructor.Tags
    //       if(s) {
    //         for(let i = 0; i < s.length; i++) 
    //           if(s[i]) c.push(s[i])
    //       }
    //     }
    //     setOptions(c)
    // }, [props.drawerData.training])

    function rsvListContent(rsvs) {
        if (!rsvs) return null
        return (<React.Fragment>
            {rsvs.map((r, i) => (
                <React.Fragment key={i}>
                    <ListItem>
                        <RsvInfoListItem
                            lang={props.lang}
                            instructor={true}
                            onChange={updateAll}
                            setInfo={props.setInfo}
                            rsv={r} />
                    </ListItem>
                </React.Fragment>
            ))}
        </React.Fragment>)
    }

    function rsvListGroup(title, rs) {
        if (!rs || rs.length === 0) return null
        return (<React.Fragment>
            <ListItem>
                <Typography variant="h6">{title}</Typography>
            </ListItem>
            {rsvListContent(rs)}
        </React.Fragment>)
    }

    function rsvListGroups(s, lang) {
        if (!s || !s.Reservations) {
            return (
                <React.Fragment>
                    <ListItem>
                        <Typography variant="h6">
                            {locale2.NO_RESERVATIONS[lang]}
                        </Typography>
                    </ListItem>
                </React.Fragment>
            )
        }
        let conf = []
        let act = []
        let dead = []
        for (let i = 0; i < s.Reservations.length; i++) {
            let x = s.Reservations[i]
            if (x.IsConfirmed) {
                conf.push(x)
                continue
            }
            if (x.IsActive) {
                act.push(x)
                continue
            }
            dead.push(x)
        }
        return (<React.Fragment>
            {rsvListGroup(locale2.CONFIRMED_RSVS[lang], conf)}
            {rsvListGroup(locale2.IN_PROGRESS_RSVS[lang], act)}
            {rsvListGroup(locale2.CANCELLED_RSVS[lang], dead)}
        </React.Fragment>)
    }

    function traininglistitems(lang) {
        return (<React.Fragment>
            <ListItem>
                <Typography variant="h6" >{locale2.OCCURRENCE[lang]}</Typography>
            </ListItem>
            <ListItem>
                <OccEditor
                    onChange={updateAll}
                    occs={props.drawerData.occs} lang={lang} />
            </ListItem>
            <ListItem>
                <Typography variant="h6" >{locale2.TRAINING_INFO[lang]}</Typography>
            </ListItem>
            <ListItem >
                <GenericInlineEdit
                    lang={props.lang}
                    label={locale2.NAME[lang]}
                    value={props.drawerData.training.Title}
                    onChange={async v =>
                        await editTrainingWithFields(props.drawerData.training, ["Title"], [v])} />
            </ListItem>
            <ListItem>
                <Grid container
                    style={{ marginRight: 15 }}
                    justifyContent="space-between"
                    alignItems="center"
                    direction="row">
                    <Grid item style={{ color: "gray" }}>
                        <Typography>{locale2.TRAINING_SUPPORTS_DISABLED[props.lang]}</Typography>
                    </Grid>
                    <Grid item>
                        <Checkbox
                            style={{ paddingRight: 0 }}
                            onChange={async (c, v) => {
                                try {
                                    await editTrainingWithFields(
                                        props.drawerData.training,
                                        ["TrainingSupportsDisabled"],
                                        [v])
                                } catch (ex) {
                                    props.setInfo(
                                        getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
                                }
                            }}
                            checked={props.drawerData.training.TrainingSupportsDisabled} />
                    </Grid>
                </Grid>
            </ListItem>
            <ListItem>
                <Grid container
                    style={{ marginRight: 15 }}
                    justifyContent="space-between"
                    alignItems="center"
                    direction="row">
                    <Grid item style={{ color: "gray" }}>
                        <Typography>{locale2.PLACE_SUPPORTS_DISABLED[props.lang]}</Typography>
                    </Grid>
                    <Grid item>
                        <Checkbox
                            style={{ paddingRight: 0 }}
                            onChange={async (c, v) => {
                                try {
                                    await editTrainingWithFields(
                                        props.drawerData.training,
                                        ["PlaceSupportsDisabled"],
                                        [v])
                                } catch (ex) {
                                    props.setInfo(
                                        getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
                                }
                            }}
                            checked={props.drawerData.training.PlaceSupportsDisabled} />
                    </Grid>
                </Grid>
            </ListItem>
            {/* <ListItem>
                <Grid container
                    style={{ marginRight: 15 }}
                    justify="space-between"
                    alignItems="center"
                    direction="row">
                    <Grid item style={{ color: "gray" }}>
                        <Typography>{locale2.CONFIRM_RSV_MANUALLY[props.lang]}</Typography>
                    </Grid>
                    <Grid item>
                        <Checkbox
                            style={{paddingRight: 0}}
                            disabled={props.drawerData.training.AllowExpress}
                            onChange={async (c, v) => {
                                try {
                                    await editTrainingWithFields(
                                        props.drawerData.training,
                                        ["ManualConfirm"],
                                        [v])
                                } catch (ex) {
                                    props.setInfo(
                                        getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
                                }
                            }}
                            checked={props.drawerData.training.ManualConfirm} />
                    </Grid>
                </Grid>
            </ListItem> */}
            <ListItem>
                <Grid container
                    style={{ marginRight: 15 }}
                    justifyContent="space-between"
                    alignItems="center"
                    direction="row">
                    <Grid item style={{ color: "gray" }}>
                        <Typography>{locale2.ALLOW_LAST_MINUTE[props.lang]}</Typography>
                    </Grid>
                    <Grid item>
                        <Checkbox
                            style={{ paddingRight: 0 }}
                            disabled={props.drawerData.training.ManualConfirm}
                            onChange={async (c, v) => {
                                try {
                                    await editTrainingWithFields(
                                        props.drawerData.training,
                                        ["AllowExpress"],
                                        [v])
                                } catch (ex) {
                                    props.setInfo(
                                        getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
                                }
                            }}
                            checked={props.drawerData.training.AllowExpress} />
                    </Grid>
                </Grid>
            </ListItem>
            <ListItem>
                <TextAreaModalEdit
                    label={locale2.DESCRIPTION[props.lang]}
                    multiline
                    lang={props.lang}
                    title={props.drawerData.training.Title + ' - ' + locale2.DESCRIPTION[props.lang]}
                    value={props.drawerData.training.Description}
                    onChange={async v => await editTrainingWithFields(props.drawerData.training, ["Description"], [v])}
                />
            </ListItem>
            <ListItem>
                <PriceEditor
                    lang={props.lang}
                    title={props.drawerData.training.Title}
                    price={props.drawerData.training.Price}
                    currency={props.drawerData.training.Currency}
                    onChange={async (p, c) => {
                        await editTrainingWithFields(
                            props.drawerData.training,
                            ["Price", "Currency"],
                            [p, c])
                    }}
                />
            </ListItem>
            <ListItem>
                <AgeEditor
                    lang={props.lang}
                    title={props.drawerData.training.Title}
                    minAge={props.drawerData.training.MinAge}
                    maxAge={props.drawerData.training.MaxAge}
                    onChange={async (p, c) => {
                        await editTrainingWithFields(
                            props.drawerData.training,
                            ["MinAge", "MaxAge"],
                            [p, c])
                    }}
                />
            </ListItem>
            <ListItem>
                <GenericInlineEdit
                    lang={props.lang}
                    label={locale2.MAX_AMOUNT_OF_PEOPLE[props.lang]}
                    value={props.drawerData.training.Capacity}
                    onChange={async v => await editTrainingWithFields(props.drawerData.training, ["Capacity"], [Number(v)])} />
                {/* <Typography>{props.drawerData.training.Description}</Typography> */}
            </ListItem>
            <ListItem>
                <TrainingImageEditor
                    // updateAll may be overkill  updating drawer may suffice.
                    lang={props.lang}
                    onChange={updateAll}
                    setInfo={props.setInfo}
                    training={props.drawerData.training} />
            </ListItem>
            <ListItem>
                <Grid container spacing={3}>
                    <Grid item xs={4}>
                        <Typography style={{
                            color: "gray"
                        }}>
                            {locale2.LOCATION[props.lang]}
                        </Typography>
                    </Grid>
                    <Grid item xs={6}>
                        <Typography noWrap>
                            {props.drawerData.training.LocationText}
                        </Typography>
                    </Grid>
                    <Grid item xs={2}>
                        <Button color="primary" size='small' aria-label="edit"
                            onClick={() => setOpenMapDialog(!openMapDialog)}>
                            <EditIcon />
                        </Button>
                    </Grid>
                </Grid>
                <Dialog open={openMapDialog} aria-labelledby="form-dialog-title"
                    onClose={() => { setOpenMapDialog(false) }} >
                    <DialogTitle id="form-dialog-title">
                        {locale2.ENTER_NEW_ADDRESS[props.lang]}
                    </DialogTitle>
                    <DialogContent>
                        <GmapEditor
                            withConfirmButton
                            lang={props.lang}
                            setLocalizationData={async (v) => {
                                if (!v) return
                                setOpenMapDialog(false)
                                await editTrainingWithFields(
                                    props.drawerData.training,
                                    ["LocationText", "LocationLat", "LocationLng"],
                                    [
                                        v.LocationText,
                                        parseFloat(v.LocationLat),
                                        parseFloat(v.LocationLng)
                                    ]
                                )
                            }} />
                    </DialogContent>
                </Dialog>
            </ListItem>
            <ListItem >
                <AutocompleteEdit
                    label={locale2.LEVEL_OF_DIFF[props.lang]}
                    placeholder={locale2.NOT_APPLICABLE[props.lang]}
                    options={getDiffsByLang(getSupportedLanguage())}
                    value={() => {
                        let x = props.drawerData.training.Diff
                        let y = []
                        if (x) {
                            for (let i = 0; i < x.length; i++) {
                                let v = x[i]
                                y.push({ id: v, label: diffs[v][getSupportedLanguage()] })
                            }
                        }
                        return y
                    }}
                    mapValue={(v) => ({ id: v, label: diffs[v][getSupportedLanguage()] })}
                    // convert props.value to string
                    valueLabel={v => v.label}
                    onChange={async (v) => {
                        let r = []
                        for (let i = 0; i < v.length; i++) {
                            let x = v[i]
                            if (r.indexOf(x) < 0)
                                r.push(x.id)
                        }
                        await editTrainingWithFields(
                            props.drawerData.training,
                            ["Diff"],
                            [r])
                    }} />
            </ListItem>
            <ListItem>
                <DynAtcEdit
                    lang={props.lang}
                    label={locale2.TAGS[props.lang]}
                    options={options}
                    value={() => {
                        return tags || []
                    }}
                    renderOption={(props, t) => <li {...props} >
                        <Grid container direction="row" justifyContent="space-between">
                            <Grid item><strong>{getTagLabel(t)}</strong></Grid>
                            <Grid item>{getCategoryLabel(t)}</Grid>
                        </Grid>
                    </li>}
                    valueLabel={getTagLabel}
                    updateOptions={updateTagOptions}
                    optionLabel={getTagLabel}
                    noFreeSolo
                    equals={(o, v) => o.Tag.Name === v.Tag.Name}
                    onChange={async (v) => {
                        let r = []
                        for (let i = 0; i < v.length; i++) {
                            let x = v[i]
                            if (r.indexOf(x) < 0)
                                r.push(x.Tag.Name)
                        }
                        await editTrainingWithFields(
                            props.drawerData.training,
                            ["Tags"],
                            [r])
                    }} />
            </ListItem>
            <ListItem>
                <DynAtcEdit
                    lang={props.lang}
                    placeholder={locale2.NOT_APPLICABLE[props.lang]}
                    label={locale2.REQUIRED_GEAR[props.lang]}
                    atclabel={locale2.GEAR_WHICH_CUSTOMER_MUST_HAVE[props.lang]}
                    options={[]}
                    value={() => {
                        return props.drawerData.training.RequiredGear || []
                    }}
                    mapValue={(v) => v}
                    valueLabel={v => v}
                    updateOptions={null}
                    onChange={async (v) => {
                        let r = []
                        for (let i = 0; i < v.length; i++) {
                            let x = v[i]
                            if (r.indexOf(x) < 0)
                                r.push(x)
                        }
                        await editTrainingWithFields(
                            props.drawerData.training,
                            ["RequiredGear"],
                            [r])
                    }} />
            </ListItem>
            <ListItem>
                <DynAtcEdit
                    lang={props.lang}
                    placeholder={locale2.NOT_APPLICABLE[props.lang]}
                    label={locale2.RECOMMENDED_GEAR[props.lang]}
                    atclabel={locale2.GEAR_WHICH_YOU_RECOMMEND_TO_CUSTOMER[props.lang]}
                    options={[]}
                    value={() => {
                        return props.drawerData.training.RecommendedGear || []
                    }}
                    mapValue={(v) => v}
                    valueLabel={v => v}
                    updateOptions={null}
                    onChange={async (v) => {
                        let r = []
                        for (let i = 0; i < v.length; i++) {
                            let x = v[i]
                            if (r.indexOf(x) < 0)
                                r.push(x)
                        }
                        await editTrainingWithFields(
                            props.drawerData.training,
                            ["RecommendedGear"],
                            [r])
                    }} />
            </ListItem>
            <ListItem>
                <DynAtcEdit
                    lang={props.lang}
                    placeholder={locale2.NOT_APPLICABLE[props.lang]}
                    label={locale2.YOUR_GEAR[props.lang]}
                    atclabel={locale2.GEAR_WHICH_YOU_HAVE[props.lang]}
                    options={[]}
                    value={() => {
                        return props.drawerData.training.InstructorGear || []
                    }}
                    mapValue={(v) => v}
                    valueLabel={v => v}
                    updateOptions={null}
                    onChange={async (v) => {
                        let r = []
                        for (let i = 0; i < v.length; i++) {
                            let x = v[i]
                            if (r.indexOf(x) < 0)
                                r.push(x)
                        }
                        await editTrainingWithFields(
                            props.drawerData.training,
                            ["InstructorGear"],
                            [r])
                    }} />
            </ListItem>
            <ListItem>
                {props.drawerData.groups && (<React.Fragment>
                    <GroupEdit lang={props.lang} onChange={updateAll} d={props.drawerData} />
                </React.Fragment>)}
            </ListItem>
            <ListItem>
                {props.drawerData.dc && (<React.Fragment>
                    <DcEdit lang={props.lang} dc={props.drawerData.dc} />
                </React.Fragment>)}
            </ListItem>
            <ListItem>
                {props.drawerData.sm && (<React.Fragment>
                    <SmEdit lang={props.lang} sm={props.drawerData.sm} />
                </React.Fragment>)}
            </ListItem>
            {/* <br />
            {props.drawerData.occs && props.drawerData.occs.length > 0 && (
                <React.Fragment>
                    <Grid item>
                        <Typography variant="h6" >Występowanie treningu</Typography>
                    </Grid>
                    <br />
                    {props.drawerData.occs.map((occ, i) => {
                        return (<React.Fragment key={i}>
                            <OccAcordionEditor
                                id={"occAccordion" + i}
                                next={"occAccordion" + (i + 1)}
                                onChange={updateAll}
                                drawerData={props.drawerData}
                                occ={occ}
                                setInfo={props.setInfo} />
                        </React.Fragment>)
                    })}
                    <div id={"occAccordion" + props.drawerData.occs.length}>
                    </div>
                </React.Fragment>
            )} */}
            {/* <ListItem>
                <OccDisplay occs={props.drawerData.occs} 
                    lang={props.lang}
                    t={props.drawerData.training} />
            </ListItem> */}
        </React.Fragment>)
    }

    function subList(subs) {
        if (!subs || subs.length === 0) return (
            <ListItem>
                <Typography variant="h6">{locale2.NO_CARNETS_AVAILABLE[props.lang]}</Typography>
            </ListItem>)
        return (<React.Fragment>
            <ListItem>
                <Typography variant="h6">{locale2.CARNETS_AVAILABLE_FOR_TRAINING[props.lang]}</Typography>
            </ListItem>
            {subs.map((s, i) => (<ListItem>
                <Grid container style={{
                    borderLeft: "solid",
                    paddingLeft: 5,
                    borderWidth: 6,
                    borderColor: s.IsConfirmed ? MulwiColors.greenDark :
                        (s.IsActive ? MulwiColors.blueLight : MulwiColors.redError)
                }}
                    direction="row"
                    justify="space-between"
                    alignItems="center">
                    {console.log(s)}
                    <Grid item>
                        <Grid container orientation="row" alignItems="center" justify="center">
                            {s.UserInfo.AvatarUrl && <Avatar src={s.UserInfo.AvatarUrl} />}
                            <div style={{ marginLeft: 10 }}>
                                {s.UserInfo.Name}
                            </div>
                        </Grid>
                    </Grid>
                    <Grid item>
                        <Grid container direction="column">
                            <Grid item>
                                {stateDisplay(s.State, s, props.lang)}
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item>
                        <SubInteractMenu
                            instructor={props.instructor}
                            onChange={props.onChange}
                            setInfo={props.setInfo}
                            sub={s} />
                    </Grid>
                </Grid>
            </ListItem>))}
        </React.Fragment>)
    }

    return (props.drawerData && props.drawerData.training && (
        <React.Fragment>
            <List style={{ backgroundColor: MulwiColors.whiteBackground, padding: 5 }}>
                {rsvListGroups(props.drawerData && props.drawerData.sch, props.lang)}
                {subList(props.drawerData.sch && props.drawerData.sch.Subs)}
                {(props.drawerData.sch && props.drawerData.sch.IsOrphaned && (
                    <React.Fragment>
                        <ListItem>
                            <Typography color="secondary">
                                {locale2.TRAINING_GONE[props.lang]}
                            </Typography>
                        </ListItem>
                        <TrainingSummary lang={props.lang}
                            setInfo={props.setInfo}
                            onChange={props.onChange}
                            session={props.drawerData && props.drawerData.sch}
                            training={props.drawerData && props.drawerData.training} />
                    </React.Fragment>
                )) || (
                        traininglistitems(props.lang)
                    )}
            </List>
        </React.Fragment>)) || null
}
