import axios from 'axios';
import React, { useEffect, useState } from 'react';
import editor from '../parts/editor';
import cookie from 'react-cookies';
import './uploadform.css';

function uploadpost() {
    const title = document.getElementsByClassName("title_box")[0].value;
    if (title === '')
        return alert("제목을 입력해주세요.");
    const content = document.getElementById("editor").innerHTML;
    const secret = document.getElementsByName("Secret")[0].checked;
    const thumbnail = document.getElementsByName("urlinput")[0].value;
    axios.post(/*"http://127.0.0.1:36530/api/uploadpost"*/"https://anend.site:53373/api/uploadpost", {
        "title": title,
        "content": content,
        "secret": secret,
        "thumbURL": thumbnail
    }).then((res) => {
        window.location.href = res.data.redirect
    })
}

function loginadmin() {
    const id = document.getElementById('ID').value;
    const pass = document.getElementById('PASS').value;

    if (id === '')
        return alert("아이디를 입력해주세요.");
    if (pass === '')
        return alert("패스워드를 입력해주세요.");

    axios.post(/*"http://127.0.0.1:36530/api/adminauth"*/"https://anend.site:53373/api/adminauth", {
        "adminid": id,
        "adminpassword": pass
    }).then((res) => {
        cookie.save("auth", res.data.auth_cookie)
        cookie.save("time", res.data.time)
        window.location.reload()
    })
}

function editpost(id) {
    const title = document.getElementsByClassName("title_box")[0].value;
    if (title === '')
        return alert("제목을 입력해주세요.");
    const content = document.getElementById("editor").innerHTML;
    const secret = document.getElementsByName("Secret")[0].checked;
    const thumbnail = document.getElementsByName("urlinput")[0].value;
    axios.post(/*"http://127.0.0.1:36530/api/editpost"*/"https://anend.site:53373/api/editpost", {
        "id": id,
        "title": title,
        "content": content,
        "secret": secret,
        "thumbURL": thumbnail
    }).then((res) => {
        window.location.href = res.data.redirect
    })
}

function Uploader(load) {
    const [adminauth, setadminauth] = useState();

    const getadminauth = async () => {
        const auth = cookie.load("auth")
        const time = cookie.load("time")
        await axios.post(/*"http://127.0.0.1:36530/api/checkadmin"*/"https://anend.site:53373/api/checkadmin", {
            "auth_cookie": auth,
            "time" : time
        }).then((res) => {
            setadminauth(res.data.auth === "Yes");
        })
    }

    useEffect(() => {
        getadminauth();
        document.execCommand('styleWithCSS', false, true);
        document.execCommand('enableObjectResizing', false, true);
    }, []);

    if (adminauth) {
        if (load === undefined) {
            return (
                <div className="uploadformwrapper">
                    <div className="title_con"><p className="inline">제목:  </p><input className="title_box" type="text" name="posttitle" maxLength="200" /></div>
                    {editor()}
                    <div className="bottomwrapper">
                        <input type="button" name="Submit" value="업로드" onClick={() => uploadpost()} />
                        <div className="thumbnail">
                            <p className="inline">썸네일 URL: </p>
                            <input type="text" className="url_box" name="urlinput" maxLength="500" />
                        </div>
                        <div className="right"><p className="inline">비밀글:</p><input type="checkbox" name="Secret" value="yes" /></div>
                    </div>
                </div>
            );
        } else {
            return (
                <div className="uploadformwrapper">
                    <div className="title_con"><p className="inline">제목:  </p><input className="title_box" type="text" name="posttitle" id="posttitle" maxLength="200" /></div>
                    {editor()}
                    <div className="bottomwrapper">
                        <input type="button" name="Submit" value="수정" onClick={() => editpost(load)} />
                        <div className="thumbnail">
                            <p className="inline">썸네일 URL: </p>
                            <input type="text" id="urlinput" className="url_box" name="urlinput" maxLength="500" />
                        </div>
                        <div className="right"><p className="inline">비밀글:</p><input type="checkbox" name="Secret" value="yes" /></div>
                    </div>
                    <input hidden onLoad={
                        axios.get(/*"http://127.0.0.1:36530/api/posts?id="*/"https://anend.site:53373/api/posts?id=" + load).then((data) => {
                            const postdatas = data.data[0]
                            document.getElementById("posttitle").value = postdatas.title
                            document.getElementById("urlinput").value = postdatas.thumbURL
                            document.getElementById("editor").innerHTML = postdatas.content
                            document.getElementsByName("Secret")[0].checked = postdatas.secret
                        })
                    } />
                </div>
            );
        }
    } else {
        return (
            document.title = "Please Auth",
            <div className="logincontainer">
                <div className="loginblock">
                    <h1 className="setcenter">Login</h1>
                    <div className="logininput">
                        <div className="loginid">
                            <h4 className="logincenter">ID</h4>
                            <input id="ID" name="inputid" type="text" minLength="50" />
                        </div>
                        <div className="loginpass">
                            <h4 className="logincenter">PASS</h4>
                            <input id="PASS" name="inputpass" type="password" minLength="50" />
                        </div>
                    </div>
                    <button className="loginsubmitbutton" name="submit" onClick={() => loginadmin()}>제출</button>
                </div>
            </div>
        );
    }
}

export default Uploader;