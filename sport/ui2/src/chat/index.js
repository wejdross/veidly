import { Close, OpenInBrowser } from "@mui/icons-material";
import MenuIcon from '@mui/icons-material/Menu';
import { Drawer, Fab, Grid, IconButton, Typography, useMediaQuery, useTheme } from "@mui/material";
import { makeStyles } from "@mui/styles";
import clsx from 'clsx';
import React, { useEffect, useState } from "react";
import { useHistory, useLocation } from "react-router-dom/cjs/react-router-dom.min";
import { getChatroomsForUser } from "../apicalls/chat";
import { getPathFromSearchParams, isUuid, QSsetAndReturn } from "../helpers";
import { locale2 } from "../locale";
import { MulwiColors } from "../mulwiColors";
import {
    getErrorDialog, getNullDialog,
    StatusDialog
} from "../StatusDialog";
import { Chatroom } from "./chatroom";
import { ChatroomList } from "./chatrooms";
import { getOffset, offToStr } from "./commons";
import { JoinOrAddChatroom } from "./joinAddChatroom";
import { ChatroomMembers as ChatroomMemberList, InviteMember } from "./members";
import { addChatroomToMiniWindow } from "./miniChatWindow";

const useStyles = makeStyles((theme) => ({
    text: {
        padding: theme.spacing(2, 2, 0),
    },
    paper: {
        paddingBottom: 50,
    },
    list: {
        marginBottom: theme.spacing(2),
    },
    subheader: {
        backgroundColor: theme.palette.background.paper,
    },
    appBar: {
        top: 'on',
    },
    grow: {
        flexGrow: 1,
    },
    fabButton: {
        position: 'absolute',
        zIndex: 1,
        top: -30,
        left: 0,
        right: 0,
        margin: '0 auto',
    },
    white: {
        color: "white"
    },
    itemlist: {
        width: 250,
    },
    itemfullList: {
        width: 'auto',
    },
}));

