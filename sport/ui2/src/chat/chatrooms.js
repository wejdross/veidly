import {
    Grid, List,
    ListItem, ListItemIcon,
    Typography
} from "@mui/material";
import { Chat } from "@mui/icons-material";
import React, { useEffect, useRef, useState } from "react";
import { MulwiColors } from "../mulwiColors";
import { getOffset, offToStr } from "./commons";

export function ChatroomListItem(props) {

    const c = props.chatroom
    
    // 

    function hasUnreadMsgs() {
        return props.nots && props.nots[c.ChatRoomID]
    }

    function getStatusColor(key, selKey) {
        if(hasUnreadMsgs())
            return "orange"
        if(key === selKey) {
            if(props.isConnected) {
                if(hasUnreadMsgs())
                    return null
                return MulwiColors.greenDark
            }
            return MulwiColors.redError
        }
        return null
    }

    if (!c)
        return null

    return (
        <ListItem button
            onClick={() => {
                props.onClick(c)
                if(hasUnreadMsgs())
                    props.nots.rm(c.ChatRoomID)
            }}
            selected={props.selectedChatroom &&
                c.ChatRoomID === props.selectedChatroom.ChatRoomID}>
            <ListItemIcon>
                <Chat style={{
                    color: getStatusColor(
                        c.ChatRoomID,
                        props.selectedChatroom && props.selectedChatroom.ChatRoomID)
                }} />
            </ListItemIcon>
            <Grid style={{
                width: "100%"
            }} container justifyContent="space-between"
                alignItems="center">
                <Grid item>
                    <Typography>{c.Data.ChatRoomName || "?"}</Typography>
                </Grid>
                <Grid item>
                    <Grid container direction="row-reverse" alignItems="center">
                        {/* {props.nots && (
                            <Grid item>
                                {(props.nots[c.ChatRoomID] && (<Typography style={{
                                    color: "orange",
                                    marginLeft: 5
                                }}>
                                    <strong>!</strong>
                                </Typography>)) || null}
                            </Grid>
                        )} */}
                        {/* <Grid item>
                            {connIndicator(getStatusColor(
                                c.ChatRoomID,
                                props.selectedChatroom && props.selectedChatroom.ChatRoomID,
                                props.isConnected))}
                        </Grid> */}
                    </Grid>
                </Grid>
            </Grid>
        </ListItem>
    )
}

export function ChatroomList(props) {
    let chatrooms = props.chatrooms
    if (!chatrooms)
        return null
    return (<List 
            style={{overflowY: "auto", maxHeight: offToStr(getOffset(props)+40)}}>
        {chatrooms.map((c, i) => (<React.Fragment key={i}>
            <ChatroomListItem
                {...props}
                chatroom={c} />
        </React.Fragment>))}
    </List>)
}