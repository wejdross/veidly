import { CircularProgress, Grid, TextField } from "@mui/material";
import React, { useEffect, useRef, useState } from "react";
import { newWsConn, postWsToken } from "../apicalls/chat";
import { prettyPrintDate, prettyPrintHr } from "../harmonogram/trainingDetails";
import { dateToMicroEpoch, microEpochToDate, sleep } from "../helpers";
import { locale2 } from "../locale";
import { getErrorDialog } from "../StatusDialog";
import { getLastMessageCache, updateLastMessageCache } from "./chatNotification";
import { getOffset, isWindowVisible, offToStr, startNotify } from "./commons";
import InputAdornment from '@mui/material/InputAdornment';
import IconButton from '@mui/material/IconButton';
import FormControl from '@mui/material/FormControl';
import SendIcon from '@mui/icons-material/Send';
import Input from '@mui/material/Input';
import { MulwiColors } from "../mulwiColors";

export async function closeWsAndWait(ws) {
    if (!ws)
        return
    try {
        await ws.close()
        let i
        for (i = 0; i < 100; i++) {
            if (ws.readyState === 3)
                break
            await sleep(10)
        }
    } catch (ex) {
        console.log("error closing connection: " + ex)
    }
}

export function Chatroom(props) {

    const [inputMsg, setInputMsg] = useState("")
    const chatcontainer = useRef()
    const chatgrid = useRef()
    const conn = useRef(null)
    const msgMeta = useRef({})
    const memberHash = useRef()

    function isAtBottom() {
        let t = chatcontainer.current
        if (!t)
            return false
        return (((t.scrollHeight - t.scrollTop) - t.clientHeight) < 30)
    }

    function isAtTop() {
        let t = chatcontainer.current
        if (!t)
            return false
        return (t.scrollTop === 0)
    }

    function requestMsgFeed(start, end) {
        if (!conn.current)
            return
        let feedReq = {
            Type: 1,
            FeedOpts: {
                Limit: 100,
                Start: start || null,
                End: end || null,
            }
        }
        try {
            conn.current.send(JSON.stringify(feedReq))
        } catch (ex) {
            console.log(ex)
        }
    }

    async function scrolledToTop() {
        if(!props.chatroomID)
            return
        let mm = msgMeta.current[props.chatroomID]
        if (!mm || !mm.fts)
           return
        requestMsgFeed(0, mm.fts + 10000)
    }

    async function sendMsg(e) {
        e.preventDefault()

        if(inputbox.current === document.activeElement) {
            console.log("is in focus")
        }
        if (!conn.current)
            return
        try {
            conn.current.send(JSON.stringify({
                Type: 2,
                Msg: inputMsg
            }))
        } catch (ex) {
            props.setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
            console.log(ex)
        }
        setInputMsg("")
    }

    function confirmMsg(timestamp) {
        console.log("confirming last msg")
        try {
            updateLastMessageCache(props.chatroomID, timestamp)  
            conn.current.send(JSON.stringify({
                Type: 3,
                ReadMsgTimestamp: Number(timestamp)
            }))
            if(props.nots)
                props.nots.rm(props.chatroomID)
        } catch (ex) {
            props.setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
            console.log(ex)
        }
    }

    function confirmLastMsg() {
        if(!isAtBottom())
            return
        
        let lastElement = chatgrid.current.lastChild
        if(!lastElement)
            return

        let lastTimestamp = lastElement.id
        let lmc = getLastMessageCache()
        if(lmc) {
            lmc = JSON.parse(lmc)
            if(lmc) {
                lmc = lmc[props.chatroomID]
                if(!lastTimestamp || lastTimestamp <= lmc)
                    return
            }
        }
        confirmMsg(lastTimestamp)
    }

    function scrollToBottom() {
        confirmLastMsg()
    }

    async function closeConnIfExists() {
        await closeWsAndWait(conn.current)
        conn.current = null
    }

    function handleMemberReport(members) {
        let hash = {}
        for (let i = 0; i < members.length; i++) {
            hash[members[i].UserID] = members[i]
        }
        memberHash.current = hash
        let s = new Date()
        s.setHours(s.getHours() - 2)
        requestMsgFeed(dateToMicroEpoch(s), 0)
        props.onChatroomMembers && props.onChatroomMembers(members)
    }

    function printMsgTimestamp(ts) {
        let d = microEpochToDate(ts)
        let now = new Date()
        if (now.getDay() === d.getDay()
            && now.getMonth() === d.getMonth()
            && now.getFullYear() === d.getFullYear()) {
            return prettyPrintHr(d)
        } else {
            return prettyPrintDate(d)
        }
    }

    // noop since i use innerText everywhere.
    function __sanitize(html) {
        return html
        // const decoder = document.createElement('div')
        // decoder.innerHTML = html
        // return decoder.textContent
    }

    const inputbox = useRef()

    function handleMsgFeed(msgs) {

        let crid = props.chatroomID
        if(!crid)
            return

        let crData = msgMeta.current[crid]
        // TODO: dont group fts, lts by userid, store msgs instead 
        if (!crData) {
            msgMeta.current[crid] = {
                fts: null,
                lts: null
            }
            crData = msgMeta.current[crid]
        }

        let wasAppend = false
        let prependPos = null
        let wasAtBottom = isAtBottom()
        let lastTimestamp = null

        for (let i = 0; i < msgs.length; i++) {

            let newMsg = msgs[i]
            let fts = crData.fts
            let lts = crData.lts

            // duplicate check; prove that m.Msgs are always sorted by Timestamp asc
            // and you may optimize this code
            if (fts && lts && fts <= newMsg.Timestamp && lts >= newMsg.Timestamp)
                continue

            if (!fts || newMsg.Timestamp < fts)
                crData.fts = newMsg.Timestamp
            if (!lts || newMsg.Timestamp > lts)
                crData.lts = newMsg.Timestamp

            let mem = memberHash.current[newMsg.UserID]
            let author = (mem && mem.DisplayName) || "unkown"

            let colItem = document.createElement("div")
            colItem.className = "MuiGrid-root MuiGrid-item"
            colItem.id = newMsg.Timestamp
            let rowContainer = document.createElement("div")
            rowContainer.className = "MuiGrid-root MuiGrid-container MuiGrid-spacing-xs-1 MuiGrid-align-items-xs-center"

            //const lgMarginStyle = "margin-left: 10px; padding: 0px; margin-bottom: -12px;"
            const smMarginStyle = "margin-left: 10px; padding: 0px; "

            rowContainer.style = smMarginStyle
            rowContainer.id = newMsg.Timestamp + "-rowc"

            let rowItem = document.createElement("div")
            rowItem.className = "MuiGrid-root MuiGrid-item"
            let authorP = document.createElement("p")
            authorP.className = "MuiTypography-root MuiTypography-body2"
            authorP.style = "color: rgb(96, 96, 96);"
            authorP.id = newMsg.Timestamp + "-authorP"
            authorP.dataset.author = __sanitize(author)
            // authorP.style = "color: orange;"
            authorP.innerText = __sanitize(printMsgTimestamp(newMsg.Timestamp) + " " + author + ": ")
            let rowItem1 = document.createElement("div")
            rowItem1.className = rowItem.className
            let contentP = document.createElement("p")
            contentP.className = "MuiTypography-root MuiTypography-body1"
            contentP.style = "margin-bottom: 2px; margin-top: 1px; white-space: -moz-pre-wrap !important; word-break: break-all; white-space: -webkit-pre-wrap; word-wrap: break-word;"
            contentP.innerText = __sanitize(newMsg.Content)

            rowItem1.appendChild(contentP)
            rowItem.appendChild(authorP)
            rowContainer.appendChild(rowItem)
            rowContainer.appendChild(rowItem1)
            colItem.appendChild(rowContainer)

            // let neighbour = null

            if (newMsg.Timestamp >= lts) {
                wasAppend = true
                lastTimestamp = newMsg.Timestamp
                //neighbour = chatgrid.current.lastChild
                chatgrid.current.appendChild(colItem)
            } else {
                if (!prependPos)
                    prependPos = colItem
                //neighbour = chatgrid.current.firstChild
                chatgrid.current.prepend(colItem)
            }

            
            // i have no strength to deal with this mess right now.
            // to get back: uncomment neighbour assignment above
            // and then following block and remove bugs
            

            // if(neighbour) {
            //     let id = neighbour.id
            //     let lastItemAuthorP = document.getElementById(id + "-authorP")
            //     let nrc = document.getElementById(id + "-rowc")
            //     if(lastItemAuthorP && lastItemAuthorP.dataset.author === __sanitize(author)) {
            //         rowContainer.style = lgMarginStyle
            //         if(nrc)
            //             nrc.style = lgMarginStyle
            //     } else {
            //         if(nrc)
            //             nrc.style = smMarginStyle
            //     }
            // }
        }

        startNotify()

        if(lastTimestamp && wasAtBottom && isWindowVisible) {
            confirmMsg(lastTimestamp)
        }


        if (chatcontainer.current) {
            if (wasAppend && wasAtBottom) {
                chatcontainer.current.scrollTo(0, chatcontainer.current.scrollHeight);
            } else if (prependPos) {
                prependPos.scrollIntoView()
                //chatcontainer.current.scrollTop()
            }
        }
    }

    function handleWsMsg(e) {
        if(!chatgrid.current)
            return
        let m = JSON.parse(e.data)
        switch (m.Type) {
            case 1:
                handleMsgFeed(m.Msgs)
                break
            case 2:
                handleMemberReport(m.Members)
                break
            default:
                console.log("unrecognized wsMsg type")
                return
        }
    }

    async function openChatroom(forceRedirect) {

        await closeConnIfExists()
        
        let crid = props.chatroomID
        msgMeta.current = {}
        let token

        try {
            let wst = await postWsToken({
                ForceRedirectPeers: forceRedirect ? true : false,
                ChatRoomID: crid
            })
            wst = JSON.parse(wst)
            token = wst.Token
        } catch(ex) {
            if(ex == 404) {
                return openChatroom(true)
            }
            throw ex
        }

        let ws = newWsConn(token)
        conn.current = ws

        ws.addEventListener("open", () => {
            console.log("ws opened")
            props.onConnectChange && props.onConnectChange(true)
        })

        ws.addEventListener("close", () => {
            console.log("ws closed")
            props.onConnectChange && props.onConnectChange(false)
        })

        ws.addEventListener("message", e => {
            try {
                handleWsMsg(e)
            } catch (ex) {
                props.setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
                console.log("handleWsMsg err", ex)
            }
        })

        ws.addEventListener("error", (e) => {
            console.log("ws error")
        })

    }

    const interval = useRef()

    async function openChatroomErrHandle(arg) {
        try {
            await openChatroom(arg)
        } catch(ex) {
            console.log(ex)
            // props.setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    async function startConnInterval() {
        if(interval.current) {
            clearInterval(interval.current)
            interval.current = null
        }
        
        await openChatroomErrHandle()

        interval.current = setInterval(async () => {
            if(conn.current && conn.current.readyState === 1)
                return
            await openChatroomErrHandle()
        }, 1500)
    }
    
    useEffect(() => {
        if (!props.chatroomID) {
            closeConnIfExists()
                if(interval.current){
                    clearInterval(interval.current)
                interval.current = null
            }
            return
        }
        startConnInterval()
        return () => {
                if(interval.current) {
                    clearInterval(interval.current)
                interval.current = null
            }
            closeConnIfExists()
        }
    }, [props.chatroomID])

    let isConnected = props.isConnected

    if(!props.chatroomID) {
        return null
    }

    if (!isConnected)
        return (
            <React.Fragment>
                <center><CircularProgress/></center>
            </React.Fragment>
        ) 

    return (<React.Fragment>
        <div 
            onClick={confirmLastMsg}>
        <Grid
            ref={chatcontainer}
            container
            onScroll={() => {
                if (isAtBottom())
                    scrollToBottom()
                if (isAtTop())
                    scrolledToTop()
            }}
            direction="column"
            style={{
                height: props.height || offToStr(getOffset(props) + 60 + 50),
                backgroundColor: "white",
                overflowY: "auto",
                overflowX: "hidden"
            }}
            justifyContent="space-between">

            <Grid item ref={chatgrid} >

                {/* <Grid container dens style={{
                    marginLeft: 10, padding: 0, marginBottom: -20
                }} direction="row" spacing={1} alignItems="center">
                    <Grid item>
                        <Typography variant="body2" style={{
                            color: MulwiColors.subtitleTypography
                        }}>AUTHOR TIMESTAMP</Typography>
                    </Grid>
                    <Grid item>
                        <Typography style={{
                            marginBottom: 3
                        }}>CONTENT</Typography>
                    </Grid>
                </Grid> 
                <Grid container style={{
                    marginLeft: 10, padding: 0, marginBottom: -20
                }} direction="row" spacing={1} alignItems="center">
                    <Grid item>
                        <Typography variant="body2" style={{
                            color: MulwiColors.subtitleTypography
                        }}>AUTHOR TIMESTAMP</Typography>
                    </Grid>
                    <Grid item>
                        <Typography style={{
                            marginBottom: 3
                        }}>CONTENT</Typography>
                    </Grid>
                </Grid>  */}

            </Grid>
        </Grid>
            <Grid item>
                <form onSubmit={sendMsg}>

                    <FormControl fullWidth>
                        <Input
                            ref={inputbox}
                            inputProps={{ maxLength: 128 }}
                            placeholder="Message"
                            value={inputMsg}
                            onChange={e => setInputMsg(e.target.value)}
                            style={{
                                backgroundColor: "white",
                                marginLeft: 10,
                                marginRight: 10,
                            }}
                            minRows={2}
                            maxRows={2}
                            variant="outlined"
                            endAdornment={
                                <InputAdornment position="end">
                                    <IconButton
                                        aria-label="toggle password visibility"
                                        onClick={(e) => {
                                            sendMsg(e)
                                        }}
                                        style={{
                                            color: MulwiColors.blueDark
                                        }}
                                    >
                                        <SendIcon />
                                    </IconButton>
                                </InputAdornment>
                            }
                        />
                    </FormControl>
                </form>
            </Grid>
        </div>
    </React.Fragment>)
}