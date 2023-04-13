import { Button, Fab, Grid, 
        Menu, MenuItem, TextField, Typography } from '@mui/material';
import { MoreHoriz } from '@mui/icons-material';
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import {
    createQr,
    postDecision, postInstructorRefund,
    postRsvExpire,
    postUserCancel, postUserDispute, postUserRefund
} from '../apicalls/instructor.api';
import { MulwiColors } from '../mulwiColors';
import { getDialogWithOptions, getErrorDialog, getNullDialog, supportEmail } from '../StatusDialog';
import { getUntilDateLabel, instrDecisionDisplay } from './trainingDetails';
import { locale2 } from '../locale';
import { sprintf } from '../helpers';
import { makeStyles } from "@mui/styles";

const linkStyles = makeStyles(t => ({
    link: {
        textDecoration: "none",
        color: MulwiColors.blueDark
    }
}))

export function DisputeForm(props) {

    const [email, setEmail] = useState("")
    const [msg, setMsg] = useState("")

    return (<React.Fragment>
        <TextField
            style={{ marginBottom: 10 }}
            fullWidth
            variant="outlined"
            size="small"
            value={email}
            onChange={(e) => {
                setEmail(e.target.value)
            }}
            label={locale2.YOUR_EMAIL[props.lang]} />
        <TextField
            inputProps={
                { maxLength: 250 }
            }
            helperText={((msg && msg.length) || 0) + "/250"}
            style={{ marginBottom: 10 }}
            fullWidth
            value={msg}
            onChange={(e) => setMsg(e.target.value)}
            label={locale2.DESCRIBE_YOUR_ISSUE[props.lang]}
            rows={10}
            multiline
            variant="outlined"
            size="small" />
        <br />
        <Typography>
            {locale2.UNTIL_PROBLEM_SOLVED[props.lang]}
        </Typography>
        <Button fullWidth
            variant="contained"
            style={{
                backgroundColor: MulwiColors.blueDark,
                color: "white"
            }}
            onClick={() => props.onChange(email, msg)} >
            {locale2.REPORT_ISSUE[props.lang]}
        </Button>
    </React.Fragment>)
}

