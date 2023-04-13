import { Container, Grid, 
        InputAdornment, TextField, 
        Typography } from '@mui/material';
import React from 'react';
import { getSm, patchSm, postSm, 
        deleteSm, getTrainingsForSm, 
        postSmBinding, deleteSmBinding } from '../apicalls/sm';
import { TextInput } from '../CRUDCard/AddModal';
import { Card } from '../CRUDCard/Card';
import { daySuffix as dl } from '../harmonogram/trainingDetails';
import { prettyPrintCurrency } from '../helpers';
import { locale2 } from '../locale';
import { Subs } from './Subs';

function editor(lang) {
    return function(req, setReq) {
        return (<Grid container spacing={2} direction="column">
            <Grid item>
                <TextInput
                    value={req.Name || ""}
                    label={locale2.NAME[lang]}
                    onChange={e => setReq(r => ({ ...r, Name: e.target.value }))} />
            </Grid>
            <Grid item>
                <Typography>
                    {locale2.MAX_NUMBER_OF_ENTRIES[lang]}
                </Typography>
                <TextField
                    variant="outlined"
                    size="small"
                    fullWidth
                    placeholder={locale2.NOT_APPLICABLE[lang]}
                    value={req.MaxEntrances === -1 ? "" : String(req.MaxEntrances || "")}
                    onChange={e => {
                        let x = Number(e.target.value)
                        if (isNaN(x)) return
                        if (x === 0) x = -1
                        setReq(r => ({ ...r, MaxEntrances: x }))
                    }} />
            </Grid>
            <Grid item>
                <Typography>
                    {locale2.PRICE[lang]}
                </Typography>
                <TextField
                    InputProps={{
                        endAdornment: <InputAdornment position="end">
                            {prettyPrintCurrency(req.Currency)}
                        </InputAdornment>
                    }}
                    variant="outlined"
                    size="small"
                    fullWidth
                    type="number"
                    value={
                        (req.Price === 0 ) ? "whatever-what-is-not-empty-string" : req.Price / 100
                    }
                    onChange={e => {
                        var finalPrice
                        if (e.target.value === "" || isNaN(e.target.value)) {
                            finalPrice = 0
                        }
                        finalPrice = e.target.value
                        setReq(r => ({ ...r, Price: Math.round(finalPrice * 100) }))
                    }} />
            </Grid>
            <Grid item>
                <Typography>
                    {locale2.DURATION[lang]}
                </Typography>
                <TextField
                    InputProps={{
                        endAdornment: <InputAdornment position="end" style={{
                            width: 30
                        }}>
                            {dl(req.Duration || 30, lang)}
                        </InputAdornment>
                    }}
                    variant="outlined"
                    size="small"
                    fullWidth
                    value={String(req.Duration || 30)}
                    onChange={e => {
                        let x = Number(e.target.value)
                        if (isNaN(x)) return
                        setReq(r => ({ ...r, Duration: x }))
                    }} />
            </Grid>
        </Grid>)
    }
}

export function SmCard(props) {

    let lang = props.lang

    return (<React.Fragment>

        <Card lang={props.lang}
            getData={async function () {
                let d = await getSm()
                d = JSON.parse(d)
                return d
            }}

            newReq={function () {
                return {
                    Name: "",
                    MaxEntrances: -1,
                    Currency: "PLN",
                    Price: 5000,
                    Duration: 30
                }
            }}

            postData={async r => {
                let d = await postSm(r)
                return JSON.parse(d)
            }}
            patchData={async r => await patchSm(r)}
            deleteData={async r => await deleteSm(r.ID)}

            cardHeader={<React.Fragment>
                <Typography variant="h5" style={{
                    marginBottom: 5
                }}>
                    {locale2.CARNETS[props.lang]}
                </Typography>
                <Typography style={{
                    maxWidth: 800
                }}>
                    {locale2.CARNETS_ALLOW_FOR_CL[props.lang]}
                </Typography>
                <Typography style={{
                    maxWidth: 750
                }} variant="body2">
                    {locale2.CARNETS_INSTEAD_OF_SIGNING_IN[props.lang]}
                </Typography>

                <Typography style={{
                    maxWidth: 800
                }} variant="body2">
                    <strong>{locale2.VERIFY_CARNET_LIKE_RSV[props.lang]}</strong>
                </Typography>
            </React.Fragment>}

            tableColumns={[
                {
                    header: locale2.NAME[props.lang],
                    fieldSelector: c => c.Name
                }, {
                    header: locale2.NUMBER_OF_ENTRIES[props.lang],
                    fieldSelector: e => e.MaxEntrances === -1 ? 
                     locale2.NOT_APPLICABLE[props.lang] : e.MaxEntrances
                }, {
                    header: locale2.DURATION[props.lang],
                    fieldSelector: e => e.Duration + " " + dl(e.Duration, lang)
                }, {
                    header: locale2.PRICE[props.lang],
                    fieldSelector: e => e.Price / 100 + " " + prettyPrintCurrency(e.Currency)
                },
            ]}

            objectName={locale2.SUB_NAME[lang]}
            updateForm={editor(props.lang)}
            nameSelector={r => r.Name}

            getTrainings={async r => await getTrainingsForSm(r.ID)}
            createBinding={async (r, tid) => await postSmBinding(r.ID, tid)}
            deleteBinding={async (r, tid) => await deleteSmBinding(r.ID, tid)}
        />

        <br />
        <Container>
            <Subs lang={props.lang} instr user={props.user} embedded />
        </Container>
    </React.Fragment>)
}