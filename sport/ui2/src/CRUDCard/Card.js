import {
    Table, TableBody, TableCell,
    TableContainer,
    TableHead, TableRow,
} from '@mui/material';
import React, { useEffect, useState } from 'react';
import CardWithBg from '../card/cardWithBg';
import { locale2 } from '../locale';
import { getErrorDialog } from '../StatusDialog';
import { AddModal } from './AddModal';
import { DeleteModal } from './DeleteModal';
import { SetBindingModal } from './SetBindingModal';

export function Card(props) {

    const [newKey, setNewKey] = useState(null)
    const [info, setInfo] = useState(null)
    const [data, setData] = useState(null)

    async function refresh(key) {
        try {
            let d = await props.getData()
            setData(d)
            if (key) setNewKey(key)
        } catch (ex) {
            setInfo(getErrorDialog("Couldnt refresh data", ex))
        }
    }

    useEffect(() => {
        refresh()
    }, [])

    return (<React.Fragment>
        <CardWithBg img="/static/form-backgrounds/surfer.webp">

            {props.cardHeader}

            <AddModal onChange={refresh} {...props} />
            <TableContainer>
            <Table >
                <TableHead>
                    <TableRow>
                        {props.tableColumns.map(c => (<TableCell>
                            {c.header}
                        </TableCell>))}
                        <TableCell align="center">
                            {locale2.TRAININGS[props.lang]}
                        </TableCell>
                        <TableCell align="center">
                            {locale2.EDIT[props.lang]}
                        </TableCell>
                        <TableCell align="center">
                            {locale2.DELETE[props.lang]}
                        </TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {data && data.map((row, i) => (<TableRow key={i}>
                        {props.tableColumns.map(c => (<TableCell>
                            {c.fieldSelector(row)}
                        </TableCell>))}
                        <TableCell align="center">
                            <SetBindingModal lang={props.lang} 
                                    setInfo={setInfo} record={row} 
                                    newKey={newKey} setNewKey={setNewKey} {...props} />
                        </TableCell>
                        <TableCell align="center">
                            <AddModal lang={props.lang} 
                                    setInfo={setInfo} record={row} onChange={refresh} {...props} />
                        </TableCell>
                        <TableCell align="center">
                            <DeleteModal lang={props.lang} 
                                    setInfo={setInfo} onChange={refresh} record={row} {...props} /> 
                        </TableCell>
                    </TableRow>))}
                </TableBody>
            </Table>
                
                </TableContainer>

        </CardWithBg>
    </React.Fragment>)
}
