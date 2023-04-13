import { Button, Grid, TextField, Typography } from '@mui/material';
import { DatePicker as KeyboardDatePicker } from '@mui/x-date-pickers/DatePicker';
import React from 'react';
import {
    deleteDc, deleteDcBinding, getDc, getTrainingsForDc,
    patchDc, postDc, postDcBinding
} from '../apicalls/dc';
import { Card } from '../CRUDCard/Card';
import { prettyPrintDate } from '../harmonogram/trainingDetails';
import { randomString } from '../helpers';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';

export function DcCard(props) {

    let now = new Date()

    let lang = props.lang

    return (<React.Fragment>

        <Card
            lang={props.lang}

            getData={async function () {
                let d = await getDc()
                return JSON.parse(d)
            }}

            newReq={function () {
                let vstart = new Date()
                let vend = new Date(vstart)
                vend.setMonth(vend.getMonth() + 1)
                return {
                    Name: "",
                    Discount: 15,
                    Quantity: 1,
                    ValidStart: vstart,
                    ValidEnd: vend
                }
            }}

            postData={async req => JSON.parse(await postDc(req))}
            patchData={async req => await patchDc(req)}
            deleteData={async req => await deleteDc(req.ID)}

            cardHeader={<React.Fragment>
                <Typography variant="h5" style={{
                    marginBottom: 5
                }}>
                    {locale2.DCS[props.lang]}
                </Typography>
            </React.Fragment>}

            tableColumns={[
                {
                    header: "",
                    fieldSelector: e => (
                        (new Date(e.ValidStart) <= now && now <= new Date(e.ValidEnd)) && ((e.Quantity - e.RedeemedQuantity) > 0) && (
                            <div style={{ borderRadius: 10, width: 15, height: 15, backgroundColor: MulwiColors.greenDark }}></div>
                        ) || (
                            <div style={{ borderRadius: 10, width: 15, height: 15, backgroundColor: MulwiColors.redError }}></div>
                        )
                    )
                }, {
                    header: locale2.CODE[props.lang],
                    fieldSelector: e => <strong>{e.Name}</strong>
                }, {
                    header: locale2.DISCOUNT[props.lang],
                    fieldSelector: e => e.Discount + "%"
                }, {
                    header: locale2.NO_USES[props.lang],
                    fieldSelector: e => e.RedeemedQuantity
                }, {
                    header: locale2.REMAINING[props.lang],
                    fieldSelector: e => e.Quantity - e.RedeemedQuantity < 0 ? 0 : e.Quantity - e.RedeemedQuantity
                }, {
                    header: locale2.VALID_FROM[props.lang],
                    fieldSelector: e => prettyPrintDate(new Date(e.ValidStart), props.lang)
                }, {
                    header: locale2.VALID_TO[props.lang],
                    fieldSelector: e => prettyPrintDate(new Date(e.ValidEnd), props.lang)
                }
            ]}

            objectName={locale2.DC_NAME[lang]}
            updateForm={(req, setReq) => (<React.Fragment>
                <Grid container spacing={2} direction="column">
                    <Grid item>
                        <Typography variant="body2">
                            {locale2.USE_ANY_NAME_OR[props.lang]}
                        </Typography>
                        <Button onClick={() => setReq(c => ({...c, Name: randomString(12)}))}>
                            {locale2.AUTOGEN[props.lang]}
                        </Button>
                        <Typography variant="body2">
                            {locale2.CODE_ON_BOOKING[props.lang]}
                        </Typography>
                        <TextField
                            variant="outlined"
                            size="small"
                            fullWidth
                            value={req.Name}
                            onChange={e => setReq(c => ({...c, Name: e.target.value}))}
                            label={locale2.CODE[props.lang]}
                        />
                    </Grid>
                    <Grid item>
                        <Grid container spacing={2} direction="row">
                            <Grid item sm={6}>
                                <Typography>
                                    {locale2.DISCOUNT[props.lang]}
                                </Typography>
                                <TextField
                                    variant="outlined"
                                    size="small"
                                    type="number"
                                    fullWidth
                                    value={String(req.Discount)}
                                    onChange={e => {
                                        let x = Number(e.target.value)
                                        if (isNaN(x)) return
                                        setReq(c => ({...c, Discount: x}))
                                    }}
                                    InputProps={{
                                        endAdornment: (<Typography>
                                            %
                                        </Typography>)
                                    }}
                                    />
                            </Grid>
                            <Grid item sm={6}>
                                <Typography>
                                    {locale2.NO_USES[props.lang]}
                                </Typography>
                                <TextField
                                    variant="outlined"
                                    size="small"
                                    type="number"
                                    fullWidth
                                    value={String(req.Quantity)}
                                    onChange={e => {
                                        let x = Number(e.target.value)
                                        if (isNaN(x)) return
                                        setReq(c => ({...c, Quantity: x}))
                                    }}
                                    />
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item>
                        <Grid container spacing={2} direction="row">
                            <Grid item xs={12} sm={6}>
                                <KeyboardDatePicker fullWidth
                                    margin="normal"
                                    id="date-picker-dialog"
                                    label={locale2.VALID_FROM[props.lang]}
                                    format="dd/MM/yyyy"
                                    value={req.ValidStart}
                                    minDate={new Date()}
                                    onChange={e => {
                                        let x = new Date(e)
                                        x.setHours(0, 0, 0, 0)
                                        setReq(c => ({...c, ValidStart: x}))
                                    }}
                                    KeyboardButtonProps={{
                                        'aria-label': 'change date',
                                    }} />
                            </Grid>
                            <Grid item sm={12} md={6}>
                                <KeyboardDatePicker fullWidth
                                    margin="normal"
                                    id="date-picker-dialog"
                                    label={locale2.VALID_TO[props.lang]}
                                    format="dd/MM/yyyy"
                                    value={req.ValidEnd}
                                    minDate={req.ValidStart || new Date()}
                                    onChange={e => {
                                        let x = new Date(e)
                                        x.setHours(23, 59, 59, 0)
                                        setReq(c => ({...c, ValidEnd: x}))
                                    }}
                                    KeyboardButtonProps={{
                                        'aria-label': 'change date',
                                    }} />
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </React.Fragment>)}
            nameSelector={r => r.Name}

            getTrainings={async r => await getTrainingsForDc(r.ID)}
            createBinding={async (r, tid) => await postDcBinding(r.ID, tid)}
            deleteBinding={async (r, tid) => await deleteDcBinding(r.ID, tid)}
        />
    </React.Fragment>)
}