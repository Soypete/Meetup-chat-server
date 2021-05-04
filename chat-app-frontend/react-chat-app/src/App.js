import "./App.css";
import React, { Component } from "react";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "",
      timeOnLoad: "",
      currentDate: "",
    };
  }
  componentDidMount() {
    var d = new Date().toString();
    this.setState({
      currentDate: d,
    });
  }

  updateTime = () => {
    this.setState({
      timeOnLoad: new Date(),
    });
  };

  handleClick = () => {
    alert(this.state.name);
  };

  handleChange(event) {
    this.setState({ name: event.target.value });
  }

  handleSubmit(event) {
    alert("A name was submitted: " + this.state.name);
    event.preventDefault();
  }

  render() {
    return (
      <div className="App">
        <div className="chat_box_main">
          <form onSubmit={() => this.handleSubmit()}>
            <label>
              Name:
              <input
                type="text"
                value={this.state.name}
                onChange={(event) => this.handleChange(event)}
              />
            </label>
            <input type="submit" value="Submit" />
          </form>
          <button onClick={() => this.handleClick()}>
            YOU KNOW WHAT IT IS
          </button>
          <div>{this.state.currentDate} - Last load timestamp</div>
          <div>{this.state.name} - This is the current name in state</div>
        </div>
      </div>
    );
  }
}

export default App;
