import {
    Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField, Typography
} from "@mui/material";
import React, { useEffect, useState } from "react";
import { postChatroom, postJoinChatroom, storeChatToken } from "../apicalls/chat";
import { locale2 } from "../locale";
import { MulwiColors } from "../mulwiColors";
import { getRandName } from "../names";
import { GenericLangSelect } from "../navbar/LangSelect";

export const ChatroomFlags = {
    FreeJoin: 1,
    ForceRedirectEnabled: 2
}

export function JoinOrAddChatroom(props) {

    const [addChatOpen, _setAddChatOpen] = useState(false)
    
    const [joinToken, setJoinToken] = useState(null)

    const [chatroomReq, setChatroomReq] = useState({
        TokenType: "user_id",
        UserData: {
            Email: "",
            DisplayName: "",
            ChatRoomName: "",
            IconRelpath: "",
            Language: ""
        },
        Flags: ChatroomFlags.ForceRedirectEnabled
    })

    function setChatroomUserdata(field, val) {
        let cpy = { ...chatroomReq }
        cpy.UserData[field] = val
        setChatroomReq(cpy)
    }

    function setAddChatOpen(o) {
        if(o)
            setChatroomUserdata("ChatRoomName", getRandName())
        _setAddChatOpen(o)
    }

    async function createChatroom(e) {
        e.preventDefault()
        try {
            await postChatroom(chatroomReq)
            props.onChange()
            setAddChatOpen(false)
        } catch (ex) {
            console.log(ex)
        }
    }

    async function tryJoinRoom(joinReq) {
        setJoinToken(joinReq)
        try { 
            let res = await postJoinChatroom(joinReq, true)
            if(res === 202) {
                setAddChatOpen(true)
                return
            }
            let token = JSON.parse(res).Token
            storeChatToken(token)
            props.onChange()
        } catch(ex) {
            console.log(ex)
            return
        }
        return
    }

    useEffect(() => {
        let query = new URLSearchParams(window.location.search)
        let joinToken = query.get("join")
        if (!joinToken)
            return
        let crid = query.get("crid")
        if (!crid)
            return
        tryJoinRoom({
            Token: joinToken,
            ChatRoomID: crid
        })
    }, [])

    useEffect(() => {
        if (!props.user) {
            return
        }
        let cpy = { ...chatroomReq }
        cpy.UserData.Email = props.user.Email
        cpy.UserData.DisplayName = props.user.Name
        cpy.UserData.Language = props.lang
        setChatroomReq(cpy)
    }, [props.user])

    async function joinChatroom(e) {
        e.preventDefault()
        if (!joinToken)
            return
        try {
            let req = joinToken
            req.UserData = chatroomReq.UserData
            await postJoinChatroom(req)
            window.history.replaceState(null, null, window.location.pathname)
            setJoinToken(null)
            setAddChatOpen(false)
            props.onChange()
        } catch (ex) {
            console.log(ex)
        }
    }

    return (<React.Fragment>
        <Dialog open={addChatOpen} onClose={() => setAddChatOpen(false)}>
            <DialogTitle>{joinToken ? locale2.JOIN[props.lang]: locale2.ADD[props.lang]} {locale2.CHANNEL[props.lang]}</DialogTitle>
            <form onSubmit={joinToken ? joinChatroom : createChatroom}>
                <DialogContent>
                    <TextField
                        label={locale2.CHATROOM_NAME[props.lang]}
                        value={chatroomReq.UserData.ChatRoomName}
                        onChange={(e) => setChatroomUserdata('ChatRoomName', e.target.value)}
                        margin="dense" fullWidth required
                    />
                    <TextField
                        label={locale2.DISPLAY_NAME[props.lang]}
                        value={chatroomReq.UserData.DisplayName}
                        onChange={(e) => setChatroomUserdata('DisplayName', e.target.value)}
                        margin="dense" fullWidth required
                    />
                    <TextField
                        label={locale2.NOTIFY_EMAIL[props.lang]}
                        value={chatroomReq.UserData.Email}
                        onChange={(e) => setChatroomUserdata('Email', e.target.value)}
                        margin="dense" fullWidth
                    />
                    <GenericLangSelect 
                            value={chatroomReq.UserData.Language} 
                        onChange={(e) => setChatroomUserdata('Language', e)} />
                    <Typography>
                        {locale2.NOTIFY_EMAIL_DISCLAIMER[props.lang]}
                    </Typography>

                </DialogContent>
                <DialogActions>
                    <Button type="submit" style={{
                        color: "white",
                        backgroundColor: MulwiColors.greenDark
                    }}>
                        {joinToken ? locale2.JOIN[props.lang]: locale2.ADD[props.lang]}
                    </Button>
                    <Button onClick={() => setAddChatOpen(false)}>
                        {locale2.CLOSE[props.lang]}
                    </Button>
                </DialogActions>
            </form>

        </Dialog>
        <Button variant="contained"
            onClick={() => setAddChatOpen(true)}
            style={{
                color: "white",
                borderBottomRightRadius: 0,
                borderBottomLeftRadius: 0,
                backgroundColor: MulwiColors.blueDark
            }} fullWidth>
            {locale2.NEW_CHATROOM[props.lang]}
        </Button>
    </React.Fragment>)
}
