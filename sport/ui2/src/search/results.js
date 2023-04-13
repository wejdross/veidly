import Pagination from "@mui/lab/Pagination";
import {
    Backdrop, Button, Chip, CircularProgress,
    Dialog, DialogActions, DialogContent, DialogTitle, Grid, Typography,
    useMediaQuery, useTheme
} from "@mui/material";
import React, { useEffect, useState } from "react";
import { useHistory, useLocation } from "react-router";
import { search } from "../apicalls/instructor.api";
import { DrawerResponsive } from "../card/DrawerResponsive";
import { getLocaleFromNavigator, getSupportedLanguage, locale2 } from "../locale";
import { MulwiColors } from "../mulwiColors";
import Navbar from "../navbar/Navbar";
import {
    getErrorDialog, getNullDialog,
    StatusDialog
} from "../StatusDialog";
import { FilterContent } from "./filterContent";
import NoResults from "./NoResults";
import SearchBar from "./searchBar";
import SingleTraining from "./singleTraining";
// import Suggestions from "./suggestions";
import GMap from "../gmaps/maps2";
import SelectOcc from "../reservations/SelectOcc";

function activateMapItem(v) {
    let el = document.getElementById('mapt' + v.Training.ID)
    if (!el)
        return
    el.style.width = '40px'
    el.style.height = '40px'
    el.classList.add('ve-hoverable-icon-raised')
}

function deactivateMapItem(v) {
    let el = document.getElementById('mapt' + v.Training.ID)
    if (!el)
        return
    el.style.width = '30px'
    el.style.height = '30px'
    el.classList.remove('ve-hoverable-icon-raised')
}

