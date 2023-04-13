import { MulwiColors } from "../mulwiColors"
import React from "react"

export function offToStr(off) {
    return ("calc(100vh - " + off + "px)")
}

export function getOffset(props) {
    if (props.instructor && props.instructor.Config != 0)
        return 120
    return 65
}


export function connIndicator(sc) {
    return (<div style={{
        width: 15,
        height: 15,
        borderRadius: 15,
        backgroundColor: sc
    }}>
    </div>)
}

let notifyTimeout = null
export let isWindowVisible = false

export function startNotify() {
    
    if(isWindowVisible)
        return

    document.title = "(!) Veidly"

    if(notifyTimeout)
        return
    
    let audio = new Audio("ding.mp3");
    audio.play();
    
    notifyTimeout = window.setTimeout(() => {
        notifyTimeout = null
    }, 20000)
}

export function endNotify() {
    document.title = "Veidly"
}

function visibilityChanged() {
    if (document.visibilityState === 'visible') {
        endNotify()
        isWindowVisible = true
    } else {
        isWindowVisible = false
    }
}

document.removeEventListener("visibilitychange", visibilityChanged)
document.addEventListener("visibilitychange", visibilityChanged)
visibilityChanged()

