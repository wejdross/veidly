/*
    url - wiadomo
    method - wiadomo
    token - token JWT jeśli istnieje
    body - wiadomo
    hk - header key
    hv - header value
    qk - querry key
    qv - querry value
*/

export function xhr(url, method, token, body, hk, hv, qk, qv) {

    return new Promise((res, rej) => {

        let xhr = new XMLHttpRequest();
        // dont override response content type with binary
        //xhr.overrideMimeType('text/plain; charset=x-user-defined');

        if (qk && qk.length > 0) {
            if (!url.endsWith("?")) {
                url += "?";
            }
            let isFirst = true;
            for (let i = 0; i < qk.length; i++) {

                if (!isFirst) {
                    url += "&";
                }

                url += encodeURIComponent(qk[i]) + "=" + encodeURIComponent(qv[i]);

                if (isFirst)
                    isFirst = false;
            }
        }

        xhr.open(method, url);

        if (token)
            xhr.setRequestHeader("Authorization", "Bearer " + token);

        if (hk) {
            for (let i = 0; i < hk.length; i++) {
                xhr.setRequestHeader(hk[i], hv[i])
            }
        }

        xhr.onload = function (e) {

            let r = xhr.response;

            switch (xhr.status) {
                case 200:
                    res(r);
                    break;
                case 204:
                case 202:
                case 304:
                    res(xhr.status);
                    break;
                case 404:
                case 500:
                case 422:
                case 400:
                case 401:
                case 403:
                default:
                    rej(xhr.status);
                    break;
            }
            return;
        };

        xhr.send(body);
    });
}
