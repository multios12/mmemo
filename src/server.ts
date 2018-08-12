#!/usr/bin/env node
import http from "http";

// tslint:disable-next-line:no-var-requires
const app = require("./app");
// tslint:disable-next-line:no-var-requires
const debug = require("debug")("bootstraptemplate:server");

/** Get port from environment and store in Express. */
const port = normalizePort(process.env.PORT || "3000");
app.set("port", port);

/** HTTPサーバ */
const server = http.createServer(app);

// Listen on provided port, on all network interfaces.
server.listen(port);
server.on("error", onError);
server.on("listening", onListening);

/** Normalize a port into a number, string, or false. */
function normalizePort(val: string) {
  const portNumber = parseInt(val, 10);

  if (isNaN(portNumber)) {
    // named pipe
    return val;
  }

  if (portNumber >= 0) {
    // port number
    return portNumber;
  }

  return false;
}

/** Event listener for HTTP server "error" event. */
function onError(error: any) {
  if (error.syscall !== "listen") {
    throw error;
  }

  const bind = typeof port === "string"
    ? "Pipe " + port
    : "Port " + port;

  // handle specific listen errors with friendly messages
  switch (error.code) {
    case "EACCES":
      console.error(bind + " requires elevated privileges");
      process.exit(1);
      break;
    case "EADDRINUSE":
      console.error(bind + " is already in use");
      process.exit(1);
      break;
    default:
      throw error;
  }
}

/** Event listener for HTTP server "listening" event. */
function onListening() {
  const addr = server.address();
  const bind = typeof addr === "string"
    ? "pipe " + addr
    : "port " + addr.port;
  debug("Listening on " + bind);
}
