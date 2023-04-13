import { Grid, Typography } from "@mui/material";
import { Pagination } from "@mui/lab";
import React, { useEffect, useState } from "react";
import { getPubReviews } from "../apicalls/review";
import { locale2 } from "../locale";
import { ReviewContent } from "../review/userReview";

export default function RsvReviews(props) {

    const [rvs, setRvs] = useState([])

    const [page, setPage] = useState(0)
    const pageSize = 5

    async function setReviews(trainingID) {
        try {
            let d = await getPubReviews(trainingID)
            let r = JSON.parse(d) || []
            setRvs(r)
        } catch(ex) {
            props.setInfo && props.setInfo("Couldnt download reviews", ex)
        }
    }

    useEffect(() => {
        if(!props.training) return
        let id = props.training.ID
        setReviews(id)
    }, [props.training])

    return (<Grid container direction="column">
        {!rvs || rvs.length == 0 && (<Grid item>
            <Typography variant="h6">{locale2.REVIEWS[props.lang]}</Typography>
            <Grid container style={{
                marginLeft: 20,
                marginTop: 20
            }}>
                {locale2.NOONE_REVIEWED_THIS[props.lang]}
            </Grid>
        </Grid>)}

        {rvs && rvs.slice(page * pageSize, (page + 1) * pageSize).map((r,i) => (<Grid key={i} item>
            <ReviewContent content={r} />
        </Grid>))}

        {rvs && (<Pagination 
            count={Math.ceil(rvs.length / pageSize)} 
            page={page + 1}
            onChange={(_, value) => {
                setPage(value - 1)
            }}
        />)}
    </Grid>)
}
