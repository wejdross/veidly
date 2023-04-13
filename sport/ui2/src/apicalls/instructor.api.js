import {API_URL, G_API_GEODECODE, G_API_KEY} from "../conf";
import { dateToEpoch, gettoken } from "../helpers";
import { xhr } from "./api";

export function deleteUser() {
    return xhr(
        API_URL + "/api/user",
        "DELETE",
        gettoken(),
        // JSON.stringify({
        //     Password: pass
        // })
    )
}

export function deleteInstructor() {
    return xhr(
        API_URL + "/api/instructor",
        "DELETE",
        gettoken(),
        // JSON.stringify({
        //     Password: pass
        // })
    )
}

export function canDeleteInstructor() {
    return xhr(
        API_URL + "/api/instructor/can_delete",
        "GET",
        gettoken()
    )
}

export function postDecision(rsvID, decision) {
    return xhr(
        API_URL + "/api/rsv/decision",
        "POST",
        gettoken(),
        JSON.stringify({
            ReservationID: rsvID,
            Decision: decision
        }))
}

export function postInstructorRefund(rsvID) {
    return xhr(
        API_URL + "/api/rsv/refund/instructor",
        "POST",
        gettoken(),
        JSON.stringify({
            ReservationID: rsvID
        }))
}

export function postUserRefund(rsvID, at) {
    return xhr(
        API_URL + "/api/rsv/refund/user",
        "POST",
        gettoken(),
        JSON.stringify({
            ReservationID: rsvID,
            AccessToken: at || null
        }))
}

export function postRsvExpire(rsvID, at) {
    return xhr(
        API_URL + "/api/rsv/expire",
        "POST",
        gettoken(),
        JSON.stringify({
            ReservationID: rsvID,
            AccessToken: at || null
        }))
}

export function postUserCancel(rsvID, at) {
    return xhr(
        API_URL + "/api/rsv/cancel",
        "POST",
        gettoken(),
        JSON.stringify({
            ReservationID: rsvID,
            AccessToken: at || null
        }))
}

export function postUserDispute(rsvID, at, email, msg) {
    return xhr(
        API_URL + "/api/rsv/dispute/user",
        "POST",
        gettoken(),
        JSON.stringify({
            ReservationID: rsvID,
            AccessToken: at || null,
            Email: email,
            msg
        }))
}

export function getInstructor(iid) {
    if(iid) {
        return xhr(
            API_URL + "/api/instructor?instructor_id=" + iid,
            "GET")
    } else {
        return xhr(
            API_URL + "/api/instructor",
            "GET",
            gettoken(),
            )
    }
}

export function processTrainings(x) {
    x = JSON.parse(x)
    if(!x) return []
    for(let j = 0; j < x.length; j++) {
        if(x[j].Training) {
            x[j].Training.DateStart = new Date(x[j].Training.DateStart)
            x[j].Training.DateEnd = new Date(x[j].Training.DateEnd)
        }
        if(x[j].Occurrences) {
            for(let i = 0; i < x[j].Occurrences.length; i++) {
                x[j].Occurrences[i].DateStart = new Date(x[j].Occurrences[i].DateStart)
                x[j].Occurrences[i].DateEnd = new Date(x[j].Occurrences[i].DateEnd)
            }
        }
    }
    return x
}

export async function getTrainingByID(id) {
    let x = await xhr(
        API_URL + "/api/training?id=" + id,
        "GET")
    return processTrainings(x)
}

export async function searchTrainings(q) {
    if(q) {
        q = "?q=" + q
    }
    let x = await xhr(
        API_URL + "/api/training" + q,
        "GET",
        gettoken()
    )
    return processTrainings(x)
}

export async function searchPubTrainings(q, instructorID) {
    let x = await xhr(
        API_URL + "/api/training",
        "GET",
        null,
        null, null, null,
        ["q", "instructor_id"],
        [q, instructorID || ""]
    )
    return processTrainings(x)
}

export async function getTrainings(id) {
    let ids = ""
    if(id) {
        ids = "?id=" + id
    }
    let x = await xhr(
        API_URL + "/api/training" + ids,
        "GET",
        gettoken()
    )
    return processTrainings(x)
}


export function apiCreateTraining(req) {
    return xhr(
        API_URL + "/api/training",
        "POST",
        gettoken(),
        JSON.stringify(req)
    )
}

export function patchTraining(req) {
    return xhr(
        API_URL + "/api/training",
        "PATCH",
        gettoken(),
        JSON.stringify(req)
    )
}

