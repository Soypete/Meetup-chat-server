import "../../App.css";
import React, { Component } from "react";
import axios from "axios";

class Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      email: "",
      password: "",
      rememberBool: false,
      username: "",
    };
  }

  handleClick = () => {
    if (this.state.name === "Drew") {
      alert("Yesssss of course");
      return;
    }
    return alert("Please enter a valid name into the form first!");
  };

  //This is a test for now. testing to see if the database is connected and if it will write an entry to the accounts table
  handleSubmit = async () => {
    try {
      let testBody = {
        username: "CHA BOI",
        password: "Fart Face",
        email: "hugetoots5000@gmail.com",
      };
      await axios.post("localhost:8000/signup", testBody).then((payload) => {
        this.setState({
          username: payload.username,
        });
      });
      return alert("successful creation!");
    } catch (error) {
      alert(`Here is the error` + error);
    }
  };

  render() {
    return (
      <div className="login_container">
        <h1 style={{ textAlign: "center" }}>Miriah's Chat Server Login</h1>
        <div className="login_form_fields">
          <form style={{ padding: "5px" }}>
            <label>
              Email:
              <input
                type="text"
                onChange={(event) =>
                  this.setState({ email: event.target.value })
                }
              />
            </label>
          </form>
          <form>
            <label>
              Password:
              <input
                type="text"
                onChange={(event) =>
                  this.setState({ password: event.target.value })
                }
              />
            </label>
          </form>
        </div>

        <button onClick={() => this.handleSubmit()}>Login</button>
      </div>
    );
  }
}

export default Login;
