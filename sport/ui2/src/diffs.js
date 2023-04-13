export function getDiffsByLang(lang) {
    let ret = []
    for(let k in diffs) {
        ret.push({
            e: diffs[k].enabled,
            id: Number(k),
            label: diffs[k][lang] || ""
        })
    }
    return ret
}

export function getDiffsObjByLang(lang) {
    let ret = {}
    for(let k in diffs) {
        ret[Number(k)] = {
            e: diffs[k].enabled,
            id: Number(k),
            label: diffs[k][lang] || ""
        }
    }
    return ret
}

export const diffs = {
    1: {
        enabled: true,
        pl: "podstawowy",
        en: "novice"
    },
    2: {
        enabled: true,
        pl: "średnio-zaawansowany",
        en: "medium"
    },
    3: {
        enabled: true,
        pl: "zaawansowany",
        en: "advanced"
    },
}