export function ChatIndex(props) {
    useEffect(() => {
        window.scrollTo(0, 0);
      }, []);
    
    const classes = useStyles();

    const [chatrooms, setChatrooms] = useState([])
    const [selectedChatroom, setSelectedChatroom] = useState(null)

    const theme = useTheme()
    const isSmall = useMediaQuery(theme.breakpoints.down('sm'))
    const [info, setInfo] = useState(getNullDialog())

    const [isConnected, setIsConnected] = useState(false)
    const [chatroomMembers, setChatroomMembers] = useState(null)

    async function refreshChatrooms() {
        try {
            let chatrooms = await getChatroomsForUser()
            chatrooms = JSON.parse(chatrooms)
            setChatrooms(chatrooms || [])
        } catch (ex) {
            setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    const location = useLocation()
    const history = useHistory()

    async function closeActiveChatroom() {
        setSelectedChatroom(null)
        setIsConnected(false)
        setChatroomMembers(null)
        toggleDrawer(false)
    }

    async function openChatMiniwindow() {
        if (!selectedChatroom || !selectedChatroom.ChatRoomID)
            return
        let q = addChatroomToMiniWindow(selectedChatroom.ChatRoomID)
        if (!q)
            return
        history.push(window.location.pathname + "?" + q.toString())
        closeActiveChatroom()
    }

    useEffect(() => {
        if (!chatrooms)
            return
        let q = new URLSearchParams(window.location.search)
        if (!q)
            return
        let open = q.get("open")
        if (!open)
            return
        if (selectedChatroom && selectedChatroom == open) {
            return
        }
        for (let i = 0; i < chatrooms.length; i++) {
            if (chatrooms[i].ChatRoomID === open) {
                setSelectedChatroom(chatrooms[i])
                let p = getPathFromSearchParams(QSsetAndReturn("open", ""))
                history.push(p)
                return
            }
        }
    }, [location, chatrooms])

    useEffect(() => {
        if (!props.chatToken)
            return

        let ch = {}
        if (chatrooms) {
            for (let i = 0; i < chatrooms.length; i++) {
                ch[chatrooms[i].ChatRoomID] = 1
            }
        }

        if (!props.nots) {
            console.log("ChatIndex: null nots - refreshing chatrooms")
            refreshChatrooms()
            return
        }

        for (let k in props.nots) {
            // uuid
            if (!isUuid(k)) {
                continue
            }
            if (!ch[k]) {
                console.log("ChatIndex: got more chatrooms in notification - refreshing chatrooms")
                refreshChatrooms()
                return
            }
        }

    }, [props.nots, props.chatToken, props.user])

    const [open, setOpen] = React.useState(true);

    const toggleDrawer = (open) => (event) => {
        if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
            return;
        }
        setOpen(open)
    }

    return (
        <React.Fragment>
            <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />

            {isSmall ? (<React.Fragment>

                <div style={{ marginTop: 60 }}></div>

                {selectedChatroom ? (<React.Fragment>

                    <Chatroom
                        lang={props.lang}
                        isSmall={isSmall}
                        chatroomID={(selectedChatroom && selectedChatroom.ChatRoomID)}
                        isConnected={isConnected}
                        nots={props.nots}
                        height="calc(100vh - 55px - 60px)"
                        onChatroomMembers={m => setChatroomMembers(m)}
                        onConnectChange={(c) => setIsConnected(c)}
                        setInfo={setInfo}
                        instructor={props.instructor}
                    />

                    <Fab style={{
                        position: "absolute",
                        right: 5,
                        top: 60,
                        color: "white",
                        backgroundColor: MulwiColors.blueDark
                    }} onClick={() => setOpen(true)}>
                        <MenuIcon />
                    </Fab>

                    <Drawer anchor="left" 
                            open={open} onClose={toggleDrawer(false)}>
                        <div
                            className={clsx(classes.itemlist, {
                                [classes.itemfullList]: false,
                            })}
                            role="presentation"
                            onClick={toggleDrawer(false)}
                            onKeyDown={toggleDrawer(false)}>

                            {selectedChatroom && (
                                <center>
                                    <Typography variant="h6" style={{
                                        borderTop: "solid " + MulwiColors.lightGreyAddedByLukasz + " 2px"
                                    }}>
                                        {selectedChatroom && selectedChatroom.Data && selectedChatroom.Data.ChatRoomName}
                                    </Typography>
                                </center>
                            )}

                                {selectedChatroom && (<Grid container
                                        justifyContent="center"
                                        direction="row">
                                    <Grid item>
                                        <IconButton style={{
                                            color: MulwiColors.blueDark,
                                            borderRadius: 0
                                        }} onClick={openChatMiniwindow}>
                                            <OpenInBrowser />
                                        </IconButton>
                                        <IconButton style={{
                                            color: MulwiColors.redError,
                                            borderRadius: 0
                                        }} onClick={closeActiveChatroom}>
                                            <Close />
                                        </IconButton>
                                    </Grid>
                                </Grid>)}

                            <ChatroomMemberList
                                members={chatroomMembers} isConnected={isConnected} />

                            {/* <InviteMember chatroom={selectedChatroom} isConnected={isConnected} /> */}

                            {/* <Button style={{
                                marginTop: 5,
                            }} variant="contained" fullWidth>
                                {locale2.CLOSE[props.lang]}
                            </Button> */}
                        </div>
                    </Drawer>

                </React.Fragment>) : (<React.Fragment>

                    <JoinOrAddChatroom
                        lang={props.lang}
                        user={props.user}
                        onChange={refreshChatrooms} />
                    <div style={{ height: "50%", overflow: "auto" }}>
                        <ChatroomList
                            selectedChatroom={selectedChatroom}
                            isConnected={isConnected}
                            chatrooms={chatrooms}
                            nots={props.nots}
                            onClick={(c) => setSelectedChatroom(c)}
                        />
                    </div>

                </React.Fragment>)}


            </React.Fragment>) : (<React.Fragment>

                <Grid
                    container direction="row"
                    style={{
                        height: offToStr(getOffset(props)),
                        marginTop: (isSmall && 60) || null
                    }}>
                    <Grid item sm={3}>
                        <div style={{ marginTop: 2 }}></div>
                        <JoinOrAddChatroom
                            lang={props.lang}
                            user={props.user}
                            onChange={refreshChatrooms} />

                        <div style={{ height: "50%", overflow: "auto" }}>
                            <ChatroomList
                                selectedChatroom={selectedChatroom}
                                isConnected={isConnected}
                                chatrooms={chatrooms}
                                nots={props.nots}
                                onClick={(c) => setSelectedChatroom(c)}
                            />
                        </div>

                        {selectedChatroom && (
                            <center>
                                <Typography variant="h6" style={{
                                    borderTop: "solid " + MulwiColors.lightGreyAddedByLukasz + " 2px"
                                }}>
                                    {selectedChatroom && selectedChatroom.Data && selectedChatroom.Data.ChatRoomName}
                                </Typography>
                            </center>
                        )}

                        <ChatroomMemberList
                            members={chatroomMembers} isConnected={isConnected} />
                        <InviteMember chatroom={selectedChatroom} isConnected={isConnected} />
                    </Grid>
                    <Grid item sm={9}>
                        {selectedChatroom && (<Grid container direction="row">
                            <Grid item>
                                <IconButton style={{
                                    color: MulwiColors.blueDark,
                                    borderRadius: 0
                                }} onClick={openChatMiniwindow}>
                                    <OpenInBrowser />
                                </IconButton>
                                <IconButton style={{
                                    color: MulwiColors.redError,
                                    borderRadius: 0
                                }} onClick={closeActiveChatroom}>
                                    <Close />
                                </IconButton>
                            </Grid>
                        </Grid>)}
                        <Chatroom
                            lang={props.lang}
                            chatroomID={(selectedChatroom && selectedChatroom.ChatRoomID)}
                            isConnected={isConnected}
                            nots={props.nots}
                            onChatroomMembers={m => setChatroomMembers(m)}
                            onConnectChange={(c) => setIsConnected(c)}
                            setInfo={setInfo}
                            instructor={props.instructor}
                        />
                    </Grid>
                </Grid>

            </React.Fragment>)}

        </React.Fragment>
    )
}
