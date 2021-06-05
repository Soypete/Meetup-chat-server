const express = require("express");
const app = express();
const port = 8000;
const bodyParser = require("body-parser");
const session = require("express-session");
const massive = require("massive");
const path = require("path");

const db_path = "postgresql://localhost:6000/postgres";

massive(db_path).then((db) => {
  console.log("PostgreSQL Database Successfully Connected");
  app.set("db", db);
});

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`);
});

//I never really understood express session..I should get help on figuring out why this is needed, or if it's even needed at all.
app.use(
  session({
    secret: "keyboard cat",
    maxAge: 86400000,
    resave: true,
    saveUninitialized: true,
  })
);

//Hello world example
app.get("/", (req, res) => {
  res.send("Hello World!");
});

//Basic account creation example - need to pass a body from the frontend
app.post("/signup", async (req, res) => {
  try {
    const newUser = await db.accounts.insert({
      username: req.body.username,
      password: req.body.password,
      email: req.body.email,
      created_on: `NOW()`,
      last_login: `NOW()`,
    });

    delete newUser.password;
    res.send(newUser);
  } catch (error) {
    console.log(error);
    res.status(500).send(error);
  }
});
