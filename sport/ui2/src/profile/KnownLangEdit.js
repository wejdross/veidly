import { Typography
} from '@mui/material';
import React, { useEffect, useState } from 'react';
import { getLangLabel, locale2 } from '../locale';
import DynAtcEdit from './DynAtcEditor';
import { explainLangs, getLangs } from '../apicalls/user.api';
import { patchUserData } from '../apicalls/user.api';
import { MulwiColors } from '../mulwiColors';
  
export default function KnownLangEdit(props) {

    const [knownLangs, setKnownLangs] = useState([])
    const [options, setOptions] = useState([])

    useEffect(() => {
        if(!props.user) return
        (async () => {
            let langs = props.user.KnownLangs || []
            try {
                let l = await explainLangs(langs)
                l = JSON.parse(l)
                setKnownLangs(l)
            } catch(ex) {
                console.log(ex)
            }
        })()
      }, [props.user])
  
      async function save(x) {
        try {
            let u = props.user
            let cpy = []
            for(let i = 0 ; i < x.length; i++) {
                cpy.push(x[i].ISO_639_1)
            }
            u.KnownLangs = cpy
            await patchUserData(u)
            props.main.refresh()
        } catch(ex) {
            console.log(ex)
        }
      }

    async function queryLangs(q) {
        try {
            let l = await getLangs(q)
            setOptions(JSON.parse(l))
        } catch(ex) {
            console.log(ex)
        }
    }
    
    return (<DynAtcEdit 
        label={<Typography variant="body2" style={{
            color: MulwiColors.subtitleTypography
        }}>
            {locale2.KNOWN_LANGS[props.lang]}
        </Typography>}
        lang={props.lang}
        noLabelTypo
        atclabel={locale2.KNOWN_LANGS[props.lang]}
        options={options}
        value={() => (knownLangs || [])}
        updateOptions={queryLangs}
        equals={(o,v) => o.ISO_639_1===v.ISO_639_1}
        optionLabel={getLangLabel}
        valueLabel={getLangLabel}
        onChange={save} noFreeSolo/>)  
}