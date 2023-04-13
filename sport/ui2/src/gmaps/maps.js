import { Backdrop, CircularProgress } from '@mui/material';
import React, { useEffect, useState} from 'react';
import { G_API_KEY } from '../conf';
import { MulwiColors } from '../mulwiColors';

var to = null

export default function MulwiMap(props) {

  const[request, _setRequest] = useState(null)
  const [isLoading, setIsLoading] = useState(false)

  const [w, sw] = useState(0)

  function setRequest(w) {
    let r = "https://maps.googleapis.com/maps/api/staticmap?"

    // if (props.width) {
    //   request += `&size=${props.width}`
    // } else {
    //   request += 'size=400'
    // }

    r += "size=" + w
  
    if (props.height) {
      r += `x${props.height}`
    } else {
      r += 'x300'
    }
  
    if (props.mapType) {
      r += `&maptype=${props.mapType}}`
    } else {
      r += '&maptype=roadmap'
    }
  
    r += `&markers=size:mid%7C${props.center}`
    r += `&key=${G_API_KEY}`

    //console.log(r)

    _setRequest(r)
  }

  //const [to, setTo] = useState(null)

  function refresh() {
    setIsLoading(true)
    if(to) {
      window.clearTimeout(to)
    }
    to = setTimeout(() => {
          let key = props.el || "internmap"
          let el = document.getElementById(key)
          if(!el || !el.clientWidth) return
          let w  = el.clientWidth -  15
          setRequest(w)
          sw(w)
          setIsLoading(false)
      }, 200)
  }

  useEffect(() => {
    window.addEventListener("resize", refresh)
    return () => window.removeEventListener("resize", refresh)
  }, [])

  useEffect(() => {
    refresh()
  }, [props.center])

  return (
    <div id="internmap" style={{
      width: "100%"
    }}>
      <Backdrop
          sx={{ color: '#fff', zIndex: 9999999 }}
          style={{
              zIndex: 9999999,
              opacity: 0.4,
              margin: 15,
          }}
          open={isLoading}>
          <CircularProgress style={{
              color: MulwiColors.blueDark.greenDark
          }} />
      </Backdrop>
      <a href={`https://maps.google.com/maps/?q=${props.center}`}
        target="_blank" rel="noreferrer" style={{
          width: "100%"
        }}>
      {!isLoading && request && <img width={w} height={props.height} src={request} alt="google maps" />}
      </a>
    </div>
  )
}
