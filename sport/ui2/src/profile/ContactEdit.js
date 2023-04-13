import { Checkbox, FormControlLabel, Grid, 
        TextField, Typography } from '@mui/material'
import React, { useEffect, useState } from 'react'
import { patchUserContactData } from '../apicalls/user.api'
import { MulwiColors } from '../mulwiColors'
import ModalEdit from './ModalEdit'
import { locale2 } from '../locale'

export default function ContactEdit(props) {

    const [email, setEmail] = useState("")
    const [phone, setPhone] = useState("")
    const [share, setShare] = useState(false)

    async function save(del) {
        let cd = { ...props.user.ContactData }
        cd.Email = email
        cd.Phone = phone
        cd.Share = share
        await patchUserContactData(cd)
        props.main.refresh()
    }

    useEffect(() => {
        if(props.user) {
            setEmail(props.user.ContactData.Email)
            setPhone(props.user.ContactData.Phone)
            setShare(props.user.ContactData.Share)
        }
      }, [props.user])

    if (!props.user) return null

    return (<React.Fragment>
        <ModalEdit
        lang={props.lang}

            buttonProps={props.buttonProps}

            onlyButton={props.onlyButton}

            title={locale2.CONTACT_DATA[props.lang]}
            label={props.label || (<React.Fragment>
                <Typography variant="body2" style={{ color: MulwiColors.subtitleTypography }}>{locale2.PHONE[props.lang]}</Typography>
                <Typography variant="body2" style={{ color: MulwiColors.subtitleTypography }}>Email</Typography>
            </React.Fragment>)}
            custom
            onSave={save}
            value={<Grid container direction="column">
                <Typography variant="body2">
                    {props.user.ContactData.Phone}&nbsp;
                </Typography>
                <Typography variant="body2">
                    {props.user.ContactData.Email}&nbsp;
                </Typography>
            </Grid>}
            content={(<div>
                <Grid container spacing={2} direction="column">
                    <Grid item style={{
                        maxWidth: 400
                    }}>
                        {props.instructor ? (
                            <React.Fragment>
                                <Typography variant="body2" >
                                {locale2.CONTACT_PASS[props.lang]}
                            </Typography>
                            <Typography variant="body2" >
                                {locale2.CONTACT_PASS_2[props.lang]}
                            </Typography>
                            </React.Fragment>
                        ) : (<Typography variant="body2" >
                            {locale2.CONTACT_PASS_CUST[props.lang]}
                        </Typography>)}
                        <Typography variant="body2">
                            <strong>{locale2.NOONE_ELSE_WILL_SEE[props.lang]}</strong>
                        </Typography>
                    </Grid>
                    <Grid item>
                        <TextField
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            label="Email" />
                    </Grid>
                    <Grid item>
                        <TextField
                            value={phone}
                            onChange={(e) => setPhone(e.target.value)}
                            label={locale2.PHONE[props.lang]} />
                    </Grid>
                    <Grid item>
                        <FormControlLabel
                            control={<Checkbox
                                checked={share}
                                onChange={e => setShare(e.target.checked)} />}
                            label={locale2.SHARE_PRIOR_BOOKING[props.lang]}
                        />
                        {<Typography variant="body2" style={{color: MulwiColors.subtitleTypography, maxWidth: 400}}>
                        {locale2.IF_YOU_DONT_SHARE_CONTACT[props.lang]}
                        </Typography>}
                    </Grid>
                </Grid>
            </div>)} />
    </React.Fragment>)
}