export function apiDeleteTraining(id) {
    return xhr(
        API_URL + "/api/training",
        "DELETE",
        gettoken(),
        JSON.stringify({
            id: id
        })
    )
}

export function createInstructor(request) {
    return xhr(
        API_URL + "/api/instructor",
        "POST",
        gettoken(),
        JSON.stringify(request)
    );
}

export function PATCHInstructor(req) {
    return xhr(
        API_URL + "/api/instructor",
        "PATCH",
        gettoken(),
        JSON.stringify(req)
    );
}

export function GeoDecode(userInput) {
    return xhr(
        G_API_GEODECODE,
        "GET",
        null,
        null,
        null,
        null,
        ["address", "key"],
        [userInput, G_API_KEY]
    )
}

// raw api call
export function _getSchedule(start, end, trainingID) {
    return xhr(
        API_URL + "/api/schedule",
        "GET",
        gettoken(),
        null,
        null, null,
        ["start", "end", "training_id"],
        [dateToEpoch(start), dateToEpoch(end), trainingID || ""]
    )
}

export function _getRsvSchedule(start, end, trainingID) {
    return xhr(
        API_URL + "/api/schedule/rsv/t/user",
        "GET",
        gettoken(),
        null, null, null, 
        ["start", "end", "training_id"],
        [dateToEpoch(start), dateToEpoch(end), trainingID || ""]
    )
}

export function _getUserSchedule(start, end, instructorID, trainingID, smID) {
    return xhr(
        API_URL + "/api/schedule",
        "GET",
        null,
        null,
        null, null,
        ["start", "end", "instructor_id", "training_id", "smID"],
        [dateToEpoch(start), dateToEpoch(end), instructorID, trainingID || "", smID || ""]
    )
}

function processSchedule(x) {
    if(!x) return []
    for(let j = 0; j < x.length; j++) {
        if(x[j].Training) {
            x[j].Training.DateStart = new Date(x[j].Training.DateStart)
            x[j].Training.DateEnd = new Date(x[j].Training.DateEnd)
        }
        if(x[j].Occs) {
            for(let i = 0; i < x[j].Occs.length; i++) {
                x[j].Occs[i].DateStart = new Date(x[j].Occs[i].DateStart)
                x[j].Occs[i].DateEnd = new Date(x[j].Occs[i].DateEnd)
            }
        }
        if(x[j].Schedule) {
            for(let i = 0; i < x[j].Schedule.length; i++) {
                x[j].Schedule[i].Start = new Date(x[j].Schedule[i].Start)
                x[j].Schedule[i].End = new Date(x[j].Schedule[i].End)
                if(x[j].Schedule[i].Occ) {
                    x[j].Schedule[i].Occ.DateStart = new Date(x[j].Schedule[i].Occ.DateStart)
                    x[j].Schedule[i].Occ.DateEnd = new Date(x[j].Schedule[i].Occ.DateEnd)
                }
            }
        }
    }
    return x
}
// this will call _getSchedule and then process result. for example will parse dates into js objects
export async function getRsvSchedule(start, end, trainingID) {
    let x = await _getRsvSchedule(start, end, trainingID)
    x = JSON.parse(x)
    return processSchedule(x)
}

// this will call _getSchedule and then process result. for example will parse dates into js objects
export async function getSchedule(start, end, trainingID) {
    let x = await _getSchedule(start, end, trainingID)
    x = JSON.parse(x)
    return processSchedule(x)
}

// this will call _getUserSchedule and then process result. for example will parse dates into js objects
export async function getUserSchedule(start, end, instructorID, trainingID, smID) {
    let x = await _getUserSchedule(start, end, instructorID, trainingID, smID)
    x = JSON.parse(x)
    return processSchedule(x)
}


export function search(ApiSearchRequest) {
    //const ApiSearchRequest = {
    //    Query: query,
    //        Country: "PL",
    //        Lang: "pl",
    //        PriceMax: 90,
    //        Diffs: [
    //            3,1,2
    //        ],
    //        Sort: [
    //            {
    //                Column: "price",
    //                IsDesc: true,
    //            },
    //        ]
    //}
    return xhr(
        API_URL + "/api/search",
        "POST",
        null,
        JSON.stringify(ApiSearchRequest)
        //JSON.stringify({
        //    Query: query,
        //    Country: "PL",
        //    Lang: "pl",
        //    PriceMax: 90,
        //    Diffs: [
        //        3,1,2
        //    ],
        //    Sort: [
        //        {
        //            Column: "price",
        //            IsDesc: true,
        //        },
        //    ] 
        //  })
    )//
}

