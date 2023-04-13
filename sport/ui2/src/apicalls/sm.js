import { API_URL } from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";
import { processTrainings } from './instructor.api'

export function postSm(req) {
    return xhr(
        API_URL + "/api/sub/model",
        "POST",
        gettoken(),
        JSON.stringify(req)
    )
}

export function getSm() {
    return xhr(
        API_URL + "/api/sub/model",
        "GET",
        gettoken(),
    )
}

export function getSmForInstr(iid) {
    return xhr(
        API_URL + "/api/sub/model",
        "GET", null, null, null, null,
        ["instructor_id"],
        [iid]
    )
}


export function patchSm(req) {
    return xhr(
        API_URL + "/api/sub/model",
        "PATCH",
        gettoken(),
        JSON.stringify(req)
    )
}

export function deleteSm(id) {
    return xhr(
        API_URL + "/api/sub/model",
        "DELETE",
        gettoken(),
        JSON.stringify({
            ID: id
        })
    )
}


export function postSmBinding(smid, tid) {
    return xhr(
        API_URL + "/api/sub/model/binding",
        "POST",
        gettoken(),
        JSON.stringify({
            SubModelID: smid,
            TrainingID: tid
        })
    )
}

export function deleteSmBinding(smid, tid) {
    return xhr(
        API_URL + "/api/sub/model/binding",
        "DELETE",
        gettoken(),
        JSON.stringify({
            SubModelID: smid,
            TrainingID: tid
        })
    )
}

export async function getTrainingsForSm(smID, instrID) {
    
    let x = await xhr(
        API_URL + "/api/training",
        "GET",
        null,
        null, null, null,
        ["sm_id", "instructor_id"],
        [smID, instrID || ""]
    )
    return processTrainings(x)
}

export function postSub(smID) {
    return xhr(
        API_URL + "/api/sub",
        "POST",
        gettoken(),
        JSON.stringify({ SubModelID: smID })
    )
}

export function getUserSub(id) {
    return xhr(
        API_URL + "/api/sub/user",
        "GET",
        gettoken(),
        null, null, null,
        ["id"],
        [id || ""]
    )
}

export function getInstrSub(id) {
    return xhr(
        API_URL + "/api/sub/instructor",
        "GET",
        gettoken(),
        null, null, null,
        ["id"],
        [id || ""]
    )
}

export function postUserRefund(subID, at) {
    return xhr(
        API_URL + "/api/sub/refund/user",
        "POST",
        gettoken(),
        JSON.stringify({
            SubID: subID
        }))
}

export function postExpire(subID, at) {
    return xhr(
        API_URL + "/api/sub/expire",
        "POST",
        gettoken(),
        JSON.stringify({
            SubID: subID
        }))
}

export function postUserDispute(subID, at, email, msg) {
    return xhr(
        API_URL + "/api/sub/dispute",
        "POST",
        gettoken(),
        JSON.stringify({
            SubID: subID,
            Email: email,
            msg
        }))
}


export async function createQr(subID) {
    return xhr(
        API_URL + "/api/qr/sub",
        "POST",
        gettoken(),
        JSON.stringify({
            SubID: subID,
            DataUrl: true,
            Size: 256
        })
    )
}

export async function evalQr(id) {
    return xhr(
        API_URL + "/api/qr/sub/eval?id=" + id,
        "GET",
        gettoken()
    )
}

export async function confirmQr(id) {
    return xhr(
        API_URL + "/api/qr/sub/confirm?id=" + id,
        "GET",
        gettoken()
    )
}
