import { API_URL } from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";
import { processTrainings } from "./instructor.api";

export function postDc(req) {
    return xhr(
        API_URL + "/api/dc",
        "POST",
        gettoken(),
        JSON.stringify(req)
    )
}

export function getDc() {
    return xhr(
        API_URL + "/api/dc",
        "GET",
        gettoken(),
    )
}

export function patchDc(req) {
    return xhr(
        API_URL + "/api/dc",
        "PATCH",
        gettoken(),
        JSON.stringify(req)
    )
}

export function deleteDc(id) {
    return xhr(
        API_URL + "/api/dc",
        "DELETE",
        gettoken(),
        JSON.stringify({
            ID: id
        })
    )
}

export function postDcBinding(dcID, tid) {
    return xhr(
        API_URL + "/api/dc/binding",
        "POST",
        gettoken(),
        JSON.stringify({
            DcID: dcID,
            TrainingID: tid
        })
    )
}

export function deleteDcBinding(dcID, tid) {
    return xhr(
        API_URL + "/api/dc/binding",
        "DELETE",
        gettoken(),
        JSON.stringify({
            DcID: dcID,
            TrainingID: tid
        })
    )
}


export async function getTrainingsForDc(id) {
    let x = await xhr(
        API_URL + "/api/training",
        "GET",
        gettoken(),
        null, null, null,
        ["dc_id"],
        [id]
    )
    return processTrainings(x)
}

export function redeemDc(name, trainingID) {
    return xhr(
        API_URL + "/api/dc/redeem",
        "GET",
        null,
        null, null, null,
        ["name", "training_id"],
        [name, trainingID]
    )
}
