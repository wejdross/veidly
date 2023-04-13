import { API_URL } from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";

let ct = null

export function removeChatToken() {
    localStorage.removeItem("chatToken")
    ct = null
}

export function storeChatToken(token, persist) {
    if(persist) {
        ct = null
        localStorage.setItem("chatToken", token)
    } else {
        ct = token
        localStorage.setItem("chatToken", null)
    }
}

export function getChatToken() {
    return ct || localStorage.getItem("chatToken")
}

export async function setChatToken(persist) {
    let res = await xhr(
        API_URL + "/api/chat_integrator/token", 
        "GET", 
        gettoken())
    res = JSON.parse(res).Token
    if(!res)
        throw "no token"
    storeChatToken(res, persist)
}


export function apiValidateToken() {
    return xhr(API_URL + "/api/chat/token/validate", "GET", getChatToken())
}

export function postChatroom(req) {
    return xhr(
        API_URL + "/api/chat/room", 
        "POST", 
        getChatToken(),
        JSON.stringify(req)) 
}

export function getChatroomsForUser() {
    return xhr(API_URL + "/api/chat/room", "GET", getChatToken()) 
}

const chatCluster = [
    "123"
]

export function postWsToken(req) {
    return xhr(API_URL + "/api/chat/token/ws", "POST", getChatToken(), JSON.stringify(req)) 
}

export function newWsConn(wst) {
    let randIx = 0
    let _url = API_URL.replace("https", "wss").replace("http", "ws")
    _url += "/api/chat/open/" + chatCluster[randIx] + "?t=" + wst
    return new WebSocket(_url)
}

export function newNotifyWsConn() {
    let _url = API_URL.replace("https", "wss").replace("http", "ws")
    _url += "/api/chat/notify/open?t=" + getChatToken()
    return new WebSocket(_url)
}

export function getChatAccessTokens(chatroomid) {
    return xhr(API_URL + "/api/chat/room/access_token?chatRoomID=" + chatroomid, 
        "GET", getChatToken()) 
}

export function postChatAccessToken(req) {
    return xhr(
        API_URL + "/api/chat/room/access_token", 
        "POST", 
         getChatToken(),
        JSON.stringify(req)) 
}

export function postJoinChatroom(req, anon) {
    return xhr(
        API_URL + "/api/chat/room/join", 
        "POST", 
        anon ? null : getChatToken(),
        JSON.stringify(req)) 
}

export function postUserRoom(req) {
    return xhr(
        API_URL + "/api/chat_integrator/room", 
        "POST", 
        gettoken(),
        JSON.stringify(req)) 
}
