import React from 'react';
import './App.css';

class ImageInput extends React.Component{
    render(){
        return(
            <div key={"input" + this.props.id} id={"input" + this.props.id} ref={this.props.addChildRef}>
                Select a file: <input className="ui input" type="file" name="file" id="file"/>
                Set frame delay: <input className="ui input" style={{width:"45px"}} type="number" name="delay" id="delay"/>
            </div>
        );
    }
}

class App extends React.Component{
    constructor(props){
        super(props);
        this.state = {
            count: 0,
            inputRefs: [],
            inputComponents: [],
            imageLoaded: false
        }
    }

    addInputRef = (instance) => {
        this.setState({
            inputRefs: this.state.inputRefs.concat(instance)
        })
    }

    toPngBlob = (str) => {
        var hexStr = str.slice(2);
        var buf = new ArrayBuffer(hexStr.length/2);
        var byteBuf = new Uint8Array(buf);
        for (let i=0; i<hexStr.length; i+=2) {
          byteBuf[i/2] = parseInt(hexStr.slice(i,i+2),16);
        }
        var blob = new Blob([byteBuf], {type: "image/png"});
        return blob;
    };

    upload = async () => {
        const formData = new FormData();

        for(let i = 0; i < this.state.inputRefs.length;i++){
            let inputRef = this.state.inputRefs[i];
            formData.append(inputRef.getAttribute("id"), inputRef.querySelector("#file").files[0]);
            formData.append(inputRef.getAttribute("id"), inputRef.querySelector("#delay").value);
        }
        
        /*let self = this;

        fetch('http://127.0.0.1:8090/upload', {
              method: 'POST',
              body: formData
        })
        .then(function(res) {
            return res.blob();
        })
        .then(function(blob){
            let urlCreator = window.URL || window.webkitURL;
            let imageURL = urlCreator.createObjectURL(blob);
            self.setState({
                imageLoaded: true
            });
            document.querySelector("#result").src = imageURL;
        })
        .catch(function(e) {
              console.log('Error', e);
        });
        */

        try{
            let res = await fetch('http://127.0.0.1:8090/upload', {
                method: 'POST',
                body: formData
            });
            let resp = await res.json();
            let urlCreator = window.URL || window.webkitURL;
            let imageURL = urlCreator.createObjectURL(this.toPngBlob(resp.Image));
            this.setState({
                imageLoaded: true
            });
            document.querySelector("#result").src = imageURL;
        } catch(e){
            console.log('Error', e);
        }
    }

    addImage = () => {
        this.setState({
            count: this.state.count + 1,
            inputComponents: this.state.inputComponents.concat(<ImageInput key={this.state.count + 1} id={this.state.count + 1} addChildRef={this.addInputRef}/>)
        });
    }

    render(){
        return(
            <div>
                <div className="ui secondary pointing menu">
                    <div className="item">
                        <h2 className="ui header">
                            <div className="content">
                                Web Animated PNG2
                            </div>
                        </h2>
                    </div>
                </div> 
                <div className="ui page grid">
                    <div className="row">
                            <h2 className="ui header">
                                <div className="content">
                                    Animated PNG Encoder App
                                    <div className="sub header">
                                        Click Add Image to add frames, click Upload to generate the apng2
                                    </div>
                                </div>
                            </h2>
                        </div>
                    <div className="row">
                        <div className="ui buttons">
                            <button className="ui button primary" onClick={this.addImage}>Add Image2</button>
                            <button className="ui button secondary" onClick={this.upload}>Upload2</button>
                        </div>
                    </div>
                    <div className="row">
                        <div className="ui celled list">
                            {this.state.inputComponents.map((comp)=>{
                                return comp;
                            })}
                        </div>
                    </div>
                    {this.state.imageLoaded
                    ?   <div style={{marginTop:'30px'}}>
                            <h2 className="ui header">
                                <div className="content">
                                    Generated Image
                                    <div className="sub header">
                                        Right click and choose Save As to save the full resolution
                                    </div>
                                </div>
                            </h2>
                            <div className="img-container">
                                <img id="result" src="" className="img-input"/>
                            </div>
                        </div>
                    : null
                    }
                </div>
            </div>
        );
    }
}

export default App;