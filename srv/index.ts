import express from "express";
import apiMemos from "./api-memos";
// import socketIO from "socket.io";

export default (app, http) => {
  app.use(express.json());

  app.use((req, res, next) => {
    res.header("Access-Control-Allow-Origin", "*");
    res.header(
      "Access-Control-Allow-Headers",
      "Origin, X-Requested-With, Content-Type, Accept, Authorization",
    );
    res.header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE");

    next();
  });

  app.use("/api/memos/", apiMemos);
};
