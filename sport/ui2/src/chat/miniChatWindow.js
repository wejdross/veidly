import { Dialog, DialogContent, DialogTitle, Grid, IconButton, Modal, Paper } from "@mui/material";
import React, { useEffect, useRef, useState } from "react";
import { useHistory, useLocation } from "react-router-dom/cjs/react-router-dom.min";
import { connIndicator } from "./commons";
import Draggable from 'react-draggable';
import { MulwiColors } from "../mulwiColors";
import { Chatroom } from "./chatroom";
import { Close, Fullscreen, Minimize, MinimizeOutlined } from "@mui/icons-material";
import { getPathFromSearchParams, QSsetAndReturn } from "../helpers";

export function DraggablePaper(props) {
    return (
        <Draggable
            handle="#draggable-dialog-title"
            cancel={'[class*="MuiDialogContent-root"]'}
        >
            <Paper {...props} />
        </Draggable>
    );
}

function getMinichatQuerystring() {
    let q = new URLSearchParams(window.location.search)
    if (!q)
        return null
    let minichat = q.get("minichat")
    if (!minichat)
        return null
    return minichat
}

function setMinichatQuerystring(val) {
    return QSsetAndReturn("minichat", JSON.stringify(val))
}

export function rmChatroomMiniWindow(crid) {
    let minichat = getMinichatQuerystring()
    try {
        minichat = JSON.parse(decodeURIComponent(minichat))
    } catch (ex) {
        return new URLSearchParams(window.location.search)
    }
    if(!minichat)
        return new URLSearchParams(window.location.search)
    let cpy = []
    for(let i = 0; i < minichat.length; i++) {
        if(minichat[i] != crid)
            cpy.push(minichat[i])
    }
    return setMinichatQuerystring(cpy) 
}

export function addChatroomToMiniWindow(crid) {
    let minichat = getMinichatQuerystring()
    if (!minichat) {
        return setMinichatQuerystring([crid])
    }
    try {
        minichat = JSON.parse(decodeURIComponent(minichat))
    } catch (ex) {
        minichat = []
    }
    if (!minichat || minichat.length == 0) {
        return setMinichatQuerystring([crid])
    }
    for (let i = 0; i < minichat.length; i++) {
        // PROPER CHECK DUPLICATES
        if (minichat[i] == crid)
            continue
        minichat.push(crid)
        return setMinichatQuerystring(minichat) 
    }
}

export function MiniChatWindow(props) {

    const [chatRoomIDs, setChatRoomIDs] = useState([])

    async function refresh(crid) {
        setChatRoomIDs(crid)
    }

    const location = useLocation()

    useEffect(() => {

        if (!props.chatToken)
            return

        let minichat = getMinichatQuerystring()
        if (!minichat)
            return

        try {
            minichat = JSON.parse(decodeURIComponent(minichat))
        } catch (ex) {
            console.log(ex)
            return
        }

        refresh(minichat)

    }, [location, props.chatToken])

    const [peer, setPeer] = useState({})

    const height = 500

    const history = useHistory()

    return (<React.Fragment>
        {chatRoomIDs.map(cid => (<Dialog key={cid}
            disableEnforceFocus hideBackdrop
            aria-labelledby="draggable-dialog-title"
            style={{
                pointerEvents: 'none',
            }}
            PaperProps={{
                style: {
                    pointerEvents: 'auto',
                    height: height,
                    width: 400
                }
            }}
            PaperComponent={DraggablePaper} open={true}>
            <DialogTitle style={{
                backgroundColor: MulwiColors.blueDark,
                color: "white",
                cursor: "pointer"
            }} id="draggable-dialog-title">
                <Grid container direction="row"
                    alignItems="center"
                    justifyContent="space-between">

                    <Grid item>
                        <Grid container direction="row" spacing={1} alignItems="center">
                            <Grid item>
                                {connIndicator(
                                    (peer[cid] && peer[cid].IsConnected && MulwiColors.greenDark)
                                    || MulwiColors.redError)}
                            </Grid>
                            <Grid item>
                                {(peer[cid] && peer[cid].DisplayName)}
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item>
                        <Grid container direction="row" alignItems="center">
                            <IconButton size="small" style={{
                                color: "white"
                            }}>
                                <Minimize />
                            </IconButton>
                            <IconButton size="small" onClick={() => {
                                let q = new URLSearchParams(window.location.search);
                                q.set("open", encodeURIComponent(cid))
                                history.push(getPathFromSearchParams(q))
                                q = rmChatroomMiniWindow(cid)
                                history.push("/chat?" + q.toString())
                            }} style={{
                                color: "white"
                            }}>
                                <Fullscreen />
                            </IconButton>
                            <IconButton size="small" style={{
                                color: "white"
                            }} onClick={() => {
                                let q = rmChatroomMiniWindow(cid)
                                history.push(getPathFromSearchParams(q))
                            }}>
                                <Close />
                            </IconButton>
                        </Grid>
                    </Grid>
                </Grid>
            </DialogTitle>
            <DialogContent style={{
                padding: '0px 0px 0px 0px'
            }}>
                <Chatroom
                    onChatroomMembers={m => {
                        if (!m)
                            return
                        for (let i = 0; i < m.length; i++) {
                            let el = m[i]
                            if (el.You)
                                continue
                            peer[cid] = el
                            setPeer({ ...peer })
                            return
                        }
                    }}
                    chatroomID={cid}
                    isConnected={Boolean(cid)}
                    height={height - 120}
                />
            </DialogContent>
        </Dialog>))}
    </React.Fragment>)
}