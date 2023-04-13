import React, { useEffect, useRef, useState } from 'react';
import { Wrapper } from "@googlemaps/react-wrapper";
import { G_API_KEY } from '../conf';
import { arrayRepeat, sprintf } from '../helpers'
import { CircularProgress } from '@mui/material';
import { MulwiColors } from '../mulwiColors';

export default function GMap(props) {

    function distForLat(lat) {
        return 156.54303392 * Math.cos(lat * Math.PI / 180)
    }

    const [ow, _setow] = useState(null)

    const [loading, setLoading] = useState(false)
    const ref = useRef(null)
    const [map, setMap] = useState()

    const [mapCenter, setMapCenter] = useState({})
    const [mapZoom, setMapZoom] = useState(0)

    function updateOw() {
        let _lat = map.getCenter().lat()
        let _zoom = map.getZoom()
        let x = distForLat(_lat) / (1 << _zoom)
        let rpx = props.dist / x
        if (!ow || ow != rpx) {
            _setow(rpx)
        }
    }

    useEffect(() => {
        if (ref.current && !map) {
            let map = new window.google.maps.Map(ref.current, {})
            map.setOptions({
                zoom: 11,
                center: {
                    lat: 0,
                    lng: 0
                },
                disableDefaultUI: true,
                styles: [{
                    featureType: "poi",
                    elementType: "labels",
                    stylers: [
                        { visibility: "off" }
                    ]
                }]
            })
            map.addListener("center_changed", () => setMapCenter(map.getCenter()))
            map.addListener("zoom_changed", () => setMapZoom(map.getZoom()))
            setMap(map)
        }
    }, [ref, map])

    function isConsiderableCenterChange(oldCenter, newCenter) {

        let zoom = map.getZoom()
        if (!zoom)
            return false

        let oldLat = oldCenter.lat
        let oldLng = oldCenter.lng

        const rd = 1e5

        let newLat = Math.round(newCenter.lat * rd) / rd
        let newLng = Math.round(newCenter.lng * rd) / rd

        let latdf2 = Math.pow(newLat - oldLat, 2)
        let lngdf2 = Math.pow(newLng - oldLng, 2)
        let corr = Math.cos(newLat * Math.PI / 180)
        let dist = 110.25 * Math.sqrt((latdf2 + (lngdf2 * corr)))
        if (!dist || (dist / props.dist) < 0.1) {
            return false
        }

        return true
    }

    function getMapCenter() {
        let _center = map.getCenter()
        let center = {
            lat: _center.lat(),
            lng: _center.lng()
        }
        return center
    }

    function handleCenterChange(newCenter) {
        let center = getMapCenter()
        if (!center) {
            return
        }

        if (!isConsiderableCenterChange(center, newCenter))
            return

        map.setCenter(newCenter)
    }

    function adjustZoomToDist() {
        let el = document.getElementById("mwc")
        if (!el) return
        let lat = map.getCenter().lat()
        let F = distForLat(lat)
        if (!F) return
        let wpx = Math.min(el.clientWidth, el.clientHeight)
        if (!wpx) return
        let wkm = props.dist * 2
        let _zoom = Math.log2((wpx * F) / wkm)
        /* always making zoom one lower to allow for some spacing. */
        _zoom = Math.floor(_zoom - 1)
        /* updateOw(props.center.lat, _zoom, props.dist) */
        if (_zoom != map.getZoom()) {
            console.log(sprintf("adjusting zoom to %d to fit search results", _zoom))
            map.setZoom(_zoom)
        }
    }

    useEffect(() => {
        if (!map || !props.center)
            return
        handleCenterChange(props.center)
    }, [props.center, map])

    useEffect(() => {
        if (!map || !mapCenter || !mapZoom || !props.dist)
            return
        updateOw()
    }, [props.dist, mapCenter, mapZoom])

    const updatedZoom = useRef(false)
    useEffect(() => {
        if (!map || !props.dist || updatedZoom.current)
            return
        adjustZoomToDist()
        updatedZoom.current = true
    }, [map, props.dist])


    const to = useRef(null)
    function resetTimeout() {
        if (!loading)
            setLoading(true)
        if (to.current)
            window.clearTimeout(to.current)
        to.current = setTimeout(() => {
            props.onMapChange(getMapCenter())
            setLoading(false)
        }, 500);
    }

    useEffect(() => {
        if (!map || !isConsiderableCenterChange(props.center, getMapCenter()))
            return
        resetTimeout()
    }, [map, mapCenter])


    const markers = useRef([])


    function addMarkerForTraining(t, all, pr) {


        class Popup extends window.google.maps.OverlayView {

            position;
            containerDiv;

            constructor(t, userInfo, all, pr) {
                super()
                this.position 
                    = new window.google.maps.LatLng(t.LocationLat, t.LocationLng)
                this.containerDiv = document.createElement("div")
                this.containerDiv.style.position = "absolute"
                this.containerDiv.style.backgroundColor 
                    = pr ? MulwiColors.greenDark : MulwiColors.lightGreyAddedByLukasz
                this.containerDiv.style.borderRadius = "10px"
                this.containerDiv.className = 've-hoverable-icon'
                this.containerDiv.id = 'mapt' + t.ID
                this.containerDiv.style.cursor = 'pointer'

                let w = 30
                let h = 30

                this.containerDiv.style.width = w + "px"
                this.containerDiv.style.height = h + "px"

                this.containerDiv.title = t.Title

                this.containerDiv.onmouseenter = () => {
                    let el = document.getElementById('listt' + t.ID)
                    if(!el)
                        return
                    el.classList.remove('MuiPaper-elevation1')
                    el.classList.add('MuiPaper-elevation8')
                    el.scrollIntoView()
                }

                this.containerDiv.onmouseleave = () => {
                    let el = document.getElementById('listt' + t.ID)
                    if(!el)
                        return
                    el.classList.remove('MuiPaper-elevation8')
                    el.classList.add('MuiPaper-elevation1')
                }

                this.containerDiv.onclick = () => {
                    props.onSelect(all)
                }

                let inner = document.createElement('img')
                inner.style.maxWidth = '90%'
                inner.style.maxHeight = '90%'
                inner.style.width = 'auto'
                inner.style.height = 'auto'
                inner.style.position = 'absolute'
                inner.style.top = 0
                inner.style.bottom = 0
                inner.style.left = 0
                inner.style.right = 0
                inner.style.margin = 'auto'

                inner.src 
                    = userInfo.AvatarUrl 
                    //|| (userInfo.Name && userInfo.Name.length > 0 && userInfo.Name[0]) 
                    || (pr ? "/primary_pin.svg" : "/secondary_pin.svg")

                inner.alt = ':('

                // let inner = document.createElement("div")
                // inner.style.display = 'table'
                // inner.style.margin = '0 auto'
                // inner.style.fontSize = fs + "px"
                // inner.style.color = 'white'
                // inner.style.marginTop = Math.round(h/2 - 8) + "px"

                this.containerDiv.appendChild(inner)
            }

            onAdd() {
                this.getPanes().floatPane.appendChild(this.containerDiv)
            }

            onRemove() {
                if (this.containerDiv.parentElement) {
                    this.containerDiv.parentElement.removeChild(this.containerDiv)
                }
            }

            draw() {
                const divPosition = this.getProjection().fromLatLngToDivPixel(
                    this.position
                );
                // Hide the popup when it is far out of view.
                const display =
                    Math.abs(divPosition.x) < 4000 && Math.abs(divPosition.y) < 4000
                        ? "block"
                        : "none";

                if (display === "block") {
                    this.containerDiv.style.left = divPosition.x + "px";
                    this.containerDiv.style.top = divPosition.y + "px";
                }

                if (this.containerDiv.style.display !== display) {
                    this.containerDiv.style.display = display;
                }
            }
        }

        let m = new Popup(t, all.UserInfo, all, pr)
        m.setMap(map)
        // m.setAnimation(window.google.maps.Animation.BOUNCE)

        markers.current.push(m)
    }

    function clearMarkers() {

        for (let i = 0; i < markers.current.length; i++) {
            markers.current[i].setMap(null)
        }
        markers.current = []
    }

    useEffect(() => {
        if (!map)
            return
        clearMarkers()

        if (props.searchData) {
            for (let i = 0; i < props.searchData.length; i++) {
                addMarkerForTraining(props.searchData[i].Training, props.searchData[i], 1)
            }
        }
        if (props.suggestionData) {
            for (let i = 0; i < props.suggestionData.length; i++) {
                addMarkerForTraining(props.suggestionData[i].Training, props.suggestionData[i], 0)
            }
        }

            
    }, [map, props.searchData, props.suggestionData])

    return (<div id="mwc" style={{ display: "flex", height: "100%", position: "relative" }}>

        {ow && (<div style={{
            height: "100%",
            width: "100%",
            position: "absolute",
            pointerEvents: "none",
            backgroundColor: "rgba(0,0,0,0.15)",
            WebkitMaskImage: "radial-gradient(" + ow + "px at 50% 50% , transparent 93%, black 100%)",
            zIndex: 999
        }}></div>)}

        {loading && <CircularProgress style={{
            position: "absolute",
            right: 5,
            bottom: 5,
            zIndex: 1000
        }} />}

        <div ref={ref}
            style={{
                flexGrow: "1",
                height: "100%",
            }} />
    </div>)

}