export function RsvInteractMenu(props) {

    const [menuOpen, setMenuOpen] = useState(null)

    const c = linkStyles()

    let r = props.rsv


    async function postInstructorDecision(d) {
        try {
            await postDecision(r.ID, d)
            if (props.onChange) await props.onChange()
            props.setInfo(getNullDialog())
            setMenuOpen(false)
        } catch (ex) {
            switch (d) {
                case "approve":
                    props.setInfo(getErrorDialog(
                        locale2.SOMETHING_WENT_WRONG[props.lang],
                        ex,
                        <Button onClick={() => props.setInfo(getAcceptDialogOptions())} style={{
                            color: MulwiColors.blueLight
                        }}>
                            {locale2.ONCE_AGAIN[props.lang]}
                        </Button>))
                    break
                case "reject":
                    props.setInfo(getErrorDialog(
                        locale2.SOMETHING_WENT_WRONG[props.lang],
                        ex,
                        <Button onClick={() => props.setInfo(getRejectDialogOptions())} style={{
                            color: MulwiColors.blueLight
                        }}>
                            {locale2.ONCE_AGAIN[props.lang]}
                        </Button>))
                    break
                default:
                    props.setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
                    break
            }
        }
    }

    async function doInstructorRefund() {
        try {
            await postInstructorRefund(r.ID)
            if (props.onChange) await props.onChange()
            props.setInfo(getNullDialog())
            setMenuOpen(false)
        } catch (ex) {
            props.setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang],
                ex,
                <Button onClick={() => props.setInfo(getInstructorRefundDialogOptions())} style={{
                    color: MulwiColors.blueLight
                }}>
                    {locale2.ONCE_AGAIN[props.lang]}
                </Button>))
        }
    }

    async function generateQr() {
        try {
            let qr = await createQr(r.ID, props.at)
            // let downloadUrl = URL.createObjectURL(qr);
            let a = document.createElement("a")
            document.body.appendChild(a)
            a.style = "display: none"
            a.href = qr
            a.download = "qr.png"
            a.click()
            a.remove()
            props.setInfo(getNullDialog())
            setMenuOpen(false)
        } catch (ex) {
            console.log(ex)
            props.setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang],
                ex,
                <Button onClick={() => props.setInfo(getGenerateQrDialogOptions())} style={{
                    color: MulwiColors.blueLight
                }}>
                    {locale2.ONCE_AGAIN[props.lang]}
                </Button>))
        }
    }

    async function doUserRefund() {
        try {
            await postUserRefund(r.ID, props.at)
            if (props.onChange) await props.onChange()
            props.setInfo(getNullDialog())
            setMenuOpen(false)
        } catch (ex) {
            props.setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang],
                ex,
                <Button onClick={() => props.setInfo(getUserRefundDialogOptions())} style={{
                    color: MulwiColors.blueLight
                }}>
                    {locale2.ONCE_AGAIN[props.lang]}
                </Button>))
        }
    }

    async function doUserDispute(email, msg) {
        try {
            await postUserDispute(r.ID, props.at, email, msg)
            if (props.onChange) await props.onChange()
            props.setInfo(getNullDialog())
            setMenuOpen(false)
        } catch (ex) {
            props.setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang],
                ex,
                <Button onClick={() => props.setInfo(getDisputeDialogOptions())} style={{
                    color: MulwiColors.blueLight
                }}>
                    {locale2.ONCE_AGAIN[props.lang]}
                </Button>))
        }
    }

    async function doLinkExpire() {
        try {
            await postRsvExpire(r.ID, props.at)
            if (props.onChange) await props.onChange()
            props.setInfo(getNullDialog())
            setMenuOpen(false)
        } catch (ex) {
            props.setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang],
                ex,
                <Button
                    onClick={() => props.setInfo(getLinkExpDialogOptions())}
                    style={{
                        color: MulwiColors.blueLight
                    }}>
                    {locale2.ONCE_AGAIN[props.lang]}
                </Button>))
        }
    }

    async function doUserCancel() {
        try {
            await postUserCancel(r.ID, props.at)
            if (props.onChange) await props.onChange()
            props.setInfo(getNullDialog())
            setMenuOpen(false)
        } catch (ex) {
            props.setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang],
                ex,
                <Button onClick={() => props.setInfo(getUserCancelDialogOptions())} style={{
                    color: MulwiColors.blueLight
                }}>
                    {locale2.ONCE_AGAIN[props.lang]}
                </Button>))
        }
    }

    function getAcceptDialogOptions() {
        if (!r.ManualConfirm) {
            return (getDialogWithOptions(
                locale2.ACCEPTING_RSV[props.lang],
                "Rezerwacja opłacona - potwierdź",
                <Button variant="contained" style={{
                    backgroundColor: MulwiColors.greenDark,
                    color: "white"
                }} onClick={() => postInstructorDecision("approve")}>
                    {locale2.ACCEPT_USER[props.lang]} {r.UserInfo.Name}
                </Button>))
        } else {
            return (getDialogWithOptions(
                locale2.ARE_YOU_SURE[props.lang],
                locale2.ACCEPTING_RSV_IS_NOT_NECESSARY[props.lang],
                <Button variant="contained" style={{
                    backgroundColor: MulwiColors.greenDark,
                    color: "white"
                }} onClick={() => postInstructorDecision("approve")}>
                    {locale2.ACCEPT_USER[props.lang]} {r.UserInfo.Name}
                </Button>))
        }
    }

    function getRejectDialogOptions() {
        return (getDialogWithOptions(
            locale2.ARE_YOU_SURE[props.lang],
            "Anuluj tą rezerwacje",
            <Button variant="contained" style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }} onClick={() => postInstructorDecision("reject")}>
                {locale2.REJECT_USER[props.lang]} {r.UserInfo.Name}
            </Button>))
    }

    function getInstructorRefundDialogOptions() {
        return (getDialogWithOptions(
            locale2.ARE_YOU_SURE[props.lang],
            locale2.DECUCT_WARNING[props.lang],
            <Button variant="contained" style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }} onClick={() => doInstructorRefund()}>
                {locale2.YES_CANCEL_RSV[props.lang]}
            </Button>))
    }

    function getGenerateQrDialogOptions() {
        return (getDialogWithOptions(
            locale2.CREATE_QR[props.lang],
            locale2.CREATE_QR_EXTENDED[props.lang],
            <Button variant="contained" style={{
                backgroundColor: MulwiColors.greenDark,
                color: "white"
            }} onClick={() => generateQr()}>
                Generuj kod
            </Button>))
    }

    function getUserCancelDialogOptions() {
        return (getDialogWithOptions(
            locale2.ARE_YOU_SURE[props.lang],
            locale2.SOMETHING_WENT_WRONG[props.lang],
            <Button variant="contained" style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }} onClick={() => doUserCancel()}>
                {locale2.CANCEL_RSV[props.lang]}
            </Button>))
    }

    function getLinkExpDialogOptions() {
        return (getDialogWithOptions(
            locale2.ARE_YOU_SURE[props.lang],
            locale2.ARE_YOU_SURE[props.lang],
            <Button variant="contained" style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }} onClick={() => doLinkExpire()}>
                {locale2.YES[props.lang]}
            </Button>))
    }

    function getUserRefundDialogOptions() {
        return (getDialogWithOptions(
            locale2.ARE_YOU_SURE[props.lang],
            locale2.ASK_FOR_REFUND[props.lang] + ' ' + supportEmail,
            <Button variant="contained" style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }} onClick={() => doUserRefund()}>
                {locale2.YES[props.lang]}
            </Button>))
    }

    function getDisputeDialogOptions() {
        return (getDialogWithOptions(
            locale2.HOW_CAN_WE_HELP[props.lang],
            <DisputeForm lang={props.lang} onChange={(e, m) => doUserDispute(e, m)} />,
            null
        ))
    }

    async function selectInstrDecision(d) {
        switch (d) {
            case "approve":
                props.setInfo(getAcceptDialogOptions())
                break
            case "reject":
                props.setInfo(getRejectDialogOptions())
                break
        }
    }

    if (!r) {
        return null
    }

    function menuDetails() {
        return (!props.noDetails && (
            <MenuItem><Button style={{
                color: MulwiColors.blueDark
            }} fullWidth>
                <Link className={c.link}
                    to={"/rsv_details?id=" + r.ID + ((props.instructor && "&instr=1") || "")}>
                    {locale2.DETAILS[props.lang]}
                </Link>
            </Button>
            </MenuItem>
        )) || null
    }

    //

    function acceptBtn() {
        return (<Button
            variant="contained"
            fullWidth
            onClick={() => selectInstrDecision("approve")}
            style={{
                backgroundColor: MulwiColors.greenDark,
                color: "white"
            }}>
            {locale2.ACCEPT[props.lang]}
        </Button>)
    }

    function rejectBtn() {
        return (<Button
            fullWidth
            variant="contained"
            onClick={() => selectInstrDecision("reject")}
            style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }}>
            {locale2.CANCEL[props.lang]}
        </Button>)
    }

    function instrCancelBtn() {
        return (<Button
            fullWidth
            variant="contained"
            onClick={() => props.setInfo(getInstructorRefundDialogOptions())}
            style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }}>{locale2.INSTRUCTOR_CANCEL_RSV[props.lang]}</Button>)
    }

    function userCancelBtn() {
        return (<Button
            fullWidth
            variant="contained"
            onClick={() => props.setInfo(getUserCancelDialogOptions())}
            style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }}>{locale2.CANCEL_RSV[props.lang]}</Button>)
    }

    function disputeBtn() {
        return (<Button
            fullWidth
            variant="contained"
            style={{
                backgroundColor: MulwiColors.blueDark,
                color: "white"
            }}
            onClick={() => props.setInfo(getDisputeDialogOptions())}>
            {locale2.REPORT_RSV[props.lang]}
        </Button>)
    }

    function userRefund() {
        return (<Button
            fullWidth
            variant="contained"
            onClick={() => props.setInfo(getUserRefundDialogOptions())}
            style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }}>{locale2.INSTANT_REFUND[props.lang]}</Button>)
    }


    // function linkExpBtn() {
    //     return (<Button
    //         fullWidth
    //         variant="contained"
    //         onClick={() => props.setInfo(getLinkExpDialogOptions())}
    //         style={{
    //             backgroundColor: MulwiColors.redError,
    //             color: "white"
    //         }}>{locale2.CANCEL_PAYMENT[props.lang]}</Button>)
    // }

    function qrBtn() {
        return (<Button
            fullWidth
            variant="contained"
            onClick={() => props.setInfo(getGenerateQrDialogOptions())}
            style={{
                backgroundColor: MulwiColors.greenDark,
                color: "white"
            }}>{locale2.CREATE_QR[props.lang]}</Button>)
    }

    function sw() {
        if (props.instructor) {
            switch (r.State) {
                case "hold":
                    return (
                        <React.Fragment>
                            {props.grid ? (
                                <React.Fragment>
                                    {r.InstructorDecision === "unset" && (
                                        <Grid item>{acceptBtn()}</Grid>
                                    )}
                                    {r.InstructorDecision !== "reject" && (
                                        <Grid item>{rejectBtn()}</Grid>
                                    )}
                                    <Grid item>{disputeBtn()}</Grid>
                                </React.Fragment>
                            ) : (
                                <Menu
                                    id="simple-menu"
                                    keepMounted
                                    anchorEl={menuOpen}
                                    open={Boolean(menuOpen)}
                                    onClose={() => setMenuOpen(null)}>
                                    {menuDetails()}
                                    <MenuItem>
                                        <Grid container direction="column">
                                            <Typography>{locale2.MAKE_DECISION_ABOUT_RSV[props.lang]}</Typography>
                                            {(r.InstructorDecision !== "unset" && (
                                                <Typography>{locale2.CURRENT_DECISION[props.lang]} <strong>{instrDecisionDisplay(r.InstructorDecision, null, props.lang)}</strong></Typography>
                                            )) || (
                                                    r.ManualConfirm && (
                                                        <Typography variant="caption" style={{
                                                            color: MulwiColors.redError
                                                        }}>{sprintf(locale2.IF_U_DONT_CONFIRM_BEFORE_FMT[props.lang], getUntilDateLabel(new Date(r.SmTimeout), 1, props.lang))}</Typography>)
                                                )}
                                        </Grid>
                                    </MenuItem>
                                    {r.InstructorDecision === "unset" && (
                                        <MenuItem>{acceptBtn()}</MenuItem>
                                    )}
                                    {r.InstructorDecision !== "reject" && (
                                        <MenuItem>{rejectBtn()}</MenuItem>
                                    )}
                                    <MenuItem item>{disputeBtn()}</MenuItem>

                                </Menu>
                            )}
                        </React.Fragment>
                    )
                case "payout":
                case "capture":
                    return props.grid ? (
                        <React.Fragment>
                            <Grid item>{instrCancelBtn()}</Grid>
                            <Grid item>{disputeBtn()}</Grid>
                        </React.Fragment>
                    ) : ((<Menu
                        id="simple-menu"
                        keepMounted
                        anchorEl={menuOpen}
                        open={Boolean(menuOpen)}
                        onClose={() => setMenuOpen(null)}>
                        {menuDetails()}
                        <MenuItem>
                            {instrCancelBtn()}
                        </MenuItem>
                        <MenuItem item>{disputeBtn()}</MenuItem>
                    </Menu>))
                default:
                    return (props.grid ? null : (<Menu
                        id="simple-menu"
                        keepMounted
                        anchorEl={menuOpen}
                        open={Boolean(menuOpen)}
                        onClose={() => setMenuOpen(null)}>
                        {menuDetails()}
                    </Menu>))
            }
        } else {
            switch (r.State) {
                case "link":
                case "link_express":
                    return props.grid ? (
                        <React.Fragment>
                            <Grid item>{qrBtn()}</Grid>
                            {/* <Grid item>{linkExpBtn()}</Grid> */}
                            <Grid item>{disputeBtn()}</Grid>
                        </React.Fragment>
                    ) : (
                        <React.Fragment>
                            <Menu
                                id="simple-menu"
                                keepMounted
                                anchorEl={menuOpen}
                                open={Boolean(menuOpen)}
                                onClose={() => setMenuOpen(null)}>
                                {menuDetails()}
                                {/* <MenuItem>{linkExpBtn()}</MenuItem> */}
                                <MenuItem item>{disputeBtn()}</MenuItem>
                                <MenuItem item>{qrBtn()}</MenuItem>
                            </Menu>
                        </React.Fragment>
                    )
                    break
                case "hold":
                    return props.grid ? (
                        <React.Fragment>
                            <Grid item>{qrBtn()}</Grid>
                            <Grid item>{userCancelBtn()}</Grid>
                            <Grid item>{disputeBtn()}</Grid>
                        </React.Fragment>
                    ) : (
                        <React.Fragment>
                            <Menu
                                id="simple-menu"
                                keepMounted
                                anchorEl={menuOpen}
                                open={Boolean(menuOpen)}
                                onClose={() => setMenuOpen(null)}>
                                {menuDetails()}
                                <MenuItem>{userCancelBtn()}</MenuItem>
                                <MenuItem item>{disputeBtn()}</MenuItem>
                                <MenuItem item>{qrBtn()}</MenuItem>
                            </Menu>
                        </React.Fragment>
                    )
                case "payout":
                    return props.grid ? (
                        <React.Fragment>
                            <Grid item>{qrBtn()}</Grid>
                            <Grid item>{disputeBtn()}</Grid>
                        </React.Fragment>
                    ) : (
                        <React.Fragment>
                            <Menu
                                id="simple-menu"
                                keepMounted
                                anchorEl={menuOpen}
                                open={Boolean(menuOpen)}
                                onClose={() => setMenuOpen(null)}>
                                {menuDetails()}
                                <MenuItem>{disputeBtn()}</MenuItem>
                                <MenuItem item>{qrBtn()}</MenuItem>
                            </Menu>
                        </React.Fragment>
                    )
                case "pre_payout":
                case "capture":
                    return props.grid ? (
                        <React.Fragment>
                            <Grid item>{qrBtn()}</Grid>
                            {props.rsv.DateStart > new Date() ? (
                                <Grid item>{disputeBtn()}</Grid>
                            ) : (
                                <React.Fragment>
                                    <Grid item>{userRefund()}</Grid>
                                    <Grid item>{disputeBtn()}</Grid>
                                </React.Fragment>
                            )}
                        </React.Fragment>
                    ) : (
                        <Menu
                            id="simple-menu"
                            keepMounted
                            anchorEl={menuOpen}
                            open={Boolean(menuOpen)}
                            onClose={() => setMenuOpen(null)}>
                            {menuDetails()}
                            {props.rsv.DateStart > new Date() ? (
                                <MenuItem>{disputeBtn()}</MenuItem>
                            ) : (
                                <React.Fragment>
                                    <MenuItem>{userRefund()}</MenuItem>
                                    <MenuItem item>{disputeBtn()}</MenuItem>
                                    <MenuItem item>{qrBtn()}</MenuItem>
                                </React.Fragment>
                            )}
                        </Menu>)
                default:
                    return props.grid ? <React.Fragment>
                        <Grid item>{qrBtn()}</Grid>
                    </React.Fragment> : (<Menu
                        id="simple-menu"
                        keepMounted
                        anchorEl={menuOpen}
                        open={Boolean(menuOpen)}
                        onClose={() => setMenuOpen(null)}>
                        {menuDetails()}
                        <MenuItem item>{qrBtn()}</MenuItem>
                    </Menu>)
            }
        }
    }

    return (
        <React.Fragment>
            {sw()}
            {!props.grid && (
                <Fab size="small" aria-controls="simple-menu" aria-haspopup="true"
                    style={{ backgroundColor: MulwiColors.blueLight, color: "white", marginRight: 15 }}
                    onClick={(e) => setMenuOpen(e.currentTarget)}>
                    <MoreHoriz />
                </Fab>
            )}
        </React.Fragment>
    )
}