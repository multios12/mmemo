import { Button, Card, CardActions, CardContent, Grid, TextField } from "@mui/material";
import { Box } from "@mui/system";
import axios from "axios";
import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

const EditView = () => {
  const { id } = useParams();
  const [name, setName] = useState("")
  const [shop, setShop] = useState("")
  const [page, setPage] = useState("")
  const [date, setDate] = useState("")
  const [play, setPlay] = useState("")
  const [talk, setTalk] = useState("")
  var errorMessage: string = ""
  const navigate = useNavigate();

  useEffect(() => {
    axios.get(`/api/memos/${id}`).then(r => {
      if (r.data != undefined) {
        setName(r.data[0].name || "");
        setShop(r.data[0].shop || "");
        setPage(r.data[0].Page || "");
        setDate(r.data[0].date || "");
        setPlay(r.data[0].play || "");
        setTalk(r.data[0].talk || "");
      }
    })
  }, [])

  const regist = () => {
    if (id == "") {
      axios.put("../api/memos", { name, shop, page, date, play, talk })
        .catch(res => (errorMessage = "登録が失敗しました。"));
    } else {
      axios.post(`../api/memos/${id}`, { name, shop, page, date, play, talk })
        .catch(res => (errorMessage = "登録が失敗しました。"));
    }
    navigate("/");
  }

  return <Card>
    <CardContent>
      <Grid container spacing={1}>
        <Grid item xs={12}>
          <TextField label="name" variant="filled" onChange={e => setName(e.target.value)} fullWidth value={name} />
        </Grid>
        <Grid item xs={12}>
          <Box display="flex" justifyContent="space-between">
            <TextField label="shop" variant="filled" onChange={e => setShop(e.target.value)} fullWidth value={shop} />
            <TextField label="home page" variant="filled" onChange={e => setPage(e.target.value)} fullWidth value={page} />
          </Box>
        </Grid>
        <Grid item xs={12}>
          date<TextField variant="filled" type="date" onChange={e => setDate(e.target.value)} fullWidth value={date} />
        </Grid>
        <Grid item xs={12}>
          <TextField multiline label="play" variant="filled" onChange={e => setPlay(e.target.value)} fullWidth value={play} />
        </Grid>
        <Grid item xs={12}>
          <TextField multiline label="talk" variant="filled" onChange={e => setTalk(e.target.value)} fullWidth value={talk} />
        </Grid>
      </Grid>
    </CardContent>
    <CardActions>
      <Button onClick={regist}>Ok</Button>
      <Link to="/"><Button>Cancel</Button></Link>
    </CardActions>
  </Card>
}

export default EditView