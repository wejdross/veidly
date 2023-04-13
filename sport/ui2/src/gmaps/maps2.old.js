import React, { useEffect, useRef, useState } from 'react';
import { Wrapper } from "@googlemaps/react-wrapper";
import { G_API_KEY } from '../conf';
import { arrayRepeat, sprintf } from '../helpers'
import { CircularProgress } from '@mui/material';

const primaryIcon = "/primary_pin.svg"
const secondaryIcon = "/secondary_pin.svg"

function Map({ zoomChanged, 
                centerChanged, onClick, options, 
                searchData, suggestionData, 
                distData, setHovers, sdMeta }) {

    const ref = useRef(null)
    const [map, setMap] = useState()
    const [o, so] = useState(null)

    const pMarkers = useRef([])
    const sMarkers = useRef([])

    useEffect(() => {
        if (ref.current && !map) {
            setMap(new window.google.maps.Map(ref.current, {}));
        }
    }, [ref, map])

    function sdToPos(d) {
        return ({
            lat: d.Training.LocationLat,
            lng: d.Training.LocationLng
        })
    }

    function sdToI(_, primary, dist) {
        let iw = 0, ih = 0
        if(dist) {
            iw = 44
            ih = 60
        } else {
            iw = 22
            ih = 40
        }
        let iu = primary ? primaryIcon : secondaryIcon
        return {
            url: iu,
            scaledSize: {
                width: iw,
                height: ih
            }
        }
    }

    function diffMarkers(searchData, cache, primary, dd) {
        if (!searchData)
            searchData = []

        let mci = Math.min(cache.length, searchData.length)

        for (let i = 0; i < mci; i++) {
            let p = cache[i].getPosition()
            let icon = cache[i].getIcon()
            let d = searchData[i]
            if (p.lat() != d.Training.LocationLat || p.lng() != d.Training.LocationLng) {
                cache[i].setPosition(sdToPos(d))
                console.log("updated position for " + i)
            }

            let ic = sdToI(d, primary, dd == i)

            if (icon.url != ic.url || 
                        icon.scaledSize.width != ic.scaledSize.width || 
                        icon.scaledSize.height != ic.scaledSize.height) {
                cache[i].setIcon(ic)
                console.log("updated icon for " + i)
            }

            let iw = cache[i].__infoWindow

            if(d.Training.Title != iw.content) {
                iw.setContent(d.Training.Title)
            }
        }

        if (cache.length == searchData.length)
            return

        if (cache.length < searchData.length) {
            if (!map)
                return
            for (let i = mci; i < searchData.length; i++) {
                let _m = new window.google.maps.Marker()
                let d = searchData[i]
                _m.setOptions({
                    position: sdToPos(d),
                    icon: sdToI(d, primary, i == dd)
                })
                _m.setMap(map)
                const infowindow = new window.google.maps.InfoWindow({
                    content: d.Training.Title,
                })
                _m.__infoWindow = infowindow
                _m.addListener("mouseover", function() {
                    infowindow.open(_m.get("map"), _m, {shouldFocus: false})
                    let x = arrayRepeat(searchData.length, false)
                    x[i] = true
                    setHovers(x)
                })
                _m.addListener("mouseout", function() {
                    infowindow.close()
                    let x = arrayRepeat(searchData.length, false)
                    setHovers(x)
                })
                _m.addListener("click", function() {
                    sdMeta[i] && sdMeta[i].open && sdMeta[i].open(d)
                })  
                cache.push(_m)
                console.log("adding marker for " + i)
            }
        } else {
            /* points have been removed */
            for (let i = mci; i < cache.length; i++) {
                cache[i].setMap(null)
                //cache[i].removeListener("mouseover")
                //cache[i].removeListener("mouseout")
                console.log("removing marker for " + i)
            }
            cache.splice(mci, cache.length - mci)
        }
    }

    // useEffect(() => {
    //     diffMarkers(searchData, pMarkers.current, true, distData)
    // }, [map, searchData, distData])
    
    // useEffect(() => {
    //     diffMarkers(suggestionData, sMarkers.current, false, null)
    // }, [map, suggestionData])

    useEffect(() => {
        if (map) {
            if (o && JSON.stringify(o) === JSON.stringify(o))
                return
            so(options)
            map.setOptions(options);
        }
    }, [map, options])

    useEffect(() => {
        if (!map) {
            return
        }

        if (centerChanged) {
            map.addListener("center_changed", () => centerChanged(map))
        }

        if (zoomChanged) {
            map.addListener("zoom_changed", () => zoomChanged(map))
        }

        if (onClick) {
            map.addListener("click", (e) => onClick(e))
        }
    }, [map])

    return (
        <React.Fragment>
            <div ref={ref}
                style={{
                    flexGrow: "1",
                    height: "100%",
                }} />
            {/* {React.Children.map(children, (child) => {
                if (React.isValidElement(child)) {
                    return React.cloneElement(child, { map });
                }
            })} */}
        </React.Fragment>
    );
}

