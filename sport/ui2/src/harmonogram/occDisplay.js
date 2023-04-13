import { Button, Dialog, DialogActions, DialogContent, 
            Grid, IconButton, List, 
            ListItem, Typography } from '@mui/material'
import { ScheduleOutlined } from '@mui/icons-material'
import { Pagination } from '@mui/lab'
import React, { useState } from 'react'
import { MulwiColors } from '../mulwiColors'
import { prettyPrintDateRange } from './trainingDetails'
import { locale2 } from '../locale'

function occ2Elem(occ, occ2, ds, lang) {
    let start = ds || occ.DateStart
    let s = new Date(start)
    let e = new Date(start)
    s.setMinutes(s.getMinutes() + occ2.OffsetStart)
    e.setMinutes(e.getMinutes() + occ2.OffsetEnd)
    return (<React.Fragment>
        <Grid container direction="column">
            <Grid item>
                <Typography>
                    <strong>{prettyPrintDateRange(s, e, 0, 0, lang)}</strong>
                </Typography>
            </Grid>
            <Grid item>
                {occ2.Remarks}
            </Grid>
        </Grid>
    </React.Fragment>)
} 

export function Occ2List(props) {
    let occ = props.occ
    if(!occ) return null
    return (<List>
        {occ.SecondaryOccs && occ.SecondaryOccs.map((s, i) => (<ListItem style={{
            borderLeft: "4px solid " +  s.Color || MulwiColors.blueDark,
        }} key={i}>
            {occ2Elem(occ, s, props.dateStart, props.lang)}
        </ListItem>))}
    </List>)
}

export function OccDisplay(props) {
    
    const [page, setPage] = useState(0)
    const [occ, setOcc] = useState(null)
    
    if(!props.t || !props.occs) return null

    

    const pageSize = 3

    return (<React.Fragment>
        <Dialog open={Boolean(occ)} onClose={() => setOcc(null)}>
            {occ && <DialogContent>
                    <Typography>
                        <strong>{props.t.Title}</strong>
                    </Typography>
                    <Typography>
                        {prettyPrintDateRange(occ.DateStart, occ.DateEnd, occ.RepeatDays, null, props.lang)}
                    </Typography>
                    <Typography variant="h6">
                        { locale2.SCHEDULE[props.lang] }
                    </Typography>
                    <Occ2List occ={occ} />
            </DialogContent>}
            <DialogActions>
                <Button onClick={() => setOcc(null)}>
                { locale2.CLOSE[props.lang] }
                </Button>
            </DialogActions>
        </Dialog>
        <Grid container spacing={3} 
                style={{marginRight: 0}}
                alignItems="center" justify="space-between">
            <Grid item>
                <Typography style={{color:"gray"}}
                            component={'span'} >
                    { locale2.OCCURRENCE[props.lang] }
                </Typography>
            </Grid>
            <Grid item>
                <List>
                    {props.occs.slice(page * pageSize, (page + 1) * pageSize).map((o, i) => <ListItem key={i}>
                        <Grid container direction="row" justify="space-between" alignItems="center">
                            <Grid item>
                                {prettyPrintDateRange(o.DateStart, o.DateEnd, o.RepeatDays, null, props.lang)}
                            </Grid>
                            {o.SecondaryOccs && o.SecondaryOccs.length > 0 && (<Grid item>
                                <IconButton size="small" 
                                    onClick={() => setOcc(o)}
                                    style={{
                                        marginBottom: 3,
                                        color: MulwiColors.blueDark
                                    }}><ScheduleOutlined/>
                                </IconButton>
                            </Grid>)}
                        </Grid>
                    </ListItem>)}
                </List>
                <Pagination count={Math.ceil(props.occs.length/pageSize)} page={page+1} onChange={(e,v) => {
                    setPage(v-1)
                }} />
            </Grid>
        </Grid>
    </React.Fragment>)
}