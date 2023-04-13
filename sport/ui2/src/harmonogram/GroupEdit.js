import { Button, Grid, List, ListItem, Typography } from '@mui/material'
import { ArrowLeft, ArrowRight } from '@mui/icons-material'
import React, { useEffect, useState } from 'react'
import { deleteGroupBinding, getGroups, putGroupBinding } from '../apicalls/grp'
import ModalEdit from '../profile/ModalEdit'
import { locale2 } from '../locale'

export function GroupEditContent(props) {
    
    const [groups, setGroups] = useState([])

    async function refresh() {
        try {
            let g = await getGroups()
            g = JSON.parse(g)
            let _g = []
            if(props.d.groups) {
                let ix = {}
                for(let i = 0; i < props.d.groups.length; i++) {
                    ix[props.d.groups[i].ID] = props.d.groups[i]
                }
                for(let i = 0; i < g.length; i++) {
                    if(!ix[g[i].ID])
                        _g.push(g[i])
                }
            } else {
                _g = g
            }
            setGroups(_g)
        } catch(ex) {
            console.log(ex)
        }
    }

    useEffect(() => {
        refresh()
    }, [props.d])

    async function addBinding(gid, tid) {
        try {
            await putGroupBinding(gid, tid)
            props.onChange()
        } catch(ex) {
            console.log(ex)
        }
    }

    async function rmBinding(gid, tid) {
        try {
            await deleteGroupBinding(gid, tid)
            props.onChange()
        } catch(ex) {
            console.log(ex)
        }
    }

    return (<React.Fragment>
        <Grid container direction="row" spacing={2}>
            <Grid sm={6} item style={{borderRight:"solid 1px gray"}}>
                <center>
                    <Typography variant="body2"><strong>{props.d.training.Title}</strong></Typography>
                </center>
                <List>
                    {props.d.groups && props.d.groups.map(g => (<ListItem>
                            <Button onClick={() => rmBinding(g.ID, props.d.training.ID)}>
                                {g.Name} <ArrowRight/> 
                            </Button>
                        </ListItem>))}
                </List>
            </Grid>
            <Grid item sm={6}>
                <center><Typography variant="body2"><strong>{ locale2.ALL_LIMITS[props.lang] }</strong></Typography></center>
                <List>
                    {groups && groups.map(g => (<ListItem>
                            <Button onClick={() => addBinding(g.ID, props.d.training.ID)}>
                                <ArrowLeft/> {g.Name}
                            </Button>
                        </ListItem>))}
                </List>
            </Grid>
        </Grid>
    </React.Fragment>)
} 

export function GroupEdit(props) {
    if(!props.d || !props.d.groups) return null
    
    return (<ModalEdit 
        lang={props.lang}
        hideSaveButton nocontent
        title="Limit"
        label={<Typography style={{color:"gray"}}>{ locale2.LIMITS[props.lang] }</Typography>}
        value={props.d.groups.map((g,i) => (i === 0 ? (
            <span key={g.Name}>{g.Name}</span>
        ) : (
            <span key={g.Name}>, {g.Name}</span>
        )))}
        //content={<GroupEditContent d={props.d} onChange={props.onChange}/>}
        content={null}
    />)
}