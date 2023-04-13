import { MulwiColors } from "./mulwiColors";
import { createMuiTheme, createTheme } from "@mui/material";

const tokenStr = "mulwii_token"

export const defaultLogoPath = "static/logo.text.svg"

export function gettoken() {
    return localStorage.getItem(tokenStr);
}

export function settoken(tk) {
    localStorage.setItem(tokenStr, tk);
}

export function rmtoken() {
    localStorage.removeItem(tokenStr);
}

export function toUTC(date) {
    date = new Date(date)
    return Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate(),
        date.getUTCHours(), date.getUTCMinutes(), date.getUTCSeconds())
}

export function asUTC(date) {
    date = new Date(date)
    return Date.UTC(date.getFullYear(), date.getMonth(), date.getDate(),
        date.getHours(), date.getMinutes(), date.getSeconds())
}

export function dfInHours(start, end) {
    let toDate = parseInt(toUTC(end) / 1000)
    let fromDate = parseInt(toUTC(start) / 1000)
    return (toDate - fromDate) / 3600
}

export function dateIsNotZero(d) {
    if (!d) return false
    if (d.getFullYear() === 1) {
        return false
    } else {
        return true
    }
}

export function dateStartOfWeek(d) {
    d = new Date(d)
    d.setHours(0, 0, 0)
    var day = d.getDay(),
        diff = d.getDate() - day + (day === 0 ? -6 : 1) // adjust when day is sunday
    d.setDate(diff)
    return d
}

// will return differences spanning over multiple days
export function extendedDfInHours(start, end) {
    let ret = {}

    if (start.getDate() === end.getDate() &&
        start.getMonth() === end.getMonth() &&
        start.getFullYear() === end.getFullYear()) {
        return {
            [dateToEpoch(start)]: dfInHours(start, end)
        }
    }

    let cpy = new Date(start)

    while (cpy <= end) {

        let ed = new Date(cpy)
        ed.setHours(23, 59, 59, 0)

        let d = dfInHours(cpy, end)
        let dday = dfInHours(cpy, ed)

        if (d >= dday) {
            ret[dateToEpoch(cpy)] = dday
        } else {
            // end of loop
            ret[dateToEpoch(cpy)] = d
            break // may not be needed
        }

        cpy.setDate(cpy.getDate() + 1)
        cpy.setHours(0, 0, 0, 0)

        // move it to the beginning of next day
        // cpy.setDate(cpy.getDate() + 1)
        // cpy.setHours(0,0,0,0)
        // let d = new Date(cpy)

        // if((cpy.getDate()-1) == start.getDate()) {
        //     console.log(1, i)
        //     d.setDate(cpy.getDate() - 1)
        //     ret[dateToEpoch(d)] = dfInHours(start, cpy) 
        // } else if(cpy.getDate() == end.getDate()+1) {
        //     console.log(2, i)
        //     d.setDate(cpy.getDate() - 1)
        //     ret[dateToEpoch(d)] = 24-dfInHours(end, cpy)
        // } else {
        //     console.log(3, i)
        //     ret[dateToEpoch(d)] = 24
        // }          
    }
    return ret
}

export function isLoggedIn() {
    if (gettoken()) return true
    return false
}
/*
    verifyToken()
    true = verified; user has token
    false = no token 
    I don't care what kind of token, because everything further will be verified by API
    if token is present but bad, then when I receive 401 I'll clear this token and redirect to #/login
*/
export function verifyToken() {
    try {
        if (gettoken().length !== 0) {
            return true;
        } else {
            return false;
        }
    } catch {
        return false;
    }
}

export function getPathFromSearchParams(sp) {
    return window.location.pathname + "?" + sp.toString()
}

export function QSsetAndReturn(k, v) {
    let q = new URLSearchParams(window.location.search);
    q.set(k, encodeURIComponent(v))
    return q
}

export const primaryColorClass = "grey darken-2"

export function dateToEpoch(d) {
    return (d.getTime() - d.getMilliseconds()) / 1000;
    //return Math.floor(d.getTime() / 1000)
}

export function epochToDate(e) {
    let d = new Date(0)
    d.setUTCSeconds(e)
    return d
}

export function microEpochToDate(e) {
    let d = new Date(0)
    d.setUTCMilliseconds(e/1000)
    return d
}

