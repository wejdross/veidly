import { TextField, useTheme, Typography } from '@mui/material'
import { makeStyles } from "@mui/styles";
import React, { useEffect, useState } from 'react'
import { patchUserData } from '../apicalls/user.api'
import ModalEdit from './ModalEdit'
import { locale2 } from '../locale'
import { MulwiColors } from '../mulwiColors'

export default function AboutMeEdit(props) {

    const [aboutMe, setAboutMe] = useState("")
    const theme = useTheme()

    useEffect(() => {
      if(props.user)
        setAboutMe(props.user.AboutMe)
    }, [props.user])

    async function save(del) {
        await patchUserData({AboutMe: aboutMe})
        props.main.refresh()
    }

    if(!props.user) return null

    return (<React.Fragment>
        <ModalEdit
        lang={props.lang}

            buttonProps={props.buttonProps}
            onlyButton={props.onlyButton}

            title={`${locale2.ABOUT_ME[props.lang]}`}
            label={props.label || `${locale2.ABOUT_ME[props.lang]}`}
            labelStyle={{ color: MulwiColors.subtitleTypography }}

            custom
            onSave={save}
            value={
              <Typography style={{overflowWrap:"break-word", maxWidth: 320}} variant="body2">
                  {props.user.AboutMe}
              </Typography>}
            content={(<div>
                <TextField
                  multiline
                  inputProps={
                    { maxLength: 250 }
                  }
                  rows={5}
                  variant={"outlined"}
                  helperText={((aboutMe && aboutMe.length) || 0) + "/250"}
                  fullWidth={true}
                  value={aboutMe}
                  onChange={(event => (setAboutMe(event.target.value)))}
            />
            </div>)} />
    </React.Fragment>)
}

