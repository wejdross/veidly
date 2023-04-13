// import { ButtonBase, Card, 
//         CardHeader, CardMedia, 
//         List, ListItem, 
//         makeStyles, Typography } from "@mui/material";
// import React, { useEffect, useState } from "react";
// import { returnLocaleString } from "../locale";

// const useStyles = makeStyles((theme) => ({
//     media: {
//         height: 0,
//         paddingTop: '56.25%', // 16:9
//     },
//     root: {
//         width: 280,
//     },
// }))

// export const stb = "/static/form-backgrounds/surfer.webp"


// function City(props) {

//     const classes = useStyles()

//     function redir(lat, lng) {
//         try {
//             let query = new URLSearchParams(window.location.search)
//             let r = query.get("q")
//             r = JSON.parse(r)
//             r.Lat = lat
//             r.Lng = lng
//             r.DistKm = props.d
//             r.display_name = props.name
//             props.onChange(r)
//             // query.set("q", JSON.stringify(r));
//             // h.push({
//             //     search: query.toString()
//             //   })
//             //window.location.search = query.toString();
//         } catch (ex) {
//             console.log(ex)
//         }
//     }

//     return (<ButtonBase onClick={() => redir(props.lat, props.lng)}>
//         <Card className={classes.root}>
//             <CardMedia
//                 className={classes.media}
//                 image={props.img} />
//             <CardHeader
//                 title={props.name}
//                 subheader={props.sub} />
//         </Card>
//     </ButtonBase>)
// }

// export default function Cities(props) {

//     const [allowedDistance, setAllowedDistance] = useState(10)

//     useEffect(() => {
//         let query = new URLSearchParams(window.location.search)
//         let r = query.get("q")
//         r = JSON.parse(r)
//         if (r && r.DistKm)
//             setAllowedDistance(r.DistKm)
//     }, [])

//     const cities = [
//         {
//             img: "kato.jpg",
//             name: "Katowice",
//             sub: "Śląskie",
//             lat: 50.16,
//             lng: 19.01
//         },
//         {
//             img: "poz.jpg",
//             name: "Poznań",
//             sub: "Wielkopolska",
//             lat: 52.25,
//             lng: 16.58
//         },
//         {
//             img: "war.jpg",
//             name: "Warszawa",
//             sub: "Mazowsze",
//             lat: 52.22977,
//             lng: 21.0117800
//         },
//         {
//             img: "krak.jpg",
//             name: "Kraków",
//             sub: "Małopolska",
//             lat: 50.0614300,
//             lng: 19.9365800
//         },
//         {
//             img: "lodz.jpg",
//             name: "Łódź",
//             sub: "Łódzkie",
//             lat: 51.7500000,
//             lng: 19.4666700
//         },
//         {
//             img: "wroc.jpg",
//             name: "Wrocław",
//             sub: "Dolny śląsk",
//             lat: 51.10000,
//             lng: 17.03333
//         },
//         {
//             img: "yacht.jpg",
//             name: "Gdańsk",
//             sub: "Pomorze",
//             lat: 54.3520500,
//             lng: 18.6463700
//         },
//         {
//             img: "bydg.jpg",
//             name: "Bydgoszcz",
//             sub: "Kujawsko-pomorskie",
//             lat: 53.1235000,
//             lng: 18.0076200
//         },
//     ]

//     return (<React.Fragment>
//         <center>
//             <Typography variant="h5" style={{
//                 marginTop: 20,
//                 marginBottom: 20,
//             }}>
        /*
        Cities: {
            pl: [
                "Znajdź trenerów w pobliżu większego miasta"
            ],
            en: [
                "Find instructors close to the bigger cities"
            ]
        }
        */
//                 {returnLocaleString(["search", "Cities"])}
//             </Typography>
//         </center>
//             <List style={{
//                 display: 'flex',
//                 flexDirection: 'row',
//                 padding: 0,
//                 overflowX:"auto"
//             }}>
//                 {cities.map((c, i) => (
//                     <ListItem key={i} style={{
//                         width: 320,
//                         marginBottom: 20,
//                         marginRight: 10,
//                         marginLeft: 10,
//                     }}>
//                         <City onChange={props.onChange}
//                             img={stb + c.img}
//                             name={c.name}
//                             sub={c.sub}
//                             d={allowedDistance}
//                             lat={c.lat} lng={c.lng} />
//                     </ListItem>
//                 ))}
//             </List>
//     </React.Fragment>)
// }
