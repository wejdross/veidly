import { API_URL } from "../conf";
import { gettoken } from "../helpers";
import { xhr } from "./api";

export function getPricingInfo(req) {
    return xhr(
        API_URL + "/api/rsv/pricing",
        "POST",
        null,
        JSON.stringify(req)
    )
}