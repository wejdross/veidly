import {
    TextField, Typography, useMediaQuery, useTheme,
} from '@mui/material'
import React, {  useState } from 'react'
import ModalEdit from './ModalEdit';

export default function TextAreaModalEdit(props) {

    const t = useTheme()
    const isLowRes = useMediaQuery(t.breakpoints.down('sm'))

    const [val, setVal] = useState(props.value || "")

    return <ModalEdit
    lang={props.lang}
            multiline={props.multiline}
        content={
            <React.Fragment>
            <TextField
                multiline
                inputProps={
                    { maxLength: 250 }
                }
                style={{
                    width: isLowRes ? null : 500
                }}
                rows={5}
                variant="outlined"
                helperText={(val || "").length + "/250"}
                fullWidth={true}
                value={val}
                onChange={e => setVal(e.target.value)}
            />
            </React.Fragment>
        }
        label={<Typography style={{color:"gray"}}>{props.label}</Typography>}
        title={props.title}
        value={<Typography noWrap>{props.value}</Typography>}
        onOpen={() => setVal(props.value)}
        onSave={() => props.onChange && props.onChange(val)}
    />
}