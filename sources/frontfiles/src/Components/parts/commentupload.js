import axios from 'axios';
import React from 'react';
import './commentupload.css';

function Commentprocess(id) {
    const email = document.getElementById('uploademail').value;
    if(!email.includes('@') || !email.includes('.')) {
        alert('이메일 형식이 아닙니다');
        return;
    }
    const comment = document.getElementById('uploadcomment').value;
    const secret = document.getElementById('commentsecret').checked;
    
    axios.post(/*"http://127.0.0.1:36530/api/comments"*/"https://anend.site:53373/api/comments",{
        postid: id,
        email: email,
        comment: comment,
        secret: secret,
    }).then( () => {
            window.location.reload()
        }
    )
}

function Commupload(id) {
    return(
        <div className="uploadcomment">
            <input type="text" id="uploademail" placeholder="이메일" maxLength="100"/>
            <textarea type="text" id="uploadcomment" placeholder="내용" maxLength="2000"/>
            <div className="commentsecretform">
                <p className="comminline">비밀글: </p>
                <input type="checkbox" id="commentsecret"/>
            </div>
            <button className="commentsubmit" onClick={() => {Commentprocess(id)}} >제출</button>
        </div>
    );
};

export default Commupload;