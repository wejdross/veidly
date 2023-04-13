import React from 'react'
import { RsvQrEval } from './rsvQrEval'
import { SubQrEval } from './subQrEval'

export function QrEval(props) {
    let query = new URLSearchParams(window.location.search)
    let t = query.get("type")
    if(t == "sub") {
        return <SubQrEval {...props} />
    }
    return <RsvQrEval {...props} />
}