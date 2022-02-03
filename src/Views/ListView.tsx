import { Button, Card, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import axios from "axios";
import { useEffect, useState } from "react";
import DeleteIcon from '@mui/icons-material/Delete';
import { Box } from "@mui/system";
import { Link } from "react-router-dom";

const ListView = () => {
  const [memos, setMemos] = useState([])
  const showList = () => {
    axios.get("./api/memos").then(r => {
      if (r.data !== undefined) {
        setMemos(r.data)
      }
    })
  }

  useEffect(showList, [])

  const deleteMemo = (id: string) => { axios.delete(`./api/memos/${id}`).then(r => showList()) }
  const editMemo = (id: string) => { window.location.href = `#/memos/${id}` }

  const createListItems = (i: { id: string, name: string, shop: string, date: string }) => {
    return <TableRow hover key={i.id}>
      <TableCell component="th" onClick={() => editMemo(i.id)} style={{ cursor: 'pointer' }}>{i.date}</TableCell>
      <TableCell onClick={() => editMemo(i.id)} style={{ cursor: 'pointer' }}>{i.name}</TableCell>
      <TableCell onClick={() => editMemo(i.id)} style={{ cursor: 'pointer' }}>{i.shop}</TableCell>
      <TableCell><Button onClick={deleteMemo.bind(this, i.id)}><DeleteIcon /></Button></TableCell>
    </TableRow>
  }

  return <Box m={1}>
    <Card>
      <TableContainer>
        <Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
          <TableHead>
            <TableRow>
              <TableCell>date</TableCell>
              <TableCell>name</TableCell>
              <TableCell>shop</TableCell>
              <TableCell></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {memos.map((v: { id: string, name: string, shop: string, date: string }) => createListItems(v))}
          </TableBody>
        </Table>
      </TableContainer>
    </Card>
  </Box>
}

export default ListView