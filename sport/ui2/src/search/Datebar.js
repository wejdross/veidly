import Grid from "@mui/material/Grid";
import { DatePicker as KeyboardDatePicker } from '@mui/x-date-pickers/DatePicker';
import React, { useEffect, useState } from "react";
import { locale2 } from "../locale";

export function DateBar(props) {
  const [dateStart, setDateStart] = useState(new Date());
  const [dateEnd, setDateEnd] = useState(new Date());
  const [selectOption, _] = useState(2);

  function onDateChange(s, e) {
    let r = props.searchRequest
    if (!s || !e) return
    r.DateStart = s
    r.DateEnd = e
    props.onChange(r)
  }

  useEffect(() => {
    if (props.searchRequest && props.searchRequest.DateStart 
            && props.searchRequest.DateEnd) {
      setDateStart(new Date(props.searchRequest.DateStart))
      setDateEnd(new Date(props.searchRequest.DateEnd))
    } 

  }, [props.searchRequest])

  function next() {
    let newStart, newEnd
    switch (selectOption) {
      // week
      case 1:
        newStart = new Date(dateStart)
        newStart.setDate(newStart.getDate() + 7)
        newEnd = new Date(dateEnd)
        newEnd.setDate(newEnd.getDate() + 7)
        newEnd.setHours(23, 59, 59)
        setDateStart(newStart)
        setDateEnd(newEnd)
        break
      // month
      case 2:
        newStart = new Date(dateStart)
        newStart.setMonth(newStart.getMonth() + 1)
        newEnd = new Date(dateEnd)
        newEnd.setMonth(newEnd.getMonth() + 1)
        newEnd.setHours(23, 59, 59)
        setDateStart(newStart)
        setDateEnd(newEnd)
        break
    }
    onDateChange(newStart, newEnd)
  }

  function prev() {
    let newStart, newEnd
    switch (selectOption) {
      // week
      case 1:
        newStart = new Date(dateStart)
        newStart.setDate(newStart.getDate() - 7)
        newEnd = new Date(dateEnd)
        newEnd.setDate(newEnd.getDate() - 7)
        newEnd.setHours(23, 59, 59)
        setDateStart(newStart)
        setDateEnd(newEnd)
        break
      // month
      case 2:
        newStart = new Date(dateStart)
        newStart.setMonth(newStart.getMonth() - 1)
        newEnd = new Date(dateEnd)
        newEnd.setMonth(newEnd.getMonth() - 1)
        newEnd.setHours(23, 59, 59)
        setDateStart(newStart)
        setDateEnd(newEnd)
        break
    }
    onDateChange(newStart, newEnd)
  }

  let maxDate = new Date()
  maxDate.setMonth(maxDate.getMonth() + 6)

  return (
    <React.Fragment>      
      <Grid
        container
        direction={"row"}
        justify={"center"}
        alignItems={"center"}
      >
        {/* <KeyboardDatePicker
          inputProps={{
            style: { textAlign: 'center' }
          }}
          style={{
            width: 150
          }}
          InputLabelProps={{
            style: {
              textAlign: 'center', width: '78%', transformOrigin: 'center top 0px'
            }
          }}
          maxDate={maxDate}
          disableToolbar
          value={dateStart || ""}
          onChange={(date) => {
            setDateStart(date);
            onDateChange(date, dateEnd)
          }}
          format="dd/MM/yyyy"
          id="date-picker-inline"
          label={locale2.FROM[props.lang]}
        />
        <KeyboardDatePicker
          disableToolbar
          value={dateEnd || ""}
          onChange={(date) => {
            setDateEnd(date);
            onDateChange(dateStart, date)
          }}
          style={{
            width: 150
          }}
          variant="inline"
          format="dd/MM/yyyy"
          id="date-picker-inline1"
          label={locale2.TO[props.lang]}
          minDate={dateStart}
          maxDate={maxDate}
          InputLabelProps={{
            style: {
              textAlign: 'center', width: '78%', transformOrigin: 'center top 0px'
            }
          }}
          inputProps={{
            style: { textAlign: 'center' }
          }}
        /> */}
      </Grid>
    </React.Fragment>
  );
}