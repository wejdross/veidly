const { API_URL } = require("../conf");
const { xhr } = require("./api");

export function GETArticles() {
    return xhr(API_URL + "/api/articles", "GET", null, null, null, null, null);
}

export function GETArticle(id) {
    return xhr(API_URL + "/api/articles/" + id, "GET", null, null, null, null, null);
}

export function GenerateArticlesNumbers(arrLength) {
    var randoms = new Set()
    while (randoms.size < 3) {
        const num = Math.floor(Math.random() * arrLength);
        if (num < 0) {
            continue
        }
        randoms.add(num);            
    }
    const rnds = Array.from(randoms)
    return rnds
}

export function GETInstructors() {
    return xhr(API_URL + "/api/instructor/top", "GET", null, null, null, null, null);
}