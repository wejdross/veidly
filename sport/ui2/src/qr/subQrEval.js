import { Button, Card, CardContent, CardHeader, 
        CircularProgress, Grid, Typography } from '@mui/material'
import { Close } from '@mui/icons-material'
import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { confirmQr, evalQr } from '../apicalls/sm'
import { MulwiColors } from '../mulwiColors'
import { errToStr } from '../StatusDialog'
import { KeyVal } from '../sub/SubCard'
import { SubInfo } from '../sub/subDetails'
import { locale2 } from '../locale'

export function SubQrEval(props) {

    const [st, setst] = useState(null)

    async function refresh() {
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
    }

    useEffect(() => {
        refresh()
    }, [])

    const [IsConfirmed, setIsConfirmed] = useState(false)

    async function confirmEntry() {
        let query = new URLSearchParams(window.location.search)
        let id = query.get("id")
        try {
            await confirmQr(id)
            setIsConfirmed(true)
        } catch (ex) {
            if (ex == 409) {
                setst({
                    ex: ex,
                    msg: locale2.CARNET_ISNT_VALID[props.lang]
                })
            } else {
                setst({
                    ex: ex
                })
            }
        }
    }

    return (
        <React.Fragment>
            <Grid container style={{
                width: "100%",
                height: "80vh",

            }} justify="center" alignItems="center">
                <Grid item>
                    <Card>
                        <CardHeader title={<center>
                            {locale2.VER_OF_SUB_QR_CODE[props.lang]}
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
                                        {locale2.FAILED_TO_CONFIRM_CARNET[props.lang]}
                                    </Typography>
                                    <Typography variant="body2" style={{
                                        maxWidth: 400
                                    }}>
                                        {st.msg}
                                    </Typography>
                                    <br />
                                    {errToStr(st.ex)}
                                    <Button onClick={refresh}>
                                        {locale2.RETURN[props.lang]}
                                    </Button>
                                </center></CardContent>
                        )}
                        {st && st.res && (
                            <center>
                                <CardContent>
                                    {st.res && (
                                    <React.Fragment>
                                        <SubInfo lang={props.lang} sub={st.res} instr />
                                        <KeyVal 
                                            k={locale2.CARNET_ISSUED_FOR[props.lang]} 
                                            v={st.res.UserInfo.Name || locale2.ANON[props.lang]} />
                                        <Button style={{ 
                                                marginTop: 20, color: MulwiColors.blueDark }} fullWidth>
                                            <Link style={{
                                                color: MulwiColors.blueDark,
                                                textDecoration: "none",
                                            }} to={"/sub_details?instr=1&id=" + st.res.ID}>
                                                {locale2.GOTO_CARNET[props.lang]}
                                            </Link>
                                        </Button>
                                        {!IsConfirmed && <Button variant="contained" fullWidth style={{
                                            color: "white",
                                            marginTop: 20,
                                            backgroundColor: MulwiColors.greenDark
                                        }} onClick={confirmEntry}>
                                            {locale2.CONFIRM_ENTRANCE[props.lang]}
                                        </Button>}
                                    </React.Fragment>)}
                                </CardContent>
                            </center>
                        )}
                    </Card>
                </Grid>
            </Grid>
        </React.Fragment>
    )
}