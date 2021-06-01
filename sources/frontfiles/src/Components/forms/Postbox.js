import React from 'react';
import cookie from 'react-cookies';
import styled from 'styled-components';
import axios from 'axios';
import './Postbox.css';

const NullImage = styled.div`
    display: inline-block;
    width: 140px;
    height: 140px;
    border-style: solid;
    border-width: 1px;
    margin: 5px;
    background-color: grey;
`;

function deletepost(jsondata) {
    var pass = prompt('비밀번호를 입력해주세요.')
    if(pass === "3695") {
        axios.post(/*"http://127.0.0.1:36530/api/deletepost?postid="*/"https://anend.site:53373/api/deletepost?postid=" + jsondata.id).then( (res) => {
            window.location.href = res.data.redirect;
        });
    }else{
        alert('비밀번호가 다릅니다.');
    }
}



function Postbox(jsondata) {
    console.log(jsondata)
    if(cookie.load("anendblogAdminauthcookied") === "yes") {
        if (jsondata.secret === false) {
            return(
                <div className="postform">
                <a className = "openpost" href = {"/pages?id="+jsondata.id}>
                {jsondata.thumbURL === "" ? <NullImage/> : <img className = "thumbimgset" src={jsondata.thumbURL} alt=''/>}
                <div className = "postinfobox">
                    <h2 className = "posttitle">{jsondata.title}</h2>
                    <p className = "postbody">{jsondata.body}...</p>
                    <p className = "postdate">{jsondata.postuploaddate}</p>
                </div>
                </a>
                <button className = "delpost" onClick={() => {deletepost(jsondata)}}>삭제</button>
            </div>
            );
        }else {
            return(
                <div className="postform">
                <a className = "openpost" href = {"/pages?id="+jsondata.id}>
                {jsondata.thumbURL === "" ? <NullImage/> : <img className = "thumbimgset" src={jsondata.thumbURL} alt=''/>}
                <div className = "postinfobox">
                    <h2 className = "posttitle">{jsondata.title}</h2>
                    <p className = "postbody">...</p>
                    <p className = "postdate">{jsondata.postuploaddate}</p>
                    <p className = "postbemil">비밀글</p>
                </div>
                </a>
                <button className = "delpost" onClick={() => {deletepost(jsondata)}}>삭제</button>
            </div>
            );
        }
    }else {
        if (jsondata.secret === false) {
            return(
                <div className="postform">
                <a className = "openpost" href = {"/pages?id="+jsondata.id}>
                {jsondata.thumbURL === "" ? <NullImage/> : <img className = "thumbimgset" src={jsondata.thumbURL} alt=''/>}
                <div className = "postinfobox">
                    <h2 className = "posttitle">{jsondata.title}</h2>
                    <p className = "postbody">{jsondata.body}...</p>
                    <p className = "postdate">{jsondata.postuploaddate}</p>
                </div>
                </a>
            </div>
            );
        }else {
            return(
                <div className="postform">
                <a className = "openpost" href = {"/pages?id="+jsondata.id}>
                {jsondata.thumbURL === "" ? <NullImage/> : <img className = "thumbimgset" src={jsondata.thumbURL} alt=''/>}
                <div className = "postinfobox">
                    <h2 className = "posttitle">{jsondata.title}</h2>
                    <p className = "postbody">...</p>
                    <p className = "postdate">{jsondata.postuploaddate}</p>
                    <p className = "postbemil">비밀글</p>
                </div>
                </a>
            </div>
            );
        }
    }
}

export default Postbox;