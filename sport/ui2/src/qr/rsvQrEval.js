import {
    Button, Card, CardContent, CardHeader,
    CircularProgress,
    Grid, Typography
} from '@mui/material'
import { Check, Close } from '@mui/icons-material'
import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { evalQr } from '../apicalls/instructor.api'
import { MulwiColors } from '../mulwiColors'
import { errToStr } from '../StatusDialog'
import { locale2 } from '../locale'

export function RsvQrEval(props) {

    const [st, setst] = useState(null)

    useEffect(() => {
        (async function () {

            let query = new URLSearchParams(window.location.search)

            let id = query.get("id")
            if (!id) {
                setst({
                    ex: 400,
                })
                return
            }

            try {
                let res = await evalQr(id)
                res = JSON.parse(res)
                setst({
                    res: res
                })
            } catch (ex) {
                if (ex == 404) {
                    setst({
                        ex: ex,
                        msg: locale2.TO_CONFIRM_YOU_MUST_BE_LOGGED_IN[props.lang]
                    })
                } else {
                    setst({
                        ex: ex
                    })
                }
            }
        })()
    }, [])

    return (
        <React.Fragment>
            <Grid container style={{
                width: "100%",
                height: "80vh",

            }} justify="center" alignItems="center">
                <Grid item>
                    <Card>
                        <CardHeader title={<center>
                            {locale2.VER_OF_RSV_QR_CODE[props.lang]}
                        </center>} />
                        {!st && (
                            <center>
                                <CardContent>
                                    <CircularProgress style={{
                                        width: 70,
                                        height: 70,
                                        color: MulwiColors.blueDark
                                    }} />
                                </CardContent>
                            </center>
                        )}
                        {st && st.ex && (
                            <CardContent>
                                <center>
                                    <Close style={{
                                        width: 70,
                                        height: 70,
                                        color: MulwiColors.redError
                                    }} />
                                    <Typography style={{ color: MulwiColors.redError }}>
                                        {locale2.CANT_CONFIRM_RSV[props.lang]}
                                    </Typography>
                                    <Typography variant="body2" style={{
                                        maxWidth: 400
                                    }}>
                                        {st.msg}
                                    </Typography>
                                    <br />
                                    {errToStr(st.ex)}
                                </center></CardContent>
                        )}
                        {st && st.res && (
                            <center>
                                <CardContent>
                                    {(st.res.ConfirmCode === 0 && (<React.Fragment>
                                        <Check style={{
                                            width: 70,
                                            height: 70,
                                            color: MulwiColors.greenDark
                                        }} />
                                        <Typography>
                                            {locale2.RSV_CONFIRMED[props.lang]}
                                        </Typography>
                                    </React.Fragment>)) || (<React.Fragment>
                                        <Close style={{
                                            width: 70,
                                            height: 70,
                                            color: MulwiColors.redError
                                        }} />
                                        {((st.res.ConfirmCode & 1) && (<Typography>
                                            {locale2.QR_HAS_BEEN_USED[props.lang]}
                                        </Typography>)) || null}
                                        {((st.res.ConfirmCode & 2) && (<Typography>
                                            {locale2.RSV_NOT_PAYED[props.lang]}
                                        </Typography>)) || null}
                                    </React.Fragment>
                                        )}
                                    {st.res && (<Button style={{ color: MulwiColors.blueDark }} fullWidth>
                                        <Link style={{
                                            color: MulwiColors.blueDark,
                                            textDecoration: "none"
                                        }} to={"/rsv_details?instr=1&id=" + st.res.RsvID}>
                                            {locale2.GOTO_RSV[props.lang]}
                                        </Link>
                                    </Button>)}
                                </CardContent>
                            </center>
                        )}
                    </Card>
                </Grid>
            </Grid>
        </React.Fragment>
    )
}