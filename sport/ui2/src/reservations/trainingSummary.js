import {
    Chip, Grid, List, ListItem,
    Typography, useMediaQuery, useTheme
} from '@mui/material';
import { Rating } from '@mui/lab';
import React, { useEffect, useState } from 'react';
import { getDiffsObjByLang } from '../diffs';
import MulwiMap from '../gmaps/maps';
import { Occ2List } from '../harmonogram/occDisplay';
import {
    getTagLabel, prettyPrintDateRange,
    RsvInfoListItem
} from '../harmonogram/trainingDetails';
import { MulwiColors } from '../mulwiColors';
import { ImgGrid } from '../training/ImgGrid';
import { locale2 } from '../locale';
import { explainTags } from '../apicalls/user.api';
import { DisabledSupport } from './disabled';


export function TrainingSummary(props) {

    function printTrainingStatus(count, capacity) {
        if (count >= capacity) {
            return <Typography style={{ color: MulwiColors.redError }}>
                {locale2.NO_AVAILABLE_SLOTS_FOR_THIS_TRAINING[props.lang]}
            </Typography>
        }
        return <Typography>
            <strong>{capacity - count} / {capacity}</strong>
        </Typography>
    }

    const t = useTheme()
    const isLowRes = useMediaQuery(t.breakpoints.down('sm'))

    const [tags, setTags] = useState(null)

    async function refreshTags() {
        if (!props.training) {
            return
        }
        if (!props.training.Tags || props.training.Tags.length == 0) return
        try {
            let ts = JSON.parse(await explainTags(props.training.Tags))
            setTags(ts)
        } catch (ex) {
            console.log(ex)
        }
    }

    useEffect(() => {
        refreshTags()
    }, [props.training])

    const [w, _sw] = useState(10)
    function sw() {
        let el = document.getElementById("mapcontainer1")
        if (!el) return
        let w = el.clientWidth - 15
        if (isLowRes) {
            w = "70vw"
        }
        _sw(w)
    }
    useEffect(() => {
        window.addEventListener("resize", sw)
        sw()
        return () => window.removeEventListener("resize", sw)
    }, [props.resize])


    function rsvInfo(r) {
        return <Grid item>
            <Typography variant="h6" style={{ marginBottom: 15 }}>
                {locale2.YOUR_RSV[props.lang]}
            </Typography>
            <RsvInfoListItem
                lang={props.lang} rsv={r}
                setInfo={props.setInfo} onChange={props.onChange} />
        </Grid>
    }

    const diffs = getDiffsObjByLang("pl")

    let contentMargin = isLowRes ? 0 : 20

    function hgi(t) {
        return (
            <Grid item>
                <Typography style={{
                    textAlign: "left",
                    color: MulwiColors.subtitleTypography
                }}>
                    <strong>{t}</strong>
                </Typography>
            </Grid>)
    }


    function record(t, v) {
        return (
            <Grid item style={{ width: "100%" }}>
                <Grid container direction="row" alignItems="center"
                    justify="space-between">
                    <Grid item>
                        <Typography style={{
                            textAlign: "left",
                            color: MulwiColors.subtitleTypography
                        }}>
                            <strong>{t}</strong>
                        </Typography>
                    </Grid>
                    <Grid item style={{ marginLeft: 5 }}>
                        {v}
                    </Grid>
                </Grid>
            </Grid>)
    }

    return (props.training && (<React.Fragment>
        <Grid container
            spacing={2}
            style={{
                paddingLeft: isLowRes ? 0 : 30,
                paddingRight: isLowRes ? 0 : 30
            }}
            direction="column"
            alignItems={isLowRes ? "center" : null}
            justify="center" >
            <Grid item>
                {props.usrRsv && props.sch &&
                    props.sch.Reservations &&
                    props.sch.Reservations[0] &&
                    rsvInfo(props.sch.Reservations[0])}
            </Grid>
            <Grid item>
                <Grid container direction='row' justifyContent='space-between' alignContent='center'>
                    <Grid item>
                        <Typography variant="h5">
                            {props.training.Title}
                        </Typography>
                    </Grid>
                    <Grid item>
                            <DisabledSupport training={props.training} lang={props.lang} />
                    </Grid>
                </Grid>
            </Grid>
            <Grid item style={{ marginLeft: contentMargin }}>
                <Grid container
                    alignItems={isLowRes ? "center" : null}
                    direction={isLowRes ? "column" : "row"}
                    spacing={1}>
                    {tags && tags.map((t, i) => (
                        <Grid item>
                            <Chip
                                key={i}
                                label={getTagLabel(t)}
                                style={{
                                    marginRight: 5,
                                    color: "white",
                                    backgroundColor: MulwiColors.blueDark
                                }} />
                        </Grid>
                    ))}
                    {!tags && props.training.Tags && props.training.Tags.map((t, id) => (
                        <Grid item>
                            <Chip
                                key={id}
                                label={t}
                                style={{
                                    marginRight: 5,
                                    color: "white",
                                    backgroundColor: MulwiColors.blueDark
                                }} />
                        </Grid>
                    ))}
                </Grid>
            </Grid>
            {(props.training.MainImgUrl
                || (props.training.SecondaryImgUrls &&
                    props.training.SecondaryImgUrls.length > 0)) && (
                    <Grid item style={{ display: "table", margin: "0 auto" }}>
                        <ImgGrid
                            MainImgUrl={props.training.MainImgUrl}
                            SecondaryImgUrls={props.training.SecondaryImgUrls} />
                    </Grid>
                )}
            {props.sch && (
                <Grid item style={{ marginBottom: 20 }}>
                    {props.sch.Occ && props.sch.Occ.length > 0 && (
                        <Typography>
                            {prettyPrintDateRange(
                                props.sch.Start,
                                props.sch.End,
                                0, 0, props.lang)}
                        </Typography>
                    )}
                    <Occ2List dateStart={props.dateStart} occ={props.sch.Occ} />
                </Grid>
            )}
            {props.training.Description && <React.Fragment>
                {hgi(locale2.DESCRIPTION[props.lang])}
                <Grid item style={{
                    width: w
                }}>
                    <Typography
                        style={{
                            textAlign: isLowRes ? "center" : "left",
                            wordWrap: "break-word",
                            whiteSpace: "pre-wrap",
                        }}>
                        {props.training.Description}
                    </Typography>
                </Grid>
            </React.Fragment>}

            {record(locale2.REVIEWS[props.lang],
                <Grid item style={{ marginLeft: contentMargin }}>
                    <Grid container direction="row">
                        <Grid item>
                            <Rating readOnly value={props.training.AvgMark} max={6} />
                        </Grid>
                        <Grid item style={{ marginLeft: 5 }}>
                            <Typography>
                                ({props.training.NumberReviews})
                            </Typography>
                        </Grid>
                    </Grid>
                </Grid>)}

            {!props.usrRsv && (
                <React.Fragment>
                    {record(locale2.PRICE[props.lang],
                        (<React.Fragment>
                            <Grid item>
                                <Chip style={{
                                    marginLeft: contentMargin,
                                    padding: 10,
                                    color: "white",
                                    backgroundColor: MulwiColors.priceColor
                                }} label={<Typography>
                                    {props.training.Price / 100} {props.training.Currency}
                                </Typography>} />
                            </Grid>
                        </React.Fragment>))}

                    {props.sch && props.sch.Count >= 0 && record(
                        locale2.AVAILABLE_SPOTS[props.lang],
                        <React.Fragment>
                            {printTrainingStatus(
                                props.sch.Count,
                                props.training.Capacity)}
                        </React.Fragment>)}
                </React.Fragment>
            )}

            {props.training.Diff && props.training.Diff.length > 0 && (
                record(locale2.LEVEL[props.lang], (<Grid item>
                    <Typography style={{ marginLeft: contentMargin }}>
                        {props.training.Diff.map((di, i) => (
                            <Chip
                                style={{
                                    color: "white",
                                    backgroundColor: MulwiColors.blueDark
                                }}
                                label={diffs[di].val} key={i} />
                        ))}
                    </Typography>
                </Grid>))
            )}

            {(props.training.MinAge && props.training.MaxAge && (
                <React.Fragment>
                    {record(locale2.AGE[props.lang],
                        <Grid item style={{ marginLeft: contentMargin }}>
                            <Typography>
                                {props.training.MinAge} - {props.training.MaxAge}
                            </Typography>
                        </Grid>)}
                </React.Fragment>
            )) || null}

            {hgi(locale2.GEAR[props.lang])}
            <Grid item style={{ marginLeft: contentMargin }}>
                <Grid container direction={isLowRes ? "column" : "row"} spacing={2} justify="space-around">
                    <Grid item>
                        {hgi(locale2.REQUIRED_GEAR[props.lang])}
                        <List>
                            {(!props.training.RequiredGear
                                || props.training.RequiredGear.length === 0) && (
                                    <Typography variant="body2">
                                        {locale2.NONE[props.lang]}
                                    </Typography>
                                )}
                            {props.training.RequiredGear
                                && props.training.RequiredGear.map((t, i) =>
                                (<ListItem key={i}>
                                    <strong>{t}</strong>
                                </ListItem>))}
                        </List>
                    </Grid>
                    <Grid item>
                        {hgi(locale2.RECOMMENDED_GEAR[props.lang])}
                        <List>
                            {(!props.training.RequiredGear
                                || props.training.RequiredGear.length === 0) && (
                                    <Typography variant="body2">
                                        {locale2.NONE[props.lang]}
                                    </Typography>
                                )}
                            {props.training.RecommendedGear
                                && props.training.RecommendedGear.map((t, i) =>
                                (<ListItem key={i}>
                                    {t}
                                </ListItem>))}
                        </List>
                    </Grid>
                    <Grid item>
                        {hgi(locale2.INSTRUCTOR_GEAR[props.lang])}
                        <List>
                            {(!props.training.RequiredGear
                                || props.training.RequiredGear.length === 0) && (
                                    <Typography variant="body2">
                                        {locale2.NONE[props.lang]}
                                    </Typography>
                                )}
                            {props.training.InstructorGear
                                && props.training.InstructorGear.map((t, i) =>
                                (<ListItem key={i}>
                                    {t}
                                </ListItem>))}
                        </List>
                    </Grid>
                </Grid>
            </Grid>
            {record(locale2.LOCATION[props.lang], (
                <Grid item>
                    <Typography style={{ textAlign: "left", marginLeft: contentMargin }}>
                        {props.training.LocationText}
                    </Typography>
                </Grid>))}

            <Grid item id="mapcontainer1" style={{
                width: "100%",
            }}>
                <MulwiMap el="mapcontainer1"
                    width={isLowRes ? 300 : null}
                    height={isLowRes ? 300 : null}
                    center={props.training.LocationText} />
            </Grid>
        </Grid>
    </React.Fragment>)) || null
}