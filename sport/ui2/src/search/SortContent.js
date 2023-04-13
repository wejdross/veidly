import { FormControl, MenuItem, Select } from "@mui/material";
import React, { useEffect, useState } from "react";
import { locale2 } from "../locale";

export function SortContent(props) {

    const [sort, setSort] = useState(-1)

    useEffect(() => {
      let r = props.searchRequest || {}
      if(r.Sort) {
        for(let i = 0; i < r.Sort.length; i++) {
          let x = r.Sort[i]
          switch(x.Column) {
          case "price":
            if(x.IsDesc) {
              setSort(1)
            } else {
              setSort(0)
            }
            break
          case "capacity":
            setSort(2)
            break
            //s.Dostępność = true
          }
        }
      }
    }, [props.searchRequest])

    return (
    <React.Fragment>
          <FormControl >
            <Select
              labelId="demo-simple-select-helper-label"
              id="demo-simple-select-helper"
              onChange={(e) => {
                setSort(e.target.value)
                let s = []
                switch(e.target.value) {
                  case 0:
                      s.push({
                          Column: "price",
                          IsDesc: false
                      })
                      break
                  case 1:
                      s.push({
                          Column: "price",
                          IsDesc: true
                      })
                      break
                  case 2:
                      s.push({
                          Column: "capacity",
                          IsDesc: true
                      })
                      break
                  default:
                    break
                  }
                let r = {...props.searchRequest}
                r.Sort = s
                props.onChange(r)
              }}
              value={sort}
            >
              <MenuItem value={-1}>
                {locale2.ACCURACY[props.lang]}
              </MenuItem>
              <MenuItem value={0}>{locale2.PRICE_ASC[props.lang]}</MenuItem>
              <MenuItem value={1}>{locale2.PRICE_DESC[props.lang]}</MenuItem>
              <MenuItem value={2}>{locale2.MAX_PEOPLE_DESC[props.lang]}</MenuItem>
            </Select>
          </FormControl>
    </React.Fragment>
      );
}