import { AppBar, Button, Container, Toolbar, Typography } from "@mui/material";
import { Box } from "@mui/system";
import React from "react";
import { Link } from "react-router-dom";

class TopBar extends React.Component {
  render() {
    return (
      <AppBar position="static">
        <Container maxWidth="xl">
          <Toolbar disableGutters>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              HMemo
            </Typography>
            <Box>
              <Link to="/add"><Button variant="contained">Add</Button></Link>
            </Box>
          </Toolbar>
        </Container>
      </AppBar>
    );
  }
}

export default TopBar;  