import LocationOnIcon from '@mui/icons-material/LocationOn';
import Autocomplete from '@mui/material/Autocomplete';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import { useTheme } from '@mui/material/styles';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import { debounce } from '@mui/material/utils';
import parse from 'autosuggest-highlight/parse';
import * as React from 'react';
import { useEffect } from "react";
import { googlePlaceIDtoJSON } from '../apicalls/user.api';
import { G_API_KEY } from '../conf';
import { locale2 } from "../locale";


const GOOGLE_MAPS_API_KEY = G_API_KEY;

function loadScript(src, position, id) {
  if (!position) {
    return;
  }

  const script = document.createElement('script');
  script.setAttribute('async', '');
  script.setAttribute('id', id);
  script.src = src;
  position.appendChild(script);
}

const autocompleteService = { current: null };

export default function GoogleMaps(props) {
  const [value, setValue] = React.useState(null);
  const [inputValue, setInputValue] = React.useState('');
  const [options, setOptions] = React.useState([]);

  const fetch = React.useMemo(
    () =>
      debounce((request, callback) => {
        autocompleteService.current.getPlacePredictions(request, callback);
      }, 400),
    [],
  );

  React.useEffect(() => {
    let active = true;

    if (!autocompleteService.current && window.google) {
      autocompleteService.current =
        new window.google.maps.places.AutocompleteService();
    }
    if (!autocompleteService.current) {
      return undefined;
    }

    if (inputValue === '') {
      setOptions(value ? [value] : []);
      return undefined;
    }

    fetch({ input: inputValue }, (results) => {
      if (active) {
        let newOptions = [];

        if (value) {
          newOptions = [value];
        }

        if (results) {
          newOptions = [...newOptions, ...results];
        }
        setOptions(newOptions);
      }
    });
    return () => {
        active = false;
    };
}, [value, inputValue, fetch]);
    // below func ensures that search engine receives data in correct form
    useEffect(() => {
        var ret = {}
        if (!value) {
            return
        }
        if (value && !value.place_id) {
            return
        }
        return (
            async () => {
                let x1 = await googlePlaceIDtoJSON(value.place_id)
                ret ={
                    lat: JSON.parse(x1).results[0].geometry.location.lat,
                    lon: JSON.parse(x1).results[0].geometry.location.lng,
                    display_name: JSON.parse(x1).results[0].formatted_address
                }
                props.setLocation(ret)
            }
            ) ()
    }, [value])
    // below useEffect sets localisation so it's correctly passed to search bar
    useEffect(() => {
      if (props.location && props.location.display_name && props.location.display_name.length > 0) {
        setValue(props.location.display_name)
      }
    }, [props.location])

  return (
    <Autocomplete
      id="google-map-autocomplete"
      getOptionLabel={(option) =>
        typeof option === 'string' ? option : option.description
      }
      className={props.class}
      filterOptions={(x) => x}
      options={options}
      size={props.size}
      autoComplete
      includeInputInList
      filterSelectedOptions
      value={value}
      noOptionsText={locale2.GIVE_ME_ADDRESS[props.lang]}
      onChange={(event, newValue) => {
        if (!newValue || newValue.length === 0) {
          newValue = ''
        }
        setOptions(newValue ? [newValue, ...options] : options);
        setValue(newValue);
      }}
      onInputChange={(event, newInputValue) => {
        setInputValue(newInputValue);
      }}
      renderInput={(params) => (
        <TextField 
            {...params}
            error={Boolean(props.error)}
            helperText={props.errorText}
            placeholder={props.label || locale2.WHERE[props.lang]} 
            variant="outlined"
            />
      )}
      renderOption={(props, option) => {
        const matches =
          option.structured_formatting.main_text_matched_substrings || [];

        const parts = parse(
          option.structured_formatting.main_text,
          matches.map((match) => [match.offset, match.offset + match.length]),
        );
        return (
          <li {...props}>
            <Grid container alignItems="center">
              <Grid item sx={{ display: 'flex', width: 44 }}>
                <LocationOnIcon sx={{ color: 'text.secondary' }} />
              </Grid>
              <Grid item sx={{ width: 'calc(100% - 44px)', wordWrap: 'break-word' }}>
                {parts.map((part, index) => (
                  <Box
                    key={index}
                    component="span"
                    sx={{ fontWeight: part.highlight ? 'bold' : 'regular' }}
                  >
                    {part.text}
                  </Box>
                ))}

                <Typography variant="body2" color="text.secondary">
                  {option.structured_formatting.secondary_text}
                </Typography>
              </Grid>
            </Grid>
          </li>
        );
      }}
    />
  );
}