export function postReservation(rr) {
    return xhr(
        API_URL + "/api/rsv",
        "POST",
        gettoken(),
        JSON.stringify(rr)
    )
}

export function readRsvByToken(token) {
    return xhr(
        API_URL + "/api/rsv/t/token?access_token=" + token,
        "GET",
    )
}

export function getUserRsvs() {
    return xhr(
        API_URL + "/api/rsv/t/user",
        "GET",
        gettoken()
    )
}

export function getInstructorRsvs() {
    return xhr(
        API_URL + "/api/rsv/t/instructor",
        "GET",
        gettoken()
    )
}


export function getUserRsvByID(id) {
    return xhr(
        API_URL + "/api/rsv/t/user",
        "GET",
        gettoken(),
        null,
        null,
        null,
        ["id"],
        [id]
    )
}

export function getInstrRsvByID(id) {
    return xhr(
        API_URL + "/api/rsv/t/instructor",
        "GET",
        gettoken(),
        null,
        null,
        null,
        ["id"],
        [id]
    )
}

export function getUserRsvByAccessToken(at) {
    return xhr(
        API_URL + "/api/rsv/t/token",
        "GET",
        null,
        null,
        null,
        null,
        ["access_token"],
        [at]
    )
}

export function updatePayoutData(objPayoutData) {
    return xhr(
        API_URL + "/api/instructor/payout",
        "PATCH",
        gettoken(),
        JSON.stringify(objPayoutData)
    )
}

export function getPayoutData() {
    return xhr(
        API_URL + "/api/instructor/payout",
        "GET",
        gettoken(),
    )
}

export function deletePayoutData() {
    return xhr(
        API_URL + "/api/instructor/payout",
        "DELETE",
        gettoken(),
    )
}
export function postImg(fd) {
    return xhr(
        API_URL + "/api/training/img",
        "POST",
        gettoken(),
        fd
    )
}

export function deleteImg(id, trainingID) {
    return xhr(
        API_URL + "/api/training/img",
        "DELETE",
        gettoken(),
        JSON.stringify({
            ID: id,
            TrainingID: trainingID
        })
    )
}

export function postVacation(req) {
    return xhr(
        API_URL + "/api/instructor/vacation",
        "POST",
        gettoken(),
        JSON.stringify(req)
    )
}

export function deleteVacation(id) {
    return xhr(
        API_URL + "/api/instructor/vacation",
        "DELETE",
        gettoken(),
        JSON.stringify({
            ID: id
        })
    )
}

export async function getVacations(instrID) {
    let d = await xhr(
        API_URL + "/api/instructor/vacation",
        "GET",
        instrID ? null : gettoken(),
        null, null, null,
        ["instructor_id"],
        [instrID || ""]
    )
    d = JSON.parse(d)
    for(let i = 0 ; i < d.length; i++) {
        d[i].DateStart = new Date(d[i].DateStart)
        d[i].DateEnd = new Date(d[i].DateEnd)
    }
    return d
}

export async function getInstrContact(instructorID, accessToken) {
    return xhr(
        API_URL + "/api/rsv/instr/contact",
        "GET",
        accessToken ? null : gettoken(),
        null,
        null,
        null,
        ["instructor_id", "access_token"],
        [instructorID, accessToken || ""]
    )
}

export async function createQr(rsvID, accessToken) {
    return xhr(
        API_URL + "/api/qr/rsv",
        "POST",
        accessToken ? null : gettoken(),
        JSON.stringify({
            RsvID: rsvID,
            AccessToken: accessToken || null,
            DataUrl: true,
            Size: 256
        })
    )
}

export async function evalQr(id) {
    return xhr(
        API_URL + "/api/qr/rsv/eval?id=" + id,
        "GET",
        gettoken()
    )
}

export async function postProfileImg(fd, isPrimary) {
    return xhr(
        API_URL + "/api/instructor/profile/img" + (isPrimary ? "?primary=1" : ""),
        "POST",
        gettoken(),
        fd
    )
}

export async function deleteProfileImg(path, isPrimary) {
    return xhr(
        API_URL + "/api/instructor/profile/img" + (isPrimary ? "?primary=1" : ""),
        "DELETE",
        gettoken(),
        JSON.stringify({
            Path: path
        })
    )
}

export async function getInvoices(id) {
    return xhr(
        API_URL + "/api/invoice?id=" + (id || ""),
        "GET",
        gettoken()
    )
}

