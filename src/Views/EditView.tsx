import { Alert, Backdrop, Button, Card, CardActions, CardContent, CircularProgress, Grid, Snackbar, TextField } from "@mui/material";
import { Box } from "@mui/system";
import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

type memoType = {
  id: string | undefined,
  name: string,
  shop: string,
  page: string,
  date: string,
  play: string,
  talk: string,
}

export default function EditView() {
  const m: memoType = { id: undefined, name: "", shop: "", page: "", date: new Date().toISOString().substring(0, 10), play: "", talk: "" }
  const { id } = useParams();
  const [memo, setMemo] = useState(m);
  const [isErr, setIsErr] = useState(false)
  const [errMessage, setErrMessage] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate();

  useEffect(() => {
    if ((id === undefined) || (id === "")) {
      return
    }
    setIsLoading(true)
    axios.get(`../api/memos/${id}`).then(r => {
      if (r.data !== undefined) {
        setMemo(r.data[0]);
      }
    }).finally(() => setIsLoading(false))
  }, [])

  const regist = () => {
    setIsLoading(true)
    if (id === "") {
      axios.put("../api/memos", memo)
        .then(r => navigate("/"))
        .catch(res => {
          setErrMessage(res.response.data.error)
          setIsErr(true)
        }).finally(() => { setIsLoading(false) });
    } else {
      axios.post(`../api/memos/${id}`, memo)
        .then(r => navigate("/"))
        .catch(res => {
          setErrMessage(res.response.data.error)
          setIsErr(true)
        }).finally(() => { setIsLoading(false) });
    }

  }

  return <Card>
    <CardContent>
      <Snackbar open={isErr} autoHideDuration={4000} onClose={() => setIsErr(false)}>
        <Alert onClose={() => setIsErr(false)} severity="error">{errMessage}</Alert>
      </Snackbar>
      <Backdrop open={isLoading} onClick={() => setIsLoading(false)}>
        <CircularProgress color="inherit" />
      </Backdrop>
      <Grid container spacing={1}>
        <Grid item xs={12}>
          <TextField fullWidth label="name" variant="filled" value={memo.name} onChange={e => setMemo({ ...memo, name: e.target.value })}
            error={!memo.name && !isLoading} helperText={!memo.name && "input required"} />
        </Grid>
        <Grid item xs={12}>
          <Box display="flex" justifyContent="space-between">
            <TextField fullWidth sx={{ mr: 1 }} label="shop" variant="filled" value={memo.shop} onChange={e => setMemo({ ...memo, shop: e.target.value })} />
            <TextField fullWidth label="home page" variant="filled" onChange={e => setMemo({ ...memo, page: e.target.value })} value={memo.page} type="url" />
          </Box>
        </Grid>
        <Grid item xs={12}>
          <TextField fullWidth label="date" variant="filled" type="date" value={memo.date} onChange={e => setMemo({ ...memo, date: e.target.value })}
            error={!memo.date && !isLoading} helperText={!memo.date ? "input required" : ""} />
        </Grid>
        <Grid item xs={12}>
          <TextField multiline fullWidth label="play" variant="filled" value={memo.play} onChange={e => setMemo({ ...memo, play: e.target.value })} />
        </Grid>
        <Grid item xs={12}>
          <TextField multiline fullWidth label="talk" variant="filled" value={memo.talk} onChange={e => setMemo({ ...memo, talk: e.target.value })} />
        </Grid>
      </Grid>
    </CardContent>
    <CardActions>
      <Button variant="contained" color="primary" disabled={isLoading} onClick={regist}>Ok</Button>
      <Button variant="contained" color="inherit" onClick={() => navigate("/")}>Cancel</Button>
    </CardActions>
  </Card>
}

