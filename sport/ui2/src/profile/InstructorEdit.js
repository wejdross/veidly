import { Grid, Typography } from '@mui/material'
import React, { useEffect, useState } from 'react'
import { PATCHInstructor } from '../apicalls/instructor.api'
import { MulwiColors } from '../mulwiColors'
import { InstrModal } from './instrModal'
import ModalEdit from './ModalEdit'
import { locale2 } from '../locale'

export default function InstructorEdit(props) {

    const [instrInfo, setInstrInfo] = useState({})

    async function save() {
        if(!instrInfo) throw 400
        await PATCHInstructor(instrInfo)
        props.main.refreshInstructor()
    }

    useEffect(() => {
        if(props.instructor) {
            setInstrInfo({
                tags: props.instructor.Tags,
                yearExp: props.instructor.YearExp
            })
        }
      }, [props.instructor])

    if (!props.instructor) return null

    return (<React.Fragment>
        <ModalEdit
            buttonProps={{
                variant: "contained",
                style: {
                    color: "white",
                    backgroundColor: MulwiColors.greenDark,
                    marginLeft: 10
                },
                size: "small"
            }}
            onlyButton={props.onlyButton}
            lang={props.lang}
            label={props.label || (<React.Fragment>
                    <Typography variant="body2" 
                        style={{ color: MulwiColors.subtitleTypography }}>
                            {locale2.INSTRUCTOR_SINCE[props.lang]}
                    </Typography>
            </React.Fragment>)}
            custom
            onSave={save}
            value={<Grid container direction="column">
                <Typography variant="body2">
                    {props.instructor.YearExp || "2021"}&nbsp;
                </Typography>
            </Grid>}
            content={<InstrModal 
                lang={props.lang}
                instrInfo={instrInfo}
                setInstrInfo={setInstrInfo}
                instructor={props.instructor} />} />
    </React.Fragment>)
}
