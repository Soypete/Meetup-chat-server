import "./App.css";
import React, { Component } from "react";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "",
    };
  }

  handleClick = (word) => {
    alert(word);
  };

  render() {
    return (
      <div className="App">
        <div className="chat_box">
          <button onClick={() => this.handleClick("Drew")}>
            YOU KNOW WHAT IT IS
          </button>
        </div>
        <div className="chat_box">
          <button onClick={() => this.handleClick("Drew")}>
            YOU KNOW WHAT IT IS
          </button>
        </div>
      </div>
    );
  }
}

export default App;
