import { Button, List, ListItem } from "@mui/material";
import axios from "axios";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import DeleteIcon from '@mui/icons-material/Delete';

const ListView = () => {
    const [memos, setMemos] = useState([])
    const showList = () => {
        axios.get("/api/memos").then(r => {
            if (r.data != undefined) {
                setMemos(r.data)
            }
        })
    }

    useEffect(showList, [])

    const deleteMemo = (id: string) => { axios.delete(`/api/memos/${id}`).then(r => showList())}

    const createListItems = (i: { id: string, name: string, shop: string, date: string }) => {
        return <ListItem key={i.id}>{i.date} <Link to={`/memos/${i.id}`}>{i.name}</Link> {i.shop}
            <Button onClick={deleteMemo.bind(this, i.id)}><DeleteIcon /></Button>
        </ListItem>
    }
    return <List>{memos.map((v: { id: string, name: string, shop: string, date: string }) => createListItems(v))}</List>
}

export default ListView