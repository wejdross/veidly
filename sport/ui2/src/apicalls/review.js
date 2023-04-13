import {API_URL} from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";

export function getUserReview(rsvID) {
    return xhr(
        API_URL + "/api/review/user?rsv_id=" + (rsvID || ""),
        "GET",
        gettoken()
    )
}

export function postUserReview(c) {
    return xhr(
        API_URL + "/api/review",
        "POST",
        gettoken(),
        JSON.stringify(c)
    )
}

export function getPubReviews(trainingID) {
    return xhr(API_URL + "/api/review/pub?training_id=" + trainingID, "GET")
}
