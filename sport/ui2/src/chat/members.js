import { Button, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, Grid, List, ListItem, ListItemIcon, ListItemText, TextField, Typography } from "@mui/material";
import { Person } from "@mui/icons-material";
import React, { useEffect, useRef, useState } from "react";
import { postChatAccessToken } from "../apicalls/chat";
import { MulwiColors } from "../mulwiColors";
import { connIndicator } from "./commons";


export function InviteMember(props) {

    const [chatInvReq, setChatInvReq] = useState(() => {
        let exp = new Date()
        exp.setDate(exp.getDate() + 30)
        return ({
            Uses: 1,
            ExpiresOn: exp,
        })
    })
    const [inviteTokenRes, setInviteTokenRes] = useState(null)

    async function createChatInv(e) {
        if(e)
            e.preventDefault()
        if(!props.chatroom)
            return
        try {
            chatInvReq.ChatRoomID = props.chatroom.ChatRoomID
            let res = await postChatAccessToken(chatInvReq)
            setInviteTokenRes(JSON.parse(res))
        } catch (ex) {
            console.log(ex)
        }
    }

    if(props.chatroom && !props.isConnected) {
        return (
            <React.Fragment>
                <center><CircularProgress/></center>
            </React.Fragment>
        )
    }

    if(!props.chatroom || !props.isConnected)
        return null

    return (<React.Fragment>

        <Dialog fullWidth 
                    open={Boolean(inviteTokenRes)} 
                    onClose={() => setInviteTokenRes(null)}>
                <DialogTitle>
                    Your invite link
                </DialogTitle>
                <DialogContent>
                    <TextField fullWidth multiline 
                        minRows={1} maxRows={3} 
                        value={(inviteTokenRes && inviteTokenRes.JoinLink) || ""} />
                </DialogContent>
                <DialogActions>
                    <Button type="submit" style={{
                        color: "white",
                        backgroundColor: MulwiColors.blueDark
                    }} onClick={() => {
                        navigator.clipboard.writeText(inviteTokenRes.JoinLink)
                    }}>
                        Copy
                    </Button>
                    <Button onClick={() => setInviteTokenRes(null)}>
                        Close
                    </Button>
                </DialogActions>
        </Dialog>

        <Button variant="contained"
            onClick={createChatInv}
            style={{
                color: "white",
                backgroundColor: MulwiColors.blueDark
            }} fullWidth>

            Invite
        </Button>
    </React.Fragment>)
}

export function ChatroomMembers(props) {

    let ms = props.members

    if (!ms || !props.isConnected)
        return null

    return (<React.Fragment>

        <List>
            {ms.map((m, i) => (<ListItem key={i} button>
                <ListItemIcon>
                    <Person />
                </ListItemIcon>
                <ListItemText>
                    <Grid container direction="row" alignItems="center"
                        spacing={2}>
                        <Grid item>
                            {connIndicator(m.IsConnected ? MulwiColors.greenDark : MulwiColors.redError)}
                        </Grid>
                        <Grid item>
                            {m.DisplayName}
                        </Grid>
                        <Grid item>
                            {m.You && <Typography style={{
                                fontStyle: "italic",
                                color: "grey"
                            }}>
                                (you)</Typography>}
                        </Grid>

                    </Grid>
                </ListItemText>
            </ListItem>))}
        </List>
    </React.Fragment>)
}
