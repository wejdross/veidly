import { API_URL } from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";

export function donate(v, email) {
    return xhr(
        API_URL + "/api/donate?amount=" + v + "&email=" + escape(email),
        "POST"
    )
}