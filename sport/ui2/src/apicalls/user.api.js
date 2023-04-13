import { getSupportedLanguage } from "../locale";


const { API_URL, G_API_KEY } = require("../conf");
const { gettoken } = require("../helpers");
const { xhr } = require("./api");

export function validatPass(pass) {
    return xhr(
        API_URL + "/api/user/password/validate?password=" + (escape(pass) || ""),
        "GET",
        gettoken()
    )
}

export function userStat() {
    return xhr(
        API_URL + "/api/user/stat",
        "GET",
        gettoken(),
        null, null, null, null);
}

export function putImg(fd) {
    return xhr(
        API_URL + "/api/user/avatar",
        "PUT",
        gettoken(),
        fd
    )
}

export function userLogin(ur) {
    return xhr(
        API_URL + "/api/user/login",
        "POST",
        null,
        JSON.stringify(ur),
        null, true)
}

export function userRegister(userRequest, mfaReturnUrl) {
    let x = "";
    if (mfaReturnUrl) {
        x = escape(mfaReturnUrl);
    }
    return xhr(
        API_URL + "/api/user/register?mfa_return_url=" + x,
        "POST",
        null,
        JSON.stringify(userRequest))
}

export function forgotPassword(email) {
    return xhr(
        API_URL + "/api/user/password/forgot",
        "GET",
        null,
        null,
        null,
        null,
        ["email", "lang"],
        [email, getSupportedLanguage()]);
}

export function resetPassword(token, password) {
    return xhr(
        API_URL + "/api/user/password/reset",
        "POST",
        null,
        JSON.stringify({
            token: token,
            password: password
        }));
}

export function resendRegisterEmail(email, mfaReturnUrl) {
    let x = "";
    if (mfaReturnUrl) {
        x = escape(mfaReturnUrl);
    }
    return xhr(
        API_URL + "/api/user/register/resend?mfa_return_url=" + x,
        "POST",
        null,
        JSON.stringify({
            email: email
        }))
}

export function getOauthGoogleUrl(return_url) {
    let x = "";
    if (return_url) {
        x = escape(return_url)
    }
    return xhr(
        API_URL + "/api/user/oauth/url/google?mode=raw&return_url=" + x,
        "GET")
}


export function oauthActionWithCode(provider, code) {
    return xhr(
        API_URL + "/api/user/oauth/code/" + provider,
        "POST",
        null,
        JSON.stringify({
            code: code
        }));
}

export function patchUserData(objTosend) {
    return xhr(
        API_URL + "/api/user/data",
        "PATCH",
        gettoken(),
        JSON.stringify(objTosend),
        null,
        null,
        null,
        null
    );
}

export function patchUserContactData(objTosend) {
    return xhr(
        API_URL + "/api/user/contact",
        "PATCH",
        gettoken(),
        JSON.stringify(objTosend),
    );
}

export function patchUserPassword(objTosend) {
    return xhr(
        API_URL + "/api/user/password",
        "PATCH",
        gettoken(),
        JSON.stringify(objTosend),
        null,
        null,
        null,
        null
    );
}

export function deleteUser() {
    return xhr(
        API_URL + "/api/user",
        "DELETE",
        gettoken(),
        null,
        null,
        null,
        null,
        null
    );
}

export function uploadPhoto(data) {
    return xhr(
        API_URL + "/api/user/photo",
        "PUT",
        gettoken(),
        data,
        null,
        null,
        null,
        null
    );
}

export function getUserInfo() {
    return xhr(API_URL + "/api/user", "GET", gettoken(), null, null, null, null, null)
}

