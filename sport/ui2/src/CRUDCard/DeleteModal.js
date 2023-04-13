import { Typography } from '@mui/material';
import { Delete } from '@mui/icons-material';
import React from 'react';
import { MulwiColors } from '../mulwiColors';
import ModalEdit from '../profile/ModalEdit';
import { locale2 } from '../locale';
import { sprintf } from '../helpers';

export function DeleteModal(props) {

    let lang = props.lang

    return (<ModalEdit
        lang={props.lang}
        onlyButton
        btnLabel={sprintf(locale2.DELETE_FMT[lang], props.objectName)}
        buttonProps={{
            style: {
                color: MulwiColors.redError,
            },
        }}
        onSave={async () => {
            await props.deleteData(props.record)
            props.onChange()
        }}
        content={<React.Fragment>
            <center><Typography>
                <strong>{props.nameSelector(props.record)}</strong>
            </Typography></center>
        </React.Fragment>}
        title={<Typography>
            {sprintf(locale2.DELETE_CONFIRM[lang], props.objectName)}?
        </Typography>}
        label={<Delete/>}
        btnColor={MulwiColors.redError}>
    </ModalEdit>)
}