export function dateToMicroEpoch(d) {
    d = new Date(d)
    return d.getUTCMilliseconds() * 1000
}

export function dateStartOfDay() {
    return new Date(new Date().setHours(0, 0, 0, 0));
}

export function dateEndOfDay() {
    return new Date(new Date().setHours(23, 59, 59, 999));
}

export function arrayRepeat(len, value) {
    var arr = [];
    for (var i = 0; i < len; i++) {
        arr.push(value);
    }
    return arr;
}

export function randomString(length) {
    var result = '';
    var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for (var i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

export function getFromQueryString(location, key) {
    if (!location || !key) {
        return null;
    }
    let ix = location.hash.indexOf(key + "=");
    if (ix < 0) {
        return null;
    }
    let ss = location.hash.substr(ix + key.length + 1);
    let eix = ss.indexOf("&");
    if (eix < 0) {
        return ss;
    }
    return ss.substr(0, eix);
}

export function dateToHour(date) {
    if (!date) return null;
    return date.getHours() + ":" + date.getMinutes();
}

export function concatDateTime(date, hr) {
    let ret = new Date(date);
    let parts = hr.split(":")
    if (parts.length !== 2) {
        return
    }
    ret.setHours(parseInt(parts[0]), parseInt(parts[1]));
    return ret;
}

export function displayLimit(str, limit) {
    if (str.length > limit) {
        return str.slice(0, limit) + "...";
    }
    return str;
}

export function getEmailDomain(email) {
    let x = email.split("@");
    if (x.length === 2) {
        return x[1];
    }
    return null;
}

export function redirectToLogin(location, history) {
    if (!location)
        location = window.location.pathname + window.location.search + window.location.hash;
    if (history) {
        history.push("/login?return_url=" + escape(location.toString()))
    } else {
        window.location = ("/login?return_url=" + escape(location.toString()))
    }
}

export function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export function prettyPrintCurrency(c) {
    switch (c) {
        case "PLN":
            return "zł"
        default:
            return c
    }
}


export function removeFromQs(key) {
    let query = new URLSearchParams(window.location.search)
    query.delete(key)
    let q = ""
    if(query.toString()){
        q = "?" + query.toString()
    }
    var newurl = window.location.protocol + "//" 
        + window.location.host + window.location.pathname + q;
    window.history.pushState({path:newurl},'',newurl);
}

export function trainingResToDrawerData(tr, user, scheduleRecord) {
    if (!tr.Training) {
        throw "invalid training res"
    }
    return {
        training: tr.Training,
        occs: tr.Occurrences,
        groups: tr.Groups,
        dc: tr.Dcs,
        sch: scheduleRecord,
        user: user,
        sm: tr.Sms
    }
}

const avReasons = {
    1: "Sesja jest już pełna",
    2: "Trening lub Sesja nie istnieje lub została przesunięta",
    4: "Rezerwacje na ten trening są wyłączone",
    8: "Nie można się już zapisać na ten trening",
    16: "Nieprwaidłowe występowanie",
    32: "Grupa jest pełna",
    64: "Grupa jest pełna",
    128: "Instruktor ma wolne"
}

// right now r is not bit flags
export function avReasonToStr(r) {
    let x = avReasons[r]
    return x || ""
}

// 1-based days
export const dayIndex = {
    1: {
        pl: "Pn",
        en: "Mo",
        de: "Mo"
    }, 2: {
        pl: "Wt",
        en: "Tu",
        de: "Di",
    }, 3: {
        pl: "Śr",
        en: "We",
        de: "Mi"
    }, 4: {
        pl: "Cz",
        en: "Th",
        de: "Do"
    }, 5: {
        pl: "Pt",
        en: "Fr",
        de: "Fr"
    }, 6: {
        pl: "So",
        en: "Sa",
        de: "Sa",
    }, 7: {
        pl: "Nd",
        en: "Su",
        de: "So"
    },
} 

export const lMonths = [
    {
        pl: "Styczeń",
        en: "January",
        de: "Januar"
    },
    {
        pl: "Luty",
        en: "Febuary",
        de: "Februar"
    },
    {
        pl: "Marzec",
        en: "March",
        de: "März"
    },
    {
        pl: "Kwiecień",
        en: "April",
        de: "April"
    },
    {
        pl: "Maj",
        en: "May",
        de: "Mai"
    },
    {
        pl: "Czerwiec",
        en: "June",
        de: "Juni"
    },
    {
        pl: "Lipiec",
        en: "July",
        de: "Juli"
    },
    {
        pl: "Sierpień",
        en: "August",
        de: "August"
    },
    {
        pl: "Wrzesień",
        en: "September",
        de: "September"
    },
    {
        pl: "Październik",
        en: "October",
        de: "October",
    },
    {
        pl: "Listopad",
        en: "November",
        de: "November"
    },
    {
        pl: "Grudzień",
        en: "December",
        de: "Dezember"
    },
]

// obsolete
export const months = {
    0: "Styczeń",
    1: "Luty",
    2: "Marzec",
    3: "Kwiecień",
    4: "Maj",
    5: "Czerwiec",
    6: "Lipiec",
    7: "Sierpień",
    8: "Wrzesień",
    9: "Październik",
    10: "Listopad",
    11: "Grudzień",
}

export const DatePickerCustomTheme = createTheme({
    overrides: {
        MuiPickersToolbar: {
            toolbar: {
                backgroundColor: MulwiColors.greenDark,
            },
        },
        MuiPickersCalendarHeader: {
            switchHeader: {
                // backgroundColor: lightBlue.A200,
                // color: "white",
            },
        },
        MuiPickersDay: {
            day: {
                color: MulwiColors.blueLight,
                "&:hover": {
                    backgroundColor: MulwiColors.pinkAction,
                },
            },
            daySelected: {
                backgroundColor: MulwiColors.greenDark,
                "&:hover": {
                    backgroundColor: MulwiColors.pinkAction,
                },
            },
            dayDisabled: {
                color: MulwiColors.lightGreyAddedByLukasz,
            },
            current: {
                color: "black",
            },
        },
        MuiPickersModal: {
            dialogAction: {
                color: MulwiColors.greenDark,
            },
        },
        MuiPickerBasePicker: {
            pickerView: {
                color: "red",
            },
        },
        MuiPickersClockPointer: {
            thumb: {
                borderColor: MulwiColors.greenLight
            },
            pointer: {
                backgroundColor: MulwiColors.greenLight
            },
        },
        MuiPickersClock: {
            pin: {
                color: MulwiColors.greenLight + "!important"
            }
        }
    },
});

/*
    supported formats:
        %s - string
        %d - digit - note that this will work for any argument with 'toString' prototype fn
        %v - any (currently same handling as %s)
        %j - object / json
    _fmt must be string
*/
export function sprintf(_fmt, ...args) {
    if(typeof _fmt !== 'string') return null
    let carr = []
    let ai = 0
    for(let i = 0; i < _fmt.length; i++) {
        let c = _fmt[i]
        if(c !== '%' || i == (_fmt.length - 1))  {
            carr.push(c)
            continue
        }
        let n = _fmt[i + 1]
        switch(n) {
        case '%':
            carr.push(c)
            break
        case 'd':
            if(args && args.length > ai && args[ai].toString) {
                carr.push(args[ai].toString())
                ai++
            } else {
                carr.push("!%" + n)
            }
            break
        case 'v':
        case 's':
            if(args && args.length > ai) {
                carr.push(args[ai])
                ai++
            } else {
                carr.push("!%" + n)
            }
            break
        case 'j':
            if(args && args.length > ai) {
                carr.push(JSON.stringify(args[ai]))
                ai++
            } else {
                carr.push("!%" + n)
            }
            break
        default:
            break
        }
        i++
    }

    return carr.join('')
}

export function polishVar1(num) {
    if(num < 1) {
        return "będzie mogło"
    }
    if(num < 2) {
        return "będzię mógł"
    }
    if(num < 5) {
        return "będą mogły"
    }
    return "będzie mogło"
}


export function englishVar2(num) {
    if(num == 1) return "person"
    return "people"
}

export function polishVar2(num) {
    if(num < 1) {
        return "klientów"
    }
    if(num < 2) {
        return "klient"
    }
    return "klientów"
}

export function multilangVar2(lang, num) {
    switch(lang) {
    case "en":
        return englishVar2(num)
    case "pl":
        return polishVar2(num)
    }
}

export function isUuid(val) {
    return val.match('^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$')
}
