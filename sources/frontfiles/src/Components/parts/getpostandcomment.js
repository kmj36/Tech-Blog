import React, { useEffect, useState } from 'react';
import cookie from 'react-cookies';
import axios from 'axios';
import './getpostandcomment.css';
import commentupload from './commentupload';

function postdel(jsondata) {
    var pass = prompt("비밀번호를 입력해주세요.")
    if (pass === "3695") {
        axios.post(/*"http://127.0.0.1:36530/api/deletepost?postid="*/"https://anend.site:53373/api/deletepost?postid=" + jsondata.id).then((res) => {
            window.location.href = res.data.redirect;
        }
        )
    } else {
        alert("비밀번호가 다릅니다.")
    }
}

function commentdel(jsondata) {
    var pass = prompt("비밀번호를 입력해주세요.")
    if (pass === "3695") {
        axios.post(/*"http://127.0.0.1:36530/api/deletecomment?comtid="*/"https://anend.site:53373/api/deletecomment?comtid=" + jsondata.comtid).then(() => {
            window.location.reload();
        }
        )
    } else {
        alert("비밀번호가 다릅니다.")
    }
}

function Getthepost(id) {
    const [post, setpost] = useState();
    const [comment, setcomment] = useState();
    const [adminauth, setadminauth] = useState();

    const getpost = async () => {
        const data = await axios.get(/*"http://127.0.0.1:36530/api/posts?id="*/"https://anend.site:53373/api/posts?id=" + id);
        setpost(data.data);
    }

    const getcomment = async () => {
        const data = await axios.get(/*"http://127.0.0.1:36530/api/comments?id="*/"https://anend.site:53373/api/comments?id=" + id);
        setcomment(data.data);
    }

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

    useEffect(() => { getadminauth(); getpost(); getcomment(); }, []);
    return (post?.map((postdata) => {
        if (adminauth) {
            document.title = postdata.title;
            return (
                <div className="postwrapper">
                    <div className="posttop">
                        <a href="/" className="postbackbutton">&lt;돌아가기</a>
                        <p className="postopentitle">{postdata.title}</p>
                        <div className="viewsanddate">
                            <p className="postviews">조회수: {postdata.views}</p>
                            <div className="controltopbutton">
                                <button onClick={() => { window.location.href = "/upload?load=" + postdata.id }}>수정</button>
                                <button onClick={() => { postdel(postdata) }}>삭제</button>
                            </div>
                            <p className="postuploaddate">{postdata.postuploaddate}</p>
                        </div>
                    </div>
                    <div className="postmiddle" dangerouslySetInnerHTML={{ __html: postdata.content }}>
                    </div>
                    <div className="postbottom">
                        {comment?.map((commentdata) => {
                            return (
                                <div className="postcommentform">
                                    <div className="commenthead">
                                        <h4 className="commentemail">{commentdata.email}</h4>
                                        <p className="commentdate">{commentdata.date}</p>
                                    </div>
                                    <div className="comment">
                                        <p className="commentbody">{commentdata.comment}</p>
                                    </div>
                                    <button className="controlcombutton" onClick={() => { commentdel(commentdata) }}>삭제</button>
                                </div>
                            );
                        })}
                    </div>
                    {commentupload(id)}
                </div>
            );
        } else {
            if (postdata.secret) {
                return (
                    alert('비밀글 입니다.'),
                    window.location.href = "/"
                );
            } else {
                document.title = postdata.title;
                return (
                    <div className="postwrapper">
                        <div className="posttop">
                            <a href="/" className="postbackbutton">&lt;돌아가기</a>
                            <p className="postopentitle">{postdata.title}</p>
                            <div className="viewsanddate">
                                <p className="postviews">조회수: {postdata.views}</p>
                                <p className="postuploaddate">{postdata.postuploaddate}</p>
                            </div>
                        </div>
                        <div className="postmiddle" dangerouslySetInnerHTML={{ __html: postdata.content }}>
                        </div>
                        <div className="postbottom">
                            {comment?.map((commentdata) => {
                                if (commentdata.secret === false) {
                                    return (
                                        <div className="postcommentform">
                                            <div className="commenthead">
                                                <h4 className="commentemail">{commentdata.email}</h4>
                                                <p className="commentdate">{commentdata.date}</p>
                                            </div>
                                            <div className="comment">
                                                <p className="commentbody">{commentdata.comment}</p>
                                            </div>
                                        </div>
                                    );
                                } else {
                                    return (
                                        <div className="postcommentform">
                                            <p>비밀글 입니다.</p>
                                        </div>
                                    );
                                }

                            })}
                        </div>
                        {commentupload(id)}
                    </div>
                );
            }
        }
    }));
}

export default Getthepost;