export default function SearchResults(props) {

    const h = useHistory()

    const [apiRequest, setApiRequest] = useState({})
    const [data, setdata] = useState(null)
    const [sdata, setSdata] = useState(null)
    const [page, setPage] = useState(0);
    const pageSize = 20
    const [loading, setLoading] = useState(true)

    const theme = useTheme()
    const isSmall = useMediaQuery(theme.breakpoints.down('sm'))
    const [info, setInfo] = useState(getNullDialog())

    const [invalidInfo, setInvalidInfo] = useState(false)

    useEffect(() => {
        try {
            let query = new URLSearchParams(window.location.search)
            let x = query.get("q")
            x = JSON.parse(x)
            //setUdist(x.DistKm)
        } catch (ex) {
            setInfo(getErrorDialog("Couldnt mount component", ex))
        }
    }, [])

    // const [hovers, setHovers] = useState([])
    // const [sdMeta, setSdMeta] = useState([])


    function opensum(el) {
        setSelectedElem(el)
        setOccOpen(true)
    }

    async function refreshFromUrl(l) {
        setLoading(true)
        try {
            let query = new URLSearchParams(l.search)
            let x = query.get("q")
            //console.log(JSON.parse(x))
            x = JSON.parse(x)
            x.Langs = [getSupportedLanguage(), getLocaleFromNavigator()[0]]
            x.OmitEmptySchedules = true
            x.OnlyAvailable = true
            x.SugEnabled = true
            x.SugDistKm = Math.max(150, x.DistKm)
            setApiRequest(x)
            //searchSuggestions(x)
            let d = await search(x)
            d = JSON.parse(d)

            // let len = ((d && d.Data) || []).length
            // setHovers(arrayRepeat(len, null))

            // let meta = arrayRepeat(len, {})
            // if (d.Data) for (let i = 0; i < meta.length; i++) {
            //     meta[i].open = (e) => opensum(e)
            // }
            //setSdMeta(meta)

            setSdata(d.SugData)
            setdata(d)
            if(invalidInfo)
                setInvalidInfo(false)
        } catch (ex) {
            if(ex === 400) {
                setInvalidInfo(true)
            } else {
                setInfo(getErrorDialog("Couldnt perform search request", ex))
            }
            //console.log(ex)
        } finally {
            if (isSmall)
                setfopen(false)
            setLoading(false)
        }
    }

    const location = useLocation();
    React.useEffect(() => {
        refreshFromUrl(window.location)
    }, [location]);

    function onChange(c) {
        if (!c) return
        let query = new URLSearchParams(window.location.search)
        query.set("q", JSON.stringify(c));
        h.push({
            search: query.toString()
        })
    }

    async function onMapChange(center) {
        try {
            // let rg = await reverseFetchGeo(center.lat, center.lng)
            // let city = null
            // rg = JSON.parse(rg)
            // if (rg && rg.address) {
            //     if (rg.address.city) {
            //         city = rg.address.city
            //     } else if (rg.address.town) {
            //         city = rg.address.town
            //     } else if (rg.address.village) {
            //         city = rg.address.village
            //     }
            // }
            let r = apiRequest
            r.display_name = ""
            r.Lat = center.lat
            r.Lng = center.lng
            // if (radius != r.DistKm)
            //     r.DistKm = radius
            onChange(r)
        } catch (ex) {
            setInfo(getErrorDialog("Couldnt handle map change", ex))
        }
    }

    const [fopen, setfopen] = useState(false)
    const [updateToken, setUpdateToken] = useState(0)


    useEffect(() => {
        if (!isSmall) setfopen(true)
        else setfopen(false)
    }, [isSmall])

    const [selectedElem, setSelectedElem] = useState(null)
    const [occOpen, setOccOpen] = useState(false)


    const [mobileMapOpen, setMobileMapOpen] = useState(false)

    return (<React.Fragment>
        <Dialog open={invalidInfo} onClose={() => setInvalidInfo(false)}>
            <DialogTitle>{locale2.INVALID_SEARCH_QUERY[props.lang]}</DialogTitle>
            <DialogContent>
                <SearchBar lang={props.lang} simple column fullSearchBtn />
            </DialogContent>
            <DialogActions>
                <Button onClick={() => setInvalidInfo(false)} color="secondary">
                    {locale2.CLOSE[props.lang]}
                </Button>
            </DialogActions>
        </Dialog>
        {!fopen && !isSmall && (
            <Button
                variant="contained"
                size={"medium"}
                style={{
                    zIndex: 9999,
                    position: "fixed",
                    color: "white",
                    top: (props.instructor && props.instructor.Config !== 0) ? 130 : 70,
                    right: 5,
                    backgroundColor: MulwiColors.greenDark
                }}
                onClick={() => {
                    setfopen(true)
                }}>
                {locale2.FILTERING[props.lang]}
            </Button>
        )}
        {!mobileMapOpen && !fopen && isSmall && (
            <Button
                variant="contained"
                size={"medium"}
                style={{
                    zIndex: 9999,
                    position: "fixed",
                    color: "white",
                    bottom: 40,
                    left: 0,
                    width: "100%",
                    backgroundColor: MulwiColors.greenDark
                }}
                onClick={() => {
                    setfopen(true)
                }}>
                {locale2.FILTERING[props.lang]}
            </Button>
        )}
        {isSmall && (
            <Button
                variant="contained"
                size={"medium"}
                style={{
                    zIndex: 9999,
                    position: "fixed",
                    color: "white",
                    bottom: 0,
                    left: 0,
                    width: "100%",
                    backgroundColor: MulwiColors.blueDark
                }}
                onClick={() => {
                    setMobileMapOpen(!mobileMapOpen)
                }}>
                {locale2.MAP[props.lang]}
            </Button>
        )}
        <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
        <Navbar main={props.main}
            user={props.user}
            content={<React.Fragment>
                <SearchBar
                    lang={props.lang}
                    nav simple searchRequest={apiRequest} />
            </React.Fragment>
            }
            lang={props.lang} setLang={props.setLang}
            instructor={props.instructor}>
            <DrawerResponsive
                padding={1}
                width={isSmall ? "100%" : null}
                navContent={
                    <Button
                        style={{
                            color: "white",
                            backgroundColor: MulwiColors.greenDark,
                        }}
                        onClick={() => {
                            // trigger update 
                            let x = updateToken
                            x++
                            setUpdateToken(x)
                        }}
                        fullWidth variant="contained">
                        {locale2.APPLY[props.lang]}
                    </Button>
                }
                open={fopen}
                onClose={() => setfopen(false)} content={
                    <React.Fragment>
                        <FilterContent lang={props.lang}
                            searchRequest={apiRequest}
                            updateToken={updateToken}
                            onChange={onChange} />
                    </React.Fragment>
                }>

                <Grid container direction="row">
                    <Grid item md={6} style={{ position: "relative", 
                    height: isSmall ? null : ((props.instructor && props.instructor.Config !== 0)
                    ? "calc(100vh - 120px)" : "calc(100vh - 65px)"),
                    // direction:"rtl",
                    overflowY: "auto" }}>
                        <Backdrop
                            sx={{ color: '#fff', zIndex: 9999999 }}
                            style={{
                                zIndex: 9999999,
                                opacity: 0.4,
                                position: "absolute",
                                height: "100%",
                            }}
                            open={loading}>
                            <CircularProgress style={{
                                color: MulwiColors.blueDark.greenDark
                            }} />
                        </Backdrop>

                        <div style={{
                            display: (data && data.Data && data.Data.length > 0 && "flex") || null,
                            flexDirection: "column",
                            minHeight: (props.instructor && props.instructor.Config !== 0)
                                ? "calc(100vh - 120px)" : "calc(100vh - 65px)"
                        }}>

                            <Grid container direction="column">

                                {isSmall && <SearchBar belowNav simple searchRequest={apiRequest} />}

                                <Grid container direction="row" justifyContent={"center"} style={{
                                    marginTop: 5
                                }} spacing={2}>
                                    {data && data.Meta && data.Meta.Tags && data.Meta.Tags.map(t => (<Grid item>
                                        <Chip style={{
                                            color: "white",
                                            backgroundColor: MulwiColors.blueDark
                                        }} label={<span>
                                            {(t.Tag.Translations[getSupportedLanguage()] || t.Tag.Name)}
                                            <strong>{" " + (t.NumberMatches || 0)}</strong>
                                        </span>} />
                                    </Grid>))}
                                </Grid>

                                {data && data.Data && data.Data.length > 0 && (<React.Fragment>
                                    <center>
                                        <Typography variant="h5" style={{ marginTop: 20 }}>
                                            <strong>{locale2.MATCHING_TRAININGS[props.lang]}</strong>
                                        </Typography>
                                    </center>
                                    <Grid
                                        justify="center"
                                        alignItems="center"
                                        // spacing={isSmall ? 0 : 3}
                                        style={{ marginTop: 20 }}
                                        container direction="row" item xs={12}>

                                        {selectedElem && <SelectOcc
                                            lang={props.lang}
                                            user={props.user}
                                            dr={{
                                                DateStart: apiRequest.DateStart,
                                                DateEnd: apiRequest.DateEnd
                                            }}
                                            open={occOpen}
                                            elem={selectedElem}
                                            setOpen={setOccOpen}
                                        />}

                                        {data.Data.map((v, i) => (
                                            <Grid key={i} item lg={!isSmall ? 12 : 4}
                                                style={{
                                                    marginBottom: isSmall ? 20 : null,
                                                    maxWidth: !isSmall ? 900 : 390,
                                                }}>
                                                <SingleTraining lang={props.lang}
                                                    onClick={() => opensum(v)}
                                                    onMouseEnter={() => {
                                                        activateMapItem(v)
                                                    }}
                                                    hover={/*hovers[i]*/ false}
                                                    onMouseLeave={() => {
                                                        deactivateMapItem(v)
                                                    }}
                                                    list={!isSmall}
                                                    d={v}
                                                    key={i} />
                                            </Grid>
                                        ))}
                                    </Grid>
                                </React.Fragment>)}


                                {data && data.Data && data.Data.length > 0 && (<Grid
                                    container
                                    direction={"row"}
                                    justify={"center"}
                                    style={{
                                        marginTop: "auto"
                                    }}
                                    alignItems={"center"}>
                                    <Pagination
                                        style={{
                                            marginTop: 10
                                        }}
                                        count={Math.ceil(data.Data.length / pageSize)}
                                        page={page + 1}
                                        onChange={(event, value) => {
                                            setPage(value - 1);
                                        }}
                                    />
                                </Grid>)}


                                {(!data || !data.Data || data.Data.length === 0) &&
                                    (<NoResults lang={props.lang} onChange={onChange} />)}

                                {/* <Suggestions lang={props.lang} 
                                        data={sdata} apiRequest={apiRequest} /> */}


                                {sdata && sdata.length > 0 && <div style={{
                                    marginTop: "auto"
                                }}>
                                    <center>
                                        <Typography variant="h5" style={{
                                            marginTop: 20,
                                            marginBottom: 20,
                                        }}>
                                            <strong>{locale2.SIMILAR_OFFERS[props.lang]}</strong>
                                        </Typography>
                                    </center>
                                    <React.Fragment>
                                        <Grid
                                            justify="center"
                                            alignItems="center"
                                            // spacing={isSmall ? 0 : 3}
                                            style={{ marginTop: 20 }}
                                            container direction="row" item xs={12}>
                                            {sdata && sdata.map((v, i) => {
                                                return (
                                                    <Grid item key={i} lg={!isSmall ? 12 : 4}
                                                        style={{
                                                            marginBottom: isSmall ? 20 : null,
                                                            maxWidth: !isSmall ? 900 : 390,
                                                        }}>
                                                        <SingleTraining lang={props.lang}
                                                            onClick={() => opensum(v)}
                                                            user={props.user}
                                                            onMouseEnter={() => {
                                                                activateMapItem(v)
                                                            }}
                                                            onMouseLeave={() => {
                                                                deactivateMapItem(v)
                                                            }}
                                                            dr={{
                                                                DateStart: apiRequest.DateStart,
                                                                DateEnd: apiRequest.DateEnd
                                                            }}
                                                            list={!isSmall}
                                                            d={v}
                                                            key={i} />
                                                    </Grid>
                                                )
                                            })}
                                        </Grid>
                                    </React.Fragment>
                                </div>}
                            </Grid>
                        </div>
                    </Grid>
                    {!isSmall && <Grid item md={6}>
                        <div style={{ height: "100%", width: "100%" }}>
                            <GMap
                                //sdMeta={sdMeta}
                                //setHovers={setHovers}
                                //distData={distData}
                                onSelect={v => opensum(v)}

                                onMapChange={onMapChange}

                                suggestionData={sdata}
                                searchData={data && data.Data}

                                dist={apiRequest && apiRequest.DistKm}
                                center={apiRequest && { lat: apiRequest.Lat, lng: apiRequest.Lng }} />
                        </div>
                    </Grid>}
                    {isSmall && <Dialog fullScreen fullWidth={true} open={mobileMapOpen}
                        onClose={() => setMobileMapOpen(false)}>
                        <div style={{ height: "100%", width: "100%" }}>
                            <GMap
                                // sdMeta={sdMeta}
                                // setHovers={setHovers}
                                // distData={distData}

                                onSelect={v => opensum(v)}

                                onMapChange={onMapChange}

                                suggestionData={sdata}
                                searchData={data && data.Data}

                                dist={apiRequest && apiRequest.DistKm}
                                center={apiRequest && { lat: apiRequest.Lat, lng: apiRequest.Lng }} />
                        </div>
                    </Dialog>}
                </Grid>

            </DrawerResponsive>
        </Navbar>
    </React.Fragment>)
}