import { Typography } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { createInstructor } from './apicalls/instructor.api';
import Spinner from './Spinner';

export function BecomeTrainer(props) {
    const [err, setErr] = useState(null)
    const h = useHistory()
    useEffect(() => {
        let x = async () => {
            try {
                await createInstructor(null)
                props.main.refreshInstructor()
                h.push("/configure")
            } catch(ex) {
                setErr(ex)
            }
        }
        x()
    }, [props.main, h])
    
    return (err && (
        <Typography>
            {err}
        </Typography>
    )) || (
        <Spinner/>
    )
}