// ExpressJS Setup
const express = require("express");
const app = express();
var bodyParser = require("body-parser");

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "..", "network", "connection.json");
const ccpJSON = fs.readFileSync(ccpPath, "utf8");
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 8080;
const HOST = "localhost";

app.set("views", path.join(__dirname, "views"));
app.set("view engine", "ejs");

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get("/", (req, res) => {
  res.render("index.ejs");
});

async function cc_call(fn_name, args) {
  const walletPath = path.join(process.cwd(), "wallet");
  const wallet = new FileSystemWallet(walletPath);

  const userExists = await wallet.exists("user1");
  if (!userExists) {
    console.log(
      'An identity for the user "user1" does not exist in the wallet'
    );
    console.log("Run the registerUser.js application before retrying");
    return;
  }
  const gateway = new Gateway();
  await gateway.connect(ccp, {
    wallet,
    identity: "user1",
    discovery: { enabled: false },
  });
  const network = await gateway.getNetwork("mychannel");
  const contract = network.getContract("bc_game");

  var result;

  if (fn_name == "addUser")
    result = await contract.submitTransaction("addUser", args);
  else if (fn_name == "readUserInfo")
    result = await contract.evaluateTransaction("readUserInfo", args);
  else if (fn_name == "readUserList")
    result = await contract.evaluateTransaction("readUserList");
  else result = "not supported function";

  return result;
}

// create PN
app.post("/user/create", async (req, res) => {
  const name = req.body.name;
  console.log("add user: " + name);

  result = cc_call("addUser", name);

  res.render("index", {
    user: name,
    result: result,
  });
});

app.post("/user/info", function (req, res) {
  const name = req.body.name;
  result = cc_call("readUserInfo", name);
  res.status(200).send(result);
});

// add score
app.post("/user/info/:name", async (req, res) => {
  const email = req.body.email;
  const prj = req.body.project;
  const sc = req.body.score;
  console.log("add project email: " + email);
  console.log("add project name: " + prj);
  console.log("add project score: " + sc);

  var args = [email, prj, sc];
  result = cc_call("addRating", args);

  const myobj = { result: "success" };
  res.status(200).json(myobj);
});

// find mate
app.post("/donation/:name", async (req, res) => {
  const PNname = req.body.name;
  console.log("name: " + req.body.name);
  const walletPath = path.join(process.cwd(), "wallet");
  const wallet = new FileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);

  // Check to see if we've already enrolled the user.
  const userExists = await wallet.exists("user1");
  if (!userExists) {
    console.log(
      'An identity for the user "user1" does not exist in the wallet'
    );
    console.log("Run the registerUser.js application before retrying");
    return;
  }
  const gateway = new Gateway();
  await gateway.connect(ccp, {
    wallet,
    identity: "user1",
    discovery: { enabled: false },
  });
  const network = await gateway.getNetwork("mychannel");
  const contract = network.getContract("teamate");
  const result = await contract.evaluateTransaction("readRating", email);
  const myobj = JSON.parse(result);
  res.status(200).json(myobj);
});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
