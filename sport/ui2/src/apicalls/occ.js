import {API_URL} from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";

export function putOcc(req) {
    return xhr(
        API_URL + "/api/training/occ",
        "PUT",
        gettoken(),
        JSON.stringify(req)
    )
}
