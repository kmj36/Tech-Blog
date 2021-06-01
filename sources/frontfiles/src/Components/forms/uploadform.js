import axios from 'axios';
import React from 'react';
import editor from '../parts/editor';
import cookie from 'react-cookies';
import './uploadform.css';

function process() {
    const title = document.getElementsByClassName("title_box")[0].value;
    if (title === '')
        return alert("제목을 입력해주세요.");
    const content = document.getElementById("editor").innerHTML;
    const secret = document.getElementsByName("Secret")[0].checked;
    const thumbnail = document.getElementsByName("urlinput")[0].value;
    axios.post(/*"http://127.0.0.1:36530/api/upload"*/"https://anend.site:53373/api/upload",{
        "title": title,
        "content": content,
        "secret": secret,
        "thumbURL": thumbnail
    }).then((res) => {
        window.location.href = res.data.redirect
    })
}

function checkadmin() {
    const id = document.getElementById('ID').value;
    const pass = document.getElementById('PASS').value;
    if(id === "kmj36953695" && pass === "Aho@0406*0013-") {
        cookie.save("anendblogAdminauthcookied","yes");
        window.location.reload();
    }else {
        alert("입력를 다시 확인해주십시오.")
    }
}

function editpost(id) {
    const title = document.getElementsByClassName("title_box")[0].value;
    if (title === '')
        return alert("제목을 입력해주세요.");
    const content = document.getElementById("editor").innerHTML;
    const secret = document.getElementsByName("Secret")[0].checked;
    const thumbnail = document.getElementsByName("urlinput")[0].value;
    axios.post(/*"http://127.0.0.1:36530/api/editpost"*/"https://anend.site:53373/api/editpost",{
        "id": id,
        "title": title,
        "content": content,
        "secret": secret,
        "thumbURL": thumbnail
    }).then((res) => {
        window.location.href = res.data.redirect
    })
}

function uploader(load) {
    if(cookie.load("anendblogAdminauthcookied") !== "yes") {
        return(
            document.title = "Please Auth",
            <div className="logincontainer">
                <div className="loginblock">
                    <h1 className="setcenter">Login</h1>
                    <div className="logininput">
                            <div className="loginid">
                                 <h4 className="logincenter">ID</h4>
                               <input id = "ID" name = "inputid" type="text" minLength="50"/>
                            </div>
                            <div className="loginpass">
                                <h4 className="logincenter">PASS</h4>
                                 <input id = "PASS" name = "inputpass" type="password" minLength="50"/>
                            </div>
                        </div>
                        <button className="loginsubmitbutton" name="submit" onClick={() => checkadmin()}>제출</button>
                    </div>
                </div>
        );
    }else {
        if(load === undefined) {
            return(
                <div className = "uploadformwrapper">
                    <div className = "title_con"><p className = "inline">제목:  </p><input className = "title_box" type = "text" name = "posttitle" maxLength="200"/></div>
                    {editor()}
                    <div className = "bottomwrapper">
                    <input type = "button" name = "Submit" value = "업로드" onClick={() => process()}/>
                    <div className = "thumbnail">
                    <p className="inline">썸네일 URL: </p>
                    <input type = "text" className="url_box" name = "urlinput" maxLength="500"/>
                    </div>
                    <div className = "right"><p className = "inline">비밀글:</p><input type = "checkbox" name = "Secret" value = "yes"/></div>
                    </div>
                </div>
            );
        }else {
            return(
                <div className = "uploadformwrapper">
                    <div className = "title_con"><p className = "inline">제목:  </p><input className = "title_box" type = "text" name = "posttitle" id = "posttitle" maxLength="200"/></div>
                    {editor()}
                    <div className = "bottomwrapper">
                    <input type = "button" name = "Submit" value = "수정" onClick={() => editpost(load)}/>
                    <div className = "thumbnail">
                    <p className="inline">썸네일 URL: </p>
                    <input type = "text" id="urlinput" className="url_box" name = "urlinput" maxLength="500"/>
                    </div>
                    <div className = "right"><p className = "inline">비밀글:</p><input type = "checkbox" name = "Secret" value = "yes"/></div>
                    </div>
                    <input hidden onLoad = {
                       axios.get(/*"http://127.0.0.1:36530/api/posts?id="*/"https://anend.site:53373/api/posts?id=" + load).then((data) => {
                            const postdatas = data.data[0]
                            console.log(postdatas)
                            document.getElementById("posttitle").value = postdatas.title
                            document.getElementById("urlinput").value = postdatas.thumbURL
                            document.getElementById("editor").innerHTML = postdatas.content
                            document.getElementsByName("Secret")[0].checked = postdatas.secret
                        })
                    }/>
                </div>
            );
        }
    }
}

export default uploader;