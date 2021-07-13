import "./App.css";
import Login from "./components/stateful/loginpage";
import React, { Component } from "react";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      userEmail: "",
      userName: "",
    };
  }

  render() {
    return (
      <div className="App">
        <Login />
      </div>
    );
  }
}

export default App;
