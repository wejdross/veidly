const { API_URL } = require("../conf");
const { xhr } = require("./api");

export function getAllTags(input) {
    return xhr(
        API_URL + "/api/tag?page=0&size=5&query=" + escape(input),
        "GET");
}
