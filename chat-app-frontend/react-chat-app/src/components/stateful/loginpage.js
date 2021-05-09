import "../../App.css";
import React, { Component } from "react";

class Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      email: "",
      password: "",
      rememberBool: false,
    };
  }

  handleClick = () => {
    if (this.state.name === "Drew") {
      alert("Yesssss of course");
      return;
    }
    return alert("Please enter a valid name into the form first!");
  };

  handleSubmit() {
    alert(
      "Here are the credentials: " +
        this.state.email +
        " " +
        this.state.password
    );
  }

  render() {
    return (
      <div className="login_container">
        <h1>Miriah's Chat Server Login, THOU FIENDS</h1>
        <form>
          <label>
            Email:
            <input
              type="text"
              onChange={(event) => this.setState({ email: event.target.value })}
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

        <button onClick={() => this.handleSubmit()}>Login</button>
      </div>
    );
  }
}

export default Login;