export default function GMap(props) {

    const [zoom, setZoom] = React.useState(11)
    //const [points, setPoints] = React.useState([])

    function distForLat(lat) {
        return 156.54303392 * Math.cos(lat * Math.PI / 180)
    }

    function updateOw(_lat, _zoom, _dist) {
        console.log(_lat, _zoom, _dist)
        let x = distForLat(_lat) / (1 << _zoom)
        let rpx = _dist / x
        if (ow != rpx)
            setow(rpx)
    }

    const [ow, setow] = useState(null)

    const zoomIsSet = useRef(false)

    /*
        make zoom big enough to contain results
    */
    useEffect(() => {
        if(zoomIsSet.current || !props.dist || !props.center)
            return
        let el = document.getElementById("mwc")
        if (!el) return
        let F = distForLat(props.center.lat)
        if (!F) return
        let wpx = Math.min(el.clientWidth, el.clientHeight)
        if (!wpx) return
        zoomIsSet.current = true
        let wkm = props.dist * 2
        let _zoom = Math.log2((wpx * F) / wkm)
        _zoom = Math.floor(_zoom-1)
        /* updateOw(props.center.lat, _zoom, props.dist) */
        if (_zoom != zoom) {
            console.log(sprintf("adjusting zoom to %d to fit search results", _zoom))
            setZoom(_zoom)
        }
    }, [props.dist, props.center])


    // function calculateRadius(_lat, _zoom) {
    //     let F = distForLat(_lat)
    //     let x = F / (1<<_zoom)
    //     let el = document.getElementById("mwc")
    //     if (!el) return 0
    //     let wkm = x * Math.min(el.clientWidth, el.clientHeight)
    //     wkm /= 2 
    //     wkm = Math.round(wkm)
    //     if(wkm > 100) wkm = 100
    //     if(wkm < 5) wkm = 5
    //     return wkm
    // }


    // useEffect(() => {
    //     let a1 = pointsFromSD(props.searchData, true)
    //     let a2 = pointsFromSD(props.suggestionData, null)
    //     setPoints(a1.concat(a2))
    // }, [props.searchData, props.suggestionData, props.distData])

    useEffect(() => {
        if (!props.center)
            return
        updateOw(props.center.lat, zoom, props.dist)
    }, [props.center])

    useEffect(() => {
        if (!props.dist || !props.center) return
        updateOw(props.center.lat, zoom, props.dist)
    }, [props.dist])

    const [mc, setmc] = useState(null)

    function onMapChange(newCenter, newZoom) {

        if (!newCenter || !newZoom) return

        console.log("onMapChange")
        let oldLat = props.center.lat
        let oldLng = props.center.lng

        const rd = 1e5
        //let [newZoom, newCenter] = [m.newZoom, m.newCenter]

        let newLat = Math.round(newCenter.lat * rd) / rd
        let newLng = Math.round(newCenter.lng * rd) / rd

        // determine if latlng changed considerably
        let latdf2 = Math.pow(newLat - oldLat, 2)
        let lngdf2 = Math.pow(newLng - oldLng, 2)
        let corr = Math.cos(newLat * Math.PI / 180)
        let dist = 110.25 * Math.sqrt((latdf2 + (lngdf2 * corr)))
        if (!dist) {
            return
        }

        if (newZoom == zoom && (dist / props.dist) < 0.1) {
            //console.log("no effective change")
            return
        }
        console.log(JSON.stringify({ lat: oldLat, lng: oldLng }) + " -> " + JSON.stringify(newCenter))

        if (newZoom != zoom)
            setZoom(newZoom)

        //let wkm = calculateRadius(newLat, newZoom)
        //updateOw(newLat, newZoom, props.dist)

        props.onMapChange(props.dist, newCenter)
    }

    useEffect(() => {
        onMapChange(mc, zoom)
    }, [mc])

    useEffect(() => {
        onMapChange(mc, zoom)
    }, [zoom])

    const timer = useRef(null)

    const [loading, setLoading] = useState(false)

    function refreshTimeout(m) {
        if (!loading)
            setLoading(true)
        if (timer.current) {
            window.clearTimeout(timer.current)
        }
        timer.current = setTimeout(() => {
            let nc = m.getCenter()
            setmc({ lat: nc.lat(), lng: nc.lng() })
            setZoom(m.getZoom())
            //onMapChange({newZoom: m.getZoom(), newCenter: {lat: nc.lat(),lng: nc.lng()}})
            setLoading(false)
            timer.current = null
        }, 500)
    }

    function zoomChanged(m) {
        refreshTimeout(m)
    }

    function centerChanged(m) {
        refreshTimeout(m)
    }

    return (
        <div id="mwc" style={{ display: "flex", height: "100%", position: "relative" }}>
            <div style={{
                height: "100%",
                width: "100%",
                position: "absolute",
                pointerEvents: "none",
                backgroundColor: "rgba(0,0,0,0.10)",
                WebkitMaskImage: "radial-gradient(" + ow + "px at 50% 50% , transparent 95%, black 100%)",
                zIndex: 999
            }}></div>
            {loading && <CircularProgress style={{
                position: "absolute",
                right: 5,
                bottom: 5,
                zIndex: 1000
            }} />}
            {/* <Snackbar
                open={sopen}
                onClose={() => setSopen(false)}
                autoHideDuration={5000}
                message="cowardly refusing to load data for radius this big"
             /> */}
                <Map
                    centerChanged={centerChanged}
                    zoomChanged={zoomChanged}

                    sdMeta={props.sdMeta}
                    suggestionData={props.suggestionData}
                    searchData={props.searchData}
                    distData={props.distData}
                    setHovers={props.setHovers}

                    options={{
                        zoom: zoom,
                        center: props.center,
                        disableDefaultUI: true,
                        styles: [{
                            featureType: "poi",
                            elementType: "labels",
                            stylers: [
                                { visibility: "off" }
                            ]
                        }]
                    }}>
                </Map>
        </div>)
}
