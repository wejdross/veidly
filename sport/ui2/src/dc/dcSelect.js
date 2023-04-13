import { Button, Dialog, DialogContent, 
            TextField, Typography } from '@mui/material';
import React, { useState } from 'react';
import { redeemDc } from '../apicalls/dc';
import { MulwiColors } from '../mulwiColors';
import { errToStr } from '../StatusDialog';
import { locale2 } from '../locale';

export function DcSelectModal(props) {

    const [open, setOpen] = useState(false)
    const [err, setErr] = useState("")
    const [c, setc] = useState("")

    async function rc() {
        try {
            let _c = await redeemDc(c, props.trainingID)
            _c = JSON.parse(_c)
            setErr("")
            props.onChange(_c)
            setOpen(false)
        } catch (ex) {
            if(ex == 404) {
                setErr(locale2.DIDNT_FIND_DC[props.lang])
                return
            }
            setErr(errToStr(ex))
            console.log(ex)
        }
    }

    return (<React.Fragment>
        <Dialog open={open} onClose={() => setOpen(false)}>
            <DialogContent>
                <Typography>
                    {locale2.ENTER_DC[props.lang]}
                </Typography>
                <TextField 
                    fullWidth
                    size="small" variant="outlined" 
                    value={c} onChange={e => setc(e.target.value)} />
                    <Button variant="contained" style={{
                        color: "white",
                        marginTop: 5,
                        backgroundColor: MulwiColors.blueDark
                    }} fullWidth onClick={rc}>
                        {locale2.CHECK[props.lang]}
                    </Button>
                <Typography 
                    variant="body2"
                    style={{color: MulwiColors.redError}}>{err}&nbsp;</Typography>
            </DialogContent>
        </Dialog>
        <Button size="small" onClick={() => setOpen(true)}>
            {locale2.DO_YOU_HAVE_DC[props.lang]}
        </Button>
    </React.Fragment>)
}