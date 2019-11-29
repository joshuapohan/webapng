import React from 'react';
import './App.css';

class ImageInput extends React.Component{
    constructor(props){
        super(props);
        this.state = {
            files: [],
            delay: this.props.image.delay
        }
    }

    onFileChange = (e) => {
        this.setState({
            files: e.target.files
        });
        this.props.onChildFileChange(this.props.image.id, e.target.files);
    }

    onDelayChange = (e) => {
        this.setState({
            delay: e.target.value
        });
        this.props.onChildDelayChange(this.props.image.id, e.target.value);
    }

    deleteSelf = (e) => {
        e.preventDefault();
        this.props.deleteChild(this.props.image.id);
    }

    render(){
        return(
            <div key={"input" + this.props.image.id} id={"input" + this.props.image.id}>
                Select a file: <input className="ui input" 
                                    type="file" 
                                    name="file" 
                                    id="file"
                                    onChange={this.onFileChange}
                                />
                Set frame delay: <input 
                                    className="ui input" 
                                    style={{width:"45px"}} 
                                    type="number" 
                                    name="delay" 
                                    id="delay" 
                                    value={this.state.delay} 
                                    onChange={this.onDelayChange}
                                /> ms
                <button className="ui button red" style={{textAlign:"center", marginLeft:"10px", marginBottom:"10px"}} onClick={this.deleteSelf}>Remove</button>
            </div>
        );
    }
}

class App extends React.Component{
    constructor(props){
        super(props);
        this.state = {
            fileInput: [],
            imageLoaded: false
        }
    }

    upload = () => {
        const formData = new FormData();

        for(let i = 0; i < this.state.fileInput.length;i++){
            let imageFrame = this.state.fileInput[i];

            formData.append(i, imageFrame.files[0]);
            formData.append(i, parseInt(imageFrame.delay));
        }
        
        let self = this;

        fetch('/upload', {
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
    }

    addImage = () => {
        this.setState({
            fileInput: [...this.state.fileInput, {id: this.state.fileInput.length, files:[], delay:64}]
        });
    }

    deleteChild = (id) => {
        let inputList = [...this.state.fileInput].filter(input=>{
            return input.id !== id;
        });
        this.setState({
            fileInput: inputList
        });      
    }

    onChildFileChange = (id, files) => {
        let inputList = [...this.state.fileInput];
        let inputFile = {...inputList[id]};
        inputFile.files = files;
        inputList[id] = inputFile;
        this.setState({
            fileInput: inputList
        });      
    }

    onChildDelayChange = (id, delay) => {
        let inputList = [...this.state.fileInput];
        let inputFile = {...inputList[id]};
        inputFile.delay = delay;
        inputList[id] = inputFile;
        this.setState({
            fileInput: inputList
        });
    }

    render(){
        return(
            <div>
                <div className="ui secondary pointing menu">
                    <div className="item">
                        <h2 className="ui header">
                            <div className="content">
                                Web Animated PNG
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
                                        Click Add Image to add frames, click Upload to generate the apng
                                    </div>
                                </div>
                            </h2>
                        </div>
                    <div className="row">
                        <div className="ui buttons">
                            <button className="ui button primary" onClick={this.addImage}>Add Image</button>
                            <button className="ui button secondary" onClick={this.upload}>Upload</button>
                        </div>
                    </div>
                    <div className="row">
                        <div className="ui celled list">
                            {this.state.fileInput.map((comp, index)=>{
                                return <ImageInput 
                                            key={index} 
                                            image={{...comp}} 
                                            deleteChild={this.deleteChild}
                                            onChildFileChange={this.onChildFileChange} 
                                            onChildDelayChange={this.onChildDelayChange}
                                        />;
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
