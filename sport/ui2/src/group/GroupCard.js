import { Grid,  TextField, Typography } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { deleteGroup, getGroups, getTrainingsForGroups, 
        patchGroup, postGroup, deleteGroupBinding, putGroupBinding } from '../apicalls/grp';
import { getSupportedLanguage, locale2 } from '../locale';
import { Card } from '../CRUDCard/Card';
import { multilangVar2, sprintf } from '../helpers';

const BIGVAL = 9999999

export function GroupCard(props) {

    const [groups, setGroups] = useState([])

    async function refresh() {
        try {
            let d = await getGroups()
            d = JSON.parse(d)
            if (d) {
                setGroups(d)
            }
        } catch (ex) {
            console.log(ex)
        }
    }

    useEffect(() => {
        refresh()
    }, [])


    let lang = props.lang

    return (<React.Fragment>

        <Card 
            lang={props.lang}
            getData={async function () {
                let d = await getGroups()
                return JSON.parse(d)
            }}

            newReq={function () {
                return {
                    Name: "",
                    MaxPeople: BIGVAL,
                    MaxTrainings: BIGVAL
                }
            }}

            postData={async req => JSON.parse(await postGroup(req))}
            patchData={async req => await patchGroup(req)}
            deleteData={async req => await deleteGroup(req.ID)}

            cardHeader={<React.Fragment>
                <Typography variant="h5" style={{
                    marginBottom: 5
                }}>
                    {locale2.LIMITS[props.lang]}
                </Typography>
                <Typography variant="body2" style={{
                    maxWidth: 600,
                    marginBottom: 10,
                }}>
                    {locale2.LIMITS_DETAILS[props.lang]}
                </Typography>
            </React.Fragment>}

            tableColumns={[
                {
                    header: locale2.NAME[props.lang],
                    fieldSelector: e => e.Name
                },{
                    header: locale2.MAX_AMOUNT_OF_PEOPLE[props.lang],
                    fieldSelector: e => e.MaxPeople == BIGVAL ? locale2.UNLIMITED[lang] : e.MaxPeople
                },{
                    header: locale2.MAX_NO_TRAININGS[props.lang],
                    fieldSelector: e => e.MaxTrainings == BIGVAL ? locale2.UNLIMITED[lang] : e.MaxTrainings
                },
            ]}
  
            objectName={locale2.GRP_NAME[lang]}

            updateForm={(req, setReq) => (<React.Fragment>
                <Grid container spacing={2} direction="column" style={{width: 320}}>
                    <Grid item>
                        <TextField
                            fullWidth
                            value={req.Name}
                            onChange={e => setReq(c => ({...c, Name: e.target.value}))}
                            label={locale2.NAME[props.lang]}
                        />
                    </Grid>
                    <Grid item>
                        <Typography>
                            {locale2.MAX_AMOUNT_OF_PEOPLE[props.lang]}
                        </Typography>
                        <TextField
                            variant="outlined"
                            size="small"
                            type="number"
                            fullWidth
                            value={String(req.MaxPeople == BIGVAL ? "" : req.MaxPeople)}
                            onChange={e => {
                                let x = Number(e.target.value)
                                if (isNaN(x)) return
                                if(x <= 0) x = BIGVAL
                                setReq(c => ({...c, MaxPeople: x}))
                            }}
                            type="number" />
                            <Typography variant="body2" style={{maxWidth:320, minHeight: 50}}>
                                {req.MaxPeople != BIGVAL && sprintf(
                                    locale2.LIMIT_PEOPLE[props.lang], 
                                    req.MaxPeople, 
                                    multilangVar2(lang, req.MaxPeople)) || locale2.NO_LIMIT_PEOPLE[props.lang]}
                            </Typography>
                    </Grid>
                    <Grid item>
                        <Typography>
                            {locale2.MAX_NO_TRAININGS[props.lang]}
                        </Typography>
                        <TextField
                            type="number"
                            variant="outlined"
                            size="small"
                            fullWidth
                            value={String(req.MaxTrainings == BIGVAL ? "" : req.MaxTrainings)}
                            onChange={e => {
                                let x = Number(e.target.value)
                                if (isNaN(x)) return
                                if(x <= 0) x = BIGVAL
                                setReq(c => ({...c, MaxTrainings: x}))
                            }}
                            type="number" />
                        <Typography variant="body2" style={{maxWidth:320, minHeight: 80}}>
                            {req.MaxTrainings != BIGVAL && 
                                sprintf(locale2.LIMIT_TRAININGS[props.lang], req.MaxTrainings) 
                                        || locale2.NO_LIMIT_TRAININGS[props.lang]}
                        </Typography>
                    </Grid>
                </Grid>
            </React.Fragment>)}

            nameSelector={r => r.Name}

            getTrainings={async r => await getTrainingsForGroups([r.ID])}
            createBinding={async (r, tid) => await putGroupBinding(r.ID, tid)}
            deleteBinding={async (r, tid) => await deleteGroupBinding(r.ID, tid)}

        />

    </React.Fragment>)
}
