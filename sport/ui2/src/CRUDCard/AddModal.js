import { TextField, Typography } from '@mui/material';
import { Edit } from '@mui/icons-material';
import React, { useEffect, useState } from 'react';
import { sprintf } from '../helpers';
import { MulwiColors } from '../mulwiColors';
import ModalEdit from '../profile/ModalEdit';
import { locale2, } from '../locale';

export function TextInput(props) {
    return (<React.Fragment>
        <Typography>
            {props.label}
        </Typography>
        <TextField
            variant="outlined"
            size="small"
            fullWidth
            value={props.value}
            onChange={props.onChange} />
    </React.Fragment>)
}

export function AddModal(props) {

    let lang = props.lang

    const [req, setReq] = useState(null)

    useEffect(() => {
        if (!props.record) {
            if(props.newReq) setReq(props.newReq())
            return
        }
        setReq(props.record)
    }, [props.record])

    if(!req) return null

    return (<ModalEdit
        lang={props.lang}
        onlyButton
        buttonProps={props.record ? {
            style: {
                color: MulwiColors.blueDark,
            },
        } : {
            style: {
                color: "white",
                backgroundColor: MulwiColors.greenDark
            },
            variant: "contained"
        }}
        onSave={async () => {
            if (props.record) {
                await props.patchData(req)
                props.onChange(null)
            } else {
                let res = await props.postData(req)
                if (res && res.ID) {
                    props.onChange(res.ID)
                } else {
                    props.onChange(null)
                }
            }
        }}
        content={props.updateForm(req || {}, setReq)}
        title={props.record ? 
            sprintf(locale2.EDIT_FMT[lang], props.objectName) 
                : 
            sprintf(locale2.ADD_FMT[lang], props.objectName)}
        label={props.record ? <Edit /> : sprintf(locale2.ADD_FMT[lang], props.objectName)}
    ></ModalEdit>)
}
