import styled from 'styled-components';
import cookie from 'react-cookies';
import axios from 'axios';
import './smallBox.css';

function postdel(jsondata) {
    var pass = prompt('비밀번호를 입력해주세요.')
    if(pass === "3695") {
        axios.post(/*"http://127.0.0.1:36530/api/deletepost?postid="*/"https://anend.site:53373/api/deletepost?postid=" + jsondata.id).then( (res) => {
            window.location.href = res.data.redirect;
        }
        );
    }else{
        alert('비밀번호가 다릅니다.');
    }
}

const NullImage = styled.div`
    display: inline-block;
    width: 60px;
    height: 60px;
    border-style: solid;
    border-width: 1px;
    margin: 5px;
    background-color: grey;
`;

function Searchform(result){
    if(cookie.load("anendblogAdminauthcookied") === "yes") {
        return(
            <div className = "searchform">
            <a className = "openpost" href = {"/pages?id="+result.id}>
            <div className = "searchbox">
                <div className = "searchleft">
                    <p className = "searchID">{result.id}</p>
                    {result.thumbURL === "" ? <NullImage/> : <img className = "searchthumbimgset" src={result.thumbURL} alt=''/>}
                </div>
                    <h3 className = "searchtitle">{result.title}</h3>
               <div className = "searchright">
                    <p>조회수: {result.views}</p>
                    <p className = "searchdate">{result.postuploaddate}</p>
                </div>
           </div>
           </a>
           <button className = "searchdelbutton" onClick={() => {postdel(result)}}>삭제</button>
         </div>
        )
    }else {
        return(
            <div className = "searchform">
            <a className = "openpost" href = {"/pages?id="+result.id}>
            <div className = "searchbox">
                <div className = "searchleft">
                    <p className = "searchID">{result.id}</p>
                    {result.thumbURL === "" ? <NullImage/> : <img className = "searchthumbimgset" src={result.thumbURL} alt=''/>}
                </div>
                    <h3 className = "searchtitle">{result.title}</h3>
               <div className = "searchright">
                    <p>조회수: {result.views}</p>
                    <p className = "searchdate">{result.postuploaddate}</p>
                </div>
           </div>
           </a>
         </div>
        )
    }
}

export default Searchform;