import { API_URL } from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";
import { processTrainings } from "./instructor.api";


export function postGroup(req) {
    return xhr(
        API_URL + "/api/training/group",
        "POST",
        gettoken(),
        JSON.stringify(req)
    )
}

export function patchGroup(req) {
    return xhr(
        API_URL + "/api/training/group",
        "PATCH",
        gettoken(),
        JSON.stringify(req)
    )
}

export async function getTrainingsForGroups(ids) {
    let t = await xhr(
        API_URL + "/api/training",
        "GET",
        gettoken(),
        null, null, null,
        ["group_ids"],
        [JSON.stringify(ids)]
    )
    return processTrainings(t)
}

export function deleteGroup(id) {
    return xhr(
        API_URL + "/api/training/group",
        "DELETE",
        gettoken(),
        JSON.stringify({
            ID: id
        })
    )
}

export function getGroups() {
    return xhr(
        API_URL + "/api/training/group",
        "GET",
        gettoken(),
    )
}

export function putGroupBinding(grpID, tID) {
    return xhr(
        API_URL + "/api/training/group/binding",
        "PUT",
        gettoken(),
        JSON.stringify({
            GroupID: grpID,
            TrainingID: tID
        })
    )
}

export function deleteGroupBinding(grpID, tID) {
    return xhr(
        API_URL + "/api/training/group/binding",
        "DELETE",
        gettoken(),
        JSON.stringify({
            GroupID: grpID,
            TrainingID: tID
        })
    )
}
