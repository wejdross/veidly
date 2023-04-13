import { Button, Fab, Grid,
        Menu, MenuItem } from '@mui/material';
import { MoreHoriz } from '@mui/icons-material';
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { postExpire, postUserDispute, 
        postUserRefund, createQr } from '../apicalls/sm';
import { MulwiColors } from '../mulwiColors';
import { getDialogWithOptions, getErrorDialog, 
    getNullDialog, supportEmail } from '../StatusDialog';
import { DisputeForm } from './rsvInteract';
import { locale2 } from '../locale';
import { makeStyles } from "@mui/styles";
import useMediaQuery from '@mui/material/useMediaQuery';
const linkStyles = makeStyles(t => ({
    link: {
        textDecoration: "none",
        color: MulwiColors.blueDark
    }
}))


export function SubInteractMenu(props) {

    const [menuOpen, setMenuOpen] = useState(null)

    const c = linkStyles()

    let s = props.sub

    async function generateQr() {
        try {
            let qr = await createQr(s.ID, props.at)
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
            await postUserRefund(s.ID, props.at)
            if(props.onChange) await props.onChange()
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
            await postUserDispute(s.ID, props.at, email, msg)
            if(props.onChange) await props.onChange()
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
            await postExpire(s.ID, props.at)
            if(props.onChange) await props.onChange()
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

    function getLinkExpDialogOptions() {
        return (getDialogWithOptions(
            locale2.ARE_YOU_SURE[props.lang],
            locale2.ARE_YOU_SURE_YOU_WANT_TO_CANCEL[props.lang],
            <Button variant="contained" style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }} onClick={() => doLinkExpire()}>
                Tak
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
            <DisputeForm lang={props.lang} onChange={(e, m) => doUserDispute(e, m)}/>, 
            null))
    }

    if (!s) {
        return null
    }

    function menuDetails() {
        return (!props.noDetails && (
                <MenuItem><Button style={{
                    color: MulwiColors.blueDark
                }} fullWidth><Link className={c.link}
                    to={"/sub_details?id=" + s.ID + ((props.instructor && "&instr=1") || "")}>
                        {locale2.DETAILS[props.lang]}</Link></Button>
                    </MenuItem>
        )) || null
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
                {locale2.REPORT_ISSUE[props.lang]}
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


    function linkExpBtn() {
        return (<Button
            fullWidth
            variant="contained"
            onClick={() => props.setInfo(getLinkExpDialogOptions())}
            style={{
                backgroundColor: MulwiColors.redError,
                color: "white"
            }}>{locale2.CANCEL_PAYMENT[props.lang]}</Button>)
    }

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
            switch (s.State) {
                case "hold":
                    return (
                        <React.Fragment>
                            {props.grid ? (
                                <React.Fragment>
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
                                <MenuItem item>{disputeBtn()}</MenuItem>
                            </Menu>
                            )}
                        </React.Fragment>
                    )
                case "payout":
                case "capture":
                    return props.grid ? (
                        <React.Fragment>
                            {/* <Grid item>{instrCancelBtn()}</Grid> */}
                            <Grid item>{disputeBtn()}</Grid>
                        </React.Fragment>
                    ) : ( (<Menu
                            id="simple-menu"
                            keepMounted
                            anchorEl={menuOpen}
                            open={Boolean(menuOpen)}
                            onClose={() => setMenuOpen(null)}>
                       {menuDetails()}
                        {/* <MenuItem>{instrCancelBtn()}</MenuItem> */}
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
            switch (s.State) {
                case "link":
                case "link_express":
                    return props.grid ? (
                        <React.Fragment>
                            <Grid item>{qrBtn()}</Grid>
                            <Grid item>{linkExpBtn()}</Grid>
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
                                    <MenuItem>{linkExpBtn()}</MenuItem>
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
                            {props.sub.DateStart > new Date() ? (
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
                                {props.sub.DateStart > new Date() ? (
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
                    style={{ backgroundColor: MulwiColors.blueLight, color: "white" }}
                    onClick={(e) => setMenuOpen(e.currentTarget)}>
                    <MoreHoriz />
                </Fab>
            )}
        </React.Fragment>
    )
}