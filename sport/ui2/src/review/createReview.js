import {
    Button, CircularProgress,
    Dialog, DialogActions,
    DialogContent, DialogTitle,
    TextField, Typography
} from '@mui/material';
import { Rating } from '@mui/lab';
import React, { useState } from 'react';
import { postUserReview } from '../apicalls/review';
import { MulwiColors } from '../mulwiColors';
import { locale2 } from '../locale';

export function CreateReview(props) {

    const [value, setValue] = useState(0);
    const [c, setc] = useState("");

    const [state, setState] = useState(0)
    const [err, setErr] = useState(null)

    async function doCreateReview() {
        setState(1)
        try {
            let r = {
                Mark: value || "",
                Review: c || "",
                AccessToken: props.accessToken
            }
            await postUserReview(r)
            props.onChange()
        } catch(ex) {
            setErr(ex)
            setState(3)
        }
    }
    
    function loading() {
        if(state !== 1) return null
        return (<DialogContent>
           <center>
            <CircularProgress/>
           </center>
        </DialogContent>)
    }

    function errorContent() {
        if(state !== 3) return null
        return (<React.Fragment>
            <DialogContent>
            <Typography>{locale2.SOMETHING_WENT_WRONG[props.lang]}</Typography>
            <Typography style={{
                color: MulwiColors.redError
            }}>{String(err)}</Typography>
        </DialogContent>
            <DialogActions>
                <Button onClick={() => setState(0)}>
                    {locale2.ONCE_AGAIN[props.lang]}
                </Button>
                <Button onClick={() => props.setOpen(false)}>
                    {locale2.CLOSE[props.lang]}
                </Button>
            </DialogActions>
        </React.Fragment>)
    }

    function successContent() {
        if(state !== 2) return null
        return (<React.Fragment>
            <DialogContent>
            {locale2.REVIEW_HAS_BEEN_SAVED[props.lang]}
        </DialogContent>
            <DialogActions>
                <Button onClick={() => props.setOpen(false)}>
                    {locale2.ADD_REVIEW[props.lang]}
                </Button>
            </DialogActions>
        </React.Fragment>)
    }

    function initContent() {
        if(state !== 0) return null
        return (<React.Fragment>
            <DialogContent>
                <center>
                    <Typography variant="h6" style={{
                        color: MulwiColors.blueDark
                    }}>
                        {props.training}
                    </Typography>
                    <br/>
                    <Rating 
                        max={6}
                        value={value} 
                        onChange={(e,v)=>setValue(v)} />
                </center>
                <TextField  
                    multiline
                    inputProps={
                        { maxLength: 250 }
                    }
                    rows={5}
                    variant={"outlined"}
                    helperText={((c && c.length) || 0) + "/250"}
                    fullWidth={true}
                    value={c}
                    onChange={(event => (setc(event.target.value)))}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={() => props.setOpen(false)}>
                    {locale2.CLOSE[props.lang]}
                </Button>
                <Button onClick={doCreateReview} variant="contained" style={{
                    color: "white",
                    backgroundColor: MulwiColors.greenDark
                }}>
                    {locale2.SAVE[props.lang]}
                </Button>
            </DialogActions>
        </React.Fragment>)
    }

    return (<React.Fragment>
        <Dialog open={props.open} onClose={() => props.setOpen(false)}>
            <DialogTitle>{locale2.ADD_REVIEW[props.lang]}</DialogTitle>
            {initContent()}
            {loading()}
            {successContent()}
            {errorContent()}
        </Dialog>
    </React.Fragment>)
}

