import { useEffect, useRef, useState } from "react";
import { getChatToken, newNotifyWsConn } from "../apicalls/chat";
import { closeWsAndWait } from "./chatroom";

const LMKey = "LastMessageCache"

export function updateLastMessageCache(chatroomID, timestamp) {
    timestamp = Number(timestamp)
    let cache = localStorage.getItem(LMKey)
    if(!cache)
        cache = {}
    else 
        cache = JSON.parse(cache)
    cache[chatroomID] = timestamp
    localStorage.setItem(LMKey, JSON.stringify(cache))
}

export function getLastMessageCache() {
    return localStorage.getItem(LMKey)
}

export function hasUnreadMessages(not) {
    for(let k in not) {
        if(not[k] > 0)
            return true
    }
    return false
}

export function ChatNotifications(props) {
    
    const [hasMessages, setHasMessages] = useState(false)
    const nots = useRef()

    const conn = useRef()

    async function closeConnIfExists() {
        await closeWsAndWait(conn.current)
        nots.current = null
        conn.current = null
    }

    function rm(chatRoomID) {
        if(!chatRoomID || !nots.current)
            return
        let cpy = {...nots.current}
        cpy[chatRoomID] = 0
        nots.current = cpy 
        props.onNotification && props.onNotification(nots.current)
    }

    function handleNotifyWsMsg(e) {
        let d = JSON.parse(e.data)
        let cc = JSON.parse(getLastMessageCache())
        if(!cc)
            cc = {} 
        for(let k in d) {
            if(!cc[k]) {
                continue
            }
            if(d[k].Count > 0 && d[k].LastTimestamp <= cc[k]) {
                d[k].Count = 0
            }
        }
        if(!nots.current) {
            nots.current = {}
            nots.current.rm = rm
        }
        for(let k in d) {
            if(!nots.current[k])
                nots.current[k] = d[k].Count
            else
                nots.current[k] += d[k].Count
        }
        // console.log(nots.current)
        props.onNotification && props.onNotification({...nots.current})
        // console.log("nots", nots.current)
        setHasMessages(hasUnreadMessages(nots.current))
    }

    function createWs() {
        let ws = newNotifyWsConn()
        nots.current = null
        ws.addEventListener("open", () => {
            ws.send(getLastMessageCache())
            console.log("notification ws opened")
        })
        ws.addEventListener("close", () => {
            console.log("notification ws closed")
        })
        ws.addEventListener("error", (e) => {
            console.log("notification ws closed with error", e)
        })
        ws.addEventListener("message", e => {
            try {
                handleNotifyWsMsg(e)
            } catch (ex) {
                console.log("handleNotifyWsMsg err", ex)
            }
        })
        
        return ws
    }

    function resetWs() {
        try {
            let ws = createWs()
            conn.current = ws
        } catch(ex) {
            console.log(ex)
        }
    }

    const interval = useRef()
    const lastToken = useRef()

    async function testConnAndReset() {
        let t = getChatToken()
        if(!t)
            return
        if(t === lastToken.current && conn.current && conn.current.readyState === 1) {
            // console.log("already established")
            return
        }
        lastToken.current = t
        await closeConnIfExists()
        resetWs()
    }

    async function startConnInterval() {

        if(interval.current)
            clearInterval(interval.current)

        await testConnAndReset()

        interval.current = setInterval(async () => {
            await testConnAndReset()
        }, 10000)
        
    }

    useEffect(() => {
        if(!props.chatToken)
            return
        startConnInterval()
    }, [props.chatToken])

    useEffect(() => {
        return () => {
            if(interval.current)
                clearInterval(interval.current)
            closeConnIfExists()
        }
    }, [])

    return null

    // if(!props.chatToken)
    //     return null
    // return (<Snackbar open={hasMessages}
    //     onClose={() => {
    //         setHasMessages(false)
    //     }}>
    //         <Alert severity="info">
    //             <strong>You have unread messages, dude</strong>
    //         </Alert>
    //     </Snackbar>)
}