export async function fetchGeo(input_string, simple) {

    // let x = await xhr("https://maps.googleapis.com/maps/api/place/textsearch/json",
    // "GET", null, null, null, null,
    // ["query", "key", "language"],
    // [
    //     input_string, 
    //     G_API_KEY, 
    //     getSupportedLanguage()
    // ])
    // console.log(x)

    let request = {
        query: input_string,
        fields: ['name', 'geometry']
    };

    /*
     * this is pure fucking atrocity
     * but im done fighting it
     * google wins
     * dont try to resist
     * no more
    */
    return new Promise((res, rej) => {
        const service = new window.google.maps.places.PlacesService(document.createElement('div'))
        service.findPlaceFromQuery(request, function (results, status) {
            if (status !== window.google.maps.places.PlacesServiceStatus.OK) {
                return rej(status)
            }

            let arr = []
            for(let i  =0; i < results.length; i++) {
                arr.push({
                    display_name: results[i].name,
                    lat: results[i].geometry.location.lat(),
                    lon: results[i].geometry.location.lng()
                })
            }

            return res(JSON.stringify(arr))
        });
    })


    // let returnArray = []
    // let lang = "pl"
    // //let lang = getSupportedLanguage()
    // const results = await xhr(
    //     "https://git.mulwi.cloud/geo/search?q=" + encodeURI(input_string) + "&accept-language=" + lang + "&addressdetails=1",
    //     "GET",
    //     null,
    //     null,
    //     "Access-Control-Allow-Origin",
    //     "*",
    //     null,
    //     null
    // )
    // let dc = {}
    // JSON.parse(results).map((val, idx) => {
    //     let tmp = {
    //         display_name: '',
    //         lat: val.lat,
    //         lon: val.lon,
    //     }
    //     //let currentlyInResults = false

    //     if (simple) {
    //         if (val.address.city) {
    //             tmp.display_name = `${val.address.city}`
    //         } else if (val.address.town) {
    //             tmp.display_name = `${val.address.town}`
    //         }
    //     } else {
    //         // if there is a street
    //         if (val.address.road && val.address.house_number && val.address.town) {
    //             tmp.display_name += `${val.address.road} ${val.address.house_number}, ${val.address.town}`
    //         } else if (val.address.road && val.address.house_number && val.address.city) {
    //             tmp.display_name += `${val.address.road} ${val.address.house_number}, ${val.address.city}`
    //         } else if (val.address.road && val.address.town) {
    //             tmp.display_name += `${val.address.road}, ${val.address.town}`
    //         } else if (val.address.road && val.address.city) {
    //             tmp.display_name += `${val.address.road}, ${val.address.city}`
    //         } else if (val.address.city) {
    //             tmp.display_name = `${val.address.city}`
    //         } else if (val.address.town) {
    //             tmp.display_name = `${val.address.town}`
    //         } else if (val.address.village) {
    //             tmp.display_name = val.address.village
    //         }

    //     }

    //     if (!tmp.display_name || dc[tmp.display_name]) return

    //     dc[tmp.display_name] = 1

    //     // returnArray.map((val, idx) => {
    //     //     if (val.display_name === tmp.display_name) {
    //     //         currentlyInResults = true
    //     //     }
    //     // })
    //     // if (!currentlyInResults && tmp.display_name) {
    //     //     returnArray.push(tmp)
    //     // }
    //     returnArray.push(tmp)
    // })
    // return JSON.stringify(returnArray)
}

export async function reverseFetchGeo(lat, lng) {
    
    const service = new window.google.maps.Geocoder()
    let r = await service.geocode({
        location: {
            lat: lat,
            lng: lng
        }
    })

    for(let i = 0; i < r.results.length; i++) {
        let record = r.results[i]
        for(let j = 0; j < record.address_components.length; j++) {
            let ac = record.address_components[j]
            for(let z = 0; z < ac.types.length; z++) {
                if(ac.types[z] == 'locality')
                    return JSON.stringify({
                        address: {
                            city: ac.long_name
                        }
                    })
            }
        }
    }

    return null
    
    // return xhr(
    //     "https://git.mulwi.cloud/geo/reverse?format=json&lat=" + lat + "&lon=" + lng,
    //     "GET",
    //     null,
    //     null,
    //     "Access-Control-Allow-Origin",
    //     "*",
    //     null,
    //     null
    // )
}

export function fetchCharts() {
    return xhr(
        API_URL + "/api/charts",
        "GET",
        gettoken()
    )
}

export function getLangs(q) {
    return xhr(
        API_URL + "/api/lang?q=" + q + "&l=" + getSupportedLanguage(),
        "GET")
}

export function explainLangs(l) {
    return xhr(
        API_URL + "/api/lang/explain?t=" + getSupportedLanguage() + "&l=" + JSON.stringify(l),
        "GET")
}


export function getTags(q) {
    return xhr(
        API_URL + "/api/tag?q=" + q + "&l=" + getSupportedLanguage(),
        "GET")
}

export function explainTags(tags) {
    return xhr(
        API_URL + "/api/tag/explain?t=" + getSupportedLanguage() + "&s=" + JSON.stringify(tags),
        "GET")
}

export function googlePlaceIDtoJSON(place_id) {
    return xhr(
        `https://maps.googleapis.com/maps/api/geocode/json?place_id=${place_id}&key=${G_API_KEY}&callback=Function.prototype`,
        "GET")
}
